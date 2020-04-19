# HTTP Monitoring

[![pipeline](https://github.com/kj187/k8s-http-monitoring/workflows/pipeline/badge.svg?branch=master)](https://github.com/kj187/k8s-http-monitoring/actions?query=workflow%3Apipeline)

Kubernetes Ingress and Service monitoring of HTTP/HTTPS requests

These application will continuously check your application availability in a specific interval. 
You can check one or multiple endpoints for `Ingress` or/and `Service` resources. 
Beside a success or fail status, this application exposed also the HTTP status code, the request duration in seconds and the last execution time of the check itself as a metric endoint which can be scaped by Prometheus. 

## Installing

### Installing via Helm Chart

```bash
$ git clone https://github.com/kj187/k8s-http-monitoring.git
$ cd http-monitoring/chart
$ helm upgrade --install http-monitoring ./http-monitoring
```

### Installing via kubectl

```bash
$ git clone https://github.com/kj187/k8s-http-monitoring.git
$ cd http-monitoring
$ kubectl apply -f deploy/
```

## Usage

Works with Kubernetes `Ingress` and `Service` resources.
With the `Ingress` resource you can check the external access and with the `Service` resource you can check the internal access of you application.

Add a annotation to one of these resources.

Example:
```yaml
  annotation:
    kj187.de/http-monitoring: /    # Root page
```

Or if you want to check multiple endpoints of your application
```yaml
  annotation:
    kj187.de/http-monitoring: |
      /
      /custeromer-area
      /whatever/else
```

### ENV vars

Available ENV vars which could be injected (e.g. via Helm Chart)

```bash
MONITORING_INTERVAL_SECONDS="SECONDS_AS_INT"  # default is 30 seconds
```

### Metrics

The metric endpoint is not available from outside (if you need this, you have to create a `Ingress` resource)
Inside the cluster the endpoint is available under `http-monitoring/metrics` (<SERVICE_NAME>/metrics)

Exposed metrics

```bash
http_monitoring_probe_success           # 1 = success, 0 = failed
http_monitoring_probe_status_code       # HTTP status code
http_monitoring_probe_duration_seconds  # Duration of request in seconds
http_monitoring_last_execution_time     # Number of seconds since 1970 of last garbage collection.
```

All metrics have these lables: 
- `monitor_app` is the name of your application
- `monitor_namespace` is the namespace where your application lives
- `monitor_network` ingress or service 
- `monitor_url` the URL of the check

### Prometheus scraping

Prometheus need some information about the http-monitoring metric endpoint. 
If you are not using the Prometheus-Operator, just the normal Prometheus, you need the following annotations in the service resource, which is already available per default.

```yaml
  annotations:
    prometheus.io/scrape: "true"
    prometheus.io/path: /metrics
    prometheus.io/port: "8080"
``` 

If you are using the Prometheus-Operator, just create a new `ServiceMonitor` resource:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: http-monitoring-servicemonitor
  labels:
    prometheus: kube-prometheus
spec:
  jobLabel: app
  targetLabels:
    - app
  namespaceSelector:
    any: true
  selector:
    matchLabels:
      prometheus-operator/scrape: "default"
  endpoints:
    - port: metrics
      interval: 10s
```

Now Prometheus should scrape the http-monitoring metrics continuously.

### Self-check

To check if the http-monitoring application itself works properly, you can check the `http_monitoring_last_execution_time` metric which exposed a unix timestamp of the latest execution time

If you are not using the Prometheus-Operator, create a new file called `prometheus_values.yaml` and add the following content
```
serverFiles:
  alerting_rules.yml:
    groups:
      - name: http-monitoring
        rules:
          - alert: HttpMonitoringChecksAreOlderThan1Day
            expr: http_monitoring_last_execution_time < (time()-86400)
            for: 10m
            labels:
              severity: high
            annotations:
              summary: "HTTP monitoring checks are older than 1 day!"
              urgency: high
```

Now you have to update your Prometheus Helm chart
```
helm upgrade prometheus stable/prometheus -f prometheus_values.yaml
```

If you are using the Prometheus-Operator, just add the following new `PrometheusRule` resource

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: http-monitoring-rules
  labels:
    prometheus: kube-prometheus
    release: {{ .Release.Name }}
    cronjob: http-monitoring-rules
spec:
  groups:
  - name: http-monitoring
    rules:
    - alert: HttpMonitoringChecksAreOlderThan1Day
      expr: http_monitoring_last_execution_time < (time()-86400)
      labels:
        severity: high
      annotations:
        summary: "HTTP monitoring checks are older than 1 day!"
        urgency: high
```