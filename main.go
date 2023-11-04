package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/metrics"

	"github.com/olekukonko/tablewriter"
)

func main() {
	w := tableWriter{W: os.Stdout}
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: %s [flags]\n\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Prints a list of all runtime/metrics and their properties.\n\n")
		flag.PrintDefaults()
	}
	flag.BoolVar(&w.CSV, "csv", false, "output in CSV format")
	flag.Parse()

	w.SetHeader([]string{"Name", "Kind", "Cumulative", "Description"})
	descs := metrics.All()
	for _, d := range descs {
		w.Append([]string{
			d.Name,
			valueKindString(d.Kind),
			fmt.Sprintf("%v", d.Cumulative),
			d.Description,
		})
	}
	w.Flush()
}

type tableWriter struct {
	W      io.Writer
	CSV    bool
	header []string
	rows   [][]string
}

func (t *tableWriter) SetHeader(header []string) {
	t.header = header
}

func (t *tableWriter) Append(row []string) {
	t.rows = append(t.rows, row)
}

func (t *tableWriter) Flush() {
	if t.CSV {
		cw := csv.NewWriter(t.W)
		cw.Write(t.header)
		for _, row := range t.rows {
			cw.Write(row)
		}
		cw.Flush()
		return
	}
	table := tablewriter.NewWriter(t.W)
	table.SetHeader(t.header)
	for _, row := range t.rows {
		table.Append(row)
	}
	table.Render()
}

func valueKindString(v metrics.ValueKind) string {
	switch v {
	case metrics.KindBad:
		return "KindBad"
	case metrics.KindUint64:
		return "KindUint64"
	case metrics.KindFloat64:
		return "KindFloat64"
	case metrics.KindFloat64Histogram:
		return "KindFloat64Histogram"
	default:
		return "Unknown"
	}
}
