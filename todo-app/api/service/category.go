package service

import "todo-api/model"

type CategoryService struct {
	repo model.CategoryRepository
}

func NewCategoryService(repo model.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll() ([]model.Category, error) {
	return s.repo.FindAll()
}

func (s *CategoryService) Create(name string) (*model.Category, error) {
	return s.repo.Create(name)
}

func (s *CategoryService) Update(id int, name string) (*model.Category, error) {
	return s.repo.Update(id, name)
}

func (s *CategoryService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *CategoryService) Reorder(ids []int) error {
	return s.repo.Reorder(ids)
}
