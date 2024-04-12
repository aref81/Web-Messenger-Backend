package model

import (
	"context"
	"time"
)

type Group struct {
	GroupID     uint64 `gorm:"primaryKey;auto_increment;not null" json:"groupID"`
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Description string `gorm:"type:varchar(255);not null" json:"description"`
	Creator     uint64 `gorm:"foreignKey;not null" json:"creator"`
}

type GroupInterface struct {
	ID        *uint64
	Name      *string
	CreatorID *uint64
}

type GroupDTO struct {
	Group
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GroupRepository interface {
	Get(ctx context.Context, gi GroupInterface) ([]Group, error)
	Create(ctx context.Context, group Group) (uint64, error)
	Update(ctx context.Context, group Group) error
	Delete(ctx context.Context, gi GroupInterface) error
}

func (g *GroupDTO) ToGroup() *Group {
	return &Group{
		GroupID:     g.GroupID,
		Name:        g.Name,
		Description: g.Description,
		Creator:     g.Creator,
	}
}
