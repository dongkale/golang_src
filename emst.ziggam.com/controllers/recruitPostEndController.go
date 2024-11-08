package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"

	"emst.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

type RecruitPostEndController struct {
	beego.Controller
}

func (c *RecruitPostEndController) Post() {

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
	pEntpMemNo := c.GetString("entp_mem_no") // 기업회원번호
	pRecrutSn := c.GetString("recrut_sn")    // 문의종류코드

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Recruit Stat End
	log.Debug("CALL SP_EMS_RECRUIT_END_PROC('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_RECRUIT_END_PROC('%v', '%v', '%v', :1)",
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

	rtnRecruitPostEnd := models.RtnRecruitPostEnd{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnRecruitPostEnd = models.RtnRecruitPostEnd{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}

	// End : Recruit Stat End

	c.Data["json"] = &rtnRecruitPostEnd
	c.ServeJSON()
}
