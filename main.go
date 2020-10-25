package main

import (
	"log"
	"os"
	"sync"

	"github.com/serviceMonitor/m/alert"
	k8s "github.com/serviceMonitor/m/k8s"
)

func main() {

	/*
		for test
	*/

	var k k8s.KubeConfig
	k.Path = os.Getenv("k8s_config_path")
	if k.Path == "" {
		log.Fatal("Must set kubernetes config file path <k8s_config_path> env")
	}
	ns := os.Getenv("k8s_service_namespace")
	if ns == "" {
		log.Fatal("Must set service namespace <k8s_service_namespace> env")
	}
	labels := os.Getenv("k8s_service_lables")
	if labels == "" {
		log.Fatal("Must set service lables <k8s_service_lables> env")
	}

	webhook := os.Getenv("dingtalk_webhook")
	if webhook == "" {
		log.Fatal("Must set dingtalk webhook <dingtalk_webhook> env")
	}

	clientset := k.NewClient()
	m := make(map[string]string)
	mux := &sync.RWMutex{}
	wg := &sync.WaitGroup{}
	ch := make(chan k8s.Msg)
	// handle expection message
	except := make(chan k8s.ExceptMsg)

	defer close(ch)
	go k8s.LiveLoop(clientset, ns, labels, ch)
	go k8s.WriteLoop(m, mux, ch)
	go k8s.ReadLoop(m, mux, wg, except)
	go alert.Exception(webhook, except)
	// stop program from exiting, must be killed
	printer := make(chan struct{})
	<-printer
}
