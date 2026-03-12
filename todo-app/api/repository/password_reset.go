package repository

import (
	"time"
	"todo-api/model"

	"gorm.io/gorm"
)

type passwordResetTokenRepository struct {
	db *gorm.DB
}

func NewPasswordResetTokenRepository(db *gorm.DB) model.PasswordResetTokenRepository {
	return &passwordResetTokenRepository{db: db}
}

func (r *passwordResetTokenRepository) Create(userID int, tokenHash string, expiresAt time.Time) error {
	return r.db.Exec(
		"INSERT INTO password_reset_tokens (user_id, token_hash, expires_at) VALUES (?, ?, ?)",
		userID, tokenHash, expiresAt,
	).Error
}

func (r *passwordResetTokenRepository) FindByHash(tokenHash string) (*model.PasswordResetToken, error) {
	var prt model.PasswordResetToken
	err := r.db.Raw(
		"SELECT id, user_id, token_hash, expires_at, created_at FROM password_reset_tokens WHERE token_hash = ? AND expires_at > NOW()",
		tokenHash,
	).Scan(&prt).Error
	if err != nil {
		return nil, err
	}
	if prt.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &prt, nil
}

func (r *passwordResetTokenRepository) DeleteByHash(tokenHash string) error {
	return r.db.Exec(
		"DELETE FROM password_reset_tokens WHERE token_hash = ?",
		tokenHash,
	).Error
}

func (r *passwordResetTokenRepository) DeleteByUserID(userID int) error {
	return r.db.Exec(
		"DELETE FROM password_reset_tokens WHERE user_id = ?",
		userID,
	).Error
}
