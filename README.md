# Snip ✂️

![Snip Demo](./demo.gif)

Snip is a dead simple CLI tool that manages shell snippets as zsh aliases and functions, and copies them to your clipboard on demand.

Ideal for keeping long commands one short key away.

## Quick Overview

- **Aliases + Functions**: Plain snippets become `alias k='v'`; snippets with `$1`, `$2`, ... become shell functions.
- **Auto-sourced**: Snippets live in one managed file, sourced from your `.zshrc` automatically.
- **Copy on demand**: `snip key` expands arguments and copies the result to your clipboard.
- **Pretty Output**: Every command prints a branded `✂️ SNIP` message, and `snip ls` shows an easy-to-read table with alternating row colors.

## Requirements

- Go 1.22 or later
- A clipboard tool: `pbcopy` (macOS), `wl-copy` (Wayland), or `xclip` (X11)

## Installation

```
go install github.com/iamhectorsosa/snip@latest
```

Or from source:

```
git clone https://github.com/iamhectorsosa/snip.git
cd snip
go install
```

## Usage

### Basic Usage

```
# Save a snippet
snip gs='git status'

# Save a snippet with positional arguments
snip cpto='cp $1 $2'

# Copy a snippet to your clipboard
snip gs

# Copy a snippet with arguments filled in
snip cpto a.txt b.txt

# List snippets
snip ls

# Delete a snippet
snip rm gs
```

### Output

Snippets are stored in `~/.config/snip/snippets.zsh`:

```
# managed by snip, do not edit by hand
alias gs='git status'
cpto() { cp $1 $2; }
```

Run `source ~/.zshrc` (or open a new shell) after saving to use them directly.

### Commands

- `snip key='value'`: Save a snippet.
- `snip key [args...]`: Copy a snippet, replacing `$1`, `$2`, ... with args.
- `snip ls`: List snippets.
- `snip rm key`: Delete a snippet.

### Environment

- `SNIP_FILE`: Snippets file (default: `~/.config/snip/snippets.zsh`)
- `SNIP_RC`: Shell rc to source it from (default: `~/.zshrc`)

## Development

```
go mod tidy
go test ./...
go build
```

```
main.go            # CLI wiring
internal/snip/     # parse, format, expand (pure)
internal/store/    # file I/O + rc hook
internal/clip/     # clipboard
internal/ui/       # lipgloss styling
```
