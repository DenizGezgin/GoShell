package shell

import (
	"os"
	"path/filepath"
	"strings"

	"shell/pkg/prefix_tree"
)

// CommandRepository provides access to both internal and external command names and paths
type CommandRepository interface {
	GetAllCommandNamesGroupedByPrefix(prefix string) [][]string
	GetCommand(name string) Command
	GetExternalCommandPath(name string) string
	GetCommandPath(name string) string
}

type CommandRepositoryImpl struct {
	commandTree *prefix_tree.PrefixTree
	commands    map[string]Command
}

func NewCommandRepository() *CommandRepositoryImpl {
	r := &CommandRepositoryImpl{
		commandTree: prefix_tree.NewTree(),
		commands:    make(map[string]Command),
	}
	r.initializeInternalCommands()
	r.initializeExternalCommands()
	return r
}

func (r *CommandRepositoryImpl) initializeInternalCommands() {
	internalCommands := []string{
		"cd",
		"exit",
		"type",
		"pwd",
	}
	for _, command := range internalCommands {
		r.commandTree.Insert(command)
		r.commands[command] = NewCommand(command, "a shell builtin")
	}
}

func (r *CommandRepositoryImpl) initializeExternalCommands() {
	var paths = strings.Split(os.Getenv("PATH"), ":")
	for _, path := range paths {
		files, err := filepath.Glob(filepath.Join(path, "*"))
		if err != nil {
			continue
		}
		for _, file := range files {
			baseName := filepath.Base(file)
			if r.commands[baseName] != nil {
				continue
			}
			r.commandTree.Insert(baseName)
			r.commands[baseName] = NewCommand(baseName, file)
		}
	}
}

func (r *CommandRepositoryImpl) GetAllCommandNamesGroupedByPrefix(prefix string) [][]string {
	return r.commandTree.GetAllWordsStartingWithGroupedByChildren(prefix)
}

func (r *CommandRepositoryImpl) GetCommand(name string) Command {
	return r.commands[name]
}

func (r *CommandRepositoryImpl) GetExternalCommandPath(name string) string {
	return r.commands[name].GetPath()
}

func (r *CommandRepositoryImpl) GetCommandPath(name string) string {
	if cmd := r.commands[name]; cmd != nil {
		return cmd.GetPath()
	}
	return ""
}
