package cmd

import "C"
import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/alex-held/gofile/internal/gofile"
)

/*
var (
	// list
	 listCmd    = app.Command("list", "lists the go binaries of the file")
	gofilepath = listCmd.Arg("file", "path of the Gofile.Example.").
		Required().
		String()

	dir = kingpin.Flag("directory", "install directory").
		Short('d').
		Default("$HOME").
		ExistingDir()
)*/

type ListCommand struct {
	*kingpin.CmdClause
	GofilePath *string
}

type ListCmd struct {
	File string
	vars struct {
		filePtr *string
	}
}

func (cmd *ListCmd) run(ctx *kingpin.ParseContext) (err error) {

	f, err := gofile.OpenFromPath(cmd.File)
	if err != nil {
		return err
	}
	err = f.PrettyPrintW(os.Stderr)
	return err
}

func configureListCommand(app *CLI) {
	const fileArg = "file"
	c := &ListCmd{
		vars: struct{ filePtr *string }{filePtr: new(string)},
	}

	cmd := app.Command("list", "lists the go binaries of the file")

	cmd.Arg(fileArg, "path of the Gofile").
		Default("Gofile").
		StringVar(c.vars.filePtr)

	cmd.PreAction(func(ctx *kingpin.ParseContext) (err error) {
		c.File, err = gofile.ResolveGofilePath(*c.vars.filePtr)
		return err
	})

	cmd.Validate(func(clause *kingpin.CmdClause) error {
		v := &argValidator{fns: map[string]ArgPredicateFn{
			"file": func(model *kingpin.ArgModel, value kingpin.Value) (err error) {
				absolute, err := filepath.Abs(value.String())
				if err != nil {
					return fmt.Errorf("argument 'file' is invalid.\treason=could not resolve absolute path.;\terr=%v;\n", err)
				}
				if _, err = os.Stat(absolute); err != nil {
					return fmt.Errorf("argument 'file' is invalid. the resolved file does not exist. file=%s; resolved:%s;\n", value.String(), absolute)
				}
				return nil
			},
		},
		}
		return v.Validate(clause)
	})

	cmd.Action(c.run)
}

type ValidateArgFn func(clause *kingpin.ArgClause, name string, predicateFn func(model *kingpin.ArgModel) (err error)) (err error)
type ArgPredicateFn func(model *kingpin.ArgModel, value kingpin.Value) (err error)

type argValidator struct {
	fns map[string]ArgPredicateFn
}

func (v *argValidator) AddArgPredicateFn(name string, fn ArgPredicateFn) *argValidator {
	v.fns[name] = fn
	return v
}

func (v *argValidator) Validate(clause *kingpin.CmdClause) error {
	for _, arg := range clause.Model().Args {
		if argPredicateFn, ok := v.fns[arg.Name]; ok {
			return argPredicateFn(arg, arg.Value)
		}
	}
	return nil
}

func ValidateFileArg() ValidateArgFn {
	return func(clause *kingpin.ArgClause, name string, predicateFn func(*kingpin.ArgModel) (err error)) (err error) {
		m := clause.Model()
		if m.Name != name {
			return nil
		}
		return predicateFn(m)
	}
}
