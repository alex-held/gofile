package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("gofile", "A command line application to install go binaries using a config file.")

	// global flags
	verbose = app.Flag("verbose", "verbose mode.").
		Short('v').
		Default("false").
		Bool()

	// list
	listCmd    = app.Command("list", "lists the go binaries of the file")
	gofilepath = listCmd.Arg("file", "path of the Gofile.Example.").
		Required().
		String()

	dir = kingpin.Flag("directory", "install directory").
		Short('d').
		Default("$HOME").
		ExistingDir()
)

// completions:
// bash: eval "$(your-cli-tool --completion-script-bash)"
// zsh: eval "$(your-cli-tool --completion-script-zsh)"
func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	// list
	case listCmd.FullCommand():
		abs, err := filepath.Abs(*gofilepath)
		if err != nil {
			handleErrAndExit(err)
			return
		}
		gofilepath = &abs

		f, err := os.Open(*gofilepath)
		if err != nil {
			handleErrAndExit(err)
			return
		}

		gofile, err := NewFromFile(f)
		if err != nil {
			handleErrAndExit(err)
			return

		}
		gofile.PrettyPrintW(os.Stdout)

	default:
		app.FatalUsage("no command")
	}
}

func handleErrAndExit(err error) {
	fmt.Printf("[ERROR]\t%v\n", err)
	os.Exit(1)
}
