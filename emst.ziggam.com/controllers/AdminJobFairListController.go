package controllers

import (
	"fmt"
	"math"
	"strconv"

	"emst.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type AdminJobFairListController struct {
	BaseController
}

func (c *AdminJobFairListController) Get() {

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
	pKndCd := c.GetString("knd_cd")
	if pKndCd == "" {
		pKndCd = "A00"
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

	// Start : Admin JobFair List
	log.Debug("CALL SP_EMS_ADMIN_EVENT_LIST_R('%v', %v, %v, '%v', :1)",
		pLang, pOffSet, pLimit, pKndCd)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_EVENT_LIST_R('%v', %v, %v, '%v', :1)",
		pLang, pOffSet, pLimit, pKndCd),
		ora.I64, /* TOT_CNT */
		ora.S,   /* REG_DT */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* MEM_ID */
		ora.S,   /* NM */
		ora.S,   /* SEX */
		ora.S,   /* EMAIL */
		ora.S,   /* BRTH_YMD */
		ora.I64, /* AGE */
		ora.S,   /* OS_GBN */
		ora.I64, /* VD_CNT */
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

	adminEventList := make([]models.AdminEventList, 0)

	var (
		totCnt  int64
		regDt   string
		ppMemNo string
		memId   string
		nm      string
		sex     string
		email   string
		brthYmd string
		age     int64
		osGbn   string
		vdCnt   int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			totCnt = procRset.Row[0].(int64)
			regDt = procRset.Row[1].(string)
			ppMemNo = procRset.Row[2].(string)
			memId = procRset.Row[3].(string)
			nm = procRset.Row[4].(string)
			sex = procRset.Row[5].(string)
			email = procRset.Row[6].(string)
			brthYmd = procRset.Row[7].(string)
			age = procRset.Row[8].(int64)
			osGbn = procRset.Row[9].(string)
			vdCnt = procRset.Row[10].(int64)

			adminEventList = append(adminEventList, models.AdminEventList{
				TotCnt:  totCnt,
				RegDt:   regDt,
				PpMemNo: ppMemNo,
				MemId:   memId,
				Nm:      nm,
				Sex:     sex,
				Email:   email,
				BrthYmd: brthYmd,
				Age:     age,
				OsGbn:   osGbn,
				VdCnt:   vdCnt,
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

	// End : Admin Event List

	c.Data["AdminEventList"] = adminEventList
	c.Data["Pagination"] = pagination
	c.Data["PageNo"] = pageNo
	c.Data["KndCd"] = pKndCd
	c.Data["TotCnt"] = totCnt
	c.Data["MenuId"] = "08"
	//c.Data["SubMenuId"] = "02"

	c.TplName = "admin/jobfair_list.html"
}

func (c *AdminJobFairListController) Post() {

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
	pKndCd := c.GetString("knd_cd")
	if pKndCd == "" {
		pKndCd = "A00"
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

	// Start : Admin Event List
	log.Debug("CALL SP_EMS_ADMIN_EVENT_LIST_R('%v', %v, %v, '%v', :1)",
		pLang, pOffSet, pLimit, pKndCd)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_EVENT_LIST_R('%v', %v, %v, '%v', :1)",
		pLang, pOffSet, pLimit, pKndCd),
		ora.I64, /* TOT_CNT */
		ora.S,   /* REG_DT */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* MEM_ID */
		ora.S,   /* NM */
		ora.S,   /* SEX */
		ora.S,   /* EMAIL */
		ora.S,   /* BRTH_YMD */
		ora.I64, /* AGE */
		ora.S,   /* OS_GBN */
		ora.I64, /* VD_CNT */
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

	rtnAdminEventList := models.RtnAdminEventList{}
	adminEventList := make([]models.AdminEventList, 0)

	var (
		totCnt  int64
		regDt   string
		ppMemNo string
		memId   string
		nm      string
		sex     string
		email   string
		brthYmd string
		age     int64
		osGbn   string
		vdCnt   int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			totCnt = procRset.Row[0].(int64)
			regDt = procRset.Row[1].(string)
			ppMemNo = procRset.Row[2].(string)
			memId = procRset.Row[3].(string)
			nm = procRset.Row[4].(string)
			sex = procRset.Row[5].(string)
			email = procRset.Row[6].(string)
			brthYmd = procRset.Row[7].(string)
			age = procRset.Row[8].(int64)
			osGbn = procRset.Row[9].(string)
			vdCnt = procRset.Row[10].(int64)

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

			adminEventList = append(adminEventList, models.AdminEventList{
				TotCnt:     totCnt,
				RegDt:      regDt,
				PpMemNo:    ppMemNo,
				MemId:      memId,
				Nm:         nm,
				Sex:        sex,
				Email:      email,
				BrthYmd:    brthYmd,
				Age:        age,
				OsGbn:      osGbn,
				VdCnt:      vdCnt,
				Pagination: pagination,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnAdminEventList = models.RtnAdminEventList{
			RtnAdminEventListData: adminEventList,
		}
		// End : Admin Event List

		c.Data["json"] = &rtnAdminEventList
		c.ServeJSON()
	}
}
