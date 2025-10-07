package logger

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	source := getCaller(2)
	GinLogger.WithField("source", source).Errorf(format, args...)
}

func Warnf(format string, args ...any) {
	source := getCaller(2)
	GinLogger.WithField("source", source).Warnf(format, args...)
}

func Infof(format string, args ...any) {
	source := getCaller(2)
	GinLogger.WithField("source", source).Infof(format, args...)
}

func Debugf(format string, args ...any) {
	source := getCaller(2)

	// 只有当前日志级别不会输出 Debug 时，才添加到 backtrace
	if IsBacktraceEnabled() && GinLogger.Level > logrus.DebugLevel {
		msg := fmt.Sprintf(format, args...)
		fields := logrus.Fields{"source": source}
		addToBacktrace(logrus.DebugLevel, msg, fields, source, nil)
	}

	GinLogger.WithField("source", source).Debugf(format, args...)
}

// DebugTraced forced to add into backtrace
func DebugTraced(format string, args ...any) {
	source := getCaller(2)
	msg := fmt.Sprintf(format, args...)
	fields := logrus.Fields{"source": source}

	// always backtrace
	if IsBacktraceEnabled() {
		addToBacktrace(logrus.DebugLevel, msg, fields, source, nil)
	}

	GinLogger.WithField("source", source).Debugf(format, args...)
}

// func ErrorCtx(ctx *gin.Context, format string, args ...any) {
// 	GinLogger.WithContext(ctx).Errorf(format, args...)
// }
// func WarnCtx(ctx *gin.Context, format string, args ...any) {
// 	GinLogger.WithContext(ctx).Warnf(format, args...)
// }
// func InfoCtx(ctx *gin.Context, format string, args ...any) {
// 	GinLogger.WithContext(ctx).Infof(format, args...)
// }
// func DebugCtx(ctx *gin.Context, format string, args ...any) {
// 	GinLogger.WithContext(ctx).Debugf(format, args...)
// }
