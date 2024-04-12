package userChatRepoImpl

import (
	"backend/internal/model"
	"context"
	"gorm.io/gorm"
	"time"
)

type Repository struct {
	db *gorm.DB
}

type UserChatDTO struct {
	model.Chat
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (u *UserChatDTO) ToUserChat() *model.Chat {
	return &model.Chat{
		ChatID:     u.ChatID,
		UserID:     u.UserID,
		ReceiverID: u.ReceiverID,
	}
}

func ToUserChatDTO(userChat model.Chat) *UserChatDTO {
	return &UserChatDTO{
		Chat: model.Chat{
			ChatID:     userChat.ChatID,
			UserID:     userChat.UserID,
			ReceiverID: userChat.ReceiverID,
		},
		CreatedAt: time.Now(),
	}
}

func (u *Repository) Create(ctx context.Context, userChat model.Chat) error {
	userChatDTO := ToUserChatDTO(userChat)

	result := u.db.WithContext(ctx).Create(userChatDTO)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (u *Repository) Get(ctx context.Context, ci model.ChatInterface) ([]model.Chat, error) {
	var userChatDTOs []UserChatDTO
	var condition UserChatDTO

	if ci.ID != nil {
		condition.ChatID = *ci.ID
	}
	if ci.UserID != nil {
		condition.UserID = *ci.UserID
	}
	if ci.ReceiverID != nil {
		condition.ReceiverID = *ci.ReceiverID
	}

	result := u.db.WithContext(ctx).Where(&condition).Find(&userChatDTOs)
	if result.Error != nil {
		return nil, result.Error
	}

	userChats := make([]model.Chat, len(userChatDTOs))

	for i, userChatDTO := range userChatDTOs {
		userChats[i] = *userChatDTO.ToUserChat()
	}

	return userChats, nil
}

func (u *Repository) Update(ctx context.Context, userChat model.Chat) error {
	var condition UserChatDTO
	condition.ChatID = userChat.ChatID

	dto := UserChatDTO{
		Chat:      userChat,
		UpdatedAt: time.Now(),
	}

	result := u.db.WithContext(ctx).Where(&condition).Updates(&dto)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (u *Repository) Delete(ctx context.Context, ci model.ChatInterface) error {
	var condition UserChatDTO

	if ci.ID != nil {
		condition.ChatID = *ci.ID
	}
	if ci.UserID != nil {
		condition.UserID = *ci.UserID
	}
	if ci.ReceiverID != nil {
		condition.ReceiverID = *ci.ReceiverID
	}

	result := u.db.WithContext(ctx).Where(&condition).Delete(&UserChatDTO{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
