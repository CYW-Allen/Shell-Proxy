package data_defs

import (
	"log"
	"sync"
	"time"
)

const TIMEFORMAT = "2006-01-02 15:04:05"
const COOLINGDOWN = 30 * time.Second

type ReqParams struct {
	ShellName string   `form:"name" json:"name"`
	CmdOpts   []string `form:"opts" json:"opts"`
}

type ExeShell struct {
	Rw     sync.RWMutex
	Status string
	Logs   []string
	TTL    time.Time
}

func (es *ExeShell) GetCurExecStatus() string {
	es.Rw.RLock()
	defer es.Rw.RUnlock()
	return es.Status
}

func (es *ExeShell) SetCurExecStatus(newStatus string, isInit bool) {
	es.Rw.Lock()
	defer es.Rw.Unlock()

	es.Status = newStatus

	if isInit {
		es.Logs = []string{}
	}

	if newStatus == "complete" || newStatus == "fail" {
		es.TTL = time.Now().Add(COOLINGDOWN)
	} else {
		es.TTL = time.Now()
	}
}

func (es *ExeShell) GetShellLogs() []string {
	es.Rw.RLock()
	defer es.Rw.RUnlock()
	return es.Logs
}

func (es *ExeShell) UpdateShellLogs(newLog string) {
	es.Rw.Lock()
	defer es.Rw.Unlock()
	es.Logs = append(es.Logs, newLog)
}

func (es *ExeShell) GetShellTTL() time.Time {
	es.Rw.RLock()
	defer es.Rw.RUnlock()
	return es.TTL
}

func (es *ExeShell) refreshShellTTL() {
	es.Rw.Lock()
	defer es.Rw.Unlock()

	es.TTL = time.Now().Add(COOLINGDOWN)
	log.Printf(
		"curTime: %v, new ttl: %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		es.TTL.Format("2006-01-02 15:04:05"),
	)
}

func (es *ExeShell) CheckShellTTL() bool {
	stillInCoolingDown := es.TTL.After(time.Now())
	if stillInCoolingDown {
		es.refreshShellTTL()
	}
	return stillInCoolingDown
}

var WorkList map[string]*ExeShell

func init() {
	WorkList = make(map[string]*ExeShell)
}
