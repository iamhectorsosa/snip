package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

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
					return log.Error("database.New, err=%v", err)
				}

				if err = db.Create(key, value); err != nil {
					return log.Error("db.Create, err=%v", err)
				}

				log.Info("Snippet successfully created, key=%q value=%q.", key, value)
				return nil
			}

			db, cleanup, err := database.New()
			defer cleanup()
			if err != nil {
				return log.Error("database.New, err=%v", err)
			}

			key := input
			snippet, err := db.Read(key)
			if err != nil {
				return log.Error("db.Read, err=%v", err)
			}

			value := snippet.Value
			for i, arg := range args[1:] {
				placeholder := fmt.Sprintf("$%d", i+1)
				value = strings.ReplaceAll(value, placeholder, arg)
			}

			cmd := exec.Command("pbcopy")
			cmd.Stdin = bytes.NewReader([]byte(value))
			if err := cmd.Run(); err != nil {
				return log.Error("pbcopy in cmd.Run, err=%v", err)
			}

			log.Info("Copied to clipboard, value=%q", value)
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
			return log.Error("database.New, err=%v", err)
		}

		snippets, err := db.ReadAll()
		if err != nil {
			return log.Error("db.ReadAll, err=%v", err)
		}

		log.Info("Found %d snippets...", len(snippets))

		if len(snippets) == 0 {
			return nil
		}

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
			return log.Error("database.New, err=%v", err)
		}

		snippet, err := db.Read(key)
		if err != nil {
			return log.Error("db.Read, err=%v", err)
		}

		if err = db.Update(store.Snippet{
			Id:    snippet.Id,
			Key:   key,
			Value: newValue,
		}); err != nil {
			return log.Error("db.Update, err=%v", err)
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
			return log.Error("database.New, err=%v", err)
		}

		if err = db.Delete(key); err != nil {
			return log.Error("db.Delete, err=%v", err)
		}

		log.Info("Snippet successfully deleted, key=%q.", key)
		return nil
	},
}

var reset = &cobra.Command{
	Use:   "reset",
	Short: "Reset all snippets",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logger.New()
		db, cleanup, err := database.New()
		defer cleanup()
		if err != nil {
			return log.Error("database.New, err=%v", err)
		}

		if err = db.Reset(); err != nil {
			return log.Error("db.Reset, err=%v", err)
		}

		log.Info("Snippets have been successfully reset")
		return nil
	},
}

var (
	exportPath     string
	importFilePath string
	importUrlPath  string
)

var export = &cobra.Command{
	Use:   "export",
	Short: "Export all snippets",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logger.New()
		db, cleanup, err := database.New()
		defer cleanup()
		if err != nil {
			return log.Error("database.New, err=%v", err)
		}

		snippets, err := db.ReadAll()
		if err != nil {
			return log.Error("db.ReadAll, err=%v", err)
		}

		log.Info("Generating report with %d snippets...", len(snippets))

		filename := filepath.Join(exportPath, fmt.Sprintf("snip-%s.csv", time.Now().Format("2006-01-02")))
		file, err := os.Create(filename)
		if err != nil {
			return log.Error("os.Create, err=%v", err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		if err := writer.Write([]string{"Key", "Value"}); err != nil {
			return fmt.Errorf("writer.Write, err=%v", err)
		}

		for _, snippet := range snippets {
			if err := writer.Write([]string{snippet.Key, snippet.Value}); err != nil {
				return fmt.Errorf("writer.Write, err=%v", err)
			}
		}

		log.Info("CSV file successfully created at path=%q", filename)
		return nil
	},
}

var importc = &cobra.Command{
	Use:   "import",
	Short: "Import snippets",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logger.New()

		if importFilePath == "" && importUrlPath == "" {
			return log.Error("a valid path or url is required, path=%q, url=%q", importFilePath, importUrlPath)
		}

		var reader io.Reader
		if importFilePath != "" {
			file, err := os.Open(importFilePath)
			if err != nil {
				return log.Error("os.Open, err=%v", err)
			}
			defer file.Close()
			reader = file
		} else {
			resp, err := http.Get(importUrlPath)
			if err != nil {
				return log.Error("http.Get, err=%v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return log.Error("resp.StatusCode, status=%v", resp.Status)
			}
			reader = resp.Body
		}

		csvReader := csv.NewReader(reader)
		records, err := csvReader.ReadAll()
		if err != nil {
			return log.Error("csvReader.ReadAll, err=%v", err)
		}

		var snippets []store.Snippet
		for _, record := range records[1:] {
			if len(record) < 2 {
				continue
			}
			snippet := store.Snippet{
				Key:   record[0],
				Value: record[1],
			}
			snippets = append(snippets, snippet)
		}

		if len(snippets) == 0 {
			return log.Error("no valid snippets where found")
		}

		db, cleanup, err := database.New()
		defer cleanup()
		if err != nil {
			return log.Error("database.New, err=%v", err)
		}

		if err := db.Import(snippets); err != nil {
			return log.Error("db.Import, err=%v", err)
		}

		source := importFilePath
		if source == "" {
			source = importUrlPath
		}
		log.Info("CSV file successfully imported from %q", source)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(ls)
	rootCmd.AddCommand(update)
	rootCmd.AddCommand(delete)
	rootCmd.AddCommand(reset)
	rootCmd.AddCommand(export)
	rootCmd.AddCommand(importc)
	export.Flags().StringVarP(&exportPath, "path", "p", ".", "Path to directory for CSV output")
	importc.Flags().StringVarP(&importFilePath, "path", "p", "", "Path to directory of your CSV file")
	importc.Flags().StringVarP(&importUrlPath, "url", "u", "", "URL of your remote CSV file")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
}
