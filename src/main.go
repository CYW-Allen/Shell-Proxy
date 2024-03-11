package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

type ReqParams struct {
	ShellName string   `form:"name" json:"name"`
	CmdOpts   []string `form:"opts" json:"opts"`
}

func sendResponse(ctx *gin.Context, statusCode int, result string, shellLogs []string, proxyLog string) {
	ctx.JSON(statusCode, gin.H{
		"result": result,
		"logs":   shellLogs,
	})
	log.Println(proxyLog)
}

func examReq(reqParams *ReqParams, ctx *gin.Context) bool {
	if err := ctx.BindQuery(reqParams); err != nil || reqParams.ShellName == "" {
		sendResponse(
			ctx,
			400,
			"Invalid parameters for shell execution",
			nil,
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
		var reqParams ReqParams

		if examReq(&reqParams, ctx) {
			sendResponse(
				ctx,
				200,
				fmt.Sprintf("Success to run the shell %s", reqParams.ShellName),
				nil,
				fmt.Sprintf("Run the script: %s %s", reqParams.ShellName, strings.Join(reqParams.CmdOpts, " ")),
			)

		}
	})

	router.Run(":8080")
}
