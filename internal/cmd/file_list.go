package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/alex-held/gofile/internal/gofile"
)

type FileListCmd struct {
	File   string
	Format string
	vars   struct {
		filePtr   *string
		formatPtr *string
	}
}

func (cmd *FileListCmd) run(ctx *kingpin.ParseContext) (err error) {
	f, err := gofile.OpenFromPath(cmd.File)
	if err != nil {
		return err
	}

	return printWithFormat(cmd.Format, f)
}

func printWithFormat(format string, f *gofile.GoFile) (err error) {
	switch format {
	case "plain":
		if err = f.PlainPrintW(os.Stdout); err != nil {
			return err
		}
	case "table":
		if err = f.TablePrintW(os.Stdout); err != nil {
			return err
		}
	}
	return err
}

func ConfigureFileListCommand(fileCmd *kingpin.CmdClause) {
	const fileArg = "file"
	c := &FileListCmd{
		vars: struct {
			filePtr   *string
			formatPtr *string
		}{filePtr: new(string), formatPtr: new(string)},
	}

	cmd := fileCmd.Command("list", "lists the go binaries of the file")
	cmd.Flag("format", "specify the output format").
		Default("table").
		EnumVar(c.vars.formatPtr, "plain", "table")

	cmd.Flag(fileArg, "path of the Gofile").
		Default("Gofile").
		Short('f').
		StringVar(c.vars.filePtr)

	cmd.PreAction(func(ctx *kingpin.ParseContext) (err error) {
		c.File, err = gofile.ResolveGofilePath(*c.vars.filePtr)
		c.Format = *c.vars.formatPtr
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
