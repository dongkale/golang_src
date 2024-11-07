package controllers

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"bizt.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

// InviteSendListController ...
type InviteSendListController struct {
	BaseController
}

// Get ...
func (c *InviteSendListController) Get() {

	session := c.StartSession()

	memNO := session.Get(c.Ctx.Request.Context(), "mem_no")
	if memNO == nil {
		c.Ctx.Redirect(302, "/login")
	}

	memSn := session.Get(c.Ctx.Request.Context(), "mem_sn")
	if memSn == nil {
		c.Ctx.Redirect(302, "/login")
	}

	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := memNO
	pPpChrgSn := memSn
	//imgServer, _  := beego.AppConfig.String("viewpath")

	fmt.Printf(pPpChrgSn.(string))

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : invite recruit List
	pGbnCd := c.GetString("gbn_cd") // 구분코드(A:전체, I:채용중, E:종료)

	if pGbnCd == "" {
		pGbnCd = "00"
	}

	pOffSetR := "0"
	pLimitR := "100"
	pSortGbnR := "03"

	fmt.Printf(fmt.Sprintf("CALL ZSP_INVITE_RECRUIT_LIST_TYPE('%v', %v, %v, '%v', '%v', '%v', :1)",
		pLang, pOffSetR, pLimitR, pEntpMemNo, pGbnCd, pSortGbnR))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_INVITE_RECRUIT_LIST_TYPE('%v', %v, %v, '%v', '%v', '%v', :1)",
		pLang, pOffSetR, pLimitR, pEntpMemNo, pGbnCd, pSortGbnR),
		ora.I64, /* TOT_CNT */
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* PRGS_STAT */
		ora.S,   /* RECRUT_TITLE */
		ora.S,   /* UP_JOB_GRP */
		ora.S,   /* JOB_GRP */
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

	recruitSubList := make([]models.RecruitSubList, 0)

	var (
		sTotCnt      int64
		sEntpMemNo   string
		sRecrutSn    string
		sPrgsStat    string
		sRecrutTitle string
		sUpJobGrp    string
		sJobGrp      string
		//sTrimRecrutTitle string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			sTotCnt = procRset.Row[0].(int64)
			sEntpMemNo = procRset.Row[1].(string)
			sRecrutSn = procRset.Row[2].(string)
			sPrgsStat = procRset.Row[3].(string)
			sRecrutTitle = procRset.Row[4].(string)
			sUpJobGrp = procRset.Row[5].(string)
			sJobGrp = procRset.Row[6].(string)

			recruitSubList = append(recruitSubList, models.RecruitSubList{
				STotCnt:      sTotCnt,
				SEntpMemNo:   sEntpMemNo,
				SRecrutSn:    sRecrutSn,
				SPrgsStat:    sPrgsStat,
				SRecrutTitle: sRecrutTitle,
				SUpJobGrp:    sUpJobGrp,
				SJobGrp:      sJobGrp,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}

	// End : invite recruit List

	// Start : Recruit Apply List
	pRecrutSn := c.GetString("recrut_sn") // 채용일련번호
	pSdy := c.GetString("sdy")
	pEdy := c.GetString("edy")

	if pSdy == "" {
		pSdy = "20200101"
	}

	if pEdy == "" {
		pEdy = "20201201"
		nowTime := time.Now()
		pEdy = nowTime.Format("20060102")
	}

	// if pRecrutSn == "" {
	// 	pRecrutSn = "A"
	// }

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

	pageSize, err = strconv.ParseInt(pPageSize, 10, 64)
	if err != nil {
		//
	}
	pOffSet := (pageNo - 1) * pageSize
	pLimit := pageSize
	pageList = 10

	// Start : invite Send List
	fmt.Printf(fmt.Sprintf("CALL ZSP_INVITE_SEND_LIST('%v', %v, %v, '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pRecrutSn, pSdy, pEdy))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_INVITE_SEND_LIST('%v', %v, %v, '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pRecrutSn, pSdy, pEdy),
		ora.I64, /* TOT_CNT */
		ora.S,   /* SEND_DT */
		ora.S,   /* SEND_DT_FMT */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* RECRUT_TITLE */
		ora.S,   /* PP_CHRG_NM */
		ora.I64, /* CNT */
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

	inviteSendList := make([]models.InviteSendList, 0)

	var (
		rslTotCnt      int64
		rslSendDt      string
		rslSendDtFmt   string
		rslRecrutSn    string
		rslRecrutTitle string
		rslSenderName  string
		rslCnt         int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			rslTotCnt = procRset.Row[0].(int64)
			rslSendDt = procRset.Row[1].(string)
			rslSendDtFmt = procRset.Row[2].(string)
			rslRecrutSn = procRset.Row[3].(string)
			rslRecrutTitle = procRset.Row[4].(string)
			rslSenderName = procRset.Row[5].(string)
			rslCnt = procRset.Row[6].(int64)

			inviteSendList = append(inviteSendList, models.InviteSendList{
				RslTotCnt:      rslTotCnt,
				RslSendDt:      rslSendDt,
				RslSendDtFmt:   rslSendDtFmt,
				RslRecrutSn:    rslRecrutSn,
				RslRecrutTitle: rslRecrutTitle,
				RslSenderName:  rslSenderName,
				RslCnt:         rslCnt,
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
			pagination += "<a href='javascript:void(0);' class='next disabled' id='next' data-page='" + strconv.Itoa(int(nextPageNo)) + "'>다음</a>"
		} else {
			pagination += "<a href='javascript:void(0);' class='next goPage' id='next' data-page='" + strconv.Itoa(int(nextPageNo)) + "'>다음</a>"
		}
	}

	// End : Invite Send List

	c.Data["RecruitList"] = recruitSubList

	c.Data["RslTotCnt"] = rslTotCnt
	c.Data["InviteSendList"] = inviteSendList

	c.Data["EvlPrgsStat"] = "A"

	//c.Data["ApplySortCd"] = pApplySortCd
	//c.Data["ApplySortWay"] = pApplySortWay

	c.Data["Pagination"] = pagination

	c.Data["TMenuId"] = "R00"
	c.Data["SMenuId"] = "R03"
	c.TplName = "invite/invite_send_list.html"
}

func (c *InviteSendListController) Post() {

	session := c.StartSession()

	memNO := session.Get(c.Ctx.Request.Context(), "mem_no")
	if memNO == nil {
		//c.Ctx.Redirect(302, "/login")
		c.Data["json"] = &models.RtnInviteSendList{}
		c.ServeJSON()
		return
	}

	memSn := session.Get(c.Ctx.Request.Context(), "mem_sn")
	if memSn == nil {
		//c.Ctx.Redirect(302, "/login")
		c.Data["json"] = &models.RtnInviteSendList{}
		c.ServeJSON()
		return
	}

	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := memNO // 기업회원번호

	pRecrutSn := c.GetString("recrut_sn") // 채용일련번호
	pSdy := c.GetString("sdy")
	pEdy := c.GetString("edy")

	if pSdy == "" {
		pSdy = "20200101"
	}

	if pEdy == "" {
		//pEdy = "20201201"
		nowTime := time.Now()
		pEdy = nowTime.Format("20060102")
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
	fmt.Printf(fmt.Sprintf("CALL ZSP_INVITE_SEND_LIST('%v', %v, %v, '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pRecrutSn, pSdy, pEdy))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_INVITE_SEND_LIST('%v', %v, %v, '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pRecrutSn, pSdy, pEdy),
		ora.I64, /* TOT_CNT */
		ora.S,   /* SEND_DT */
		ora.S,   /* SEND_DT_FMT */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* RECRUT_TITLE */
		ora.S,   /* PP_CHRG_NM */
		ora.I64, /* CNT */
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

	rtnInviteSendList := models.RtnInviteSendList{}
	inviteSendList := make([]models.InviteSendList, 0)

	var (
		rslTotCnt      int64
		rslSendDt      string
		rslSendDtFmt   string
		rslRecrutSn    string
		rslRecrutTitle string
		rslSenderName  string
		rslCnt         int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			rslTotCnt = procRset.Row[0].(int64)
			rslSendDt = procRset.Row[1].(string)
			rslSendDtFmt = procRset.Row[2].(string)
			rslRecrutSn = procRset.Row[3].(string)
			rslRecrutTitle = procRset.Row[4].(string)
			rslSenderName = procRset.Row[5].(string)
			rslCnt = procRset.Row[6].(int64)

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

			inviteSendList = append(inviteSendList, models.InviteSendList{
				RslTotCnt:      rslTotCnt,
				RslSendDt:      rslSendDt,
				RslSendDtFmt:   rslSendDtFmt,
				RslRecrutSn:    rslRecrutSn,
				RslRecrutTitle: rslRecrutTitle,
				RslSenderName:  rslSenderName,
				RslCnt:         rslCnt,
				Pagination:     pagination,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnInviteSendList = models.RtnInviteSendList{
			RtnInviteSendListData: inviteSendList,
		}
		// End : Recruit Stat List

		c.Data["json"] = &rtnInviteSendList
		c.ServeJSON()
	}
}
