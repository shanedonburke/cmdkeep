package cmdkeep

import (
	"cmdkeep/cli"
	"cmdkeep/model"
	"fmt"
)

type RemoveCommand struct{}

func (rm *RemoveCommand) Run(cl *cli.CLI, m *model.Model) {
	config := cl.RM
	keys := config.Keys

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
