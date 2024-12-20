package controllers

import (
	"fmt"

	"emst.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	"gopkg.in/rana/ora.v4"
)

type AdminInquiryInsertController struct {
	beego.Controller
}

func (c *AdminInquiryInsertController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/common/login")
	}

	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := c.GetString("entp_mem_no")
	pInqGbnCd := c.GetString("inq_gbn_cd")
	pRegDy := c.GetString("reg_dy")
	pAnsCont := c.GetString("ans_cont")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Inquiry Answer Process

	log.Debug("CALL SP_EMS_INQUIRY_ANS_PROC('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pInqGbnCd, pRegDy, pAnsCont)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_INQUIRY_ANS_PROC('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pInqGbnCd, pRegDy, pAnsCont),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
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

	var (
		rtnCd  int64
		rtnMsg string
	)

	rtnAdminInquiryInsert := models.RtnAdminInquiryInsert{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnAdminInquiryInsert = models.RtnAdminInquiryInsert{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}

	// End : Inquiry Answer Process

	c.Data["json"] = &rtnAdminInquiryInsert
	c.ServeJSON()
}
