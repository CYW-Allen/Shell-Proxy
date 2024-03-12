package main

import (
	"fmt"
	"strings"

	"shellProxy/data_defs"
	"shellProxy/handlers"

	"github.com/gin-gonic/gin"
)

func examReq(reqParams *data_defs.ReqParams, ctx *gin.Context) bool {
	if err := ctx.BindQuery(reqParams); err != nil || reqParams.ShellName == "" {
		handlers.SendResponse(
			ctx,
			400,
			"Invalid parameters for shell execution",
			[]string{},
			fmt.Sprintf("%s - [examReq] Got invalid parameters", ctx.ClientIP()),
		)
		return false
	}
	return true
}

func main() {
	gin.SetMode((gin.ReleaseMode))
	router := gin.New()
	router.SetTrustedProxies(nil)

	router.GET("/shell", func(ctx *gin.Context) {
		var reqParams data_defs.ReqParams

		if examReq(&reqParams, ctx) {
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
