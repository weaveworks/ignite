package cmd

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/pkg/errutils"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/providers/ignite"
)

func NewCmdDaemon(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "daemon",
		Short:  "For now causes Ignite to hang indefinitely. Used for testing purposes.",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize the daemon providers (e.g. ManifestStorage)
			errutils.Check(providers.Populate(ignite.DaemonProviders))

			// Wait for Ctrl + C
			var endWaiter sync.WaitGroup
			endWaiter.Add(1)

			signalChannel := make(chan os.Signal, 1)
			signal.Notify(signalChannel, os.Interrupt)

			go func() {
				<-signalChannel
				endWaiter.Done()
			}()

			endWaiter.Wait()

			// Close the Storage's watcher threads
			fmt.Println("Closing...")
			providers.Storage.Close()
		},
	}

	return cmd
}
