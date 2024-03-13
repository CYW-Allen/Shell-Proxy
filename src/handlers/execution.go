package handlers

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"runtime"
	"shellProxy/models"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const TIMEFORMAT = "2006-01-02 15:04:05"

func connectShellPipe(shellName string, es *models.ExeShell, stdPipe io.ReadCloser, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(stdPipe)

	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("[%s] LOG: %s\n", shellName, line)
		es.UpdateShellLogs(line)
	}
}

func execShell(shellName string, es *models.ExeShell, cmdOpts []string, c *gin.Context) {
	var wg sync.WaitGroup
	var executor string

	if runtime.GOOS == "windows" {
		executor = "powershell.exe"
	} else {
		executor = "/bin/sh"
	}

	es.SetCurExecStatus("running", true)

	shellCmd := exec.Command(executor, append([]string{shellName}, cmdOpts...)...)
	shellStdoutPipe, stdoutPipeErr := shellCmd.StdoutPipe()
	shellStderrPipe, stderrPipeErr := shellCmd.StderrPipe()
	if stdoutPipeErr != nil || stderrPipeErr != nil {
		log.Println("(execShell) Fail to create std pipe: ", stdoutPipeErr, stderrPipeErr)
		return
	}

	wg.Add(2)
	go connectShellPipe(shellName, es, shellStdoutPipe, &wg)
	go connectShellPipe(shellName, es, shellStderrPipe, &wg)

	log.Printf("%s - (execShell) LOG: Start to run shell [%s]\n", c.ClientIP(), shellName)

	shellErr := shellCmd.Run()

	if shellErr != nil {
		es.SetCurExecStatus("fail", false)
		log.Printf("%s - (execShell) ERR: Shell [%s] execution failed; %s\n", c.ClientIP(), shellName, shellErr.Error())
	} else {
		es.SetCurExecStatus("complete", false)
		log.Printf("%s - (execShell) LOG: Finish shell [%s] execution\n", c.ClientIP(), shellName)
	}

	log.Printf("%s - (execShell) LOG: Lock shell [%s] for 30s\n", c.ClientIP(), shellName)
	wg.Wait()
}

func StartExecution(es *models.ExeShell, ctx *gin.Context, rp models.ReqParams) {
	go execShell(rp.ShellName, es, rp.CmdOpts, ctx)
	SendResponse(
		ctx,
		200,
		fmt.Sprintf("Script [%s] is start to run", rp.ShellName),
		[]string{},
		fmt.Sprintf("%s - (startExecution) Log: Accept request [%s]; status: start", ctx.ClientIP(), rp.ShellName),
	)
}

func CheckTTLBeforeExec(es *models.ExeShell, ctx *gin.Context, rp models.ReqParams, execResult string) {
	log.Printf(
		"now: %v, ttl: %v\n",
		time.Now().Format(TIMEFORMAT),
		es.GetShellTTL().Format(TIMEFORMAT),
	)
	if !es.CheckShellTTL() {
		log.Println("Cooling down finished!")
		StartExecution(es, ctx, rp)
	} else {
		SendResponse(
			ctx,
			200,
			fmt.Sprintf("The execution [%s] is %s", rp.ShellName, execResult),
			es.GetShellLogs(),
			fmt.Sprintf("%s - Log: The execution [%s] is %s; Deny reason: cooling down", ctx.ClientIP(), rp.ShellName, execResult),
		)
	}
}
