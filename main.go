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

	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("fail to init server: %s\n", err.Error())
		panic(err)
	}

}
