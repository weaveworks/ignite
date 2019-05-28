package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/util"
	"github.com/spf13/pflag"
	"io"
	"os"
	"path/filepath"

	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
)

type startOptions struct {
	attachOptions
	interactive bool
}

// NewCmdStart starts a Firecracker VM
func NewCmdStart(out io.Writer) *cobra.Command {
	so := &startOptions{}

	cmd := &cobra.Command{
		Use:   "start [vm]",
		Short: "Start a Firecracker VM",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if so.vm, err = matchSingleVM(args[0]); err != nil {
					return err
				}
				return RunStart(so)
			}())
		},
	}

	addStartFlags(cmd.Flags(), so)
	return cmd
}

func addStartFlags(fs *pflag.FlagSet, so *startOptions) {
	addInteractiveFlag(fs, &so.interactive)
}

func RunStart(so *startOptions) error {
	// Check if the given VM is already running
	if so.vm.Running() {
		return fmt.Errorf("%s is already running", so.vm.ID)
	}

	igniteBinary, _ := filepath.Abs(os.Args[0])

	dockerArgs := []string{
		"run",
		"-itd",
		"--rm",
		"--name",
		so.vm.ID,
		fmt.Sprintf("-v=%s:/ignite/ignite", igniteBinary),
		fmt.Sprintf("-v=%s:%s", constants.DATA_DIR, constants.DATA_DIR),
		fmt.Sprintf("--stop-timeout=%d", constants.STOP_TIMEOUT+constants.IGNITE_TIMEOUT),
		"--privileged",
		"--device=/dev/kvm",
		"ignite",
		so.vm.ID,
	}

	// Start the VM in docker
	if _, err := util.ExecuteCommand("docker", dockerArgs...); err != nil {
		return fmt.Errorf("failed to start container for VM %q: %v", so.vm.ID, err)
	}

	// If starting interactively, attach after starting
	if so.interactive {
		if err := RunAttach(&so.attachOptions); err != nil {
			return err
		}
	} else {
		// Print the ID of the started VM
		fmt.Println(so.vm.ID)
	}

	return nil
}
