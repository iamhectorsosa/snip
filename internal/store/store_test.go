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
