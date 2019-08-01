package reconcile

import (
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/prometheus"
)

func startMetricsThread() {
	_, server := prometheus.New()
	go func() {
		// create a new registry and http.Server. don't register custom metrics to the registry quite yet
		metricsSocket := path.Join(constants.DATA_DIR, constants.DAEMON_SOCKET)
		if err := prometheus.ServeOnSocket(server, metricsSocket); err != nil {
			log.Errorf("prometheus server was stopped with error: %v", err)
		}
	}()
}
