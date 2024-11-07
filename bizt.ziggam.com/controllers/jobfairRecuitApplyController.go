package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	ora "gopkg.in/rana/ora.v4"
)

type JobFairRecruitApplyController struct {
	beego.Controller
}

func (c *JobFairRecruitApplyController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	fmt.Printf("JobFairApplyController")

	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := c.GetString("entp_mem_no")
	pRecrutSn := c.GetString("recrut_sn")
	pPpMemNo := c.GetString("pp_mem_no")

	fmt.Printf(pLang)
	fmt.Printf(pEntpMemNo)
	fmt.Printf(pRecrutSn)
	fmt.Printf(pPpMemNo)

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	var (
		rtnCd  int64
		rtnMsg string
	)

	// 채용 공고에 지원 시작 -->
	fmt.Printf(fmt.Sprintf("CALL MSP_APPLY_START_PROC('%v', '%v', '%v', '%v', :1)", pLang, pEntpMemNo, pRecrutSn, pPpMemNo))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL MSP_APPLY_START_PROC('%v', '%v', '%v', '%v', :1)",
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

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// <--

	if rtnCd != 1 {
		logs.Error("CALL MSP_APPLY_START_PROC('%v', '%v', '%v', '%v', :1) ==> rtnCd != 1 : %v", pLang, pEntpMemNo, pRecrutSn, pPpMemNo)
		panic(fmt.Sprintf("CALL MSP_APPLY_START_PROC('%v', '%v', '%v', '%v', :1) Error!!", pLang, pEntpMemNo, pRecrutSn, pPpMemNo))
	}

	// 채용 공고에 라이브 지원 -->
	fmt.Printf(fmt.Sprintf("CALL MSP_APPLY_LIV_PROC('%v', '%v', '%v', '%v', :1)", pLang, pEntpMemNo, pRecrutSn, pPpMemNo))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL MSP_APPLY_LIV_PROC('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
	)
	defer stmtProcCall.Close()
	if err != nil {
		panic(err)
	}
	procRset = &ora.Rset{}
	_, err = stmtProcCall.Exe(procRset)

	if err != nil {
		panic(err)
	}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// <--

	if rtnCd != 1 {
		logs.Error("CALL MSP_APPLY_LIV_PROC('%v', '%v', '%v', '%v', :1) ==> rtnCd != 1 : %v", pLang, pEntpMemNo, pRecrutSn, pPpMemNo)
		panic(fmt.Sprintf("CALL MSP_APPLY_LIV_PROC('%v', '%v', '%v', '%v', :1) Error!!", pLang, pEntpMemNo, pRecrutSn, pPpMemNo))
	}

	rtnDefault := models.DefaultResult{}

	rtnDefault = models.DefaultResult{
		RtnCd:  rtnCd,
		RtnMsg: rtnMsg,
	}

	c.Data["json"] = &rtnDefault
	c.ServeJSON()
}
