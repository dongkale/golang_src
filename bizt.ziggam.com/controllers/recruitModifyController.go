package controllers

import (
	"fmt"
	"strings"

	"bizt.ziggam.com/models"
	"bizt.ziggam.com/tables"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type RecruitModifyController struct {
	BaseController
}

func (c *RecruitModifyController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	//log.SetLogger(logs.AdapterFile, `{"filename":"logs/biz_f.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
	// end : log

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/login")
	}

	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := mem_no //c.GetString("entp_mem_no")
	pRecrutSn := c.GetString("recrut_sn")

	pMode := c.GetString("mode")
	pStep := c.GetString("step")

	// LDK 2020/11/12 : 추가(채용 공고 기간 수정) -->
	pPeriod := c.GetString("period")
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

	// Start : Recruit Modify
	// LDK 2020/08/26 : 채용 정보 코드화, 추가 -->
	fmt.Printf(fmt.Sprintf("CALL ZSP_RECRUIT_INFO_R_V2('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_RECRUIT_INFO_R_V2('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn),
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* RECRUT_TITLE */
		ora.S,   /* UP_JOB_GRP */
		ora.S,   /* JOB_GRP */
		ora.S,   /* RECRUT_GBN_CD */
		ora.I64, /* RECRUT_CNT */
		ora.S,   /* ROL */
		ora.S,   /* APLY_QUFCT */
		ora.S,   /* PERFER_TRTM */
		ora.S,   /* SDY */
		ora.S,   /* EDY */
		ora.S,   /* VD_TITLE_UPT_YN */
		ora.S,   /* PRGS_STAT_CD */
		ora.S,   /* REG_DT */
		ora.S,   /* PP_CHRG_BP_NM */
		ora.S,   /* PP_CHRG_NM */
		ora.S,   /* UP_JOB_GRP_CD */
		ora.S,   /* JOB_GRP_CD */
		ora.S,   /* DCMNT_EVL_USE_CD */
		ora.S,   /* ONWY_INTRV_USE_CD */
		ora.S,   /* LIVE_INTRV_USE_CD */
		ora.S,   /* RECRUT_PROC_CD */
		ora.S,   /* JF_MNG_CD */
		ora.S,   /* JF_TITLE */
		ora.S,   /* CARR_GBN_CD */    // LDK 2020/08/26 : 추가
		ora.S,   /* ENTP_ADDR */      // LDK 2020/08/26 : 추가
		ora.S,   /* EMPL_TYP_CD */    // LDK 2020/08/26 : 추가
		ora.S,   /* LST_EDU_GBN_CD */ // LDK 2020/08/26 : 추가
		ora.S,   /* PRGS_STAT_STEP */ // LDK 2020/08/26 : 추가
		ora.S,   /* ANNUAL_SALARY */  // LDK 2020/08/26 : 추가
		ora.S,   /* WORK_DAYS */      // LDK 2020/08/26 : 추가
		ora.S,   /* WELFARE */        // LDK 2020/08/26 : 추가
	)
	// <--

	defer stmtProcCall.Close()
	if err != nil {
		panic(err)
	}
	procRset := &ora.Rset{}
	_, err = stmtProcCall.Exe(procRset)

	if err != nil {
		panic(err)
	}

	recruitModify := make([]models.RecruitModify, 0)

	var (
		entpMemNo      string
		recrutSn       string
		recrutTitle    string
		upJobGrp       string
		jobGrp         string
		recrutGbnCd    string
		recrutCnt      int64
		rol            string
		aplyQufct      string
		perferTrtm     string
		sdy            string
		edy            string
		vdTitleUptYn   string
		prgsStatCd     string // LDK 2020/11/12 : 추가(채용 공고 기간 수정) <-->
		regDT          string // LDK 2020/11/12 : 추가(채용 공고 기간 수정) <-->
		dcmntEvlUseCd  string
		onwyIntrvUseCd string
		liveIntrvUseCd string
		jobfairMngCd   string // LDK 2020/08/26 : 추가
		carrGbnCd      string // LDK 2020/08/26 : 추가
		entpAddr       string // LDK 2020/08/26 : 추가
		emplTypeCd     string // LDK 2020/08/26 : 추가
		lstEduGbnCd    string // LDK 2020/08/26 : 추가
		prgsStatStep   string // LDK 2020/08/26 : 추가
		annualSalary   string // LDK 2020/08/26 : 추가
		workDays       string // LDK 2020/08/26 : 추가
		welfare        string // LDK 2020/08/26 : 추가
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			entpMemNo = procRset.Row[0].(string)
			recrutSn = procRset.Row[1].(string)
			recrutTitle = procRset.Row[2].(string)
			upJobGrp = procRset.Row[3].(string)
			jobGrp = procRset.Row[4].(string)
			recrutGbnCd = procRset.Row[5].(string)
			recrutCnt = procRset.Row[6].(int64)
			rol = procRset.Row[7].(string)
			aplyQufct = procRset.Row[8].(string)
			perferTrtm = procRset.Row[9].(string)
			sdy = procRset.Row[10].(string)
			edy = procRset.Row[11].(string)
			vdTitleUptYn = procRset.Row[12].(string)

			prgsStatCd = procRset.Row[13].(string) // LDK 2020/11/12 : 추가(채용 공고 기간 수정) <-->
			regDT = procRset.Row[14].(string)      // LDK 2020/11/12 : 추가(채용 공고 기간 수정) <-->

			dcmntEvlUseCd = procRset.Row[19].(string)
			onwyIntrvUseCd = procRset.Row[20].(string)
			liveIntrvUseCd = procRset.Row[21].(string)

			// LDK 2020/08/26 : 채용 정보 코드화, 추가 -->
			jobfairMngCd = procRset.Row[23].(string)
			/* jobfairTitle = procRset.Row[24].(string) */
			carrGbnCd = procRset.Row[25].(string)
			entpAddr = procRset.Row[26].(string)
			emplTypeCd = procRset.Row[27].(string)
			lstEduGbnCd = procRset.Row[28].(string)
			prgsStatStep = procRset.Row[29].(string)
			annualSalary = procRset.Row[30].(string)
			workDays = procRset.Row[31].(string)
			welfare = procRset.Row[32].(string)
			// <--

			recruitModify = append(recruitModify, models.RecruitModify{
				EntpMemNo:      entpMemNo,
				RecrutSn:       recrutSn,
				RecrutTitle:    recrutTitle,
				UpJobGrp:       upJobGrp,
				JobGrp:         jobGrp,
				RecrutGbnCd:    recrutGbnCd,
				RecrutCnt:      recrutCnt,
				Rol:            rol,
				AplyQufct:      aplyQufct,
				PerferTrtm:     perferTrtm,
				Sdy:            sdy,
				Edy:            edy,
				VdTitleUptYn:   vdTitleUptYn,
				DcmntEvlUseCd:  dcmntEvlUseCd,
				OnwyIntrvUseCd: onwyIntrvUseCd,
				LiveIntrvUseCd: liveIntrvUseCd,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Recruit Modify
	
	// Start : Recruit Question List
	// 콘솔 출력	
	fmt.Printf(fmt.Sprintf("CALL ZSP_RECRUIT_QST_LIST_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_RECRUIT_QST_LIST_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn),
		ora.I64, /* TOT_CNT */
		ora.S,   /* QST_SN */
		ora.S,   /* VD_TITLE */

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

	recruitQuestionList := make([]models.RecruitQuestionList, 0)

	var (
		qstTotCnt int64
		qstSn     string
		vdTitle   string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			qstTotCnt = procRset.Row[0].(int64)
			qstSn = procRset.Row[1].(string)
			vdTitle = procRset.Row[2].(string)

			recruitQuestionList = append(recruitQuestionList, models.RecruitQuestionList{
				QstTotCnt: qstTotCnt,
				QstSn:     qstSn,
				VdTitle:   vdTitle,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Recruit Question List

	// Start : Region Group List

	// LDK 2020/08/26 : 채용 정보 코드화, 추가 -->
	fmt.Printf(fmt.Sprintf("CALL ZSP_RGN_LIST_R_V2('%v', NULL, :1)", pLang))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_RGN_LIST_R_V2('%v', NULL, :1)", pLang),
		ora.S, /* RGN_GRP_CD */
		ora.S, /* RGN_GRP_NM */
	)
	// <--

	defer stmtProcCall.Close()
	if err != nil {
		panic(err)
	}
	procRset = &ora.Rset{}
	_, err = stmtProcCall.Exe(procRset)

	if err != nil {
		panic(err)
	}

	rgnGrpList := make([]models.RgnGrp, 0)

	var (
		rgnGrpCd  string
		rgnGrpNm  string
		rgnGrpFNm string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			rgnGrpCd = procRset.Row[0].(string)
			rgnGrpNm = procRset.Row[1].(string)
			rgnGrpFNm = procRset.Row[2].(string)

			rgnGrpList = append(rgnGrpList, models.RgnGrp{
				Code:     rgnGrpCd,
				Name:     rgnGrpNm,
				FullName: rgnGrpFNm,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Job Group List

	// Start : Jobfair List
	fmt.Printf(fmt.Sprintf("CALL ZSP_ENTP_JOBFAIR_LIST('%v', '%v', :1)", pLang, pEntpMemNo))
	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_ENTP_JOBFAIR_LIST('%v', '%v', :1)",
		pLang, pEntpMemNo),
		ora.S, /* MNG_CD */
		ora.S, /* TITLE */
		ora.S, /* SDY */
		ora.S, /* EDY */
		ora.S, /* HOST_INSTITUTION */
		ora.S, /* MANAGE_AGENCY */
		ora.S, /* AGREEMENT_URL */
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

	jobFailrInfoList := make([]models.JobfairInfo, 0)

	var (
		jfMngCd           string
		jfTitle           string
		jfSdy             string
		jfEdy             string
		jfHostInstitution string
		jfManageAgency    string
	)

	if procRset.IsOpen() {
		for procRset.Next() {

			jfMngCd = procRset.Row[0].(string)
			jfTitle = procRset.Row[1].(string)
			jfSdy = procRset.Row[2].(string)
			jfEdy = procRset.Row[3].(string)
			jfHostInstitution = procRset.Row[4].(string)
			jfManageAgency = procRset.Row[5].(string)

			jobFailrInfoList = append(jobFailrInfoList, models.JobfairInfo{
				MngCd:           jfMngCd,
				Title:           jfTitle,
				Sdy:             jfSdy,
				Edy:             jfEdy,
				HostInstitution: jfHostInstitution,
				ManageAgency:    jfManageAgency,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Jobfair List

	c.Data["QstTotCnt"] = qstTotCnt
	c.Data["RecruitQuestionList"] = recruitQuestionList

	c.Data["EntpMemNo"] = entpMemNo
	c.Data["RecrutSn"] = recrutSn
	c.Data["RecrutTitle"] = recrutTitle
	c.Data["UpJobGrp"] = upJobGrp
	c.Data["JobGrp"] = jobGrp
	c.Data["RecrutGbnCd"] = recrutGbnCd
	c.Data["RecrutCnt"] = recrutCnt
	c.Data["Rol"] = strings.Replace(rol, "<br>", "\n", -1)
	c.Data["AplyQufct"] = strings.Replace(aplyQufct, "<br>", "\n", -1)
	c.Data["PerferTrtm"] = strings.Replace(perferTrtm, "<br>", "\n", -1)
	c.Data["Sdy"] = sdy
	c.Data["Edy"] = edy
	c.Data["VdTitleUptYn"] = vdTitleUptYn
	c.Data["DcmntEvlUseCd"] = dcmntEvlUseCd
	c.Data["OnwyIntrvUseCd"] = onwyIntrvUseCd
	c.Data["LiveIntrvUseCd"] = liveIntrvUseCd

	// LDK 2020/08/26 : 채용 정보 코드화, 추가 -->
	c.Data["JobfairMngCd"] = jobfairMngCd
	c.Data["CarrGbnCd"] = carrGbnCd
	c.Data["EntpAddr"] = entpAddr
	c.Data["EmplTypeCd"] = emplTypeCd
	c.Data["LstEduGbnCd"] = lstEduGbnCd
	c.Data["PrgsStatStep"] = prgsStatStep
	c.Data["AnnualSalary"] = annualSalary
	c.Data["WorkDays"] = workDays
	c.Data["Welfare"] = welfare

	c.Data["RgnGrpList"] = rgnGrpList
	c.Data["JobFairList"] = jobFailrInfoList

	c.Data["CarrGbnMap"] = tables.MapCarrGbnCd
	c.Data["EmplTypMap"] = tables.MapEmplTypCd
	c.Data["LstEduGbnMap"] = tables.MapLstEduGbnCd
	// <--

	// LDK 2020/11/12 : 추가(채용 공고 기간 수정) <-->
	c.Data["PrgsStatCd"] = prgsStatCd // 진행 상황(01: 대기 / 02: 진행중 / 03: 마감);
	c.Data["RegDT"] = regDT
	c.Data["Period"] = pPeriod
	// <--
	c.Data["Mode"] = pMode
	c.Data["Step"] = pStep

	c.Data["TMenuId"] = "R00"
	c.Data["SMenuId"] = "R00"

	c.TplName = "recruit/recruit_modify.html"
}
