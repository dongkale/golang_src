package controllers

import (
	"github.com/astaxie/beego/logs"
)

type ErrorController struct {
	BaseController
}

func (c *ErrorController) Error404() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	c.TplName = "error/error.html"
}
