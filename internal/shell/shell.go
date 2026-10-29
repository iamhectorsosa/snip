package shell

import (
	"cmp"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Shell struct {
	Name string
	RC   string
}

var Supported = []string{"zsh", "bash", "ksh"}

var rcFile = map[string]string{
	"zsh":  ".zshrc",
	"bash": ".bashrc",
	"ksh":  ".kshrc",
}

func Detect() (Shell, error) {
	raw := cmp.Or(os.Getenv("SNIP_SHELL"), os.Getenv("SHELL"))
	supported := strings.Join(Supported, ", ")
	if raw == "" {
		return Shell{}, fmt.Errorf("could not detect your shell ($SHELL is empty): snip supports %s (set SNIP_SHELL to one of them to override)", supported)
	}

	name := strings.ToLower(filepath.Base(raw))
	name = strings.TrimSuffix(name, ".exe")

	rc, ok := rcFile[name]
	if !ok {
		return Shell{}, fmt.Errorf("unsupported shell %q: snip supports %s (set SNIP_SHELL to one of them to override)", name, supported)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return Shell{}, err
	}

	s := Shell{Name: name, RC: filepath.Join(home, rc)}
	if snipRC := os.Getenv("SNIP_RC"); snipRC != "" {
		s.RC = snipRC
	}
	return s, nil
}
