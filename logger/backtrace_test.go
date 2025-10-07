package logger

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func init() {
	// 初始化 GinLogger，以防它没有被初始化
	if GinLogger == nil {
		GinLogger = logrus.New()
		GinLogger.SetLevel(logrus.DebugLevel)
	}
}

func TestBacktraceEnableDisable(t *testing.T) {
	// 保存原始状态
	originalEnabled := IsBacktraceEnabled()
	defer func() {
		if originalEnabled {
			EnableBacktrace(100)
		} else {
			DisableBacktrace()
		}
	}()

	// 测试禁用状态
	DisableBacktrace()
	if IsBacktraceEnabled() {
		t.Error("Expected backtrace to be disabled")
	}

	// 测试启用状态
	EnableBacktrace(10)
	if !IsBacktraceEnabled() {
		t.Error("Expected backtrace to be enabled")
	}

	// 再次禁用
	DisableBacktrace()
	if IsBacktraceEnabled() {
		t.Error("Expected backtrace to be disabled after second disable")
	}
}

func TestBacktraceBufferSizes(t *testing.T) {
	defer DisableBacktrace()

	// 测试不同缓冲区大小
	sizes := []int{5, 10, 50, 100}
	for _, size := range sizes {
		EnableBacktrace(size)
		if !IsBacktraceEnabled() {
			t.Errorf("Expected backtrace to be enabled with size %d", size)
		}
	}
}

func TestBacktraceLogging(t *testing.T) {
	// 保存原始日志级别
	originalLevel := GinLogger.Level
	defer func() {
		GinLogger.Level = originalLevel
		DisableBacktrace()
	}()

	// 设置日志级别为 Info，这样 Debug 日志不会输出但会被 trace
	GinLogger.Level = logrus.InfoLevel
	EnableBacktrace(5)

	// 记录一些 debug 日志
	Debugf("Debug message 1")
	Debugf("Debug message 2")
	DebugTraced("Traced debug message")

	// 验证功能正常运行
	if !IsBacktraceEnabled() {
		t.Error("Expected backtrace to remain enabled")
	}
}

func TestBacktraceDump(t *testing.T) {
	defer DisableBacktrace()

	// 测试在禁用状态下 dump
	DisableBacktrace()
	// 这应该不会 panic，只会输出警告
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("DumpBacktrace panicked when disabled: %v", r)
			}
		}()
		DumpBacktrace()
	}()

	// 测试在启用状态下 dump
	EnableBacktrace(10)
	DebugTraced("Test message for dump")
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("DumpBacktrace panicked when enabled: %v", r)
			}
		}()
		DumpBacktrace()
	}()
}

func TestBacktraceWithContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	defer DisableBacktrace()

	EnableBacktrace(5)
	c, _ := gin.CreateTestContext(nil)

	// 测试带 context 的 dump
	DebugTraced("Context test message")
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("DumpBacktraceCtx panicked: %v", r)
			}
		}()
		DumpBacktraceCtx(c)
	}()

	// 测试 nil context
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("DumpBacktraceCtx with nil context panicked: %v", r)
			}
		}()
		var nilCtx *gin.Context
		DumpBacktraceCtx(nilCtx)
	}()
}

func TestBacktraceErrorWithBacktrace(t *testing.T) {
	defer DisableBacktrace()

	EnableBacktrace(5)
	DebugTraced("Test error trace")

	// 测试错误日志和自动 dump
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("ErrorWithBacktrace panicked: %v", r)
			}
		}()
		ErrorWithBacktrace("Test error: %s", "sample error")
	}()
}

func TestBacktraceConcurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrency test in short mode")
	}

	defer DisableBacktrace()
	EnableBacktrace(20)

	// 并发测试
	done := make(chan bool, 5)

	// 并发写入
	for i := 0; i < 5; i++ {
		go func(id int) {
			defer func() { done <- true }()
			for j := 0; j < 10; j++ {
				DebugTraced("Goroutine %d message %d", id, j)
				time.Sleep(time.Millisecond)
			}
		}(i)
	}

	// 等待完成
	for i := 0; i < 5; i++ {
		<-done
	}

	// 最终 dump 应该不会 panic
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Final DumpBacktrace panicked: %v", r)
			}
		}()
		DumpBacktrace()
	}()
}

func TestBacktraceCircularBuffer(t *testing.T) {
	defer DisableBacktrace()

	// 测试小缓冲区的循环覆盖
	EnableBacktrace(3)

	// 添加超过缓冲区大小的消息
	for i := 0; i < 10; i++ {
		DebugTraced("Message %d", i)
	}

	// Dump 应该正常工作
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("DumpBacktrace panicked with circular buffer: %v", r)
			}
		}()
		DumpBacktrace()
	}()
}

func TestBacktraceEdgeCases(t *testing.T) {
	defer DisableBacktrace()

	// 测试零大小缓冲区
	EnableBacktrace(0)
	DebugTraced("Zero size buffer test")
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Zero size buffer panicked: %v", r)
			}
		}()
		DumpBacktrace()
	}()

	// 测试空消息
	EnableBacktrace(5)
	DebugTraced("")
	Debugf("")
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Empty message panicked: %v", r)
			}
		}()
		DumpBacktrace()
	}()
}

// BenchmarkBacktraceDisabled 性能基准测试 - 禁用状态
func BenchmarkBacktraceDisabled(b *testing.B) {
	DisableBacktrace()
	GinLogger.Level = logrus.InfoLevel

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Debugf("Benchmark message %d", i)
	}
}

// BenchmarkBacktraceEnabled 性能基准测试 - 启用状态
func BenchmarkBacktraceEnabled(b *testing.B) {
	EnableBacktrace(100)
	GinLogger.Level = logrus.InfoLevel

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Debugf("Benchmark message %d", i)
	}
	b.StopTimer()
	DisableBacktrace()
}

// BenchmarkBacktraceTraced 性能基准测试 - 强制跟踪
func BenchmarkBacktraceTraced(b *testing.B) {
	EnableBacktrace(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DebugTraced("Benchmark traced message %d", i)
	}
	b.StopTimer()
	DisableBacktrace()
}
