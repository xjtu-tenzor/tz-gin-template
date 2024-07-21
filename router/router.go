package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
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
	v1 := r.Group("v1")
	{
		v1.GET("/ping", func(context *gin.Context) {
			context.JSON(200, gin.H{
				"message": "pong",
			})
			fmt.Fprintf(os.Stderr, "这是一个故意生成的 stderr 输出\n")
		})
		v1.GET("/error", func(context *gin.Context) {
			context.JSON(404, gin.H{
				"message": "error",
			})
		})
	}
}
