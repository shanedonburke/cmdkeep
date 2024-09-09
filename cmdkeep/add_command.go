package cmdkeep

import (
	"cmdkeep/lev"
	"cmdkeep/model"
	"fmt"
	"os"
	"slices"
)

type AddCommand struct{}

func (ac *AddCommand) Run(cli *CLI) {
	config := cli.Add
	key := config.Key

	if slices.Contains(lev.GetReservedWords(), key) {
		fmt.Fprintf(os.Stderr, "Key '%s' is reserved - try a different name.\n", key)
		os.Exit(1)
	}

	m := model.ReadModel()
	command := model.NewCommand(config.Command)
	m.AddCommand(key, command)
	model.WriteModel(m)
	fmt.Printf("Added command: %s\n", key)
}
