package ui

import (
	"strings"
	"testing"

	"github.com/iamhectorsosa/snip/internal/snip"
)

func TestList(t *testing.T) {
	out := List([]snip.Snip{{Key: "gs", Value: "git status"}, {Key: "gp", Value: "git push"}})
	for _, want := range []string{"SNIP", "2 snippets", "KEY", "VALUE", "gs", "git status", "gp", "git push"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in %q", want, out)
		}
	}
}

func TestListEmpty(t *testing.T) {
	if out := List(nil); !strings.Contains(out, "SNIP") || !strings.Contains(out, "no snippets") {
		t.Fatalf("got %q", out)
	}
}

func TestMessagesAreBranded(t *testing.T) {
	cases := map[string]string{
		Saved("gs", "~/.zshrc"): "saved",
		Copied("git status"):    "copied",
		Removed("gs"):           "removed",
	}
	for out, verb := range cases {
		if !strings.Contains(out, "SNIP") || !strings.Contains(out, verb) {
			t.Errorf("got %q, want brand and %q", out, verb)
		}
	}
}
