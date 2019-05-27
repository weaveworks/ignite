package util

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"
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

func (o *output) Write(input ...interface{}) {
	var sb strings.Builder
	for i, data := range input {
		switch data.(type) {
		case string:
			sb.WriteString(fmt.Sprintf("%s", data))
		case int64:
			sb.WriteString(fmt.Sprintf("%d", data))
		case time.Time:
			sb.WriteString(data.(time.Time).Format(time.UnixDate))
		default:
			sb.WriteString(fmt.Sprintf("%v", data))
		}

		if i+1 < len(input) {
			sb.WriteString("\t")
		} else {
			sb.WriteString("\n")
		}
	}

	fmt.Fprint(o.writer, sb.String())
}

func (o *output) Flush() {
	o.writer.Flush()
}
