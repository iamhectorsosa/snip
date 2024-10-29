package main

import (
	"database/sql"
	"fmt"

	"github.com/iamhectorsosa/snippets/internal/database"
	"github.com/iamhectorsosa/snippets/internal/store"
	_ "github.com/tursodatabase/go-libsql"
)

type Snippet struct {
	Id   int
	Name string
	Text string
}

type apiConfig struct {
	store store.Store
}

func newDb(filename string) (*sql.DB, error) {
	return sql.Open("libsql", "file:./"+filename+".db")
}

func main() {
	db, cleanup := database.New()
	defer cleanup()

	err := db.Create("Unwanted snippet", "echo \"This snippet should be deleted\"")
	if err != nil {
		fmt.Println("Error create:", err)
	}

	snippet, err := db.Read(2)
	if err != nil {
		fmt.Println("Error read:", err)
	}

	err = db.Update(store.Snippet{
		Id:   snippet.Id,
		Name: fmt.Sprintf("Updated %s", snippet.Name),
		Text: snippet.Text,
	})
	if err != nil {
		fmt.Println("Error update:", err)
	}

	err = db.Delete(3)
	if err != nil {
		fmt.Println("Error delete:", err)
	}

	snippets, err := db.ReadAll()
	if err != nil {
		fmt.Println("Error readAll:", err)
	}

	for _, snippet := range snippets {
		fmt.Println(snippet.Id, snippet.Name, snippet.Text)
	}
}
