package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
)

type JobFairLiveRequestController struct {
	beego.Controller
}

func (c *JobFairLiveRequestController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	log.Debug("JobFairLiveApplyController")

	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := c.GetString("entp_mem_no")
	pRecrutSn := c.GetString("recrut_sn")
	pPpMemNo := c.GetString("pp_mem_no")
	pMsgGbnCd := c.GetString("msg_gbn_cd")
	pMsgCont := c.GetString("msg_cont")
	pLiveItvSdt := c.GetString("live_itv_sdt")
	pArrPpChrgSn := c.GetString("arr_pp_chrg_sn")
	pUrlStr := c.GetString("url_str")

	log.Debug(pLang)
	log.Debug(pEntpMemNo)
	log.Debug(pRecrutSn)
	log.Debug(pPpMemNo)
	log.Debug(pMsgGbnCd)
	log.Debug(pMsgCont)
	log.Debug(pLiveItvSdt)
	log.Debug(pArrPpChrgSn)
	log.Debug(pUrlStr)

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	/*
		// 지원자에게 라이브 신청 -->
		log.Debug("CALL ZSP_MSG_SEND_PROC_V2('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
			pLang, pEntpMemNo, pRecrutSn, pPpMemNo, pMsgGbnCd, pMsgCont, pLiveItvSdt, pArrPpChrgSn, "")

		log.Debug("CALL ZSP_MSG_SEND_PROC_V2('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
			pLang, pEntpMemNo, pRecrutSn, pPpMemNo, pMsgGbnCd, pMsgCont, pLiveItvSdt, pArrPpChrgSn, "")

		stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_MSG_SEND_PROC_V2('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
			pLang, pEntpMemNo, pRecrutSn, pPpMemNo, pMsgGbnCd, pMsgCont, pLiveItvSdt, pArrPpChrgSn, ""),
			ora.I64, ///* RTN_CD
			ora.S,   ///* RTN_MSG
			ora.S,   ///* RTN_DATA
			ora.S,   ///* RTN_LIVE_SN
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
			rtnCd     int64
			rtnMsg    string
			rtnData   string
			rtnLiveSn string
		)

		rtnMessageSend := models.RtnMessageSend_v2{}

		if procRset.IsOpen() {
			for procRset.Next() {
				rtnCd = procRset.Row[0].(int64)
				rtnMsg = procRset.Row[1].(string)
				rtnData = procRset.Row[2].(string)
				rtnLiveSn = procRset.Row[3].(string)
			}
			if err := procRset.Err(); err != nil {
				panic(err)
			}

			rtnMessageSend = models.RtnMessageSend_v2{
				RtnCd:     rtnCd,
				RtnMsg:    rtnMsg,
				RtnData:   rtnData,
				RtnLiveSn: rtnLiveSn,
			}
		}
	*/
	// <--

	//curl -X POST "localhost:8080/messenger/message/insert" --header "Content-Type: application/json" --data "{\"entp_mem_no\":\"E2018102500001\",\"recrut_sn\":\"2020081406\",\"my_mem_no\":\"P2020070100701\",\"msg_gbn_cd\":\"07\",\"msg_cont\":\"\",\"live_sn\":\"202008071443001\"}"

	var rtnLiveSn string = "202008071443001"

	var resultString string
	// s1 = fmt.Sprintf("curl -X POST \"localhost:8080/messenger/message/insert\" --header \"Content-Type: application/json\" --data \"{\\\"entp_mem_no\\\":\\\"%v\\\",\\\"recrut_sn\\\":\\\"%v\\\",\\\"my_mem_no\\\":\\\"%v\\\",\\\"msg_gbn_cd\\\":\\\"%v\\\",\\\"msg_cont\\\":\\\"%v\\\",\\\"live_sn\\\":\\\"%v\\\"}\""),
	// 			pEntpMemNo, pRecrutSn, pPpMemNo, pMsgGbnCd, pMsgCont, "")
	//beego.Deb

	//s1 = fmt.Sprintf("curl -X POST \"localhost:8080/messenger/message/insert\" --header \"Content-Type: application/json\" --data \"{\\\"entp_mem_no\\\":\\\"E2018102500001\\\",\\\"recrut_sn\\\":\\\"2020081406\\\",\\\"my_mem_no\\\":\\\"P2020070100701\\\",\\\"msg_gbn_cd\\\":\\\"07\\\",\\\"msg_cont\\\":\\\"Test\\\",\\\"live_sn\\\":\\\"202008071443001\\\"}\"")
	resultString = fmt.Sprintf("curl -X POST \"%s/messenger/message/insert\" --header \"Content-Type: application/json\" --data \"{\\\"entp_mem_no\\\":\\\"%s\\\",\\\"recrut_sn\\\":\\\"%s\\\",\\\"my_mem_no\\\":\\\"%s\\\",\\\"msg_gbn_cd\\\":\\\"%s\\\",\\\"msg_cont\\\":\\\"%s\\\",\\\"live_sn\\\":\\\"%s\\\"}\"", pUrlStr, pEntpMemNo, pRecrutSn, pPpMemNo, pMsgGbnCd, pMsgCont, rtnLiveSn)

	log.Debug(resultString)

	//curl -X POST "localhost:8080/messenger/message/insert" --header "Content-Type: application/json" --data "{\"entp_mem_no\":\"E2018102500001\",\"recrut_sn\":\"2020081406\",\"my_mem_no\":\"P2020070100701\",\"msg_gbn_cd\":\"07\",\"msg_cont\":\"\",\"live_sn\":\"202008071443001\"}"

	// memNo := pPpMemNo
	// val := "rtnData"
	// emn := pEntpMemNo
	// rsn := pRecrutSn
	// pmn := pPpMemNo
	// lsn := rtnLiveSn
	// gbn := "1006"

	// memNo := "P2020070100701"
	// val := "rtnData"
	// emn := "E2018102500001"
	// rsn := "2020081422"
	// pmn := "P2020070100701"
	// lsn := "202008071443001"
	// gbn := "2002"

	//go MessageSendFCM(memNo, gbn, val, emn, rsn, pmn, lsn)

	//c.Data["json"] = &rtnMessageSend
	c.ServeJSON()
}
