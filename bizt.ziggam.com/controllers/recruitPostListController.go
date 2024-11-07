package controllers

import (
	"fmt"
	"math"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"bizt.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

type RecruitPostListController struct {
	BaseController
}

func (c *RecruitPostListController) Get() {

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
	pEntpMemNo := mem_no // 기업회원번호(세션)

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Recruit Main List
	fmt.Printf(fmt.Sprintf("CALL ZSP_RECRUIT_MAIN_LIST_R('%v', '%v', :1)",
		pLang, pEntpMemNo))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_RECRUIT_MAIN_LIST_R('%v', '%v', :1)",
		pLang, pEntpMemNo),
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
		ora.I64, /* TOT_CNT */
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

	recruitPostList := make([]models.RecruitPostList, 0)

	var (
		entpMemNo   string
		recrutSn    string
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
		totCnt      int64
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
			totCnt = procRset.Row[13].(int64)

			recruitPostList = append(recruitPostList, models.RecruitPostList{
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
				TotCnt:      totCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}

	// End : Recruit Main List

	// Start : Recruit Stat Info
	fmt.Printf(fmt.Sprintf("CALL ZSP_RECRUIT_STAT_INFO_R('%v', '%v', :1)",
		pLang, pEntpMemNo))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_RECRUIT_STAT_INFO_R('%v', '%v', :1)",
		pLang, pEntpMemNo),
		ora.I64, /* RECRUIT_TOT_CNT */
		ora.I64, /* RECRUIT_ING_CNT */
		ora.I64, /* RECRUIT_WAIT_CNT */
		ora.I64, /* RECRUIT_END_CNT */
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

	recruitStatInfo := make([]models.RecruitStatInfo, 0)

	var (
		recrutTotCnt  int64
		recrutIngCnt  int64
		recrutWaitCnt int64
		recrutEndCnt  int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			recrutTotCnt = procRset.Row[0].(int64)
			recrutIngCnt = procRset.Row[1].(int64)
			recrutWaitCnt = procRset.Row[2].(int64)
			recrutEndCnt = procRset.Row[3].(int64)

			recruitStatInfo = append(recruitStatInfo, models.RecruitStatInfo{
				RecrutTotCnt:  recrutTotCnt,
				RecrutIngCnt:  recrutIngCnt,
				RecrutEndCnt:  recrutEndCnt,
				RecrutWaitCnt: recrutWaitCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}

	pEmplTypCd := "01"

	// Start : Job Group List
	fmt.Printf(fmt.Sprintf("CALL ZSP_RECRUIT_JOB_GRP_LIST_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pEmplTypCd))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_RECRUIT_JOB_GRP_LIST_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pEmplTypCd),
		ora.S, /* EMPL_TYP_CD */
		ora.S, /* JOB_GRP_CD */
		ora.S, /* JOB_GRP_NM */

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

	recruitMainJobGrpList := make([]models.RecruitMainJobGrpList, 0)

	var (
		rEmplTypCd string
		rJobGrpCd  string
		rJobGrpNm  string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			rEmplTypCd = procRset.Row[0].(string)
			rJobGrpCd = procRset.Row[1].(string)
			rJobGrpNm = procRset.Row[2].(string)

			recruitMainJobGrpList = append(recruitMainJobGrpList, models.RecruitMainJobGrpList{
				REmplTypCd: rEmplTypCd,
				RJobGrpCd:  rJobGrpCd,
				RJobGrpNm:  rJobGrpNm,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Job Group List

	pGbnCd := c.GetString("gbn_cd")        // 구분코드(A:전체, I:채용중, E:종료)
	pEmplTyp := "R"                        // 고용형태(A:아르바이트, R:정규직/인턴)
	pJobGrpCd := c.GetString("job_grp_cd") // 2차직무코드
	pSortGbn := c.GetString("sort_gbn")    // 정렬구분(01:등록일순, 02:마감일순, 03:최신순)
	pKeyword := c.GetString("keyword")     // 검색어
	pSday := c.GetString("sday")           // 검색시작일자
	pEday := c.GetString("eday")           // 검색종료일자

	/* Parameter */
	pmKeyword := c.GetString("p_keyword")     // 검색어
	pmJobGrpCd := c.GetString("p_job_grp_cd") // 2차직무코드
	pmSortGbn := c.GetString("p_sort_gbn")    // 정렬구분
	pmGbnCd := c.GetString("p_gbn_cd")        // 탭구분
	pmSday := c.GetString("p_sday")           // 검색시작일자
	pmEday := c.GetString("p_eday")           // 검색종료일자
	pmPageNo := c.GetString("p_page_no")      // 페이지번호

	if pGbnCd == "" {
		if pmGbnCd == "" {
			pGbnCd = "A"
		} else {
			pGbnCd = pmGbnCd
		}
	}

	if pSortGbn == "" {
		pSortGbn = "03" //최신순
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
		if pmPageNo == "" {
			pPageNo = "1"
		} else {
			pPageNo = pmPageNo
		}
	}
	pageNo, err = strconv.ParseInt(pPageNo, 10, 64)
	if err != nil {
		//
	}

	pPageSize := c.GetString("size")
	if pPageSize == "" {
		pPageSize = "10"
	}

	//pPageSize = strconv.FormatInt(imPageNo, 16)

	pageSize, err = strconv.ParseInt(pPageSize, 10, 64)
	if err != nil {
		//
	}
	pOffSet := (pageNo - 1) * pageSize
	pLimit := pageSize
	pageList = 10

	// Start : Recruit sub List

	fmt.Printf(fmt.Sprintf("CALL ZSP_RECRUIT_SUB_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pGbnCd, pEmplTyp, pJobGrpCd, pSortGbn, pKeyword, pSday, pEday))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_RECRUIT_SUB_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pGbnCd, pEmplTyp, pJobGrpCd, pSortGbn, pKeyword, pSday, pEday),
		ora.I64, /* TOT_CNT */
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
		ora.S,   /* REG_DT */
		ora.S,   /* REG_ID */
		ora.S,   /* PP_CHRG_BP_NM */
		ora.S,   /* PP_CHRG_NM */
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

	recruitSubList := make([]models.RecruitSubList, 0)

	var (
		sTotCnt      int64
		sEntpMemNo   string
		sRecrutSn    string
		sPrgsStat    string
		sRecrutTitle string
		sEmplTyp     string
		sUpJobGrp    string
		sJobGrp      string
		sRecrutDy    string
		sRecrutEdt   string
		sApplyCnt    int64
		sIngCnt      int64
		sPassCnt     int64
		sFailCnt     int64
		sRegDt       string
		sRegId       string
		sPpChrgBpNm  string
		sPpChrgNm    string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			sTotCnt = procRset.Row[0].(int64)
			sEntpMemNo = procRset.Row[1].(string)
			sRecrutSn = procRset.Row[2].(string)
			sPrgsStat = procRset.Row[3].(string)
			sRecrutTitle = procRset.Row[4].(string)
			sEmplTyp = procRset.Row[5].(string)
			sUpJobGrp = procRset.Row[6].(string)
			sJobGrp = procRset.Row[7].(string)
			sRecrutDy = procRset.Row[8].(string)
			sRecrutEdt = procRset.Row[9].(string)
			sApplyCnt = procRset.Row[10].(int64)
			sIngCnt = procRset.Row[11].(int64)
			sPassCnt = procRset.Row[12].(int64)
			sFailCnt = procRset.Row[13].(int64)
			sRegDt = procRset.Row[14].(string)
			sRegId = procRset.Row[15].(string)
			sPpChrgBpNm = procRset.Row[16].(string)
			sPpChrgNm = procRset.Row[17].(string)

			recruitSubList = append(recruitSubList, models.RecruitSubList{
				STotCnt:      sTotCnt,
				SEntpMemNo:   sEntpMemNo,
				SRecrutSn:    sRecrutSn,
				SPrgsStat:    sPrgsStat,
				SRecrutTitle: sRecrutTitle,
				SEmplTyp:     sEmplTyp,
				SUpJobGrp:    sUpJobGrp,
				SJobGrp:      sJobGrp,
				SRecrutDy:    sRecrutDy,
				SRecrutEdt:   sRecrutEdt,
				SApplyCnt:    sApplyCnt,
				SIngCnt:      sIngCnt,
				SPassCnt:     sPassCnt,
				SFailCnt:     sFailCnt,
				SRegDt:       sRegDt,
				SRegId:       sRegId,
				SPpChrgBpNm:  sPpChrgBpNm,
				SPpChrgNm:    sPpChrgNm,
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

	finalPage = (sTotCnt + (pageSize - 1)) / pageSize // 마지막 페이지
	if pageNo > finalPage {                           // 기본값 설정
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

	t := float64(sTotCnt) / float64(pageSize)
	totalPage = int64(math.Ceil(t))

	if endPageNo > totalPage {
		endPageNo = totalPage
	}

	var pagination string

	// LDK 2020/11/12 : 채용 공고 리스트 페이지 처리 오류(inviteSendListController.go,applicantListController.go 참고) -> 주석 처리 -->
	if sTotCnt == 0 {
		pagination += "<a href='javascript:void(0);' class='prev disabled'>이전</a>"
		//pagination += "<span>"
		pagination += "<a href='javascript:void(0);' class='disabled'>1</a>"
		//pagination += "</span>"
		pagination += "<a href='javascript:void(0);' class='next disabled'>다음</a>"
	} else {
		if prevPageNo == pageNo {
			pagination += "<a href='javascript:void(0);' class='prev disabled' id='prev' data-page='" + strconv.Itoa(int(prevPageNo)) + "'>이전</a>"
			//pagination += "<span>"
		} else {
			pagination += "<a href='javascript:void(0);' class='prev goPage' id='prev' data-page='" + strconv.Itoa(int(prevPageNo)) + "'>이전</a>"
			//pagination += "<span>"
		}
		for i := startPageNo; i <= endPageNo; i++ {
			if i == pageNo {
				pagination += "<a href='javascript:void(0);' class='active num goPage' id='num' data-page='" + strconv.Itoa(int(i)) + "'>" + strconv.Itoa(int(i)) + "</a>"
			} else {
				pagination += "<a href='javascript:void(0);' class='num goPage' id='num' data-page='" + strconv.Itoa(int(i)) + "'>" + strconv.Itoa(int(i)) + "</a>"
			}
		}
		if nextPageNo == pageNo {
			//pagination += "</span>"
			//pagination += "<a href='javascript:void(0);' class='btnNext next' id='next' data-page='" + strconv.Itoa(int(nextPageNo)) + "'>다음</a>"
			pagination += "<a href='javascript:void(0);' class='next disabled' id='prev' data-page='" + strconv.Itoa(int(nextPageNo)) + "'>다음</a>"
		} else {
			//pagination += "</span>"
			//pagination += "<a href='javascript:void(0);' class='btnNext next' id='next' data-page='" + strconv.Itoa(int(nextPageNo)) + "'>다음</a>"
			pagination += "<a href='javascript:void(0);' class='next goPage' id='prev' data-page='" + strconv.Itoa(int(nextPageNo)) + "'>다음</a>"
		}
	}
	// <--

	// End : Recruit Sub List

	c.Data["RecrutTotCnt"] = recrutTotCnt
	c.Data["RecrutIngCnt"] = recrutIngCnt
	c.Data["RecrutWaitCnt"] = recrutWaitCnt
	c.Data["RecrutEndCnt"] = recrutEndCnt

	c.Data["RecruitSubList"] = recruitSubList
	c.Data["TotCnt"] = totCnt
	c.Data["STotCnt"] = sTotCnt
	c.Data["GbnCd"] = pGbnCd

	c.Data["RecruitMainJobGrpList"] = recruitMainJobGrpList

	c.Data["RecruitPostList"] = recruitPostList
	c.Data["Pagination"] = pagination
	c.Data["PageNo"] = pageNo

	/* Parameter */
	c.Data["pKeyword"] = pmKeyword
	c.Data["pJobGrpCd"] = pmJobGrpCd
	c.Data["pSortGbn"] = pmSortGbn
	c.Data["pGbnCd"] = pmGbnCd
	c.Data["pSday"] = pmSday
	c.Data["pEday"] = pmEday
	c.Data["pPageNo"] = pmPageNo

	c.Data["TMenuId"] = "R00"
	c.Data["SMenuId"] = "R01"
	c.TplName = "recruit/recruit_post_list.html"
}

func (c *RecruitPostListController) Post() {
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
	pEntpMemNo := mem_no                   // 기업회원번호(세션)
	pGbnCd := c.GetString("gbn_cd")        // 구분코드(A:전체, I:채용중, E:종료)
	pEmplTyp := "R"                        // 고용형태(R:정규직)
	pJobGrpCd := c.GetString("job_grp_cd") // 2차직무코드
	pSortGbn := c.GetString("sort_gbn")    // 정렬구분(01:등록일순, 02:마감일순)
	pKeyword := c.GetString("keyword")     // 검색어
	pSday := c.GetString("sday")           // 검색시작일자
	pEday := c.GetString("eday")           // 검색종료일자

	pmGbnCd := c.GetString("p_gbn_cd")   // 탭구분
	pmPageNo := c.GetString("p_page_no") // 페이지번호

	if pGbnCd == "" {
		if pmGbnCd == "" {
			pGbnCd = "A"
		} else {
			pGbnCd = pmGbnCd
		}
	}

	if pSortGbn == "" {
		pSortGbn = "03" //최신순
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
		if pmPageNo == "" {
			pPageNo = "1"
		} else {
			pPageNo = pmPageNo
		}
	}
	pageNo, err = strconv.ParseInt(pPageNo, 10, 64)
	if err != nil {
		//
	}

	pPageSize := c.GetString("size")
	if pPageSize == "" {
		pPageSize = "10"
	}

	//pPageSize = strconv.FormatInt(imPageNo, 16)

	pageSize, err = strconv.ParseInt(pPageSize, 10, 64)
	if err != nil {
		//
	}
	pOffSet := (pageNo - 1) * pageSize
	pLimit := pageSize
	pageList = 10

	fmt.Printf(fmt.Sprintf("CALL ZSP_RECRUIT_SUB_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pGbnCd, pEmplTyp, pJobGrpCd, pSortGbn, pKeyword, pSday, pEday))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_RECRUIT_SUB_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pGbnCd, pEmplTyp, pJobGrpCd, pSortGbn, pKeyword, pSday, pEday),
		ora.I64, /* TOT_CNT */
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
		ora.S,   /* REG_DT */
		ora.S,   /* REG_ID */
		ora.S,   /* PP_CHRG_BP_NM */
		ora.S,   /* PP_CHRG_NM */
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

	rtnRecruitSubList := models.RtnRecruitSubList{}
	recruitSubList := make([]models.RecruitSubList, 0)

	var (
		sTotCnt      int64
		sEntpMemNo   string
		sRecrutSn    string
		sPrgsStat    string
		sRecrutTitle string
		sEmplTyp     string
		sUpJobGrp    string
		sJobGrp      string
		sRecrutDy    string
		sRecrutEdt   string
		sApplyCnt    int64
		sIngCnt      int64
		sPassCnt     int64
		sFailCnt     int64
		sRegDt       string
		sRegId       string
		sPpChrgBpNm  string
		sPpChrgNm    string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			sTotCnt = procRset.Row[0].(int64)
			sEntpMemNo = procRset.Row[1].(string)
			sRecrutSn = procRset.Row[2].(string)
			sPrgsStat = procRset.Row[3].(string)
			sRecrutTitle = procRset.Row[4].(string)
			sEmplTyp = procRset.Row[5].(string)
			sUpJobGrp = procRset.Row[6].(string)
			sJobGrp = procRset.Row[7].(string)
			sRecrutDy = procRset.Row[8].(string)
			sRecrutEdt = procRset.Row[9].(string)
			sApplyCnt = procRset.Row[10].(int64)
			sIngCnt = procRset.Row[11].(int64)
			sPassCnt = procRset.Row[12].(int64)
			sFailCnt = procRset.Row[13].(int64)
			sRegDt = procRset.Row[14].(string)
			sRegId = procRset.Row[15].(string)
			sPpChrgBpNm = procRset.Row[16].(string)
			sPpChrgNm = procRset.Row[17].(string)

			var (
				prevPageNo  int64 // 이전 페이지 번호
				nextPageNo  int64 // 다음 페이지 번호
				startPageNo int64
				endPageNo   int64
				totalPage   int64
			)

			prevPageNo = 0
			nextPageNo = 0

			finalPage = (sTotCnt + (pageSize - 1)) / pageSize // 마지막 페이지
			if pageNo > finalPage {                           // 기본값 설정
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

			t := float64(sTotCnt) / float64(pageSize)
			totalPage = int64(math.Ceil(t))

			if endPageNo > totalPage {
				endPageNo = totalPage
			}

			var pagination string

			// LDK 2020/11/12 : 채용 공고 리스트 페이지 처리 오류(inviteSendListController.go,applicantListController.go 참고) -> 주석 처리 -->
			if sTotCnt == 0 {
				pagination += "<a href='javascript:void(0);' class='prev disabled'>이전</a>"
				//pagination += "<span>"
				pagination += "<a href='javascript:void(0);' class='disabled'>1</a>"
				//pagination += "</span>"
				pagination += "<a href='javascript:void(0);' class='next disabled'>다음</a>"
			} else {
				if prevPageNo == pageNo {
					pagination += "<a href='javascript:void(0);' class='prev disabled' id='prev' data-page='" + strconv.Itoa(int(prevPageNo)) + "'>이전</a>"
					//pagination += "<span>"
				} else {
					pagination += "<a href='javascript:void(0);' class='prev goPage' id='prev' data-page='" + strconv.Itoa(int(prevPageNo)) + "'>이전</a>"
					//pagination += "<span>"
				}
				for i := startPageNo; i <= endPageNo; i++ {
					if i == pageNo {
						pagination += "<a href='javascript:void(0);' class='active num goPage' id='num' data-page='" + strconv.Itoa(int(i)) + "'>" + strconv.Itoa(int(i)) + "</a>"
					} else {
						pagination += "<a href='javascript:void(0);' class='num goPage' id='num' data-page='" + strconv.Itoa(int(i)) + "'>" + strconv.Itoa(int(i)) + "</a>"
					}
				}
				if nextPageNo == pageNo {
					//pagination += "</span>"
					//pagination += "<a href='javascript:void(0);' class='btnNext next' id='next' data-page='" + strconv.Itoa(int(nextPageNo)) + "'>다음</a>"
					pagination += "<a href='javascript:void(0);' class='next disabled' id='next' data-page='" + strconv.Itoa(int(nextPageNo)) + "'>다음</a>"
				} else {
					//pagination += "</span>"
					//pagination += "<a href='javascript:void(0);' class='btnNext next' id='next' data-page='" + strconv.Itoa(int(nextPageNo)) + "'>다음</a>"
					pagination += "<a href='javascript:void(0);' class='next goPage' id='next' data-page='" + strconv.Itoa(int(nextPageNo)) + "'>다음</a>"
				}
			}
			// <--

			recruitSubList = append(recruitSubList, models.RecruitSubList{
				STotCnt:      sTotCnt,
				SEntpMemNo:   sEntpMemNo,
				SRecrutSn:    sRecrutSn,
				SPrgsStat:    sPrgsStat,
				SRecrutTitle: sRecrutTitle,
				SEmplTyp:     sEmplTyp,
				SUpJobGrp:    sUpJobGrp,
				SJobGrp:      sJobGrp,
				SRecrutDy:    sRecrutDy,
				SRecrutEdt:   sRecrutEdt,
				SApplyCnt:    sApplyCnt,
				SIngCnt:      sIngCnt,
				SPassCnt:     sPassCnt,
				SFailCnt:     sFailCnt,
				SRegDt:       sRegDt,
				SRegId:       sRegId,
				SPpChrgBpNm:  sPpChrgBpNm,
				SPpChrgNm:    sPpChrgNm,
				Pagination:   pagination,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnRecruitSubList = models.RtnRecruitSubList{
			RtnRecruitSubListData: recruitSubList,
		}

		c.Data["json"] = &rtnRecruitSubList
		c.ServeJSON()
	}
}
