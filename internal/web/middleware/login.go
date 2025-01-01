package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct {
}

func (middleware *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if path == "/users/signup" || path == "/users/signin" {
			return
		}

		session := sessions.Default(c)
		if session.Get("userId") == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
