package middleware

import (
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

// Secure creates the middleware that sets up the necessary HTTP security headers.
func Secure() gin.HandlerFunc {
	return secure.New(secure.Config{
		SSLRedirect: false,
	})
}
