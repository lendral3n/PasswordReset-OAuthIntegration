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
	Verified     bool
	RegistrationType string
	Code string
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
	SelectByEmail(email string) (*Core, error)
	ResetPasswordLink(userId int, newPassword string) error
	VerifyEmailLink(userId int, verification bool) error
	CreateCode(email, code string) error
	CheckCode(email string) (bool, error)
	DeleteCode(email string)error
	VerifyCode(email, code string) error
	VerifyEmailCode(email string, verification bool) error
	ResetPasswordCode(email, newPassword string) error
}

// interface untuk Service Layer
type UserServiceInterface interface {
	Create(input Core) error
	GetById(userId int) (*Core, error)
	Update(userId int, input CoreUpdate) error
	Delete(userId int) error
	Login(email, password string) (data *Core, token string, err error)
	ChangePassword(userId int, oldPassword, newPassword string) error
	ForgotPassword(email string) (data *Core, token string, err error)
	ResetPassword(userId int, newPassword string) error
	VerifyEmailLink(userId int) error
	SelectByEmail(email string) (*Core, error)
	RequestCode(email, code string) (data *Core, err error)
	VerifyEmailCode(email string, code string) error
	ResetPasswordCode(email, newPassword, code string) error
	RegisterGoogle(input Core) error
}