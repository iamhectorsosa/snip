package store

type Snippet struct {
	Id    int
	Key   string
	Value string
}

type Store interface {
	Create(key, value string) error
	Read(key string) (Snippet, error)
	ReadAll() ([]Snippet, error)
	Update(snippet Snippet) error
	Delete(key string) error
	Reset() error
	Import([]Snippet) error
}
