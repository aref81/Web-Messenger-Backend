package model

import (
	"context"
	"time"
)

type UserGroup struct {
	UserGroupID uint64 `gorm:"primaryKey;autoIncrement;not null" json:"userGroupID"`
	UserID      uint64 `gorm:"foreignKey;not null" json:"userID"`
	GroupID     uint64 `gorm:"foreignKey;not null" json:"groupID"`
}

type UserGroupInterface struct {
	ID      *uint64
	UserID  *uint64
	GroupID *uint64
}

type UserGroupDTO struct {
	UserGroup
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserGroupRepository interface {
	Create(ctx context.Context, userGroup UserGroup) error
	Get(ctx context.Context, ugi UserGroupInterface) ([]UserGroup, error)
	Update(ctx context.Context, userGroup UserGroup) error
	Delete(ctx context.Context, ugi UserGroupInterface) error
	GetGroupWithUserGroups(ctx context.Context, groupID uint64) ([]GroupDTO, []UserGroupDTO, error)
}
