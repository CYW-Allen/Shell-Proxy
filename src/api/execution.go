package api

import (
	"fmt"

	"shellProxy/handlers"
	"shellProxy/models"

	"github.com/gin-gonic/gin"
)

func RunScript(ctx *gin.Context) {
	var reqParams models.ReqParams

	if handlers.ExamReq(&reqParams, ctx) {
		if models.WorkList[reqParams.ShellName] == nil {
			models.WorkList[reqParams.ShellName] = new(models.ExeShell)
		}
		curExec := models.WorkList[reqParams.ShellName]

		switch execStatus := curExec.GetCurExecStatus(); execStatus {
		case "":
			handlers.StartExecution(models.WorkList[reqParams.ShellName], ctx, reqParams)
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
}

func GetStatus(ctx *gin.Context) {
	var reqParams models.ReqParams

	if handlers.ExamReq(&reqParams, ctx) {
		reqExec := models.WorkList[reqParams.ShellName]
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
}
