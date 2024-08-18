package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"template/logger"
	"time"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		w := &logger.ResponseBodyWriter{Body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		// Read and store the request body
		var requestBody []byte
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				requestBody = make([]byte, len(bodyBytes))
				copy(requestBody, bodyBytes)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Reset request body
			}
		}

		// Create a copy of the context for logging
		logContext := *c

		c.Next()

		go func() {
			status := logContext.Writer.Status()
			path := logContext.Request.URL.Path
			query := logContext.Request.URL.RawQuery
			cost := time.Since(start)
			method := logContext.Request.Method
			clientIP := logContext.ClientIP()
			userAgent := logContext.Request.UserAgent()

			level := logrus.InfoLevel
			if status >= 400 && status < 500 {
				level = logrus.WarnLevel
			} else if status >= 500 {
				level = logrus.ErrorLevel
			}

			if logger.GinLogger.Level == logrus.DebugLevel {
				if logContext.Writer != nil {
					responseHeaders := logContext.Writer.Header()
					responseBody := w.Body.Bytes()
					requestHeaders, _ := httputil.DumpRequest(logContext.Request, false)

					logger.GinLogger.WithFields(logrus.Fields{
						"method":           method,
						"url":              path,
						"query":            query,
						"clientIP":         clientIP,
						"userAgent":        userAgent,
						"status":           status,
						"duration":         cost,
						"request_headers":  string(requestHeaders),
						"request_body":     string(requestBody),
						"response_headers": responseHeaders,
						"response_body":    string(responseBody),
					}).Debug("Debug level log with detailed information")
				}
			} else {
				logger.GinLogger.Log(level,
					"method:", method, ";"+
						" url:", path, ";"+
						" query:", query, "; "+
						"ClientIP:", clientIP, "; "+
						"UserAgent:", userAgent, "; "+
						"Status:", status, "; "+
						"Duration:", cost)
			}
		}()
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
					logger.GinLogger.Error("broken pipe: ", err, ". Request: ", string(httpRequest))
					c.Abort()
					return
				}
				if stack {
					logger.GinLogger.Error("panic recovered: ", err, ". Request: ", string(httpRequest), ". Stack: ", string(debug.Stack()))
				} else {
					logger.GinLogger.Error("panic recovered: ", err, ". Request: ", string(httpRequest))
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
