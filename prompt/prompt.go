package prompt

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func ConfirmOrExit(prompt string) {
	termState := makeTermRaw()

	fmt.Print(prompt)

	confirmed := readYesNo()

	if confirmed {
		restoreTermState(termState)
		fmt.Print("\n")
	} else {
		os.Exit(0)
	}
}

func readYesNo() bool {
	buf := make([]byte, 1)
	if _, err := os.Stdin.Read(buf); err != nil {
		confirmationError(err)
	}
	answer := strings.ToLower(string(buf[0]))
	return answer == "y"
}

func confirmationError(err error) {
	fmt.Fprintf(os.Stderr, "Error: Confirmation failed: %v\n", err)
	os.Exit(1)
}

func makeTermRaw() *term.State {
	state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		confirmationError(err)
	}
	return state
}

func restoreTermState(state *term.State) {
	term.Restore(int(os.Stdin.Fd()), state)
}
