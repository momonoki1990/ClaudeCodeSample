package repository

import (
	"errors"
	"todo-api/model"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) model.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindAll() ([]model.Category, error) {
	var categories []model.Category
	result := r.db.Order("position, id").Find(&categories)
	return categories, result.Error
}

func (r *categoryRepository) Create(name string) (*model.Category, error) {
	var maxPos int
	r.db.Model(&model.Category{}).Select("COALESCE(MAX(position), -1)").Scan(&maxPos)
	category := &model.Category{Name: name, Position: maxPos + 1}
	result := r.db.Create(category)
	return category, result.Error
}

func (r *categoryRepository) Update(id int, name string) (*model.Category, error) {
	var category model.Category
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	category.Name = name
	if err := r.db.Save(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Delete(id int) error {
	var category model.Category
	if err := r.db.First(&category, id).Error; err != nil {
		return err
	}
	if category.IsSystem {
		return errors.New("cannot delete system category")
	}
	return r.db.Delete(&model.Category{}, id).Error
}

func (r *categoryRepository) Reorder(ids []int) error {
	for i, id := range ids {
		if err := r.db.Model(&model.Category{}).Where("id = ?", id).Update("position", i).Error; err != nil {
			return err
		}
	}
	return nil
}
