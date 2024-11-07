package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type RecruitUpdateController struct {
	beego.Controller
}

func (c *RecruitUpdateController) Post() {

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
	pVdTitleUptYn := c.GetString("vd_title_upt_yn")
	pEntpMemNo := mem_no
	pRecrutSn := c.GetString("recrut_sn")
	pPpChrgSn := c.GetString("pp_chrg_sn")
	pRecrutGbnCd := c.GetString("recrut_gbn_cd")
	pRecrutCnt := c.GetString("recrut_cnt")
	pRol := c.GetString("rol")
	pAplyQufct := c.GetString("aply_qufct")
	pPerferTrtm := c.GetString("perfer_trtm")
	pRecrutTitle := c.GetString("recrut_title")
	pArrQstTitle := c.GetString("qst_title_arr")

	pDcmntEvlUseCd := c.GetString("dcmnt_evl_use_cd")
	pOnwyIntrvUseCd := c.GetString("onwy_intrv_use_cd")
	pLiveIntrvUseCd := c.GetString("live_intrv_use_cd")

	fmt.Printf("pArrQstTitle: " + pArrQstTitle)

	// LDK 2020/08/24 채용 정보 코드화 -->
	pCarrGbnCd := c.GetString("carr_gbn_cd")
	pEntpAddr := c.GetString("entp_addr")
	pEmplTypCd := c.GetString("empl_typ_cd")
	pLstEduGbnCd := c.GetString("lst_edu_gbn_cd")
	pPrgsStatStep := c.GetString("prgs_stat_step")
	pAnnualSalary := c.GetString("annual_salary")
	pWorkDays := c.GetString("work_days")
	pWelfare := c.GetString("welfare")
	pJobfair := c.GetString("jobfair")

	fmt.Printf("pCarrGbnCd: " + pCarrGbnCd)
	fmt.Printf("pEntpAddr: " + pEntpAddr)
	fmt.Printf("pEmplTypCd: " + pEmplTypCd)
	fmt.Printf("pLstEduGbnCd: " + pLstEduGbnCd)
	fmt.Printf("pPrgsStatStep: " + pPrgsStatStep)
	fmt.Printf("pAnnualSalary: " + pAnnualSalary)
	fmt.Printf("pWorkDays: " + pWorkDays)
	fmt.Printf("pWelfare: " + pWelfare)
	fmt.Printf("pJobfair: " + pJobfair)
	// <--

	// LDK 2020/11/09: 채용 공고 수정에 날짜 추가 -->
	pSdy := c.GetString("sdy")
	pEdy := c.GetString("edy")

	fmt.Printf("pSdy: " + pSdy)
	fmt.Printf("pEdy: " + pEdy)

	// convTime1, _ := time.Parse("2006/01/02 15:04", pSdy)
	// convTime2, _ := time.Parse("2006/01/02 15:04", pEdy)

	// fmt.Printf("pSdy Date: " + convTime1.Format("20060102"))
	// fmt.Printf("pEdy Date: " + convTime2.Format("20060102"))

	// fmt.Printf("pSdy Time: " + convTime1.Format("1504"))
	// fmt.Printf("pEdy Time: " + convTime2.Format("1504"))
	// <--

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
	fmt.Printf(fmt.Sprintf("CALL ZSP_RECRUIT_UPT_PROC_V3('%v', '%v', '%v', '%v', '%v', '%v',  %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pVdTitleUptYn, pEntpMemNo, pRecrutSn, pPpChrgSn, pRecrutGbnCd, pRecrutCnt, pRol, pAplyQufct, pPerferTrtm, pRecrutTitle, pSdy, pEdy, pArrQstTitle, pDcmntEvlUseCd, pOnwyIntrvUseCd, pLiveIntrvUseCd, pJobfair, pCarrGbnCd, pEntpAddr, pEmplTypCd, pLstEduGbnCd, pPrgsStatStep, pAnnualSalary, pWorkDays, pWelfare))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_RECRUIT_UPT_PROC_V3('%v', '%v', '%v', '%v', '%v', '%v', %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pVdTitleUptYn, pEntpMemNo, pRecrutSn, pPpChrgSn, pRecrutGbnCd, pRecrutCnt, pRol, pAplyQufct, pPerferTrtm, pRecrutTitle, pSdy, pEdy, pArrQstTitle, pDcmntEvlUseCd, pOnwyIntrvUseCd, pLiveIntrvUseCd, pJobfair, pCarrGbnCd, pEntpAddr, pEmplTypCd, pLstEduGbnCd, pPrgsStatStep, pAnnualSalary, pWorkDays, pWelfare),
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

	rtnRecruitUpdate := models.RtnRecruitUpdate{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnRecruitUpdate = models.RtnRecruitUpdate{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}
	// End : Recruit Update Process

	c.Data["json"] = &rtnRecruitUpdate
	c.ServeJSON()
}
