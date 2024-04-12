package model

import (
	"context"
)

type Chat struct {
	ChatID     uint64 `gorm:"primaryKey;autoIncrement;not null" json:"chatID"`
	UserID     uint64 `gorm:"foreignKey;not null" json:"userID"`
	ReceiverID uint64 `gorm:"foreignKey;not null" json:"receiverID"`
}

type ChatInterface struct {
	ID         *uint64
	UserID     *uint64
	ReceiverID *uint64
}

type UserChatRepository interface {
	Create(ctx context.Context, userChat Chat) error
	Get(ctx context.Context, ci ChatInterface) ([]Chat, error)
	Update(ctx context.Context, userChat Chat) error
	Delete(ctx context.Context, ci ChatInterface) error
}
