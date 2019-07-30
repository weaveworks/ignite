package cmd

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/pkg/providers"
)

func NewCmdDaemon(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "daemon",
		Short:  "For now causes Ignite to hang indefinitely. Used for testing purposes.",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize the daemon providers (e.g. ManifestStorage)
			cmdutil.CheckErr(providers.Populate(providers.DaemonProviders))

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

			// Close the SyncStorage's watcher threads
			if providers.SyncStorage != nil {
				fmt.Println("Closing...")
				providers.SyncStorage.Close()
			}
		},
	}

	return cmd
}
