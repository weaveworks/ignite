package prometheus

import (
	"net/http"
	"net"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ServeMetrics(socketPath string) error {
	// Create a registry to register metrics into
	registry := prometheus.NewRegistry()
	// Register the default collectors with this registry
	registry.MustRegister(prometheus.NewBuildInfoCollector())
	registry.MustRegister(prometheus.NewGoCollector())
	registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))

	// Create a HTTP server to handle metrics requests
	server := http.Server{
		Handler: promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
	}

	// Listen on the unix socket and serve the web server
	unixListener, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}
	return server.Serve(unixListener)
}
