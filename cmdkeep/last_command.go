package cmdkeep

import (
	"cmdkeep/model"
	"cmdkeep/runner"
	"fmt"
	"os"
	"strings"
)

type LastCommand struct{}

func (lc *LastCommand) Run(cli *CLI) {
	config := cli.Last
	m := model.ReadModel()
	lastCommand := m.Last

	if lastCommand == "" {
		fmt.Fprint(os.Stderr, "no commands have been executed - try `ck run`\n")
		os.Exit(1)
	} else if strings.HasPrefix(lastCommand, "key:") {
		runner.NewRunner().RunKey(m, strings.Split(lastCommand, ":")[1], config.Args, config.PrintCommand, !config.NoPrompt)
	} else if strings.HasPrefix(lastCommand, "template:") {
		runner.NewRunner().RunTemplate(m, strings.Split(lastCommand, ":")[1], config.Args, config.PrintCommand, !config.NoPrompt)
	} else {
		fmt.Fprintf(os.Stderr, "Invalid last command: %s", lastCommand)
		os.Exit(1)
	}
}
