package controllers

import (
	"github.com/beego/beego/v2/core/logs"
)

type WithDrawController struct {
	BaseController
}

func (c *WithDrawController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/login")
	}
	auth_cd := session.Get(c.Ctx.Request.Context(), "auth_cd")
	if auth_cd != "01" {
		c.Ctx.Redirect(302, "/")
	}

	c.Data["TMenuId"] = "M00"
	c.Data["SMenuId"] = "M00"

	c.TplName = "setting/withdraw.html"
}
