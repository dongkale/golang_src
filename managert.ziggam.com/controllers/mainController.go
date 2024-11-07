package controllers

import (
	"github.com/astaxie/beego/logs"
)

type MainController struct {
	BaseController
}

func (c *MainController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	pBnrGrpTypCd := c.GetString("bnr_grp_typ_cd")

	c.Data["BnrGrpTypCd"] = pBnrGrpTypCd

	c.TplName = "main/main.html"
}
