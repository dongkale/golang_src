package controllers

import (
	"bytes"
	"fmt"

	"bizt.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
	ora "gopkg.in/rana/ora.v4"
)

type ApiMessageSendController struct {
	beego.Controller
}

func (c *ApiMessageSendController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	//session := c.StartSession()
	fmt.Printf("ApiMessageSendController")

	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := c.GetString("entp_mem_no")
	pRecrutSn := c.GetString("recrut_sn")
	pPpMemNo := c.GetString("pp_mem_no")
	pMsgGbnCd := c.GetString("msg_gbn_cd")
	pMsgCont := c.GetString("msg_cont")
	pLiveItvSdt := c.GetString("live_itv_sdt")
	pArrPpChrgSn := c.GetString("arr_pp_chrg_sn")
	pLiveSn := c.GetString("live_sn")

	// 한글 변환
	var bufs bytes.Buffer
	wr := transform.NewWriter(&bufs, korean.EUCKR.NewDecoder())
	wr.Write([]byte(pMsgCont))
	wr.Close()

	pMsgCont = bufs.String()

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Message Send Process

	// LDK: 합격 처리, 합격 메세지 등록(ZSP_MSG_SEND_PROC)
	fmt.Printf(fmt.Sprintf("CALL ZSP_MSG_SEND_PROC('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo, pMsgGbnCd, pMsgCont, pLiveItvSdt, pArrPpChrgSn, pLiveSn))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_MSG_SEND_PROC('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo, pMsgGbnCd, pMsgCont, pLiveItvSdt, pArrPpChrgSn, pLiveSn),
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
		rtnCd   int64
		rtnMsg  string
		rtnData string
	)

	rtnMessageSend := models.RtnMessageSend{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
			rtnData = procRset.Row[2].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnMessageSend = models.RtnMessageSend{
			RtnCd:   rtnCd,
			RtnMsg:  rtnMsg,
			RtnData: rtnData,
		}

		// 메시지전송(공통)
		memNo := pPpMemNo
		val := rtnData
		emn := pEntpMemNo
		rsn := pRecrutSn
		pmn := pPpMemNo
		lsn := pLiveSn

		if pMsgGbnCd == "02" { //대화종료
			gbn := "1004"
			go MessageSendFCM(memNo, gbn, val, emn, rsn, pmn, lsn, "")
		} else if pMsgGbnCd == "03" { //다시 대화 시작
			gbn := "1005"
			go MessageSendFCM(memNo, gbn, val, emn, rsn, pmn, lsn, "")
		} else if pMsgGbnCd == "06" { //라이브 인터뷰 요청
			gbn := "1006"
			go MessageSendFCM(memNo, gbn, val, emn, rsn, pmn, lsn, "")
		} else if pMsgGbnCd == "09" { //라이브 인터뷰 확정 취소
			gbn := "1008"
			go MessageSendFCM(memNo, gbn, val, emn, rsn, pmn, lsn, "")
		} else if pMsgGbnCd == "11" { //라이브 인터뷰 요청 취소
			gbn := "1007"
			go MessageSendFCM(memNo, gbn, val, emn, rsn, pmn, lsn, "")
		} else if pMsgGbnCd == "99" { //일반대화
			gbn := "1003"
			go MessageSendFCM(memNo, gbn, val, emn, rsn, pmn, lsn, "")
		} else if pMsgGbnCd == "04" { //합격
			gbn := "1002"
			go MessageSendFCM(memNo, gbn, val, emn, rsn, pmn, lsn, "")
		} else if pMsgGbnCd == "05" { //불합격
			gbn := "1002"
			go MessageSendFCM(memNo, gbn, val, emn, rsn, pmn, lsn, "")
		}
	}
	// End : Message Send Process

	c.Data["json"] = &rtnMessageSend
	c.ServeJSON()
}
