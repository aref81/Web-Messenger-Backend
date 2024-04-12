package groupRepoImpl

import (
	"backend/internal/model"
	"context"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func ToGroupDTO(group model.Group) *model.GroupDTO {
	return &model.GroupDTO{
		Group: model.Group{
			GroupID:     group.GroupID,
			Name:        group.Name,
			Description: group.Description,
			Creator:     group.Creator,
		},
		CreatedAt: time.Now(),
	}
}

func (g *Repository) Create(ctx context.Context, group model.Group) (uint64, error) {
	groupDTO := ToGroupDTO(group)

	result := g.db.WithContext(ctx).Create(groupDTO)
	if result.Error != nil {
		return 0, result.Error
	}

	return groupDTO.GroupID, nil
}

func (g *Repository) Get(ctx context.Context, gi model.GroupInterface) ([]model.Group, error) {
	var groupDTOs []model.GroupDTO
	var condition model.GroupDTO

	if gi.ID != nil {
		condition.GroupID = *gi.ID
	}
	if gi.Name != nil {
		condition.Name = *gi.Name
	}
	if gi.CreatorID != nil {
		condition.Creator = *gi.CreatorID
	}

	result := g.db.WithContext(ctx).Where(&condition).Find(&groupDTOs)
	if result.Error != nil {
		return nil, result.Error
	}

	groups := make([]model.Group, len(groupDTOs))
	for i, userDTO := range groupDTOs {
		groups[i] = *userDTO.ToGroup()
	}

	return groups, nil
}

func (g *Repository) Update(ctx context.Context, group model.Group) error {
	var condition model.GroupDTO
	condition.GroupID = group.GroupID

	dto := model.GroupDTO{
		Group:     group,
		UpdatedAt: time.Now(),
	}

	result := g.db.WithContext(ctx).Where(&condition).Updates(&dto)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (g *Repository) Delete(ctx context.Context, gi model.GroupInterface) error {
	var condition model.GroupDTO

	if gi.ID != nil {
		condition.GroupID = *gi.ID
	}
	if gi.Name != nil {
		condition.Name = *gi.Name
	}

	result := g.db.WithContext(ctx).Where(&condition).Delete(&model.GroupDTO{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
