package run

type RunOptions struct {
	CreateOptions
	StartOptions
}

func Run(ro *RunOptions) error {
	if err := Create(&ro.CreateOptions); err != nil {
		return err
	}

	ro.StartOptions.VM = ro.CreateOptions.vm

	if err := Start(&ro.StartOptions); err != nil {
		return err
	}

	return nil
}
