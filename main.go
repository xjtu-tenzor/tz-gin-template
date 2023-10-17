package main

import (
	"fmt"
	"template/config"
	"template/router"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(config.Config.AppMode)
	fmt.Println(config.Config.AppSecret)
	srv := router.NewServer()

	srv.ListenAndServe()

}
