package main

import (
	"fmt"
	"log"

	"shellProxy/api"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	log.SetOutput(&lumberjack.Logger{
		Filename: "./server.log",
		MaxSize:  100,
		MaxAge:   30,
	})

	gin.SetMode((gin.ReleaseMode))
	router := gin.New()
	router.SetTrustedProxies(nil)

	router.GET("/shell", api.RunScript)
	router.GET("/status", api.GetStatus)

	fmt.Println("The api server is running on port 8080")
	router.Run(":8080")
}
