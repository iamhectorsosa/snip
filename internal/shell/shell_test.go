package shell

import (
	"strings"
	"testing"
)

func TestDetect(t *testing.T) {
	const home = "/home/user"
	t.Setenv("HOME", home)

	cases := []struct {
		name    string
		env     map[string]string
		want    Shell
		wantErr string
	}{
		{
			name: "zsh from SHELL",
			env:  map[string]string{"SHELL": "/bin/zsh"},
			want: Shell{Name: "zsh", RC: "/home/user/.zshrc"},
		},
		{
			name: "bash from SHELL",
			env:  map[string]string{"SHELL": "/bin/bash"},
			want: Shell{Name: "bash", RC: "/home/user/.bashrc"},
		},
		{
			name: "ksh from SHELL",
			env:  map[string]string{"SHELL": "/bin/ksh"},
			want: Shell{Name: "ksh", RC: "/home/user/.kshrc"},
		},
		{
			name: "SNIP_SHELL overrides SHELL",
			env:  map[string]string{"SHELL": "/bin/zsh", "SNIP_SHELL": "bash"},
			want: Shell{Name: "bash", RC: "/home/user/.bashrc"},
		},
		{
			name: "SNIP_RC overrides rc path only",
			env:  map[string]string{"SHELL": "/bin/zsh", "SNIP_RC": "/x/rc"},
			want: Shell{Name: "zsh", RC: "/x/rc"},
		},
		{
			name:    "unsupported shell fish",
			env:     map[string]string{"SHELL": "/usr/bin/fish"},
			wantErr: "unsupported",
		},
		{
			name:    "unsupported shell sh",
			env:     map[string]string{"SHELL": "/bin/sh"},
			wantErr: "unsupported",
		},
		{
			name:    "unsupported shell dash",
			env:     map[string]string{"SHELL": "/bin/dash"},
			wantErr: "unsupported",
		},
		{
			name:    "SNIP_SHELL unsupported",
			env:     map[string]string{"SHELL": "/bin/zsh", "SNIP_SHELL": "fish"},
			wantErr: "unsupported",
		},
		{
			name:    "no shell detected",
			env:     map[string]string{"SHELL": "", "SNIP_SHELL": ""},
			wantErr: "could not detect",
		},
		{
			name: "uppercase SHELL normalised",
			env:  map[string]string{"SHELL": "/bin/ZSH"},
			want: Shell{Name: "zsh", RC: "/home/user/.zshrc"},
		},
		{
			name: "bash.exe normalised",
			env:  map[string]string{"SHELL": "/mnt/c/tools/bash.exe"},
			want: Shell{Name: "bash", RC: "/home/user/.bashrc"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for _, k := range []string{"SHELL", "SNIP_SHELL", "SNIP_RC"} {
				t.Setenv(k, c.env[k])
			}

			got, err := Detect()
			if c.wantErr != "" {
				if err == nil {
					t.Fatalf("want error containing %q, got nil (shell=%+v)", c.wantErr, got)
				}
				if !strings.Contains(err.Error(), c.wantErr) {
					t.Fatalf("got error %q, want it to contain %q", err.Error(), c.wantErr)
				}
				if !strings.Contains(err.Error(), "zsh, bash, ksh") {
					t.Fatalf("error %q should mention supported shells", err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != c.want {
				t.Fatalf("got %+v, want %+v", got, c.want)
			}
		})
	}
}
