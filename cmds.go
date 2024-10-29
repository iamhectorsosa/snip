package main

import (
	"fmt"
	"strconv"

	"github.com/iamhectorsosa/snippets/internal/database"
	"github.com/iamhectorsosa/snippets/internal/store"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "snippets [command]",
	Short: "Snippets is a terminal tool for managing your snippets",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
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

		db, cleanup := database.New()
		defer cleanup()

		err := db.Create(name, text)
		if err != nil {
			return fmt.Errorf("Error Create: %v", err)
		}
		return nil
	},
}

var view = &cobra.Command{
	Use:   "view [id]",
	Short: "View a snipppet",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return fmt.Errorf("Need to provide a minimum of 1 arguments: %v", err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		idStr := args[0]
		db, cleanup := database.New()
		defer cleanup()

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fmt.Errorf("Error Atoi: %v", err)
		}

		snippet, err := db.Read(id)
		if err != nil {
			return fmt.Errorf("Error Read: %v", err)
		}

		fmt.Println(snippet.Id, snippet.Name, snippet.Text)
		return nil
	},
}

var list = &cobra.Command{
	Use:   "list",
	Short: "List all snippets",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		db, cleanup := database.New()
		defer cleanup()

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
	Use:   "update [id] [name] [text]",
	Short: "Update a snipppet",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(3)(cmd, args); err != nil {
			return fmt.Errorf("Need to provide a minimum of 3 arguments: %v", err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		idStr, name, text := args[0], args[1], args[2]
		db, cleanup := database.New()
		defer cleanup()

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fmt.Errorf("Error Atoi: %v", err)
		}

		snippet := store.Snippet{Id: id,
			Name: name,
			Text: text}

		err = db.Update(snippet)

		if err != nil {
			return fmt.Errorf("Error Update: %v", err)
		}

		fmt.Println(snippet.Id, snippet.Name, snippet.Text)
		return nil
	},
}

var delete = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a snippet",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return fmt.Errorf("Need to provide a minimum of 1 arguments: %v", err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		idStr := args[0]
		db, cleanup := database.New()
		defer cleanup()

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fmt.Errorf("Error Atoi: %v", err)
		}

		err = db.Delete(id)
		if err != nil {
			return fmt.Errorf("Error Delete: %v", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(add)
	rootCmd.AddCommand(view)
	rootCmd.AddCommand(list)
	rootCmd.AddCommand(update)
	rootCmd.AddCommand(delete)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
