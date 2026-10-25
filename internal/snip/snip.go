package snip

import (
	"fmt"
	"regexp"
	"strings"
)

type Snip struct{ Key, Value string }

var (
	keyRe = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)
	posRe = regexp.MustCompile(`\$[0-9]`)
)

func Valid(key, value string) bool {
	return keyRe.MatchString(key) && value != "" && !strings.Contains(value, "\n")
}

func Find(ss []Snip, key string) int {
	for i, s := range ss {
		if s.Key == key {
			return i
		}
	}
	return -1
}

func Set(ss []Snip, key, value string) []Snip {
	if i := Find(ss, key); i >= 0 {
		ss[i].Value = value
		return ss
	}
	return append(ss, Snip{key, value})
}

func Remove(ss []Snip, i int) []Snip { return append(ss[:i], ss[i+1:]...) }

// Expand replaces $1, $2, ... with args.
func Expand(value string, args []string) string {
	for n, a := range args {
		value = strings.ReplaceAll(value, fmt.Sprintf("$%d", n+1), a)
	}
	return value
}

func Format(ss []Snip) string {
	var b strings.Builder
	b.WriteString("# managed by snip, do not edit by hand\n")
	for _, s := range ss {
		if posRe.MatchString(s.Value) {
			fmt.Fprintf(&b, "%s() { %s; }\n", s.Key, s.Value)
		} else {
			fmt.Fprintf(&b, "alias %s='%s'\n", s.Key, strings.ReplaceAll(s.Value, "'", `'\''`))
		}
	}
	return b.String()
}

// Parse is the inverse of Format.
func Parse(data string) []Snip {
	var out []Snip
	for _, l := range strings.Split(data, "\n") {
		switch {
		case strings.HasPrefix(l, "alias "):
			k, v, _ := strings.Cut(strings.TrimPrefix(l, "alias "), "=")
			v = strings.TrimSuffix(strings.TrimPrefix(v, "'"), "'")
			out = append(out, Snip{k, strings.ReplaceAll(v, `'\''`, "'")})
		case strings.HasSuffix(l, "; }") && strings.Contains(l, "() { "):
			k, v, _ := strings.Cut(strings.TrimSuffix(l, "; }"), "() { ")
			out = append(out, Snip{k, v})
		}
	}
	return out
}
