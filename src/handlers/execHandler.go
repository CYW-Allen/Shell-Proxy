package handlers

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"runtime"
	"shellProxy/data_defs"
	"sync"

	"github.com/gin-gonic/gin"
)

func connectShellPipe(wl *data_defs.WorkList, shellName string, stdPipe io.ReadCloser, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(stdPipe)

	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("[%s] LOG: %s\n", shellName, line)
		wl.Executions[shellName] = data_defs.ExeShell{
			Status: "running",
			Logs:   append(wl.Executions[shellName].Logs, line),
		}
	}
}

func ExecShell(wl *data_defs.WorkList, reqParams data_defs.ReqParams, c *gin.Context) {
	var wg sync.WaitGroup
	var executor string

	if runtime.GOOS == "windows" {
		executor = "powershell.exe"
	} else {
		executor = "/bin/sh"
	}

	wl.SetCurExecStatus(reqParams.ShellName, "running", true)

	shellCmd := exec.Command(executor, append([]string{reqParams.ShellName}, reqParams.CmdOpts...)...)
	shellStdoutPipe, stdoutPipeErr := shellCmd.StdoutPipe()
	shellStderrPipe, stderrPipeErr := shellCmd.StderrPipe()
	if stdoutPipeErr != nil || stderrPipeErr != nil {
		log.Println("(ExecShell) Fail to create std pipe: ", stdoutPipeErr, stderrPipeErr)
		return
	}

	wg.Add(2)
	go connectShellPipe(wl, reqParams.ShellName, shellStdoutPipe, &wg)
	go connectShellPipe(wl, reqParams.ShellName, shellStderrPipe, &wg)

	log.Printf("%s - (ExecShell) LOG: Start to run shell [%s]\n", c.ClientIP(), reqParams.ShellName)

	shellErr := shellCmd.Run()

	if shellErr != nil {
		wl.SetCurExecStatus(reqParams.ShellName, "fail", false)
		log.Printf("%s - (ExecShell) ERR: Shell [%s] execution failed; %s\n", c.ClientIP(), reqParams.ShellName, shellErr.Error())
	} else {
		wl.SetCurExecStatus(reqParams.ShellName, "complete", false)
		log.Printf("%s - (ExecShell) LOG: Finish shell [%s] execution\n", c.ClientIP(), reqParams.ShellName)
	}

	log.Printf("%s - (ExecShell) LOG: Lock shell [%s] for 30s\n", c.ClientIP(), reqParams.ShellName)
	wg.Wait()
}
