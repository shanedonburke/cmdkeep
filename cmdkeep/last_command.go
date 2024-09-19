package cmdkeep

import (
	"cmdkeep/cli"
	"cmdkeep/model"
	"cmdkeep/runner"
	"fmt"
	"os"
	"strings"
)

type LastCommand struct{}

func (lc *LastCommand) Run(cl *cli.CLI, m *model.Model) {
	config := cl.Last
	lastCommand := m.Last

	var mode runner.ExecMode = runner.Prompt
	if config.PrintOnly {
		mode = runner.Print
	} else if config.NoPrompt {
		mode = runner.Execute
	}

	if lastCommand == "" {
		fmt.Fprint(os.Stderr, "Error: No commands have been executed - try `ck run`\n")
		os.Exit(1)
	} else if strings.HasPrefix(lastCommand, "key:") {
		runner.NewRunner().RunKey(m, strings.Split(lastCommand, ":")[1], config.Args, config.UseDefaults, mode)
	} else if strings.HasPrefix(lastCommand, "template:") {
		runner.NewRunner().RunTemplate(m, strings.Split(lastCommand, ":")[1], config.Args, config.UseDefaults, mode)
	} else {
		fmt.Fprintf(os.Stderr, "Error: Invalid last command: %s", lastCommand)
		os.Exit(1)
	}
}
