package model

import "time"

type PasswordResetToken struct {
	ID        int
	UserID    int
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type PasswordResetTokenRepository interface {
	Create(userID int, tokenHash string, expiresAt time.Time) error
	FindByHash(tokenHash string) (*PasswordResetToken, error)
	DeleteByHash(tokenHash string) error
	DeleteByUserID(userID int) error
}
