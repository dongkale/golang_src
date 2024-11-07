package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
)

type AdminSimplePassController struct {
	beego.Controller
}

func (c *AdminSimplePassController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	mem_id := c.GetString("mem_id")
	c.Data["MemID"] = mem_id

	c.TplName = "common/pass.html"
}
