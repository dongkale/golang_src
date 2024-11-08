package controllers

import (
	"github.com/beego/beego/v2/core/logs"
)

type Error404Controller struct {
	BaseController
}

func (c *Error404Controller) Get() {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	// start : session
	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/common/login")
	}
	//pLang, _ := beego.AppConfig.String("lang")
	// end : session
	//log.Debug("url : %v", request.Header.Get("Referer"))

	// pErrorMessage := c.GetString("em")
	// if pErrorMessage == "" {
	// 	pErrorMessage = "해당 페이지를 찾을 수 없습니다."
	// }

	// c.Data["ErrorMessage"] = pErrorMessage
	c.TplName = "error/404.html"
}
