package print

import (
	"io"
	"strings"
)

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
