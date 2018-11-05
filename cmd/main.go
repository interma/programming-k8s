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
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" // for google cloud auth
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig string // kube config path
	cpuCR	string // cpu custom resource
)

//usage:
//go run main.go -kubeconfig=$HOME/.kube/config
func main() {
	flag.Parse()
	masterURL := ""

	stopCh := make(chan struct{})
	defer close(stopCh)

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
		log.Fatalf("Error building customresource clientset: %s", err.Error())
	}

	c := controller.CreatePodsStatsController(KubeClient, CrClient, cpuCR)
	go c.Run(stopCh)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm

	log.Printf("main exited")
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&cpuCR, "cpuCR", "cpu-sample", "the name of cpu custom resource.")
}
