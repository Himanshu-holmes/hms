package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/himanshu-holmes/hms/internal/authorization"
	"github.com/himanshu-holmes/hms/internal/model"
)
var (
	secret = []byte(os.Getenv("SECRET"))
)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is missing",
			})
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			return
		}
		tokenString := parts[1]
		auth := authorization.AuthWithOutDuration(secret)
		info,err := auth.Authorize(c.Request.Context(),tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.Set("info",info)

	}
}

func RoleMiddleware(requiredRole model.UserRole)gin.HandlerFunc{
	return func(c *gin.Context){
		info := c.MustGet("info").(authorization.Info)
		if info.Type != authorization.AccessToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token type",
			})
			return
		}
		if info.Role != requiredRole {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}
	    c.Next()
	}
}

