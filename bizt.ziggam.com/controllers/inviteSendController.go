package controllers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"bizt.ziggam.com/models"
	"bizt.ziggam.com/tables"
	"bizt.ziggam.com/utils"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

// InviteSendController ...
type InviteSendController struct {
	BaseController
}

// Get ...
func (c *InviteSendController) Get() {
	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	memId := session.Get(c.Ctx.Request.Context(), "mem_id")
	if memId == nil {
		c.Ctx.Redirect(302, "/login")
	}

	memNo := session.Get(c.Ctx.Request.Context(), "mem_no")
	if memNo == nil {
		c.Ctx.Redirect(302, "/login")
	}

	memSn := session.Get(c.Ctx.Request.Context(), "mem_sn")
	if memSn == nil {
		c.Ctx.Redirect(302, "/login")
	}

	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := memNo
	pPpChrgSn := memSn

	pGbnCd := c.GetString("gbn_cd") // 구분코드(A:전체, I:채용중, E:종료)

	pInitRecrutSn := c.GetString("initRecrutSn")
	pInitList := c.GetString("initList")

	log.Debug(pInitRecrutSn)
	log.Debug(pInitList)

	if pGbnCd == "" {
		pGbnCd = "00"
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

	log.Debug(pPpChrgSn.(string))

	pOffSet := "0"
	pLimit := "100"
	pSortGbn := "03"

	recruitSubList := make([]models.RecruitSubList, 0)

	if pInitRecrutSn != "" {
		// Start : invite recruit info
		log.Debug(fmt.Sprintf("CALL ZSP_INVITE_RECRUIT_INFO('%v', '%v', :1)",
			pLang, pInitRecrutSn))

		stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_INVITE_RECRUIT_INFO('%v', '%v', :1)",
			pLang, pInitRecrutSn),
			ora.S, /* RECRUT_SN */
			ora.S, /* PRGS_STAT */
			ora.S, /* RECRUT_TITLE */
			ora.S, /* TRIM_RECRUT_TITLE*/
			ora.S, /* UP_JOB_GRP */
			ora.S, /* JOB_GRP */
			ora.S, /* END_DT */
			ora.S, /* SDY */
			ora.S, /* EDY */

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

		var (
			sRecrutSn        string
			sPrgsStat        string
			sRecrutTitle     string
			sTrimRecrutTitle string
			sUpJobGrp        string
			sJobGrp          string
			sEndDt           string
			sSdy             string
			sEdy             string
		)

		if procRset.IsOpen() {
			for procRset.Next() {
				sRecrutSn = procRset.Row[0].(string)
				sPrgsStat = procRset.Row[1].(string)
				sRecrutTitle = procRset.Row[2].(string)
				sTrimRecrutTitle = procRset.Row[3].(string)
				sUpJobGrp = procRset.Row[4].(string)
				sJobGrp = procRset.Row[5].(string)
				sEndDt = procRset.Row[6].(string)
				sSdy = procRset.Row[7].(string)
				sEdy = procRset.Row[8].(string)

				recruitSubList = append(recruitSubList, models.RecruitSubList{
					STotCnt:      0,                   // 미사용
					SEntpMemNo:   pEntpMemNo.(string), // 미사용
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

			log.Debug("sRecrutSn:" + sRecrutSn)
			log.Debug("sPrgsStat:" + sPrgsStat)
			log.Debug("sRecrutTitle:" + sRecrutTitle)
			log.Debug("sTrimRecrutTitle:" + sTrimRecrutTitle)
			log.Debug("sUpJobGrp:" + sUpJobGrp)
			log.Debug("sJobGrp:" + sJobGrp)
			log.Debug("sEndDt:" + sEndDt)
			log.Debug("sSdy:" + sSdy)
			log.Debug("sEdy:" + sEdy)

			if sPrgsStat == "ING" {
				pGbnCd = "02"
			} else if sPrgsStat == "END" {
				pGbnCd = "03"
			} else if sPrgsStat == "WAIT" {
				pGbnCd = "01"
			} else {
				pGbnCd = "00"
			}
		}
		// End : invite recruit info
	}

	log.Debug("pGbnCd:" + pGbnCd)

	// Start : invite recruit info list
	log.Debug(fmt.Sprintf("CALL ZSP_INVITE_RECRUIT_LIST_TYPE('%v', %v, %v, '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pGbnCd, pSortGbn))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_INVITE_RECRUIT_LIST_TYPE('%v', %v, %v, '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pGbnCd, pSortGbn),
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

	//recruitSubList := make([]models.RecruitSubList, 0)

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

			if pInitRecrutSn != sRecrutSn {
				recruitSubList = append(recruitSubList, models.RecruitSubList{
					STotCnt:      sTotCnt,    // 미사용
					SEntpMemNo:   sEntpMemNo, // 미사용
					SRecrutSn:    sRecrutSn,
					SPrgsStat:    sPrgsStat,
					SRecrutTitle: sRecrutTitle,
					SUpJobGrp:    sUpJobGrp,
					SJobGrp:      sJobGrp,
				})
			}
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : invite recruit info list

	// Start : invite base info
	log.Debug(fmt.Sprintf("CALL ZSP_INVITE_INFO('%v', '%v', '%v', :1)", pLang, pEntpMemNo, pPpChrgSn))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_INVITE_INFO('%v', '%v', '%v', :1)", pLang, pEntpMemNo, pPpChrgSn),
		ora.S, /* ENTP_KO_NM */
		ora.S, /* PP_CHRG_NM */
		ora.S, /* PP_CHRG_GBN_CD */
		ora.S, /* PP_CHRG_BP_NM */
		ora.S, /* PP_CHRG_TEL_NO */
		ora.S, /* EMAIL */
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
		entpKoName string
		ppChrgName string

		ppChrgGbnCd string
		ppChrgBpNm  string
		ppChrgTelNo string
		ppChrgEmail string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			entpKoName = procRset.Row[0].(string)
			ppChrgName = procRset.Row[1].(string)
			ppChrgGbnCd = procRset.Row[2].(string)
			ppChrgBpNm = procRset.Row[3].(string)
			ppChrgTelNo = procRset.Row[4].(string)
			ppChrgEmail = procRset.Row[5].(string)
		}

		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : invite recruit info

	c.Data["GbnCd"] = pGbnCd
	c.Data["EntpMemNo"] = pEntpMemNo
	c.Data["RecruitList"] = recruitSubList

	c.Data["EntpKoNm"] = entpKoName
	c.Data["PpChrgName"] = ppChrgName

	var convTel string
	telNoFmt, _ := regexp.Compile("(^02.{0}|^01.{1}|[0-9]{3})([0-9]+)([0-9]{4})")
	telNoCov := telNoFmt.ReplaceAllString(ppChrgTelNo, "$1-$2-$3")

	matched, _ := regexp.MatchString("^01(?:0|1|[6-9])-(?:\\d{3}|\\d{4})-\\d{4}$", telNoCov)
	if matched == true {
		convTel = telNoCov
	} else {
		convTel = "-"
	}

	c.Data["ppChrgGbnCd"] = ppChrgGbnCd
	c.Data["ppChrgBpNm"] = ppChrgBpNm
	c.Data["ppChrgTelNo"] = convTel
	c.Data["ppChrgEmail"] = ppChrgEmail

	c.Data["InitList"] = pInitList
	c.Data["InitRecrutSn"] = pInitRecrutSn

	c.Data["TMenuId"] = "R00"
	c.Data["SMenuId"] = "R00"

	c.TplName = "invite/invite_send_write.html"
}

// Post ...
func (c *InviteSendController) Post() {
	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	memId := session.Get(c.Ctx.Request.Context(), "mem_id")
	if memId == nil {
		c.Data["json"] = &models.DefaultResult{RtnCd: 99, RtnMsg: "/login"}
		c.ServeJSON()
		return
	}

	memNO := session.Get(c.Ctx.Request.Context(), "mem_no")
	if memNO == nil {
		c.Data["json"] = &models.DefaultResult{RtnCd: 99, RtnMsg: "/login"}
		c.ServeJSON()
		return
	}

	memSn := session.Get(c.Ctx.Request.Context(), "mem_sn")
	if memSn == nil {
		c.Data["json"] = &models.DefaultResult{RtnCd: 99, RtnMsg: "/login"}
		c.ServeJSON()
		return
	}

	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := memNO
	pPpChrgSn := memSn

	pRecruitSn := c.GetString("recruit_sn")
	pSendList := c.GetString("send_list")

	pSendTitle := c.GetString("send_title")
	pSendMsg := c.GetString("send_msg")
	pIsTest := c.GetString("is_test")

	var isTest bool
	if pIsTest == "true" {
		isTest = true
	} else {
		isTest = false
	}

	log.Debug(pIsTest)
	log.Debug(fmt.Sprintf("%v", isTest))
	log.Debug(pSendTitle)
	log.Debug(pSendMsg)

	nowTime := time.Now()
	var sendTime string = nowTime.Format("20060102150405")

	// 리스트 Json 데이터 -> array 로
	var retSendData []models.InviteMember
	err := json.Unmarshal([]byte(pSendList), &retSendData)
	if err != nil {
		panic(err)
	}

	//var convSendData []models.InviteMember

	// for _, val := range retSendData {

	// 	if len(val.Name) <= 0 || len(val.Email) <= 0 || len(val.Phone) <= 0 {

	// 		log.Debug(fmt.Sprintf("[Invite][Send][Error] SendTime:%v, EntpMemNo: %v, SendTitle:%v, SendMsg:%v, RecruitSn: %v, Name:%v, Email:%v, Phone:%v -> Invalid Data",
	// 			sendTime, pEntpMemNo, pSendTitle, pSendMsg, pRecruitSn, val.Name, val.Email, val.Phone))
	// 		continue
	// 	}

	// 	convSendData = append(convSendData, models.InviteMember{
	// 		Name:  val.Name,
	// 		Email: val.Email,
	// 		Phone: val.Phone,
	// 	})
	// }

	// retSendData = append(retSendData, models.InviteMember{
	// 	Name:  "이동관",
	// 	Email: "dongkale@naver.com",
	// 	Phone: "010-5226-2107",
	// })

	//retSendData = removeDuplicateValues(retSendData)

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : invite send info
	log.Debug(fmt.Sprintf("CALL ZSP_INVITE_SEND_INFO('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecruitSn))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_INVITE_SEND_INFO('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecruitSn),
		ora.S, // ENTP_MEM_NO
		ora.S, // ENTP_KO_NM
		ora.S, // RECRUT_SN
		ora.S, // RECRUT_TITLE
		ora.S, // UP_JOB_GRP
		ora.S, // JOB_GRP
		ora.S, // ENTP_VD_URL
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

	var (
		entpMemNo   string
		entpKoNm    string
		recrutSn    string
		recrutTitle string
		upJobGrp    string
		jobGrp      string
		entpVdUrl   string
	)

	rtnData := models.DefaultResult{}

	if procRset.IsOpen() {
		for procRset.Next() {
			entpMemNo = procRset.Row[0].(string)
			entpKoNm = procRset.Row[1].(string)
			recrutSn = procRset.Row[2].(string)
			recrutTitle = procRset.Row[3].(string)
			upJobGrp = procRset.Row[4].(string)
			jobGrp = procRset.Row[5].(string)
			entpVdUrl = procRset.Row[6].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		// log.Debug(fmt.Sprintf("entpMemNo: %v, entpKoNm: %v, recrutSn: %v, recrutTitle: %v, entpVdUrl: %v",
		// 	entpMemNo, entpKoNm, recrutSn, recrutTitle, entpVdUrl))

		if entpMemNo == pEntpMemNo && recrutSn == pRecruitSn {
			rtnData = models.DefaultResult{
				RtnCd:  1,
				RtnMsg: "Success",
			}

			log.Debug(fmt.Sprintf("[Invite][Send][Step 1] SendTime:%v, pEntpMemNo: %v, pSendTitle:%v, pSendMsg:%v, RecruitSn: %v, upJobGrp:%v,, jobGrp:%v, PpChrgSn:%v, sendList: %v, sendCount: %v, ",
				sendTime, pEntpMemNo, pSendTitle, pSendMsg, pRecruitSn, upJobGrp, jobGrp, pPpChrgSn, retSendData, len(retSendData)))

			go InviteSend(isTest, sendTime, entpMemNo, entpKoNm, pSendTitle, pSendMsg, pRecruitSn, recrutTitle, upJobGrp, jobGrp, pPpChrgSn.(string), entpVdUrl, retSendData)
		} else {
			rtnData = models.DefaultResult{
				RtnCd:  2,
				RtnMsg: "Invalid [EntpMemNo, RecruitSn]",
			}

			log.Debug(fmt.Sprintf("[Invite][Send][Step 1] Error -> SendTime:%v, pEntpMemNo: %v, RecruitSn: %v, PpChrgSn:%v, SendMsg: %v, sendList: %v",
				sendTime, pEntpMemNo, pRecruitSn, pPpChrgSn, pSendMsg, retSendData))
		}
	}
	// End : Message Detail Top Info

	c.Data["json"] = &rtnData
	c.ServeJSON()
}

// func removeDuplicateValues(checkArray []models.InviteMember) []models.InviteMember {
// 	keys := make(map[models.InviteMember]bool)
// 	list := []models.InviteMember{}

// 	// If the key(values of the slice) is not equal
// 	// to the already present value in new slice (list)
// 	// then we append it. else we jump on another element.
// 	for _, entry := range checkArray {
// 		if _, value := keys[entry]; !value {
// 			keys[entry] = true
// 			list = append(list, entry)
// 		}
// 	}
// 	return list
// }

// InviteSend ...
func InviteSend(isTest bool, sendTime string, entpMemNo string, entpKoNm string, title string, msg string, recruitSn string, recruitTitle string, upJobGrp string, jobGrp string, ppChrgSn string, entpVdUrl string, sendList []models.InviteMember) {
	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	pLang, _ := beego.AppConfig.String("lang")

	imgServer, _  := beego.AppConfig.String("viewpath")

	googleStore, _ := beego.AppConfig.String("googlestore")
	appleStore, _ := beego.AppConfig.String("applestore")
	mailTo, _ := beego.AppConfig.String("mailto")

	siteUrl, _ := beego.AppConfig.String("siteurl")

	bridgeUrl, _ := beego.AppConfig.String("bridgeUrl")

	var resultEntpVdUrl string
	var respData utils.NaverShortUrlResp
	respData = utils.NaverShortUrlReq(entpVdUrl)
	if respData.IsOk() == false {
		log.Debug(fmt.Sprintf("[Invite][Send][Step 2] EntpVideo Short URL Fail(%v, %v) -> SendTime:%v, EntpMemNo:%v, RecruitSn:%v",
			respData.Code, respData.Message, sendTime, entpMemNo, recruitSn))

		resultEntpVdUrl = entpVdUrl
	} else {
		resultEntpVdUrl = respData.Result.Url
	}

	if isTest == true {
		InviteSendSmsAll(true, bridgeUrl, sendTime, entpMemNo, entpKoNm, title, msg, recruitSn, recruitTitle, upJobGrp, jobGrp, resultEntpVdUrl, sendList)

		InviteSendEmailAll(true, bridgeUrl, imgServer, mailTo, "support@ziggam.com", googleStore, appleStore, siteUrl, sendTime, entpMemNo, entpKoNm, title, msg, recruitSn, recruitTitle, upJobGrp, jobGrp, resultEntpVdUrl, sendList)
	} else {
		// Start : Oracle DB Connection
		env, srv, ses, err := GetRawConnection()
		defer env.Close()
		defer srv.Close()
		defer ses.Close()
		if err != nil {
			panic(err)
		}
		// End : Oracle DB Connection

		var nameList string
		var emailList string
		var phoneList string
		for index, val := range sendList {
			nameList += val.Name
			emailList += val.Email
			phoneList += val.Phone

			if len(sendList)-1 > index {
				nameList += ","
				emailList += ","
				phoneList += ","
			}
		}

		log.Debug(fmt.Sprintf("CALL ZSP_INVITE_SEND_LIST_REG('%v', '%v', '%v', '%v', %v, '%v', '%v', '%v', '%v', :1)",
			pLang, entpMemNo, recruitSn, sendTime, len(sendList), nameList, emailList, phoneList, ppChrgSn))
		stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_INVITE_SEND_LIST_REG('%v', '%v', '%v', '%v', %v, '%v', '%v', '%v', '%v', :1)",
			pLang, entpMemNo, recruitSn, sendTime, len(sendList), nameList, emailList, phoneList, ppChrgSn),
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

		var (
			rtnCd  int64
			rtnMsg string
		)

		//rtnInvite := make([]models.RtnInvite, 0)

		if procRset.IsOpen() {
			for procRset.Next() {
				rtnCd = procRset.Row[0].(int64)
				rtnMsg = procRset.Row[1].(string)

				if rtnCd == 1 {
					if len(sendList) > 0 {
						InviteSendSmsAll(false, bridgeUrl, sendTime, entpMemNo, entpKoNm, title, msg, recruitSn, recruitTitle, upJobGrp, jobGrp, resultEntpVdUrl, sendList)

						InviteSendEmailAll(false, bridgeUrl, imgServer, mailTo, "support@ziggam.com", googleStore, appleStore, siteUrl, sendTime, entpMemNo, entpKoNm, title, msg, recruitSn, recruitTitle, upJobGrp, jobGrp, resultEntpVdUrl, sendList)
					}
				} else {
					log.Debug(fmt.Sprintf("[Invite][Send][Step 2] Error rtnCd:%v, rtnMsg:%v", rtnCd, rtnMsg))
				}
			}

			if err := procRset.Err(); err != nil {
				panic(err)
			}

			log.Debug(fmt.Sprintf(" ===> rtnCd:%v, rtnMsg:%v", rtnCd, rtnMsg))

			log.Debug(fmt.Sprintf("[Invite][Send][Step 2] rtnInvite:%v, rtnCd:%v, rtnMsg:%v", sendList, rtnCd, rtnMsg))
		}
	}
}

// InviteSendSmsAll ...
func InviteSendSmsAll(isTest bool, bridgeUrl string, sendTime string, entpMemNo string, entpKoNm string, title string, msg string, recruitSn string, recruitTitle string, upJobGrp string, jobGrp string, entpVdUrl string, sendList []models.InviteMember) {

	var sendTitle = title

	massList := []utils.AligoSendMass{}

	for _, val := range sendList {

		if val.Phone == tables.InviteSmsEmailSelectChar {
			logs.Debug(fmt.Sprintf("[Invite][Send(Sms)][Step 3] Skip -> SendTime:%v, EntpMemNo:%v, RecruitSn:%v, Count:%v, Name:%v, Phone:%v, Email:%v",
				sendTime, entpMemNo, recruitSn, len(sendList), val.Name, val.Phone, val.Email))
			continue
		}

		var resultRecruitUrl string

		var recruitUrl = bridgeUrl + "?" + url.QueryEscape(fmt.Sprintf("entpmemno=%v&recruitsn=%v&reqname=%v&reqemail=%v&reqmono=%v", entpMemNo, recruitSn, val.Name, val.Email, val.Phone))

		var respData utils.NaverShortUrlResp
		respData = utils.NaverShortUrlReq(recruitUrl)
		if respData.IsOk() == false {
			logs.Debug(fmt.Sprintf("[Invite][Send(Sms)][Step 3] Short URL Fail(%v, %v) -> SendTime:%v, EntpMemNo:%v, RecruitSn:%v, Count:%v, Name:%v, Phone:%v, Email:%v",
				respData.Code, respData.Message, sendTime, entpMemNo, recruitSn, len(sendList), val.Name, val.Phone, val.Email))

			resultRecruitUrl = recruitUrl
		} else {
			resultRecruitUrl = respData.Result.Url
		}

		/*
			var sendMsg = fmt.Sprintf("%v 님, 안녕하세요.\\n"+
				"%v 에서 직감을 통한 영상 지원을 요청하셨습니다.\\n"+
				"\\n"+
				"※ 하단 링크를 통해 앱 설치 후 자세한 채용공고 내용을 확인하실 수 있습니다.\\n"+
				"※ 채용공고 마감일 이후에는 지원이 불가하니 유의해주세요.\\n"+
				"\\n"+
				"▶ 기업명: %v\\n"+
				"▶ 채용공고: %v\\n"+
				"▶ 직무: %v > %v\\n"+
				"▶ 바로가기: %v\\n"+
				"\\n"+
				"* 직감이란? 직감은 영상 기반의 채용 서비스를 제공하는 플랫폼입니다. 모바일 어플리케이션을 다운로드하여 진행하실 수 있습니다.\\n"+
				"직감 영상 인터뷰 가이드 바로가기: %v\\n"+
				"* 수신자 정보는 30일 후 자동 삭제됩니다. 즉시 삭제를 원하실 경우 [바로가기] 링크를 통해 삭제 요청을 하실 수 있습니다.\\n"+
				"* 영상 인터뷰 관련 문의 : support@ziggam.com\\n",
				val.Name,
				entpKoNm,
				entpKoNm,
				recruitTitle,
				upJobGrp, jobGrp,
				resultRecruitUrl,
				entpVdUrl)
		*/

		var convMsg = strings.Replace(msg, "{지원자명}", val.Name, 100)

		convMsg = strings.Replace(convMsg, "{기업명}", entpKoNm, 100)
		convMsg = strings.Replace(convMsg, "{채용공고 제목}", recruitTitle, 100)
		convMsg = strings.Replace(convMsg, "{채용 공고 제목}", recruitTitle, 100)
		convMsg = strings.Replace(convMsg, "{공고 제목}", recruitTitle, 100)
		convMsg = strings.Replace(convMsg, "{1차직군}", upJobGrp, 100)
		convMsg = strings.Replace(convMsg, "{1차 직군}", upJobGrp, 100)
		convMsg = strings.Replace(convMsg, "{2차직군}", jobGrp, 100)
		convMsg = strings.Replace(convMsg, "{2차 직군}", jobGrp, 100)
		convMsg = strings.Replace(convMsg, "{URL1}", resultRecruitUrl, 100)
		convMsg = strings.Replace(convMsg, "{URL2}", entpVdUrl, 100)

		logs.Debug(convMsg)

		// var sendMsg = fmt.Sprintf("%v\\n"+
		// 	"\\n"+
		// 	"※ 하단 링크를 통해 앱 설치 후 자세한 채용공고 내용을 확인하실 수 있습니다.\\n"+
		// 	"※ 채용공고 마감일 이후에는 지원이 불가하니 유의해주세요.\\n"+
		// 	"\\n"+
		// 	"▶ 기업명: %v\\n"+
		// 	"▶ 채용공고: %v\\n"+
		// 	"▶ 직무: %v > %v\\n"+
		// 	"▶ 바로가기: %v\\n"+
		// 	"\\n"+
		// 	"* 직감이란? 직감은 영상 기반의 채용 서비스를 제공하는 플랫폼입니다. 모바일 어플리케이션을 다운로드하여 진행하실 수 있습니다.\\n"+
		// 	"직감 영상 인터뷰 가이드 바로가기: %v\\n"+
		// 	"* 수신자 정보는 30일 후 자동 삭제됩니다. 즉시 삭제를 원하실 경우 [바로가기] 링크를 통해 삭제 요청을 하실 수 있습니다.\\n"+
		// 	"* 영상 인터뷰 관련 문의 : support@ziggam.com\\n",
		// 	convMsg,
		// 	entpKoNm,
		// 	recruitTitle,
		// 	upJobGrp, jobGrp,
		// 	resultRecruitUrl,
		// 	entpVdUrl)

		var sendMsg = fmt.Sprintf("%v\\n", convMsg)

		massList = append(massList, utils.AligoSendMass{
			PhoneNum: val.Phone,
			Message:  sendMsg,
		})

		logs.Debug(fmt.Sprintf("[Invite][Send(Sms)][Step 3] Result -> SendTime:%v, EntpMemNo:%v, RecruitSn:%v, Count:%v, Name:%v, Phone:%v, Email:%v",
			sendTime, entpMemNo, recruitSn, len(sendList), val.Name, val.Phone, val.Email))
	}

	// SMS 전송
	var respData utils.AligoSendSmsResp

	respData = utils.AligoSendMassSms(massList, sendTitle, "LMS")

	if respData.IsOk() == false {
		logs.Debug(fmt.Sprintf("[Invite][Send(Sms)][Step 4] AligoSendMassSms Error!! -> respData:%v", respData))
	} else {
		logs.Debug(fmt.Sprintf("[Invite][Send(Sms)][Step 4] Result -> respData:%v", respData))
	}

	// Result
	// var respData2 utils.AligoSendSmsHisResp

	// respData2 = utils.AligoSendSmsHis(respData.MsgID, 1, 30)

	// beego.Trace(respData2)

	if isTest == false {

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

		var (
			rtnCd  int64
			rtnMsg string
		)

		// Start : SMS Send Result(State) Update

		// type AligoSendSmsResp struct {
		// 	ResultCode string `json:"result_code"`
		// 	Message    string `json:"message"`
		// 	MsgID      string `json:"msg_id"`
		// 	SuccessCnt uint32 `json:"success_cnt"`
		// 	ErrorCnt   uint32 `json:"error_cnt"`
		// 	MsgType    string `json:"msg_type"`
		// }

		logs.Debug(fmt.Sprintf("CALL ZSP_INVITE_SEND_SMS_RET_REG('%v', '%v', '%v', '%v', '%v', %v, '%v', '%v', %v, %v, :1)",
			pLang, entpMemNo, recruitSn, sendTime, respData.MsgID, len(sendList), respData.ResultCode, respData.Message, respData.SuccessCnt, respData.ErrorCnt))
		stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_INVITE_SEND_SMS_RET_REG('%v', '%v', '%v', '%v', '%v', %v, '%v', '%v', %v, %v, :1)",
			pLang, entpMemNo, recruitSn, sendTime, respData.MsgID, len(sendList), respData.ResultCode, respData.Message, respData.SuccessCnt, respData.ErrorCnt),
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

			logs.Debug(fmt.Sprintf(" ===> rtnCd:%v, rtnMsg:%v", rtnCd, rtnMsg))
		}

		logs.Debug(fmt.Sprintf("[Invite][Send(Sms)][Step 5] SendTime:%v, EntpMemNo:%v, RecruitSn:%v, Count:%v, MsgID:%v, rtnCd:%v, rtnMsg:%v",
			sendTime, entpMemNo, recruitSn, len(sendList), respData.MsgID, rtnCd, rtnMsg))
		// End : SMS Send Result(State) Update

		// Start : SMS Send Result Update
		var nameList string
		var emailList string
		var phoneList string
		var resultMid string
		var resultMsg string
		for index, val := range sendList {

			nameList += val.Name
			emailList += val.Email
			phoneList += val.Phone

			if val.Phone == tables.InviteSmsEmailSelectChar {
				resultMid += "0"
				resultMsg += "미발송"
			} else {
				if respData.IsOk() == true {
					resultMid += respData.MsgID
					resultMsg += "대기중"
				} else {
					resultMid += respData.MsgID
					resultMsg += respData.Message
				}
			}

			if len(sendList)-1 > index {
				nameList += ","
				emailList += ","
				phoneList += ","
				resultMid += ","
				resultMsg += ","
			}
		}

		logs.Debug(fmt.Sprintf("CALL ZSP_INVITE_SEND_LIST_SMS_UPT('%v', '%v', '%v', '%v', %v, '%v', '%v', '%v', '%v', '%v', :1)",
			pLang, entpMemNo, recruitSn, sendTime, len(sendList), nameList, emailList, phoneList, resultMid, resultMsg))
		stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_INVITE_SEND_LIST_SMS_UPT('%v', '%v', '%v', '%v', %v, '%v', '%v', '%v', '%v', '%v', :1)",
			pLang, entpMemNo, recruitSn, sendTime, len(sendList), nameList, emailList, phoneList, resultMid, resultMsg),
			ora.I64, /* RTN_CD */
			ora.S,   /* RTN_MSG */
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

		if procRset.IsOpen() {
			for procRset.Next() {
				rtnCd = procRset.Row[0].(int64)
				rtnMsg = procRset.Row[1].(string)
			}

			if err := procRset.Err(); err != nil {
				panic(err)
			}

			logs.Debug(fmt.Sprintf(" ===> rtnCd:%v, rtnMsg:%v", rtnCd, rtnMsg))
		}

		logs.Debug(fmt.Sprintf("[Invite][Send(Sms)][Step 6] SendTime:%v, EntpMemNo:%v, RecruitSn:%v, Count:%v, rtnCd:%v, rtnMsg:%v",
			sendTime, entpMemNo, recruitSn, len(sendList), rtnCd, rtnMsg))
		// End : SMS Send Result Update
	}
}

// InviteSendEmailAll ...
func InviteSendEmailAll(isTest bool, bridgeUrl string, imgServer string, mailTo string, supportMail string, googleStore string, appleStore string, siteUrl string, sendTime string, entpMemNo string, entpKoNm string, title string, msg string, recruitSn string, recruitTitle string, upJobGrp string, jobGrp string, entpVdUrl string, sendList []models.InviteMember) {

	for _, val := range sendList {

		if val.Email == tables.InviteSmsEmailSelectChar {
			logs.Debug(fmt.Sprintf("[Invite][Send(Email)][Step 4] Skip -> SendTime:%v, EntpMemNo:%v, RecruitSn:%v, Count:%v, Name:%v, Email:%v, Phone:%v",
				sendTime, entpMemNo, recruitSn, len(sendList), val.Name, val.Email, val.Phone))
			continue
		}

		var resultRecruitUrl string

		var recruitUrl = bridgeUrl + "?" + url.QueryEscape(fmt.Sprintf("entpmemno=%v&recruitsn=%v&reqname=%v&reqemail=%v&reqmono=%v", entpMemNo, recruitSn, val.Name, val.Email, val.Phone))

		var respData utils.NaverShortUrlResp
		respData = utils.NaverShortUrlReq(recruitUrl)
		if respData.IsOk() == false {
			logs.Debug(fmt.Sprintf("[Invite][Send(Email)][Step 3] Short URL Fail(%v, %v -> SendTime:%v, EntpMemNo:%v, RecruitSn:%v, Count:%v, Name:%v, Phone:%v, Email:%v",
				respData.Code, respData.Message, sendTime, entpMemNo, recruitSn, len(sendList), val.Name, val.Phone, val.Email))

			resultRecruitUrl = recruitUrl
		} else {
			resultRecruitUrl = respData.Result.Url
		}

		var convMsg = strings.Replace(msg, "{지원자명}", val.Name, 100)
		convMsg = strings.Replace(convMsg, "\n", "<br>", 200)

		convMsg = strings.Replace(convMsg, "{기업명}", entpKoNm, 100)
		convMsg = strings.Replace(convMsg, "{채용공고 제목}", recruitTitle, 100)
		convMsg = strings.Replace(convMsg, "{채용 공고 제목}", recruitTitle, 100)
		convMsg = strings.Replace(convMsg, "{공고 제목}", recruitTitle, 100)
		convMsg = strings.Replace(convMsg, "{1차직군}", upJobGrp, 100)
		convMsg = strings.Replace(convMsg, "{1차 직군}", upJobGrp, 100)
		convMsg = strings.Replace(convMsg, "{2차직군}", jobGrp, 100)
		convMsg = strings.Replace(convMsg, "{2차 직군}", jobGrp, 100)
		convMsg = strings.Replace(convMsg, "{URL1}", fmt.Sprintf("<a href=\"%v\">%v</a>", resultRecruitUrl, resultRecruitUrl), 10)
		convMsg = strings.Replace(convMsg, "{URL2}", fmt.Sprintf("<a href=\"%v\">%v</a>", entpVdUrl, entpVdUrl), 10)
		convMsg = strings.Replace(convMsg, supportMail, fmt.Sprintf("<a href=\"mailto:%v\">%v</a>", supportMail, mailTo), 10)
		//convMsg = strings.Replace(convMsg, "{URL2}", fmt.Sprintf("<a href=\"%v\">%v</a>", entpVdUrl, entpVdUrl), 10)

		logs.Debug(convMsg)

		var htmlString = `
			<html>
			<head>
			<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
			<meta name="viewport" content="width=device-width, initial-scale=1.0" />
			<title>직감 채용을 편하게, 면접을 영상으로</title>
			</head>
			<style type="text/css">
				// CLIENT-SPECIFIC RESETS
				// Outlook.com(Hotmail)의 전체 너비 및 적절한 줄 높이를 허용
				.ReadMsgBody{ width: 100%; }
				.ExternalClass{ width: 100%; }
				.ExternalClass, .ExternalClass p, .ExternalClass span, .ExternalClass font, .ExternalClass td, .ExternalClass div { line-height: 100%; }
				// Outlook 2007 이상에서 Outlook이 추가하는 테이블 주위의 간격을 제거
				table, td { mso-table-lspace: 0pt; mso-table-rspace: 0pt; }
				// Internet Explorer에서 크기가 조정된 이미지를 렌더링하는 방식을 수정
				img { -ms-interpolation-mode: bicubic; }
				// Webkit 및 Windows 기반 클라이언트가 텍스트 크기를 자동으로 조정하지 않도록 수정
				body, table, td, p, a, li, blockquote { -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%;}
			</style>
			<body>
			<!-- OUTERMOST CONTAINER TABLE -->
			<table border="0" cellpadding="0" cellspacing="0" width="100%" id="bodyTable" style="font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif">
				<!-- header -->
				<tr>
					<td style="padding:26px 0px 20px;border-bottom:4px solid #2ad0c7;">
						<table border="0" cellspacing="0" cellpadding="0" style="width:100%;">
							<tr>
								<td style="float:left">
									<a href="https://www.ziggam.com" target="_blank"><img style="display:block;border:0;" src="` + imgServer + `/mail/ic-main-logo.png" alt="직감 채용을 편하게,면접을 영상으로"></a>
								</td>
							</tr>
						</table>
					</td>
				</tr>
				<!-- //header -->
				<tr>
				<td>
					<!-- 600px - 800px CONTENTS CONTAINER TABLE -->
					<table border="0" cellpadding="0" cellspacing="0" width="800">
					<tr>
						<td style="text-align:left">
							<!-- 내용 -->
							<table border="0" cellspacing="0" cellpadding="0" style="width:100%">
								<!-- 메인 타이틀 -->
								<tr>
									<td style="padding-top:20px;width:100%;font-size:27px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:43px;color:#171717;letter-spacing:-1px">` + title + `</td>
								</tr>
								<!-- 내용 텍스트 -->
								<tr>
									<td style="padding-top:20px;width:100%;font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:25px;color:#171717;">` + convMsg + `</td>
								</tr>
								<tr>
									<td style="padding-top:20px;width:100%;font-size:22px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:43px;color:#171717;letter-spacing:-1px"> 지원방법</td>
								</tr>
								<tr>
									<td style="padding-top:20px;width:100%;font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:25px;color:#171717;">
										<p>1. 아래 ‘바로 지원하기’ 버튼을 클릭하시거나 구글플레이/앱스토어/원스토어를 통해 ‘직감’어플리케이션을 설치합니다.</p>
										<p>2. 앱을 실행 후 회원가입을 진행합니다.</p>
										<p>3. 메일로 돌아오셔서 아래 ‘직감 바로 가기’ 버튼을 다시 클릭하시거나 채용공고 제목을 검색합니다.</p>
										<p>4. 채용공고 내용을 확인 후 하단 [지원하기]를 선택해주세요.</p>
									</td>
								</tr>
								<tr>
									<td style="padding-top:20px;width:100%;font-size:22px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:43px;color:#171717;letter-spacing:-1px">유의사항</td>
								</tr>
								<tr>
									<td style="padding-top:20px;width:100%;font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:25px;color:#171717;">
										<p>• 직감의 영상 지원은 모바일 어플리케이션을 통해서만 진행하실 수 있습니다.</p>
										<p>• 영상 촬영은 세로 화면으로만 진행됩니다. 모바일을 가로로 돌리지 않도록 유의해주세요.</p>
										<p>• 채용공고 마감일 이후에는 지원이 불가능하니 유의해주세요.</p>
										<p>* 수신자 정보는 30일 후 자동 삭제됩니다. 즉시 삭제를 원하실 경우 [바로 지원하기] 링크를 통해 삭제 요청을 하실 수 있습니다.</p>
										<p>* 영상 인터뷰 관련 문의 : <a href="mailto:` + mailTo + `">` + supportMail + `</a></p>
									</td>
								</tr>								
								<table border="0" style="border:0px solid darkgray;width:100%;border-collapse:collapse;text-align:center;height:150px">
								<tr>
									<td>
									<a href="` + entpVdUrl + `"
										style="height:20px;border: 1px solid darkgray;width:180px;border-radius:3px;background-color: white;color: black; padding: 10px 10px; text-align: center; text-decoration: none;font-size: 16px; margin: 4px 2px; cursor: pointer;">
										직감 영상 인터뷰 가이드
									</a>
									<a href=` + resultRecruitUrl + `
										style="height:20px;border: 1px solid darkgray;width:180px;border-radius:3px;background-color: black; color: white; padding: 10px 10px; text-align: center; text-decoration: none;font-size: 16px; margin: 4px 2px; cursor: pointer;">
										바로 지원하기
									</a>
									</td>
								<tr>
								</table>
							</table>
							<!-- //내용 -->
						</td>
					</tr>
					</table>
				</td>
				</tr>
			</table>
			<table border="0" cellpadding="0" cellspacing="0" width="800">
			<!-- footer -->
				<tr>
					<td style="padding:20px 0px 20px 0;background-color:#ced3d6;">
						<table border="0" cellspacing="0" cellpadding="0" style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center">
							<tr>
								<td style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center;color:#4d5256">
									본 메일은 발신 전용 메일입니다. 문의 사항은 <a href="mailto:` + mailTo + `" style="color:#4d5256;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;">support@ziggam.com</a>을 이용해주세요.
								</td>
							</tr>
							<tr>
								<td style="padding:10px 0 20px 0;width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:20px;text-align:center;color:#4d5256">
									©<b>Qrate</b> Corp. All rights reserved&nbsp;&nbsp;&nbsp;I&nbsp;&nbsp;&nbsp;경기도 성남시 분당구 황새울로 311번길 36, 2F<br>
									사업자등록번호 : 850-88-00704&nbsp;&nbsp;&nbsp;I&nbsp;&nbsp;&nbsp;직업정보제공사업신고 : J151602019001
								</td>
							</tr>
							<tr>
								<td style="text-align:center">
									<table border="0" cellspacing="0" cellpadding="0" style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center">
										<tr>
											<td style="width:266px">&nbsp;</td>
											<td><a href="https://www.ziggam.com" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-home.png" alt="HOME"></a></td>
											<td style="width:16px">&nbsp;</td>
											<td><a href="` + siteUrl + `" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-ems.png" alt="EMS"></a></td>
											<td style="width:16px">&nbsp;</td>
											<td><a href="` + googleStore + `" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-andapp.png" alt="Android"></a></td>
											<td style="width:16px">&nbsp;</td>
											<td><a href="` + appleStore + `" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-iosapp.png" alt="IOS"></a></td>
											<td style="width:266px">&nbsp;</td>
										</tr>
									</table>
								</td>
							</tr>
						</table>
					</td>
				</tr>
				<!-- //footer -->
				</table>
			</body>
			</html>`

		isEmailSend := utils.SendMailEx(val.Email, val.Name, "no-reply@ziggam.com", fmt.Sprintf("%v by 직감", entpKoNm), title, htmlString)
		if isEmailSend != nil {
			logs.Debug(fmt.Sprintf("[Invite][Send(Email)][Step 4] Result Error(%v)!! -> SendTime:%v, EntpMemNo:%v, RecruitSn:%v, Count:%v, Name:%v, Email:%v, Phone:%v, result:%v",
				isEmailSend, sendTime, entpMemNo, recruitSn, len(sendList), val.Name, val.Email, val.Phone, isEmailSend))
		} else {
			logs.Debug(fmt.Sprintf("[Invite][Send(Email)][Step 4] Result -> SendTime:%v, EntpMemNo:%v, RecruitSn:%v, Count:%v, Name:%v, Email:%v, Phone:%v, result:%v",
				sendTime, entpMemNo, recruitSn, len(sendList), val.Name, val.Email, val.Phone, isEmailSend))
		}
	}

	if isTest == false {
		// Start : Oracle DB Connection
		env, srv, ses, err := GetRawConnection()
		defer env.Close()
		defer srv.Close()
		defer ses.Close()
		if err != nil {
			panic(err)
		}
		// End : Oracle DB Connection

		pLang, _ := beego.AppConfig.String("lang")

		var nameList string
		var emailList string
		var phoneList string
		var resultMid string
		var resultMsg string
		for index, val := range sendList {

			nameList += val.Name
			emailList += val.Email
			phoneList += val.Phone

			if val.Email == tables.InviteSmsEmailSelectChar {
				resultMid += "0"
				resultMsg += "미발송"
			} else {
				resultMid += "0"
				resultMsg += "발송완료"
			}

			if len(sendList)-1 > index {
				nameList += ","
				emailList += ","
				phoneList += ","
				resultMid += ","
				resultMsg += ","
			}
		}

		logs.Debug(fmt.Sprintf("CALL ZSP_INVITE_SEND_LIST_EMAIL_UPT('%v', '%v', '%v', '%v', %v, '%v', '%v', '%v', '%v', '%v', :1)",
			pLang, entpMemNo, recruitSn, sendTime, len(sendList), nameList, emailList, phoneList, resultMid, resultMsg))
		stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_INVITE_SEND_LIST_EMAIL_UPT('%v', '%v', '%v', '%v', %v, '%v', '%v', '%v', '%v', '%v', :1)",
			pLang, entpMemNo, recruitSn, sendTime, len(sendList), nameList, emailList, phoneList, resultMid, resultMsg),
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

		var (
			rtnCd  int64
			rtnMsg string
		)

		if procRset.IsOpen() {
			for procRset.Next() {
				rtnCd = procRset.Row[0].(int64)
				rtnMsg = procRset.Row[1].(string)
			}

			if err := procRset.Err(); err != nil {
				panic(err)
			}

			logs.Debug(fmt.Sprintf(" ===> rtnCd:%v, rtnMsg:%v", rtnCd, rtnMsg))
		}

		logs.Debug(fmt.Sprintf("[Invite][Send(Email)][Step 5] SendTime:%v, EntpMemNo:%v, RecruitSn:%v, Count:%v, rtnCd:%v, rtnMsg:%v",
			sendTime, entpMemNo, recruitSn, len(sendList), rtnCd, rtnMsg))
	}
}

/*
	var htmlString = `
			<html>
			<head>
			<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
			<meta http-equiv="X-UA-Compatible" content="ie=edge">
			<title>직감 채용을 편하게, 면접을 영상으로</title>
			</head>
			<body style="margin:0;padding:0;background-color:#f5f6f9;">
			<table border="0" cellspacing="0" cellpadding="0" align="center" style="font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif">
				<!-- header -->
				<tr>
					<td style="padding:26px 0px 20px;border-bottom:4px solid #2ad0c7;">
						<table border="0" cellspacing="0" cellpadding="0" style="width:100%;">
							<tr>
								<td style="float:left">
									<a href="https://www.ziggam.com" target="_blank"><img style="display:block;border:0;" src="` + imgServer + `/mail/ic-main-logo.png" alt="직감 채용을 편하게,면접을 영상으로"></a>
								</td>
							</tr>
						</table>
					</td>
				</tr>
				<!-- //header -->
				<!-- contents -->
				<tr>
					<td style="padding:40px 70px;background-color:#ffffff;">
						<table border="0" cellspacing="0" cellpadding="0" style="width:100%;">
							<tr>
								<td style="text-align:left">
									<!-- 내용 -->
									<table border="0" cellspacing="0" cellpadding="0" style="width:100%">
										<!-- 메인 타이틀 -->
										<tr>
											<td style="width:100%;font-size:27px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:43px;color:#171717;letter-spacing:-1px">` + title + `</td>
										</tr>
										<!-- 내용 텍스트 -->
										<tr>
											<td style="padding:35px 0 1px 0;width:100%;font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:25px;color:#171717;">` + convMsg + `</td>
										</tr>
									</table>
									<!-- //내용 -->
								</td>
							</tr>
						</table>
						<br><br>
						<p style="font-size:25px;float:left">채용정보</p>
						<br>
						<table border="1" style="border:1px solid darkgray;width:100%;border-collapse:collapse;">
							<tr>
								<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>기업명</b></td>
								<td style="padding:5px;font-size:16px;letter-spacing:0px">` + entpKoNm + `</td>
							</tr>
							<tr>
								<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>채용공고</b></td>
								<td style="padding:5px;font-size:16px;letter-spacing:0px">` + recruitTitle + `</td>
							</tr>
							<tr>
								<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>직무</b></td>
								<td style="padding:5px;font-size:16px;letter-spacing:0px">` + upJobGrp + `>` + jobGrp + `</td>
							</tr>
							<tr>
								<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>지원방법</b></td>
								<td style="padding:5px;font-size:16px;letter-spacing:0.3px"><p>1. 아래 ‘바로 지원하기’ 버튼을 클릭하시거나 구글플레이/앱스토어/원스토어를 통해 ‘직감’어플리케이션을 설치합니다.</p>
									<p>2. 앱을 실행 후 회원가입을 진행합니다.</p>
									<p>3. 메일로 돌아오셔서 아래 ‘직감 바로가기’ 버튼을 다시 클릭하시거나 채용공고 제목을 검색합니다.</p>
									<p>4. 채용공고 내용을 확인 후 하단 [지원하기]를 선택해주세요.</p>
								</td>
							</tr>
							<tr>
								<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>유의사항</b></td>
								<td style="padding:5px;font-size:16px;letter-spacing:0.3px"><p>• 직감의 영상 지원은 모바일 어플리케이션을 통해서만 진행하실 수 있습니다.</p>
									<p>• 영상 촬영은 세로 화면으로만 진행됩니다. 모바일을 가로로 돌리지 않도록 유의해주세요.</p>
									<p>• 채용공고 마감일 이후에는 지원이 불가능하니 유의해주세요.</p>
									<p>* 수신자 정보는 30일 후 자동 삭제됩니다. 즉시 삭제를 원하실 경우 [바로가기] 링크를 통해 삭제 요청을 하실 수 있습니다.</p>
									<p>* 영상 인터뷰 관련 문의 : <a href="mailto:` + mailTo + `">` + supportMail + `</a></p>
								</td>
							</tr>
						</table>
						<br>
						<div style="display: flex;text-align:center;justify-content: center;">
							<div></div>
							<a href="` + entpVdUrl + `"
							style="height:20px;border: 1px solid darkgray;width:180px;border-radius:3px;background-color: lightgray;color: black; padding: 10px 10px; text-align: center; text-decoration: none;font-size: 16px; margin: 4px 2px; cursor: pointer">
								직감 영상 인터뷰 가이드
							</a>
							<a href=` + resultRecruitUrl + `
							style="height:20px;border: 1px solid darkgray;width:180px;border-radius:3px;background-color: darkgray; color: black; padding: 10px 10px; text-align: center; text-decoration: none;font-size: 16px; margin: 4px 2px; cursor: pointer">
								바로 지원하기
							</a>
							<div></div>
						</div>
					</td>
				</tr>
				<!-- //contents -->
				<!-- footer -->
				<tr>
					<td style="padding:20px 0px 20px 0;background-color:#ced3d6;">
						<table border="0" cellspacing="0" cellpadding="0" style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center">
							<tr>
								<td style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center;color:#4d5256">
									본 메일은 발신 전용 메일입니다. 문의 사항은 <a href="mailto:` + mailTo + `" style="color:#4d5256;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;">support@ziggam.com</a>을 이용해주세요.
								</td>
							</tr>
							<tr>
								<td style="padding:10px 0 20px 0;width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:20px;text-align:center;color:#4d5256">
									©<b>Qrate</b> Corp. All rights reserved&nbsp;&nbsp;&nbsp;I&nbsp;&nbsp;&nbsp;경기도 성남시 분당구 황새울로 311번길 36, 2F<br>
									사업자등록번호 : 850-88-00704&nbsp;&nbsp;&nbsp;I&nbsp;&nbsp;&nbsp;직업정보제공사업신고 : J151602019001
								</td>
							</tr>
							<tr>
								<td style="text-align:center">
									<table border="0" cellspacing="0" cellpadding="0" style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center">
										<tr>
											<td style="width:266px">&nbsp;</td>
											<td><a href="https://www.ziggam.com" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-home.png" alt="HOME"></a></td>
											<td style="width:16px">&nbsp;</td>
											<td><a href="` + siteUrl + `" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-ems.png" alt="EMS"></a></td>
											<td style="width:16px">&nbsp;</td>
											<td><a href="` + googleStore + `" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-andapp.png" alt="Android"></a></td>
											<td style="width:16px">&nbsp;</td>
											<td><a href="` + appleStore + `" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-iosapp.png" alt="IOS"></a></td>
											<td style="width:266px">&nbsp;</td>
										</tr>
									</table>
								</td>
							</tr>
						</table>
					</td>
				</tr>
				<!-- //footer -->
			</table>
			</body>
			</html>`
*/

/*
	var htmlString = `
	<html>
	<head>
	<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
	<title>직감 채용을 편하게, 면접을 영상으로</title>
	</head>
	<style type="text/css">
		// CLIENT-SPECIFIC RESETS
		// Outlook.com(Hotmail)의 전체 너비 및 적절한 줄 높이를 허용
		.ReadMsgBody{ width: 100%; }
		.ExternalClass{ width: 100%; }
		.ExternalClass, .ExternalClass p, .ExternalClass span, .ExternalClass font, .ExternalClass td, .ExternalClass div { line-height: 100%; }
		// Outlook 2007 이상에서 Outlook이 추가하는 테이블 주위의 간격을 제거
		table, td { mso-table-lspace: 0pt; mso-table-rspace: 0pt; }
		// Internet Explorer에서 크기가 조정된 이미지를 렌더링하는 방식을 수정
		img { -ms-interpolation-mode: bicubic; }
		// Webkit 및 Windows 기반 클라이언트가 텍스트 크기를 자동으로 조정하지 않도록 수정
		body, table, td, p, a, li, blockquote { -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%;}
	</style>
	<body>
	<!-- OUTERMOST CONTAINER TABLE -->
	<table border="0" cellpadding="0" cellspacing="0" width="100%" id="bodyTable" style="font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif">
		<!-- header -->
		<tr>
			<td style="padding:26px 0px 20px;border-bottom:4px solid #2ad0c7;">
				<table border="0" cellspacing="0" cellpadding="0" style="width:100%;">
					<tr>
						<td style="float:left">
							<a href="https://www.ziggam.com" target="_blank"><img style="display:block;border:0;" src="` + imgServer + `/mail/ic-main-logo.png" alt="직감 채용을 편하게,면접을 영상으로"></a>
						</td>
					</tr>
				</table>
			</td>
		</tr>
		<!-- //header -->
		<tr>
		<td>
			<!-- 600px - 800px CONTENTS CONTAINER TABLE -->
			<table border="0" cellpadding="0" cellspacing="0" width="700">
			<tr>
				<td style="text-align:left">
					<!-- 내용 -->
					<table border="0" cellspacing="0" cellpadding="0" style="width:100%">
						<!-- 메인 타이틀 -->
						<tr>
							<td style="padding-top:20px;width:100%;font-size:27px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:43px;color:#171717;letter-spacing:-1px">` + title + `</td>
						</tr>
						<!-- 내용 텍스트 -->
						<tr>
							<td style="padding-top:20px;width:100%;font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:25px;color:#171717;">` + convMsg + `</td>
						</tr>
						<tr>
							<td style="padding-top:20px;padding-bottom:10px;width:100%;font-size:27px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:43px;color:#171717;letter-spacing:-1px">채용 정보</td>
						</tr>
						<table border="1" style="border:1px solid darkgray;width:100%;border-collapse:collapse;">
							<tr>
								<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>기업명</b></td>
								<td style="padding:5px;font-size:16px;letter-spacing:0px">` + entpKoNm + `</td>
							</tr>
							<tr>
								<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>채용공고</b></td>
								<td style="padding:5px;font-size:16px;letter-spacing:0px">` + recruitTitle + `</td>
							</tr>
							<tr>
								<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>직무</b></td>
								<td style="padding:5px;font-size:16px;letter-spacing:0px">` + upJobGrp + `>` + jobGrp + `</td>
							</tr>
							<tr>
								<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>지원방법</b></td>
								<td style="padding:5px;font-size:16px;letter-spacing:0.3px"><p>1. 아래 ‘바로 지원하기’ 버튼을 클릭하시거나 구글플레이/앱스토어/원스토어를 통해 ‘직감’어플리케이션을 설치합니다.</p>
									<p>2. 앱을 실행 후 회원가입을 진행합니다.</p>
									<p>3. 메일로 돌아오셔서 아래 ‘직감 바로 지원하기’ 버튼을 다시 클릭하시거나 채용공고 제목을 검색합니다.</p>
									<p>4. 채용공고 내용을 확인 후 하단 [바로 지원하기]를 선택해주세요.</p>
								</td>
							</tr>
							<tr>
								<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>유의사항</b></td>
								<td style="padding:5px;font-size:16px;letter-spacing:0.3px"><p>• 직감의 영상 지원은 모바일 어플리케이션을 통해서만 진행하실 수 있습니다.</p>
									<p>• 영상 촬영은 세로 화면으로만 진행됩니다. 모바일을 가로로 돌리지 않도록 유의해주세요.</p>
									<p>• 채용공고 마감일 이후에는 지원이 불가능하니 유의해주세요.</p>
									<p>* 수신자 정보는 30일 후 자동 삭제됩니다. 즉시 삭제를 원하실 경우 [바로 지원하기] 링크를 통해 삭제 요청을 하실 수 있습니다.</p>
									<p>* 영상 인터뷰 관련 문의 : <a href="mailto:` + mailTo + `">` + supportMail + `</a></p>
								</td>
							</tr>
						</table>
						<table border="0" style="border:0px solid darkgray;width:100%;border-collapse:collapse;text-align:center;height:100px">
						<tr>
							<td>
							<a href="` + entpVdUrl + `"
								style="height:20px;border: 1px solid darkgray;width:180px;border-radius:3px;background-color: lightgray;color: black; padding: 10px 10px; text-align: center; text-decoration: none;font-size: 16px; margin: 4px 2px; cursor: pointer;">
								직감 영상 인터뷰 가이드
							</a>
							<a href=` + resultRecruitUrl + `
								style="height:20px;border: 1px solid darkgray;width:180px;border-radius:3px;background-color: darkgray; color: black; padding: 10px 10px; text-align: center; text-decoration: none;font-size: 16px; margin: 4px 2px; cursor: pointer;">
								바로 지원하기
							</a>
							</td>
						<tr>
						</table>
					</table>
					<!-- //내용 -->
				</td>
			</tr>
			</table>
		</td>
		</tr>
	</table>
	<table border="0" cellpadding="0" cellspacing="0" width="700">
	<!-- footer -->
		<tr>
			<td style="padding:20px 0px 20px 0;background-color:#ced3d6;">
				<table border="0" cellspacing="0" cellpadding="0" style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center">
					<tr>
						<td style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center;color:#4d5256">
							본 메일은 발신 전용 메일입니다. 문의 사항은 <a href="mailto:` + mailTo + `" style="color:#4d5256;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;">support@ziggam.com</a>을 이용해주세요.
						</td>
					</tr>
					<tr>
						<td style="padding:10px 0 20px 0;width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:20px;text-align:center;color:#4d5256">
							©<b>Qrate</b> Corp. All rights reserved&nbsp;&nbsp;&nbsp;I&nbsp;&nbsp;&nbsp;경기도 성남시 분당구 황새울로 311번길 36, 2F<br>
							사업자등록번호 : 850-88-00704&nbsp;&nbsp;&nbsp;I&nbsp;&nbsp;&nbsp;직업정보제공사업신고 : J151602019001
						</td>
					</tr>
					<tr>
						<td style="text-align:center">
							<table border="0" cellspacing="0" cellpadding="0" style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center">
								<tr>
									<td style="width:266px">&nbsp;</td>
									<td><a href="https://www.ziggam.com" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-home.png" alt="HOME"></a></td>
									<td style="width:16px">&nbsp;</td>
									<td><a href="` + siteUrl + `" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-ems.png" alt="EMS"></a></td>
									<td style="width:16px">&nbsp;</td>
									<td><a href="` + googleStore + `" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-andapp.png" alt="Android"></a></td>
									<td style="width:16px">&nbsp;</td>
									<td><a href="` + appleStore + `" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-iosapp.png" alt="IOS"></a></td>
									<td style="width:266px">&nbsp;</td>
								</tr>
							</table>
						</td>
					</tr>
				</table>
			</td>
		</tr>
		<!-- //footer -->
		</table>
	</body>
	</html>`
*/
