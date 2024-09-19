package cmdkeep

import (
	"cmdkeep/cli"
	"cmdkeep/model"
	"fmt"
	"os"
	"slices"
	"strings"
)

type AddCommand struct{}

func (ac *AddCommand) Run(cl *cli.CLI, m *model.Model) {
	config := cl.Add
	key := config.Key

	if slices.Contains(cli.CommandStrings, strings.ToLower(key)) {
		fmt.Fprintf(os.Stderr, "Error: Key '%s' is reserved - try a different name.\n", key)
		os.Exit(1)
	}

	command := model.NewCommand(config.Command)
	m.AddCommand(key, command)
	model.WriteModel(m)
	fmt.Printf("Added command: %s\n", key)
}
