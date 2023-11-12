package controller

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

func SessionGet(c *gin.Context, name string) any {
	session := sessions.Default(c)
	return session.Get(name)
}

func SessionSet(c *gin.Context, name string, body any) {
	session := sessions.Default(c)
	if body == nil {
		return
	}
	gob.Register(body)
	session.Set(name, body)

}

func SessionUpdate(c *gin.Context, name string, body any) {
	SessionSet(c, name, body)
}

func SessionClear(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
}

func SessionDelete(c *gin.Context, name string) {
	session := sessions.Default(c)
	session.Delete(name)
}
