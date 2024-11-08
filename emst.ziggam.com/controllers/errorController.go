package controllers

import (
	"github.com/beego/beego/v2/core/logs"
)

type ErrorController struct {
	BaseController
}

func (c *ErrorController) Error404() {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	// start : session
	// 1. dvc_id checking
	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		mem_no = ""
	}
	//pLang, _ := beego.AppConfig.String("lang")
	// end : session
	//log.Debug("url : %v", request.Header.Get("Referer"))

	// pErrorMessage := c.GetString("em")
	// if pErrorMessage == "" {
	// 	pErrorMessage = "해당 페이지를 찾을 수 없습니다."
	// }

	// c.Data["ErrorMessage"] = pErrorMessage
	c.TplName = "error/error.html"
}
