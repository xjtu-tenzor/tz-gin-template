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

func Errorf(ctx *gin.Context, format string, args ...interface{}) {
	GinLogger.WithContext(ctx).Errorf(format, args...)
}
func Warnf(ctx *gin.Context, format string, args ...interface{}) {
	GinLogger.WithContext(ctx).Warnf(format, args...)
}
func Infof(ctx *gin.Context, format string, args ...interface{}) {
	GinLogger.WithContext(ctx).Infof(format, args...)
}
func Debugf(ctx *gin.Context, format string, args ...interface{}) {
	GinLogger.WithContext(ctx).Debugf(format, args...)
}
