package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/pkg/providers"
	"io"
	"os"
	"os/signal"
	"sync"
)

func NewCmdHang(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "hang",
		Short:  "Causes Ignite to hang indefinitely. Used for testing purposes.",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			var endWaiter sync.WaitGroup
			endWaiter.Add(1)

			signalChannel := make(chan os.Signal, 1)
			signal.Notify(signalChannel, os.Interrupt)

			go func() {
				<-signalChannel
				endWaiter.Done()
			}()

			endWaiter.Wait()

			if providers.SS != nil {
				fmt.Println("Closing...")
				providers.SS.Close()
			}
		},
	}

	return cmd
}
