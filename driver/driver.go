package driver

import "cmdkeep/cmdkeep"

const ENV_VAR = "__ck"

type Driver struct{}

func NewDriver() *Driver {
	return &Driver{}
}

func (d *Driver) Run(command cmdkeep.CKCommand, cli *cmdkeep.CLI) {
	command.Run(cli)
}
