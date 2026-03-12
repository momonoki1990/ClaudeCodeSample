package service

import "todo-api/model"

type TodoService struct {
	repo model.TodoRepository
}

func NewTodoService(repo model.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) GetAll(userID int, categoryID *int) ([]model.Todo, error) {
	return s.repo.FindAll(userID, categoryID)
}

func (s *TodoService) Create(userID int, text string, categoryID *int) (*model.Todo, error) {
	return s.repo.Create(userID, text, categoryID)
}

func (s *TodoService) Update(userID, id int, text string, done bool, categoryID *int) (*model.Todo, error) {
	return s.repo.Update(userID, id, text, done, categoryID)
}

func (s *TodoService) Delete(userID, id int) error {
	return s.repo.Delete(userID, id)
}

func (s *TodoService) DeleteDone(userID int, categoryID *int) error {
	return s.repo.DeleteDone(userID, categoryID)
}

func (s *TodoService) Reorder(userID int, ids []int) error {
	return s.repo.Reorder(userID, ids)
}
