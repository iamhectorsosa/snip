package store

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/iamhectorsosa/snip/internal/snip"
)

func TestSaveLoadAndHookOnce(t *testing.T) {
	dir := t.TempDir()
	file, rc := filepath.Join(dir, "a", "snippets.zsh"), filepath.Join(dir, ".zshrc")
	ss := []snip.Snip{{Key: "gs", Value: "git status"}}

	for range 2 {
		if err := Save(file, rc, ss); err != nil {
			t.Fatal(err)
		}
	}
	got, err := Load(file)
	if err != nil || !reflect.DeepEqual(got, ss) {
		t.Fatalf("got %v, %v", got, err)
	}
	data, _ := os.ReadFile(rc)
	if n := strings.Count(string(data), "source"); n != 1 {
		t.Fatalf("rc sourced %d times", n)
	}
}

func TestLoadMissingFile(t *testing.T) {
	if ss, err := Load(filepath.Join(t.TempDir(), "nope")); ss != nil || err != nil {
		t.Fatalf("got %v, %v", ss, err)
	}
}

func TestPaths(t *testing.T) {
	const home = "/home/user"
	defaultFile := "/home/user/.config/snip/snippets.sh"

	cases := []struct {
		name     string
		env      map[string]string
		wantFile string
		wantRC   string
		wantErr  string
	}{
		{
			name:     "default file and rc",
			env:      map[string]string{"SHELL": "/bin/zsh"},
			wantFile: defaultFile,
			wantRC:   "/home/user/.zshrc",
		},
		{
			name:     "SNIP_FILE and SNIP_RC overrides honoured",
			env:      map[string]string{"SHELL": "/bin/zsh", "SNIP_FILE": "/x/file", "SNIP_RC": "/x/rc"},
			wantFile: "/x/file",
			wantRC:   "/x/rc",
		},
		{
			name:    "unsupported shell errors even with overrides set",
			env:     map[string]string{"SHELL": "/usr/bin/fish", "SNIP_FILE": "/x/file", "SNIP_RC": "/x/rc"},
			wantErr: "unsupported",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("HOME", home)
			for _, k := range []string{"SHELL", "SNIP_SHELL", "SNIP_FILE", "SNIP_RC"} {
				t.Setenv(k, c.env[k])
			}

			file, rc, err := Paths()
			if c.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), c.wantErr) {
					t.Fatalf("got err %v, want it to contain %q", err, c.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if file != c.wantFile || rc != c.wantRC {
				t.Fatalf("got file=%q rc=%q, want file=%q rc=%q", file, rc, c.wantFile, c.wantRC)
			}
		})
	}
}
