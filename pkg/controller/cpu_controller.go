package controller

import (
	"fmt"
	"log"
	"time"

	api_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	clientset "github.com/interma/programming-k8s/pkg/client/clientset/versioned"
)

const maxRetries = 5
const namespace = "default"

// Controller object
type StatsController struct {
	KubeClient kubernetes.Interface
	CrClient   clientset.Interface
	CpuCrName  string

	queue    workqueue.RateLimitingInterface
	informer cache.SharedIndexInformer
}

func CreatePodsStatsController(kc kubernetes.Interface, cc clientset.Interface, crName string) *StatsController {
	// pod event informer
	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				return kc.CoreV1().Pods(namespace).List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				return kc.CoreV1().Pods(namespace).Watch(options)
			},
		},
		&api_v1.Pod{},
		0, //Skip resync
		cache.Indexers{},
	)

	c := newPodsStatsController(kc, cc, informer, crName)

	return c
}

func newPodsStatsController(kc kubernetes.Interface, cc clientset.Interface,
	informer cache.SharedIndexInformer, crName string) *StatsController {

	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// add handler to informer
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			log.Printf("processing pod add: %s", key)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(old, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			log.Printf("processing pod update: %s", key)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			log.Printf("processing pod delete: %s", key)
			if err == nil {
				queue.Add(key)
			}
		},
	})

	return &StatsController{kc, cc, crName, queue, informer}
}

func (c *StatsController) HasSynced() bool {
	return c.informer.HasSynced()
}

// Run starts the kubewatch controller
func (c *StatsController) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	log.Printf("controller start")
	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("sync timeout"))
		return
	}

	log.Printf("controller synced and ready")

	// only one worker
	wait.Until(c.runWorker, time.Second, stopCh)
}

func (c *StatsController) runWorker() {
	for c.processNextItem() {
		// continue looping
	}
}

func (c *StatsController) processNextItem() bool {
	key, quit := c.queue.Get()

	if quit {
		return false
	}
	defer c.queue.Done(key)

	err := c.processItem(key.(string))
	if err == nil {
		// No error, reset the ratelimit counters
		c.queue.Forget(key)
	} else if c.queue.NumRequeues(key) < maxRetries {
		log.Printf("Error processing %s (will retry): %v", key, err)
		c.queue.AddRateLimited(key)
	} else {
		// err != nil and too many retries
		log.Printf("Error processing %s (giving up): %v", key, err)
		c.queue.Forget(key)
		utilruntime.HandleError(err)
	}

	return true
}

func (c *StatsController) processItem(key string) error {
	_, name, err := cache.SplitMetaNamespaceKey(key)
	obj, exists, err := c.informer.GetIndexer().GetByKey(key)
	if err != nil {
		return fmt.Errorf("Error fetching object with key %s from store: %v", key, err)
	}
	//log.Printf("processItem: %v\n", obj)

	statsObj, err := c.CrClient.StatsV1alpha1().Cpus(namespace).Get(c.CpuCrName, meta_v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("get custom resource[%s] failed: %v", c.CpuCrName, err)
	}

	statsObj = statsObj.DeepCopy()
	if statsObj.Status.Requests == nil {
		statsObj.Status.Requests = make(map[string]string)
	}

	if exists {
		// add or update event
		pod, _ := obj.(*api_v1.Pod)

		requestCpu := pod.Spec.Containers[0].Resources.Requests.Cpu() //TODO assuming only one container here
		log.Printf("add/update: %s, request cpu: %v\n", name, requestCpu)

		statsObj.Status.Requests[pod.Name] = requestCpu.String()
		_, err = c.CrClient.StatsV1alpha1().Cpus(namespace).Update(statsObj) //TODO loop get-update here
	} else {
		// delete event
		log.Printf("delete: %s", name)

		delete(statsObj.Status.Requests, name)
		_, err = c.CrClient.StatsV1alpha1().Cpus(namespace).Update(statsObj) //TODO loop get-update here
	}
	if err != nil {
		return err
	}

	return nil
}
