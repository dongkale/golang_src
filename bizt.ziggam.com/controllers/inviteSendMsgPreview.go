package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
)

// InviteSendMsgPreviewController ...
type InviteSendMsgPreviewController struct {
	BaseController
}

// Get ...
func (c *InviteSendMsgPreviewController) Get() {

	session := c.StartSession()

	memNO := session.Get(c.Ctx.Request.Context(), "mem_no")
	if memNO == nil {
		c.Ctx.Redirect(302, "/login")
	}

	memSn := session.Get(c.Ctx.Request.Context(), "mem_sn")
	if memSn == nil {
		c.Ctx.Redirect(302, "/login")
	}

	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := memNO

	pEntpKoNm := c.GetString("entp_ko_nm")
	pPpChrgName := c.GetString("pp_chrg_nm")
	pUpJobGrp := c.GetString("up_job_grp")
	pJobGrp := c.GetString("job_grp")
	pRecruitTitle := c.GetString("recruit_title")

	pInviteTitle := c.GetString("invite_tit")
	pInviteMsg := c.GetString("invite_msg")

	fmt.Printf(pLang)
	fmt.Printf(pEntpMemNo.(string))
	fmt.Printf(pEntpKoNm)
	fmt.Printf(pPpChrgName)
	fmt.Printf(pUpJobGrp)
	fmt.Printf(pJobGrp)
	fmt.Printf(pRecruitTitle)
	fmt.Printf(pInviteTitle)
	fmt.Printf(pInviteMsg)

	//var convMsg = strings.Replace(pInviteMsg, "{지원자명}", val.Name, 100)

	// var convMsg = strings.Replace(pInviteMsg, "{지원자명}", pEntpKoNm, 100)
	// convMsg = strings.Replace(convMsg, "{채용공고 제목}", pRecruitTitle, 100)
	// convMsg = strings.Replace(convMsg, "{1차직군}", pUpJobGrp, 100)
	// convMsg = strings.Replace(convMsg, "{2차직군}", pJobGrp, 100)
	// // convMsg = strings.Replace(convMsg, "{URL1}", resultRecruitUrl, 100)
	// convMsg = strings.Replace(convMsg, "{URL2}", entpVdUrl, 100)

	c.Data["EntpKoNm"] = pEntpKoNm
	c.Data["PpChrgName"] = pPpChrgName
	c.Data["UpJobGrp"] = pUpJobGrp
	c.Data["JobGrp"] = pJobGrp
	c.Data["RecruitTitle"] = pRecruitTitle
	c.Data["InviteTitle"] = pInviteTitle
	c.Data["InviteMsg"] = pInviteMsg

	c.TplName = "invite/invite_send_msg_preview.html"
}
