package handler

import (
	"emailnotifl3n/features/user"
	"emailnotifl3n/utils/email"
	"emailnotifl3n/utils/middlewares"
	oauthfacebook "emailnotifl3n/utils/oauthFacebook"
	"emailnotifl3n/utils/oauthGoogle"
	"emailnotifl3n/utils/responses"
	"emailnotifl3n/utils/upload"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService user.UserServiceInterface
	s3          upload.S3UploaderInterface
	email       email.EmailInterface
	oauthGoogle oauthGoogle.GoogleInterface
	oauthFB     oauthfacebook.FacebookInterface
}

func New(service user.UserServiceInterface, s3Uploader upload.S3UploaderInterface, email email.EmailInterface, google oauthGoogle.GoogleInterface, fb oauthfacebook.FacebookInterface) *UserHandler {
	return &UserHandler{
		userService: service,
		s3:          s3Uploader,
		email:       email,
		oauthGoogle: google,
		oauthFB:     fb,
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

func (handler *UserHandler) ForgotPassword(c echo.Context) error {
	var ForgotReq = ForgotPasswordRequest{}
	errBind := c.Bind(&ForgotReq)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data not valid", nil))
	}

	user, token, err := handler.userService.ForgotPassword(ForgotReq.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(err.Error(), nil))
	}

	errForgot := handler.email.SendResetPasswordLink(user, token)
	if errForgot != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error sending reset password email - "+errForgot.Error(), nil))
	}
	return c.JSON(http.StatusOK, responses.WebResponse("reset password email sent", nil))
}

func (handler *UserHandler) ResetPassword(c echo.Context) error {
	token := c.QueryParam("token")

	userId, err := middlewares.ExtractUserIdFromResetPasswordToken(token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error extracting user id from reset password token. "+err.Error(), nil))
	}

	var resetPasswordRequest = ResetPasswordRequest{}
	errBind := c.Bind(&resetPasswordRequest)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data not valid", nil))
	}

	if resetPasswordRequest.ConfirmPassword != resetPasswordRequest.NewPassword {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("new password and confirm password do not match", nil))
	}

	errReset := handler.userService.ResetPassword(userId, resetPasswordRequest.NewPassword)
	if errReset != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error reset password. "+errReset.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success reset password", nil))
}

func (handler *UserHandler) SendVerifyEmail(c echo.Context) error {
	var ForgotReq = ForgotPasswordRequest{}
	errBind := c.Bind(&ForgotReq)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data not valid", nil))
	}

	user, token, err := handler.userService.ForgotPassword(ForgotReq.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(err.Error(), nil))
	}

	errForgot := handler.email.SendVerificationLink(user, token)
	if errForgot != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error sending reset password email - "+errForgot.Error(), nil))
	}
	return c.JSON(http.StatusOK, responses.WebResponse("verification email sent", nil))
}

func (handler *UserHandler) VerifyEmailLink(c echo.Context) error {
	token := c.QueryParam("token")

	userId, err := middlewares.ExtractUserIdFromResetPasswordToken(token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error extracting user id from reset password token. "+err.Error(), nil))
	}

	errReset := handler.userService.VerifyEmailLink(userId)
	if errReset != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error reset password. "+errReset.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success verification email", nil))
}

func (handler *UserHandler) RequestCodePassword(c echo.Context) error {
	var reqData = CodeRequest{}
	errBind := c.Bind(&reqData)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data, data not valid", nil))
	}

	userCore := CoderequestToCore(reqData)
	user, err := handler.userService.RequestCode(userCore.Email, userCore.Code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(err.Error(), nil))
	}

	errForgot := handler.email.SendCodeResetPassword(user, userCore.Code)
	if errForgot != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error sending code - "+errForgot.Error(), nil))
	}
	return c.JSON(http.StatusOK, responses.WebResponse("code email sent", nil))
}

func (handler *UserHandler) RequestCodeVerify(c echo.Context) error {
	var reqData = CodeRequest{}
	errBind := c.Bind(&reqData)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data, data not valid", nil))
	}

	userCore := CoderequestToCore(reqData)
	user, err := handler.userService.RequestCode(userCore.Email, userCore.Code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(err.Error(), nil))
	}

	errForgot := handler.email.SendCodeResetEmail(user, userCore.Code)
	if errForgot != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error sending code - "+errForgot.Error(), nil))
	}
	return c.JSON(http.StatusOK, responses.WebResponse("code email sent", nil))
}

func (handler *UserHandler) ResetPasswordCode(c echo.Context) error {
	code := c.QueryParam("code")

	var resetPasswordRequest = ResetPasswordRequestCode{}
	errBind := c.Bind(&resetPasswordRequest)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data not valid", nil))
	}

	if resetPasswordRequest.ConfirmPassword != resetPasswordRequest.NewPassword {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("new password and confirm password do not match", nil))
	}

	errReset := handler.userService.ResetPasswordCode(resetPasswordRequest.Email, resetPasswordRequest.NewPassword, code)
	if errReset != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error reset password. "+errReset.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success reset password", nil))
}

func (handler *UserHandler) VerifyEmailCode(c echo.Context) error {
	code := c.QueryParam("code")

	var verifyReq = ForgotPasswordRequest{}
	errBind := c.Bind(&verifyReq)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data not valid", nil))
	}

	errVerify := handler.userService.VerifyEmailCode(verifyReq.Email, code)
	if errVerify != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(errVerify.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success verification email", nil))
}

func (handler *UserHandler) GoogleLoginRedirect(c echo.Context) error {
	url := handler.oauthGoogle.GetAuthURL()
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (handler *UserHandler) RegisterWithGoogle(c echo.Context) error {
	code := c.QueryParam("code")

	googleOauthToken, err := handler.oauthGoogle.GetGoogleOauthToken(code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error getting Google OAuth token: "+err.Error(), nil))
	}

	googleUser, err := handler.oauthGoogle.GetGoogleUser(googleOauthToken.Access_token, googleOauthToken.Id_token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error getting Google user: "+err.Error(), nil))
	}

	errInsert := handler.userService.RegisterGoogle(*googleUser)
	if errInsert != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error insert data. "+errInsert.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success register user", googleUser))
}

func (handler *UserHandler) FacebookRedirect(c echo.Context) error {
	url := handler.oauthFB.GetAuthURL()
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (handler *UserHandler) RegisterWithFacebook(c echo.Context) error {
	code := c.QueryParam("code")

	fbOauthToken, err := handler.oauthFB.GetFacebookOauthToken(code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error getting fb OAuth token: "+err.Error(), nil))
	}

	fbUser, err := handler.oauthFB.GetFacebookUser(fbOauthToken.Access_token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error getting fb user: "+err.Error(), nil))
	}

	errInsert := handler.userService.RegisterGoogle(*fbUser)
	if errInsert != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error insert data. "+errInsert.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success register user", fbUser))
}