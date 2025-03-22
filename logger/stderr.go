package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"runtime"
)

type StdWriter struct {
	*logrus.Logger
}

func (sw StdWriter) Write(p []byte) (n int, err error) {
	// Retrieve the file and line number where the error occurred
	pc, file, line, _ := runtime.Caller(2)
	func_name := runtime.FuncForPC(pc).Name()

	sw.Logger.Errorf("Find stderr: %s Location: File: %s, Line: %d, Function Name: %s", string(p), file, line, func_name)
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
