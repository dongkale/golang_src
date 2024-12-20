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

type InquiryListController struct {
	BaseController
}

func (c *InquiryListController) Get() {

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
	//pEntpMemNo := c.GetString("entp_mem_no")
	pEntpMemNo := mem_no //"E2018102500001"
	//E2018102500001

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Notice List
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

	// Start : Inquiry List
	logs.Debug("CALL SP_EMS_INQUIRY_LIST_R('%v', '%v', %v, %v, :1)",
		pLang, pEntpMemNo, pOffSet, pLimit)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_INQUIRY_LIST_R('%v', '%v', %v, %v, :1)",
		pLang, pEntpMemNo, pOffSet, pLimit),
		ora.I64, /* TOT_CNT */
		ora.I64, /* BRD_NO */
		ora.S,   /* INQ_GBN_NM */
		ora.S,   /* REG_DY */
		ora.S,   /* INQ_TITLE */
		ora.S,   /* INQ_CONT */
		ora.S,   /* EMAIL */
		ora.S,   /* ANS_YN */
		ora.S,   /* ANS_YN_NM */
		ora.S,   /* ANS_DT */
		ora.S,   /* ANS_CONT */
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

	inquiryList := make([]models.InquiryList, 0)

	var (
		totCnt   int64
		brdNo    int64
		inqGbnNm string
		regDy    string
		inqTitle string
		inqCont  string
		email    string
		ansYn    string
		ansYnNm  string
		ansDt    string
		ansCont  string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			totCnt = procRset.Row[0].(int64)
			brdNo = procRset.Row[1].(int64)
			inqGbnNm = procRset.Row[2].(string)
			regDy = procRset.Row[3].(string)
			inqTitle = procRset.Row[4].(string)
			inqCont = procRset.Row[5].(string)
			email = procRset.Row[6].(string)
			ansYn = procRset.Row[7].(string)
			ansYnNm = procRset.Row[8].(string)
			ansDt = procRset.Row[9].(string)
			ansCont = procRset.Row[10].(string)

			inquiryList = append(inquiryList, models.InquiryList{
				TotCnt:   totCnt,
				BrdNo:    brdNo,
				InqGbnNm: inqGbnNm,
				RegDy:    regDy,
				InqTitle: inqTitle,
				InqCont:  inqCont,
				Email:    email,
				AnsYn:    ansYn,
				AnsYnNm:  ansYnNm,
				AnsDt:    ansDt,
				AnsCont:  ansCont,
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
				if i == startPageNo {
					if i == endPageNo {
						pagination += "<a href='javascript:void(0);' class='active num goPage' id='num' data-page='" + strconv.Itoa(int(i)) + "'>" + strconv.Itoa(int(i)) + "</a>"
					} else {
						pagination += "<a href='javascript:void(0);' class='active num goPage' id='num' data-page='" + strconv.Itoa(int(i)) + "'>" + strconv.Itoa(int(i)) + "</a>" + " | "
					}
				} else {
					if i == endPageNo {
						pagination += "<a href='javascript:void(0);' class='active num goPage' id='num' data-page='" + strconv.Itoa(int(i)) + "'>" + strconv.Itoa(int(i)) + "</a>"
					} else {
						pagination += "<a href='javascript:void(0);' class='active num goPage' id='num' data-page='" + strconv.Itoa(int(i)) + "'>" + strconv.Itoa(int(i)) + "</a>" + " | "
					}
				}
			} else {
				if i == startPageNo {
					if i == endPageNo {
						pagination += "<a href='javascript:void(0);' class='num goPage' id='num' data-page='" + strconv.Itoa(int(i)) + "'>" + strconv.Itoa(int(i)) + "</a>"
					} else {
						pagination += "<a href='javascript:void(0);' class='num goPage' id='num' data-page='" + strconv.Itoa(int(i)) + "'>" + strconv.Itoa(int(i)) + "</a>" + " | "
					}
				} else {
					if i == endPageNo {
						pagination += "<a href='javascript:void(0);' class='num goPage' id='num' data-page='" + strconv.Itoa(int(i)) + "'>" + strconv.Itoa(int(i)) + "</a>"
					} else {
						pagination += "<a href='javascript:void(0);' class='num goPage' id='num' data-page='" + strconv.Itoa(int(i)) + "'>" + strconv.Itoa(int(i)) + "</a>" + " | "
					}
				}
			}
		}
		if nextPageNo == pageNo {
			pagination += "<a href='javascript:void(0);' class='btnNext disabled' id='prev' data-page='" + strconv.Itoa(int(nextPageNo)) + "'>다음</a>"
		} else {
			pagination += "<a href='javascript:void(0);' class='btnNext goPage' id='prev' data-page='" + strconv.Itoa(int(nextPageNo)) + "'>다음</a>"
		}
	}

	// End : Inquiry List

	c.Data["InquiryList"] = inquiryList
	c.Data["Pagination"] = pagination
	c.Data["PageNo"] = pageNo

	c.TplName = "inquiry/inquiry_list.html"
}
