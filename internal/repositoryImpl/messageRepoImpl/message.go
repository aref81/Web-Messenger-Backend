package messageRepoImpl

import (
	"backend/internal/model"
	"context"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type MessageDTO struct {
	model.Message
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (m *MessageDTO) ToMessage() *model.Message {
	return &model.Message{
		MessageID: m.MessageID,
		ChatID:    m.ChatID,
		SenderID:  m.SenderID,
		Content:   m.Content,
		Type:      m.Type,
		IsRead:    m.IsRead,
	}
}

func ToMessageDTO(message model.Message) *MessageDTO {
	return &MessageDTO{
		Message: model.Message{
			MessageID: message.MessageID,
			ChatID:    message.ChatID,
			SenderID:  message.SenderID,
			Content:   message.Content,
			Type:      message.Type,
			IsRead:    message.IsRead,
		},
		CreatedAt: time.Now(),
	}
}

func (m *Repository) Create(ctx context.Context, message model.Message) (uint64, error) {
	messageDTO := ToMessageDTO(message)

	result := m.db.WithContext(ctx).Create(messageDTO)
	if result.Error != nil {
		return 0, result.Error
	}

	return messageDTO.MessageID, nil
}

func (m *Repository) Get(ctx context.Context, mi model.MessageInterface) ([]model.Message, error) {
	var messageDTOs []MessageDTO
	var condition MessageDTO

	if mi.ID != nil {
		condition.MessageID = *mi.ID
	}
	if mi.ChatID != nil {
		condition.ChatID = *mi.ChatID
	}
	if mi.SenderID != nil {
		condition.SenderID = *mi.SenderID
	}
	if mi.Type != nil {
		condition.Type = *mi.Type
	}
	if mi.IsRead != nil {
		condition.IsRead = *mi.IsRead
	}

	result := m.db.WithContext(ctx).Where(&condition).Find(&messageDTOs)
	if result.Error != nil {
		return nil, result.Error
	}

	messages := make([]model.Message, len(messageDTOs))

	for i, message := range messageDTOs {
		messages[i] = *message.ToMessage()
	}

	return messages, nil
}

func (m *Repository) Update(ctx context.Context, message model.Message) error {
	var condition MessageDTO
	condition.MessageID = message.MessageID

	dto := MessageDTO{
		Message:   message,
		UpdatedAt: time.Now(),
	}

	result := m.db.WithContext(ctx).Where(&condition).Updates(&dto)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (m *Repository) Delete(ctx context.Context, mi model.MessageInterface) error {
	var condition MessageDTO

	if mi.ID != nil {
		condition.MessageID = *mi.ID
	}
	if mi.ChatID != nil {
		condition.ChatID = *mi.ChatID
	}
	if mi.SenderID != nil {
		condition.SenderID = *mi.SenderID
	}
	if mi.Type != nil {
		condition.Type = *mi.Type
	}

	result := m.db.WithContext(ctx).Where(&condition).Delete(&MessageDTO{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (m *Repository) GetDto(ctx context.Context, mi model.MessageInterface) ([]model.MessageDTO, error) {
	var messageDTOs []MessageDTO
	var condition MessageDTO

	if mi.ID != nil {
		condition.MessageID = *mi.ID
	}
	if mi.ChatID != nil {
		condition.ChatID = *mi.ChatID
	}
	if mi.SenderID != nil {
		condition.SenderID = *mi.SenderID
	}
	if mi.Type != nil {
		condition.Type = *mi.Type
	}

	result := m.db.WithContext(ctx).Where(&condition).Find(&messageDTOs)
	if result.Error != nil {
		return nil, result.Error
	}

	otherMessageDTO := make([]model.MessageDTO, len(messageDTOs))

	for i, message := range messageDTOs {
		m.Update(ctx, model.Message{
			MessageID: message.MessageID,
			IsRead:    "true",
		})
		otherMessageDTO[i] = model.MessageDTO{
			Message:   message.Message,
			CreatedAt: message.CreatedAt,
			UpdatedAt: message.UpdatedAt,
		}
	}

	return otherMessageDTO, nil
}
