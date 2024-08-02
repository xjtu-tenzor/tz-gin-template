package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"template/logger"
	"time"
)

type LogConfig struct {
	LogLevel   string `json:"log_level"`
	LogOutput  string `json:"log_output"`
	GinLogFile string `json:"gin_log_file"`
	DbLogFile  string `json:"db_log_file"`
	LogMaxSize int    `json:"log_max_size"`
	TimeFormat string `json:"time_format"`
	MaxAge     int    `json:"max_age"`
	MaxBackups int    `json:"max_backups"`
	Compress   bool   `json:"compress"`
}

func createLogger(logFilePrefix, logOutputDir string, config LogConfig) *logrus.Logger {
	logDir := filepath.Join(".", logOutputDir)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		logrus.Errorf("Failed to create log directory: %v", err)
		return nil
	}

	dateString := time.Now().Format(config.TimeFormat)
	logFilePath := filepath.Join(logDir, fmt.Sprintf("%s.%s.log", logFilePrefix, dateString))
	logger := logrus.New()

	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logrus.Errorf("Invalid log level: %v", err)
		return nil
	}
	logger.SetLevel(level)

	logOutput := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    config.LogMaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	// If the log level is debug, log to both file and console
	if Config.AppMode == "debug" {
		logger.Out = io.MultiWriter(logOutput, os.Stdout)
	} else {
		logger.Out = logOutput
	}

	return logger
}

func initLogger() {
	config := LogConfig{
		LogLevel:   Config.LogLevel,
		LogOutput:  "./log",
		GinLogFile: "gin",
		DbLogFile:  "database",
		LogMaxSize: 512,
		TimeFormat: "20060102",
		MaxAge:     7,
		MaxBackups: 5,
		Compress:   true,
	}
	logger.DatabaseLogger = createLogger(config.DbLogFile, config.LogOutput, config)
	logger.GinLogger = createLogger(config.GinLogFile, config.LogOutput, config)

	stderrLogger := createLogger("stderr", config.LogOutput, config)
	if stderrLogger != nil {
		stderrWriter := &logger.StdWriter{Logger: stderrLogger}
		logger.RedirectStderr(stderrWriter)
	} else {
		logrus.Errorf("Failed to create stderr logger")
	}
}
