package main

import (
	"fmt"
	"log"
	"time"

	"shellProxy/data_defs"
	"shellProxy/handlers"

	"github.com/gin-gonic/gin"
)

func startExecution(wl *data_defs.WorkList, ctx *gin.Context, rp data_defs.ReqParams) {
	go handlers.ExecShell(wl, rp, ctx)
	handlers.SendResponse(
		ctx,
		200,
		fmt.Sprintf("Script [%s] is start to run", rp.ShellName),
		[]string{},
		fmt.Sprintf("%s - (startExecution) Log: Accept request [%s]; status: start", ctx.ClientIP(), rp.ShellName),
	)
}

func checkTTLBeforeExec(wl *data_defs.WorkList, ctx *gin.Context, rp data_defs.ReqParams, execResult string) {
	log.Printf(
		"now: %v, ttl: %v\n",
		time.Now().Format(data_defs.TIMEFORMAT),
		wl.Executions[rp.ShellName].TTL.Format(data_defs.TIMEFORMAT),
	)
	if !wl.CheckShellTTL(rp.ShellName) {
		log.Println("Cooling down finished!")
		startExecution(wl, ctx, rp)
	} else {
		handlers.SendResponse(
			ctx,
			200,
			fmt.Sprintf("The execution [%s] is %s", rp.ShellName, execResult),
			wl.Executions[rp.ShellName].Logs,
			fmt.Sprintf("%s - Log: The execution [%s] is %s; Deny reason: cooling down", ctx.ClientIP(), rp.ShellName, execResult),
		)
	}
}

func main() {
	gin.SetMode((gin.ReleaseMode))
	router := gin.New()
	router.SetTrustedProxies(nil)

	workList := data_defs.WorkList{Executions: make(map[string]data_defs.ExeShell)}

	router.GET("/shell", func(ctx *gin.Context) {
		var reqParams data_defs.ReqParams

		if handlers.ExamReq(&reqParams, ctx) {
			switch execStatus := workList.GetCurExecStatus(reqParams.ShellName); execStatus {
			case "":
				startExecution(&workList, ctx, reqParams)
			case "running":
				handlers.SendResponse(
					ctx,
					200,
					fmt.Sprintf("Script [%s] is still running", reqParams.ShellName),
					workList.Executions[reqParams.ShellName].Logs,
					fmt.Sprintf("%s - Log: Deny request [%s]; Deny reason: still running", ctx.ClientIP(), reqParams.ShellName),
				)
			case "complete":
				checkTTLBeforeExec(&workList, ctx, reqParams, "complete")
			case "fail":
				checkTTLBeforeExec(&workList, ctx, reqParams, "failed")
			}
		}
	})

	router.Run(":8080")
}
