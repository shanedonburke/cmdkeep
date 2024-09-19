package prompt

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func ConfirmOrExit(prompt string) {
	termState := makeTermRaw()

	// We need to capture SIGINT (Ctrl+C) so we can restore the terminal
	// to its non-raw state when triggered, i.e. before exiting.
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	cancelChan := make(chan bool, 1)

	go func() {
		for {
			select {
			case <-interruptChan:
				// SIGINT
				restoreTermState(termState)
				fmt.Println()
				os.Exit(0)
			case <-cancelChan:
				return
			}
		}
	}()

	defer func() {
		// Non-SIGINT case:
		// Restore terminal state and cancel SIGINT capture/goroutine
		restoreTermState(termState)
		signal.Stop(interruptChan)
		signal.Reset(syscall.SIGINT)
		cancelChan <- true
	}()

	fmt.Print(prompt)
	confirmed := readYesNo()

	if confirmed {
		// Continue execution
		fmt.Println()
	} else {
		restoreTermState(termState)
		fmt.Println()
		signal.Stop(interruptChan)
		signal.Reset(syscall.SIGINT)
		cancelChan <- true
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
