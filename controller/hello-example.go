package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Hello struct {
}

func (s *Hello) Hello(c *gin.Context) {
	var form struct {
		Msg string `form:"msg" binding:"required"`
		PagerForm
	}
	if err := c.ShouldBindQuery(&form); err != nil {
		fmt.Printf("controller %v\n", err)
		c.Error(ErrNew(err, ParamErr))
		return
	}

	resp, err := srv.Hello.Hello(form.Msg)

	if err != nil {
		fmt.Printf("controller %v\n", err)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, ResponseNew(c, resp))
}

func (s *Hello) HelloTime(c *gin.Context) {
	var form struct {
		Date time.Time `form:"date" binding:"required,timing" time_format:"2006-01-02"`
	}
	if err := c.ShouldBindQuery(&form); err != nil {
		fmt.Printf("controller %v\n", err)
		c.Error(ErrNew(err, ParamErr))
		return
	}
	resp := srv.Hello.HelloTime(form.Date)

	c.JSON(http.StatusOK, ResponseNew(c, resp))
}
