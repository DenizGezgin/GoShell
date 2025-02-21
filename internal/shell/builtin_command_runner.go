package shell

import (
	"fmt"
	"io"
)

// CommandRunner defines the interface for executing shell commands
type CommandRunner interface {
	Run(args []string)
}

type builtinCommandRunner struct {
	Stderr      io.Writer
	Stdout      io.Writer
	commandName string
}

func NewBuiltinCommandRunner(commandName string, stdout, stderr io.Writer) *builtinCommandRunner {
	return &builtinCommandRunner{
		commandName: commandName,
		Stdout:      stdout,
		Stderr:      stderr,
	}
}

func (b *builtinCommandRunner) WriteToStd(result string, err error) {
	if err != nil {
		b.Stderr.Write([]byte(fmt.Sprintf("%s: %v\r\n", b.commandName, err)))
		return
	}
	b.Stdout.Write([]byte(fmt.Sprintf("%s\r\n", result)))
}

func (b *builtinCommandRunner) validateAtLeastOneArgs(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	return nil
}
