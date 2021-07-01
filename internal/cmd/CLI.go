package cmd

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

type GlobalFlags struct {
	Verbose *bool
}

type CLI struct {
	*kingpin.Application
	GlobalFlags *GlobalFlags
}

func (cli *CLI) Run() {
	kingpin.MustParse(cli.Parse(os.Args[1:]))
}

func (cli *CLI) LogF(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (cli *CLI) ConfigureGlobals() {
	cli.GlobalFlags.Verbose = cli.Flag("verbose", "verbose mode.").
		Short('v').
		Default("false").
		Bool()
}

func (cli *CLI) ConfigureCommands() {
	configureListCommand(cli)
}

func New() (cli *CLI) {
	cli = &CLI{
		Application: kingpin.New("gofile", "A command line application to install go binaries using a config file."),
		GlobalFlags: &GlobalFlags{
			Verbose: new(bool),
		},
	}

	cli.Author("Alexander Held")
	cli.Version("v0.0.1")

	cli.ConfigureGlobals()
	cli.ConfigureCommands()

	return cli
}
