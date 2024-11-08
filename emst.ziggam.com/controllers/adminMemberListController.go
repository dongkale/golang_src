package controllers

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"emst.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	ora "gopkg.in/rana/ora.v4"
)

type AdminMemberListController struct {
	BaseController
}

func (c *AdminMemberListController) Get() {
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
	pSex := c.GetString("sex")                      //성별 (무관:A, 남성:M, 여성:F)
	pAge := c.GetString("age")                      //연령 (전체:00, 19세이하:01, 20~29:02, 30~39:03, 40~49:04, 50~59:05, 60세이상:06)
	pVpYn := c.GetString("vp_yn")                   //영상프로필 (전체:9, 있음:1, 없음:0)
	pKeyword := c.GetString("keyword")              // 검색어
	pSdy := c.GetString("sdy")                      // 검색시작일자
	pEdy := c.GetString("edy")                      // 검색종료일자
	pMemStat := c.GetString("mem_stat")             // 회원상태
	pOsGbn := c.GetString("os_gbn")                 //OS구분 (전체:A, 안드로이드:AD, 애플:IS)
	pMemJoinGbnCd := c.GetString("mem_join_gbn_cd") //회원가입유형 (전체:A, 일반:00, 페이스북:01, 카카오:02)
	pJobFair := c.GetString("jf_mng_cd")

	pLoginSdy := c.GetString("login_sdy") // 로그인 검색시작일자
	pLoginEdy := c.GetString("login_edy") // 로그인 검색종료일자

	/* Parameter */
	pmMemStat := c.GetString("p_mem_stat")
	pmSex := c.GetString("p_sex")
	pmAge := c.GetString("p_age")
	pmVpYn := c.GetString("p_vp_yn")
	pmOsGbn := c.GetString("p_os_gbn")
	pmMemJoinGbnCd := c.GetString("p_mem_join_gbn_cd")
	pmSdy := c.GetString("p_sdy")
	pmEdy := c.GetString("p_edy")
	pmPageNo := c.GetString("p_page_no")
	pmKeyword := c.GetString("p_keyword")
	pmSize := c.GetString("p_size")

	pmLoginSdy := c.GetString("p_login_sdy")
	pmLoginEdy := c.GetString("p_login_edy")

	if pMemStat == "" {
		pMemStat = "00"
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

	if pOsGbn == "" {
		pOsGbn = "A"
	}

	if pMemJoinGbnCd == "" {
		pMemJoinGbnCd = "A"
	}

	if pSdy == "" {
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

	// Start : Admin Member List
	logs.Debug(fmt.Sprintf("CALL SP_EMS_ADMIN_MEMBER_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pSex, pAge, pVpYn, pSdy, pEdy, pKeyword, pMemStat, pOsGbn, pMemJoinGbnCd, pJobFair, pLoginSdy, pLoginEdy))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_MEMBER_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pSex, pAge, pVpYn, pSdy, pEdy, pKeyword, pMemStat, pOsGbn, pMemJoinGbnCd, pJobFair, pLoginSdy, pLoginEdy),
		ora.I64, /* TOT_CNT */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* MEM_STAT_CD */
		ora.S,   /* MEM_STAT_NM */
		ora.S,   /* MEM_ID */
		ora.S,   /* NM */
		ora.S,   /* SEX */
		ora.S,   /* EMAIL */
		ora.S,   /* BRTH_YMD */
		ora.I64, /* AGE */
		ora.S,   /* OS_GBN */
		ora.S,   /* M_REG_PRGS_STAT_CD */
		ora.S,   /* VP_YN */
		ora.S,   /* REG_DT */
		ora.S,   /* REG_DY */
		ora.S,   /* ARR_VD_PATH */
		ora.S,   /* MEM_JOIN_GBN_CD */
		ora.S,   /* MEM_JOIN_GBN_NM */
		ora.S,   /* MO_NO */
		ora.I64,   /* AH_TOT_CNT */
		ora.S,   /* JOBFAIR_MNG_CDS */
		ora.S,   /* LOGIN_DT */
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

	adminMemberList := make([]models.AdminMemberList, 0)

	var (
		totCnt         int64
		ppMemNo        string
		memStatCd      string
		memStatNm      string
		memId          string
		nm             string
		sex            string
		email          string
		brthYmd        string
		age            int64
		osGbn          string
		mregPrgsStatCd string
		vpYn           string
		regDt          string
		regDy          string
		arrVdPath      string
		memJoinGbnCd   string
		memJoinGbnNm   string
		moNo		   string
		ahTotCnt	   int64
		jobfairCds     string
		loginDt        string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			totCnt = procRset.Row[0].(int64)
			ppMemNo = procRset.Row[1].(string)
			memStatCd = procRset.Row[2].(string)
			memStatNm = procRset.Row[3].(string)
			memId = procRset.Row[4].(string)
			nm = procRset.Row[5].(string)
			sex = procRset.Row[6].(string)
			email = procRset.Row[7].(string)
			brthYmd = procRset.Row[8].(string)
			age = procRset.Row[9].(int64)
			osGbn = procRset.Row[10].(string)
			mregPrgsStatCd = procRset.Row[11].(string)
			vpYn = procRset.Row[12].(string)
			regDt = procRset.Row[13].(string)
			regDy = procRset.Row[14].(string)
			arrVdPath = procRset.Row[15].(string)
			memJoinGbnCd = procRset.Row[16].(string)
			memJoinGbnNm = procRset.Row[17].(string)
			moNo = procRset.Row[18].(string)
			ahTotCnt = procRset.Row[19].(int64)
			jobfairCds = procRset.Row[20].(string)
			loginDt = procRset.Row[21].(string)

			adminMemberList = append(adminMemberList, models.AdminMemberList{
				TotCnt:         totCnt,
				PpMemNo:        ppMemNo,
				MemStatCd:      memStatCd,
				MemStatNm:      memStatNm,
				MemId:          memId,
				Nm:             nm,
				Sex:            sex,
				Email:          email,
				BrthYmd:        brthYmd,
				Age:            age,
				OsGbn:          osGbn,
				MregPrgsStatCd: mregPrgsStatCd,
				VpYn:           vpYn,
				RegDt:          regDt,
				RegDy:          regDy,
				ArrVdPath:      arrVdPath,
				MemJoinGbnCd:   memJoinGbnCd,
				MemJoinGbnNm:   memJoinGbnNm,
				MoNo:   		moNo,
				AhTotCnt: 		ahTotCnt,
				JobFairCdsArr:  strings.Split(jobfairCds, ","),
				LoginDt:        loginDt,
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

	// End : Admin Member List

	// Start : Common YYYY List
	pGbnCd := "YY"

	log.Debug("CALL SP_EMS_DATE_LIST_R('%v', '%v', :1)",
		pLang, pGbnCd)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_DATE_LIST_R('%v', '%v', :1)",
		pLang, pGbnCd),
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
	pGbnCd = "MM"

	log.Debug("CALL SP_EMS_DATE_LIST_R('%v', '%v', :1)",
		pLang, pGbnCd)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_DATE_LIST_R('%v', '%v', :1)",
		pLang, pGbnCd),
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
	pGbnCd = "DD"

	log.Debug("CALL SP_EMS_DATE_LIST_R('%v', '%v', :1)",
		pLang, pGbnCd)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_DATE_LIST_R('%v', '%v', :1)",
		pLang, pGbnCd),
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

	// Start : Member Top Info

	log.Debug("CALL SP_EMS_MEMBER_TOP_INFO_R('%v', :1)",
		pLang)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_MEMBER_TOP_INFO_R('%v', :1)",
		pLang),
		ora.I64, /* TOT_MEM_CNT */
		ora.I64, /* RUN_MEM_CNT */
		ora.I64, /* WTD_MEM_CNT */
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

	memberTopInfo := make([]models.MemberTopInfo, 0)

	var (
		totMemCnt int64
		runMemCnt int64
		wtdMemCnt int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			totMemCnt = procRset.Row[0].(int64)
			runMemCnt = procRset.Row[1].(int64)
			wtdMemCnt = procRset.Row[2].(int64)

			memberTopInfo = append(memberTopInfo, models.MemberTopInfo{
				TotMemCnt: totMemCnt,
				RunMemCnt: runMemCnt,
				WtdMemCnt: wtdMemCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Member Top Info

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

	c.Data["CommonYYList"] = commonYYList
	c.Data["CommonMMList"] = commonMMList
	c.Data["CommonDDList"] = commonDDList

	c.Data["TotMemCnt"] = totMemCnt
	c.Data["RunMemCnt"] = runMemCnt
	c.Data["WtdMemCnt"] = wtdMemCnt

	c.Data["MenuId"] = "05"
	c.Data["MemStat"] = pMemStat
	c.Data["TotCnt"] = totCnt
	c.Data["Pagination"] = pagination
	c.Data["PageNo"] = pageNo
	c.Data["AdminMemberList"] = adminMemberList

	c.Data["JobFairList"] = jobFairList

	// 웹페이지 출력시 초기값, html 이랑 값을 맞추기 위해 -->
	c.Data["Sdy"] = pSdy
	c.Data["Edy"] = pEdy

	c.Data["LoginSdy"] = pLoginSdy
	c.Data["LoginEdy"] = pLoginEdy
	// <--

	// 웹페이지 출력시 입력값 -->
	/* Parameter Value */
	c.Data["pMemStat"] = pmMemStat
	c.Data["pSex"] = pmSex
	c.Data["pAge"] = pmAge
	c.Data["pVpYn"] = pmVpYn
	c.Data["pOsGbn"] = pmOsGbn
	c.Data["pMemJoinGbnCd"] = pmMemJoinGbnCd
	c.Data["pSdy"] = pmSdy
	c.Data["pEdy"] = pmEdy
	c.Data["pPageNo"] = pmPageNo
	c.Data["pKeywords"] = pmKeyword
	c.Data["pSize"] = pmSize

	c.Data["pLoginSdy"] = pmLoginSdy
	c.Data["pLoginEdy"] = pmLoginEdy
	// <--

	c.TplName = "admin/member_list.html"
}

func (c *AdminMemberListController) Post() {

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
	pSex := c.GetString("sex")                      //성별 (무관:A, 남성:M, 여성:F)
	pAge := c.GetString("age")                      //연령 (전체:00, 19세이하:01, 20~29:02, 30~39:03, 40~49:04, 50~59:05, 60세이상:06)
	pVpYn := c.GetString("vp_yn")                   //영상프로필 (전체:9, 있음:1, 없음:0)
	pKeyword := c.GetString("keyword")              // 검색어
	pSdy := c.GetString("sdy")                      // 검색시작일자
	pEdy := c.GetString("edy")                      // 검색종료일자
	pMemStat := c.GetString("mem_stat")             // 회원상태
	pOsGbn := c.GetString("os_gbn")                 //OS구분 (전체:A, 안드로이드:AD, 애플:IS)
	pMemJoinGbnCd := c.GetString("mem_join_gbn_cd") //회원가입유형 (전체:A, 일반:00, 페이스북:01, 카카오:02)
	pJobFair := c.GetString("jf_mng_cd")

	pLoginSdy := c.GetString("login_sdy") // 로그인 검색시작일자
	pLoginEdy := c.GetString("login_edy") // 로그인 검색종료일자

	if pMemStat == "" {
		pMemStat = "00"
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

	if pOsGbn == "" {
		pOsGbn = "A"
	}

	if pMemJoinGbnCd == "" {
		pMemJoinGbnCd = "A"
	}

	if pSdy == "" {
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

	// Start : Admin Member List
	logs.Debug(fmt.Sprintf("CALL SP_EMS_ADMIN_MEMBER_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pSex, pAge, pVpYn, pSdy, pEdy, pKeyword, pMemStat, pOsGbn, pMemJoinGbnCd, pJobFair, pLoginSdy, pLoginEdy))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_MEMBER_LIST_R('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pSex, pAge, pVpYn, pSdy, pEdy, pKeyword, pMemStat, pOsGbn, pMemJoinGbnCd, pJobFair, pLoginSdy, pLoginEdy),
		ora.I64, /* TOT_CNT */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* MEM_STAT_CD */
		ora.S,   /* MEM_STAT_NM */
		ora.S,   /* MEM_ID */
		ora.S,   /* NM */
		ora.S,   /* SEX */
		ora.S,   /* EMAIL */
		ora.S,   /* BRTH_YMD */
		ora.I64, /* AGE */
		ora.S,   /* OS_GBN */
		ora.S,   /* M_REG_PRGS_STAT_CD */
		ora.S,   /* VP_YN */
		ora.S,   /* REG_DT */
		ora.S,   /* REG_DY */
		ora.S,   /* ARR_VD_PATH */
		ora.S,   /* MEM_JOIN_GBN_CD */
		ora.S,   /* MEM_JOIN_GBN_NM */
		ora.S,   /* MO_NO */
		ora.I64,   /* AH_TOT_CNT */
		ora.S,   /* JOBFAIR_MNG_CDS */
		ora.S,   /* LOGIN_DT */
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

	rtnAdminMemberList := models.RtnAdminMemberList{}
	adminMemberList := make([]models.AdminMemberList, 0)

	var (
		totCnt         int64
		ppMemNo        string
		memStatCd      string
		memStatNm      string
		memId          string
		nm             string
		sex            string
		email          string
		brthYmd        string
		age            int64
		osGbn          string
		mregPrgsStatCd string
		vpYn           string
		regDt          string
		regDy          string
		arrVdPath      string
		memJoinGbnCd   string
		memJoinGbnNm   string
		moNo		   string
		ahTotCnt	   int64
		jobfairCds     string
		loginDt        string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			totCnt = procRset.Row[0].(int64)
			ppMemNo = procRset.Row[1].(string)
			memStatCd = procRset.Row[2].(string)
			memStatNm = procRset.Row[3].(string)
			memId = procRset.Row[4].(string)
			nm = procRset.Row[5].(string)
			sex = procRset.Row[6].(string)
			email = procRset.Row[7].(string)
			brthYmd = procRset.Row[8].(string)
			age = procRset.Row[9].(int64)
			osGbn = procRset.Row[10].(string)
			mregPrgsStatCd = procRset.Row[11].(string)
			vpYn = procRset.Row[12].(string)
			regDt = procRset.Row[13].(string)
			regDy = procRset.Row[14].(string)
			arrVdPath = procRset.Row[15].(string)
			memJoinGbnCd = procRset.Row[16].(string)
			memJoinGbnNm = procRset.Row[17].(string)
			moNo = procRset.Row[18].(string)
			ahTotCnt = procRset.Row[19].(int64)
			jobfairCds = procRset.Row[20].(string)
			loginDt = procRset.Row[21].(string)

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

			adminMemberList = append(adminMemberList, models.AdminMemberList{
				TotCnt:         totCnt,
				PpMemNo:        ppMemNo,
				MemStatCd:      memStatCd,
				MemStatNm:      memStatNm,
				MemId:          memId,
				Nm:             nm,
				Sex:            sex,
				Email:          email,
				BrthYmd:        brthYmd,
				Age:            age,
				OsGbn:          osGbn,
				MregPrgsStatCd: mregPrgsStatCd,
				VpYn:           vpYn,
				RegDt:          regDt,
				RegDy:          regDy,
				ArrVdPath:      arrVdPath,
				MemJoinGbnCd:   memJoinGbnCd,
				MemJoinGbnNm:   memJoinGbnNm,
				MoNo:			moNo,
				AhTotCnt:		ahTotCnt,
				JobFairCdsArr:  strings.Split(jobfairCds, ","),
				LoginDt:        loginDt,
				Pagination:     pagination,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnAdminMemberList = models.RtnAdminMemberList{
			RtnAdminMemberListData: adminMemberList,
		}
		// End : Applicant List

		c.Data["json"] = &rtnAdminMemberList
		c.ServeJSON()
	}
}
