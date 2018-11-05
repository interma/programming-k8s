package main

import (
	"flag"
	"log"

	"os"
	"os/signal"
	"syscall"

	clientset "github.com/interma/programming-k8s/pkg/client/clientset/versioned"
	"github.com/interma/programming-k8s/pkg/controller"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig string
	KubeClient kubernetes.Interface // global kubeClient
	CrClient   clientset.Interface  // global crClient
)

//usage:
//go run main.go -kubeconfig=$HOME/.kube/config
func main() {
	flag.Parse()
	masterURL := ""

	// init kubeclient by kubeconfig
	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	KubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	CrClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("Error building hawqcluster clientset: %s", err.Error())
	}

	c := controller.CreatePodsStatsController(KubeClient, CrClient)

	stopCh := make(chan struct{})
	defer close(stopCh)

	go c.Run(stopCh)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm

	log.Printf("main exited")
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
}
