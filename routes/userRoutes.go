package routes

import (
	"gin-user-management/controllers"
	"gin-user-management/middleware"

	"github.com/gin-gonic/gin"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewUserRouteController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController: userController}
}

func (uc *UserRouteController) UserRouter(rg *gin.RouterGroup) {
	router := rg.Group("users")
	router.GET("/profile", middleware.RequireAuth(), uc.userController.Userprofile)
}