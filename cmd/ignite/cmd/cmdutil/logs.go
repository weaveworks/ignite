package cmdutil

import "fmt"

func PrintMachineReadableID(id string, err error) error {
	if err != nil {
		return err
	}
	fmt.Println(id)
	return nil
}
