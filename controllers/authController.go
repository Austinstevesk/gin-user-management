package controllers

import (
	"fmt"
	"gin-user-management/helpers"
	"gin-user-management/initializers"
	"gin-user-management/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{
		DB: DB,
	}
}

func (ac *AuthController) SignUp(ctx *gin.Context)  {
	var payload *models.SignUpInput

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"message": err.Error(),
		})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"message": "passwords do not match",
		})
		return
	}

	hashedPassword, err := helpers.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status": "error",
			"message": err.Error(),
		})
	}

	newUser := models.User{
		Name: payload.Name,
		Email: payload.Email,
		Password: hashedPassword,
		Role: "user",
		Verified: false,
		Photo: payload.Photo,
		Provider: "local",
	}

	result := ac.DB.Create(&newUser)
	if result.Error != nil && strings.Contains(
		result.Error.Error(), "duplicate key value violates unique") {
			ctx.JSON(http.StatusConflict, gin.H{
				"status": "failed",
				"message": "User wit that email address already exists",
			})
			return
		} else if result.Error != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{
				"status": "failed",
				"message": "Could not create user",
			})
			return
		}

		initializers.LoadEnvVariables()
		// generate code
		code := randstr.String(20)
		verificationCode := helpers.Encode(code)

		// update user in db
		newUser.VerificationCode = verificationCode
		ac.DB.Save(newUser)

		firstName := newUser.Name
		if strings.Contains(firstName, " ") {
			firstName = strings.Split(firstName, " ")[1]
		}

		// send email
		emailData := helpers.EmailData{
			URL: os.Getenv("CLIENT_ORIGIN") + "/verifyemail/" + code,
			FirstName: firstName,
			Subject: "Your account verification code",
		}

		helpers.SendEmail(&newUser, &emailData, "verificationCode.html")
		message := "We sent an email with verification code to  " + newUser.Email

		ctx.JSON(http.StatusCreated, gin.H{
			"status": "success",
			"message": message,
		})
}

func (ac *AuthController) VerifyEmail (ctx *gin.Context) {
	code := ctx.Params.ByName("verificationCode")
	verificationCode := helpers.Encode(code)


	var updatedUser models.User
	result := ac.DB.First(&updatedUser, "verification_code = ?", verificationCode)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"message": "Invalid verification code",
		})
		return
	}

	if updatedUser.Verified {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"message": "User already verified",
		})
		return
	}

	updatedUser.VerificationCode = ""
	updatedUser.Verified = true
	ac.DB.Save(&updatedUser)

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Email verified successfully",
	})
}

func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var payload *models.SignInInput

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"message": "Invalid email or password",
		})
		return
	}
	//result := ac.DB.First(&models.User{}, "email")
	var user models.User
	result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"message": "Invalid email or password",
		})
		return
	}

	if err := helpers.VerifyPassword(user.Password, payload.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"message": "Invalid email or password",
		})
		return
	}
	
	if !user.Verified {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status": "failed",
			"message": "Please verify your account",
		})
		return
	}
	accessToken, err := helpers.CreateToken(time.Hour, user.ID, os.Getenv("ACCESS_TOKEN_PRIVATE_KEY"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"message": err.Error(),
		})
		return
	}
	refreshToken, err := helpers.CreateToken(time.Hour * 24 * 30, user.ID, os.Getenv("REFRESH_TOKEN_PRIVATE_KEY"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"message": err.Error(),
		})
		return
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("accessToken", accessToken, 3600, "/", "", false, true)
	ctx.SetCookie("refreshToken", refreshToken, 3600 * 24 * 30, "/", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"accessToken": accessToken,
		"refreshToken": refreshToken,
	})
}


func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"

	cookie, err := ctx.Cookie("refreshToken")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"message": message,
		})
		return
	}

	sub, err := helpers.ValidateToken(cookie, os.Getenv("REFRESH_TOKEN_PUBLIC_KEY"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"message": err.Error(),
		})
		return
	}

	var user models.User
	result := ac.DB.First(&user, fmt.Sprint(sub))
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"message": "the user belonging to this token no longer exists",
		})
		return
	}

	accessToken, err := helpers.CreateToken(time.Hour, user.ID, os.Getenv("ACCESS_TOKEN_PRIVATE_KEY"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"message": err.Error(),
		})
		return
	}
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("accessToken", accessToken, 3600, "", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"accessToken": accessToken,
	})
}


func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("accessToken", "", -1, "", "", false, true)
	ctx.SetCookie("refreshToken", "", -1, "", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Logged out successfully",
	})
}