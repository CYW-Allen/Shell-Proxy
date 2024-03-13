package main

import (
	"fmt"

	"shellProxy/api"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode((gin.ReleaseMode))
	router := gin.New()
	router.SetTrustedProxies(nil)

	router.GET("/shell", api.RunScript)
	router.GET("/status", api.GetStatus)

	fmt.Println("The api server is running on port 8080")
	router.Run(":8080")
}
