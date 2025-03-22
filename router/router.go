package router

import (
	"fmt"
	"os"
	"template/middleware"

	"github.com/gin-gonic/gin"
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

	test_r := r.Group("test")
	{
		test_r.GET("/v1", func(c *gin.Context) {
			fmt.Fprintln(os.Stderr, "test")
		})
	}
}
