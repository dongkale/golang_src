package controllers

import (
	// beego "github.com/beego/beego/v2/server/web"
)

type ErrorController struct {
	BaseController
}

func (c *ErrorController) Error404() {
	// log := logs.NewLogger()
	// log.SetLogger(logs.AdapterConsole)

	// start : session
	// 1. dvc_id checking
	//session := c.StartSession()
	//pLang := beego.AppConfig.String("lang")
	// end : session
	//log.Debug("url : %v", request.Header.Get("Referer"))

	/*
		pErrorMessage := c.GetString("em")
		if pErrorMessage == "" {
			pErrorMessage = "해당 페이지를 찾을 수 없습니다."
		}
	*/

	c.Data["ErrorMessage"] = "error"
	c.TplName = "error/error.html"
}
