package database

import (
	"database/sql"

	"github.com/iamhectorsosa/snippets/internal/store"

	_ "github.com/tursodatabase/go-libsql"
)

type Store struct {
	db *sql.DB
}

func New() (store *Store, cleanup func() error, err error) {
	db, err := sql.Open("libsql", "file:./local.db")
	cleanup = func() error {
		if err := db.Close(); err != nil {
			return err
		}
		return nil
	}

	if err != nil {
		return nil, cleanup, err
	}

	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS snippets (id INTEGER PRIMARY KEY, name TEXT UNIQUE, text TEXT)`); err != nil {
		return nil, cleanup, err
	}

	return &Store{db}, cleanup, nil
}

func (s *Store) Create(name, text string) error {
	_, err := s.db.Exec(`INSERT INTO snippets (name, text) VALUES (?, ?)`, name, text)
	return err
}

func (s *Store) Read(name string) (store.Snippet, error) {
	var snippet store.Snippet
	err := s.db.QueryRow("SELECT * FROM snippets WHERE name = ?", name).Scan(
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
	_, err := s.db.Exec(`UPDATE snippets SET name = ?, text = ? WHERE id = ?`, snippet.Name, snippet.Text, snippet.Id)
	return err
}

func (s *Store) Delete(name string) error {
	_, err := s.db.Exec("DELETE FROM snippets WHERE name = ?", name)
	return err
}
