package models

type ReqParams struct {
	ShellName string   `form:"name" json:"name"`
	CmdOpts   []string `form:"opts" json:"opts"`
}
