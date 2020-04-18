package prometheus

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// DefaultMetricsPort ...
const DefaultMetricsPort = "1877"

var (
	// ProbeSuccess ...
	ProbeSuccess = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "http_monitoring_probe_success"}, []string{"monitor_app", "monitor_namespace", "om3_cloud_maintainer", "monitor_network", "monitor_url"})

	// ProbeStatusCode ...
	ProbeStatusCode = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "http_monitoring_probe_status_code"}, []string{"monitor_app", "monitor_namespace", "om3_cloud_maintainer", "monitor_network", "monitor_url"})

	// ProbeDurationSeconds ...
	ProbeDurationSeconds = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "http_monitoring_probe_duration_seconds"}, []string{"monitor_app", "monitor_namespace", "om3_cloud_maintainer", "monitor_network", "monitor_url"})

	// ProbeLastExecutionTime ...
	ProbeLastExecutionTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "http_monitoring_last_execution_time"}, []string{"monitor_app", "monitor_namespace", "om3_cloud_maintainer", "monitor_network", "monitor_url"})
)

// Prometheus Struct
type Prometheus struct{}

// Init ...
func (prom *Prometheus) Init() {
	fmt.Printf("Initializing Prometheus Metrics Server (host:%v/metrics)\n", DefaultMetricsPort)
	prometheus.MustRegister(ProbeSuccess, ProbeStatusCode, ProbeDurationSeconds, ProbeLastExecutionTime)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":"+DefaultMetricsPort, nil)
}
