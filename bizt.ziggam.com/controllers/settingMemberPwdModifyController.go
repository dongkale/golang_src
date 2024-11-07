package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type SettingMemberPwdModifyController struct {
	BaseController
}

func (c *SettingMemberPwdModifyController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/login")
	}
	pLang, _ := beego.AppConfig.String("lang")

	mem_sn := session.Get(c.Ctx.Request.Context(), "mem_sn")

	pEntpMemNo := mem_no //c.GetString("entp_mem_no")
	pPpChrgSn := mem_sn  //c.GetString("entp_mem_no")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Setting Member Modify
	fmt.Printf(fmt.Sprintf("CALL ZSP_TEAM_MEM_INFO_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pPpChrgSn))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_TEAM_MEM_INFO_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pPpChrgSn),
		ora.S, /* PP_CHRG_NM */
		ora.S, /* PP_CHRG_BP_NM */
		ora.S, /* PP_CHRG_TEL_NO */
		ora.S, /* SMS_RECV_YN */
		ora.S, /* EMAIL */
		ora.S, /* EMAIL_CERT_YN */
		ora.S, /* ENTP_MEM_ID */
	)
	defer stmtProcCall.Close()
	if err != nil {
		panic(err)
	}
	procRset := &ora.Rset{}
	_, err = stmtProcCall.Exe(procRset)

	if err != nil {
		panic(err)
	}

	settingMemberInfo := make([]models.SettingMemberInfo, 0)

	var (
		ppChrgNm    string
		ppChrgBpNm  string
		ppChrgTelNo string
		smsRecvYn   string
		email       string
		emailRecvYn string
		entpMemId   string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			ppChrgNm = procRset.Row[0].(string)
			ppChrgBpNm = procRset.Row[1].(string)
			ppChrgTelNo = procRset.Row[2].(string)
			smsRecvYn = procRset.Row[3].(string)
			email = procRset.Row[4].(string)
			emailRecvYn = procRset.Row[5].(string)
			entpMemId = procRset.Row[6].(string)

			settingMemberInfo = append(settingMemberInfo, models.SettingMemberInfo{
				PpChrgNm:    ppChrgNm,
				PpChrgBpNm:  ppChrgBpNm,
				PpChrgTelNo: ppChrgTelNo,
				SmsRecvYn:   smsRecvYn,
				Email:       email,
				EmailRecvYn: emailRecvYn,
				EntpMemId:   entpMemId,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Setting Member Modify
	c.Data["PpChrgSn"] = pPpChrgSn
	c.Data["EntpMemId"] = entpMemId

	c.Data["TMenuId"] = "M00"
	c.Data["SMenuId"] = "M00"

	c.TplName = "setting/password_modify.html"
}
