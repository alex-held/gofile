package cmd

/*
import "C"
import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/alex-held/gofile/internal/gofile"
)

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
	err = f.TablePrintW(os.Stderr)
	return err
}

func configureListCommand(app *CLI) {
	const fileArg = "file"
	c := &ListCmd{
		vars: struct{ filePtr *string }{filePtr: new(string)},
	}

	cmd := app.Command("list", "lists the go binaries of the file")

	cmd.Flag(fileArg, "path of the Gofile").
		Default("Gofile").
		Short('f').
		StringVar(c.vars.filePtr)

	cmd.PreAction(func(ctx *kingpin.ParseContext) (err error) {
		c.File, err = gofile.ResolveGofilePath(*c.vars.filePtr)
		return err
	})

	cmd.Validate(func(clause *kingpin.CmdClause) error {
		v := NewValidator().
			AddFlagPredicateFn("file", func(model *kingpin.FlagModel, value kingpin.Value) (err error) {
				absolute, err := filepath.Abs(value.String())
				if err != nil {
					return fmt.Errorf("argument 'file' is invalid.\treason=could not resolve absolute path.;\terr=%v;\n", err)
				}
				if _, err = os.Stat(absolute); err != nil {
					return fmt.Errorf("argument 'file' is invalid. the resolved file does not exist. file=%s; resolved:%s;\n", value.String(), absolute)
				}
				return nil
			}).
			AddFlagPredicateFn("dir", func(model *kingpin.FlagModel, value kingpin.Value) (err error) {
				absolute, err := filepath.Abs(value.String())
				if err != nil {
					return fmt.Errorf("argument 'dir' is invalid.\treason=could not resolve absolute path.;\terr=%v;\n", err)
				}
				if _, err = os.Stat(absolute); err != nil {
					return fmt.Errorf("argument 'dir' is invalid. the resolved directory does not exist. file=%s; resolved:%s;\n", value.String(), absolute)
				}
				return nil
			})

		return v.Validate(clause)
	})

	cmd.Action(c.run)
}

 */
