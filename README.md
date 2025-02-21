# GoShell

A implementation of a UNIX shell that I worked on while learning Go

## Features

  - Built-in Commands (cd, pwd, echo, type, exit)
  - External Command Execution from PATH
  - Tab completion for commands
  - Output redirection (`>`, `>>`, `1>`, `1>>`, `2>`, `2>>`)
  - Support for quoted strings (both single and double quotes)
  - Escape character support (`\`)

## Installation

### Prerequisites
- Go 1.22.1 or higher

### Building from Source

1. Clone the repository:

```bash
git clone git@github.com:DenizGezgin/GoShell.git
cd goshell
```

2. Run the shell:
```bash
./run.sh
```

### Key Bindings
- `Tab` - Auto-complete command
- `Ctrl+C` - Cancel current input
- `Ctrl+D` - Exit shell

## Project Structure

- `cmd/shell/`
- `internal/shell/`
  - Command runners
  - Shell interface
  - Command repository
- `pkg/`
  - Prefix tree implementation for command completion


