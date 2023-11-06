package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"template/controller"

	vl "template/service/validator"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Error(c *gin.Context) {
	c.Next()
	if len(c.Errors) != 0 {
		err := c.Errors.Last().Err
		switch err := err.(type) {
		case validator.ValidationErrors:
			errorHandle(c, err.Translate(vl.Trans))
		case *strconv.NumError, *json.UnmarshalTypeError:
			errorHandle(c, errors.New("错误的传入参数"))
		default:
			errorHandle(c, err)
		}
	}
}

func errorHandle(c *gin.Context, err any) {
	errMsg := fmt.Sprintf("%v: %v\n", controller.ErrorMapper[uint64(c.Errors.Last().Type)], err)
	c.JSON(http.StatusOK, controller.Response{
		Success: false,
		Message: errMsg,
		Code:    uint64(c.Errors.Last().Type),
	})
}
