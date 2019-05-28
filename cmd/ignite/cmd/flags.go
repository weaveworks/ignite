package cmd

// This file contains a collection of variables for common flags
var (
	interactive bool
	all         bool
)

type createOptions struct {
	cpus   int64
	memory int64
}
