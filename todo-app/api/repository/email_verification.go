package repository

import (
	"time"
	"todo-api/model"

	"gorm.io/gorm"
)

type emailVerificationTokenRepository struct {
	db *gorm.DB
}

func NewEmailVerificationTokenRepository(db *gorm.DB) model.EmailVerificationTokenRepository {
	return &emailVerificationTokenRepository{db: db}
}

func (r *emailVerificationTokenRepository) Create(userID int, tokenHash string, expiresAt time.Time) error {
	return r.db.Exec(
		"INSERT INTO email_verification_tokens (user_id, token_hash, expires_at) VALUES (?, ?, ?)",
		userID, tokenHash, expiresAt,
	).Error
}

func (r *emailVerificationTokenRepository) FindByHash(tokenHash string) (*model.EmailVerificationToken, error) {
	var evt model.EmailVerificationToken
	err := r.db.Raw(
		"SELECT id, user_id, token_hash, expires_at, created_at FROM email_verification_tokens WHERE token_hash = ? AND expires_at > NOW()",
		tokenHash,
	).Scan(&evt).Error
	if err != nil {
		return nil, err
	}
	if evt.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &evt, nil
}

func (r *emailVerificationTokenRepository) DeleteByHash(tokenHash string) error {
	return r.db.Exec(
		"DELETE FROM email_verification_tokens WHERE token_hash = ?",
		tokenHash,
	).Error
}
