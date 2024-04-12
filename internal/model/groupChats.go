package model

import "context"

type GroupChat struct {
	GroupChatID uint64 `gorm:"primaryKey;auto_increment;not null" json:"groupChatID"`
	GroupID     uint64 `gorm:"foreignKey;not null" json:"groupID"`
	MessageID   uint64 `gorm:"foreignKey;not null" json:"messageID"`
}

type GroupChatInterface struct {
	ID      *uint64
	GroupID *uint64
}

type GroupChatRepository interface {
	Create(ctx context.Context, groupChat GroupChat) error
	Get(ctx context.Context, gci GroupChatInterface) ([]GroupChat, error)
	Update(ctx context.Context, groupChat GroupChat) error
	Delete(ctx context.Context, gci GroupChatInterface) error
}
