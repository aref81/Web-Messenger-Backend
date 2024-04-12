package model

import (
	"context"
	"time"
)

type MessageType string

const (
	TypePV MessageType = "PV"
	TypeGP MessageType = "GP"
)

type Message struct {
	MessageID uint64      `gorm:"primaryKey;autoIncrement;not null" json:"messageID"`
	ChatID    uint64      `gorm:"foreignKey;not null" json:"chatID"`
	SenderID  uint64      `gorm:"foreignKey;not null" json:"senderID"`
	Type      MessageType `gorm:"not null" json:"type"`
	Content   string      `gorm:"type:varchar(5000);not null" json:"content"`
	IsRead    string      `gorm:"type:varchar(10);" json:"isRead"`
}

type MessageInterface struct {
	ID       *uint64
	ChatID   *uint64
	SenderID *uint64
	Type     *MessageType
	IsRead   *string
}

type MessageRepository interface {
	Create(ctx context.Context, message Message) (uint64, error)
	Get(ctx context.Context, mi MessageInterface) ([]Message, error)
	Update(ctx context.Context, message Message) error
	Delete(ctx context.Context, mi MessageInterface) error
	GetDto(ctx context.Context, mi MessageInterface) ([]MessageDTO, error)
}

type MessageDTO struct {
	Message
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
