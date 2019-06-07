package run

type RunOptions struct {
	CreateOptions
	StartOptions
}

func Run(ro *RunOptions) (string, error) {
	if _, err := Create(&ro.CreateOptions); err != nil {
		return "", err
	}

	ro.StartOptions.VM = ro.CreateOptions.vm
	return Start(&ro.StartOptions)
}
