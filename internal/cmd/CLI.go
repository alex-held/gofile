package cli

import (
	"os"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"
)

func getGofilePath(v *string) (path string, err error)  {
	abs, err := filepath.Abs(*v)
	if err != nil {
		return "", err
	}
	return abs, nil
}


func (cli *CLI) Run() (err error) {
	switch kingpin.MustParse(cli.Parse(os.Args[1:])) {

	// list
	case cli.ListCmd.FullCommand():
		f, err := openGofile()
		if err != nil {
			return err
		}

		gofile, err := NewFromFile(f)
		if err != nil {
			return err

		}
		err = gofile.PrettyPrintW(os.Stdout)
		cli.FatalIfError(err, "[ERROR]\t%v\n", err)
	default:
		cli.FatalUsage("no command")
	}
}

type CLI struct {
	*kingpin.Application
	ListCmd         *ListCommand
	Version, Author string
	GlobalFlags     *GlobalFlags
}
