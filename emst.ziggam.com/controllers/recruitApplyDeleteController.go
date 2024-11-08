package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"

	"emst.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

type RecruitApplyDeleteController struct {
	BaseController
}

func (c *RecruitApplyDeleteController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// ￡ 파일로 남기게 처리
	// end : log

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/common/login")
	}

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := c.GetString("entp_mem_no")
	pRecrutSn := c.GetString("recrut_sn")
	pPpMemNo := c.GetString("pp_mem_no")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Recruit Apply Detail
	log.Debug("CALL SP_EMS_RECRUT_APPlY_DEL_PROC('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_RECRUT_APPlY_DEL_PROC('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo),
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

	recrutApplyDelete := models.RecrutApplyDelete{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		recrutApplyDelete = models.RecrutApplyDelete{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}
	// End : Favorite Apply Member Process

	c.Data["json"] = &recrutApplyDelete
	c.ServeJSON()
}
