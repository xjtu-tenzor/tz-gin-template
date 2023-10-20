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
		return ErrNew(errors.New("body is nil"), OpErr)
	}
	gob.Register(body)
	session.Set(name, body)
	return session.Save()
}

func sessionUpdate(c *gin.Context, name string, body interface{}) error {
	if sessionGet(c, name) == nil {
		return ErrNew(errors.New("cannot find session, Set please"), OpErr)
	}
	return sessionSet(c, name, body)
}

func sessionClear(c *gin.Context) error {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		return ErrNew(err, OpErr)
	}
	return nil
}

func sessionDelete(c *gin.Context, name string) error {
	session := sessions.Default(c)
	if sessionGet(c, name) == nil {
		return errors.New("cannot find session, Set please")
	}
	session.Delete(name)
	if err := session.Save(); err != nil {
		return ErrNew(err, OpErr)
	}
	return nil
}
