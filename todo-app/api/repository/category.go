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

func (r *categoryRepository) FindAll(userID int) ([]model.Category, error) {
	var categories []model.Category
	result := r.db.Where("user_id = ?", userID).Order("position, id").Find(&categories)
	return categories, result.Error
}

func (r *categoryRepository) Create(userID int, name string) (*model.Category, error) {
	var maxPos int
	r.db.Model(&model.Category{}).Where("user_id = ?", userID).Select("COALESCE(MAX(position), -1)").Scan(&maxPos)
	category := &model.Category{UserID: userID, Name: name, Position: maxPos + 1}
	result := r.db.Create(category)
	return category, result.Error
}

func (r *categoryRepository) Update(userID, id int, name string) (*model.Category, error) {
	var category model.Category
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&category).Error; err != nil {
		return nil, err
	}
	category.Name = name
	if err := r.db.Save(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Delete(userID, id int) error {
	var category model.Category
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&category).Error; err != nil {
		return err
	}
	if category.IsSystem {
		return errors.New("cannot delete system category")
	}
	return r.db.Where("user_id = ?", userID).Delete(&model.Category{}, id).Error
}

func (r *categoryRepository) Reorder(userID int, ids []int) error {
	for i, id := range ids {
		if err := r.db.Model(&model.Category{}).Where("id = ? AND user_id = ?", id, userID).Update("position", i).Error; err != nil {
			return err
		}
	}
	return nil
}
