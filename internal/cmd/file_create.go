package cmd

import (
	"bytes"
	"os"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/alex-held/gofile/internal/log"
	"github.com/alex-held/gofile/internal/search"
)

type FileCreateCmd struct {
	Out  string
	Dir  string
	Dry  bool
	vars struct {
		dirPtr *string
		outPtr *string
		dryPtr *bool
	}
}

var globalFlags *GlobalFlags

func (cmd *FileCreateCmd) run(ctx *kingpin.ParseContext) (err error) {
	pkgNames, err := cmd.getPackages(err)
	if err != nil {
		return err
	}
	sb := &bytes.Buffer{}
	count := 0

	searchResults := search.SearchPackages( pkgNames...)
	for _, result := range searchResults {
		if result.Err != nil || result.Repo == "" {
			log.Log.Error(err)
			continue
		}

		log.Log.Debugf("found url '%s' for pkg '%s'\n", result.Repo, result.Package)

		count++
		sb.WriteString(result.Repo + "\n")
	}

	sb.WriteString("\n")
	content := sb.String()

	if cmd.Dry {
		log.Log.Infof("would create Gofile at '%s' with content:\n---\n%s\n---\n", cmd.Out, content)
		return nil
	}

	outFile, err := os.OpenFile(cmd.Out, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}

	defer outFile.Close()
	_, err = sb.WriteTo(outFile)
	log.Log.Infof("found %d/%d packages\n", count, len(pkgNames))
	return err
}

func (cmd *FileCreateCmd) getPackages(err error) ([]string, error) {
	if cmd.Dir == "" {
		cmd.Dir = os.Getenv("GOPATH") + "/bin"
	}

	entries, err := os.ReadDir(cmd.Dir)
	if err != nil {
		return nil, err
	}

	var pkgNames []string
	for _, dirE := range entries {
		if dirE.IsDir() {
			continue
		}
		log.Log.Debugf("found package '%s'\n",dirE.Name())
		pkgNames = append(pkgNames, dirE.Name())
	}

	log.Log.Debugf("found %d packages in '$GOBIN' directory '%s'\n", len(pkgNames), cmd.Dir)
	return pkgNames, nil
}

func ConfigureFileCreateCommand(fileCmd *kingpin.CmdClause) {
	cmd := fileCmd.Command("create", "creates a Gofile")

	c := &FileCreateCmd{
		vars: struct {
			dirPtr *string
			outPtr *string
			dryPtr *bool
		}{dirPtr: new(string), outPtr: new(string), dryPtr: new(bool)},
	}

	cmd.Flag("out", "output a Gofile at this path").
		Default("Gofile").
		Short('o').
		StringVar(c.vars.outPtr)

	cmd.Flag("dir", "path of the bin dir").
		Default("").
		Short('d').
		StringVar(c.vars.dirPtr)

	cmd.Flag("dry", "enables / disables dry mode").
		Default("false").
		Short('n').
		BoolVar(c.vars.dryPtr)

	cmd.PreAction(func(ctx *kingpin.ParseContext) (err error) {
		c.Out, err = filepath.Abs(*c.vars.outPtr)
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

	cmd.Action(c.run)
}
