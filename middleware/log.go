package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"template/config"
	"time"
)

const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
)

func logWithLevel(level string, message string) {
	fmt.Printf("Logging level: %s, message: %s\n", level, message) // Debug print
	timestamp := time.Now().Format(time.RFC3339)
	logMessage := fmt.Sprintf("[%s] %s: %s\n", timestamp, level, message)
	if _, err := config.GinLogger.Write([]byte(logMessage)); err != nil {
		fmt.Printf("Failed to write to log file: %v\n", err) // Debug print for file write failure
		return
	}
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next() // Process request

		status := c.Writer.Status()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		cost := time.Since(start)

		level := INFO // Default to INFO level
		if status >= 400 && status < 500 {
			level = WARN // Set to WARN for client errors
		} else if status >= 500 {
			level = ERROR // Set to ERROR for server errors
		}

		// Log the request with the determined level
		logWithLevel(level, fmt.Sprintf("method:%s ;url:%s ;query:%s ;ClientIP:%s ;UserAgent:%s ;Status:%d ;Duration:%v",
			c.Request.Method, path, query, c.ClientIP(), c.Request.UserAgent(), status, cost))
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
					logWithLevel(ERROR, fmt.Sprintf("broken pipe: %s. Request: %s", err, string(httpRequest)))
					c.Abort()
					return
				}
				if stack {
					logWithLevel(ERROR, fmt.Sprintf("panic recovered: %s. Request: %s. Stack: %s", err, string(httpRequest), string(debug.Stack())))
				} else {
					logWithLevel(ERROR, fmt.Sprintf("panic recovered: %s. Request: %s", err, string(httpRequest)))
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
