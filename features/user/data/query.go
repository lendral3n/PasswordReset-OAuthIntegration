package data

import (
	"context"
	"emailnotifl3n/app/cache"
	"emailnotifl3n/features/user"
	"errors"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type userQuery struct {
	db    *gorm.DB
	redis cache.Redis
}

func New(db *gorm.DB, redis cache.Redis) user.UserDataInterface {
	return &userQuery{
		db:    db,
		redis: redis,
	}
}

// Insert implements user.UserDataInterface.
func (repo *userQuery) Insert(input user.Core) error {
	dataGorm := CoreToModel(input)

	tx := repo.db.Create(&dataGorm)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("insert failed, row affected = 0")
	}
	return nil
}

// SelectById implements user.UserDataInterface.
func (repo *userQuery) SelectById(userId int) (*user.Core, error) {
	var userDataGorm User
	tx := repo.db.First(&userDataGorm, userId)
	if tx.Error != nil {
		return nil, tx.Error
	}

	result := userDataGorm.ModelToCore()
	return &result, nil
}

// Update implements user.UserDataInterface.
func (repo *userQuery) Update(userId int, input user.CoreUpdate) error {
	dataGorm := CoreToModelUpdate(input)
	tx := repo.db.Model(&User{}).Where("id = ?", userId).Updates(dataGorm)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("error record not found ")
	}
	return nil
}

// Delete implements user.UserDataInterface.
func (repo *userQuery) Delete(userId int) error {
	tx := repo.db.Delete(&User{}, userId)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("error record not found")
	}
	return nil
}

// Login implements user.UserDataInterface.
func (repo *userQuery) Login(email string, password string) (data *user.Core, err error) {
	var userGorm User
	tx := repo.db.Where("email = ?", email).First(&userGorm)
	if tx.Error != nil {
		// return nil, tx.Error
		return nil, errors.New(" Invalid email or password")
	}
	result := userGorm.ModelToCore()
	return &result, nil
}

// ChangePassword implements user.UserDataInterface.
func (repo *userQuery) ChangePassword(userId int, oldPassword, newPassword string) error {
	var userGorm User
	userGorm.Password = newPassword
	tx := repo.db.Model(&User{}).Where("id = ?", userId).Updates(&userGorm)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("error record not found ")
	}
	return nil
}

// ResetPassword implements user.UserDataInterface.
func (repo *userQuery) ResetPasswordLink(userId int, newPassword string) error {
	var userGorm User
	userGorm.Password = newPassword

	tx := repo.db.Model(&User{}).Where("id = ?", userId).Updates(&userGorm)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// SelectByEmail implements user.UserDataInterface.
func (repo *userQuery) SelectByEmail(email string) (*user.Core, error) {
	var userGorm User
	tx := repo.db.Where(" email = ?", email).First(&userGorm)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("email tidak ada")
		}
		return nil, tx.Error
	}

	result := userGorm.ModelToCore()
	return &result, nil
}

// VerifyEmailLink implements user.UserDataInterface.
func (repo *userQuery) VerifyEmailLink(userId int, verification bool) error {
	var userGorm User
	userGorm.Verified = verification

	tx := repo.db.Model(&User{}).Where("id = ?", userId).Updates(&userGorm)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// CreateCode implements user.UserDataInterface.
func (repo *userQuery) CreateCode(email, code string) error {
	ctx := context.Background()
	
	err := repo.redis.Delete(ctx, email)
	if err != nil {
		return err
	}
	
	err = repo.redis.Set(ctx, email, code)
	return err
}

// CheckCode implements user.UserDataInterface.
func (repo *userQuery) CheckCode(email string) (bool, error) {
	ctx := context.Background()
	_, err := repo.redis.Get(ctx, email)
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// VerifyEmailCode implements user.UserDataInterface.
func (repo *userQuery) VerifyEmailCode(email string, code string) error {
	panic("unimplemented")
}

// ResetPasswordCode implements user.UserDataInterface.
func (repo userQuery) ResetPasswordCode(newPassword string) error {
	panic("unimplemented")
}
