package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"template/config"
	"time"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		cost := time.Since(start)

		level := logrus.InfoLevel
		if status >= 400 && status < 500 {
			level = logrus.WarnLevel
		} else if status >= 500 {
			level = logrus.ErrorLevel
		}

		// Log the request with the determined level
		config.GinLogger.Log(level, "method:", c.Request.Method, "; url:", path, "; query:", query, "; ClientIP:", c.ClientIP(), "; UserAgent:", c.Request.UserAgent(), "; Status:", status, "; Duration:", cost)
	}
}

func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					config.GinLogger.Error("broken pipe: ", err, ". Request: ", string(httpRequest))
					c.Abort()
					return
				}
				if stack {
					config.GinLogger.Error("panic recovered: ", err, ". Request: ", string(httpRequest), ". Stack: ", string(debug.Stack()))
				} else {
					config.GinLogger.Error("panic recovered: ", err, ". Request: ", string(httpRequest))
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
