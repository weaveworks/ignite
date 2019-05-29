package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdImageRm removes an images
// TODO: Support removing multiple images at once
func NewCmdImageRm(out io.Writer) *cobra.Command {
	ro := &run.RmiOptions{}

	cmd := &cobra.Command{
		Use:   "rm [image]",
		Short: "Remove a VM base image",
		Long:  "TODO", // TODO: Long description
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if ro.Image, err = matchSingleImage(args[0]); err != nil {
					return err
				}
				return run.Rmi(ro)
			}())
		},
	}

	return cmd
}
