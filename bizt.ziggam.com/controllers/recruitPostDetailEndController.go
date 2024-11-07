package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type RecruitPostDetailEndController struct {
	BaseController
}

func (c *RecruitPostDetailEndController) Post() {

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

	pEntpMemNo := mem_no //c.GetString("entp_mem_no")
	pRecrutSn := c.GetString("recrut_sn")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Recruit Modify
	fmt.Printf(fmt.Sprintf("CALL ZSP_RECRUT_APPLY_END_PROC('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_RECRUT_APPLY_END_PROC('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn),
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

	rtnRecruitStart := models.DefaultResult{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnRecruitStart = models.DefaultResult{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}

	// End : Recruit Insert Process

	c.Data["json"] = &rtnRecruitStart
	c.ServeJSON()
}
