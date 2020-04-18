package kubernetes

import (
	"fmt"
	"log"
	"time"

	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	k8s_api_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	k8s_go "k8s.io/client-go/kubernetes"
	k8s_clientcmd "k8s.io/client-go/tools/clientcmd"
)

var (
	// Ingresses with monitoring annotation
	Ingresses = make(map[types.UID]*v1beta1.Ingress)

	// Services with monitoring annotation
	Services = make(map[types.UID]*k8s_api_v1.Service)
)

// Kubernetes Struct
type Kubernetes struct {
	clientset *k8s_go.Clientset
}

// Init ...
func (k8s *Kubernetes) Init() {
	if k8s.clientset != nil {
		return
	}

	fmt.Println("Initializing Kubernetes")
	kubeconfig := k8s_clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		k8s_clientcmd.NewDefaultClientConfigLoadingRules(),
		&k8s_clientcmd.ConfigOverrides{},
	)

	restconfig, err := kubeconfig.ClientConfig()
	if err != nil {
		log.Fatal(fmt.Sprintf("[ERROR][#1563220691] %v", err))
		return
	}

	k8s.clientset, err = k8s_go.NewForConfig(restconfig)
	if err != nil {
		log.Fatal(fmt.Sprintf("[ERROR][#1563220710] %v", err))
		return
	}
}

// IngressWatcherLoop ...
func (k8s *Kubernetes) IngressWatcherLoop() {
	for {
		k8s.IngressWatcher()
		time.Sleep(1 * time.Second)
	}
}

// IngressWatcher ...
func (k8s *Kubernetes) IngressWatcher() {
	watcher, err := k8s.clientset.ExtensionsV1beta1().Ingresses("").Watch(v1.ListOptions{})
	if err != nil {
		log.Fatal(fmt.Sprintf("[ERROR][#1576760541] %v", err))
		return
	}

	ch := watcher.ResultChan()
	for event := range ch {
		ingress, ok := event.Object.(*v1beta1.Ingress)
		if !ok {
			log.Fatal("[ERROR][#1576760491] unexpected type")
			continue
		}

		switch event.Type {
		case watch.Added, watch.Modified:
			Ingresses[ingress.UID] = ingress.DeepCopy()

		case watch.Deleted:
			delete(Ingresses, ingress.UID)
		}
	}
}

// ServiceWatcherLoop ...
func (k8s *Kubernetes) ServiceWatcherLoop() {
	for {
		k8s.ServiceWatcher()
		time.Sleep(1 * time.Second)
	}
}

// ServiceWatcher ...
func (k8s *Kubernetes) ServiceWatcher() {
	watcher, err := k8s.clientset.CoreV1().Services("").Watch(v1.ListOptions{})
	if err != nil {
		log.Fatal(fmt.Sprintf("[ERROR][#1576828226] %v", err))
		return
	}

	ch := watcher.ResultChan()
	for event := range ch {
		service, ok := event.Object.(*k8s_api_v1.Service)
		if !ok {
			log.Fatal("[ERROR][#1576828234] unexpected type")
			continue
		}

		switch event.Type {
		case watch.Added, watch.Modified:
			Services[service.UID] = service.DeepCopy()
		case watch.Deleted:
			delete(Services, service.UID)
		}
	}
}
