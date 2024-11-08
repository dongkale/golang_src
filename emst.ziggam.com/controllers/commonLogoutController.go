package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type CommonLogoutController struct {
	beego.Controller
}

func (c *CommonLogoutController) Get() {

	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	// start : session
	// 1. dvc_id checking
	//pLang, _ := beego.AppConfig.String("lang")
	session := c.StartSession()

	// Delete the session
	session.Delete(c.Ctx.Request.Context(), "mem_no")
	session.Delete(c.Ctx.Request.Context(), "mem_id")

	// 메인으로 이동
	c.Ctx.Redirect(302, "/common/login")
}
