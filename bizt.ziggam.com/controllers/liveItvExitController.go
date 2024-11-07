package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	ora "gopkg.in/rana/ora.v4"
)

type LiveItvExitController struct {
	BaseController
}

func (c *LiveItvExitController) Post() {

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

	pEntpMemNo := mem_no
	//pEntpMemSn := c.GetString("entp_mem_sn")
	pRecrutSn := c.GetString("recrut_sn")
	pPpMemNo := c.GetString("pp_mem_no")
	pPpMemNm := c.GetString("pp_mem_nm")
	pLiveSn := c.GetString("live_sn")

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
	fmt.Printf(fmt.Sprintf("CALL MSP_LIVE_OUT_PROC('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo, pLiveSn))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL MSP_LIVE_OUT_PROC('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo, pLiveSn),
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

	//rtnLiveInvExit := models.DefaultResult{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		// rtnLiveInvExit = models.DefaultResult{
		// 	RtnCd:  rtnCd,
		// 	RtnMsg: rtnMsg,
		// }
	}

	if rtnCd != 1 {
		logs.Error("CALL MSP_LIVE_OUT_PROC('%v', '%v', '%v', '%v', '%v', '%v', :1) ==> rtnCd != 1 : %v",
			pLang, "E", pEntpMemNo, pRecrutSn, pPpMemNo, pLiveSn, rtnMsg)
	}

	// 결과
	fmt.Printf(fmt.Sprintf("CALL MSP_LIVE_END_RSLT_R('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo, pLiveSn))

	// EntpMemNo : 기업회원번호
	// RecrutSn : 채용일련번호
	// PpMemNo : 개인회원번호
	// LiveSn : 라이브인터뷰 일련번호
	// EntpNm : 기업명
	// LiveItvSday : 라이브 인터뷰 시작 일자
	// LiveItvStime : 라이브 인터뷰 시작 시간(초)
	// LiveItvEday : 라이브 인터뷰 종료 일자
	// LiveItvEtime : 라이브 인터뷰 종료 시간(초)
	// LiveItvJt : 총 인터뷰 시간
	// Nm : 인터뷰 대상자 명 (개인회원명)
	// SdtTstmp : 라이브 인터뷰 TimeStamp

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL MSP_LIVE_END_RSLT_R('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo, pLiveSn),
		ora.I64, /* RTN_CD */
		ora.S,   /* ENTP_MEM_NO */ // 기업회원번호
		ora.S,   /* RECRUT_SN */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* LIVE_SN */
		ora.S,   /* ENTP_NM */
		ora.S,   /* NM */
		ora.S,   /* LIVE_ITV_SDAY */
		ora.S,   /* LIVE_ITV_STIME */
		ora.S,   /* LIVE_ITV_EDAY */
		ora.S,   /* LIVE_ITV_ETIME */
		ora.S,   /* LIVE_ITV_JT */
		ora.S,   /* SDT_TSTMP */
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

	var (
		rtnEntpMemNo    string
		rtnRecrutSn     string
		rtnPpMemNo      string
		rtnLiveSn       string
		rtnEntpNm       string
		rtnNm           string
		rtnLiveItvSday  string
		rtnLiveItvSTime string
		rtnLiveItvEday  string
		rtnLiveItvETime string
		rtnLiveItvJt    string
		rtnSdtTstmp     string
	)

	rtnLiveInvResult := models.LiveInvResult{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnEntpMemNo = procRset.Row[1].(string)
			rtnRecrutSn = procRset.Row[2].(string)
			rtnPpMemNo = procRset.Row[3].(string)
			rtnLiveSn = procRset.Row[4].(string)
			rtnEntpNm = procRset.Row[5].(string)
			rtnNm = procRset.Row[6].(string)
			rtnLiveItvSday = procRset.Row[7].(string)
			rtnLiveItvSTime = procRset.Row[8].(string)
			rtnLiveItvEday = procRset.Row[9].(string)
			rtnLiveItvETime = procRset.Row[10].(string)
			rtnLiveItvJt = procRset.Row[11].(string)
			rtnSdtTstmp = procRset.Row[12].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		if rtnCd == 1 {
			rtnLiveInvResult = models.LiveInvResult{
				LirEntpMemNo:    rtnEntpMemNo,
				LirRecrutSn:     rtnRecrutSn,
				LirPpMemNo:      rtnPpMemNo,
				LirPpMemNm:      pPpMemNm,
				LirLiveSn:       rtnLiveSn,
				LirEntpNm:       rtnEntpNm,
				LirNm:           rtnNm,
				LirLiveItvSday:  rtnLiveItvSday,
				LirLiveItvSTime: rtnLiveItvSTime,
				LirLiveItvEday:  rtnLiveItvEday,
				LirLiveItvETime: rtnLiveItvETime,
				LirLiveItvJt:    rtnLiveItvJt,
				LirSdtTstmp:     rtnSdtTstmp,
			}
		}
	}

	c.Data["json"] = &rtnLiveInvResult
	c.ServeJSON()
}
