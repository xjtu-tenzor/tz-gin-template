package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"template/logger"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 这个是给TraceHook使用的
// 设置了debug模式转换成trace模式
// 设置是否要在trace下输出gin框架网络请求的日志
var SkipSignalChan = make(chan struct{})

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

// 自定义日志输出样式
type CustomFormatter struct {
	logrus.TextFormatter
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = 36 // Cyan
	case logrus.InfoLevel:
		levelColor = 32 // Green
	case logrus.WarnLevel:
		levelColor = 33 // Yellow
	case logrus.ErrorLevel:
		levelColor = 31 // Red
	case logrus.FatalLevel, logrus.PanicLevel:
		levelColor = 35 // Magenta
	}

	// 设置日志消息的颜色
	entry.Message = fmt.Sprintf("\033[%dm%s\033[0m", levelColor, entry.Message)

	// 设置字段名称的颜色
	coloredData := make(logrus.Fields)
	for k, v := range entry.Data {
		var fieldColor int
		switch k {
		case "\nmethod":
			fieldColor = 34 // Blue
		case "\nurl":
			fieldColor = 32 // Green

		case "\nclient_ip":
			fieldColor = 36 // Cyan
		case "\nuser_agent":
			fieldColor = 35 // Magenta

		case "\nstatus":
			fieldColor = 31 // Red

		case "\nrequest_headers":
			fieldColor = 33 // Yellow
		case "\nrequest_body":
			fieldColor = 33

		case "\nresponse_headers":
			fieldColor = 34
		case "\nresponse_body":
			fieldColor = 34

		case "\nduration":
			fieldColor = 111 // Bright Yellow
		default:
			fieldColor = levelColor // 使用日志级别的颜色
		}
		coloredData[fmt.Sprintf("\033[%dm%s\033[0m", fieldColor, k)] = v
	}
	entry.Data = coloredData

	return f.TextFormatter.Format(entry)
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
	logger.SetFormatter(
		&CustomFormatter{
			TextFormatter: logrus.TextFormatter{
				ForceColors:     true,
				FullTimestamp:   true,
				TimestampFormat: config.TimeFormat,
			},
		},
	)

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

	// Add hooks
	//调试时候可以使用这个,将debug转换为trace模式，有堆栈信息
	//logger.AddHook(&TraceHook{})
	//推送日志到服务器
	//logger.AddHook(&RemoteHook{Endpoint: "http://localhost:8080/log"})
	//...

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
