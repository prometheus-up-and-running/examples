import json
import re
import time
from urllib.request import urlopen

from prometheus_client.core import GaugeMetricFamily, CounterMetricFamily
from prometheus_client.core import SummaryMetricFamily, REGISTRY
from prometheus_client import start_http_server


def sanitise_name(s):
    return re.sub(r"[^a-zA-Z0-9:_]", "_", s)

class ConsulCollector(object):
  def collect(self):
    out = urlopen("http://localhost:8500/v1/agent/metrics").read()
    metrics = json.loads(out.decode("utf-8"))

    for g in metrics["Gauges"]:
      yield GaugeMetricFamily(sanitise_name(g["Name"]),
          "Consul metric " + g["Name"], g["Value"])

    for c in metrics["Counters"]:
      yield CounterMetricFamily(sanitise_name(c["Name"]) + "_total",
          "Consul metric " + c["Name"], c["Count"])

    for s in metrics["Samples"]:
      yield SummaryMetricFamily(sanitise_name(s["Name"]) + "_seconds",
          "Consul metric " + s["Name"],
          count_value=c["Count"], sum_value=s["Sum"] / 1000)

if __name__ == '__main__':
  REGISTRY.register(ConsulCollector())
  start_http_server(8000)
  while True:
    time.sleep(1)



