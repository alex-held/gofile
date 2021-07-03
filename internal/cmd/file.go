package cmd

func ConfigureFileCommand(cli *CLI) {
	cmd := cli.Command("file", "manages Gofile")
	ConfigureFileCreateCommand(cmd)
	ConfigureFileListCommand(cmd)
	ConfigureFileInitCommand(cmd)
}
