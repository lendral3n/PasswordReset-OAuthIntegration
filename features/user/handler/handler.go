package handler

import (
	"emailnotifl3n/features/user"
	"emailnotifl3n/utils/upload"
	"emailnotifl3n/utils/middlewares"
	"emailnotifl3n/utils/responses"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService user.UserServiceInterface
	s3         upload.S3UploaderInterface
}

func New(service user.UserServiceInterface, s3Uploader upload.S3UploaderInterface) *UserHandler {
	return &UserHandler{
		userService: service,
		s3:         s3Uploader,
	}
}

func (handler *UserHandler) RegisterUser(c echo.Context) error {
	newUser := UserRequest{}
	errBind := c.Bind(&newUser)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data, data not valid", nil))
	}

	userCore := RequestToCore(newUser)
	errInsert := handler.userService.Create(userCore)
	if errInsert != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error insert data. "+errInsert.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success insert user", nil))
}

func (handler *UserHandler) Login(c echo.Context) error {
	var reqData = LoginRequest{}
	errBind := c.Bind(&reqData)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data, data not valid", nil))
	}
	result, token, err := handler.userService.Login(reqData.Email, reqData.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error login. "+err.Error(), nil))
	}
	responseData := map[string]any{
		"token": token,
		"nama":  result.Name,
	}
	return c.JSON(http.StatusOK, responses.WebResponse("success login", responseData))
}

func (handler *UserHandler) GetUser(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)

	result, errSelect := handler.userService.GetById(userIdLogin)
	if errSelect != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error read data. "+errSelect.Error(), nil))
	}

	var userResult = CoreToResponse(result)
	return c.JSON(http.StatusOK, responses.WebResponse("success read data", userResult))
}

func (handler *UserHandler) UpdateUser(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)

	var userData = UserRequest{}
	errBind := c.Bind(&userData)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data not valid", nil))
	}

	fileData, err := c.FormFile("photo_profile")
	if err != nil && err != http.ErrMissingFile {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error retrieving the file", nil))
	}

	var imageURL string
	if fileData != nil {
		imageURL, err = handler.s3.UploadImage(fileData)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.WebResponse("error uploading the image "+err.Error(), nil))
		}
	}

	userCore := UpdateRequestToCoreUpdate(userData, imageURL)
	errUpdate := handler.userService.Update(userIdLogin, userCore)
	if errUpdate != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error update data. "+errUpdate.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success update data", nil))
}

func (handler *UserHandler) DeleteUser(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)

	errDelete := handler.userService.Delete(userIdLogin)
	if errDelete != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error delete data. "+errDelete.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success delete data", nil))
}

func (handler *UserHandler) ChangePassword(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)

	var passwords = ChangePasswordRequest{}
	errBind := c.Bind(&passwords)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data not valid", nil))
	}

	errChange := handler.userService.ChangePassword(userIdLogin, passwords.OldPassword, passwords.NewPassword)
	if errChange != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error change password. "+errChange.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success change password", nil))
}
