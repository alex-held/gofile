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
	return val.TablePrintW(w)
}

type PrettyPrintable interface {
	TablePrint() (str string, err error)
	TablePrintW(w io.Writer) (err error)
}
