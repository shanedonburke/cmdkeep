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

	keySpecified := config.Key != ""
	commandSpecified := strings.TrimSpace(config.Command) != ""

	if keySpecified && commandSpecified {
		fmt.Fprintln(os.Stderr, "Error: Cannot specify both a key and `--command` - try `ck run -h`")
		os.Exit(1)
	} else if !keySpecified && !commandSpecified {
		fmt.Fprintln(os.Stderr, "Error: `ck run` must specify key or command - try `ck run -h`")
		os.Exit(1)
	}

	r := runner.NewRunner()

	var mode runner.ExecMode = runner.Execute
	if config.PrintOnly {
		mode = runner.Print
	}

	if keySpecified {
		r.RunKey(m, config.Key, config.Args, config.UseDefaults, mode)
	} else {
		r.RunTemplate(m, config.Command, config.Args, config.UseDefaults, mode)
	}
}
