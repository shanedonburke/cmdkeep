package cmdkeep

import (
	"cmdkeep/model"
	"fmt"
)

type AddCommand struct{}

func (ac *AddCommand) Run(cli *CLI) {
	config := cli.Add
	m := model.ReadModel()
	command := model.NewCommand(config.Command)
	m.AddCommand(config.Key, command)
	model.WriteModel(m)
	fmt.Printf("Added command: %s\n", config.Key)
}

func (ac *AddCommand) addPersistedCommand(key string, template string) {

}
