package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type CommonStandbyEmailController struct {
	beego.Controller
}

func (c *CommonStandbyEmailController) Get() {

	pMemNo := c.GetString("mem_no")

	if pMemNo == "" {
		c.Ctx.Redirect(302, "/common/login")
	}

	c.Data["MemNo"] = pMemNo
	c.TplName = "common/email_auth_standby.html"
}
