package userRepoImpl

import (
	"backend/internal/model"
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type UserDTO struct {
	model.User
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (u *UserDTO) ToUser() *model.User {
	return &model.User{
		Name:           u.Name,
		UserID:         u.UserID,
		Username:       u.Username,
		Password:       u.Password,
		Phone:          u.Phone,
		IsActive:       u.IsActive,
		Biography:      u.Biography,
		ProfilePicture: u.ProfilePicture,
		IsFirtsLogin:   u.IsFirtsLogin,
	}
}

func ToUserDTO(user model.User) *UserDTO {
	return &UserDTO{
		User: model.User{
			Name:           user.Name,
			UserID:         user.UserID,
			Username:       user.Username,
			Password:       user.Password,
			Phone:          user.Phone,
			IsActive:       user.IsActive,
			Biography:      user.Biography,
			ProfilePicture: user.ProfilePicture,
			IsFirtsLogin:   user.IsFirtsLogin,
		},
		CreatedAt: time.Now(),
	}
}

func (u *Repository) Create(ctx context.Context, user model.User) error {
	userDTO := ToUserDTO(user)

	result := u.db.WithContext(ctx).Create(userDTO)
	if result.Error != nil {
		log.Warnln(result.Error)
		return result.Error
	}

	return nil
}

func (u *Repository) Get(ctx context.Context, ui model.UserInterface) ([]model.User, error) {
	var userDTOs []UserDTO
	var condition UserDTO

	if ui.ID != nil {
		condition.UserID = *ui.ID
	}
	if ui.Username != nil {
		condition.Username = *ui.Username
	}
	if ui.Phone != nil {
		condition.Phone = *ui.Phone
	}
	if ui.IsActive != nil {
		condition.IsActive = *ui.IsActive
	}

	result := u.db.WithContext(ctx).Where(&condition).Find(&userDTOs)
	if result.Error != nil {
		return nil, result.Error
	}

	users := make([]model.User, len(userDTOs))
	for i, userDTO := range userDTOs {
		users[i] = *userDTO.ToUser()
	}

	return users, nil
}

func (u *Repository) Update(ctx context.Context, user model.User) error {
	var condition UserDTO
	condition.Username = user.Username
	condition.UserID = user.UserID
	user.IsFirtsLogin = "false"

	dto := UserDTO{
		User:      user,
		UpdatedAt: time.Now(),
	}

	result := u.db.WithContext(ctx).Where(&condition).Updates(&dto)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (u *Repository) Delete(ctx context.Context, ui model.UserInterface) error {
	var condition UserDTO

	if ui.ID != nil {
		condition.UserID = *ui.ID
	}
	if ui.Username != nil {
		condition.Username = *ui.Username
	}
	if ui.Phone != nil {
		condition.Phone = *ui.Phone
	}
	if ui.IsActive != nil {
		condition.IsActive = *ui.IsActive
	}

	result := u.db.WithContext(ctx).Where(&condition).Delete(&UserDTO{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
