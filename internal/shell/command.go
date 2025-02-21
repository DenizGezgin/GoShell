package shell

import (
	"io"
	"os"
)

// Command represents a shell command that can be executed
type Command interface {
	Execute(args []string)
	GetName() string
	GetPath() string
	SetOutput(stdout, stderr io.Writer)
}

type command struct {
	name   string
	path   string
	stdout io.Writer
	stderr io.Writer
}

func NewCommand(name string, path string) Command {
	return &command{
		name:   name,
		path:   path,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

func (c *command) SetOutput(stdout, stderr io.Writer) {
	c.stdout = stdout
	c.stderr = stderr
}

func (c *command) Execute(args []string) {
	commandRunnerBuilder := &CommandRunnerBuilder{}
	runner, err := commandRunnerBuilder.Build(c.name, c.stdout, c.stderr)
	if err != nil {
		c.stderr.Write([]byte(err.Error() + "\r\n"))
		return
	}
	runner.Run(args)
}

func (c *command) GetName() string {
	return c.name
}

func (c *command) GetPath() string {
	return c.path
}
