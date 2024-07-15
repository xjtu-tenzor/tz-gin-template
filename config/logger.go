package config

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

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
	fmt.Println(logFilePath)
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
	return lw.Logger.Out.Write(p)
}

func initLogger() {
	config := readLogConfig()
	DatabaseLogger = createLogger(config.DbLogFile, config.LogOutput, config)
	GinLogger = createLogger(config.GinLogFile, config.LogOutput, config)
}
