package gofile

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

type GoFile struct {
	Installables []Installable
}

type Installable struct {
	Name, URI    string
	author, host string
}

func New(installables ...Installable) (f *GoFile, err error) {
	return &GoFile{Installables: installables}, nil
}

func SaveFile(path string, file *GoFile) (f *os.File, err error) {
	buffer := &bytes.Buffer{}

	for _, installable := range file.Installables {
		if _, err = buffer.WriteString(installable.URI + "\n"); err != nil {
			return nil, err
		}
	}

	f, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return f, err
	}
	defer f.Close()

	_, err = f.Write(buffer.Bytes())
	return f, err
}

func ResolveGofilePath(rawPath string) (path string, err error) {
	abs, err := filepath.Abs(rawPath)
	if err != nil {
		return "", err
	}
	return abs, nil
}

func OpenFromPath(rawPath string) (f *GoFile, err error) {
	path, err := ResolveGofilePath(rawPath)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return NewFromFile(file)

}
func NewFromFile(file *os.File) (f *GoFile, err error) {
	bytes, ioErr := ioutil.ReadAll(file)
	if ioErr != nil {
		return nil, errors.Wrapf(ioErr, "failed while reading bytes from file. path=%s", file.Name())
	}
	gofileStr := string(bytes)

	fileScanner := bufio.NewScanner(strings.NewReader(gofileStr))
	var installables []Installable

	for fileScanner.Scan() {
		line := fileScanner.Text()
		rawUrl := line
		i, iErr := NewInstallable(rawUrl)
		if iErr != nil {
			return f, errors.Wrapf(ioErr, "failed while parsing line from Gofile.Example\nLine=%s\nGoFile=%s", line, gofileStr)
		}
		installables = append(installables, *i)
	}

	return New(installables...)
}

func NewInstallable(rawUrl string) (i *Installable, err error) {
	uri, err := url.Parse(rawUrl)
	if err != nil {
		return i, errors.Wrapf(errors.Wrapf(err, "failed while parsing rawUrl\nURL=%s", rawUrl), "failed while creating a new Installable from rawUrl")
	}

	segments := strings.Split(uri.Path, "/")
	i = &Installable{
		host:   segments[0],
		author: segments[1],
		Name:   segments[2],
		URI:    rawUrl,
	}
	return i, err
}

func (g GoFile) TablePrint() (string, error) {
	sb := &strings.Builder{}
	err := g.TablePrintW(sb)
	return sb.String(), err
}

func (g GoFile) PlainPrint() (str string, err error) {
	sb := &strings.Builder{}
	err = g.PlainPrintW(sb)
	return sb.String(), err
}

func (g GoFile) PlainPrintW(w io.Writer) (err error) {
	buffer := &bytes.Buffer{}
	for _, installable := range g.Installables {
		if _, err = buffer.WriteString(installable.URI + "\n"); err != nil {
			return err
		}
	}
	_, err = buffer.WriteTo(w)
	return err
}

func (g GoFile) TablePrintW(w io.Writer) (err error) {
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
