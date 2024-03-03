package handler

import "emailnotifl3n/features/user"

type UserRequest struct {
	Name         string `json:"name" form:"name"`
	Email        string `json:"email" form:"email"`
	Password     string `json:"password" form:"password"`
	PhotoProfile string `json:"photo_profile" form:"photo_profile"`
}

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" form:"old_password"`
	NewPassword string `json:"new_password" form:"new_password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" form:"email"`
}

type ResetPasswordRequest struct {
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

func RequestToCore(input UserRequest) user.Core {
	return user.Core{
		Name:         input.Name,
		Email:        input.Email,
		Password:     input.Password,
		PhotoProfile: input.PhotoProfile,
	}
}

func UpdateRequestToCore(input UserRequest, imageURL string) user.Core {
	return user.Core{
		Name:         input.Name,
		Email:        input.Email,
		PhotoProfile: imageURL,
	}
}

func UpdateRequestToCoreUpdate(input UserRequest, imageURL string) user.CoreUpdate {
	return user.CoreUpdate{
		Name:         input.Name,
		Email:        input.Email,
		PhotoProfile: imageURL,
	}
}
