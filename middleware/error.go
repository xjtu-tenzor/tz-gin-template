package middleware

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"template/common"
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
			var errs string
			for _, v := range err.Translate(vl.Trans) {
				errs = fmt.Sprintf("%v,%v", errs, v)
			}
			errorHandle(c, strings.Replace(errs, ",", "", 1))
		case *strconv.NumError, *json.UnmarshalTypeError, *time.ParseError, *xml.SyntaxError:
			errorHandle(c, errors.New("错误或非法的传入参数"))
		default:
			errorHandle(c, err)
		}
	}
}

func errorHandle(c *gin.Context, err any) {
	errMsg := fmt.Sprintf("%v: %v\n", common.ErrorMapper[uint64(c.Errors.Last().Type)], err)
	c.JSON(http.StatusOK, controller.Response{
		Success: false,
		Message: errMsg,
		Code:    uint64(c.Errors.Last().Type),
	})
}
