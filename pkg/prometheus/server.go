package prometheus

import (
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New() (*prometheus.Registry, *http.Server) {
	// Create a registry to register metrics into
	registry := prometheus.NewRegistry()
	// Register the default collectors with this registry
	registry.MustRegister(prometheus.NewBuildInfoCollector())
	registry.MustRegister(prometheus.NewGoCollector())
	registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))

	// Create a HTTP server to handle metrics requests
	return registry, &http.Server{
		Handler: promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
	}
}

func ServeOnSocket(server *http.Server, socketPath string) error {
	// Listen on the unix socket and serve the web server
	unixListener, err := net.Listen("unix", socketPath)
	if err != nil {
		return err
	}

	return server.Serve(unixListener)
}
