package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/term"
)

// Shell represents an interactive command-line interface
type Shell interface {
	Run()
	GetWorkingDirectory() string
	SetWorkingDirectory(dir string)
	GetCommandRepository() CommandRepository
}

var globalShell Shell

func GetShell() Shell {
	if globalShell == nil {
		globalShell = NewShell()
	}
	return globalShell
}

func NewShell() Shell {
	workingDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return &shell{
		reader:            bufio.NewReader(os.Stdin),
		workingDirectory:  workingDirectory,
		commandRepository: NewCommandRepository(),
	}
}

type shell struct {
	reader            *bufio.Reader
	workingDirectory  string
	commandRepository CommandRepository
	autoCompleteCache []string
}

func (s *shell) GetCommandRepository() CommandRepository {
	return s.commandRepository
}

func (s *shell) GetWorkingDirectory() string {
	return s.workingDirectory
}

func (s *shell) SetWorkingDirectory(dir string) {
	s.workingDirectory = dir
}

func (s *shell) Run() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	for {
		input, shouldContinue := s.prompt()
		if !shouldContinue {
			break
		}

		commandStr, stdout, stderr := s.parseRedirections(input)
		commandName, args := s.parseCommand(commandStr)
		if commandName == "" {
			continue
		}
		command := s.commandRepository.GetCommand(commandName)
		if command == nil {
			fmt.Fprintf(stderr, "command not found: %s\r\n", commandName)
			continue
		}
		command.SetOutput(stdout, stderr)
		command.Execute(args)
	}
}

func newFileWriter(path string, isAppending bool) io.Writer {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil
	}

	flag := os.O_CREATE | os.O_WRONLY
	if isAppending {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}

	file, err := os.OpenFile(path, flag, 0644)
	if err != nil {
		return nil
	}
	return file
}

func (s *shell) prompt() (string, bool) {
	fmt.Printf("$ ")

	var input []byte
	for {
		b, err := s.reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				fmt.Printf("\r\n")
				return "", false
			}
			panic(err)
		}

		if b != '\t' {
			s.autoCompleteCache = nil
		}

		switch b {
		case '\x03': // Ctrl+C
			fmt.Printf("^C\r\n$ ")
			input = []byte{}
		case '\x04': // Ctrl+D
			if len(input) == 0 {
				fmt.Printf("\r\n")
				return "", false
			}
		case '\x7f', '\b': // Backspace
			if len(input) > 0 {
				input = input[:len(input)-1]
				fmt.Printf("\b \b")
			}
		case '\r': // Enter
			fmt.Printf("\r\n")
			return strings.TrimSpace(string(input)), true
		case '\t': // Tab
			if s.autoCompleteCache != nil {
				// Second tab press - display cached matches
				fmt.Printf("\r\n%s\r\n", strings.Join(s.autoCompleteCache, "  "))
				fmt.Printf("$ %s", string(input))
				s.autoCompleteCache = nil
			} else {
				// First tab press
				completedCommands := NewAutoComplete(string(input)).CompleteCommand()
				switch len(completedCommands) {
				case 0:
					fmt.Printf("\a")
				case 1:
					fmt.Printf("\r\033[K$ %s ", completedCommands[0])
					input = []byte(completedCommands[0] + " ")
				default:
					fmt.Printf("\a")
					s.autoCompleteCache = completedCommands
				}
			}
		default:
			input = append(input, b)
			fmt.Printf("%c", b)
		}
	}
}

func (s *shell) parseRedirections(input string) (string, io.Writer, io.Writer) {
	stderrRedirection := regexp.MustCompile(`2>>|2>`)
	stderrMatches := stderrRedirection.FindStringIndex(input)
	if len(stderrMatches) == 2 {
		command := strings.TrimSpace(input[:stderrMatches[0]])
		outputFile := strings.TrimSpace(input[stderrMatches[1]:])
		isAppending := input[stderrMatches[0]:stderrMatches[1]] == "2>>"
		stderr := newFileWriter(outputFile, isAppending)
		return command, os.Stdout, stderr
	}

	stdoutRedirection := regexp.MustCompile(`1>>|>>|>|1>`)
	stdoutMatches := stdoutRedirection.FindStringIndex(input)
	if len(stdoutMatches) == 2 {
		command := strings.TrimSpace(input[:stdoutMatches[0]])
		outputFile := strings.TrimSpace(input[stdoutMatches[1]:])
		operator := input[stdoutMatches[0]:stdoutMatches[1]]
		isAppending := operator == "1>>" || operator == ">>"
		stdout := newFileWriter(outputFile, isAppending)
		return command, stdout, os.Stderr
	}

	return input, os.Stdout, os.Stderr
}

func (s *shell) parseCommand(input string) (string, []string) {
	var tokens []string
	i := 0

	for i < len(input) {
		// Skip whitespace
		for i < len(input) && unicode.IsSpace(rune(input[i])) {
			i++
		}
		if i >= len(input) {
			break
		}

		var tokenBuilder strings.Builder
		// Parse until next whitespace
		for i < len(input) && !unicode.IsSpace(rune(input[i])) {
			if input[i] == '\'' {
				i++ // Skip opening quote

				// Add the content between quotes to token
				for i < len(input) && input[i] != '\'' {
					tokenBuilder.WriteByte(input[i])
					i++
				}
			} else if input[i] == '"' {
				i++ // Skip opening quote

				// Add the content between quotes to token
				for i < len(input) && input[i] != '"' {
					if input[i] == '\\' {
						i++
						switch input[i] {
						case '"':
							tokenBuilder.WriteByte('"')
						case '\\':
							tokenBuilder.WriteByte('\\')
						case '\n':
							tokenBuilder.WriteByte('\n')
						case '$':
							tokenBuilder.WriteByte('$')
						default:
							tokenBuilder.WriteByte('\\')
							tokenBuilder.WriteByte(input[i])
						}
					} else {
						tokenBuilder.WriteByte(input[i])
					}
					i++
				}
			} else {
				if input[i] == '\\' {
					i++
				}
				tokenBuilder.WriteByte(input[i])
			}

			i++
		}

		tokens = append(tokens, tokenBuilder.String())
	}

	if len(tokens) == 0 {
		return "", nil
	}
	return tokens[0], tokens[1:]
}
