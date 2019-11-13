package reconcile

import (
	"path"

	go_prom "github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/prometheus"
)

var (
	vmCreated = go_prom.NewCounter(go_prom.CounterOpts{
		Name: "vm_create_counter",
		Help: "The count of VMs created",
	})
	vmDeleted = go_prom.NewCounter(go_prom.CounterOpts{
		Name: "vm_delete_counter",
		Help: "The count of VMs deleted",
	})
	vmStarted = go_prom.NewCounter(go_prom.CounterOpts{
		Name: "vm_start_counter",
		Help: "The count of VMs started",
	})
	vmStopped = go_prom.NewCounter(go_prom.CounterOpts{
		Name: "vm_stop_counter",
		Help: "The count of VMs stopped",
	})
	kindIgnored = go_prom.NewCounter(go_prom.CounterOpts{
		Name: "kind_ignored_counter",
		Help: "A counter of non-vm manifests ignored",
	})
)

func startMetricsThread() {
	reg, server := prometheus.New()
	reg.MustRegister(vmCreated, vmDeleted, vmStarted, vmStopped, kindIgnored)

	go func() {
		// create a new registry and http.Server. don't register custom metrics to the registry quite yet
		metricsSocket := path.Join(constants.DATA_DIR, constants.DAEMON_SOCKET)
		if err := prometheus.ServeOnSocket(server, metricsSocket); err != nil {
			log.Errorf("prometheus server was stopped with error: %v", err)
		}
	}()
}
