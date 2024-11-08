package controllers

import (
	"fmt"
	"math"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"

	"emst.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

type RecruitStatListController struct {
	BaseController
}

func (c *RecruitStatListController) Get() {

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
	pEntpMemNo := mem_no                        // 기업회원번호
	pChkEntpMemNo := c.GetString("entp_mem_no") // 체크 기업회원번호

	if mem_no != pChkEntpMemNo {
		c.Ctx.Redirect(302, "/error/404")
	}

	pRecrutSn := c.GetString("recrut_sn")        // 채용일련번호
	pEvlPrgsStat := c.GetString("evl_prgs_stat") // 평가진행상태 (00:전체, 02:대기, 03:합격, 04:불합격)
	pSex := c.GetString("sex")                   //성별 (무관:A, 남성:M, 여성:F)
	pAge := c.GetString("age")                   //연령 (전체:00, 19세이하:01, 20~29:02, 30~39:03, 40~49:04, 50~59:05, 60세이상:06)
	pVpYn := c.GetString("vp_yn")                //영상프로필 (전체:9, 있음:1, 없음:0)
	pFavrAplyPp := c.GetString("favr_aply_pp")   //관심지원자 (전체:9, 있음:1, 없음:0)
	pSortGbn := c.GetString("sort_gbn")          // 정렬구분(01:신규순, 02:마감순)
	pKeyword := c.GetString("keyword")           // 검색어

	/* Parameter */
	pmKeyword := c.GetString("p_keyword")     // 검색어
	pmEmplTyp := c.GetString("p_empl_typ")    // 고용형태코드
	pmJobGrpCd := c.GetString("p_job_grp_cd") // 직군코드
	pmSortGbn := c.GetString("p_sort_gbn")    // 정렬구분
	pmGbnCd := c.GetString("p_gbn_cd")        // 구분코드
	pmPageNo := c.GetString("p_page_no")      // 페이지번호

	log.Debug("pmEmplTyp : %v", pmEmplTyp)

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
		pPageSize = "9"
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
	log.Debug("CALL SP_EMS_APPLY_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pRecrutSn, pEvlPrgsStat, pSex, pAge, pVpYn, pFavrAplyPp, pSortGbn, pKeyword)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_APPLY_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pRecrutSn, pEvlPrgsStat, pSex, pAge, pVpYn, pFavrAplyPp, pSortGbn, pKeyword),
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
		ora.S,   /* LEFT_DY */
		ora.S,   /* TM */
		ora.S,   /* VP_YN */
		ora.S,   /* PP_MEM_NO */
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

	recruitStatList := make([]models.RecruitStatList, 0)

	var (
		totCnt         int64
		entpMemNo      string
		recrutSn       string
		favrAplyPpYn   string
		nm             string
		sex            string
		age            string
		regDt          string
		applyDt        string
		evlStatDt      string
		evlPrgsStatCd  string
		rcrtAplyStatCd string
		entpCfrmYn     string
		leftDy         string
		tm             string
		vpYn           string
		ppMemNo        string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			totCnt = procRset.Row[0].(int64)
			entpMemNo = procRset.Row[1].(string)
			recrutSn = procRset.Row[2].(string)
			favrAplyPpYn = procRset.Row[3].(string)
			nm = procRset.Row[4].(string)
			sex = procRset.Row[5].(string)
			age = procRset.Row[6].(string)
			regDt = procRset.Row[7].(string)
			applyDt = procRset.Row[8].(string)
			evlStatDt = procRset.Row[9].(string)
			evlPrgsStatCd = procRset.Row[10].(string)
			rcrtAplyStatCd = procRset.Row[11].(string)
			entpCfrmYn = procRset.Row[12].(string)
			leftDy = procRset.Row[13].(string)
			tm = procRset.Row[14].(string)
			vpYn = procRset.Row[15].(string)
			ppMemNo = procRset.Row[16].(string)

			recruitStatList = append(recruitStatList, models.RecruitStatList{
				TotCnt:         totCnt,
				EntpMemNo:      entpMemNo,
				RecrutSn:       recrutSn,
				FavrAplyPpYn:   favrAplyPpYn,
				Nm:             nm,
				Sex:            sex,
				Age:            age,
				RegDt:          regDt,
				ApplyDt:        applyDt,
				EvlStatDt:      evlStatDt,
				EvlPrgsStatCd:  evlPrgsStatCd,
				RcrtAplyStatCd: rcrtAplyStatCd,
				EntpCfrmYn:     entpCfrmYn,
				LeftDy:         leftDy,
				Tm:             tm,
				VpYn:           vpYn,
				PpMemNo:        ppMemNo,
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

	finalPage = (totCnt + (pageSize - 1)) / pageSize // 마지막 페이지
	if pageNo > finalPage {                          // 기본값 설정
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

	t := float64(totCnt) / float64(pageSize)
	totalPage = int64(math.Ceil(t))

	if endPageNo > totalPage {
		endPageNo = totalPage
	}

	var pagination string

	if totCnt == 0 {
		pagination += "<a href='javascript:void(0);' class='btnPrev disabled'>이전</a>"
		pagination += "<a href='javascript:void(0);' class='disabled'>1</a>"
		pagination += "<a href='javascript:void(0);' class='btnNext disabled'>다음</a>"
	} else {
		if prevPageNo == pageNo {
			pagination += "<a href='javascript:void(0);' class='btnPrev disabled' id='prev' data-page='" + strconv.Itoa(int(prevPageNo)) + "'>이전</a>"
		} else {
			pagination += "<a href='javascript:void(0);' class='btnPrev goPage' id='prev' data-page='" + strconv.Itoa(int(prevPageNo)) + "'>이전</a>"
		}
		for i := startPageNo; i <= endPageNo; i++ {
			if i == pageNo {
				pagination += "<a href='javascript:void(0);' class='active num goPage' id='num' data-page='" + strconv.Itoa(int(i)) + "'>" + strconv.Itoa(int(i)) + "</a>"
			} else {
				pagination += "<a href='javascript:void(0);' class='num goPage' id='num' data-page='" + strconv.Itoa(int(i)) + "'>" + strconv.Itoa(int(i)) + "</a>"
			}
		}
		if nextPageNo == pageNo {
			pagination += "<a href='javascript:void(0);' class='btnNext disabled' id='prev' data-page='" + strconv.Itoa(int(nextPageNo)) + "'>다음</a>"
		} else {
			pagination += "<a href='javascript:void(0);' class='btnNext goPage' id='prev' data-page='" + strconv.Itoa(int(nextPageNo)) + "'>다음</a>"
		}
	}

	// End : Recruit Stat List

	// Start : Recruit Stat Info
	log.Debug("CALL SP_EMS_APPLY_TOP_INFO_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_APPLY_TOP_INFO_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn),
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* PRGS_STAT */
		ora.S,   /* RECRUT_TITLE */
		ora.S,   /* EMPL_TYP */
		ora.S,   /* UP_JOB_GRP */
		ora.S,   /* JOB_GRP */
		ora.S,   /* RECRUT_DY */
		ora.S,   /* RECRUT_EDT */
		ora.I64, /* APPLY_CNT */
		ora.I64, /* ING_CNT */
		ora.I64, /* PASS_CNT */
		ora.I64, /* FAIL_CNT */
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
		prgsStat    string
		recrutTitle string
		emplTyp     string
		upJobGrp    string
		jobGrp      string
		recrutDy    string
		recrutEdt   string
		applyCnt    int64
		ingCnt      int64
		passCnt     int64
		failCnt     int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			entpMemNo = procRset.Row[0].(string)
			recrutSn = procRset.Row[1].(string)
			prgsStat = procRset.Row[2].(string)
			recrutTitle = procRset.Row[3].(string)
			emplTyp = procRset.Row[4].(string)
			upJobGrp = procRset.Row[5].(string)
			jobGrp = procRset.Row[6].(string)
			recrutDy = procRset.Row[7].(string)
			recrutEdt = procRset.Row[8].(string)
			applyCnt = procRset.Row[9].(int64)
			ingCnt = procRset.Row[10].(int64)
			passCnt = procRset.Row[11].(int64)
			failCnt = procRset.Row[12].(int64)

			recruitStatTopInfo = append(recruitStatTopInfo, models.RecruitStatTopInfo{
				EntpMemNo:   entpMemNo,
				RecrutSn:    recrutSn,
				PrgsStat:    prgsStat,
				RecrutTitle: recrutTitle,
				EmplTyp:     emplTyp,
				UpJobGrp:    upJobGrp,
				JobGrp:      jobGrp,
				RecrutDy:    recrutDy,
				RecrutEdt:   recrutEdt,
				ApplyCnt:    applyCnt,
				IngCnt:      ingCnt,
				PassCnt:     passCnt,
				FailCnt:     failCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}

	/* Parameter */
	c.Data["pKeyword"] = pmKeyword
	c.Data["pEmplTyp1"] = pmEmplTyp
	c.Data["pJobGrpCd"] = pmJobGrpCd
	c.Data["pSortGbn"] = pmSortGbn
	c.Data["pGbnCd"] = pmGbnCd
	c.Data["pPageNo"] = pmPageNo

	c.Data["MemNo"] = mem_no
	c.Data["RecrutSn"] = recrutSn
	c.Data["PrgsStat"] = prgsStat
	c.Data["RecrutTitle"] = recrutTitle
	c.Data["EmplTyp"] = emplTyp
	c.Data["UpJobGrp"] = upJobGrp
	c.Data["JobGrp"] = jobGrp
	c.Data["RecrutDy"] = recrutDy
	c.Data["RecrutEdt"] = recrutEdt
	c.Data["ApplyCnt"] = applyCnt
	c.Data["IngCnt"] = ingCnt
	c.Data["PassCnt"] = passCnt
	c.Data["FailCnt"] = failCnt

	c.Data["EvlPrgsStat"] = pEvlPrgsStat
	c.Data["MenuId"] = "02"

	c.Data["RecruitStatList"] = recruitStatList
	c.Data["Pagination"] = pagination
	c.Data["PageNo"] = pageNo

	c.TplName = "recruit/recruit_stat_list.html"
}

func (c *RecruitStatListController) Post() {

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
	pEntpMemNo := mem_no                         // 기업회원번호
	pRecrutSn := c.GetString("recrut_sn")        // 채용일련번호
	pEvlPrgsStat := c.GetString("evl_prgs_stat") // 평가진행상태 (00:전체, 02:대기, 03:합격, 04:불합격)
	pSex := c.GetString("sex")                   //성별 (무관:A, 남성:M, 여성:F)
	pAge := c.GetString("age")                   //연령 (전체:00, 19세이하:01, 20~29:02, 30~39:03, 40~49:04, 50~59:05, 60세이상:06)
	pVpYn := c.GetString("vp_yn")                //영상프로필 (전체:9, 있음:1, 없음:0)
	pFavrAplyPp := c.GetString("favr_aply_pp")   //관심지원자 (전체:9, 있음:1, 없음:0)
	pSortGbn := c.GetString("sort_gbn")          // 정렬구분(01:신규순, 02:마감순)
	pKeyword := c.GetString("keyword")           // 검색어

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
		pPageSize = "6"
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
	log.Debug("CALL SP_EMS_APPLY_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pRecrutSn, pEvlPrgsStat, pSex, pAge, pVpYn, pFavrAplyPp, pSortGbn, pKeyword)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_APPLY_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pRecrutSn, pEvlPrgsStat, pSex, pAge, pVpYn, pFavrAplyPp, pSortGbn, pKeyword),
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
		ora.S,   /* LEFT_DY */
		ora.S,   /* TM */
		ora.S,   /* VP_YN */
		ora.S,   /* PP_MEM_NO */
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
		totCnt         int64
		entpMemNo      string
		recrutSn       string
		favrAplyPpYn   string
		nm             string
		sex            string
		age            string
		regDt          string
		applyDt        string
		evlStatDt      string
		evlPrgsStatCd  string
		rcrtAplyStatCd string
		entpCfrmYn     string
		leftDy         string
		tm             string
		vpYn           string
		ppMemNo        string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			totCnt = procRset.Row[0].(int64)
			entpMemNo = procRset.Row[1].(string)
			recrutSn = procRset.Row[2].(string)
			favrAplyPpYn = procRset.Row[3].(string)
			nm = procRset.Row[4].(string)
			sex = procRset.Row[5].(string)
			age = procRset.Row[6].(string)
			regDt = procRset.Row[7].(string)
			applyDt = procRset.Row[8].(string)
			evlStatDt = procRset.Row[9].(string)
			evlPrgsStatCd = procRset.Row[10].(string)
			rcrtAplyStatCd = procRset.Row[11].(string)
			entpCfrmYn = procRset.Row[12].(string)
			leftDy = procRset.Row[13].(string)
			tm = procRset.Row[14].(string)
			vpYn = procRset.Row[15].(string)
			ppMemNo = procRset.Row[16].(string)

			var (
				prevPageNo  int64 // 이전 페이지 번호
				nextPageNo  int64 // 다음 페이지 번호
				startPageNo int64
				endPageNo   int64
				totalPage   int64
			)

			prevPageNo = 0
			nextPageNo = 0

			finalPage = (totCnt + (pageSize - 1)) / pageSize // 마지막 페이지
			if pageNo > finalPage {                          // 기본값 설정
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

			t := float64(totCnt) / float64(pageSize)
			totalPage = int64(math.Ceil(t))

			if endPageNo > totalPage {
				endPageNo = totalPage
			}

			var pagination string

			if totCnt == 0 {
				pagination += "<a href='javascript:void(0);' class='btnPrev disabled'>이전</a>"
				pagination += "<a href='javascript:void(0);' class='disabled'>1</a>"
				pagination += "<a href='javascript:void(0);' class='btnNext disabled'>다음</a>"
			} else {
				if prevPageNo == pageNo {
					pagination += "<a href='javascript:void(0);' class='btnPrev disabled' id='prev' data-page='" + strconv.Itoa(int(prevPageNo)) + "'>이전</a>"
				} else {
					pagination += "<a href='javascript:void(0);' class='btnPrev goPage' id='prev' data-page='" + strconv.Itoa(int(prevPageNo)) + "'>이전</a>"
				}
				for i := startPageNo; i <= endPageNo; i++ {
					if i == pageNo {
						pagination += "<a href='javascript:void(0);' class='active num goPage' id='num' data-page='" + strconv.Itoa(int(i)) + "'>" + strconv.Itoa(int(i)) + "</a>"
					} else {
						pagination += "<a href='javascript:void(0);' class='num goPage' id='num' data-page='" + strconv.Itoa(int(i)) + "'>" + strconv.Itoa(int(i)) + "</a>"
					}
				}
				if nextPageNo == pageNo {
					pagination += "<a href='javascript:void(0);' class='btnNext disabled' id='prev' data-page='" + strconv.Itoa(int(nextPageNo)) + "'>다음</a>"
				} else {
					pagination += "<a href='javascript:void(0);' class='btnNext goPage' id='prev' data-page='" + strconv.Itoa(int(nextPageNo)) + "'>다음</a>"
				}
			}

			recruitStatList = append(recruitStatList, models.RecruitStatList{
				TotCnt:         totCnt,
				EntpMemNo:      entpMemNo,
				RecrutSn:       recrutSn,
				FavrAplyPpYn:   favrAplyPpYn,
				Nm:             nm,
				Sex:            sex,
				Age:            age,
				RegDt:          regDt,
				ApplyDt:        applyDt,
				EvlStatDt:      evlStatDt,
				EvlPrgsStatCd:  evlPrgsStatCd,
				RcrtAplyStatCd: rcrtAplyStatCd,
				EntpCfrmYn:     entpCfrmYn,
				LeftDy:         leftDy,
				Tm:             tm,
				VpYn:           vpYn,
				PpMemNo:        ppMemNo,
				Pagination:     pagination,
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
