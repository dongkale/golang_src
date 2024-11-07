package controllers

import (
	"github.com/beego/beego/v2/core/logs"
)

type TeamMemberWriteController struct {
	BaseController
}

func (c *TeamMemberWriteController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	mem_id := session.Get(c.Ctx.Request.Context(), "mem_id")
	if mem_id == nil {
		c.Ctx.Redirect(302, "/login")
	}

	c.Data["TMenuId"] = "T00"
	c.Data["SMenuId"] = "T00"

	c.TplName = "team/team_member_write.html"
}
