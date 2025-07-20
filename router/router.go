package router

import (
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

		// std package examples
		// begin
		stdRouter := apiRouter.Group("/std")
		{
			stdRouter.GET("/function", ctr.Std.FunctionDemo)
			stdRouter.GET("/bind", ctr.Std.BindDemo)
			stdRouter.GET("/forward", ctr.Std.ForwardDemo)
			stdRouter.GET("/functional", ctr.Std.FunctionalDemo)
			stdRouter.GET("/advanced", ctr.Std.AdvancedDemo)
		}
		// end
	}
}
