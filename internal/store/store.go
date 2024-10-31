package store

type Snippet struct {
	Id    int
	Key   string
	Value string
}

type Store interface {
	Create(key, value string) error
	Read(id int) (Snippet, error)
	ReadAll() ([]Snippet, error)
	Update(snippet Snippet) error
	Delete(id int) error
}
