package middleware

import (
	"fmt"
	"gin-user-management/helpers"
	"gin-user-management/initializers"
	"gin-user-management/models"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func (ctx *gin.Context)  {
		var accessToken string
		cookie, err := ctx.Cookie("accessToken")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			accessToken = fields[1]
		} else if err == nil {
			accessToken = cookie
		}

		if accessToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "failed",
				"message": "You are logged out",
			})
			return
		}

		sub, err := helpers.ValidateToken(accessToken, os.Getenv("ACCESS_TOKEN_PUBLIC_KEY"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "failed",
				"message": err.Error(),
			})
			return
		}

		var user models.User
		result := initializers.DB.First(&user, fmt.Sprint(sub))
		if result.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "failed",
				"message": "the user belonging to this token no longer exists",
			})
		}

		ctx.Set("currentUser", user)
		ctx.Next()

	}
}

func RequireAdmin() gin.HandlerFunc {
	return func (ctx *gin.Context)  {
		user := ctx.MustGet("currentUser").(models.User)
		if user.Role != "admin"{
			ctx.AbortWithStatusJSON(http.StatusFailedDependency, gin.H{
				"status": "Failed",
				"message": "Not enough permissions",
			})
			return
		}
		ctx.Next()
		
	}
}