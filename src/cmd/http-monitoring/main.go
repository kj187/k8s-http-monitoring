package main

import (
	"github.com/kj187/http-monitoring/src/pkg/kubernetes"
	"github.com/kj187/http-monitoring/src/pkg/monitoring"
	"github.com/kj187/http-monitoring/src/pkg/prometheus"

	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

func main() {
	k8s := kubernetes.Kubernetes{}
	prom := prometheus.Prometheus{}
	mon := monitoring.Monitoring{}
	k8s.Init()
	mon.Init()
	go prom.Init()
	go k8s.IngressWatcherLoop()
	go k8s.ServiceWatcherLoop()
	go mon.MonitorIngressesLoop()
	go mon.MonitorServicesLoop()
	select {} // Wait forever
}
