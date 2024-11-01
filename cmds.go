package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/iamhectorsosa/snip/internal/database"
	"github.com/iamhectorsosa/snip/internal/logger"
	"github.com/iamhectorsosa/snip/internal/store"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "snip [key] [...$1] | [key='value']",
	Short: "Snip is a CLI tool for managing your snippets.",
	Long: `Snip is a CLI tool for managing your snippets.

To get a snippet, use: snip [key] [...$1]
To add snippets, use: snip [key='value']`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logger.New()
		if len(args) >= 1 {
			input := args[0]
			if strings.Contains(input, "=") {
				inputSlice := strings.SplitN(input, "=", 2)
				if len(inputSlice) != 2 {
					return log.Error("invalid format. Use: key='value'")
				}

				key := inputSlice[0]
				value := strings.TrimSpace(strings.Trim(inputSlice[1], "'"))

				db, cleanup, err := database.New()
				defer cleanup()
				if err != nil {
					return log.Error("database.New: %v", err)
				}

				if err = db.Create(key, value); err != nil {
					return log.Error("db.Create: %v", err)
				}

				log.Info("Snippet successfully created, key=%q value=%q.", key, value)
				return nil
			}

			db, cleanup, err := database.New()
			defer cleanup()
			if err != nil {
				return log.Error("database.New: %v", err)
			}

			key := input
			snippet, err := db.Read(key)
			if err != nil {
				return log.Error("db.Read: err=%v", err)
			}

			value := snippet.Value
			for i, arg := range args[1:] {
				placeholder := fmt.Sprintf("$%d", i+1)
				value = strings.ReplaceAll(value, placeholder, arg)
			}

			cmd := exec.Command("pbcopy")
			cmd.Stdin = bytes.NewReader([]byte(value))
			if err := cmd.Run(); err != nil {
				return log.Error("pbcopy in cmd.Run: %v", err)
			}

			log.Info("Copied to clipboard: value=%q", value)
			return nil
		} else {
			return cmd.Help()
		}
	},
}

var ls = &cobra.Command{
	Use:   "ls",
	Short: "List all snippets",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logger.New()
		db, cleanup, err := database.New()
		defer cleanup()
		if err != nil {
			return log.Error("database.New: %v", err)
		}

		snippets, err := db.ReadAll()
		if err != nil {
			return log.Error("db.ReadAll: %v", err)
		}

		log.Info("Found %d snippets...", len(snippets))

		maxKeyLen, maxValueLen := 0, 0
		for _, s := range snippets {
			if len(s.Key) > maxKeyLen {
				maxKeyLen = len(s.Key) + 6
			}
			if len(s.Value) > maxValueLen {
				maxValueLen = len(s.Value)
			}
		}

		snippets = append([]store.Snippet{store.Snippet{
			Id:    0,
			Key:   "KEY",
			Value: "VALUE",
		}}, snippets...)

		evenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
		for i, s := range snippets {
			key := fmt.Sprintf("%-*s", maxKeyLen, s.Key)
			value := fmt.Sprintf("%-*s", maxValueLen, s.Value)
			if i%2 == 0 {
				key = evenStyle.Render(key)
				value = evenStyle.Render(value)
			}
			fmt.Println(key, value)
		}

		return nil
	},
}

var update = &cobra.Command{
	Use:   "update [key='new_value']",
	Short: "Update a snipppet",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logger.New()
		input := args[0]
		inputSlice := strings.SplitN(input, "=", 2)
		if len(inputSlice) != 2 {
			return log.Error("invalid format. Use: key='new_value'")
		}

		key := inputSlice[0]
		newValue := strings.TrimSpace(strings.Trim(inputSlice[1], "'"))

		db, cleanup, err := database.New()
		defer cleanup()
		if err != nil {
			return log.Error("database.New: %v", err)
		}

		snippet, err := db.Read(key)
		if err != nil {
			return log.Error("db.Read: %v", err)
		}

		if err = db.Update(store.Snippet{
			Id:    snippet.Id,
			Key:   key,
			Value: newValue,
		}); err != nil {
			return log.Error("db.Update: %v", err)
		}

		log.Info("Snippet successfully updated, key=%q value=%q.", key, newValue)
		return nil
	},
}

var delete = &cobra.Command{
	Use:   "delete [key]",
	Short: "Delete a snippet",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logger.New()
		key := args[0]
		db, cleanup, err := database.New()
		defer cleanup()
		if err != nil {
			return log.Error("database.New: %v", err)
		}

		if err = db.Delete(key); err != nil {
			return log.Error("db.Delete: %v", err)
		}

		log.Info("Snippet successfully deleted, key=%q.", key)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(ls)
	rootCmd.AddCommand(update)
	rootCmd.AddCommand(delete)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
}
