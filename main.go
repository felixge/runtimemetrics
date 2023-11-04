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
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: %s [flags]\n\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Prints a list of all runtime/metrics and their properties.\n\n")
		flag.PrintDefaults()
	}
	formatF := flag.String("format", "markdown", "output format")
	flag.Parse()

	w, err := newTableWriter(os.Stdout, format(*formatF))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

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

type format string

const (
	formatCSV      format = "csv"
	formatMarkdown format = "markdown"
)

func newTableWriter(w io.Writer, format format) (*tableWriter, error) {
	switch format {
	case formatCSV, formatMarkdown:
	default:
		return nil, fmt.Errorf("unknown format %q", format)
	}
	return &tableWriter{w: w, format: format}, nil
}

type tableWriter struct {
	w      io.Writer
	format format
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
	switch t.format {
	case formatCSV:
		cw := csv.NewWriter(t.w)
		cw.Write(t.header)
		for _, row := range t.rows {
			cw.Write(row)
		}
		cw.Flush()
	case formatMarkdown:
		table := tablewriter.NewWriter(t.w)
		table.SetAutoWrapText(false)
		table.SetHeader(t.header)
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
		for _, row := range t.rows {
			table.Append(row)
		}
		table.Render()
	}
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
