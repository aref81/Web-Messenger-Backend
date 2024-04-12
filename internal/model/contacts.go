package model

import (
	"context"
)

type Status string

const (
	Pending  Status = "pending"
	Accepted Status = "accepted"
	Blocked  Status = "blocked"
)

type Contact struct {
	ContactID     uint64 `gorm:"primaryKey;autoIncrement;not null" json:"contactID"`
	UserID        uint64 `gorm:"foreignKey;not null" json:"userID"`
	ContactUserID uint64 `gorm:"foreignKey;not null" json:"contactUserID"`
	Status        Status `gorm:"default:'pending'" json:"status"`
}

type ContactInterface struct {
	ID            *uint64
	UserID        *uint64
	ContactUserID *uint64
	Status        *Status
}

type ContactRepository interface {
	Create(ctx context.Context, contact Contact) error
	Get(ctx context.Context, ci ContactInterface) ([]Contact, error)
	Update(ctx context.Context, contact Contact) error
	Delete(ctx context.Context, ci ContactInterface) error
}
