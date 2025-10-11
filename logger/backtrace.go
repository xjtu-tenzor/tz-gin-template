package logger

import (
	"container/ring"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LogEntry struct {
	Level   logrus.Level
	Time    time.Time
	Message string
	Fields  logrus.Fields
	Source  string
	Context *gin.Context
}

type BacktraceManager struct {
	mu      sync.RWMutex
	buffer  *ring.Ring
	size    int
	enabled bool
}

var backtrace *BacktraceManager

func init() {
	backtrace = &BacktraceManager{
		buffer:  ring.New(100), // 默认缓存100条
		size:    100,
		enabled: false,
	}
}

// EnableBacktrace 启用回溯
func EnableBacktrace(size int) {
	backtrace.mu.Lock()
	defer backtrace.mu.Unlock()
	backtrace.enabled = true
	if size <= 0 {
		size = 1 // 最小缓冲区大小
	}
	backtrace.buffer = ring.New(size)
	backtrace.size = size
}

// DisableBacktrace 禁用回溯
func DisableBacktrace() {
	backtrace.mu.Lock()
	defer backtrace.mu.Unlock()
	backtrace.enabled = false
}

// IsBacktraceEnabled 检查回溯是否启用
func IsBacktraceEnabled() bool {
	backtrace.mu.RLock()
	defer backtrace.mu.RUnlock()
	return backtrace.enabled
}

// addToBacktrace 添加日志到回溯缓冲区
func addToBacktrace(level logrus.Level, msg string, fields logrus.Fields, source string, ctx *gin.Context) {
	if !backtrace.enabled {
		return
	}

	backtrace.mu.Lock()
	defer backtrace.mu.Unlock()

	if backtrace.buffer == nil {
		return
	}

	entry := LogEntry{
		Level:   level,
		Time:    time.Now(),
		Message: msg,
		Fields:  fields,
		Source:  source,
		Context: ctx,
	}

	backtrace.buffer.Value = entry
	backtrace.buffer = backtrace.buffer.Next()
}

// DumpBacktrace 输出回溯日志
func DumpBacktrace() {
	if !backtrace.enabled {
		GinLogger.Warn("Backtrace is not enabled")
		return
	}

	backtrace.mu.RLock()
	defer backtrace.mu.RUnlock()

	GinLogger.Info("***** BACKTRACE START *****")

	backtrace.buffer.Do(func(v any) {
		if v != nil {
			entry := v.(LogEntry)
			logEntry := GinLogger.WithFields(entry.Fields).WithField("source", entry.Source)

			if entry.Context != nil {
				logEntry = logEntry.WithContext(entry.Context)
			}

			switch entry.Level {
			case logrus.DebugLevel:
				logEntry.Debug(entry.Message)
			case logrus.InfoLevel:
				logEntry.Info(entry.Message)
			case logrus.WarnLevel:
				logEntry.Warn(entry.Message)
			case logrus.ErrorLevel:
				logEntry.Error(entry.Message)
			case logrus.FatalLevel:
				logEntry.Fatal(entry.Message)
			case logrus.PanicLevel:
				logEntry.Panic(entry.Message)
			}
		}
	})

	GinLogger.Info("***** BACKTRACE END *****")
}

// DumpBacktraceCtx 带 context 输出回溯日志
func DumpBacktraceCtx(ctx *gin.Context) {
	if !backtrace.enabled {
		GinLogger.WithContext(ctx).Warn("Backtrace is not enabled")
		return
	}

	GinLogger.WithContext(ctx).Info("***** BACKTRACE START *****")
	DumpBacktrace()
	GinLogger.WithContext(ctx).Info("***** BACKTRACE END *****")
}

// ErrorWithBacktrace 记录错误并自动输出回溯
func ErrorWithBacktrace(format string, args ...any) {
	Errorf(format, args...)
	DumpBacktrace()
}

// // ErrorCtxWithBacktrace 带 context 记录错误并自动输出回溯
// func ErrorCtxWithBacktrace(ctx *gin.Context, format string, args ...interface{}) {
// 	ErrorCtx(ctx, format, args...)
// 	DumpBacktraceCtx(ctx)
// }
