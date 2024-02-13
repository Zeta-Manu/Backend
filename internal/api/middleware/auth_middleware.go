package middleware

import (
	"github.com/gin-gonic/gin"
)

func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func isAuthenticated(c *gin.Context) bool {
	return true
}
