package controllers

import (
	"fmt"

	"emst.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	"gopkg.in/rana/ora.v4"
)

type RecruitPostUpdateController struct {
	beego.Controller
}

func (c *RecruitPostUpdateController) Post() {

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
	pEntpMemNo := mem_no
	pRecrutSn := c.GetString("recrut_sn")
	pRecrutCnt := c.GetString("recrut_cnt")
	pSex := c.GetString("sex")
	pRol := c.GetString("rol")
	pAplyQufct := c.GetString("aply_qufct")
	pRecrutTitle := c.GetString("recrut_title")
	pPrgsMsg := c.GetString("prgs_msg")
	pArrQstTitle := c.GetString("qst_title_arr")
	pUptYn := c.GetString("upt_yn")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Recruit Update Process

	log.Debug("CALL SP_EMS_RECRUIT_UPT_PROC('%v', '%v', '%v',  %v , '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pRecrutCnt, pSex, pRol, pAplyQufct, pRecrutTitle, pPrgsMsg, pArrQstTitle, pUptYn)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_RECRUIT_UPT_PROC('%v', '%v', '%v',  %v , '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pRecrutCnt, pSex, pRol, pAplyQufct, pRecrutTitle, pPrgsMsg, pArrQstTitle, pUptYn),
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

	rtnRecruitPostUpdate := models.RtnRecruitPostUpdate{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnRecruitPostUpdate = models.RtnRecruitPostUpdate{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}

	// End : Recruit Update Process

	c.Data["json"] = &rtnRecruitPostUpdate
	c.ServeJSON()
}
