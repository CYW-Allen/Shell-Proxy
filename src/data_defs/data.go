package data_defs

import (
	"log"
	"sync"
	"time"
)

type ReqParams struct {
	ShellName string   `form:"name" json:"name"`
	CmdOpts   []string `form:"opts" json:"opts"`
}

type ExeShell struct {
	Status string
	Logs   []string
	TTL    time.Time
}

type WorkList struct {
	Rw         sync.RWMutex
	Executions map[string]ExeShell
}

var TIMEFORMAT = "2006-01-02 15:04:05"
var COOLINGDOWN = 30 * time.Second

func (wl *WorkList) GetCurExecStatus(name string) string {
	wl.Rw.RLock()
	defer wl.Rw.RUnlock()

	return wl.Executions[name].Status
}

func (wl *WorkList) SetCurExecStatus(name string, status string, isInit bool) {
	wl.Rw.Lock()
	defer wl.Rw.Unlock()

	exec := wl.Executions[name]
	exec.Status = status

	if isInit {
		exec.Logs = []string{}
	} else {
		exec.Logs = wl.Executions[name].Logs
	}

	if status == "complete" || status == "fail" {
		exec.TTL = time.Now().Add(COOLINGDOWN)
	} else {
		exec.TTL = time.Now()
	}

	wl.Executions[name] = exec
}

func (wl *WorkList) RemoveCurExec(name string) {
	wl.Rw.Lock()
	defer wl.Rw.Unlock()

	delete(wl.Executions, name)
}

func (wl *WorkList) refreshShellTTL(name string) {
	wl.Rw.Lock()
	defer wl.Rw.Unlock()

	exec := wl.Executions[name]
	exec.TTL = time.Now().Add(COOLINGDOWN)
	wl.Executions[name] = exec
	log.Printf(
		"curTime: %v, new ttl: %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		wl.Executions[name].TTL.Format("2006-01-02 15:04:05"),
	)
}

func (wl *WorkList) CheckShellTTL(name string) bool {
	stillInCoolingDown := wl.Executions[name].TTL.After(time.Now())
	if stillInCoolingDown {
		wl.refreshShellTTL(name)
	}
	return stillInCoolingDown
}
