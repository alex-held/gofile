package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/alex-held/gofile/internal/log"
)

func init() {
	Globals = GlobalFlags{
		Verbose:  new(bool),
		LogLevel: new(string),
	}
}

type Registerer interface {
	Register(cmdClause *kingpin.CmdClause)
}

var Globals GlobalFlags

type GlobalFlags struct {
	Verbose  *bool
	LogLevel *string
}

type CLI struct {
	*kingpin.Application
	GlobalFlags *GlobalFlags
}

func (cli *CLI) Register(name, help string) *kingpin.CmdClause {
	return cli.Application.Command(name, help)
}

func (cli *CLI) Run() {
	kingpin.MustParse(cli.Parse(os.Args[1:]))
}

func (cli *CLI) LogF(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (cli *CLI) ConfigureGlobals() {
	Globals.LogLevel = cli.Flag("level", "loglevel").
		Short('l').
		Default("info").
		Action(func(ctx *kingpin.ParseContext) error {
			switch *Globals.LogLevel {
			case "vv":
				log.Log = log.Log.WithDebug()
				log.Logger.SetLevel(logrus.TraceLevel)
			case "v":
				log.Log = log.Log.WithDebug()
				log.Logger.SetLevel(logrus.DebugLevel)
			case "i":
				log.Log = log.Log.WithoutDebug()
				log.Logger.SetLevel(logrus.InfoLevel)
			case "w":
				log.Log = log.Log.WithoutDebug()
				log.Logger.SetLevel(logrus.WarnLevel)
			case "e":
				log.Log = log.Log.WithoutDebug()
				log.Logger.SetLevel(logrus.ErrorLevel)
			case "q":
				log.Log = log.Log.WithoutDebug()
				log.Logger.SetLevel(logrus.PanicLevel)
			}
			return nil
		}).
		String()

	Globals.Verbose = cli.Flag("verbose", "verbose mode.").
		Short('v').
		Default("false").
		Action(func(_ *kingpin.ParseContext) error {
			log.Logger.SetLevel(logrus.DebugLevel)
			return nil
		}).
		Bool()
}

func (cli *CLI) ConfigureCommands() {
	ConfigureFileCommand(cli)
	ConfigureInstallCommand(cli)
}

func New() (cli *CLI) {
	cli = &CLI{
		Application: kingpin.New("gofile", "A command line application to install go binaries using a config file."),
	}
	cli.Author("Alexander Held")
	cli.Version("v0.0.1")

	cli.ConfigureGlobals()
	cli.ConfigureCommands()

	return cli
}
