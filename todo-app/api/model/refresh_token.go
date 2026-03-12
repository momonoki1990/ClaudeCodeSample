package model

import "time"

type RefreshToken struct {
	ID        int       `gorm:"primaryKey"`
	UserID    int
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type RefreshTokenRepository interface {
	Create(userID int, tokenHash string, expiresAt time.Time) error
	FindByHash(tokenHash string) (*RefreshToken, error)
	DeleteByHash(tokenHash string) error
	DeleteByUserID(userID int) error
}
