package router

import (
	"emailnotifl3n/app/cache"
	ud "emailnotifl3n/features/user/data"
	uh "emailnotifl3n/features/user/handler"
	us "emailnotifl3n/features/user/service"
	"emailnotifl3n/utils/email"
	"emailnotifl3n/utils/encrypts"
	"emailnotifl3n/utils/middlewares"
	"emailnotifl3n/utils/upload"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB, e *echo.Echo, rds cache.Redis) {
	hash := encrypts.New()
	s3Uploader := upload.New()
	email := email.New()

	userData := ud.New(db, rds)
	userService := us.New(userData, hash)
	userHandlerAPI := uh.New(userService, s3Uploader, email)

	// define routes/ endpoint USER
	e.POST("/login", userHandlerAPI.Login)
	e.POST("/users", userHandlerAPI.RegisterUser)
	e.GET("/users", userHandlerAPI.GetUser, middlewares.JWTMiddleware())
	e.PUT("/users", userHandlerAPI.UpdateUser, middlewares.JWTMiddleware())
	e.DELETE("/users", userHandlerAPI.DeleteUser, middlewares.JWTMiddleware())
	e.PUT("/change-password", userHandlerAPI.ChangePassword, middlewares.JWTMiddleware())
	e.POST("forgot-password", userHandlerAPI.ForgotPassword)
	e.PATCH("reset-password", userHandlerAPI.ResetPassword)
	e.POST("verification", userHandlerAPI.SendVerifyEmail)
	e.PATCH("verification", userHandlerAPI.VerifyEmailLink)
	e.POST("request-code-password", userHandlerAPI.RequestCodePassword)
	e.PATCH("reset-password-code", userHandlerAPI.ResetPasswordCode)
	e.POST("request-code-verify", userHandlerAPI.RequestCodeVerify)
	e.PATCH("verification-email", userHandlerAPI.VerifyEmailCode)
}