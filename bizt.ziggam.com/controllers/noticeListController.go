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

type NoticeListController struct {
	BaseController
}

func (c *NoticeListController) Get() {

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

	// Start : Notice List
	fmt.Printf(fmt.Sprintf("CALL ZSP_NOTICE_LIST_R('%v', %v, %v, :1)",
		pLang, pOffSet, pLimit))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_NOTICE_LIST_R('%v', %v, %v, :1)",
		pLang, pOffSet, pLimit),
		ora.I64, /* TOT_CNT */
		ora.I64, /* SN */
		ora.S,   /* GBN_NM */
		ora.S,   /* TITLE */
		ora.S,   /* REG_DT */
		ora.S,   /* NEW_YN */
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

	noticeList := make([]models.NoticeList, 0)

	var (
		totCnt int64
		sn     int64
		gbnNm  string
		title  string
		regDt  string
		newYn  string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			totCnt = procRset.Row[0].(int64)
			sn = procRset.Row[1].(int64)
			gbnNm = procRset.Row[2].(string)
			title = procRset.Row[3].(string)
			regDt = procRset.Row[4].(string)
			newYn = procRset.Row[5].(string)

			noticeList = append(noticeList, models.NoticeList{
				TotCnt: totCnt,
				Sn:     sn,
				GbnNm:  gbnNm,
				Title:  title,
				RegDt:  regDt,
				NewYn:  newYn,
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

	// End : Notice List

	c.Data["NoticeList"] = noticeList
	c.Data["Pagination"] = pagination
	c.Data["PageNo"] = pageNo

	c.Data["TMenuId"] = "M00"
	c.Data["SMenuId"] = "M00"

	c.TplName = "setting/notice_list.html"
}

func (c *NoticeListController) Post() {
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

	// Start : Notice List
	fmt.Printf(fmt.Sprintf("CALL ZSP_NOTICE_LIST_R('%v', %v, %v, :1)",
		pLang, pOffSet, pLimit))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_NOTICE_LIST_R('%v', %v, %v, :1)",
		pLang, pOffSet, pLimit),
		ora.I64, /* TOT_CNT */
		ora.I64, /* SN */
		ora.S,   /* GBN_NM */
		ora.S,   /* TITLE */
		ora.S,   /* REG_DT */
		ora.S,   /* NEW_YN */
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

	rtnNoticeList := models.RtnNoticeList{}
	noticeList := make([]models.NoticeList, 0)

	var (
		totCnt int64
		sn     int64
		gbnNm  string
		title  string
		regDt  string
		newYn  string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			totCnt = procRset.Row[0].(int64)
			sn = procRset.Row[1].(int64)
			gbnNm = procRset.Row[2].(string)
			title = procRset.Row[3].(string)
			regDt = procRset.Row[4].(string)
			newYn = procRset.Row[5].(string)

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

			noticeList = append(noticeList, models.NoticeList{
				TotCnt:     totCnt,
				Sn:         sn,
				GbnNm:      gbnNm,
				Title:      title,
				RegDt:      regDt,
				NewYn:      newYn,
				Pagination: pagination,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnNoticeList = models.RtnNoticeList{
			RtnNoticeListData: noticeList,
		}
		// End : Notice List

		c.Data["json"] = &rtnNoticeList
		c.ServeJSON()
	}

}
