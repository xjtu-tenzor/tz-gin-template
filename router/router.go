package router

import (
	"template/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.Use(middleware.Error)

	apiRouter := r.Group("/api")
	{
		// example
		// begin
		apiRouter.GET("/", ctr.Hello.Hello)
		apiRouter.GET("/time", ctr.Hello.HelloTime)
		// end
	}

}
