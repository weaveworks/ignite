package client

import (
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/providers"
)

func SetClient() (err error) {
	log.Trace("Initializing the Client provider...")
	providers.Client = client.NewClient(providers.Storage)
	return
}
