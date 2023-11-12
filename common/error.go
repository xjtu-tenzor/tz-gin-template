package common

import "github.com/gin-gonic/gin"

const (
	ParamErr gin.ErrorType = iota + 3
	SysErr
	OpErr
	AuthErr
	LevelErr
)

var ErrorMapper = map[uint64]string{
	1: "内部错误",
	2: "公开错误",
	3: "参数错误",
	4: "系统错误",
	5: "操作错误",
	6: "鉴权错误",
	7: "权限错误",
}

func ErrNew(err error, errType gin.ErrorType) error {
	err = &gin.Error{
		Err:  err,
		Type: errType,
	}
	return err
}
