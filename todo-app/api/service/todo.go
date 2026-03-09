package service

import "todo-api/model"

type TodoService struct {
	repo model.TodoRepository
}

func NewTodoService(repo model.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) GetAll() ([]model.Todo, error) {
	return s.repo.FindAll()
}

func (s *TodoService) Create(text string) (*model.Todo, error) {
	return s.repo.Create(text)
}

func (s *TodoService) Update(id int, done bool) (*model.Todo, error) {
	return s.repo.Update(id, done)
}

func (s *TodoService) Delete(id int) error {
	return s.repo.Delete(id)
}
