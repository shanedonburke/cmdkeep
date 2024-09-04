package cmdkeep

import (
	"cmdkeep/model"
	"fmt"
)

type RemoveCommand struct{}

func (rm *RemoveCommand) Run(cli *CLI) {
	config := cli.RM
	keys := config.Keys

	m := model.ReadModel()
	for _, key := range keys {
		if _, ok := m.Commands[key]; ok {
			delete(m.Commands, key)
			model.WriteModel(m)
			fmt.Printf("Removed command: %s\n", key)
		} else {
			fmt.Printf("No such command: %s\n", key)
		}
	}
}
