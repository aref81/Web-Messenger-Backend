package groupChatRepoImpl

import (
	"backend/internal/model"
	"context"
	"gorm.io/gorm"
	"time"
)

type Repository struct {
	db *gorm.DB
}

type GroupChatDTO struct {
	model.GroupChat
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (u *GroupChatDTO) ToGroupChat() *model.GroupChat {
	return &model.GroupChat{
		GroupChatID: u.GroupChatID,
		GroupID:     u.GroupID,
		MessageID:   u.MessageID,
	}
}

func ToGroupChatDTO(groupChat model.GroupChat) *GroupChatDTO {
	return &GroupChatDTO{
		GroupChat: model.GroupChat{
			GroupChatID: groupChat.GroupChatID,
			GroupID:     groupChat.GroupID,
			MessageID:   groupChat.MessageID,
		},
		CreatedAt: time.Now(),
	}
}

func (u *Repository) Create(ctx context.Context, groupChat model.GroupChat) error {
	groupChatDTO := ToGroupChatDTO(groupChat)

	result := u.db.WithContext(ctx).Create(groupChatDTO)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (u *Repository) Get(ctx context.Context, gci model.GroupChatInterface) ([]model.GroupChat, error) {
	var GroupChatDTOs []GroupChatDTO
	var condition GroupChatDTO

	if gci.ID != nil {
		condition.GroupChatID = *gci.ID
	}
	if gci.GroupID != nil {
		condition.GroupID = *gci.GroupID
	}

	result := u.db.WithContext(ctx).Where(&condition).Find(&GroupChatDTOs)
	if result.Error != nil {
		return nil, result.Error
	}

	groupChats := make([]model.GroupChat, len(GroupChatDTOs))
	for i, userDTO := range GroupChatDTOs {
		groupChats[i] = *userDTO.ToGroupChat()
	}

	return groupChats, nil
}

func (u *Repository) Update(ctx context.Context, groupChat model.GroupChat) error {
	var condition GroupChatDTO
	condition.GroupChatID = groupChat.GroupChatID

	dto := GroupChatDTO{
		GroupChat: groupChat,
		UpdatedAt: time.Now(),
	}

	result := u.db.WithContext(ctx).Where(&condition).Updates(&dto)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (u *Repository) Delete(ctx context.Context, gci model.GroupChatInterface) error {
	var condition GroupChatDTO

	if gci.ID != nil {
		condition.GroupChatID = *gci.ID
	}
	if gci.GroupID != nil {
		condition.GroupID = *gci.GroupID
	}

	result := u.db.WithContext(ctx).Where(&condition).Delete(&GroupChatDTO{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
