package model

import "time"

type User struct {
	ID            int       `json:"id" gorm:"primaryKey"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"-" gorm:"column:email_verified"`
	PasswordHash  string    `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
}

type UserRepository interface {
	FindByEmail(email string) (*User, error)
	Create(email, passwordHash string) (*User, error)
	SetEmailVerified(userID int) error
	UpdatePassword(userID int, passwordHash string) error
}
