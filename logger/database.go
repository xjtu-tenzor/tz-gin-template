package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

var DatabaseLogger *logrus.Logger

// implement io.Writer for logrus, to compatant with gorm logger
type DataLogger struct {
	*logrus.Logger
}

func (dl DataLogger) Write(p []byte) (n int, err error) {
	dl.Out.Write([]byte(fmt.Sprintf("%s", string(p))))
	return len(p), nil
}
