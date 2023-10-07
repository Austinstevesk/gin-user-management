package routes

import (
	"gin-user-management/controllers"
	"gin-user-management/initializers"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouterController AuthRouteController

	UserController      controllers.UserController
	UserRouterController UserRouteController
)


func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	AuthController= controllers.NewAuthController(initializers.DB)
	AuthRouterController = NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(initializers.DB)
	UserRouterController = NewUserRouteController(UserController)

	server = gin.Default()
}

func NewRoutes() *gin.Engine {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	// set up v1
	r := server.Group("/api/v1")
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})

	// include the routes
	AuthRouterController.AuthRouter(r)
	UserRouterController.UserRouter(r)	
	return server
}

type routes struct {
	server *gin.Engine
}

func (r *routes) Run(addr ...string) error {
	return r.server.Run()
}