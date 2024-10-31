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
	Use:   "snip [name] | [name='text']",
	Short: "Snip is a CLI tool for managing your snippets.",
	Long: `Snip is a CLI tool for managing your snippets.

To get a snippet, use: snip [name]
To add snippets, use: snip [name='text']`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			input := args[0]
			if strings.Contains(input, "=") {
				inputSlice := strings.SplitN(input, "=", 2)
				if len(inputSlice) != 2 {
					return fmt.Errorf("Invalid format. Use: name='text'")
				}

				name := inputSlice[0]
				text := strings.TrimSpace(strings.Trim(inputSlice[1], "'"))

				db, cleanup, err := database.New()
				defer cleanup()
				if err != nil {
					return fmt.Errorf("Error database.New: %v", err)
				}

				if err = db.Create(name, text); err != nil {
					return fmt.Errorf("Error Create: %v", err)
				}

				return nil
			}

			db, cleanup, err := database.New()
			defer cleanup()
			if err != nil {
				return fmt.Errorf("Error database.New: %v", err)
			}

			name := input
			snippet, err := db.Read(name)
			if err != nil {
				return fmt.Errorf("Error Read: %v", err)
			}

			cmd := exec.Command("pbcopy")
			cmd.Stdin = bytes.NewReader([]byte(snippet.Text))
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("Error cmd.Run: %v", err)
			}

			fmt.Printf("Copied to clipboard: %q\n", snippet.Text)
			return nil
		} else {
			return cmd.Help()
		}
	},
}

var list = &cobra.Command{
	Use:   "list",
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

		writer := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
		defer writer.Flush()

		fmt.Fprintln(writer, "ID\tName\tText\tShare")

		for _, snippet := range snippets {
			fmt.Fprintf(writer, "%d\t%s\t%s\t%s\n", snippet.Id, snippet.Name, snippet.Text, fmt.Sprintf("%s='%s'", snippet.Name, snippet.Text))
		}

		return nil
	},
}

var update = &cobra.Command{
	Use:   "update [name='new_text']",
	Short: "Update a snipppet",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]
		inputSlice := strings.SplitN(input, "=", 2)
		if len(inputSlice) != 2 {
			return fmt.Errorf("Invalid format. Use: name='next_text'")
		}

		name := inputSlice[0]
		newText := strings.TrimSpace(strings.Trim(inputSlice[1], "'"))

		db, cleanup, err := database.New()
		defer cleanup()
		if err != nil {
			return fmt.Errorf("Error database.New: %v", err)
		}

		snippet, err := db.Read(name)
		if err != nil {
			return fmt.Errorf("Error Read: %v", err)
		}

		if err = db.Update(store.Snippet{
			Id:   snippet.Id,
			Name: name,
			Text: newText,
		}); err != nil {
			return fmt.Errorf("Error Update: %v", err)
		}

		return nil
	},
}

var delete = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a snippet",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		db, cleanup, err := database.New()
		defer cleanup()
		if err != nil {
			return fmt.Errorf("Error database.New: %v", err)
		}

		if err = db.Delete(name); err != nil {
			return fmt.Errorf("Error Delete: %v", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(list)
	rootCmd.AddCommand(update)
	rootCmd.AddCommand(delete)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
