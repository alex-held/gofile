package main

import (
	"bufio"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

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
