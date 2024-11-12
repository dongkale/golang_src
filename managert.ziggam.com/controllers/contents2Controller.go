package controllers

import (
	"github.com/beego/beego/v2/core/logs"
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
