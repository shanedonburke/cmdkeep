package cmdkeep

import (
	"cmdkeep/cli"
	"cmdkeep/model"
)

type CKCommand interface {
	Run(cl *cli.CLI, m *model.Model)
}

func GetCKCommand(cmdStr string) CKCommand {
	switch cmdStr {
	case "add":
		return &AddCommand{}
	case "rm":
		return &RemoveCommand{}
	case "commands":
		return &ListCommand{}
	case "run":
		return &RunCommand{}
	case "last":
		return &LastCommand{}
	default:
		return nil
	}
}
