package middleware

import (
	"errors"
	"template/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SessionUse(c *gin.Context) {
	if _, ok := c.Get("init"); !ok {
		c.Set("init", 1)
		c.Next()
	}
	if _, ok := c.Get("session_used"); ok {
		session := sessions.Default(c)
		if err := session.Save(); err != nil {
			c.Error(service.ErrNew(errors.New("fail to save session"), service.OpErr))
			c.Abort()
			return
		}
		return
	}
}
