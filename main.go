package main

import (
	"github.com/alex-held/gofile/internal/cmd"
)


// completions:
// bash: eval "$(your-cli-tool --completion-script-bash)"
// zsh: eval "$(your-cli-tool --completion-script-zsh)"
func main() {
	cli := cmd.New()
	cli.Run()
}
