package shell

type PrintWorkingDirectoryCommandRunner struct {
	builtinCommandRunner
}

func (e *PrintWorkingDirectoryCommandRunner) Run(args []string) {
	e.WriteToStd(GetShell().GetWorkingDirectory(), nil)
}
