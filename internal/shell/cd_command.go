package shell

import (
	"fmt"
	"os"
	"path/filepath"
)

type ChangeDirectoryCommandRunner struct {
	builtinCommandRunner
}

func (c *ChangeDirectoryCommandRunner) Run(args []string) {
	err := c.validateAtLeastOneArgs(args)
	if err != nil {
		c.WriteToStd("", err)
		return
	}

	var targetDirAbsolutePath = getTargetDirAbsolutePath(args[0])
	fileInfo, err := os.Stat(targetDirAbsolutePath)
	if err != nil {
		c.WriteToStd("", fmt.Errorf("no such file or directory: %s", targetDirAbsolutePath))
		return
	}
	if !fileInfo.IsDir() {
		c.WriteToStd("", fmt.Errorf("not a directory: %s", targetDirAbsolutePath))
		return
	}

	GetShell().SetWorkingDirectory(targetDirAbsolutePath)
}

func getTargetDirAbsolutePath(pathArg string) string {
	if pathArg[0] == '~' {
		return filepath.Join(os.Getenv("HOME"), pathArg[1:])
	}
	if filepath.IsAbs(pathArg) {
		return filepath.Clean(pathArg)
	}
	return filepath.Join(GetShell().GetWorkingDirectory(), filepath.Clean(pathArg))
}
