package database

import (
	"database/sql"
	"embed"
	"fmt"
	"os"

	"github.com/iamhectorsosa/snippets/internal/store"
	"github.com/pressly/goose/v3"
)

type Store struct {
	db *sql.DB
}

//go:embed sql/*.sql
var embedMigrations embed.FS

func New() (store *Store, cleanup func()) {
	db, err := sql.Open("libsql", "file:./local.db")
	if err != nil {
		fmt.Println("Error newDb:", err)
		os.Exit(1)
	}

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("turso"); err != nil {
		fmt.Println("Error goose:", err)
		os.Exit(1)
	}

	if err := goose.Up(db, "sql"); err != nil {
		fmt.Println("Error goose:", err)
		os.Exit(1)
	}

	cleanup = func() {
		db.Close()
		if err := goose.Down(db, "sql"); err != nil {
			fmt.Println("Error goose:", err)
			os.Exit(1)
		}
	}

	return &Store{db}, cleanup
}

func (s *Store) Create(name, text string) error {
	_, err := s.db.Exec(`INSERT INTO snippets (name, text) VALUES ($1, $2)`, name, text)
	return err
}

func (s *Store) Read(id int) (store.Snippet, error) {
	var snippet store.Snippet
	err := s.db.QueryRow("SELECT * FROM snippets WHERE id = ?", id).Scan(
		&snippet.Id,
		&snippet.Name,
		&snippet.Text,
	)
	return snippet, err
}

func (s *Store) ReadAll() ([]store.Snippet, error) {
	rows, err := s.db.Query("SELECT * FROM snippets")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var snippets []store.Snippet
	for rows.Next() {
		var snippet store.Snippet
		if err := rows.Scan(&snippet.Id, &snippet.Name, &snippet.Text); err != nil {
			return nil, err
		}
		snippets = append(snippets, snippet)
	}
	if err := rows.Err(); err != nil {
		return snippets, err
	}
	return snippets, nil
}

func (s *Store) Update(snippet store.Snippet) error {
	_, err := s.db.Exec(`UPDATE snippets SET id = $1, name = $2, text = $3 WHERE id = $1;)`, snippet.Id, snippet.Name, snippet.Text)
	return err
}

func (s *Store) Delete(id int) error {
	_, err := s.db.Exec("DELETE FROM snippets WHERE id = ?", id)
	return err
}
