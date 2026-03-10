package model

type Category struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Position int    `json:"position"`
	IsSystem bool   `json:"is_system"`
}

type Todo struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	Text       string    `json:"text"`
	Done       bool      `json:"done"`
	CategoryID *int      `json:"category_id"`
	Category   *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Position   int       `json:"position"`
}

type TodoRepository interface {
	FindAll(categoryID *int) ([]Todo, error)
	Create(text string, categoryID *int) (*Todo, error)
	Update(id int, text string, done bool, categoryID *int) (*Todo, error)
	Delete(id int) error
	DeleteDone(categoryID *int) error
	Reorder(ids []int) error
}

type CategoryRepository interface {
	FindAll() ([]Category, error)
	Create(name string) (*Category, error)
	Update(id int, name string) (*Category, error)
	Delete(id int) error
	Reorder(ids []int) error
}
