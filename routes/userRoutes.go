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
	router.Use(middleware.RequireAuth())
	router.GET("/profile", uc.userController.Userprofile)
	router.GET("/admin", middleware.RequireAdmin(), uc.userController.Userprofile)
}