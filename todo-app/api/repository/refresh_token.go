package repository

import (
	"time"
	"todo-api/model"

	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) model.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(userID int, tokenHash string, expiresAt time.Time) error {
	rt := &model.RefreshToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}
	return r.db.Create(rt).Error
}

func (r *refreshTokenRepository) FindByHash(tokenHash string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	err := r.db.Where("token_hash = ? AND expires_at > NOW()", tokenHash).First(&rt).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *refreshTokenRepository) DeleteByHash(tokenHash string) error {
	return r.db.Where("token_hash = ?", tokenHash).Delete(&model.RefreshToken{}).Error
}

func (r *refreshTokenRepository) DeleteByUserID(userID int) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.RefreshToken{}).Error
}
