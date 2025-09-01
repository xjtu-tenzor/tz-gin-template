package middleware

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"template/config"
	"template/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
		// logContext := *c
		logContext := struct {
			Writer   gin.ResponseWriter
			Request  *http.Request
			ClientIP func() string
		}{
			Writer:   c.Writer,
			Request:  c.Request,
			ClientIP: c.ClientIP,
		}

		c.Next()
		go func() {
			select {
			case <-config.SkipSignalChan:
				return
			default:
				goto log
			}
		log:
			status := logContext.Writer.Status()
			path := logContext.Request.URL.Path
			query := logContext.Request.URL.RawQuery
			cost := time.Since(start)
			method := logContext.Request.Method
			clientIP := logContext.ClientIP()
			userAgent := logContext.Request.UserAgent()

			if logger.GinLogger.Level == logrus.DebugLevel {
				if logContext.Writer != nil {
					responseHeaders := logContext.Writer.Header()
					responseBody := w.Body.Bytes()
					requestHeaders, _ := httputil.DumpRequest(logContext.Request, false)

					logger.GinLogger.WithFields(logrus.Fields{
						"\nmethod":           method,
						"\nurl":              path,
						"\nquery":            query,
						"\nclient_ip":        clientIP,
						"\nuser_agent":       userAgent,
						"\nstatus":           status,
						"\nduration":         cost,
						"\nrequest_headers":  string(requestHeaders),
						"\nrequest_body":     string(requestBody),
						"\nresponse_headers": responseHeaders,
						"\nresponse_body":    string(responseBody),
					}).Debug("Debug level log with detailed information")
				}
			} else {
				switch {
				case status >= http.StatusInternalServerError:
					logger.GinLogger.WithFields(logrus.Fields{
						"\nmethod":     method,
						"\nurl":        path,
						"\nquery":      query,
						"\nclient_ip":  clientIP,
						"\nuser_agent": userAgent,
						"\nStatus":     status,
						"\nduration":   cost}).Error("Error level log with brief information")
				case status >= http.StatusBadRequest:
					logger.GinLogger.WithFields(logrus.Fields{
						"\nmethod":     method,
						"\nurl":        path,
						"\nquery":      query,
						"\nclient_ip":  clientIP,
						"\nuser_agent": userAgent,
						"\nstatus":     status,
						"\nduration":   cost}).Warn("Warn level log with brief information")
				default:
					logger.GinLogger.WithFields(logrus.Fields{
						"\nmethod":     method,
						"\nurl":        path,
						"\nquery":      query,
						"\nclient_ip:": clientIP,
						"\nuser_agent": userAgent,
						"\nstatus":     status,
						"\nduration":   cost}).Info("Info level log with brief information")
				}
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

				//deeper stack
				pc, file, line, _ := runtime.Caller(3)
				func_name := runtime.FuncForPC(pc).Name()

				if stack {
					logger.GinLogger.Errorf("panic recovered: %v. Request: %s. File: %s, Line: %d, Function: %s. Stack: %s", err, string(httpRequest), file, line, func_name, string(debug.Stack()))
				} else {
					logger.GinLogger.Errorf("panic recovered: %v. Request: %s. File: %s, Line: %d, Function: %s", err, string(httpRequest), file, line, func_name)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
