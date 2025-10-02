package logger

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

var GinLogger *logrus.Logger

type ResponseBodyWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func getCaller(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown:0"
	}
	file = file[strings.LastIndex(file, "/")+1:] // Get only the file name
	return fmt.Sprintf("%s:%d", file, line)
}

func (w ResponseBodyWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Errorf(format string, args ...any) {
	GinLogger.WithField("source", getCaller(2)).Errorf(format, args...)
}
func Warnf(format string, args ...any) {
	GinLogger.WithField("source", getCaller(2)).Warnf(format, args...)
}
func Infof(format string, args ...any) {
	GinLogger.WithField("source", getCaller(2)).Infof(format, args...)
}
func Debugf(format string, args ...any) {
	GinLogger.WithField("source", getCaller(2)).Debugf(format, args...)
}

func ErrorCtx(ctx *gin.Context, format string, args ...any) {
	GinLogger.WithContext(ctx).Errorf(format, args...)
}
func WarnCtx(ctx *gin.Context, format string, args ...any) {
	GinLogger.WithContext(ctx).Warnf(format, args...)
}
func InfoCtx(ctx *gin.Context, format string, args ...any) {
	GinLogger.WithContext(ctx).Infof(format, args...)
}
func DebugCtx(ctx *gin.Context, format string, args ...any) {
	GinLogger.WithContext(ctx).Debugf(format, args...)
}
