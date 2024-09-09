package cmdkeep

import (
	"cmdkeep/model"
	"cmdkeep/runner"
)

type RunCommand struct{}

func (rc *RunCommand) Run(cli *CLI) {
	config := cli.Run
	m := model.ReadModel()
	r := runner.NewRunner()

	var mode runner.RunMode = runner.Execute
	if config.PrintOnly {
		mode = runner.Print
	}

	if config.Key != "" {
		r.RunKey(m, config.Key, config.Args, mode)
	} else {
		r.RunTemplate(m, config.Command, config.Args, mode)
	}
}
