package shell

import (
	"io"
)

type CommandRunnerBuilder struct{}

func (b *CommandRunnerBuilder) Build(name string, stdout, stderr io.Writer) (CommandRunner, error) {
	builtinRunner := NewBuiltinCommandRunner(name, stdout, stderr)
	switch name {
	case "echo":
		return &EchoCommandRunner{*builtinRunner}, nil
	case "exit":
		return &ExitCommandRunner{*builtinRunner}, nil
	case "type":
		return &TypeCommandRunner{*builtinRunner}, nil
	case "pwd":
		return &PrintWorkingDirectoryCommandRunner{*builtinRunner}, nil
	case "cd":
		return &ChangeDirectoryCommandRunner{*builtinRunner}, nil
	default:
		return NewExternalCommandRunner(name, stdout, stderr)
	}
}
