package runner

import (
	"bufio"
	"cmdkeep/model"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/google/shlex"
	"golang.org/x/term"
)

type RunMode int

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

func (r *Runner) RunKey(m *model.Model, key string, cliArgs []string, mode RunMode) {
	command := m.Commands[key]
	if command == nil {
		fmt.Fprintf(os.Stderr, "No such command: %s\n", key)
		os.Exit(1)
	}

	if mode == Prompt {
		r.promptForKey(key)
	} else {
		fmt.Printf("Using last command: '%s'\n", key)
	}

	outCommand, exitCode := r.runCommand(command, cliArgs, mode)
	m.AddCommand(key, outCommand)
	m.Last = "key:" + key
	m.LastArgs = outCommand.LastArgs
	model.WriteModel(m)
	os.Exit(exitCode)
}

func (r *Runner) RunTemplate(m *model.Model, template string, cliArgs []string, mode RunMode) {
	command := model.NewCommand(template)

	if m.Last == fmt.Sprintf("template:%s", template) {
		command.LastArgs = m.LastArgs
	}

	if mode == Prompt {
		r.promptForTemplate(template)
	} else {
		fmt.Printf("Using last command: `%s`\n", template)
	}

	outCommand, exitCode := r.runCommand(command, cliArgs, mode)
	m.Last = "template:" + template
	m.LastArgs = outCommand.LastArgs
	model.WriteModel(m)
	os.Exit(exitCode)
}

func (r *Runner) promptForKey(key string) {
	r.confirmOrExit(fmt.Sprintf("Run command '%s'? (y or n): ", key))
}

func (r *Runner) promptForTemplate(template string) {
	r.confirmOrExit(fmt.Sprintf("Run command `%s`? (y or n): ", template))
}

func (r *Runner) confirmOrExit(prompt string) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		r.confirmationError(err)
	}

	fmt.Print(prompt)

	buf := make([]byte, 1)
	if _, err = os.Stdin.Read(buf); err != nil {
		r.confirmationError(err)
	}
	answer := strings.ToLower(string(buf[0]))

	if answer == "y" {
		term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Print("\n")
	} else {
		os.Exit(0)
	}
}

func (r *Runner) confirmationError(err error) {
	fmt.Fprintf(os.Stderr, "Command confirmation failed: %v\n", err)
	os.Exit(1)
}

func (r *Runner) runCommand(command *model.Command, cliArgs []string, mode RunMode) (*model.Command, int) {
	procTemplate := r.processTemplate(command, cliArgs)

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
		fmt.Fprintf(os.Stderr, "Failed to parse command: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	exitCode := 0
	if err := cmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			exitCode = exiterr.ExitCode()
		} else {
			fmt.Fprintf(os.Stderr, "Command execution failed: %v\n", err)
			exitCode = 1
		}
	}
	return outCommand, exitCode
}

func (r *Runner) processTemplate(command *model.Command, cliArgs []string) *ProcessedTemplate {
	template := command.Template
	lastArgs := command.LastArgs

	args := []string{}
	argIdx := 0

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
			varName := ""
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
							varName += string(nameNext)
							i++
						} else {
							// The next char is not escapable,so just put a backslash
							cmdString += string(nameCh)
						}
					} else {
						// Backslash was the last character in the template
						varName += string(nameCh)
					}
				} else {
					varName += string(nameCh)
				}

				i++
				nameCh, ok = getCharByIndex(template, i)
			}
			if !ok {
				// The template ended without a closing brace.
				// Add the opening brace and subsequent chars to command.
				cmdString += "{" + varName
				break
			}
			// Reached closing brace

			if varName == "" {
				// Variable was '{}' - name it according to its index
				varName = fmt.Sprintf("%d", argIdx+1)
			}

			lastValue := ""
			if argIdx < len(lastArgs) {
				// An argument has been provided for this var before
				lastValue = lastArgs[argIdx]
			}

			if !hasPrintedHelp && argIdx >= len(cliArgs) {
				fmt.Println("Enter command arguments: (Enter/Return for default if shown, '-' for blank)")
				hasPrintedHelp = true
			}

			var argValue string
			if argIdx < len(cliArgs) {
				// An argument for this var was specified through the CLI
				argValue = cliArgs[argIdx]
			} else {
				argValue = r.promptForValue(varName, lastValue)
			}
			args = append(args, argValue)
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

func (r *Runner) promptForValue(varName string, def string) string {
	if def != "" {
		fmt.Printf("- %s [%s]: ", varName, def)
	} else {
		fmt.Printf("- %s: ", varName)
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Value prompt failed: %v\n", err)
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
