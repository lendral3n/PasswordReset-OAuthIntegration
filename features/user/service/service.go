package service

import (
	"emailnotifl3n/features/user"
	"emailnotifl3n/utils/encrypts"
	"emailnotifl3n/utils/middlewares"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
)

type userService struct {
	userData    user.UserDataInterface
	hashService encrypts.HashInterface
	validate    *validator.Validate
	m           sync.Map
}

// dependency injection
func New(repo user.UserDataInterface, hash encrypts.HashInterface) user.UserServiceInterface {
	return &userService{
		userData:    repo,
		hashService: hash,
		validate:    validator.New(),
	}
}

// Create implements user.UserServiceInterface.
func (service *userService) Create(input user.Core) error {
	errValidate := service.validate.Struct(input)
	if errValidate != nil {
		return errValidate
	}

	if input.Password != "" {
		hashedPass, errHash := service.hashService.HashPassword(input.Password)
		if errHash != nil {
			return errors.New("error hash password")
		}
		input.Password = hashedPass
	}

	err := service.userData.Insert(input)
	return err
}

// GetById implements user.UserServiceInterface.
func (service *userService) GetById(userId int) (*user.Core, error) {
	result, err := service.userData.SelectById(userId)
	return result, err
}

// Update implements user.UserServiceInterface.
func (service *userService) Update(userId int, input user.CoreUpdate) error {
	errValidate := service.validate.Struct(input)
	if errValidate != nil {
		return errValidate
	}
	if userId <= 0 {
		return errors.New("invalid id")
	}

	err := service.userData.Update(userId, input)
	return err
}

// Delete implements user.UserServiceInterface.
func (service *userService) Delete(userId int) error {
	if userId <= 0 {
		return errors.New("invalid id")
	}
	err := service.userData.Delete(userId)
	return err
}

// Login implements user.UserServiceInterface.
func (service *userService) Login(email string, password string) (data *user.Core, token string, err error) {
	if email == "" && password == "" {
		return nil, "", errors.New("email dan password wajib diisi")
	}
	if email == "" {
		return nil, "", errors.New("email wajib diisi")
	}
	if password == "" {
		return nil, "", errors.New("password wajib diisi")
	}

	data, err = service.userData.Login(email, password)
	if err != nil {
		return nil, "", err
	}

	isValid := service.hashService.CheckPasswordHash(data.Password, password)
	if !isValid {
		return nil, "", errors.New("password tidak sesuai")
	}

	token, errJwt := middlewares.CreateTokenLogin(int(data.ID))
	if errJwt != nil {
		return nil, "", errJwt
	}
	return data, token, err
}

// ChangePassword implements user.UserServiceInterface.
func (service *userService) ChangePassword(userId int, oldPassword, newPassword string) error {
	if oldPassword == "" {
		return errors.New("please input current password")
	}

	if newPassword == "" {
		return errors.New("please input new password")
	}

	hashedNewPass, errHash := service.hashService.HashPassword(newPassword)
	if errHash != nil {
		return errors.New("error hash password")
	}

	err := service.userData.ChangePassword(userId, oldPassword, hashedNewPass)
	return err

}

// ForgotPassword implements user.UserServiceInterface.
func (service *userService) ForgotPassword(email string) (data *user.Core, token string, err error) {
	user, err := service.userData.SelectByEmail(email)
	if err != nil {
		return nil, "", err
	}

	token, err = middlewares.CreateResetPasswordToken(int(user.ID))
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// ResetPassword implements user.UserServiceInterface.
func (service *userService) ResetPassword(userId int, newPassword string) error {
	hashedNewPass, errHash := service.hashService.HashPassword(newPassword)
	if errHash != nil {
		return errors.New("error hash password")
	}

	err := service.userData.ResetPasswordLink(userId, hashedNewPass)
	if err != nil {
		return err
	}
	return nil
}

// VerifyEmailLink implements user.UserServiceInterface.
func (service *userService) VerifyEmailLink(userId int) error {
	verification := true

	err := service.userData.VerifyEmailLink(userId, verification)
	if err != nil {
		return err
	}
	return nil
}

// RequestCode implements user.UserServiceInterface.
func (service *userService) RequestCode(email string, code string) (data *user.Core, err error) {
	if email == "" {
		return nil, errors.New("email harus di isi")
	}

	mail, err := service.userData.SelectByEmail(email)
	if err != nil {
		return nil, err
	}

	isValid, _ := service.userData.CheckCode(email)
	if isValid {
		if creationTime, ok := service.m.Load(email); ok {
			elapsed := time.Since(creationTime.(time.Time))
			if elapsed < 1*time.Minute {
				remaining := 1*time.Minute - elapsed
				return nil, fmt.Errorf("coba lagi minta kode dalam %.0f detik", remaining.Seconds())
			}
		}
	}
	if !isValid {
		err = service.userData.CreateCode(email, code)
		if err != nil {
			return nil, err
		}
		service.m.Store(email, time.Now())
		return mail, nil
	}
	return nil, errors.New("coba lagi")
}

// ResetPasswordCode implements user.UserServiceInterface.
func (*userService) ResetPasswordCode(newPassword string) error {
	panic("unimplemented")
}

// VerifyEmailCode implements user.UserServiceInterface.
func (*userService) VerifyEmailCode(email string, code string) error {
	panic("unimplemented")
}
