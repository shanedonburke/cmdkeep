package driver

import "cmdkeep/cmdkeep"

type Driver struct{}

func NewDriver() *Driver {
	return &Driver{}
}

func (d *Driver) Run(command cmdkeep.CKCommand, cli *cmdkeep.CLI) {
	command.Run(cli)
}
