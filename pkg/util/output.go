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

func ByteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
