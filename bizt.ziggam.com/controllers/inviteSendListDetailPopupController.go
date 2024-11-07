package controllers

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"

	"bizt.ziggam.com/models"
	"bizt.ziggam.com/utils"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

// InviteSendListDetailPopupController ...
type InviteSendListDetailPopupController struct {
	BaseController
}

// Get ...
func (c *InviteSendListDetailPopupController) Get() {

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

	pRecrutSn := c.GetString("recrut_sn")
	pSendDt := c.GetString("send_dt")
	pKeyword := c.GetString("keyword") // 검색어

	fmt.Printf(pLang)
	fmt.Printf(pEntpMemNo.(string))
	fmt.Printf(pRecrutSn)
	fmt.Printf(pSendDt)

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Recruit Apply List
	// pSdy := c.GetString("sdy")
	// pEdy := c.GetString("edy")

	// if pSdy == "" {
	// 	pSdy = "20200101"
	// }

	// if pEdy == "" {
	// 	pEdy = "20201201"
	// 	nowTime := time.Now()
	// 	pEdy = nowTime.Format("20060102")
	// }

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
		pPageSize = "100"
	}

	pageSize, err = strconv.ParseInt(pPageSize, 10, 64)
	if err != nil {
		//
	}
	pOffSet := (pageNo - 1) * pageSize
	pLimit := pageSize
	pageList = 10

	// Start : invite Send List
	fmt.Printf(fmt.Sprintf("CALL ZSP_INVITE_SEND_LIST_DETAIL('%v', %v, %v, '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pRecrutSn, pSendDt, pKeyword))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_INVITE_SEND_LIST_DETAIL('%v', %v, %v, '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pRecrutSn, pSendDt, pKeyword),
		ora.I64, /* TOT_CNT */
		ora.I64, /* ROWNO */
		ora.S,   /* SEND_DT */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* NM */
		ora.S,   /* EMAIL */
		ora.S,   /* MO_NO */
		ora.S,   /* EMAIL_MID */
		ora.S,   /* EMAIL_RESULT */
		ora.S,   /* EMAIL_DT */
		ora.S,   /* SMS_MID */
		ora.S,   /* SMS_RESULT */
		ora.S,   /* SMS_DT */
		ora.S,   /* LIST_REFUSE_YN */
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

	inviteSendListDetail := make([]models.InviteSendListDetail, 0)

	var (
		rslTotCnt      int64
		rslRowNo       int64
		rslSendDt      string
		rslRecrutSn    string
		rslName        string
		rslEmail       string
		rslPhone       string
		rslEmailMid    string
		rslEmailResult string
		rslEmailDt     string
		rslSmsMid      string
		rslSmsResult   string
		rslSmsDt       string
		rslListYN      string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			rslTotCnt = procRset.Row[0].(int64)
			rslRowNo = procRset.Row[1].(int64)
			rslSendDt = procRset.Row[2].(string)
			rslRecrutSn = procRset.Row[3].(string)

			rslName = procRset.Row[4].(string)
			rslEmail = procRset.Row[5].(string)
			rslPhone = procRset.Row[6].(string)

			rslEmailMid = procRset.Row[7].(string)
			rslEmailResult = procRset.Row[8].(string)
			rslEmailDt = procRset.Row[9].(string)

			rslSmsMid = procRset.Row[10].(string)
			rslSmsResult = procRset.Row[11].(string)
			rslSmsDt = procRset.Row[12].(string)

			rslListYN = procRset.Row[13].(string)

			inviteSendListDetail = append(inviteSendListDetail, models.InviteSendListDetail{
				RslTotCnt:      rslTotCnt,
				RslRowNo:       rslRowNo,
				RslSendDt:      rslSendDt,
				RslRecrutSn:    rslRecrutSn,
				RslName:        rslName,
				RslEmail:       rslEmail,
				RslPhone:       rslPhone,
				RslEmailMid:    rslEmailMid,
				RslEmailResult: rslEmailResult,
				RslEmailDt:     rslEmailDt,
				RslSmsMid:      rslSmsMid,
				RslSmsResult:   rslSmsResult,
				RslSmsDt:       rslSmsDt,
				RslListYN:      rslListYN,
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
			pagination += "<a href='javascript:void(0);' class='next disabled' id='prev' data-page='" + strconv.Itoa(int(nextPageNo)) + "'>다음</a>"
		} else {
			pagination += "<a href='javascript:void(0);' class='next goPage' id='prev' data-page='" + strconv.Itoa(int(nextPageNo)) + "'>다음</a>"
		}
	}

	// End : Invite Send List

	// Start : invite Send List SMS Check info
	fmt.Printf(fmt.Sprintf("CALL ZSP_INVITE_SEND_SMS_RET_CHK('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pSendDt))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_INVITE_SEND_SMS_RET_CHK('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pSendDt),
		ora.S, /* SEND_MID */
		ora.S, /* SEND_RESULT_CD */
		ora.S, /* SEND_RESULT_MSG */
		ora.S, /* SUCCESS_CNT */
		ora.S, /* ERROR_CNT */
		ora.S, /* DETAIL_CHECK */
		ora.S, /* UPT_DT */
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

	var (
		rslSendMid       string
		rslSendResultCd  string
		rslSendResultMsg string
		rslSuccessCnt    string
		rslErrorCnt      string
		rslDetailCheck   string
		rslUptDt         string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			rslSendMid = procRset.Row[0].(string)
			rslSendResultCd = procRset.Row[1].(string)
			rslSendResultMsg = procRset.Row[2].(string)
			rslSuccessCnt = procRset.Row[3].(string)
			rslErrorCnt = procRset.Row[4].(string)
			rslDetailCheck = procRset.Row[5].(string)
			rslUptDt = procRset.Row[6].(string)
		}

		if err := procRset.Err(); err != nil {
			panic(err)
		}

		fmt.Printf(fmt.Sprintf("rslSendMid:%v, rslSendResultCd:%v, rslSendResultMsg:%v, rslSuccessCnt:%v, rslErrorCnt:%v, rslDetailCheck:%v",
			rslSendMid, rslSendResultCd, rslSendResultMsg, rslSuccessCnt, rslErrorCnt, rslDetailCheck))

		if rslDetailCheck == "N" {

			var (
				rtnCd  int64
				rtnMsg string
			)

			var isCheckOk bool

			if rslSendMid != "" { // 정상적인 요청이 이루어 졌으면
				// 알리고에 결과 요청
				var respData utils.AligoSendSmsHisResp

				respData = utils.AligoSendSmsHis(rslSendMid, 1, 100)

				fmt.Printf(fmt.Sprintf("[AligoSendSmsHis] Result:%v", respData))

				var (
					phoneList     string
					smsResultList string
					smsSendDtList string
				)

				var isOkCount = 0 // 발송 완료 계산을 위해
				for index, val := range respData.List {

					// https://blog.acronym.co.kr/243
					retReg, _ := regexp.Compile("(^02.{0}|^01.{1}|[0-9]{3})([0-9]+)([0-9]{4})")

					convPhoneNum := retReg.ReplaceAllString(val.Receiver, "$1-$2-$3")

					convtime, _ := time.Parse("2006-01-02 15:04:05", val.SendDate)
					var smsTime string = convtime.Format("20060102150405")

					fmt.Printf(fmt.Sprintf("[AligoSendSmsHis][List] Receiver:%v, SmsState:%v, SendDate:%v, MoNo:%v, SendTime:%v", val.Receiver, val.SmsState, val.SendDate, convPhoneNum, smsTime))

					if val.SmsState == "발송완료" {
						isOkCount++
					}

					phoneList += convPhoneNum
					smsResultList += val.SmsState
					smsSendDtList += smsTime

					if len(respData.List)-1 > index {
						phoneList += ","
						smsResultList += ","
						smsSendDtList += ","
					}
				}

				isCheckOk = isOkCount == len(respData.List)

				// sms 상세 결과 재 갱신
				fmt.Printf(fmt.Sprintf("CALL ZSP_INVITE_SEND_LIST_SMS_UPT2('%v', '%v', '%v', '%v', '%v', %v, '%v', '%v', '%v', :1)",
					pLang, pEntpMemNo, pRecrutSn, pSendDt, rslSendMid, len(respData.List), phoneList, smsResultList, smsSendDtList))

				stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_INVITE_SEND_LIST_SMS_UPT2('%v', '%v', '%v', '%v', '%v', %v, '%v', '%v', '%v', :1)",
					pLang, pEntpMemNo, pRecrutSn, pSendDt, rslSendMid, len(respData.List), phoneList, smsResultList, smsSendDtList),
					ora.I64, /* RTN_CD */
					ora.S,   /* RTN_MSG */
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

				if procRset.IsOpen() {
					for procRset.Next() {
						rtnCd = procRset.Row[0].(int64)
						rtnMsg = procRset.Row[1].(string)
					}

					if err := procRset.Err(); err != nil {
						panic(err)
					}

					fmt.Printf(fmt.Sprintf(" ===> rtnCd:%v, rtnMsg:%v", rtnCd, rtnMsg))
				}
			}

			loc, _ := time.LoadLocation("Local")
			uptTime, _ := time.ParseInLocation("20060102150405", rslUptDt, loc)

			var uptFlag string
			if isCheckOk == true || time.Now().Sub(uptTime).Minutes() > 10 {
				uptFlag = "Y"

				// sms 상세 결과 재 갱신
				fmt.Printf(fmt.Sprintf("CALL ZSP_INVITE_SEND_SMS_RET_UPT('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
					pLang, pEntpMemNo, pRecrutSn, pSendDt, rslSendMid, "respData.ResultCode", "respData.Message", uptFlag))

				stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_INVITE_SEND_SMS_RET_UPT('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v',:1)",
					pLang, pEntpMemNo, pRecrutSn, pSendDt, rslSendMid, "respData.ResultCode", "respData.Message", uptFlag),
					ora.I64, /* RTN_CD */
					ora.S,   /* RTN_MSG */
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

				if procRset.IsOpen() {
					for procRset.Next() {
						rtnCd = procRset.Row[0].(int64)
						rtnMsg = procRset.Row[1].(string)
					}

					if err := procRset.Err(); err != nil {
						panic(err)
					}

					fmt.Printf(fmt.Sprintf(" ===> rtnCd:%v, rtnMsg:%v", rtnCd, rtnMsg))
				}
			}
		}
	}
	// End : invite Send List SMS Check info

	c.Data["RslTotCnt"] = rslTotCnt
	c.Data["InviteSendListDetail"] = inviteSendListDetail

	c.Data["Pagination"] = pagination

	c.Data["RecrutSn"] = pRecrutSn
	c.Data["SendDt"] = pSendDt

	c.Data["EvlPrgsStat"] = "A"

	c.Data["ApplySortCd"] = "01"
	c.Data["ApplySortWay"] = "01"

	c.Data["TMenuId"] = "R00"
	c.Data["SMenuId"] = "R03"

	c.TplName = "invite/invite_send_list_detail_popup.html"
}

// Post ...
func (c *InviteSendListDetailPopupController) Post() {

	session := c.StartSession()

	memNO := session.Get(c.Ctx.Request.Context(), "mem_no")
	if memNO == nil {
		c.Data["json"] = &models.RtnInviteSendListDetail{}
		c.ServeJSON()
		return
	}

	memSn := session.Get(c.Ctx.Request.Context(), "mem_sn")
	if memSn == nil {
		c.Data["json"] = &models.RtnInviteSendListDetail{}
		c.ServeJSON()
		return
	}

	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := memNO

	pRecrutSn := c.GetString("recrut_sn")
	pSendDt := c.GetString("send_dt")
	pKeyword := c.GetString("keyword") // 검색어

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
		pPageSize = "100"
	}

	//pPageSize = strconv.FormatInt(imPageNo, 16)

	pageSize, err = strconv.ParseInt(pPageSize, 10, 64)
	if err != nil {
		//
	}
	pOffSet := (pageNo - 1) * pageSize
	pLimit := pageSize
	pageList = 10

	// Start : invite Send List
	fmt.Printf(fmt.Sprintf("CALL ZSP_INVITE_SEND_LIST_DETAIL('%v', %v, %v, '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pRecrutSn, pSendDt, pKeyword))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_INVITE_SEND_LIST_DETAIL('%v', %v, %v, '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pRecrutSn, pSendDt, pKeyword),
		ora.I64, /* TOT_CNT */
		ora.I64, /* ROWNO */
		ora.S,   /* SEND_DT */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* NM */
		ora.S,   /* EMAIL */
		ora.S,   /* MO_NO */
		ora.S,   /* EMAIL_MID */
		ora.S,   /* EMAIL_RESULT */
		ora.S,   /* EMAIL_DT */
		ora.S,   /* SMS_MID */
		ora.S,   /* SMS_RESULT */
		ora.S,   /* SMS_DT */
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

	rtnInviteSendListDetail := models.RtnInviteSendListDetail{}
	inviteSendListDetail := make([]models.InviteSendListDetail, 0)

	var (
		rslTotCnt      int64
		rslRowNo       int64
		rslSendDt      string
		rslRecrutSn    string
		rslName        string
		rslEmail       string
		rslPhone       string
		rslEmailMid    string
		rslEmailResult string
		rslEmailDt     string
		rslSmsMid      string
		rslSmsResult   string
		rslSmsDt       string
		rslListYN      string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			rslTotCnt = procRset.Row[0].(int64)
			rslRowNo = procRset.Row[1].(int64)
			rslSendDt = procRset.Row[2].(string)
			rslRecrutSn = procRset.Row[3].(string)

			rslName = procRset.Row[4].(string)
			rslEmail = procRset.Row[5].(string)
			rslPhone = procRset.Row[6].(string)

			rslEmailMid = procRset.Row[7].(string)
			rslEmailResult = procRset.Row[8].(string)
			rslEmailDt = procRset.Row[9].(string)

			rslSmsMid = procRset.Row[10].(string)
			rslSmsResult = procRset.Row[11].(string)
			rslSmsDt = procRset.Row[12].(string)

			rslListYN = procRset.Row[13].(string)

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

			inviteSendListDetail = append(inviteSendListDetail, models.InviteSendListDetail{
				RslTotCnt:      rslTotCnt,
				RslRowNo:       rslRowNo,
				RslSendDt:      rslSendDt,
				RslRecrutSn:    rslRecrutSn,
				RslName:        rslName,
				RslEmail:       rslEmail,
				RslPhone:       rslPhone,
				RslEmailMid:    rslEmailMid,
				RslEmailResult: rslEmailResult,
				RslEmailDt:     rslEmailDt,
				RslSmsMid:      rslSmsMid,
				RslSmsResult:   rslSmsResult,
				RslSmsDt:       rslSmsDt,
				RslListYN:      rslListYN,
				Pagination:     pagination,
			})
		}

		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnInviteSendListDetail = models.RtnInviteSendListDetail{
			RtnInviteSendListDetailData: inviteSendListDetail,
		}
		// End : Recruit Stat List

		c.Data["json"] = &rtnInviteSendListDetail
		c.ServeJSON()
	}
}
