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

func (r *todoRepository) FindAll(categoryID *int) ([]model.Todo, error) {
	var todos []model.Todo
	q := r.db.Preload("Category").Order("position, id")
	if categoryID != nil {
		q = q.Where("category_id = ?", *categoryID)
	}
	return todos, q.Find(&todos).Error
}

func (r *todoRepository) Create(text string, categoryID *int) (*model.Todo, error) {
	todo := &model.Todo{Text: text, Done: false, CategoryID: categoryID}
	if err := r.db.Create(todo).Error; err != nil {
		return nil, err
	}
	r.db.Preload("Category").First(todo, todo.ID)
	return todo, nil
}

func (r *todoRepository) Update(id int, text string, done bool, categoryID *int) (*model.Todo, error) {
	var todo model.Todo
	if err := r.db.First(&todo, id).Error; err != nil {
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

func (r *todoRepository) Delete(id int) error {
	return r.db.Delete(&model.Todo{}, id).Error
}

func (r *todoRepository) DeleteDone(categoryID *int) error {
	q := r.db.Where("done = ?", true)
	if categoryID != nil {
		q = q.Where("category_id = ?", *categoryID)
	}
	return q.Delete(&model.Todo{}).Error
}

func (r *todoRepository) Reorder(ids []int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i, id := range ids {
			if err := tx.Model(&model.Todo{}).Where("id = ?", id).Update("position", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
