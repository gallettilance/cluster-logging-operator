package alerts

const FluentdPrometheusAlert = `
"groups":
- "name": "logging_fluentd.alerts"
  "rules":
  - "alert": "FluentdNodeDown"
    "annotations":
      "message": "Prometheus could not scrape fluentd {{ $labels.container }} for more than 10m."
      "summary": "Fluentd cannot be scraped"
    "expr": |
      up{job = "collector", container = "collector"} == 0 or absent(up{job="collector", container="collector"}) == 1
    "for": "10m"
    "labels":
      "service": "collector"
      "severity": "critical"
      namespace: "openshift-logging"
  - "alert": "FluentdQueueLengthIncreasing"
    "annotations":
      "message": "For the last hour, fluentd {{ $labels.pod }} output '{{ $labels.plugin_id }}' average buffer queue length has increased continuously."
      "summary": "Fluentd pod {{ $labels.pod }} is unable to keep up with traffic over time for forwarder output {{ $labels.plugin_id }}."
    "expr": |
      sum by (pod,plugin_id) ( 0 * (deriv(fluentd_output_status_emit_records[1m] offset 1h)))  + on(pod,plugin_id)  ( deriv(fluentd_output_status_buffer_queue_length[10m]) > 0 and delta(fluentd_output_status_buffer_queue_length[1h]) > 1 )
    "for": "1h"
    "labels":
      "service": "collector"
      "severity": "Warning"
      namespace: "openshift-logging"
  - alert: FluentDHighErrorRate
    annotations:
      message: |-
        {{ $value }}% of records have resulted in an error by fluentd {{ $labels.instance }}.
      summary: FluentD output errors are high
    expr: |
      100 * (
        sum by(instance)(rate(fluentd_output_status_num_errors[2m]))
      /
        sum by(instance)(rate(fluentd_output_status_emit_records[2m]))
      ) > 10
    for: 15m
    labels:
      severity: warning
      namespace: "openshift-logging"
  - alert: FluentDVeryHighErrorRate
    annotations:
      message: |-
        {{ $value }}% of records have resulted in an error by fluentd {{ $labels.instance }}.
      summary: FluentD output errors are very high
    expr: |
      100 * (
        sum by(instance)(rate(fluentd_output_status_num_errors[2m]))
      /
        sum by(instance)(rate(fluentd_output_status_emit_records[2m]))
      ) > 25
    for: 15m
    labels:
      severity: critical
      namespace: "openshift-logging"
- "name": "logging_clusterlogging_telemetry.rules"
  "rules":
  - "expr": |
      sum by(cluster)(log_collected_bytes_total)
    "record": "cluster:log_collected_bytes_total:sum"
  - "expr": |
      sum by(cluster)(log_logged_bytes_total)
    "record": "cluster:log_logged_bytes_total:sum"
`
