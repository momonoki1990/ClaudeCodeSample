package model

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type TodoRepository interface {
	FindAll() ([]Todo, error)
	Create(text string) (*Todo, error)
	Update(id int, done bool) (*Todo, error)
	Delete(id int) error
}
