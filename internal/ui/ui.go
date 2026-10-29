package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/iamhectorsosa/snip/internal/snip"
)

var (
	accent = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	grey   = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	brand  = accent.Render("✂️ SNIP")
)

func Saved(k, rc string) string {
	return fmt.Sprintf("%s saved %s\n%s", brand, accent.Render(k),
		grey.Render(fmt.Sprintf("Run `. %s` or open a new shell.", rc)))
}

func Copied(value string) string { return fmt.Sprintf("%s copied %s", brand, value) }

func Removed(k string) string { return fmt.Sprintf("%s removed %s", brand, accent.Render(k)) }

func List(ss []snip.Snip) string {
	if len(ss) == 0 {
		return fmt.Sprintf("%s no snippets yet, try: snip key='value'", brand)
	}
	w := len("KEY")
	for _, s := range ss {
		w = max(w, lipgloss.Width(s.Key))
	}
	row := func(k, v string) string { return fmt.Sprintf("%-*s  %s", w, k, v) }

	out := []string{
		fmt.Sprintf("%s %d snippets\n", brand, len(ss)),
		accent.Render(row("KEY", "VALUE")),
	}
	for i, s := range ss {
		line := row(s.Key, s.Value)
		if i%2 == 1 {
			line = grey.Render(line)
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}
