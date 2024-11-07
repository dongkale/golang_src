package controllers

import (
	"fmt"
	"math"
	"strconv"

	"bizt.ziggam.com/models"
	"bizt.ziggam.com/utils"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type RecruitPostDetailController struct {
	BaseController
}

func (c *RecruitPostDetailController) Get() {

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
	imgServer, _  := beego.AppConfig.String("viewpath")

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
		ora.S,   /* PP_CHRG_NM */
		ora.S,   /* UP_JOB_GRP_CD */
		ora.S,   /* JOB_GRP_CD */
		ora.S,   /* DCMNT_EVL_USE_CD */
		ora.S,   /* ONWY_INTRV_USE_CD */
		ora.S,   /* LIVE_INTRV_USE_CD */
		ora.S,   /* RECRUT_PROC_CD */
		ora.S,   /* JF_MNG_CD */
		ora.S,   /* JF_TITLE */
		ora.S,   /* CARR_GBN_CD */
		ora.S,   /* ENTP_ADDR */
		ora.S,   /* EMPL_TYP_CD */
		ora.S,   /* LST_EDU_GBN_CD */
		ora.S,   /* PRGS_STAT_STEP */
		ora.S,   /* ANNUAL_SALARY */
		ora.S,   /* WORK_DAYS */
		ora.S,   /* WELFARE */
		ora.S,   /* RECRUT_EDT */ // LDK 2020/11/25: 즉시 마감 날짜 출력
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

	recruitModify := make([]models.RecruitModify, 0)

	var (
		entpMemNo    string
		recrutSn     string
		recrutTitle  string
		upJobGrp     string
		jobGrp       string
		recrutGbnCd  string
		recrutCnt    int64
		rol          string
		aplyQufct    string
		perferTrtm   string
		sdy          string
		edy          string
		vdTitleUptYn string
		prgsStatCd   string
		regDt        string
		ppChrgBpNm   string
		ppChrgNm     string
		recrutProdCd string
		recrutEdt    string // LDK 2020/11/25: 즉시 마감 날짜 출력
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
			prgsStatCd = procRset.Row[13].(string)
			regDt = procRset.Row[14].(string)
			ppChrgBpNm = procRset.Row[15].(string)
			ppChrgNm = procRset.Row[16].(string)
			recrutProdCd = procRset.Row[22].(string)

			recrutEdt = procRset.Row[33].(string) // LDK 2020/11/25: 즉시 마감 날짜 출력

			recruitModify = append(recruitModify, models.RecruitModify{
				EntpMemNo:    entpMemNo,
				RecrutSn:     recrutSn,
				RecrutTitle:  recrutTitle,
				UpJobGrp:     upJobGrp,
				JobGrp:       jobGrp,
				RecrutGbnCd:  recrutGbnCd,
				RecrutCnt:    recrutCnt,
				Rol:          rol,
				AplyQufct:    aplyQufct,
				PerferTrtm:   perferTrtm,
				Sdy:          sdy,
				Edy:          edy,
				VdTitleUptYn: vdTitleUptYn,
				PrgsStatCd:   prgsStatCd,
				RegDt:        regDt,
				PpChrgBpNm:   ppChrgBpNm,
				PpChrgNm:     ppChrgNm,
				RecrutProdCd: recrutProdCd,
				RecrutEdt:    recrutEdt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Recruit Modify

	// Start : Recruit Question List
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

	// Start : Recruit Top Stat Info
	fmt.Printf(fmt.Sprintf("CALL ZSP_APPLY_TOP_INFO_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_APPLY_TOP_INFO_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn),
		ora.I64, /* APPLY_CNT */
		ora.I64, /* ING_CNT */
		ora.I64, /* PASS_CNT */
		ora.I64, /* FAIL_CNT */
		ora.I64, /* DCMNT_PASS_CNT */
		ora.I64, /* DCMNT_FAIL_CNT */
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

	recruitStatTopInfo := make([]models.RecruitStatTopInfo, 0)

	var (
		applyCnt     int64
		ingCnt       int64
		passCnt      int64
		failCnt      int64
		dcmntPassCnt int64
		dcmntFailCnt int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			applyCnt = procRset.Row[0].(int64)
			ingCnt = procRset.Row[1].(int64)
			passCnt = procRset.Row[2].(int64)
			failCnt = procRset.Row[3].(int64)
			dcmntPassCnt = procRset.Row[4].(int64)
			dcmntFailCnt = procRset.Row[5].(int64)

			recruitStatTopInfo = append(recruitStatTopInfo, models.RecruitStatTopInfo{
				ApplyCnt:     applyCnt,
				IngCnt:       ingCnt,
				PassCnt:      passCnt,
				FailCnt:      failCnt,
				DcmntPassCnt: dcmntPassCnt,
				DcmntFailCnt: dcmntFailCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Recruit Top Stat Info

	// Start : Recruit Apply List
	pEvlPrgsStat := c.GetString("evl_prgs_stat")      // 평가진행상태 (00:전체, 02:대기, 03:합격, 04:불합격)
	pSex := c.GetString("sex")                        //성별 (무관:A, 남성:M, 여성:F)
	pAge := c.GetString("age")                        //연령 (전체:00, 19세이하:01, 20~29:02, 30~39:03, 40~49:04, 50~59:05, 60세이상:06)
	pVpYn := c.GetString("vp_yn")                     //영상프로필 (전체:9, 있음:1, 없음:0)
	pFavrAplyPp := c.GetString("favr_aply_pp")        //관심지원자 (전체:9, 있음:1, 없음:0)
	pSortGbn := c.GetString("sort_gbn")               // 정렬구분(01:신규순, 02:마감순)
	pKeyword := c.GetString("keyword")                // 검색어
	pLiveReqStatCd := c.GetString("live_req_stat_cd") // 라이브요청상태코드(전체:A, 03:예정)

	pApplySortCd := c.GetString("apply_sort_cd")
	pApplySortWay := c.GetString("apply_sort_way")

	if pApplySortCd == "" {
		pApplySortCd = "01"
		// LDK 2021/01/29 : Z_ORDER 정렬 설정은 pApplySortCd = "00"
	}
	if pApplySortWay == "" {
		pApplySortWay = "DESC"
		// LDK 2021/01/29 : Z_ORDER 정렬 설정은 pApplySortWay = "ASC"
	}

	/* Parameter */
	/*
		pmKeyword := c.GetString("p_keyword")     // 검색어
		pmEmplTyp := c.GetString("p_empl_typ")    // 고용형태코드
		pmJobGrpCd := c.GetString("p_job_grp_cd") // 직군코드
		pmSortGbn := c.GetString("p_sort_gbn")    // 정렬구분
		pmGbnCd := c.GetString("p_gbn_cd")        // 구분코드
		pmPageNo := c.GetString("p_page_no")      // 페이지번호
	*/

	if pEvlPrgsStat == "" {
		pEvlPrgsStat = "00"
	}

	if pSortGbn == "" {
		pSortGbn = "01"
	}

	if pSex == "" {
		pSex = "A"
	}

	if pAge == "" {
		pAge = "00"
	}

	if pVpYn == "" {
		pVpYn = "9"
	}

	if pFavrAplyPp == "" {
		pFavrAplyPp = "9"
	}

	if pLiveReqStatCd == "" {
		pLiveReqStatCd = "A"
	}

	// Get Parameter
	var (
		pageNo    int64
		pageSize  int64
		finalPage int64
		pageList  int64
	)

	pPageNo := c.GetString("pn")
	if pPageNo == "" {
		pPageNo = "1"
	}
	pageNo, err = strconv.ParseInt(pPageNo, 10, 64)
	if err != nil {
		//
	}

	pPageSize := c.GetString("size")
	if pPageSize == "" {
		pPageSize = "30"
	}

	//pPageSize = strconv.FormatInt(imPageNo, 16)

	pageSize, err = strconv.ParseInt(pPageSize, 10, 64)
	if err != nil {
		//
	}
	pOffSet := (pageNo - 1) * pageSize
	pLimit := pageSize
	pageList = 10

	// Start : Recruit Stat List
	fmt.Printf(fmt.Sprintf("CALL ZSP_APPLY_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v' :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pRecrutSn, pEvlPrgsStat, pSex, pAge, pVpYn, pFavrAplyPp, pSortGbn, pKeyword, pLiveReqStatCd, pApplySortCd, pApplySortWay))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_APPLY_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pRecrutSn, pEvlPrgsStat, pSex, pAge, pVpYn, pFavrAplyPp, pSortGbn, pKeyword, pLiveReqStatCd, pApplySortCd, pApplySortWay),
		ora.I64, /* TOT_CNT */
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* FAVR_APLY_PP_YN */
		ora.S,   /* NM */
		ora.S,   /* SEX */
		ora.S,   /* AGE */
		ora.S,   /* REG_DT */
		ora.S,   /* APPLY_DT */
		ora.S,   /* EVL_STAT_DT */
		ora.S,   /* EVL_PRGS_STAT_CD */
		ora.S,   /* RCRT_APLY_STAT_CD */
		ora.S,   /* ENTP_CFRM_YN */
		ora.S,   /* VP_YN */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* LIVE_REQ_STAT_CD */
		ora.I64, /* ROWNO */
		ora.S,   /* PTO_PATH */
		ora.S,   /* DCMNT_EVL_STAT_CD */
		ora.S,   /* ONWY_INTRV_EVL_STAT_CD */
		ora.S,   /* LIVE_INTRV_EVL_STAT_CD */
		ora.S,   /* READ_END_DAY */ // LDK 2020-11-16 : 합격/불합격 처리 오류, 90일 체크 <-->
		ora.S,   /* Z_ORDER */      // LDK 2021-03-02 : 정렬 <-->
		ora.I64, /* SCORE_VALUE */  // LDK 2021-03-29 : 평가 점수 <-->
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

	recruitStatList := make([]models.RecruitStatList, 0)

	var (
		rslTotCnt          int64
		rslEntpMemNo       string
		rslRecrutSn        string
		rslFavrAplyPpYn    string
		rslNm              string
		rslSex             string
		rslAge             string
		rslRegDt           string
		rslApplyDt         string
		rslEvlStatDt       string
		rslEvlPrgsStatCd   string
		rslRcrtAplyStatCd  string
		rslEntpCfrmYn      string
		rslVpYn            string
		rslPpMemNo         string
		rslLiveReqStatCd   string
		rslRowNo           int64
		rslPtoPath         string
		fullPtoPath        string
		dcmntEvlStatCd     string
		onwyIntrvEvlStatCd string
		liveIntrvEvlStatCd string
		readEndDay         string // LDK 2020-11-16 : 합격/불합격 처리 오류, 90일 체크 <-->
		zOrder             string // LDK 2021-03-02 : 정렬 <-->
		scoreValue         int64  // LDK 2021-03-29 : 평가 점수 <-->
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			rslTotCnt = procRset.Row[0].(int64)
			rslEntpMemNo = procRset.Row[1].(string)
			rslRecrutSn = procRset.Row[2].(string)
			rslFavrAplyPpYn = procRset.Row[3].(string)
			rslNm = procRset.Row[4].(string)
			rslSex = procRset.Row[5].(string)
			rslAge = procRset.Row[6].(string)
			rslRegDt = procRset.Row[7].(string)
			rslApplyDt = procRset.Row[8].(string)
			rslEvlStatDt = procRset.Row[9].(string)
			rslEvlPrgsStatCd = procRset.Row[10].(string)
			rslRcrtAplyStatCd = procRset.Row[11].(string)
			rslEntpCfrmYn = procRset.Row[12].(string)
			rslVpYn = procRset.Row[13].(string)
			rslPpMemNo = procRset.Row[14].(string)
			rslLiveReqStatCd = procRset.Row[15].(string)
			rslRowNo = procRset.Row[16].(int64)
			rslPtoPath = procRset.Row[17].(string)
			if rslPtoPath == "" {
				fullPtoPath = rslPtoPath
			} else {
				fullPtoPath = imgServer + rslPtoPath
			}
			dcmntEvlStatCd = procRset.Row[18].(string)
			onwyIntrvEvlStatCd = procRset.Row[19].(string)
			liveIntrvEvlStatCd = procRset.Row[20].(string)

			readEndDay = procRset.Row[21].(string) // "N" 정상 "Y" 90일이 넘었므로 비정상

			zOrder = procRset.Row[22].(string)
			scoreValue = procRset.Row[23].(int64)

			// var isReadEndDay string
			// if readEndDay >= "0" {
			// 	isReadEndDay = "1"
			// } else {
			// 	isReadEndDay = "0"
			// }

			recruitStatList = append(recruitStatList, models.RecruitStatList{
				RslTotCnt:          rslTotCnt,
				RslEntpMemNo:       rslEntpMemNo,
				RslRecrutSn:        rslRecrutSn,
				RslFavrAplyPpYn:    rslFavrAplyPpYn,
				RslNm:              rslNm,
				RslSex:             rslSex,
				RslAge:             rslAge,
				RslRegDt:           rslRegDt,
				RslApplyDt:         rslApplyDt,
				RslEvlStatDt:       rslEvlStatDt,
				RslEvlPrgsStatCd:   rslEvlPrgsStatCd,
				RslRcrtAplyStatCd:  rslRcrtAplyStatCd,
				RslEntpCfrmYn:      rslEntpCfrmYn,
				RslVpYn:            rslVpYn,
				RslPpMemNo:         rslPpMemNo,
				RslLiveReqStatCd:   rslLiveReqStatCd,
				RslRowNo:           rslRowNo,
				RslPtoPath:         fullPtoPath,
				DcmntEvlStatCd:     dcmntEvlStatCd,
				OnwyIntrvEvlStatCd: onwyIntrvEvlStatCd,
				LiveIntrvEvlStatCd: liveIntrvEvlStatCd,
				ReadEndDay:         readEndDay,
				RslZOrder:          zOrder,
				RslScoreValue:      scoreValue,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}

	var (
		prevPageNo  int64 // 이전 페이지 번호
		nextPageNo  int64 // 다음 페이지 번호
		startPageNo int64
		endPageNo   int64
		totalPage   int64
	)

	prevPageNo = 0
	nextPageNo = 0

	finalPage = (rslTotCnt + (pageSize - 1)) / pageSize // 마지막 페이지
	if pageNo > finalPage {                             // 기본값 설정
		pageNo = finalPage
	}

	if pageNo < 0 || pageNo > finalPage { // 현재 페이지 유효성체크
		pageNo = 1
	}

	var (
		isNowFirst bool
		isNowFinal bool
	)

	if pageNo == 1 {
		isNowFirst = true
	} else {
		isNowFirst = false
	}

	if pageNo == finalPage {
		isNowFinal = true
	} else {
		isNowFinal = false
	}

	if isNowFirst { // 이전페이지 계산
		prevPageNo = 1
	} else {
		if (pageNo - 1) < 1 {
			prevPageNo = 1
		} else {
			prevPageNo = pageNo - 1
		}
	}

	if isNowFinal { // 다음페이지 계산
		nextPageNo = finalPage
	} else {
		if (pageNo + 1) > finalPage {
			nextPageNo = finalPage
		} else {
			nextPageNo = pageNo + 1
		}
	}

	d := float64(pageNo) / float64(pageList)
	blockSize := int64(math.Ceil(d))

	startPageNo = ((blockSize - 1) * pageList) + 1
	endPageNo = startPageNo + pageList - 1

	t := float64(rslTotCnt) / float64(pageSize)
	totalPage = int64(math.Ceil(t))

	if endPageNo > totalPage {
		endPageNo = totalPage
	}

	var pagination string

	if rslTotCnt == 0 {
		pagination += "<a href='javascript:void(0);' class='prev disabled'>이전</a>"
		pagination += "<a href='javascript:void(0);' class='disabled'>1</a>"
		pagination += "<a href='javascript:void(0);' class='next disabled'>다음</a>"
	} else {
		if prevPageNo == pageNo {
			pagination += "<a href='javascript:void(0);' class='prev disabled' id='prev' data-page='" + strconv.Itoa(int(prevPageNo)) + "'>이전</a>"
		} else {
			pagination += "<a href='javascript:void(0);' class='prev goPage' id='prev' data-page='" + strconv.Itoa(int(prevPageNo)) + "'>이전</a>"
		}
		for i := startPageNo; i <= endPageNo; i++ {
			if i == pageNo {
				pagination += "<a href='javascript:void(0);' class='active num goPage' id='num' data-page='" + strconv.Itoa(int(i)) + "'>" + strconv.Itoa(int(i)) + "</a>"
			} else {
				pagination += "<a href='javascript:void(0);' class='num goPage' id='num' data-page='" + strconv.Itoa(int(i)) + "'>" + strconv.Itoa(int(i)) + "</a>"
			}
		}
		if nextPageNo == pageNo {
			pagination += "<a href='javascript:void(0);' class='next disabled' id='prev' data-page='" + strconv.Itoa(int(nextPageNo)) + "'>다음</a>"
		} else {
			pagination += "<a href='javascript:void(0);' class='next goPage' id='prev' data-page='" + strconv.Itoa(int(nextPageNo)) + "'>다음</a>"
		}
	}

	// End : Recruit Apply List

	c.Data["RslTotCnt"] = rslTotCnt
	c.Data["RecruitQuestionList"] = recruitQuestionList

	c.Data["ApplyCnt"] = applyCnt
	c.Data["IngCnt"] = ingCnt
	c.Data["PassCnt"] = passCnt
	c.Data["FailCnt"] = failCnt
	c.Data["DcmntPassCnt"] = dcmntPassCnt
	c.Data["DcmntFailCnt"] = dcmntFailCnt

	c.Data["EntpMemNo"] = entpMemNo
	c.Data["RecrutSn"] = recrutSn
	c.Data["RecrutTitle"] = utils.ConvStringEllipsis(recrutTitle, 0, 40) // 말 줄임표 표시
	c.Data["UpJobGrp"] = upJobGrp
	c.Data["JobGrp"] = jobGrp
	c.Data["RecrutGbnCd"] = recrutGbnCd
	c.Data["RecrutCnt"] = recrutCnt
	c.Data["Rol"] = rol
	c.Data["AplyQufct"] = aplyQufct
	c.Data["PerferTrtm"] = perferTrtm
	c.Data["Sdy"] = sdy
	c.Data["Edy"] = edy
	c.Data["VdTitleUptYn"] = vdTitleUptYn
	c.Data["PrgsStatCd"] = prgsStatCd
	c.Data["RegDt"] = regDt
	c.Data["PpChrgBpNm"] = ppChrgBpNm
	c.Data["PpChrgNm"] = ppChrgNm
	c.Data["EvlPrgsStat"] = pEvlPrgsStat
	c.Data["DcmntEvlStatCd"] = dcmntEvlStatCd
	c.Data["OnwyIntrvEvlStatCd"] = onwyIntrvEvlStatCd
	c.Data["LiveIntrvEvlStatCd"] = liveIntrvEvlStatCd
	c.Data["RecrutProdCd"] = recrutProdCd

	c.Data["RecrutEdt"] = recrutEdt // LDK 2020/11/25: 즉시 마감 날짜 출력

	c.Data["ApplySortCd"] = pApplySortCd
	c.Data["ApplySortWay"] = pApplySortWay

	c.Data["RecruitStatList"] = recruitStatList
	c.Data["Pagination"] = pagination

	c.Data["TMenuId"] = "R00"
	c.Data["SMenuId"] = "R00"

	c.TplName = "recruit/recruit_detail.html"
}

func (c *RecruitPostDetailController) Post() {
	// start : log
	// log := logs.NewLogger()
	// log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/login")
	}

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := mem_no                              // 기업회원번호
	pRecrutSn := c.GetString("recrut_sn")             // 채용일련번호
	pEvlPrgsStat := c.GetString("evl_prgs_stat")      // 평가진행상태 (00:전체, 02:대기, 03:합격, 04:불합격)
	pSex := c.GetString("sex")                        //성별 (무관:A, 남성:M, 여성:F)
	pAge := c.GetString("age")                        //연령 (전체:00, 19세이하:01, 20~29:02, 30~39:03, 40~49:04, 50~59:05, 60세이상:06)
	pVpYn := c.GetString("vp_yn")                     //영상프로필 (전체:9, 있음:1, 없음:0)
	pFavrAplyPp := c.GetString("favr_aply_pp")        //관심지원자 (전체:9, 있음:1, 없음:0)
	pSortGbn := c.GetString("sort_gbn")               // 정렬구분(01:신규순, 02:마감순)
	pKeyword := c.GetString("keyword")                // 검색어
	pLiveReqStatCd := c.GetString("live_req_stat_cd") // 라이브요청상태코드(전체:A, 03:예정)
	pViewType := c.GetString("view_type")             // 보기형시(L:리스트, G:그리드)

	pApplySortCd := c.GetString("apply_sort_cd")
	pApplySortWay := c.GetString("apply_sort_way")

	if pApplySortCd == "" {
		pApplySortCd = "01"
		// LDK 2021/01/29 : Z_ORDER 정렬 설정은 pApplySortCd = "00"
	}
	if pApplySortWay == "" {
		pApplySortWay = "DESC"
		// LDK 2021/01/29 : Z_ORDER 정렬 설정은 pApplySortWay = "ASC"
	}

	imgServer, _  := beego.AppConfig.String("viewpath")

	if pEvlPrgsStat == "" {
		pEvlPrgsStat = "00"
	}

	if pSortGbn == "" {
		pSortGbn = "01"
	}

	if pSex == "" {
		pSex = "A"
	}

	if pAge == "" {
		pAge = "00"
	}

	if pVpYn == "" {
		pVpYn = "9"
	}

	if pFavrAplyPp == "" {
		pFavrAplyPp = "9"
	}

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Get Parameter
	var (
		pageNo    int64
		pageSize  int64
		finalPage int64
		pageList  int64
	)

	pPageNo := c.GetString("pn")
	if pPageNo == "" {
		pPageNo = "1"
	}
	pageNo, err = strconv.ParseInt(pPageNo, 10, 64)
	if err != nil {
		//
	}

	pPageSize := c.GetString("size")
	if pPageSize == "" {
		if pViewType == "L" {
			pPageSize = "30"
		} else {
			pPageSize = "12"
		}
	}

	//pPageSize = strconv.FormatInt(imPageNo, 16)

	pageSize, err = strconv.ParseInt(pPageSize, 10, 64)
	if err != nil {
		//
	}
	pOffSet := (pageNo - 1) * pageSize
	pLimit := pageSize
	pageList = 10

	// Start : Recruit Stat List
	fmt.Printf(fmt.Sprintf("CALL ZSP_APPLY_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pRecrutSn, pEvlPrgsStat, pSex, pAge, pVpYn, pFavrAplyPp, pSortGbn, pKeyword, pLiveReqStatCd, pApplySortCd, pApplySortWay))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_APPLY_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pRecrutSn, pEvlPrgsStat, pSex, pAge, pVpYn, pFavrAplyPp, pSortGbn, pKeyword, pLiveReqStatCd, pApplySortCd, pApplySortWay),
		ora.I64, /* TOT_CNT */
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* FAVR_APLY_PP_YN */
		ora.S,   /* NM */
		ora.S,   /* SEX */
		ora.S,   /* AGE */
		ora.S,   /* REG_DT */
		ora.S,   /* APPLY_DT */
		ora.S,   /* EVL_STAT_DT */
		ora.S,   /* EVL_PRGS_STAT_CD */
		ora.S,   /* RCRT_APLY_STAT_CD */
		ora.S,   /* ENTP_CFRM_YN */
		ora.S,   /* VP_YN */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* LIVE_REQ_STAT_CD */
		ora.I64, /* ROWNO */
		ora.S,   /* PTO_PATH */
		ora.S,   /* DCMNT_EVL_STAT_CD */
		ora.S,   /* ONWY_INTRV_EVL_STAT_CD */
		ora.S,   /* LIVE_INTRV_EVL_STAT_CD */
		ora.S,   /* READ_END_DAY */
		ora.S,   /* Z_ORDER */
		ora.I64, /* SCORE_VALUE */
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

	rtnRecruitStatList := models.RtnRecruitStatList{}
	recruitStatList := make([]models.RecruitStatList, 0)

	var (
		rslTotCnt          int64
		rslEntpMemNo       string
		rslRecrutSn        string
		rslFavrAplyPpYn    string
		rslNm              string
		rslSex             string
		rslAge             string
		rslRegDt           string
		rslApplyDt         string
		rslEvlStatDt       string
		rslEvlPrgsStatCd   string
		rslRcrtAplyStatCd  string
		rslEntpCfrmYn      string
		rslVpYn            string
		rslPpMemNo         string
		rslLiveReqStatCd   string
		rslRowNo           int64
		rslPtoPath         string
		fullPtoPath        string
		dcmntEvlStatCd     string
		onwyIntrvEvlStatCd string
		liveIntrvEvlStatCd string
		readEndDay         string // LDK 2020-11-16 : 합격/불합격 처리 오류, 90일 체크 <-->
		zOrder             string
		scoreValue         int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			rslTotCnt = procRset.Row[0].(int64)
			rslEntpMemNo = procRset.Row[1].(string)
			rslRecrutSn = procRset.Row[2].(string)
			rslFavrAplyPpYn = procRset.Row[3].(string)
			rslNm = procRset.Row[4].(string)
			rslSex = procRset.Row[5].(string)
			rslAge = procRset.Row[6].(string)
			rslRegDt = procRset.Row[7].(string)
			rslApplyDt = procRset.Row[8].(string)
			rslEvlStatDt = procRset.Row[9].(string)
			rslEvlPrgsStatCd = procRset.Row[10].(string)
			rslRcrtAplyStatCd = procRset.Row[11].(string)
			rslEntpCfrmYn = procRset.Row[12].(string)
			rslVpYn = procRset.Row[13].(string)
			rslPpMemNo = procRset.Row[14].(string)
			rslLiveReqStatCd = procRset.Row[15].(string)
			rslRowNo = procRset.Row[16].(int64)
			rslPtoPath = procRset.Row[17].(string)
			dcmntEvlStatCd = procRset.Row[18].(string)
			onwyIntrvEvlStatCd = procRset.Row[19].(string)
			liveIntrvEvlStatCd = procRset.Row[20].(string)

			readEndDay = procRset.Row[21].(string) // "N" 정상 "Y" 90일이 넘었므로 비정상
			zOrder = procRset.Row[22].(string)
			scoreValue = procRset.Row[23].(int64)

			// var isReadEndDay string
			// if readEndDay >= "0" {
			// 	isReadEndDay = "1"
			// } else {
			// 	isReadEndDay = "0"
			// }

			if rslPtoPath == "" {
				fullPtoPath = rslPtoPath
			} else {
				fullPtoPath = imgServer + rslPtoPath
			}

			var (
				prevPageNo  int64 // 이전 페이지 번호
				nextPageNo  int64 // 다음 페이지 번호
				startPageNo int64
				endPageNo   int64
				totalPage   int64
			)

			prevPageNo = 0
			nextPageNo = 0

			finalPage = (rslTotCnt + (pageSize - 1)) / pageSize // 마지막 페이지
			if pageNo > finalPage {                             // 기본값 설정
				pageNo = finalPage
			}

			if pageNo < 0 || pageNo > finalPage { // 현재 페이지 유효성체크
				pageNo = 1
			}

			var (
				isNowFirst bool
				isNowFinal bool
			)

			if pageNo == 1 {
				isNowFirst = true
			} else {
				isNowFirst = false
			}

			if pageNo == finalPage {
				isNowFinal = true
			} else {
				isNowFinal = false
			}

			if isNowFirst { // 이전페이지 계산
				prevPageNo = 1
			} else {
				if (pageNo - 1) < 1 {
					prevPageNo = 1
				} else {
					prevPageNo = pageNo - 1
				}
			}

			if isNowFinal { // 다음페이지 계산
				nextPageNo = finalPage
			} else {
				if (pageNo + 1) > finalPage {
					nextPageNo = finalPage
				} else {
					nextPageNo = pageNo + 1
				}
			}

			d := float64(pageNo) / float64(pageList)
			blockSize := int64(math.Ceil(d))

			startPageNo = ((blockSize - 1) * pageList) + 1
			endPageNo = startPageNo + pageList - 1

			t := float64(rslTotCnt) / float64(pageSize)
			totalPage = int64(math.Ceil(t))

			if endPageNo > totalPage {
				endPageNo = totalPage
			}

			var pagination string

			if rslTotCnt == 0 {
				pagination += "<a href='javascript:void(0);' class='prev disabled'>이전</a>"
				pagination += "<a href='javascript:void(0);' class='disabled'>1</a>"
				pagination += "<a href='javascript:void(0);' class='next disabled'>다음</a>"
			} else {
				if prevPageNo == pageNo {
					pagination += "<a href='javascript:void(0);' class='prev disabled' id='prev' data-page='" + strconv.Itoa(int(prevPageNo)) + "'>이전</a>"
				} else {
					pagination += "<a href='javascript:void(0);' class='prev goPage' id='prev' data-page='" + strconv.Itoa(int(prevPageNo)) + "'>이전</a>"
				}
				for i := startPageNo; i <= endPageNo; i++ {
					if i == pageNo {
						pagination += "<a href='javascript:void(0);' class='active num goPage' id='num' data-page='" + strconv.Itoa(int(i)) + "'>" + strconv.Itoa(int(i)) + "</a>"
					} else {
						pagination += "<a href='javascript:void(0);' class='num goPage' id='num' data-page='" + strconv.Itoa(int(i)) + "'>" + strconv.Itoa(int(i)) + "</a>"
					}
				}
				if nextPageNo == pageNo {
					pagination += "<a href='javascript:void(0);' class='next disabled' id='next' data-page='" + strconv.Itoa(int(nextPageNo)) + "'>다음</a>"
				} else {
					pagination += "<a href='javascript:void(0);' class='next goPage' id='next' data-page='" + strconv.Itoa(int(nextPageNo)) + "'>다음</a>"
				}
			}

			recruitStatList = append(recruitStatList, models.RecruitStatList{
				RslTotCnt:          rslTotCnt,
				RslEntpMemNo:       rslEntpMemNo,
				RslRecrutSn:        rslRecrutSn,
				RslFavrAplyPpYn:    rslFavrAplyPpYn,
				RslNm:              rslNm,
				RslSex:             rslSex,
				RslAge:             rslAge,
				RslRegDt:           rslRegDt,
				RslApplyDt:         rslApplyDt,
				RslEvlStatDt:       rslEvlStatDt,
				RslEvlPrgsStatCd:   rslEvlPrgsStatCd,
				RslRcrtAplyStatCd:  rslRcrtAplyStatCd,
				RslEntpCfrmYn:      rslEntpCfrmYn,
				RslVpYn:            rslVpYn,
				RslPpMemNo:         rslPpMemNo,
				RslLiveReqStatCd:   rslLiveReqStatCd,
				RslRowNo:           rslRowNo,
				RslPtoPath:         fullPtoPath,
				DcmntEvlStatCd:     dcmntEvlStatCd,
				OnwyIntrvEvlStatCd: onwyIntrvEvlStatCd,
				LiveIntrvEvlStatCd: liveIntrvEvlStatCd,
				ReadEndDay:         readEndDay,
				RslZOrder:          zOrder,
				RslScoreValue:      scoreValue,
				Pagination:         pagination,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnRecruitStatList = models.RtnRecruitStatList{
			RtnRecruitStatListData: recruitStatList,

			ApplySortCd:  pApplySortCd,
			ApplySortWay: pApplySortWay,
		}
		// End : Recruit Stat List

		c.Data["json"] = &rtnRecruitStatList
		c.ServeJSON()
	}
}
