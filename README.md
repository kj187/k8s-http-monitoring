# HTTP Monitoring

[![pipeline](https://github.com/kj187/http-monitoring/workflows/pipeline/badge.svg?branch=master)](https://github.com/kj187/http-monitoring/actions?query=workflow%3Apipeline)

Monitoring of HTTP requests.

Works with Kubernetes `Ingress` and `Service` resources. Just add a annotation to one of these resources. 

Example:

```yaml
  annotation:
    kj187.de/http-monitoring: |
      /
      /en-EN
      /search/...
```


Available ENV vars which could be injected (via Helm Chart)
```bash
MONITORING_INTERVAL_SECONDS # default is 30 seconds
```

**Metrics**

The following metrics have all these lables available: `monitor_app`, `monitor_namespace`, `monitor_network`, `monitor_url`

```bash
http_monitoring_probe_success           # 1 = success, 0 = failed
http_monitoring_probe_status_code       # HTTP status code
http_monitoring_probe_duration_seconds  # Duration of request in seconds
http_monitoring_last_execution_time     # Number of seconds since 1970 of last garbage collection.
```

To check if the application itself works properly you can check another metric with a unix timestamp as the latest execution time
Prometheus Rule example:

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

**Metrics Endpoint** 

- Reachable for local development: `127.0.0.1:8080/metrics`
- Reachable inside cluster as service: `http-monitoring-metrics.k28s-infrastructure/metrics`

## Installing

### Installing via Helm Chart

```
$ git clone https://github.com/kj187/http-monitoring.git
$ cd http-monitoring/chart
$ helm lint http-monitoring ./http-monitoring
$ helm upgrade --install http-monitoring ./http-monitoring
```

### Installing via kubectl

```
$ git clone https://github.com/kj187/http-monitoring.git
$ cd http-monitoring
$ kubectl apply -f deploy/
```