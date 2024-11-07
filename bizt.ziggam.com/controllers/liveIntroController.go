package controllers

import (
	"github.com/beego/beego/v2/core/logs"
)

type LiveIntroController struct {
	BaseController
}

func (c *LiveIntroController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/login")
	}

	c.Data["TMenuId"] = "L00"
	c.Data["SMenuId"] = "L00"

	c.TplName = "live/intro.html"
}
