package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/alex-held/gofile/internal/gofile"
)

type InstallCmd struct {
	File string
	Dir  string
	Dry  bool
	vars struct {
		filePtr *string
		dirPtr  *string
		dryPtr  *bool
	}
}

func (cmd *InstallCmd) run(ctx *kingpin.ParseContext) (err error) {
	f, err := gofile.OpenFromPath(cmd.File)
	if err != nil {
		return err
	}

	if cmd.Dir != "" {
		defaultDir := os.Getenv("GOPATH") + "/bin"
		cmd.Dir = defaultDir
	}

	for _, installable := range f.Installables {
		command := exec.Command("go", "get", installable.URI)
		command.Env = os.Environ()
		command.Env = append(command.Env, "GOBIN="+cmd.Dir)

		if cmd.Dry {
			fmt.Printf("[DRY] would execute '%s'\n", command.String())
			continue
		}
		if err = command.Run(); err != nil {
			return err
		}

	}
	return err
}

func ConfigureInstallCommand(app *CLI) {
	c := &InstallCmd{
		vars: struct {
			filePtr *string
			dirPtr  *string
			dryPtr  *bool
		}{filePtr: new(string), dirPtr: new(string), dryPtr: new(bool)},
	}

	cmd := app.Command("install", "installs the go binaries of the file")

	cmd.Flag("file", "path of the Gofile").
		Default("Gofile").
		Short('f').
		StringVar(c.vars.filePtr)

	cmd.Flag("dir", "path of the bin dir").
		Default("").
		Short('d').
		StringVar(c.vars.dirPtr)

	cmd.Flag("dry", "enables / disables dry mode").
		Default("false").
		Short('n').
		BoolVar(c.vars.dryPtr)

	cmd.PreAction(func(ctx *kingpin.ParseContext) (err error) {
		c.File, err = gofile.ResolveGofilePath(*c.vars.filePtr)
		if err != nil {
			return err
		}
		c.Dir, err = filepath.Abs(*c.vars.dirPtr)
		if err != nil {
			return err
		}
		c.Dry = *c.vars.dryPtr
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
