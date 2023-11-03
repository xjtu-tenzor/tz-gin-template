package service

import (
	"encoding/gob"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserSession struct {
	ID       int
	Username string
	Level    int
}

func sessionGet(c *gin.Context, name string) interface{} {
	session := sessions.Default(c)
	return session.Get(name)
}

func sessionSet(c *gin.Context, name string, body interface{}) {
	c.Set("session_used", 1)
	session := sessions.Default(c)
	if body == nil {
		return
	}
	gob.Register(body)
	session.Set(name, body)

}

func sessionUpdate(c *gin.Context, name string, body interface{}) {
	sessionSet(c, name, body)
}

func sessionClear(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
}

func sessionDelete(c *gin.Context, name string) {
	session := sessions.Default(c)
	session.Delete(name)
}
