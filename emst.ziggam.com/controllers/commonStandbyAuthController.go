package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type CommonStandbyAuthController struct {
	beego.Controller
}

func (c *CommonStandbyAuthController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "common/last_auth_standby.html"
}
