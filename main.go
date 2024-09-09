package main

import (
	"cmdkeep/cmdkeep"
	"cmdkeep/driver"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

func main() {
	var cli cmdkeep.CLI
	ctx := kong.Parse(&cli,
		kong.Name("ck"),
		kong.Description("A tool for saving and reusing shell commands."))

	d := driver.NewDriver()
	switch ctx.Command() {
	case "add <key> <command>":
		d.Run(&cmdkeep.AddCommand{}, &cli)
	case "run":
		if cli.Run.Command == "" {
			fmt.Fprintln(os.Stderr, "Error: no command specified - try `ck run -h`")
			os.Exit(1)
		}
		d.Run(&cmdkeep.RunCommand{}, &cli)
	case "run <key>":
		d.Run(&cmdkeep.RunCommand{}, &cli)
	case "last":
		d.Run(&cmdkeep.LastCommand{}, &cli)
	case "commands":
		d.Run(&cmdkeep.ListCommand{}, &cli)
	case "rm <keys>":
		d.Run(&cmdkeep.RemoveCommand{}, &cli)
	default:
		fmt.Fprintf(os.Stderr, "Unknown CLI command: %s\n", ctx.Command())
		os.Exit(1)
	}
}
