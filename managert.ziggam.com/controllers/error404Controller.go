package controllers

import (
	"github.com/beego/beego/v2/core/logs"
)

type Error404Controller struct {
	BaseController
}

func (c *Error404Controller) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	c.TplName = "error/404.html"
}
