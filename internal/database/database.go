package database

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/iamhectorsosa/snip/internal/store"

	_ "github.com/tursodatabase/go-libsql"
)

type Store struct {
	db *sql.DB
}

func New() (store *Store, cleanup func() error, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, nil, err
	}

	dbPath := filepath.Join(homeDir, ".config", "snip", "local.db")
	if err := os.MkdirAll(filepath.Dir(dbPath), os.ModePerm); err != nil {
		return nil, nil, err
	}

	db, err := sql.Open("libsql", "file:"+dbPath)
	if err != nil {
		return nil, nil, err
	}

	cleanup = db.Close
	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS snippets (id INTEGER PRIMARY KEY, key TEXT UNIQUE, value TEXT)`); err != nil {
		return nil, cleanup, err
	}

	return &Store{db}, cleanup, nil
}

func NewInMemory() (store *Store, cleanup func() error, err error) {
	db, err := sql.Open("libsql", ":memory:")
	if err != nil {
		return nil, nil, err
	}
	cleanup = db.Close

	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS snippets (id INTEGER PRIMARY KEY, key TEXT UNIQUE, value TEXT)`); err != nil {
		return nil, cleanup, err
	}
	return &Store{db}, cleanup, nil
}

func (s *Store) Create(key, value string) error {
	_, err := s.db.Exec(`INSERT INTO snippets (key, value) VALUES (?, ?)`, key, value)
	return err
}

func (s *Store) Read(key string) (store.Snippet, error) {
	var snippet store.Snippet
	err := s.db.QueryRow("SELECT * FROM snippets WHERE key = ?", key).Scan(
		&snippet.Id,
		&snippet.Key,
		&snippet.Value,
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
		if err := rows.Scan(&snippet.Id, &snippet.Key, &snippet.Value); err != nil {
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
	_, err := s.db.Exec(`UPDATE snippets SET key = ?, value = ? WHERE key = ?`, snippet.Key, snippet.Value, snippet.Key)
	return err
}

func (s *Store) Delete(Key string) error {
	_, err := s.db.Exec("DELETE FROM snippets WHERE key = ?", Key)
	return err
}

func (s *Store) Reset() error {
	_, err := s.db.Exec("DELETE FROM snippets")
	return err
}

func (s *Store) Import(snippets []store.Snippet) error {
	if len(snippets) == 0 {
		return nil
	}

	insertQuery := "INSERT OR IGNORE INTO snippets (key, value) VALUES "
	var args []interface{}
	for i, snippet := range snippets {
		if i > 0 {
			insertQuery += ", "
		}
		insertQuery += "(?, ?)"
		args = append(args, snippet.Key, snippet.Value)
	}

	_, err := s.db.Exec(insertQuery, args...)
	return err
}
