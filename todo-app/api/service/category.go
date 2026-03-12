package service

import "todo-api/model"

type CategoryService struct {
	repo model.CategoryRepository
}

func NewCategoryService(repo model.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll(userID int) ([]model.Category, error) {
	return s.repo.FindAll(userID)
}

func (s *CategoryService) Create(userID int, name string) (*model.Category, error) {
	return s.repo.Create(userID, name)
}

func (s *CategoryService) Update(userID, id int, name string) (*model.Category, error) {
	return s.repo.Update(userID, id, name)
}

func (s *CategoryService) Delete(userID, id int) error {
	return s.repo.Delete(userID, id)
}

func (s *CategoryService) Reorder(userID int, ids []int) error {
	return s.repo.Reorder(userID, ids)
}
