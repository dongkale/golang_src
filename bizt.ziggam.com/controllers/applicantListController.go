package controllers

import (
	"fmt"
	"math"
	"strconv"

	"bizt.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type ApplicantListController struct {
	BaseController
}

func (c *ApplicantListController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	mem_id := session.Get(c.Ctx.Request.Context(), "mem_id")
	if mem_id == nil {
		c.Ctx.Redirect(302, "/login")
	}
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")

	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := mem_no //c.GetString("entp_mem_no")
	imgServer, _ := beego.AppConfig.String("viewpath")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Recruit List
	fmt.Printf(fmt.Sprintf("CALL ZSP_CM_RECRUIT_LIST_R('%v', '%v', :1)",
		pLang, pEntpMemNo))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_CM_RECRUIT_LIST_R('%v', '%v', :1)",
		pLang, pEntpMemNo),
		ora.S, /* ENTP_MEM_NO */
		ora.S, /* RECRUT_SN */
		ora.S, /* RECRUT_TITLE */
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

	cmRecrutList := make([]models.CmRecrutList, 0)

	var (
		cmEntpMemNo   string
		cmRecrutSn    string
		cmRecrutTitle string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			cmEntpMemNo = procRset.Row[0].(string)
			cmRecrutSn = procRset.Row[1].(string)
			cmRecrutTitle = procRset.Row[2].(string)

			cmRecrutList = append(cmRecrutList, models.CmRecrutList{
				CmEntpMemNo:   cmEntpMemNo,
				CmRecrutSn:    cmRecrutSn,
				CmRecrutTitle: cmRecrutTitle,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// Recruit List

	// Start : Job Group List
	fmt.Printf(fmt.Sprintf("CALL ZSP_APPLY_JOB_GRP_LIST_R('%v', '%v', :1)",
		pLang, pEntpMemNo))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_APPLY_JOB_GRP_LIST_R('%v', '%v', :1)",
		pLang, pEntpMemNo),
		ora.S,   /* JOB_GRP_NM */
		ora.S,   /* JOB_GRP_CD */
		ora.I64, /* APPLY_CNT */
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

	apJobGrpList := make([]models.ApJobGrpList, 0)

	var (
		apJobGrpNm string
		apJobGrpCd string
		apApplyCnt int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			apJobGrpNm = procRset.Row[0].(string)
			apJobGrpCd = procRset.Row[1].(string)
			apApplyCnt = procRset.Row[2].(int64)

			apJobGrpList = append(apJobGrpList, models.ApJobGrpList{
				ApJobGrpNm: apJobGrpNm,
				ApJobGrpCd: apJobGrpCd,
				ApApplyCnt: apApplyCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// Job Group List

	// Start : Recruit Top Stat Info
	fmt.Printf(fmt.Sprintf("CALL ZSP_APPLICANT_TOP_INFO_R('%v', '%v', :1)",
		pLang, pEntpMemNo))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_APPLICANT_TOP_INFO_R('%v', '%v', :1)",
		pLang, pEntpMemNo),
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
	pRecrutSn := c.GetString("recrut_sn")             // 채용일련번호
	pEvlPrgsStat := c.GetString("evl_prgs_stat")      // 평가진행상태 (00:전체, 02:대기, 03:합격, 04:불합격)
	pSex := c.GetString("sex")                        //성별 (무관:A, 남성:M, 여성:F)
	pAge := c.GetString("age")                        //연령 (전체:00, 19세이하:01, 20~29:02, 30~39:03, 40~49:04, 50~59:05, 60세이상:06)
	pVpYn := c.GetString("vp_yn")                     //영상프로필 (전체:9, 있음:1, 없음:0)
	pFavrAplyPp := c.GetString("favr_aply_pp")        //관심지원자 (전체:9, 있음:1, 없음:0)
	pSortGbn := c.GetString("sort_gbn")               // 정렬구분(01:신규순, 02:마감순)
	pKeyword := c.GetString("keyword")                // 검색어
	pLiveReqStatCd := c.GetString("live_req_stat_cd") // 라이브요청상태코드(전체:A, 03:예정)
	pJobGrpCd := c.GetString("job_grp_cd")            // 직군코드

	pApplySortCd := c.GetString("apply_sort_cd")
	pApplySortWay := c.GetString("apply_sort_way")

	if pApplySortCd == "" {
		pApplySortCd = "01"
	}

	if pApplySortWay == "" {
		pApplySortWay = "DESC"
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

	if pRecrutSn == "" {
		pRecrutSn = "A"
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

	if pJobGrpCd == "" {
		pJobGrpCd = "A"
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
	fmt.Printf(fmt.Sprintf("CALL ZSP_APPLICANT_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v' :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pRecrutSn, pEvlPrgsStat, pSex, pAge, pVpYn, pFavrAplyPp, pSortGbn, pKeyword, pLiveReqStatCd, pJobGrpCd, pApplySortCd, pApplySortWay))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_APPLICANT_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pRecrutSn, pEvlPrgsStat, pSex, pAge, pVpYn, pFavrAplyPp, pSortGbn, pKeyword, pLiveReqStatCd, pJobGrpCd, pApplySortCd, pApplySortWay),
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
		ora.S,   /* RECRUT_TITLE */
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
		recrutTitle        string // LDK 2021-03-30 : 채용 공고명 <-->
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

			recrutTitle = procRset.Row[22].(string)

			// var isReadEndDay string
			// if readEndDay >= "0" {
			// 	isReadEndDay = "1"
			// } else {
			// 	isReadEndDay = "0"
			// }

			//fmt.Printf(fmt.Sprintf("=== %v,%v => %v:%v", rslNm, rslEvlStatDt, readEndDay, isReadEndDay))

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
				RslRecrutTitle:     recrutTitle,
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
	c.Data["CmRecrutList"] = cmRecrutList
	c.Data["RecruitStatList"] = recruitStatList
	c.Data["ApJobGrpList"] = apJobGrpList

	c.Data["ApplyCnt"] = applyCnt
	c.Data["IngCnt"] = ingCnt
	c.Data["PassCnt"] = passCnt
	c.Data["FailCnt"] = failCnt
	c.Data["DcmntPassCnt"] = dcmntPassCnt
	c.Data["DcmntFailCnt"] = dcmntFailCnt

	c.Data["EvlPrgsStat"] = pEvlPrgsStat

	c.Data["ApplySortCd"] = pApplySortCd
	c.Data["ApplySortWay"] = pApplySortWay

	c.Data["Pagination"] = pagination

	c.Data["TMenuId"] = "R00"
	c.Data["SMenuId"] = "R02"
	c.TplName = "applicant/applicant_list.html"
}

func (c *ApplicantListController) Post() {
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
	pJobGrpCd := c.GetString("job_grp_cd")            //직군코드

	pApplySortCd := c.GetString("apply_sort_cd")
	pApplySortWay := c.GetString("apply_sort_way")

	if pApplySortCd == "" {
		pApplySortCd = "01"
	}

	if pApplySortWay == "" {
		pApplySortWay = "DESC"
	}

	imgServer, _ := beego.AppConfig.String("viewpath")

	if pEvlPrgsStat == "" {
		pEvlPrgsStat = "00"
	}

	if pSortGbn == "" {
		pSortGbn = "01"
	}

	if pRecrutSn == "" {
		pRecrutSn = "A"
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

	if pJobGrpCd == "" {
		pJobGrpCd = "A"
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
	fmt.Printf(fmt.Sprintf("CALL ZSP_APPLICANT_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pRecrutSn, pEvlPrgsStat, pSex, pAge, pVpYn, pFavrAplyPp, pSortGbn, pKeyword, pLiveReqStatCd, pJobGrpCd, pApplySortCd, pApplySortWay))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_APPLICANT_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pRecrutSn, pEvlPrgsStat, pSex, pAge, pVpYn, pFavrAplyPp, pSortGbn, pKeyword, pLiveReqStatCd, pJobGrpCd, pApplySortCd, pApplySortWay),
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
		ora.S,   /* RECRUT_TITLE */
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
		recrutTitle        string // LDK 2021-03-30 : 채용 공고명 <-->
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

			readEndDay = procRset.Row[21].(string)

			recrutTitle = procRset.Row[22].(string)

			// var isReadEndDay string
			// if readEndDay >= "0" {
			// 	isReadEndDay = "1"
			// } else {
			// 	isReadEndDay = "0"
			// }

			fmt.Printf(fmt.Sprintf("=== RslApplyDt:%v, RecrutSn:%v, Name:%v => readEndDay:%v", rslApplyDt, rslRecrutSn, rslNm, readEndDay))

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
				RslTotCnt:         rslTotCnt,
				RslEntpMemNo:      rslEntpMemNo,
				RslRecrutSn:       rslRecrutSn,
				RslFavrAplyPpYn:   rslFavrAplyPpYn,
				RslNm:             rslNm,
				RslSex:            rslSex,
				RslAge:            rslAge,
				RslRegDt:          rslRegDt,
				RslApplyDt:        rslApplyDt,
				RslEvlStatDt:      rslEvlStatDt,
				RslEvlPrgsStatCd:  rslEvlPrgsStatCd,
				RslRcrtAplyStatCd: rslRcrtAplyStatCd,
				RslEntpCfrmYn:     rslEntpCfrmYn,
				RslVpYn:           rslVpYn,
				RslPpMemNo:        rslPpMemNo,
				RslLiveReqStatCd:  rslLiveReqStatCd,
				RslRowNo:          rslRowNo,
				RslPtoPath:        fullPtoPath,

				DcmntEvlStatCd:     dcmntEvlStatCd,
				OnwyIntrvEvlStatCd: onwyIntrvEvlStatCd,
				LiveIntrvEvlStatCd: liveIntrvEvlStatCd,
				ReadEndDay:         readEndDay,

				RslRecrutTitle: recrutTitle,

				Pagination: pagination,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnRecruitStatList = models.RtnRecruitStatList{
			RtnRecruitStatListData: recruitStatList,
		}
		// End : Recruit Stat List

		c.Data["json"] = &rtnRecruitStatList
		c.ServeJSON()
	}
}
