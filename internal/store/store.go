package store

import (
	"bytes"
	"cmp"
	"fmt"
	"os"
	"path/filepath"

	"github.com/iamhectorsosa/snip/internal/snip"
)

func Paths() (file, rc string) {
	home, _ := os.UserHomeDir()
	file = cmp.Or(os.Getenv("SNIP_FILE"), filepath.Join(home, ".config", "snip", "snippets.zsh"))
	rc = cmp.Or(os.Getenv("SNIP_RC"), filepath.Join(home, ".zshrc"))
	return
}

func Load(file string) ([]snip.Snip, error) {
	data, err := os.ReadFile(file)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return snip.Parse(string(data)), nil
}

func Save(file, rc string, ss []snip.Snip) error {
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(file, []byte(snip.Format(ss)), 0o644); err != nil {
		return err
	}
	return hook(file, rc)
}

func hook(file, rc string) error {
	data, _ := os.ReadFile(rc)
	if bytes.Contains(data, []byte(file)) {
		return nil
	}
	f, err := os.OpenFile(rc, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "\n# snip\n[ -f %q ] && source %q\n", file, file)
	return err
}
