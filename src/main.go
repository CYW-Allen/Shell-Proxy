package main

import (
	"fmt"

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
			if data_defs.WorkList[reqParams.ShellName] == nil {
				data_defs.WorkList[reqParams.ShellName] = new(data_defs.ExeShell)
			}
			curExec := data_defs.WorkList[reqParams.ShellName]

			switch execStatus := curExec.GetCurExecStatus(); execStatus {
			case "":
				handlers.StartExecution(data_defs.WorkList[reqParams.ShellName], ctx, reqParams)
			case "running":
				handlers.SendResponse(
					ctx,
					200,
					fmt.Sprintf("Script [%s] is still running", reqParams.ShellName),
					curExec.GetShellLogs(),
					fmt.Sprintf("%s - Log: Deny request [%s]; Deny reason: still running", ctx.ClientIP(), reqParams.ShellName),
				)
			case "complete":
				handlers.CheckTTLBeforeExec(curExec, ctx, reqParams, "complete")
			case "fail":
				handlers.CheckTTLBeforeExec(curExec, ctx, reqParams, "failed")
			}
		}
	})

	router.GET("/status", func(ctx *gin.Context) {
		var reqParams data_defs.ReqParams

		if handlers.ExamReq(&reqParams, ctx) {
			reqExec := data_defs.WorkList[reqParams.ShellName]
			if reqExec == nil {
				handlers.SendResponse(
					ctx,
					404,
					fmt.Sprintf("There is no record of the execution [%s]", reqParams.ShellName),
					[]string{},
					fmt.Sprintf("%s - Log: Fail to get the record of the execution [%s]", ctx.ClientIP(), reqParams.ShellName),
				)
			} else {
				handlers.SendResponse(
					ctx,
					200,
					fmt.Sprintf(
						"Current execution [%s] status: %s",
						reqParams.ShellName,
						reqExec.GetCurExecStatus(),
					),
					reqExec.GetShellLogs(),
					fmt.Sprintf("%s - Log: Get the execution [%s] status", ctx.ClientIP(), reqParams.ShellName),
				)
			}
		}
	})

	fmt.Println("The api server is running on port 8080")
	router.Run(":8080")
}
