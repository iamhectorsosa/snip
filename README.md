# Snippets (snip)

Snippets (snip) is a command-line tool for managing your code snippets using a local SQLite database.

## Usage

```bash
go install github.com/iamhectorsosa/snip@latest
```

This command installs the binary in your GOBIN directory (default: `~/go/bin`).

## Commands

Run the help command to get an updated list of all commands.

```bash
❯ snip help

Snippets is a terminal tool for managing your snippets.

To get a snippet, use: snip [name]
To add snippets, use: snip [name='text']

Usage:
  snip [name] | [name='text'] [flags]
  snip [command]

Available Commands:
  delete      Delete a snippet
  help        Help about any command
  list        List all snippets
  update      Update a snipppet
```

## Installation

1. Clone the repository:

```bash
git clone github.com/iamhectorsosa/snip
cd snip
```

2. Build the project:

```bash
CGO_ENABLED=1 go build -v -o snip .
```
