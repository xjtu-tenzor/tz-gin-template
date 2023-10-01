package middleware

import (
	"fmt"
	"net/http"

	"template/controller"

	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context) {
	c.Next()
	if len(c.Errors) != 0 {
		errMsg := fmt.Sprintf("%v: %v\n", controller.ErrorMapper[uint64(c.Errors.Last().Type)], c.Errors.Last().Error())
		c.JSON(http.StatusOK, controller.Response{
			Success: false,
			Message: errMsg,
			Code:    uint64(c.Errors.Last().Type),
		})
	}
}
