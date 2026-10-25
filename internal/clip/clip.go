package clip

import (
	"errors"
	"os/exec"
	"strings"
)

var tools = [][]string{{"pbcopy"}, {"wl-copy"}, {"xclip", "-selection", "clipboard"}}

func Copy(s string) error {
	for _, c := range tools {
		if _, err := exec.LookPath(c[0]); err == nil {
			cmd := exec.Command(c[0], c[1:]...)
			cmd.Stdin = strings.NewReader(s)
			return cmd.Run()
		}
	}
	return errors.New("no clipboard tool found")
}
