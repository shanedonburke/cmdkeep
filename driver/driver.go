package driver

import (
	"cmdkeep/cli"
	"cmdkeep/cmdkeep"
	"cmdkeep/model"
	"cmdkeep/prompt"
	"cmdkeep/suggest"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/alecthomas/kong"
)

type Driver struct{}

func NewDriver() *Driver {
	return &Driver{}
}

func (d *Driver) Run() {
	m := model.ReadModel()
	d.preCLIParse(m)

	var cl cli.CLI
	ctx := kong.Parse(&cl,
		kong.Name("ck"),
		kong.Description("A tool for saving and reusing shell commands."))

	cmdStr := strings.Split(ctx.Command(), " ")[0]
	ckCmd := cmdkeep.GetCKCommand(cmdStr)

	if ckCmd != nil {
		ckCmd.Run(&cl, m)
	} else {
		fmt.Fprintf(os.Stderr, "Error: Unknown CLI command: %s\n", ctx.Command())
		os.Exit(1)
	}
}

func (d *Driver) preCLIParse(m *model.Model) {
	if ckCmdSuggestion := d.suggestCKCommand(m); ckCmdSuggestion != suggest.NoSuggestion {
		newProcArgs := []string{ckCmdSuggestion}

		if len(os.Args) > 2 {
			newProcArgs = append(newProcArgs, os.Args[2:]...)
		}

		newProcCmd := exec.Command(os.Args[0], newProcArgs...)
		newProcCmd.Stdin = os.Stdin
		newProcCmd.Stdout = os.Stdout
		newProcCmd.Stderr = os.Stderr

		exitCode := 0
		if err := newProcCmd.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			} else {
				exitCode = 1
			}
		}
		os.Exit(exitCode)
	}
}

func (d *Driver) suggestCKCommand(m *model.Model) string {
	suggestion := suggest.NoSuggestion

	if cmdStr, found := getSuggestableArg(); found {
		if d.isValidCommand(m, cmdStr) {
			return suggest.NoSuggestion
		}

		pool := d.getSuggestionStringPool(m)
		suggestion = suggest.FindClosestMatch(cmdStr, pool)

		if suggestion != suggest.NoSuggestion {
			fmt.Printf("Error: Unknown command '%s'\n", cmdStr)
			d.promptForSuggestion(suggestion)
		}
	}
	return suggestion
}

func (d *Driver) isValidCommand(m *model.Model, cmdStr string) bool {
	return slices.Contains(cli.CommandStrings, cmdStr) || m.Commands[cmdStr] != nil
}

func (d *Driver) getSuggestionStringPool(m *model.Model) []string {
	pool := cli.CommandStrings
	for existingKey := range maps.Keys(m.Commands) {
		pool = append(pool, existingKey)
	}
	return pool
}

func (d *Driver) promptForSuggestion(suggestion string) {
	var promptStr string

	if slices.Contains(cli.CommandStrings, suggestion) {
		promptStr = fmt.Sprintf("Did you mean 'ck %s'? (y/n): ", suggestion)
	} else {
		// User may be attempting to run a defined command key
		promptStr = fmt.Sprintf("Did you mean 'ck run %s'? (y/n): ", suggestion)
	}
	prompt.ConfirmOrExit(promptStr)
}

func getSuggestableArg() (string, bool) {
	for _, arg := range os.Args[1:] {
		if isCommandFlag(arg) {
			// If the user specified `--command`, we want to passthrough and
			// treat the execution as a `ck run` invocation.
			// Otherwise, `ck -c template` tries to suggest for `template`.
			return "", false
		} else if !strings.HasPrefix(arg, "-") {
			// Arg is not a flag = candidate for suggestion
			return arg, true
		}
	}
	// No argument for which a suggestion can be generated
	return "", false
}

func isCommandFlag(arg string) bool {
	return strings.HasPrefix(arg, "-c") || strings.HasPrefix(arg, "--command")
}
