package controllers

import (
	"github.com/astaxie/beego/logs"
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
