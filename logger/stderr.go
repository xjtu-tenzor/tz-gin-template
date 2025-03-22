package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"runtime/debug"
)

type StdWriter struct {
	*logrus.Logger
}

func (sw StdWriter) Write(p []byte) (n int, err error) {
	if GinLogger.Level == logrus.DebugLevel {
		// Capture the full stack trace
		stack := debug.Stack()
		sw.Logger.Errorf("Find stderr: %s \nStack Trace: %s", string(p), stack)
		return len(p), nil
	} else {
		sw.Logger.Errorf("Find stderr: %s", string(p))
	}
	return len(p), nil
}

// capture stderr to log file
func RedirectStderr(logWriter *StdWriter) {
	r, w, err := os.Pipe()
	if err != nil {
		logrus.Errorf("Failed to create pipe for stderr redirection: %v", err)
		return
	}
	os.Stderr = w

	go func() {
		_, err := io.Copy(logWriter, r)
		if err != nil {
			logrus.Errorf("Failed to copy stderr to log writer: %v", err)
		}
		err = r.Close()
		if err != nil {
			logrus.Errorf("Failed to close pipe reader: %v", err)
		}
	}()
}
