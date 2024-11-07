package controllers

import (
	"github.com/beego/beego/v2/core/logs"
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

// ErrorTest ...
// c.Abort("DataBase")
func (c *ErrorController) ErrorDataBase() {

	//c.Data["content"] = "database is now down"
	c.TplName = "error/error.html"
}
