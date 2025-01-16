package middleware

import (
	"github.com/gin-gonic/gin"
)

// Attach sets up the necessary middleware for the given router.
func Attach(router *gin.Engine) *gin.Engine {
	router.Use(Cors())
	router.Use(Secure())

	return router
}
