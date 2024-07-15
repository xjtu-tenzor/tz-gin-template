package config

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"

	"os"
	"path/filepath"
	"time"
)

type LogConfig struct {
	LogLevel   string `json:"log_level"`
	LogOutput  string `json:"log_output"`
	GinLogFile string `json:"gin_log_file"`
	DbLogFile  string `json:"db_log_file"`
	LogMaxSize int    `json:"log_max_size"`
	TimeFormat string `json:"time_format"`
	LogFormat  string `json:"log_format"`
	MaxAge     int    `json:"max_age"`
	MaxBackups int    `json:"max_backups"`
	Compress   bool   `json:"compress"`
}

func readLogConfig() LogConfig {
	file, err := os.ReadFile("./log.json")
	if err != nil {
		logrus.Fatalf("Failed to read log config: %v", err)
	}

	var config LogConfig
	err = json.Unmarshal(file, &config)
	if err != nil {
		logrus.Fatalf("Failed to unmarshal log config: %v", err)
	}
	return config
}

func createLogger(logFilePrefix, logOutputDir string, config LogConfig) *logrus.Logger {
	logDir := filepath.Join(".", logOutputDir)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		logrus.Errorf("Failed to create log directory: %v", err)
		return nil // Return nil to avoid application termination
	}

	dateString := time.Now().Format(config.TimeFormat)
	logFilePath := filepath.Join(logDir, fmt.Sprintf("%s.%s.log", logFilePrefix, dateString))
	logger := logrus.New()

	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logrus.Errorf("Invalid log level: %v", err)
		return nil // Return nil to avoid application termination
	}
	logger.SetLevel(level)

	logger.Out = &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    config.LogMaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	return logger
}

var GinLogger *logrus.Logger
var DatabaseLogger *logrus.Logger

// implement io.Writer for logrus, to compatant with gorm logger
type LogWriter struct {
	*logrus.Logger
}

func (lw LogWriter) Write(p []byte) (n int, err error) {
	lw.Logger.Info(string(p))
	return len(p), nil
}

func initLogger() {
	config := readLogConfig()
	DatabaseLogger = createLogger(config.DbLogFile, config.LogOutput, config)
	GinLogger = createLogger(config.GinLogFile, config.LogOutput, config)

	stderrLogger := createLogger("stderr", config.LogOutput, config)
	if stderrLogger != nil {
		stderrWriter := &LogWriter{Logger: stderrLogger}
		redirectStderr(stderrWriter)
	} else {
		logrus.Errorf("Failed to create stderr logger")
	}
}

// capture stderr to log file
func redirectStderr(logWriter *LogWriter) {
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
