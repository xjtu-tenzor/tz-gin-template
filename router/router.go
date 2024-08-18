package router

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"template/middleware"
)

func InitRouter(r *gin.Engine) {
	r.Use(middleware.Error)
	r.Use(middleware.GinLogger(), middleware.GinRecovery(true))
	apiRouter := r.Group("/api")
	{
		// example
		// begin
		apiRouter.GET("/", ctr.Hello.Hello)
		apiRouter.GET("/time", ctr.Hello.HelloTime)
		// end
	}
	v1 := apiRouter.Group("/v1")
	{
		v1.POST("/test", func(c *gin.Context) {
			var req TestRequest
			if c.Request.Body == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is empty"})
				return
			}
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 重置请求体
			if err := c.ShouldBind(&req); err != nil {
				logrus.Errorf("Binding error: %v", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})
	}
}

type TestRequest struct {
	Name string `json:"name" binding:"required"`
}
