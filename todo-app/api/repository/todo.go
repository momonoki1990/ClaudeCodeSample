package repository

import (
	"todo-api/model"

	"gorm.io/gorm"
)

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) model.TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) FindAll(userID int, categoryID *int) ([]model.Todo, error) {
	var todos []model.Todo
	q := r.db.Preload("Category").Where("user_id = ?", userID).Order("position, id")
	if categoryID != nil {
		q = q.Where("category_id = ?", *categoryID)
	}
	return todos, q.Find(&todos).Error
}

func (r *todoRepository) Create(userID int, text string, categoryID *int) (*model.Todo, error) {
	var minPos int
	r.db.Model(&model.Todo{}).Where("user_id = ?", userID).Select("COALESCE(MIN(position), 1)").Scan(&minPos)
	todo := &model.Todo{UserID: userID, Text: text, Done: false, CategoryID: categoryID, Position: minPos - 1}
	if err := r.db.Create(todo).Error; err != nil {
		return nil, err
	}
	r.db.Preload("Category").First(todo, todo.ID)
	return todo, nil
}

func (r *todoRepository) Update(userID, id int, text string, done bool, categoryID *int) (*model.Todo, error) {
	var todo model.Todo
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&todo).Error; err != nil {
		return nil, err
	}
	todo.Text = text
	todo.Done = done
	todo.CategoryID = categoryID
	if err := r.db.Save(&todo).Error; err != nil {
		return nil, err
	}
	r.db.Preload("Category").First(&todo, id)
	return &todo, nil
}

func (r *todoRepository) Delete(userID, id int) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.Todo{}, id).Error
}

func (r *todoRepository) DeleteDone(userID int, categoryID *int) error {
	q := r.db.Where("done = ? AND user_id = ?", true, userID)
	if categoryID != nil {
		q = q.Where("category_id = ?", *categoryID)
	}
	return q.Delete(&model.Todo{}).Error
}

func (r *todoRepository) Reorder(userID int, ids []int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i, id := range ids {
			if err := tx.Model(&model.Todo{}).Where("id = ? AND user_id = ?", id, userID).Update("position", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
