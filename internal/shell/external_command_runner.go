package shell

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// ExternalCommandRunner executes system commands from the PATH
type ExternalCommandRunner interface {
	Run(args []string)
}

type externalCommandRunner struct {
	absolutePath string
	stdout       io.Writer
	stderr       io.Writer
}

func NewExternalCommandRunner(name string, stdout, stderr io.Writer) (ExternalCommandRunner, error) {
	absolutePath := GetShell().GetCommandRepository().GetExternalCommandPath(name)
	if absolutePath == "" {
		return nil, fmt.Errorf("%s: command not found", name)
	}

	return &externalCommandRunner{
		absolutePath: absolutePath,
		stdout:       stdout,
		stderr:       stderr,
	}, nil
}

func (e *externalCommandRunner) Run(args []string) {
	cmd := exec.Command(e.absolutePath, args...)

	var stderrBuf strings.Builder
	cmd.Stderr = &stderrBuf
	out, _ := cmd.Output()

	if stderrBuf.Len() > 0 {
		fmt.Fprintf(e.stderr, "%s\r\n", strings.TrimSpace(stderrBuf.String()))
		return
	}

	cmdOutput := strings.TrimSpace(string(out))
	cmdOutput = strings.ReplaceAll(cmdOutput, "\n", "\r\n")
	if cmdOutput != "" {
		fmt.Fprintf(e.stdout, "%s\r\n", cmdOutput)
	}
}
