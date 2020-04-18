package monitoring

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kj187/http-monitoring/src/pkg/kubernetes"
	"github.com/kj187/http-monitoring/src/pkg/prometheus"
)

type (
	// Monitoring Struct
	Monitoring struct {
		IntervalSeconds time.Duration
	}
	requestInput struct {
		networkType string
		url         string
		name        string
		namespace   string
		maintainer  string
		headers     []map[string]string
	}
)

// Constants
const (
	ListenOnAnnotationKey  = "kj187.de/http-monitoring"
	RequestProtocol        = "http"
	TimeoutSeconds         = 15
	DefaultIntervalSeconds = 30
)

// Init ...
func (mon *Monitoring) Init() {
	fmt.Println("Initializing HTTP Monitoring")
	mon.IntervalSeconds = DefaultIntervalSeconds
	if intervalSeconds := os.Getenv("MONITORING_INTERVAL_SECONDS"); intervalSeconds != "" {
		i, _ := strconv.ParseInt(intervalSeconds, 10, 64)
		mon.IntervalSeconds = time.Duration(i)
	}
}

// MonitorIngressesLoop ...
func (mon *Monitoring) MonitorIngressesLoop() {
	for {
		mon.MonitorIngresses()
		time.Sleep(mon.IntervalSeconds * time.Second)
	}
}

// MonitorIngresses ...
func (mon *Monitoring) MonitorIngresses() {
	var requestHeaders []map[string]string

	// TODO add generic headers, + add to docu
	idpToken := os.Getenv("MONITORING_X_IDP_TOKEN")
	//if idpToken == "" {
	//	fmt.Println("ENV MONITORING_X_IDP_TOKEN is not set")
	//}

	for _, i := range kubernetes.Ingresses {
		if _, ok := i.Annotations[ListenOnAnnotationKey]; ok == false {
			continue
		}
		host := i.Spec.Rules[0].Host // We have just 1 rule in our project
		paths := bufio.NewScanner(strings.NewReader(i.Annotations[ListenOnAnnotationKey]))
		if i.Annotations["kubernetes.io/ingress.class"] != "public" { // use IDP token if ingress is not public
			requestHeaders = append(requestHeaders, map[string]string{"X-Idp-Token": idpToken})
		}
		requestInput := requestInput{
			networkType: "ingress",
			name:        i.Name,
			namespace:   i.Namespace,
			maintainer:  i.Labels["om3.cloud/maintainer"],
			headers:     requestHeaders,
		}
		for paths.Scan() {
			path := strings.TrimPrefix(strings.TrimSpace(paths.Text()), "/")
			requestInput.url = fmt.Sprintf("%v://%v/%v", RequestProtocol, host, path)
			err := mon.executeRequest(requestInput)
			if err != nil {
				fmt.Println("[ERROR][#1576760441] ", err)
			}
		}
	}
}

// MonitorServicesLoop ...
func (mon *Monitoring) MonitorServicesLoop() {
	for {
		mon.MonitorServices()
		time.Sleep(mon.IntervalSeconds * time.Second)
	}
}

// MonitorServices ...
func (mon *Monitoring) MonitorServices() {
	var requestHeaders []map[string]string
	for _, i := range kubernetes.Services {
		if _, ok := i.Annotations[ListenOnAnnotationKey]; ok == false {
			continue
		}
		requestInput := requestInput{
			networkType: "service",
			name:        i.Name,
			namespace:   i.Namespace,
			maintainer:  i.Labels["om3.cloud/maintainer"],
			headers:     requestHeaders,
		}
		paths := bufio.NewScanner(strings.NewReader(i.Annotations[ListenOnAnnotationKey]))
		for paths.Scan() {
			path := strings.TrimPrefix(strings.TrimSpace(paths.Text()), "/")
			requestInput.url = fmt.Sprintf("%v://%v.%v/%v", RequestProtocol, i.Name, i.Namespace, path)
			err := mon.executeRequest(requestInput)
			if err != nil {
				fmt.Println("[ERROR][#1576828658] ", err)
			}
		}
	}
}

func (mon *Monitoring) executeRequest(input requestInput) error {
	success := float64(1)
	start := time.Now()
	client := &http.Client{Timeout: TimeoutSeconds * time.Second}
	request, err := http.NewRequest("GET", input.url, nil)
	if err != nil {
		return err
	}

	if len(input.headers) > 0 {
		for _, header := range input.headers {
			for key, value := range header {
				request.Header.Set(key, value)
			}
		}
	}

	r, err := client.Do(request)
	if err != nil {
		return err
	}

	durationSeconds := time.Since(start).Seconds()
	if err != nil {
		fmt.Println(err)
		success = float64(0)
	}

	if r.StatusCode != 200 {
		success = float64(0)
	}

	prometheus.ProbeSuccess.WithLabelValues(input.name, input.namespace, input.maintainer, input.networkType, input.url).Set(success)
	prometheus.ProbeStatusCode.WithLabelValues(input.name, input.namespace, input.maintainer, input.networkType, input.url).Set(float64(r.StatusCode))
	prometheus.ProbeDurationSeconds.WithLabelValues(input.name, input.namespace, input.maintainer, input.networkType, input.url).Set(durationSeconds)
	prometheus.ProbeLastExecutionTime.WithLabelValues(input.name, input.namespace, input.maintainer, input.networkType, input.url).SetToCurrentTime()

	//fmt.Printf("[%v][%v] %v\n", input.networkType, r.StatusCode, input.url)
	return nil
}
