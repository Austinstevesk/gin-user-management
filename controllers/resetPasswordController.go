package controllers

import (
	"gin-user-management/helpers"
	"gin-user-management/initializers"
	"gin-user-management/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
)

func (ac *AuthController) ForgotPassword(ctx *gin.Context)  {
	var payload *models.ForgotPasswordInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"message": err.Error(),
			})
			return
	}

	message := "An email has been sent to email " + strings.ToLower(payload.Email)

	var user models.User
	result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"message": "User does not exist",
		})
		return
	}

	if !user.Verified {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"message": "Account is yet to be verified",
		})
		return
	}

	resetToken := randstr.String(20)
	passwordResetToken := helpers.Encode(resetToken)
	
	user.PasswordResetToken = passwordResetToken
	user.PasswordResetAt = time.Now().Add(time.Minute * 15)
	ac.DB.Save(&user)

	var firstName = user.Name
	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	initializers.LoadEnvVariables()

	emailData := helpers.EmailData {
		URL: os.Getenv("CLIENT_ORIGIN") + "/resetpassword/" + resetToken,
		FirstName: user.Name,
		Subject: "Your password reset token (Valid for 10 min)",
	}

	helpers.SendEmail(&user, &emailData, "resetPassword.html")

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": message,
	})

}

func (ac *AuthController) ResetPassword(ctx *gin.Context)  {
	var payload *models.ResetPasswordInput
	resetToken := ctx.Params.ByName("resetToken")

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"message": err.Error(),
			})
			return
	}

	if payload.Password != payload.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"message": "passwords provided to not match",
		})
		return
	}

	passwordResetToken := helpers.Encode(resetToken)

	var updatedUser models.User
	result := ac.DB.First(&updatedUser, "password_reset_token = ? AND password_reset_at > ?", passwordResetToken, time.Now())
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"message": "The reset token is invalid or has expired",
		})
		return
	}
	hashedPassword, _:= helpers.HashPassword(payload.Password)
	updatedUser.Password = hashedPassword
	ac.DB.Save(&updatedUser)

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("accessToken", "", -1, "", "", false, true)
	ctx.SetCookie("refreshToken", "", -1, "", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "password data updated successfully",
	})
}