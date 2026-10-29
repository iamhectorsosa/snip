package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/iamhectorsosa/snip/internal/clip"
	"github.com/iamhectorsosa/snip/internal/snip"
	"github.com/iamhectorsosa/snip/internal/store"
	"github.com/iamhectorsosa/snip/internal/ui"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:      "snip",
		Usage:     "manage shell snippets as shell aliases/functions (zsh, bash, ksh)",
		UsageText: "snip key='value'  |  snip key [args]  |  snip ls  |  snip rm key",
		Action: func(_ context.Context, c *cli.Command) error {
			args := c.Args().Slice()
			if len(args) == 0 {
				return cli.ShowRootCommandHelp(c)
			}
			file, rc, err := store.Paths()
			if err != nil {
				return err
			}
			ss, err := store.Load(file)
			if err != nil {
				return err
			}

			// Set: snip key='value'
			if key, value, ok := strings.Cut(args[0], "="); ok {
				if !snip.Valid(key, value) {
					return fmt.Errorf("invalid snippet, use: snip key='value'")
				}
				if err := store.Save(file, rc, snip.Set(ss, key, value)); err != nil {
					return err
				}
				fmt.Println(ui.Saved(key, rc))
				return nil
			}

			// Get: snip key [args...]
			i := snip.Find(ss, args[0])
			if i < 0 {
				return fmt.Errorf("snippet %q not found", args[0])
			}
			value := snip.Expand(ss[i].Value, args[1:])
			if err := clip.Copy(value); err != nil {
				return err
			}
			fmt.Println(ui.Copied(value))
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "ls",
				Usage: "list snippets",
				Action: func(_ context.Context, _ *cli.Command) error {
					file, _, err := store.Paths()
					if err != nil {
						return err
					}
					ss, err := store.Load(file)
					if err != nil {
						return err
					}
					fmt.Println(ui.List(ss))
					return nil
				},
			},
			{
				Name:      "rm",
				Usage:     "delete a snippet",
				ArgsUsage: "key",
				Action: func(_ context.Context, c *cli.Command) error {
					file, rc, err := store.Paths()
					if err != nil {
						return err
					}
					ss, err := store.Load(file)
					if err != nil {
						return err
					}
					i := snip.Find(ss, c.Args().First())
					if i < 0 {
						return fmt.Errorf("snippet %q not found", c.Args().First())
					}
					if err := store.Save(file, rc, snip.Remove(ss, i)); err != nil {
						return err
					}
					fmt.Println(ui.Removed(c.Args().First()))
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR", err)
		os.Exit(1)
	}
}
