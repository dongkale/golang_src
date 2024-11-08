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

type RecruitPostListController struct {
	BaseController
}

func (c *RecruitPostListController) Get() {

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
	pEmplTyp := c.GetString("empl_typ")    // 고용형태(A:아르바이트, R:정규직/인턴)
	pJobGrpCd := c.GetString("job_grp_cd") // 2차직무코드
	pSortGbn := c.GetString("sort_gbn")    // 정렬구분(01:등록일순, 02:마감일순)
	pKeyword := c.GetString("keyword")     // 검색어

	/* Parameter */
	pmKeyword := c.GetString("p_keyword")     // 검색어
	pmEmplTyp := c.GetString("p_empl_typ")    // 고용형태
	pmJobGrpCd := c.GetString("p_job_grp_cd") // 2차직무코드
	pmSortGbn := c.GetString("p_sort_gbn")    // 정렬구분
	pmGbnCd := c.GetString("p_gbn_cd")        // 탭구분
	pmPageNo := c.GetString("p_page_no")      // 페이지번호

	if pGbnCd == "" {
		if pmGbnCd == "" {
			pGbnCd = "A"
		} else {
			pGbnCd = pmGbnCd
		}
	}

	if pSortGbn == "" {
		pSortGbn = "01"
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

	// Start : Recruit Main List
	log.Debug("CALL SP_EMS_RECRUIT_MAIN_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pGbnCd, pEmplTyp, pJobGrpCd, pSortGbn, pKeyword)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_RECRUIT_MAIN_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pGbnCd, pEmplTyp, pJobGrpCd, pSortGbn, pKeyword),
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
		totCnt      int64
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
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			totCnt = procRset.Row[0].(int64)
			entpMemNo = procRset.Row[1].(string)
			recrutSn = procRset.Row[2].(string)
			prgsStat = procRset.Row[3].(string)
			recrutTitle = procRset.Row[4].(string)
			emplTyp = procRset.Row[5].(string)
			upJobGrp = procRset.Row[6].(string)
			jobGrp = procRset.Row[7].(string)
			recrutDy = procRset.Row[8].(string)
			recrutEdt = procRset.Row[9].(string)
			applyCnt = procRset.Row[10].(int64)
			ingCnt = procRset.Row[11].(int64)
			passCnt = procRset.Row[12].(int64)
			failCnt = procRset.Row[13].(int64)

			recruitPostList = append(recruitPostList, models.RecruitPostList{
				TotCnt:      totCnt,
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

	// End : Recruit Main List

	// Start : Recruit Stat Info
	log.Debug("CALL SP_EMS_RECRUIT_STAT_INFO_R('%v', '%v', :1)",
		pLang, pEntpMemNo)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_RECRUIT_STAT_INFO_R('%v', '%v', :1)",
		pLang, pEntpMemNo),
		ora.I64, /* RECRUIT_TOT_CNT */
		ora.I64, /* RECRUIT_ING_CNT */
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
		recrutTotCnt int64
		recrutIngCnt int64
		recrutEndCnt int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			recrutTotCnt = procRset.Row[0].(int64)
			recrutIngCnt = procRset.Row[1].(int64)
			recrutEndCnt = procRset.Row[2].(int64)

			recruitStatInfo = append(recruitStatInfo, models.RecruitStatInfo{
				RecrutTotCnt: recrutTotCnt,
				RecrutIngCnt: recrutIngCnt,
				RecrutEndCnt: recrutEndCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}

	c.Data["MemNo"] = mem_no
	c.Data["RecrutTotCnt"] = recrutTotCnt
	c.Data["RecrutIngCnt"] = recrutIngCnt
	c.Data["RecrutEndCnt"] = recrutEndCnt

	c.Data["GbnCd"] = pGbnCd
	c.Data["MenuId"] = "02"

	/* Parameter */
	c.Data["pKeyword"] = pmKeyword
	c.Data["pEmplTyp1"] = pmEmplTyp
	c.Data["pJobGrpCd"] = pmJobGrpCd
	c.Data["pSortGbn"] = pmSortGbn
	c.Data["pGbnCd"] = pmGbnCd
	c.Data["pPageNo"] = pmPageNo

	c.Data["RecruitPostList"] = recruitPostList
	c.Data["Pagination"] = pagination
	c.Data["PageNo"] = pageNo

	c.TplName = "recruit/recruit_list.html"
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
	pEntpMemNo := mem_no                // 기업회원번호(세션)
	pGbnCd := c.GetString("gbn_cd")     // 구분코드(A:전체, I:채용중, E:종료)
	pEmplTyp := c.GetString("empl_typ") // 고용형태(A:아르바이트, R:정규직/인턴)
	if pEmplTyp == "01" {
		pEmplTyp = "R"
	} else if pEmplTyp == "05" {
		pEmplTyp = "A"
	} else {
		pEmplTyp = ""
	}
	pJobGrpCd := c.GetString("job_grp_cd") // 2차직무코드
	pSortGbn := c.GetString("sort_gbn")    // 정렬구분(01:등록일순, 02:마감일순)
	pKeyword := c.GetString("keyword")     // 검색어

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
		pSortGbn = "01"
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

	// Start : Recruit Main List
	log.Debug("CALL SP_EMS_RECRUIT_MAIN_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pGbnCd, pEmplTyp, pJobGrpCd, pSortGbn, pKeyword)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_RECRUIT_MAIN_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pGbnCd, pEmplTyp, pJobGrpCd, pSortGbn, pKeyword),
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

	rtnRecruitPostList := models.RtnRecruitPostList{}
	recruitPostList := make([]models.RecruitPostList, 0)

	var (
		totCnt      int64
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
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			totCnt = procRset.Row[0].(int64)
			entpMemNo = procRset.Row[1].(string)
			recrutSn = procRset.Row[2].(string)
			prgsStat = procRset.Row[3].(string)
			recrutTitle = procRset.Row[4].(string)
			emplTyp = procRset.Row[5].(string)
			upJobGrp = procRset.Row[6].(string)
			jobGrp = procRset.Row[7].(string)
			recrutDy = procRset.Row[8].(string)
			recrutEdt = procRset.Row[9].(string)
			applyCnt = procRset.Row[10].(int64)
			ingCnt = procRset.Row[11].(int64)
			passCnt = procRset.Row[12].(int64)
			failCnt = procRset.Row[13].(int64)

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

			recruitPostList = append(recruitPostList, models.RecruitPostList{
				TotCnt:      totCnt,
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
				Pagination:  pagination,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnRecruitPostList = models.RtnRecruitPostList{
			RtnRecruitPostListData: recruitPostList,
		}

		c.Data["json"] = &rtnRecruitPostList
		c.ServeJSON()
	}
}
