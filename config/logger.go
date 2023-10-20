package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

func createLogger(prefix string) io.Writer {
	err := os.Mkdir("log", 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Println("创建文件夹失败：", err)
		panic(err)
	}
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dateString := time.Now().Format("20060102")
	fileName := path.Join(dir, "log", fmt.Sprintf("%s.log.%s", prefix, dateString))
	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	go func() {
		for {
			currentDateString := time.Now().Format("20060102")
			if currentDateString != dateString {
				dateString = currentDateString
				newFileName := path.Join(dir, "log", fmt.Sprintf("%s.log.%s", prefix, dateString))
				logFile.Close()
				logFile, err = os.OpenFile(newFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					log.Fatalf("Failed to open log file: %v", err)
				}
			}
			time.Sleep(time.Second * 10)
		}
	}()

	return logFile
}

var GinLogger io.Writer
var DatabaseLogger io.Writer

func initLogger() {
	if Config.AppProd {
		GinLogger = createLogger("gin")
		DatabaseLogger = createLogger("database")

		gin.DisableConsoleColor()
		gin.DefaultWriter = GinLogger
	}
}
