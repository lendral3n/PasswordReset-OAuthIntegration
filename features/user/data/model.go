package data

import (
	"emailnotifl3n/features/user"

	"gorm.io/gorm"
)

// struct user gorm model
type User struct {
	gorm.Model
	Name             string `gorm:"not null"`
	Email            string `gorm:"unique"`
	Password         string `gorm:"not null"`
	PhotoProfile     string
	Verified         bool
	RegistrationType string
}

func CoreToModel(input user.Core) User {
	return User{
		Name:             input.Name,
		Email:            input.Email,
		Password:         input.Password,
		PhotoProfile:     input.PhotoProfile,
		Verified:         input.Verified,
		RegistrationType: input.RegistrationType,
	}
}

func CoreToModelUpdate(input user.CoreUpdate) User {
	return User{
		Name:         input.Name,
		Email:        input.Email,
		PhotoProfile: input.PhotoProfile,
	}
}

func (u User) ModelToCore() user.Core {
	return user.Core{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		Password:     u.Password,
		PhotoProfile: u.PhotoProfile,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
