package service

import "todo-api/model"

type TodoService struct {
	repo model.TodoRepository
}

func NewTodoService(repo model.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) GetAll(categoryID *int) ([]model.Todo, error) {
	return s.repo.FindAll(categoryID)
}

func (s *TodoService) Create(text string, categoryID *int) (*model.Todo, error) {
	return s.repo.Create(text, categoryID)
}

func (s *TodoService) Update(id int, text string, done bool, categoryID *int) (*model.Todo, error) {
	return s.repo.Update(id, text, done, categoryID)
}

func (s *TodoService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *TodoService) DeleteDone(categoryID *int) error {
	return s.repo.DeleteDone(categoryID)
}

func (s *TodoService) Reorder(ids []int) error {
	return s.repo.Reorder(ids)
}
