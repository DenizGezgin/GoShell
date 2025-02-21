package shell

import (
	"fmt"
)

type TypeCommandRunner struct {
	builtinCommandRunner
}

func (t *TypeCommandRunner) Run(args []string) {
	err := t.validateAtLeastOneArgs(args)
	if err != nil {
		t.WriteToStd("", err)
		return
	}

	commandName := args[0]
	cmdPath := GetShell().GetCommandRepository().GetCommandPath(commandName)
	if cmdPath == "" {
		t.WriteToStd(fmt.Sprintf("%s: not found", commandName), nil)
		return
	}

	t.WriteToStd(fmt.Sprintf("%s is %s", commandName, cmdPath), nil)
}
