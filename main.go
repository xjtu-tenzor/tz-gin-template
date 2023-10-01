package main

import (
	"template/config"
	"template/router"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(config.Config.AppMode)

	srv := router.NewServer()

	srv.ListenAndServe()

}
