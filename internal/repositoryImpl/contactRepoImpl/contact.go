package contactRepoImpl

import (
	"backend/internal/model"
	"context"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type ContactDTO struct {
	model.Contact
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (c *ContactDTO) ToContact() *model.Contact {
	return &model.Contact{
		ContactID:     c.ContactID,
		UserID:        c.UserID,
		ContactUserID: c.ContactUserID,
		Status:        c.Status,
	}
}

func ToContactDTO(contact model.Contact) *ContactDTO {
	return &ContactDTO{
		Contact: model.Contact{
			ContactID:     contact.ContactID,
			UserID:        contact.UserID,
			ContactUserID: contact.ContactUserID,
			Status:        contact.Status,
		},
		CreatedAt: time.Now(),
	}
}

func (c *Repository) Create(ctx context.Context, contact model.Contact) error {
	contactDTO := ToContactDTO(contact)

	result := c.db.WithContext(ctx).Create(contactDTO)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (c *Repository) Get(ctx context.Context, ci model.ContactInterface) ([]model.Contact, error) {
	var contactDTOs []ContactDTO
	var condition ContactDTO

	if ci.ID != nil {
		condition.ContactID = *ci.ID
	}
	if ci.UserID != nil {
		condition.UserID = *ci.UserID
	}
	if ci.Status != nil {
		condition.Status = *ci.Status
	}

	result := c.db.WithContext(ctx).Where(&condition).Find(&contactDTOs)
	if result.Error != nil {
		return nil, result.Error
	}

	contacts := make([]model.Contact, len(contactDTOs))
	for i, contactDTO := range contactDTOs {
		contacts[i] = *contactDTO.ToContact()
	}

	return contacts, nil
}

func (c *Repository) Update(ctx context.Context, contact model.Contact) error {
	var condition ContactDTO
	condition.ContactID = contact.ContactID
	condition.UserID = contact.UserID
	condition.ContactUserID = contact.ContactUserID

	dto := ContactDTO{
		Contact:   contact,
		UpdatedAt: time.Now(),
	}

	result := c.db.WithContext(ctx).Where(&condition).Updates(&dto)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (c *Repository) Delete(ctx context.Context, ci model.ContactInterface) error {
	var condition ContactDTO

	if ci.ID != nil {
		condition.ContactID = *ci.ID
	}
	if ci.UserID != nil {
		condition.UserID = *ci.UserID
	}
	if ci.UserID != nil {
		condition.ContactUserID = *ci.ContactUserID
	}
	if ci.Status != nil {
		condition.Status = *ci.Status
	}

	result := c.db.WithContext(ctx).Where(&condition).Delete(&ContactDTO{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
