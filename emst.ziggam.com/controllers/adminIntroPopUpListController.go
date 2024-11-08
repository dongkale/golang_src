package controllers

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"emst.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type AdminIntroPopUpListController struct {
	BaseController
}

func (c *AdminIntroPopUpListController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/common/login")
	}

	imgServer, _ := beego.AppConfig.String("viewpath")

	pLang, _ := beego.AppConfig.String("lang")
	pSdy := c.GetString("sdy")
	pEdy := c.GetString("edy")

	nowDate := time.Now()
	dateFmt := fmt.Sprintf(nowDate.Format("20060102150405"))

	if pSdy == "" {
		pSdy = dateFmt[0:8]
	}

	if pEdy == "" {
		pEdy = dateFmt[0:8]
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

	// Start : Intro List
	log.Debug("CALL SP_EMS_ADMIN_INTRO_LIST_R('%v', %v, %v, '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pSdy, pEdy)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_INTRO_LIST_R('%v', %v, %v, '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pSdy, pEdy),
		ora.I64, /* TOT_CNT */
		ora.S,   /* INTRO_SN */
		ora.S,   /* LNK_GBN_CD */
		ora.S,   /* LNK_GBN_NM */
		ora.S,   /* LNK_GBN_VAL */
		ora.S,   /* BRD_GBN_CD */
		ora.I64, /* SN */
		ora.S,   /* PTO_PATH */
		ora.S,   /* REG_DT */
		ora.S,   /* USE_YN */
		ora.S,   /* LNK_GBN_VAL_NM */
		ora.S,   /* PUBL_SDY */
		ora.S,   /* PUBL_EDY */
		ora.S,   /* STB_YN */
		ora.S,   /* END_YN */
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

	adminIntroPopUpList := make([]models.AdminIntroPopUpList, 0)

	var (
		totCnt      int64
		introSn     string
		lnkGbnCd    string
		lnkGbnNm    string
		lnkGbnVal   string
		brdGbnCd    string
		sn          int64
		ptoPath     string
		regDt       string
		introTitle  string
		useYn       string
		lnkGbnValNm string
		sdy         string
		edy         string
		stbYn       string
		endYn       string
		fullPtoPath string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			totCnt = procRset.Row[0].(int64)
			introSn = procRset.Row[1].(string)
			lnkGbnCd = procRset.Row[2].(string)
			lnkGbnNm = procRset.Row[3].(string)
			lnkGbnVal = procRset.Row[4].(string)
			brdGbnCd = procRset.Row[5].(string)
			sn = procRset.Row[6].(int64)
			ptoPath = procRset.Row[7].(string)
			regDt = procRset.Row[8].(string)
			introTitle = procRset.Row[9].(string)
			useYn = procRset.Row[10].(string)
			lnkGbnValNm = procRset.Row[11].(string)
			sdy = procRset.Row[12].(string)
			edy = procRset.Row[13].(string)
			stbYn = procRset.Row[14].(string)
			endYn = procRset.Row[15].(string)

			if ptoPath == "" {
				fullPtoPath = ptoPath
			} else {
				fullPtoPath = imgServer + ptoPath
			}

			adminIntroPopUpList = append(adminIntroPopUpList, models.AdminIntroPopUpList{
				TotCnt:      totCnt,
				IntroSn:     introSn,
				LnkGbnCd:    lnkGbnCd,
				LnkGbnNm:    lnkGbnNm,
				LnkGbnVal:   lnkGbnVal,
				BrdGbnCd:    brdGbnCd,
				Sn:          sn,
				PtoPath:     fullPtoPath,
				RegDt:       regDt,
				IntroTitle:  introTitle,
				UseYn:       useYn,
				LnkGbnValNm: lnkGbnValNm,
				Sdy:         sdy,
				Edy:         edy,
				StbYn:       stbYn,
				EndYn:       endYn,
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

	// End : intro Popup List

	// Start : Intro Popup Status

	log.Debug("CALL SP_EMS_ADMIN_INTRO_STAT_R('%v', :1)",
		pLang)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_INTRO_STAT_R('%v', :1)",
		pLang),
		ora.S, /* USE_YN */
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

	adminIntroPopUpStat := make([]models.AdminIntroPopUpStat, 0)

	var (
		introUseYn string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			introUseYn = procRset.Row[0].(string)

			adminIntroPopUpStat = append(adminIntroPopUpStat, models.AdminIntroPopUpStat{
				IntroUseYn: introUseYn,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End :  Intro Popup Status

	c.Data["AdminIntroPopUpList"] = adminIntroPopUpList
	c.Data["Pagination"] = pagination
	c.Data["IntroUseYn"] = introUseYn
	c.Data["PageNo"] = pageNo
	c.Data["MenuId"] = "04"
	c.Data["SubMenuId"] = "06"
	c.TplName = "admin/intro_popup_list.html"
}
