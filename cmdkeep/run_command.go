package cmdkeep

import (
	"cmdkeep/cli"
	"cmdkeep/model"
	"cmdkeep/runner"
	"fmt"
	"os"
	"strings"
)

type RunCommand struct{}

func (rc *RunCommand) Run(cl *cli.CLI, m *model.Model) {
	config := cl.Run

	if config.Key == "" && strings.TrimSpace(config.Command) == "" {
		fmt.Fprintln(os.Stderr, "Error: `ck run` must specify key or command - try `ck run -h`")
		os.Exit(1)
	}

	r := runner.NewRunner()

	var mode runner.ExecMode = runner.Execute
	if config.PrintOnly {
		mode = runner.Print
	}

	if config.Key != "" {
		r.RunKey(m, config.Key, config.Args, config.UseDefaults, mode)
	} else {
		r.RunTemplate(m, config.Command, config.Args, config.UseDefaults, mode)
	}
}
