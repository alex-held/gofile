package cmd

import (
	"fmt"
	"os"
	path2 "path"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/alex-held/gofile/internal/gofile"
)

type FileInitCmd struct {
	Out  string
	vars struct {
		outPtr *string
	}
}

func (cmd *FileInitCmd) run(ctx *kingpin.ParseContext) (err error) {
	f, err := gofile.NewFromFile(os.Stdin)
	if err != nil {
		if f, err = gofile.New(); err != nil {
			return err
		}
	}

	path := path2.Join(cmd.Out, "Gofile")
	_, err = gofile.SaveFile(path, f)
	return err
}

func ConfigureFileInitCommand(fileCmd *kingpin.CmdClause) {
	const outArg = "out"
	c := &FileInitCmd{
		vars: struct{ outPtr *string }{outPtr: new(string)},
	}

	cmd := fileCmd.Command("init", "Inits a new `Gofile`.")

	cmd.Flag(outArg, "output path of the `Gofile`").
		Default("Gofile").
		Short('o').
		StringVar(c.vars.outPtr)

	cmd.PreAction(func(ctx *kingpin.ParseContext) (err error) {
		c.Out, err = gofile.ResolveGofilePath(*c.vars.outPtr)
		return err
	})

	cmd.Validate(func(clause *kingpin.CmdClause) error {
		v := NewValidator().
			AddFlagPredicateFn("out", func(model *kingpin.FlagModel, value kingpin.Value) (err error) {
				absolute, err := filepath.Abs(value.String())
				if err != nil {
					return fmt.Errorf("argument 'out' is invalid.\treason=could not resolve absolute path.;\terr=%v;\n", err)
				}
				if _, err = os.Stat(absolute); err != nil {
					return fmt.Errorf("argument 'out' is invalid. the resolved directory ':%s' does not exist. out=%s\n", absolute, value.String())
				}
				return nil
			})
		return v.Validate(clause)
	})

	cmd.Action(c.run)
}
