package cmd

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/pkg/providers"
)

func NewCmdDaemon(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "daemon",
		Short: "Operates in daemon mode and watches /etc/firecracker/manifests for VM specifications to run.", // TODO: Parameterize
		Run: func(cmd *cobra.Command, args []string) {
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
