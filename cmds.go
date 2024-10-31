package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"

	"github.com/iamhectorsosa/snip/internal/database"
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
		if len(args) >= 1 {
			input := args[0]
			if strings.Contains(input, "=") {
				inputSlice := strings.SplitN(input, "=", 2)
				if len(inputSlice) != 2 {
					return fmt.Errorf("Invalid format. Use: key='value'")
				}

				key := inputSlice[0]
				value := strings.TrimSpace(strings.Trim(inputSlice[1], "'"))

				db, cleanup, err := database.New()
				defer cleanup()
				if err != nil {
					return fmt.Errorf("Error database.New: %v", err)
				}

				if err = db.Create(key, value); err != nil {
					return fmt.Errorf("Error Create: %v", err)
				}

				return nil
			}

			db, cleanup, err := database.New()
			defer cleanup()
			if err != nil {
				return fmt.Errorf("Error database.New: %v", err)
			}

			key := input
			snippet, err := db.Read(key)
			if err != nil {
				return fmt.Errorf("Error Read: %v", err)
			}

			value := snippet.Value
			for i, arg := range args[1:] {
				placeholder := fmt.Sprintf("$%d", i+1)
				value = strings.ReplaceAll(value, placeholder, arg)
			}

			cmd := exec.Command("pbcopy")
			cmd.Stdin = bytes.NewReader([]byte(value))
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("Error pbcopy in cmd.Run: %v", err)
			}

			fmt.Printf("Copied to clipboard: %q\n", value)
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
		db, cleanup, err := database.New()
		defer cleanup()
		if err != nil {
			return fmt.Errorf("Error database.New: %v", err)
		}

		snippets, err := db.ReadAll()
		if err != nil {
			return fmt.Errorf("Error ReadAll: %v", err)
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
		defer writer.Flush()

		fmt.Fprintln(writer, "Key\tValue")

		for _, snippet := range snippets {
			fmt.Fprintf(writer, "%s\t%s\n", snippet.Key, snippet.Value)
		}

		return nil
	},
}

var update = &cobra.Command{
	Use:   "update [key='new_value']",
	Short: "Update a snipppet",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]
		inputSlice := strings.SplitN(input, "=", 2)
		if len(inputSlice) != 2 {
			return fmt.Errorf("Invalid format. Use: key='new_value'")
		}

		key := inputSlice[0]
		newValue := strings.TrimSpace(strings.Trim(inputSlice[1], "'"))

		db, cleanup, err := database.New()
		defer cleanup()
		if err != nil {
			return fmt.Errorf("Error database.New: %v", err)
		}

		snippet, err := db.Read(key)
		if err != nil {
			return fmt.Errorf("Error Read: %v", err)
		}

		if err = db.Update(store.Snippet{
			Id:    snippet.Id,
			Key:   key,
			Value: newValue,
		}); err != nil {
			return fmt.Errorf("Error Update: %v", err)
		}

		return nil
	},
}

var delete = &cobra.Command{
	Use:   "delete [key]",
	Short: "Delete a snippet",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		db, cleanup, err := database.New()
		defer cleanup()
		if err != nil {
			return fmt.Errorf("Error database.New: %v", err)
		}

		if err = db.Delete(key); err != nil {
			return fmt.Errorf("Error Delete: %v", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(ls)
	rootCmd.AddCommand(update)
	rootCmd.AddCommand(delete)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
