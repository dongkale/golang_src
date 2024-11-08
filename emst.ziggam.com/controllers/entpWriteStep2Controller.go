package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
)

type EntpWriteStep2Controller struct {
	beego.Controller
}

func (c *EntpWriteStep2Controller) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	//session := c.StartSession()
	//pLang, _ := beego.AppConfig.String("lang")

	pEmailRecvYn := c.GetString("email_recv_yn")

	if pEmailRecvYn != "Y" {
		pEmailRecvYn = "N"
	}

	c.Data["EmailRecvYn"] = pEmailRecvYn
	c.TplName = "entp/entp_write2.html"
}
