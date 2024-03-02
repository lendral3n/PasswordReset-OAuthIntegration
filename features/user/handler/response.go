package handler

import "emailnotifl3n/features/user"

type UserResponse struct {
	ID           uint   `json:"id" form:"id"`
	Name         string `json:"name" form:"name"`
	Email        string `json:"email" form:"email"`
	PhotoProfile string `json:"photo_profile" form:"photo_profile"`
}

type UserKosDetailResponse struct {
	ID           uint   `json:"id" form:"id"`
	Name         string `json:"name" form:"name"`
	PhotoProfile string `json:"photo_profile" form:"photo_profile"`
}

func CoreToResponse(data *user.Core) UserResponse {
	var result = UserResponse{
		ID:    data.ID,
		Name:  data.Name,
		Email: data.Email,
		PhotoProfile: data.PhotoProfile,
	}
	return result
}
