package runner

import (
	"bufio"
	"cmdkeep/model"
	"cmdkeep/prompt"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/google/shlex"
)

type ExecMode int

const (
	Execute = iota + 1
	Print
	Prompt
)

type Runner struct{}

func NewRunner() *Runner {
	return &Runner{}
}

type ProcessedTemplate struct {
	CmdString string
	Args      []string
}

func (r *Runner) RunKey(m *model.Model, key string, clArgs []string, useDefaults bool, mode ExecMode) {
	command := m.Commands[key]
	if command == nil {
		fmt.Fprintf(os.Stderr, "Error: Unknown command key: %s\n", key)
		os.Exit(1)
	}

	if mode == Prompt {
		r.promptForKey(key)
	}

	outCommand, exitCode := r.runCommand(command, clArgs, useDefaults, mode)
	m.AddCommand(key, outCommand)
	m.Last = "key:" + key
	m.LastArgs = outCommand.LastArgs
	model.WriteModel(m)
	os.Exit(exitCode)
}

func (r *Runner) RunTemplate(m *model.Model, template string, clArgs []string, useDefaults bool, mode ExecMode) {
	command := model.NewCommand(template)

	if m.Last == fmt.Sprintf("template:%s", template) {
		command.LastArgs = m.LastArgs
	}

	if mode == Prompt {
		r.promptForTemplate(template)
	}

	outCommand, exitCode := r.runCommand(command, clArgs, useDefaults, mode)
	m.Last = "template:" + template
	m.LastArgs = outCommand.LastArgs
	model.WriteModel(m)
	os.Exit(exitCode)
}

func (r *Runner) promptForKey(key string) {
	prompt.ConfirmOrExit(fmt.Sprintf("Run command '%s'? (y/n): ", key))
}

func (r *Runner) promptForTemplate(template string) {
	prompt.ConfirmOrExit(fmt.Sprintf("Run command `%s`? (y/n): ", template))
}

func (r *Runner) runCommand(command *model.Command, clArgs []string, useDefaults bool, mode ExecMode) (*model.Command, int) {
	procTemplate := r.processTemplate(command, clArgs, useDefaults)

	outCommand := &model.Command{
		Template: command.Template,
		LastArgs: procTemplate.Args,
	}

	if mode == Print {
		fmt.Println(procTemplate.CmdString)
		return outCommand, 0
	}
	// Mode is `Execute`

	cmdParts, err := shlex.Split(procTemplate.CmdString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to parse command: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	exitCode := 0
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			fmt.Fprintf(os.Stderr, "Error: Command execution failed: %v\n", err)
			exitCode = 1
		}
	}
	return outCommand, exitCode
}

func (r *Runner) processTemplate(command *model.Command, clArgs []string, useDefaults bool) *ProcessedTemplate {
	template := command.Template
	lastArgs := command.LastArgs

	args := []string{}
	argIdx := 0
	argsByName := map[string]string{}

	hasPrintedHelp := false
	cmdString := ""

	i := 0

	for i < len(template) {
		ch, _ := getCharByIndex(template, i)
		if ch == '\\' {
			// Character escape - check next char
			if next, ok := getCharByIndex(template, i+1); ok {
				if isEscapableChar(rune(next)) {
					// Append the escaped character and continue
					cmdString += string(next)
					i++
				} else {
					// The next char is not escapable,so just put a backslash
					cmdString += string(ch)
				}
			} else {
				// Backslash was the last character in the template
				cmdString += string(ch)
			}
		} else if ch == '{' {
			// Unescaped opening brace - parse template variable
			paramName := ""
			i++ // Skip opening brace
			nameCh, ok := getCharByIndex(template, i)

			if !ok {
				// Opening brace was the last character of the template.
				// Add it to the command.
				cmdString += string(ch)
				break
			}

			// Read characters of variable name
			for ok && nameCh != '}' {
				if nameCh == '\\' {
					// Character escape - check next char
					if nameNext, ok := getCharByIndex(template, i+1); ok {
						if isEscapableChar(rune(nameNext)) {
							// Append the escaped character and continue
							paramName += string(nameNext)
							i++
						} else {
							// The next char is not escapable,so just put a backslash
							cmdString += string(nameCh)
						}
					} else {
						// Backslash was the last character in the template
						paramName += string(nameCh)
					}
				} else {
					paramName += string(nameCh)
				}

				i++
				nameCh, ok = getCharByIndex(template, i)
			}
			if !ok {
				// The template ended without a closing brace.
				// Add the opening brace and subsequent chars to command.
				cmdString += "{" + paramName
				break
			}
			// Reached closing brace

			if paramName == "" {
				// Variable was '{}' - name it according to its index
				paramName = fmt.Sprintf("%d", argIdx+1)
			}

			hasLastArg := argIdx < len(lastArgs)
			lastArg := ""
			if hasLastArg {
				// An argument has been provided for this var before
				lastArg = lastArgs[argIdx]
			}

			var argValue string

			if argByName, ok := argsByName[paramName]; ok {
				// This param name matches a previous one - use previous value
				argValue = argByName
			} else if argIdx < len(clArgs) {
				// An argument for this param was specified through the CLI
				argValue = clArgs[argIdx]
			} else if hasLastArg && useDefaults {
				// `--use-defaults` was specified and this param has a default
				argValue = lastArg
			} else {
				// We don't have a value from any source
				if !hasPrintedHelp {
					fmt.Println("Enter command arguments: (Enter/Return for default if shown, '-' for blank)")
					hasPrintedHelp = true
				}

				argValue = r.promptForValue(paramName, lastArg)
			}
			args = append(args, argValue)
			argsByName[paramName] = argValue
			cmdString += argValue
			argIdx++
		} else {
			cmdString += string(ch)
		}
		i++
	}
	return &ProcessedTemplate{
		Args:      args,
		CmdString: cmdString,
	}
}

func (r *Runner) promptForValue(paramName string, def string) string {
	if def != "" {
		fmt.Printf("- %s [%s]: ", paramName, def)
	} else {
		fmt.Printf("- %s: ", paramName)
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Value prompt failed: %v\n", err)
		os.Exit(1)
	}
	input = strings.TrimSpace(input)

	switch input {
	case "":
		// User pressed Enter without entering anything
		return def
	case "-":
		// User indicator for blank input
		return ""
	default:
		return input
	}
}

func getCharByIndex(str string, idx int) (byte, bool) {
	if idx >= len(str) {
		return 0, false
	} else {
		return str[idx], true
	}
}

// Braces indicate the start/end of template variables,
// so '\\', '\{', and '\}' are provided as escapes.
func isEscapableChar(ch rune) bool {
	return ch == '\\' || ch == '{' || ch == '}'
}
