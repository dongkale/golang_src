package controllers

import (
	"github.com/beego/beego/v2/core/logs"
)

type LoginController struct {
	BaseController
}

func (c *LoginController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "login/login.html"
}
