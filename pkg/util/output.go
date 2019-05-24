package util

import (
	"fmt"
	"os"
	"text/tabwriter"
)

type output struct {
	writer *tabwriter.Writer
}

func NewOutput() *output {
	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 8, 1, '\t', 0)
	return &output{
		writer: writer,
	}
}

func (o *output) Write(s string) {
	fmt.Fprint(o.writer, s+"\n")
}

func (o *output) Flush() {
	o.writer.Flush()
}
