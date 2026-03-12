package model

import "time"

type EmailVerificationToken struct {
	ID        int
	UserID    int
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type EmailVerificationTokenRepository interface {
	Create(userID int, tokenHash string, expiresAt time.Time) error
	FindByHash(tokenHash string) (*EmailVerificationToken, error)
	DeleteByHash(tokenHash string) error
}
