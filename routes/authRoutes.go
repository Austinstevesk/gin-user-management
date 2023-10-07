package routes

import (
	"gin-user-management/controllers"

	"github.com/gin-gonic/gin"
)


type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) AuthRouteController {
	return AuthRouteController{authController: authController}
}

func (rc *AuthRouteController) AuthRouter(rg *gin.RouterGroup) {
	router := rg.Group("/auth")

	router.POST("/register", rc.authController.SignUp)
	router.POST("/login", rc.authController.SignInUser)
	router.POST("/refresh", rc.authController.RefreshAccessToken)
	router.GET("/logout", rc.authController.LogoutUser)
}