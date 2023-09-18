package cmd

import (
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

const creationDateLayout = "2006-01-02T15:04:05.000 MST"

func newTabWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
}

func parseCreateDate(ms string, localTime bool) string {
	ts, _ := strconv.ParseInt(ms, 10, 64)
	t := time.UnixMilli(ts)
	if !localTime {
		t = t.UTC()
	}
	return t.Format(creationDateLayout)
}
