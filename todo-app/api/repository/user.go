package repository

import (
	"todo-api/model"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) model.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) SetEmailVerified(userID int) error {
	return r.db.Exec("UPDATE users SET email_verified = 1 WHERE id = ?", userID).Error
}

func (r *userRepository) UpdatePassword(userID int, passwordHash string) error {
	return r.db.Exec("UPDATE users SET password_hash = ? WHERE id = ?", passwordHash, userID).Error
}

func (r *userRepository) Create(email, passwordHash string) (*model.User, error) {
	var user *model.User
	err := r.db.Transaction(func(tx *gorm.DB) error {
		user = &model.User{Email: email, PasswordHash: passwordHash}
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		category := &model.Category{
			UserID:   user.ID,
			Name:     "すべて",
			Position: 0,
			IsSystem: true,
		}
		return tx.Create(category).Error
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}
