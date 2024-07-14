package main

import (
	"fmt"
	"template/config"
	"template/router"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()
	gin.SetMode(config.Config.AppMode)
	srv := router.NewServer()

	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("fail to init server: %s\n", err.Error())
		panic(err)
	}

}
