package main

import (
  "log"
  "net/http"
  "regexp"

  "github.com/hashicorp/consul/api"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
  up = prometheus.NewDesc(
    "consul_up",
    "Was talking to Consul successful.",
    nil, nil,
  )
  invalidChars = regexp.MustCompile("[^a-zA-Z0-9:_]")
)

type ConsulCollector struct {
}

// Implements prometheus.Collector.
func (c ConsulCollector) Describe(ch chan<- *prometheus.Desc) {
  ch <- up
}

// Implements prometheus.Collector.
func (c ConsulCollector) Collect(ch chan<- prometheus.Metric) {
  consul, err := api.NewClient(api.DefaultConfig())
  if err != nil {
    ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0)
    return
  }

  metrics, err := consul.Agent().Metrics()
  if err != nil {
    ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0)
    return
  }
  ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1)

  for _, g := range metrics.Gauges {
    name := invalidChars.ReplaceAllLiteralString(g.Name, "_")
    desc := prometheus.NewDesc(name, "Consul metric "+g.Name, nil, nil)
    ch <- prometheus.MustNewConstMetric(
        desc, prometheus.GaugeValue, float64(g.Value))
  }

  for _, c := range metrics.Counters {
    name := invalidChars.ReplaceAllLiteralString(c.Name, "_")
    desc := prometheus.NewDesc(name+"_total", "Consul metric "+c.Name, nil, nil)
    ch <- prometheus.MustNewConstMetric(
        desc, prometheus.CounterValue, float64(c.Count))
  }

  for _, s := range metrics.Samples {
    // All samples are times in milliseconds, we convert them to seconds below.
    name := invalidChars.ReplaceAllLiteralString(s.Name, "_") + "_seconds"
    countDesc := prometheus.NewDesc(
        name+"_count", "Consul metric "+s.Name, nil, nil)
    ch <- prometheus.MustNewConstMetric(
        countDesc, prometheus.CounterValue, float64(s.Count))
    sumDesc := prometheus.NewDesc(
        name+"_sum", "Consul metric "+s.Name, nil, nil)
    ch <- prometheus.MustNewConstMetric(
        sumDesc, prometheus.CounterValue, s.Sum/1000)
  }
}

func main() {
  c := ConsulCollector{}
  prometheus.MustRegister(c)
  http.Handle("/metrics", promhttp.Handler())
  log.Fatal(http.ListenAndServe(":8000", nil))
}
