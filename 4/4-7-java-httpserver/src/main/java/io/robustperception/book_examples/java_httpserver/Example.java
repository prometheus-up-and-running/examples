import io.prometheus.client.Counter;
import io.prometheus.client.hotspot.DefaultExports;
import io.prometheus.client.exporter.HTTPServer;

public class Example {
  private static final Counter myCounter = Counter.build()
      .name("my_counter_total")
      .help("An example counter.").register();

  public static void main(String[] args) throws Exception {
    DefaultExports.initialize();
    HTTPServer server = new HTTPServer(8000);
    while (true) {
      myCounter.inc();
      Thread.sleep(1000);
    }
  }
}
