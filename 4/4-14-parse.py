from prometheus_client.parser import text_string_to_metric_families

for family in text_string_to_metric_families(u"counter_total 1.0\n"):
  for sample in family.samples:
    print("Name: {0} Labels: {1} Value: {2}".format(*sample))
