package service

import (
	"encoding/gob"
	"errors"

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

func sessionSet(c *gin.Context, name string, body interface{}) error {
	session := sessions.Default(c)
	if body == nil {
		return errors.New("body is nil")
	}
	gob.Register(body)
	session.Set(name, body)
	if err := session.Save(); err != nil {
		return err
	}
	return nil
}
