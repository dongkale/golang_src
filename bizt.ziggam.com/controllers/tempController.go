package controllers

import (
	"strconv"

	"bizt.ziggam.com/models"	
	"github.com/beego/beego/v2/core/logs"
)

type TempController struct {
	BaseController
}

func (c *TempController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/login")
	}

	logs.Info(strconv.Itoa(models.Room1.AuthCd))

	c.Data["TMenuId"] = "M00"
	c.Data["SMenuId"] = "M00"

	c.TplName = "main/main.html"
}
