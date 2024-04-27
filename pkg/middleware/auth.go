package middleware

import (
	"github.com/gin-gonic/gin"
	"go-finance-tracker/pkg/logger"
	"go-finance-tracker/pkg/utils"
	"net/http"
)

func RequireAuthMiddleware(c *gin.Context) {
	log := logger.GetLogger()

	var token string
	cookie, err := c.Cookie("jwt")
	if err != nil {
		log.Error("JWT token not found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT token not found"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	token = cookie

	if token == "" {
		log.Error("Invalid token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	id, username, err := utils.VerifyToken(token)
	if err != nil {
		log.Error("Token verification failed:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token verification failed"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("id", id)
	c.Set("username", username)

	log.Info("User is authenticated!")
	c.Next()
}
