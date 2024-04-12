package model

import (
	"context"
	"time"
)

type User struct {
	UserID         uint64    `gorm:"primaryKey;autoIncrement;not null" json:"user_id"`
	Name           string    `gorm:"type:varchar(255);not null" json:"name"`
	Username       string    `gorm:"type:varchar(255);not null;unique" json:"username"`
	Password       string    `gorm:"type:varchar(255);not null" json:"-"`
	Phone          string    `gorm:"type:varchar(255);not null;unique" json:"phone"`
	IsActive       string    `gorm:"type:varchar(10);not null;" json:"is_active"`
	IsFirtsLogin   string    `gorm:"type:varchar(10);" json:"is_first_login"`
	Biography      string    `gorm:"type:varchar(255)" json:"biography"`
	ProfilePicture string    `gorm:"type:varchar(255)" json:"profilePicture"`
	LastSeen       time.Time `json:"lastSeen"`
}

type UserInterface struct {
	ID       *uint64
	Username *string
	Phone    *string
	IsActive *string
}

type UserRepository interface {
	Create(ctx context.Context, user User) error
	Get(ctx context.Context, ui UserInterface) ([]User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, ui UserInterface) error
}
