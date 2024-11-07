package controllers

import (
	"github.com/beego/beego/v2/core/logs"
)

type JoinCompleteController struct {
	BaseController
}

func (c *JoinCompleteController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	c.TplName = "join/join_complete.html"
}
