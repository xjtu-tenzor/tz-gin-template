package middleware

import (
	"errors"

	"template/common"
	"template/controller"

	"github.com/gin-gonic/gin"
)

func CheckRole(min int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userSession := controller.SessionGet(c, "user-session")
		if userSession == nil {
			c.Error(common.ErrNew(errors.New("您未登录"), common.AuthErr))
			c.Abort()
			return
		}
		if userSession.(controller.UserSession).Level < min {
			c.Error(common.ErrNew(errors.New("权限不足"), common.LevelErr))
			c.Abort()
			return
		}
		c.Next()
	}
}
