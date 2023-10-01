package config

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

func createLogger(name string) io.Writer {
	err := os.Mkdir("log", 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Println("创建文件夹失败：", err)
		panic(err)
	}
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fileName := path.Join(dir, "log", name+".log")
	logger, err := rotatelogs.New(fileName+".%Y%m%d",
		// rotatelogs.WithLinkName(fileName),
		rotatelogs.WithMaxAge(30*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour))
	if err != nil {
		panic(err)
	}

	return logger
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
