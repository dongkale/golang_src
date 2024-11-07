package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
)

// BridgeController ...
type BridgeController struct {
	beego.Controller
}

// localhost:7070/bridge?entpmemno=E2019011900001&recruitsn=2019010001&reqname=이동관&reqmono=010-5226-2107&reqemail=dongkale@naver.com
// https://biz-test.ziggam.com/bridge?entpmemno=E2019011900001&recruitsn=2019010001//
// https://biz-test.ziggam.com/bridge?entpmemno=E2019011900001&recruitsn=2019010001&reqname=이동관&reqmono=010-5226-2107&reqemail=dongkale@naver.com
// localhost:7070/bridge?entpmemno=E2019011900001&recruitsn=2019010001&reqname=이동관&reqmono=010-5226-2107&reqemail=dongkale@naver.com
// http://localhost:7070/bridge?entpmemno=E2019011900001&recruitsn=2019010001&reqname=%EC%9D%B4%EB%8F%99%EA%B4%80&reqmono=010-5226-2107&reqemail=dongkale@naver.com

// Get ...
func (c *BridgeController) Get() {

	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := c.GetString("entpmemno")
	pRecruitSn := c.GetString("recruitsn")
	pName := c.GetString("reqname")
	pMono := c.GetString("reqmono")
	pEmail := c.GetString("reqemail")

	fmt.Printf(fmt.Sprintf("[BridgeController] lang:%v, entp_mem_no:%v, recruit_sn:%v, name:%v, MoNo:%v, Email:%v",
		pLang, pEntpMemNo, pRecruitSn, pName, pMono, pEmail))

	c.Data["pEntpMemNo"] = pEntpMemNo
	c.Data["pRecruitSn"] = pRecruitSn

	c.Data["RegName"] = pName
	c.Data["ReqMoNo"] = pMono
	c.Data["ReqEmail"] = pEmail

	c.TplName = "bridge/bridge.html"
}
