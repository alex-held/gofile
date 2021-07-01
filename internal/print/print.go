package main

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func (g GoFile) PrettyPrint() (string, error) {
	sb := &strings.Builder{}
	err := g.PrettyPrintW(sb)
	return sb.String(), err
}

func (g GoFile) PrettyPrintW(w io.Writer) (err error) {
	buf := &bytes.Buffer{}
	table := tablewriter.NewWriter(buf)

	var data [][]string
	for _, i := range g.Installables {
		data = append(data, []string{i.Name, i.URI})
	}

	table.SetHeader([]string{"NAME", "URI"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(false)
	table.AppendBulk(data)
	table.Render()

	r := bufio.NewReader(buf)
	for {
		line, _, err := r.ReadLine()
		// if err != nil || len(line) == 0 {
		if err != nil && err != io.EOF || len(line) == 0 {
			break
		}
		prefixed := strings.TrimPrefix(string(line), "  ")
		prefixed = prefixed + "\n"
		_, err = w.Write([]byte(prefixed))
		if err != nil {
			break
		}
	}

	return err
}

func PrettyPrint(val PrettyPrintable) (s string, err error) {
	w := &strings.Builder{}
	err = PrettyPrintW(val, w)
	if err != nil {
		return "", err
	}
	return w.String(), nil
}

func PrettyPrintW(val PrettyPrintable, w io.Writer) (err error) {
	return val.PrettyPrintW(w)
}

type PrettyPrintable interface {
	PrettyPrint() (str string, err error)
	PrettyPrintW(w io.Writer) (err error)
}
