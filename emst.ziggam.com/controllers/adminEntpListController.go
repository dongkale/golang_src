package controllers

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"emst.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type AdminEntpListController struct {
	BaseController
}

func (c *AdminEntpListController) Get() {

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
	pGbnCd := c.GetString("gbn_cd")
	if pGbnCd == "" {
		pGbnCd = "A"
	}
	pKeyword := c.GetString("keyword") // 검색어
	pSdy := c.GetString("sdy")         // 검색시작일자
	pEdy := c.GetString("edy")         // 검색종료일자
	pVdYn := c.GetString("vd_yn")      //영상프로필 (전체:A, 있음:Y, 없음:N)
	pUseYn := c.GetString("use_yn")    //영상검증여부 (전체:A, 완료:1, 대기:0)
	pOsGbn := c.GetString("os_gbn")    //유입경로 (전체:A, 웹:WB, 안드로이드:AD, 아이폰:IS)

	pJobFair := c.GetString("jf_mng_cd")

	pLoginSdy := c.GetString("login_sdy") // 로그인 검색시작일자
	pLoginEdy := c.GetString("login_edy") // 로그인 검색종료일자

	if pSdy == "" {
		//pSdy = "20000101"
		pSdy = models.DefaultSdy
	}

	if pEdy == "" {
		//pEdy = "22001231"
		nowTime := time.Now()
		pEdy = nowTime.Format("20060102")
	}

	if pLoginSdy == "" {
		pLoginSdy = models.DefaultLoginSdy
	}

	if pLoginEdy == "" {
		//pLoginEdy = "22001231"
		nowTime := time.Now()
		pLoginEdy = nowTime.Format("20060102")
	}

	logs.Debug("pJobFair:" + pJobFair)

	logs.Debug("sdy:" + pSdy)
	logs.Debug("edy:" + pEdy)

	logs.Debug("login_sdy:" + pLoginSdy)
	logs.Debug("login_edy:" + pLoginEdy)

	if pVdYn == "" {
		pVdYn = "A"
	}

	if pUseYn == "" {
		pUseYn = "A"
	}

	if pOsGbn == "" {
		pOsGbn = "A"
	}

	// Parameter
	// 상세 정보에서 목록 버튼 눌렀을때 인자값 받는 부분
	pmPageNo := c.GetString("p_page_no")
	pmSize := c.GetString("p_size")
	pmGbnCd := c.GetString("p_gbn_cd")
	pmVdYn := c.GetString("p_vd_yn")
	pmUseYn := c.GetString("p_use_yn")
	pmSdy := c.GetString("p_sdy")
	pmEdy := c.GetString("p_edy")
	pmKeyword := c.GetString("p_keyword")
	pmOsGbn := c.GetString("p_os_gbn")

	pmLoginSdy := c.GetString("p_login_sdy")
	pmLoginEdy := c.GetString("p_login_edy")

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

	// Start : Admin Entp List
	logs.Debug(fmt.Sprintf("CALL SP_EMS_ADMIN_ENTP_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pGbnCd, pSdy, pEdy, pVdYn, pUseYn, pKeyword, pOsGbn, pJobFair, pLoginSdy, pLoginEdy))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_ENTP_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pGbnCd, pSdy, pEdy, pVdYn, pUseYn, pKeyword, pOsGbn, pJobFair, pLoginSdy, pLoginEdy),
		ora.I64, /* TOT_CNT */
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* MEM_STAT */
		ora.S,   /* MEM_STAT_CD */
		ora.S,   /* ENTP_MEM_ID */
		ora.S,   /* ENTP_KO_NM */
		ora.S,   /* REP_NM */
		ora.S,   /* REG_DT */
		ora.S,   /* MEM_STAT_CD */
		ora.I64, /* TOT_APLY_CNT */
		ora.I64, /* NEW_APLYCNT */
		ora.S,   /* BIZ_REG_NO */
		ora.S,   /* EMAIL */
		ora.S,   /* PP_CHRG_NM */
		ora.S,   /* PP_CHRG_TEL_NO */
		ora.S,   /* VD_YN */
		ora.S,   /* USE_YN */
		ora.S,   /* OS_GBN */
		ora.S,   /* JOBFAIR_CODE */
		ora.S,   /* LAST_LOGIN */
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

	adminEntpList := make([]models.AdminEntpList, 0)

	var (
		totCnt      int64
		entpMemNo   string
		memStat     string
		memStatDt   string
		entpMemId   string
		entpKoNm    string
		repNm       string
		regDt       string
		memStatCd   string
		totAplyCnt  int64
		newAplyCnt  int64
		bizRegNo    string
		email       string
		ppChrgNm    string
		ppChrgTelNo string
		vdYn        string
		useYn       string
		osGbn       string
		jobfairCds  string
		lastLogin   string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			totCnt = procRset.Row[0].(int64)
			entpMemNo = procRset.Row[1].(string)
			memStat = procRset.Row[2].(string)
			memStatDt = procRset.Row[3].(string)
			entpMemId = procRset.Row[4].(string)
			entpKoNm = procRset.Row[5].(string)
			repNm = procRset.Row[6].(string)
			regDt = procRset.Row[7].(string)
			memStatCd = procRset.Row[8].(string)
			totAplyCnt = procRset.Row[9].(int64)
			newAplyCnt = procRset.Row[10].(int64)
			bizRegNo = procRset.Row[11].(string)
			email = procRset.Row[12].(string)
			ppChrgNm = procRset.Row[13].(string)
			ppChrgTelNo = procRset.Row[14].(string)
			vdYn = procRset.Row[15].(string)
			useYn = procRset.Row[16].(string)
			osGbn = procRset.Row[17].(string)
			jobfairCds = procRset.Row[18].(string)
			lastLogin = procRset.Row[19].(string)

			adminEntpList = append(adminEntpList, models.AdminEntpList{
				TotCnt:      totCnt,
				EntpMemNo:   entpMemNo,
				//JobFairCds:    jobfairCds,
				JobFairCdsArr: strings.Split(jobfairCds, ","),
				MemStat:     memStat,
				MemStatDt:   memStatDt,
				EntpMemId:   entpMemId,
				EntpKoNm:    entpKoNm,
				RepNm:       repNm,
				RegDt:       regDt,
				MemStatCd:   memStatCd,
				TotAplyCnt:  totAplyCnt,
				NewAplyCnt:  newAplyCnt,
				BizRegNo:    bizRegNo,
				Email:       email,
				PpChrgNm:    ppChrgNm,
				PpChrgTelNo: ppChrgTelNo,
				VdYn:        vdYn,
				UseYn:       useYn,
				OsGbn:       osGbn,
				LastLogin:     lastLogin,
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

	// End : Admin Entp List

	// Start : Common YYYY List
	pCmGbnCd := "YY"

	log.Debug("CALL SP_EMS_DATE_LIST_R('%v', '%v', :1)",
		pLang, pCmGbnCd)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_DATE_LIST_R('%v', '%v', :1)",
		pLang, pCmGbnCd),
		ora.S, /* DATE_VAL */
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

	commonYYList := make([]models.CommonYYList, 0)

	var (
		yyyy string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			yyyy = procRset.Row[0].(string)

			commonYYList = append(commonYYList, models.CommonYYList{
				YYYY: yyyy,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Common YYYY List

	// Start : Common MM List
	pCmGbnCd = "MM"

	log.Debug("CALL SP_EMS_DATE_LIST_R('%v', '%v', :1)",
		pLang, pCmGbnCd)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_DATE_LIST_R('%v', '%v', :1)",
		pLang, pCmGbnCd),
		ora.S, /* DATE_VAL */
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

	commonMMList := make([]models.CommonMMList, 0)

	var (
		mm string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			mm = procRset.Row[0].(string)

			commonMMList = append(commonMMList, models.CommonMMList{
				MM: mm,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Common MM List

	// Start : Common DD List
	pCmGbnCd = "DD"

	log.Debug("CALL SP_EMS_DATE_LIST_R('%v', '%v', :1)",
		pLang, pCmGbnCd)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_DATE_LIST_R('%v', '%v', :1)",
		pLang, pCmGbnCd),
		ora.S, /* DATE_VAL */
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

	commonDDList := make([]models.CommonDDList, 0)

	var (
		dd string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			dd = procRset.Row[0].(string)

			commonDDList = append(commonDDList, models.CommonDDList{
				DD: dd,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Common DD List

	// Start : Admin Entp Top Info
	log.Debug("CALL SP_EMS_ADMIN_ENTP_TOP_INFO_R('%v', :1)",
		pLang)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_ENTP_TOP_INFO_R('%v', :1)",
		pLang),
		ora.I64, /* TOT_CNT */
		ora.I64, /* COM_CNT */
		ora.I64, /* STB_CNT */
		ora.I64, /* WTD_CNT */
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

	adminEntpTopInfo := make([]models.AdminEntpTopInfo, 0)

	var (
		eTotCnt int64
		eComCnt int64
		eStbCnt int64
		eWtdCnt int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			eTotCnt = procRset.Row[0].(int64)
			eComCnt = procRset.Row[1].(int64)
			eStbCnt = procRset.Row[2].(int64)
			eWtdCnt = procRset.Row[3].(int64)

			adminEntpTopInfo = append(adminEntpTopInfo, models.AdminEntpTopInfo{
				ETotCnt: eTotCnt,
				EComCnt: eComCnt,
				EStbCnt: eStbCnt,
				EWtdCnt: eWtdCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Admin Entp Top Info

	// Start : Jobfair List
	logs.Debug(fmt.Sprintf("CALL ZSP_ENTP_JOBFAIR_LIST('%v', '%v', :1)", pLang, ""))
	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_ENTP_JOBFAIR_LIST('%v', '%v', :1)",
		pLang, ""),
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

	jobFairList := make([]models.JobfairInfo, 0)

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

			jobFairList = append(jobFairList, models.JobfairInfo{
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

	c.Data["AdminEntpList"] = adminEntpList
	c.Data["Pagination"] = pagination
	c.Data["PageNo"] = pageNo
	c.Data["GbnCd"] = pGbnCd
	c.Data["TotCnt"] = totCnt

	c.Data["CommonYYList"] = commonYYList
	c.Data["CommonMMList"] = commonMMList
	c.Data["CommonDDList"] = commonDDList

	c.Data["ETotCnt"] = eTotCnt
	c.Data["EComCnt"] = eComCnt
	c.Data["EStbCnt"] = eStbCnt
	c.Data["EWtdCnt"] = eWtdCnt
	c.Data["MenuId"] = "06"

	c.Data["JobFairList"] = jobFairList

	// 웹페이지 출력시 초기값, html 이랑 값을 맞추기 위해 -->
	c.Data["Sdy"] = pSdy
	c.Data["Edy"] = pEdy

	c.Data["LoginSdy"] = pLoginSdy
	c.Data["LoginEdy"] = pLoginEdy
	// <--

	// 웹페이지 출력시 입력값 -->
	/* Parameter Value */
	c.Data["pPageNo"] = pmPageNo
	c.Data["pSize"] = pmSize
	c.Data["pGbnCd"] = pmGbnCd
	c.Data["pVdYn"] = pmVdYn
	c.Data["pUseYn"] = pmUseYn
	c.Data["pSdy"] = pmSdy
	c.Data["pEdy"] = pmEdy
	c.Data["pKeyword"] = pmKeyword
	c.Data["pPOsGbn"] = pmOsGbn

	c.Data["pLoginSdy"] = pmLoginSdy
	c.Data["pLoginEdy"] = pmLoginEdy
	// <--

	c.Data["BizUrl"], _ = beego.AppConfig.String("bizUrl")

	runmode, _ := beego.AppConfig.String("runmode")
	c.Data["runmode"] = runmode

	c.TplName = "admin/entp_list.html"
}

func (c *AdminEntpListController) Post() {

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
	pGbnCd := c.GetString("gbn_cd")
	if pGbnCd == "" {
		pGbnCd = "A"
	}
	pKeyword := c.GetString("keyword") // 검색어
	pSdy := c.GetString("sdy")         // 검색시작일자
	pEdy := c.GetString("edy")         // 검색종료일자
	pVdYn := c.GetString("vd_yn")      //영상프로필 (전체:A, 있음:Y, 없음:N)
	pUseYn := c.GetString("use_yn")    //영상검증여부 (전체:A, 완료:1, 대기:0)
	pOsGbn := c.GetString("os_gbn")    //유입경로 (전체:A, 웹:WB, 안드로이드:AD, 아이폰:IS)

	pJobFair := c.GetString("jf_mng_cd")

	pLoginSdy := c.GetString("login_sdy") // 로그인 검색시작일자
	pLoginEdy := c.GetString("login_edy") // 로그인 검색종료일자

	if pSdy == "" {
		//pSdy = "20000101"
		pSdy = models.DefaultSdy
	}

	if pEdy == "" {
		//pEdy = "22001231"
		nowTime := time.Now()
		pEdy = nowTime.Format("20060102")
	}

	if pLoginSdy == "" {
		pLoginSdy = models.DefaultLoginSdy
	}

	if pLoginEdy == "" {
		//pLoginEdy = "22001231"
		nowTime := time.Now()
		pLoginEdy = nowTime.Format("20060102")
	}

	logs.Debug("sdy:" + pSdy)
	logs.Debug("edy:" + pEdy)

	logs.Debug("login_sdy:" + pLoginSdy)
	logs.Debug("login_edy:" + pLoginEdy)

	if pVdYn == "" {
		pVdYn = "A"
	}

	if pUseYn == "" {
		pUseYn = "A"
	}

	if pOsGbn == "" {
		pOsGbn = "A"
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

	// Start : Admin Entp List
	logs.Debug(fmt.Sprintf("CALL SP_EMS_ADMIN_ENTP_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pGbnCd, pSdy, pEdy, pVdYn, pUseYn, pKeyword, pOsGbn, pJobFair, pLoginSdy, pLoginEdy))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_ENTP_LIST_R('%v', %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pGbnCd, pSdy, pEdy, pVdYn, pUseYn, pKeyword, pOsGbn, pJobFair, pLoginSdy, pLoginEdy),
		ora.I64, /* TOT_CNT */
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* MEM_STAT */
		ora.S,   /* MEM_STAT_CD */
		ora.S,   /* ENTP_MEM_ID */
		ora.S,   /* ENTP_KO_NM */
		ora.S,   /* REP_NM */
		ora.S,   /* REG_DT */
		ora.S,   /* MEM_STAT_CD */
		ora.I64, /* TOT_APLY_CNT */
		ora.I64, /* NEW_APLYCNT */
		ora.S,   /* BIZ_REG_NO */
		ora.S,   /* EMAIL */
		ora.S,   /* PP_CHRG_NM */
		ora.S,   /* PP_CHRG_TEL_NO */
		ora.S,   /* VD_YN */
		ora.S,   /* USE_YN */
		ora.S,   /* OS_GBN */
		ora.S,   /* JOBFAIR_CODE */
		ora.S,   /* LAST_LOGIN */
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

	rtnAdminEntpList := models.RtnAdminEntpList{}
	adminEntpList := make([]models.AdminEntpList, 0)

	var (
		totCnt      int64
		entpMemNo   string
		memStat     string
		memStatDt   string
		entpMemId   string
		entpKoNm    string
		repNm       string
		regDt       string
		memStatCd   string
		totAplyCnt  int64
		newAplyCnt  int64
		bizRegNo    string
		email       string
		ppChrgNm    string
		ppChrgTelNo string
		vdYn        string
		useYn       string
		osGbn       string
		jobfairCds  string
		lastLogin   string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			totCnt = procRset.Row[0].(int64)
			entpMemNo = procRset.Row[1].(string)
			memStat = procRset.Row[2].(string)
			memStatDt = procRset.Row[3].(string)
			entpMemId = procRset.Row[4].(string)
			entpKoNm = procRset.Row[5].(string)
			repNm = procRset.Row[6].(string)
			regDt = procRset.Row[7].(string)
			memStatCd = procRset.Row[8].(string)
			totAplyCnt = procRset.Row[9].(int64)
			newAplyCnt = procRset.Row[10].(int64)
			bizRegNo = procRset.Row[11].(string)
			email = procRset.Row[12].(string)
			ppChrgNm = procRset.Row[13].(string)
			ppChrgTelNo = procRset.Row[14].(string)
			vdYn = procRset.Row[15].(string)
			useYn = procRset.Row[16].(string)
			osGbn = procRset.Row[17].(string)
			jobfairCds = procRset.Row[18].(string)
			lastLogin = procRset.Row[19].(string)

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

			adminEntpList = append(adminEntpList, models.AdminEntpList{
				TotCnt:      totCnt,
				EntpMemNo:   entpMemNo,
				//JobFairCds:    jobfairCds,
				JobFairCdsArr: strings.Split(jobfairCds, ","),
				MemStat:     memStat,
				MemStatDt:   memStatDt,
				EntpMemId:   entpMemId,
				EntpKoNm:    entpKoNm,
				RepNm:       repNm,
				RegDt:       regDt,
				MemStatCd:   memStatCd,
				TotAplyCnt:  totAplyCnt,
				NewAplyCnt:  newAplyCnt,
				BizRegNo:    bizRegNo,
				Email:       email,
				PpChrgNm:    ppChrgNm,
				PpChrgTelNo: ppChrgTelNo,
				VdYn:        vdYn,
				UseYn:       useYn,
				OsGbn:       osGbn,
				LastLogin:     lastLogin,
				Pagination:  pagination,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnAdminEntpList = models.RtnAdminEntpList{
			RtnAdminEntpListData: adminEntpList,
		}
		// End : Admin Entp List

		//logs.Debug(rtnAdminEntpList)

		c.Data["json"] = &rtnAdminEntpList
		c.ServeJSON()
	}
}
