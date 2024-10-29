package store

type Snippet struct {
	Id   int
	Name string
	Text string
}

type Store interface {
	Create(name, text string) error
	Read(id int) (Snippet, error)
	ReadAll() ([]Snippet, error)
	Update(snippet Snippet) error
	Delete(id int) error
}
