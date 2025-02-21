package shell

import "os"

type ExitCommandRunner struct {
	builtinCommandRunner
}

func (e *ExitCommandRunner) Run(args []string) {
	os.Exit(0)
}
