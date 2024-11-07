package controllers

import (
	"github.com/astaxie/beego/logs"
)

type Contents2Controller struct {
	BaseController
}

func (c *Contents2Controller) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	c.TplName = "partial/contents2.html"
}
