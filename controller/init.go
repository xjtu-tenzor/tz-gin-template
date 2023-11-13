package controller

import (
	"encoding/gob"
	"template/common"
	"template/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Code    uint64 `json:"code,omitempty"`
}

func ResponseNew(c *gin.Context, obj any) *Response {
	session := sessions.Default(c)
	if session.Save() != nil {
		return &Response{
			Success: false,
			Message: "fail to save session",
			Code:    uint64(common.SysErr),
		}
	}
	return &Response{
		Success: true,
		Data:    obj,
	}
}

var srv = service.New()

func init() {
	gob.Register(UserSession{})
}
