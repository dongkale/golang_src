package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/rana/ora.v4"
)

// 공고 지원 결과 설정 API

type ApiSetApplyResultController struct {
	beego.Controller
}

func (c *ApiSetApplyResultController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/login")
	}

	pLang, err := beego.AppConfig.String("lang")
	if err != nil {
		panic(err)
	}

	pEntpMemNo := mem_no
	pRecrutSn := c.GetString("recrut_sn")
	pPpMemNo := c.GetString("pp_mem_no")
	pApplyResult := c.GetString("apply_result") // 03 : 합격, 04 : 불합격
	//imgServer, _  := beego.AppConfig.String("viewpath")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Applicant Delete Process

	log.Debug(fmt.Sprintf("CALL ZSP_APPLY_RESULT_PROC('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo, pApplyResult))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_APPLY_RESULT_PROC('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo, pApplyResult),
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

	rtnApplyResult := models.RtnApplyResult{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnApplyResult = models.RtnApplyResult{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}
	// End : Applicant Delete Process

	c.Data["json"] = &rtnApplyResult
	c.ServeJSON()

}
