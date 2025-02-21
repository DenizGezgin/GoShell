package shell

import "strings"

type EchoCommandRunner struct {
	builtinCommandRunner
}

func (e *EchoCommandRunner) Run(args []string) {
	e.WriteToStd(strings.Join(args, " "), nil)
}
