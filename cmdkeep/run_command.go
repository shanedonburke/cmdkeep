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
	if config.Key != "" {
		r.RunKey(m, config.Key, config.Args, config.PrintCommand, false)
	} else {
		r.RunTemplate(m, config.Command, config.Args, config.PrintCommand, false)
	}
}
