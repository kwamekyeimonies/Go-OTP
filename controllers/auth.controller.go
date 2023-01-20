package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/kwamekyeimonies/Go-OTP/model"
	"github.com/pquerna/otp/totp"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

func (ac *AuthController) SignUpUser(ctx *gin.Context) {

	var payload *model.RegisterUserInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	newUser := model.User{
		UserId:   uuid.New().String(),
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
	}

	result := ac.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "Duplicate key violates unique") {
		ctx.JSON(http.StatusConflict, gin.H{
			"status":  "fail",
			"message": "Email Already exist,use a different email",
		})
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": result.Error.Error(),
		})
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"user":   result,
	})
}

func (ac *AuthController) LoginUser(ctx *gin.Context) {
	var payload *model.LoginUserInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Invalid Email or Password",
		})
		return
	}

	var user model.User

	result := ac.DB.First(&user, "email = ? ", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Invalid Email or Password",
		})
		return
	}

	userResponse := gin.H{
		"id":    user.UserId,
		"name":  user.Name,
		"email": user.Email,
		"otp":   user.Otp_enabled,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user":   userResponse,
	})
}

func (ac *AuthController) GenerateOTP(ctx *gin.Context) {

	var payload *model.OTPInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})

		return
	}

	key, err := totp.Generate(
		totp.GenerateOpts{
			Issuer:      "github.com/Kwamekyeimonies",
			AccountName: "Kwamekyeimonies",
			SecretSize:  15,
		},
	)

	if err != nil {
		panic(err)
	}

	var user model.User

	result := ac.DB.First(&user, "id = ?", payload.UserId)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Invalid email or password",
		})

		return
	}

	dateToUpdate := model.User{
		Otp_secret:   key.Secret(),
		Otp_auth_url: key.URL(),
	}

	ac.DB.Model(&user).Updates(dateToUpdate)

	otpResponse := gin.H{
		"base32":      key.Secret(),
		"otpauth_url": key.URL(),
	}

	ctx.JSON(http.StatusOK, otpResponse)

}

// https://codevoweb.com/two-factor-authentication-2fa-in-golang/
