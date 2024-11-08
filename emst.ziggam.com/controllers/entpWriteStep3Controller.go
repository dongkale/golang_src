package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
)

type EntpWriteStep3Controller struct {
	beego.Controller
}

func (c *EntpWriteStep3Controller) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	//session := c.StartSession()
	//pLang, _ := beego.AppConfig.String("lang")

	pEntpKoNm := c.GetString("entp_ko_nm")
	pBizRegNo := c.GetString("biz_reg_no")
	pRepNm := c.GetString("rep_nm")
	pPpChrgNm := c.GetString("pp_chrg_nm")
	pPpChrgTelNo := c.GetString("pp_chrg_tel_no")
	pEmailRecvYn := c.GetString("email_recv_yn")

	if pEntpKoNm == "" {
		c.Ctx.Redirect(302, "/common/login")
	}

	c.Data["EntpKoNm"] = pEntpKoNm
	c.Data["BizRegNo"] = pBizRegNo
	c.Data["RepNm"] = pRepNm
	c.Data["PpChrgNm"] = pPpChrgNm
	c.Data["PpChrgTelNo"] = pPpChrgTelNo
	c.Data["EmailRecvYn"] = pEmailRecvYn

	c.TplName = "entp/entp_write3.html"
}
