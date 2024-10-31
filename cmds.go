package main

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/iamhectorsosa/snippets/internal/database"
	"github.com/iamhectorsosa/snippets/internal/store"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "snippets [name]",
	Short: "Snippets is a terminal tool for managing your snippets",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MaximumNArgs(1)(cmd, args); err != nil {
			return fmt.Errorf("Can only provide a maximum of one arguments: %v", err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			name := args[0]
			db, cleanup, err := database.New()
			defer cleanup()
			if err != nil {
				return fmt.Errorf("Error database.New: %v", err)
			}

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

var add = &cobra.Command{
	Use:   "add [name] [text]",
	Short: "Add a snippet",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(2)(cmd, args); err != nil {
			return fmt.Errorf("Need to provide a minimum of two arguments: %v", err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		name, text := args[0], args[1]
		db, cleanup, err := database.New()
		defer cleanup()
		if err != nil {
			return fmt.Errorf("Error database.New: %v", err)
		}

		if err = db.Create(name, text); err != nil {
			return fmt.Errorf("Error Create: %v", err)
		}

		return nil
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

		for _, snippet := range snippets {
			fmt.Println(snippet.Id, snippet.Name, snippet.Text)
		}
		return nil
	},
}

var update = &cobra.Command{
	Use:   "update [name] [text]",
	Short: "Update a snipppet",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(2)(cmd, args); err != nil {
			return fmt.Errorf("Need to provide a minimum of 2 arguments: %v", err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		name, text := args[0], args[1]
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
			Text: text,
		}); err != nil {
			return fmt.Errorf("Error Update: %v", err)
		}

		return nil
	},
}

var delete = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a snippet",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return fmt.Errorf("Need to provide a minimum of 1 arguments: %v", err)
		}
		return nil
	},
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
	rootCmd.AddCommand(add)
	rootCmd.AddCommand(list)
	rootCmd.AddCommand(update)
	rootCmd.AddCommand(delete)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
