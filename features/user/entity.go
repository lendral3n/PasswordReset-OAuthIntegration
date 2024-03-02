package user

import (
	"time"
)

type Core struct {
	ID           uint
	Name         string `validate:"required"`
	Email        string `validate:"required,email"`
	Password     string `validate:"required"`
	PhotoProfile string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CoreUpdate struct {
	Name         string `validate:"required"`
	Email        string `validate:"required,email"`
	PhotoProfile string
}

// interface untuk Data Layer
type UserDataInterface interface {
	Insert(input Core) error
	SelectById(userId int) (*Core, error)
	Update(userId int, input CoreUpdate) error
	Delete(userId int) error
	Login(email, password string) (data *Core, err error)
	ChangePassword(userId int, oldPassword, newPassword string) error
}

// interface untuk Service Layer
type UserServiceInterface interface {
	Create(input Core) error
	GetById(userId int) (*Core, error)
	Update(userId int, input CoreUpdate) error
	Delete(userId int) error
	Login(email, password string) (data *Core, token string, err error)
	ChangePassword(userId int, oldPassword, newPassword string) error
}
