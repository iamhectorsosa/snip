package snip

import (
	"reflect"
	"testing"
)

func TestFormatParseRoundTrip(t *testing.T) {
	in := []Snip{{"gs", "git status"}, {"q", "echo 'hi'"}, {"greet", "echo hello $1"}}
	if got := Parse(Format(in)); !reflect.DeepEqual(got, in) {
		t.Fatalf("got %v, want %v", got, in)
	}
}

func TestFormatUsesFunctionForPositionalArgs(t *testing.T) {
	want := "# managed by snip, do not edit by hand\ng() { echo $1; }\n"
	if got := Format([]Snip{{"g", "echo $1"}}); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestExpand(t *testing.T) {
	if got := Expand("cp $1 $2", []string{"a", "b"}); got != "cp a b" {
		t.Fatalf("got %q", got)
	}
}

func TestSetFindRemove(t *testing.T) {
	ss := Set(nil, "a", "1")
	ss = Set(ss, "b", "2")
	ss = Set(ss, "a", "3")
	if len(ss) != 2 || ss[Find(ss, "a")].Value != "3" {
		t.Fatalf("unexpected %v", ss)
	}
	ss = Remove(ss, Find(ss, "a"))
	if Find(ss, "a") != -1 || len(ss) != 1 {
		t.Fatalf("unexpected %v", ss)
	}
}

func TestValid(t *testing.T) {
	cases := []struct {
		k, v string
		ok   bool
	}{
		{"gs", "git status", true},
		{"bad key", "x", false},
		{"k", "", false},
		{"k", "a\nb", false},
	}
	for _, c := range cases {
		if got := Valid(c.k, c.v); got != c.ok {
			t.Errorf("Valid(%q, %q) = %v", c.k, c.v, got)
		}
	}
}
