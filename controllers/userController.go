package controllers

import (
	"gin-user-management/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(DB *gorm.DB) UserController {
	return UserController{DB: DB}
}

func (uc *UserController) Userprofile(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	userResponse := &models.UserResponse{
		ID: int(currentUser.ID),
		Name: currentUser.Name,
		Email: currentUser.Email,
		Photo: currentUser.Photo,
		Role: currentUser.Role,
		Provider: currentUser.Provider,
		CreatedAt: currentUser.CreatedAt,
		UpdatedAt: currentUser.UpdatedAt,
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user": userResponse,
		},
	})
}