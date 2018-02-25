from prometheus_client import CollectorRegistry, Gauge, pushadd_to_gateway

registry = CollectorRegistry()
duration = Gauge('my_job_duration_seconds',
        'Duration of my batch job in seconds', registry=registry)
try:
    with duration.time():
        # Your code here. 
        pass

    # This only runs if there wasn't an exception. 
    g = Gauge('my_job_last_success_seconds',
            'Last time my batch job successfully finished', registry=registry)
    g.set_to_current_time()
finally:
    pushadd_to_gateway('localhost:9091', job='batch', registry=registry)
