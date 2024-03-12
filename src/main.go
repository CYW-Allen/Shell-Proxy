package main

import (
	"fmt"
	"strings"

	"shellProxy/data_defs"
	"shellProxy/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode((gin.ReleaseMode))
	router := gin.New()
	router.SetTrustedProxies(nil)

	router.GET("/shell", func(ctx *gin.Context) {
		var reqParams data_defs.ReqParams

		if handlers.ExamReq(&reqParams, ctx) {
			handlers.SendResponse(
				ctx,
				200,
				fmt.Sprintf("Success to run the shell %s", reqParams.ShellName),
				[]string{},
				fmt.Sprintf("Run the script: %s %s", reqParams.ShellName, strings.Join(reqParams.CmdOpts, " ")),
			)

		}
	})

	router.Run(":8080")
}
