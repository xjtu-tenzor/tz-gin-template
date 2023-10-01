package middleware

import (
	"errors"

	"template/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CheckRole(min int) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userSession := session.Get("user-session")
		if userSession == nil {
			c.Error(service.ErrNew(errors.New("您未登录"), service.AuthErr))
			c.Abort()
			return
		}
		if userSession.(service.UserSession).Level < min {
			c.Error(service.ErrNew(errors.New("权限不足"), service.LevelErr))
			c.Abort()
			return
		}
		c.Next()
	}
}
