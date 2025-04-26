package logger

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var GinLogger *logrus.Logger

type ResponseBodyWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w ResponseBodyWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Errorf(format string, args ...interface{}) {
	GinLogger.Errorf(format, args...)
}
func Warnf(format string, args ...interface{}) {
	GinLogger.Warnf(format, args...)
}
func Infof(format string, args ...interface{}) {
	GinLogger.Infof(format, args...)
}
func Debugf(format string, args ...interface{}) {
	GinLogger.Debugf(format, args...)
}

func ErrorCtx(ctx *gin.Context, format string, args ...interface{}) {
	GinLogger.WithContext(ctx).Errorf(format, args...)
}
func WarnCtx(ctx *gin.Context, format string, args ...interface{}) {
	GinLogger.WithContext(ctx).Warnf(format, args...)
}
func InfoCtx(ctx *gin.Context, format string, args ...interface{}) {
	GinLogger.WithContext(ctx).Infof(format, args...)
}
func DebugCtx(ctx *gin.Context, format string, args ...interface{}) {
	GinLogger.WithContext(ctx).Debugf(format, args...)
}
