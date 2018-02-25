import io.prometheus.client.Counter;
import io.prometheus.client.exporter.MetricsServlet;
import io.prometheus.client.hotspot.DefaultExports;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import javax.servlet.ServletException;
import org.eclipse.jetty.server.Server;
import org.eclipse.jetty.servlet.ServletContextHandler;
import org.eclipse.jetty.servlet.ServletHolder;
import java.io.IOException;


public class Example {
  static class ExampleServlet extends HttpServlet {
    private static final Counter requests = Counter.build()
        .name("hello_worlds_total")
        .help("Hello Worlds requested.").register();

    @Override
    protected void doGet(final HttpServletRequest req,
        final HttpServletResponse resp)
        throws ServletException, IOException {
      requests.inc();
      resp.getWriter().println("Hello World");
    }
  }

  public static void main(String[] args) throws Exception {
      DefaultExports.initialize();

      Server server = new Server(8000);
      ServletContextHandler context = new ServletContextHandler();
      context.setContextPath("/");
      server.setHandler(context);
      context.addServlet(new ServletHolder(new ExampleServlet()), "/");
      context.addServlet(new ServletHolder(new MetricsServlet()), "/metrics");

      server.start();
      server.join();
  }
}
