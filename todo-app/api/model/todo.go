package model

type Category struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	UserID   int    `json:"user_id"`
	Name     string `json:"name"`
	Position int    `json:"position"`
	IsSystem bool   `json:"is_system"`
}

type Todo struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	UserID     int       `json:"user_id"`
	Text       string    `json:"text"`
	Done       bool      `json:"done"`
	CategoryID *int      `json:"category_id"`
	Category   *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Position   int       `json:"position"`
}

type TodoRepository interface {
	FindAll(userID int, categoryID *int) ([]Todo, error)
	Create(userID int, text string, categoryID *int) (*Todo, error)
	Update(userID, id int, text string, done bool, categoryID *int) (*Todo, error)
	Delete(userID, id int) error
	DeleteDone(userID int, categoryID *int) error
	Reorder(userID int, ids []int) error
}

type CategoryRepository interface {
	FindAll(userID int) ([]Category, error)
	Create(userID int, name string) (*Category, error)
	Update(userID, id int, name string) (*Category, error)
	Delete(userID, id int) error
	Reorder(userID int, ids []int) error
}
