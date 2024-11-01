package logger

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var snipStyle = lipgloss.NewStyle().
	SetString("SNIP").
	Bold(true).
	Foreground(lipgloss.Color("80"))

var errorStyle = lipgloss.NewStyle().
	SetString("ERROR").
	Bold(true).
	Foreground(lipgloss.Color("204"))

type Logger struct{}

func New() *Logger { return &Logger{} }

func (l *Logger) Info(format string, a ...any) {
	fmt.Printf("%s %s\n", snipStyle, fmt.Sprintf(format, a...))
}

func (l *Logger) Error(format string, a ...any) error {
	return fmt.Errorf("%s %s", errorStyle, fmt.Sprintf(format, a...))
}
