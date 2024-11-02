# Snip

![code coverage badge](https://github.com/iamhectorsosa/snip/actions/workflows/ci.yml/badge.svg)

Snip (snippets) is a command-line tool for managing your snippets using a local SQLite database.

![demo](./demo.gif)

## Usage

This command installs the binary in your GOBIN directory (default: `~/go/bin`). It also creates a SQLite database file at a fixed location in your home directory: `~/.config/snip/local.db`.

```bash
go install github.com/iamhectorsosa/snip@latest
```

**Note:** This tool is currently implemented only for macOS (Darwin). Support for other operating systems has not been implemented yet.

Manage snippets with the same known patterns as aliases. Calling snippets automatically copies them to your system clipboard. Snip can import and export from/to CSV, supporting both local and remote paths/URLs for flexibility.

Here are some basic commands:

```bash
# Creates a snippet
❯ snip [key='value']

# Calls a snippets
❯ snip [key]

# Export snippets
❯ snip export --path ~/

# Import snippets
❯ snip import --url https://gist..

# List of snippets
❯ snip ls
SNIP Found 6 snippets...
KEY          VALUE
grep-s       grep -rn "$1" .
gh-rm-b      git branch | grep -v "^\*" | xargs git branch -d
gh-config    git config --list | grep -E "user.email|user.name|user.signingkey|commit.gpgsign"
ls-h         ls -d .*
gh-rm        rm -rf .git
go-bin       ls -l ~/go/bin
```

## Commands

Run the help command to get an updated list of all commands.

```bash
❯ snip help
Snip is a CLI tool for managing your snippets.

To get a snippet, use: snip [key] [...$1]
To add snippets, use: snip [key='value']

Usage:
  snip [key] [...$1] | [key='value'] [flags]
  snip [command]

Available Commands:
  delete      Delete a snippet
  export      Export all snippets
  help        Help about any command
  import      Import snippets
  ls          List all snippets
  reset       Reset all snippets
  update      Update a snipppet
```

## Development

1. Clone the repository:

```bash
gh repo clone iamhectorsosa/snip
cd snip
```

2. Build the project:

```bash
CGO_ENABLED=1 go build -v -o snip .
```
