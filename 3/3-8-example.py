import time
from prometheus_client import start_http_server
from prometheus_client import Gauge

TIME = Gauge('time_seconds',
             'The current time.')
TIME.set_function(lambda: time.time())

if __name__ == "__main__":
    start_http_server(8000)
    while True:
        time.sleep(1)
