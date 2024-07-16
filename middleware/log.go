package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"template/config"
	"time"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		cCp := c.Copy()
		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}

		go func() {
			status := cCp.Writer.Status()
			path := cCp.Request.URL.Path
			query := cCp.Request.URL.RawQuery
			cost := time.Since(start)
			method := cCp.Request.Method
			clientIP := cCp.ClientIP()
			userAgent := cCp.Request.UserAgent()

			level := logrus.InfoLevel
			if status >= 400 && status < 500 {
				level = logrus.WarnLevel
			} else if status >= 500 {
				level = logrus.ErrorLevel
			}

			if config.GinLogger.Level == logrus.DebugLevel {
				cCp.Writer = w
				responseHeaders := cCp.Writer.Header()
				responseBody := w.body.Bytes()
				var requestBody []byte
				requestBody, _ = ioutil.ReadAll(c.Request.Body)
				cCp.Request.Body = ioutil.NopCloser(strings.NewReader(string(requestBody)))
				requestHeaders, _ := httputil.DumpRequest(c.Request, false)

				config.GinLogger.WithFields(logrus.Fields{
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
			} else {
				config.GinLogger.Log(level,
					"method:", method, ";"+
						" url:", path, ";"+
						" query:", query, "; "+
						"ClientIP:", clientIP, "; "+
						"UserAgent:", userAgent, "; "+
						"Status:", status, "; "+
						"Duration:", cost)
			}
		}()
		c.Next()
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
