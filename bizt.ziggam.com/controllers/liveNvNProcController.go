package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"bizt.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/rana/ora.v4"
)

type LiveNvNProcController struct {
	beego.Controller
}

func (c *LiveNvNProcController) Post() {

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Data["json"] = &models.RtnMessageSend{RtnCd: 99, RtnMsg: "mem_no == nil"}
		c.ServeJSON()
		return
	}

	pLang, _ := beego.AppConfig.String("lang")
	//imgServer, _  := beego.AppConfig.String("viewpath")

	pEntpMemNo := mem_no

	pRecrutSn := c.GetString("recrut_sn")
	pMsgGbnCd := c.GetString("msg_gbn_cd")
	pMsgCont := c.GetString("msg_cont")
	pLiveItvSdt := c.GetString("live_itv_sdt")
	pArrPpChrgSn := c.GetString("arr_pp_chrg_sn")
	pLiveSn := c.GetString("live_sn")
	pApplyMemNoArr := c.GetString("apply_mem_no_arr")
	pPushApplyMemNoArr := c.GetString("push_apply_mem_sn_arr")
	iCheckApplyCnt, _ := c.GetInt("check_apply_cnt")
	iCheckMemCnt, _ := c.GetInt("check_mem_cnt")

	// var iCheckCnt int64
	// pCheckCnt := c.GetString("check_cnt")
	// if pCheckCnt == "" {
	// 	iCheckCnt = 0
	// } else {
	// 	iCheckCnt, _ = strconv.ParseInt(pCheckCnt, 10, 64)
	// }

	var retPushApplyMemList []models.LiveNvNApplyMemList
	err := json.Unmarshal([]byte(pPushApplyMemNoArr), &retPushApplyMemList)
	if err != nil {
		panic(err)
	}

	fmt.Printf(fmt.Sprintf("pEntpMemNo:%v, pRecrutSn:%v, pMsgGbnCd:%v, pMsgCont:%v, pLiveItvSdt:%v, pArrPpChrgSn:%v, pLiveSn:%v, pApplyMemNoArr:%v, pPushApplyMemNoArr:%v, iCheckApplyCnt:%v, iCheckMemCnt:%v",
		pEntpMemNo, pRecrutSn, pMsgGbnCd, pMsgCont, pLiveItvSdt, pArrPpChrgSn, pLiveSn, pApplyMemNoArr, pPushApplyMemNoArr, iCheckApplyCnt, iCheckMemCnt))

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : 다대다 요청/확정 취소/요청 취소
	fmt.Printf(fmt.Sprintf("CALL ZMSP_LIVE_NVN_MSG_SEND_PROC('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', %v, %v, :1)",
		pLang, pEntpMemNo, pRecrutSn, pApplyMemNoArr, pMsgGbnCd, pMsgCont, pLiveItvSdt, pArrPpChrgSn, pLiveSn, iCheckApplyCnt, iCheckMemCnt))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZMSP_LIVE_NVN_MSG_SEND_PROC('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', %v, %v, :1)",
		pLang, pEntpMemNo, pRecrutSn, pApplyMemNoArr, pMsgGbnCd, pMsgCont, pLiveItvSdt, pArrPpChrgSn, pLiveSn, iCheckApplyCnt, iCheckMemCnt),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
		ora.S,   /* RTN_DATA */
		ora.S,   /* RTN_DATA2 */
		ora.S,   /* RTN_DATA3 */
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
		rtnCd        int64
		rtnMsg       string
		procList     string
		liveSn       string
		rtnErrorCode string
	)

	rtnMessageSend := models.RtnMessageSend_v2{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
			procList = procRset.Row[2].(string)
			liveSn = procRset.Row[3].(string)
			rtnErrorCode = procRset.Row[4].(string)
		}

		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnMessageSend = models.RtnMessageSend_v2{
			RtnCd:        rtnCd,
			RtnMsg:       rtnMsg,
			RtnData:      procList,
			RtnLiveSn:    liveSn,
			RtnErrorCode: rtnErrorCode,
		}

		fmt.Printf(fmt.Sprintf("===> rtnCd:%v, rtnMsg:%v, procList:%v, liveSn:%v, RtnErrorCode:%v", rtnCd, rtnMsg, procList, liveSn, rtnErrorCode))

		var retSendData []models.LiveNvNApplyMemList

		if rtnCd == 1 {
			//applyMemList := strings.Split(pApplyMemNoArr, ",")
			if len(procList) > 0 {
				applyMemList := strings.Split(procList, ",")

				for _, val := range applyMemList {
					//for _, val := range applyMemList {
					//for _, val := range retPushApplyMemList {
					data := strings.Split(val, ";")
					if len(data) != 3 {
						fmt.Printf(fmt.Sprintf("[LiveItvNvN][Error] data:%v", data))
						continue
					}

					ppMemNo := data[0]
					ppMemNm := data[1]
					ppErrorCode := data[2]

					if ppErrorCode != "0" {
						fmt.Printf(fmt.Sprintf("[LiveItvNvN][Error] data:%v", data))
						continue
					}

					// if len(ppMemNo) == 0 || len(ppMemNm) == 0 {
					// 	fmt.Printf(fmt.Sprintf("[LiveItvNvN][Error] data:%v", data))
					// 	continue
					// }

					//ppMemNo := val.RslPpMemNo
					//ppMemNm := val.RslNm

					fmt.Printf(ppMemNo)
					fmt.Printf(ppMemNm)
					fmt.Printf(ppErrorCode)

					retSendData = append(retSendData, models.LiveNvNApplyMemList{
						RslPpMemNo: ppMemNo,
						RslNm:      ppMemNm,
					})

					// 메시지전송(공통)
					memNo := ppMemNo
					val := ppMemNm
					emn := pEntpMemNo.(string)
					rsn := pRecrutSn
					pmn := ppMemNo
					lsn := pLiveSn

					if pMsgGbnCd == "20" { //라이브 인터뷰 요청
						gbn := "1006" // api 푸시 번호
						go MessageSendFCM(memNo, gbn, val, emn, rsn, pmn, lsn, pMsgCont)
					} else if pMsgGbnCd == "21" { //라이브 인터뷰 확정 취소
						gbn := "1008"
						go MessageSendFCM(memNo, gbn, val, emn, rsn, pmn, lsn, pMsgCont)
					} else if pMsgGbnCd == "22" { //라이브 인터뷰 요청 취소
						gbn := "1007"
						go MessageSendFCM(memNo, gbn, val, emn, rsn, pmn, lsn, pMsgCont)
					}
				}
			}

			// if pMsgGbnCd == "20" { //라이브 인터뷰 요청
			// 	gbn := "1006" // api 푸시 번호
			// 	go LiveNvNSendPushFCM(pEntpMemNo.(string), pRecrutSn, pLiveSn, gbn, pMsgCont, retSendData)
			// } else if pMsgGbnCd == "21" { //라이브 인터뷰 확정 취소
			// 	gbn := "1008"
			// 	go LiveNvNSendPushFCM(pEntpMemNo.(string), pRecrutSn, pLiveSn, gbn, pMsgCont, retSendData)
			// } else if pMsgGbnCd == "22" { //라이브 인터뷰 요청 취소
			// 	gbn := "1007"
			// 	go LiveNvNSendPushFCM(pEntpMemNo.(string), pRecrutSn, pLiveSn, gbn, pMsgCont, retSendData)
			// }

		} else {
			fmt.Printf(fmt.Sprintf("[LiveItvNvN] Error rtnCd:%v, rtnMsg:%v", rtnCd, rtnMsg))
		}
	}
	// End : Message Send Process

	c.Data["json"] = &rtnMessageSend
	c.ServeJSON()
}

// LiveNvNSendPushFCM ...
func LiveNvNSendPushFCM(entpMemNo string, recrutSn string, liveSn string, gbn string, msgCont string, sendList []models.LiveNvNApplyMemList) {

	for _, val := range sendList {

		ppMemNo := val.RslPpMemNo
		ppMemNm := val.RslNm

		memNo := ppMemNo
		val := ppMemNm
		emn := entpMemNo
		rsn := recrutSn
		pmn := ppMemNo
		lsn := liveSn

		MessageSendFCM(memNo, gbn, val, emn, rsn, pmn, lsn, msgCont)
	}
}
