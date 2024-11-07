package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"gopkg.in/rana/ora.v4"
)

type MessageSendController struct {
	beego.Controller
}

func (c *MessageSendController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()
	mem_no := session.Get(context.Background(), "mem_no")
	// if mem_no == nil {
	// 	c.Data["json"] = &models.RtnMessageSend{RtnCd:99,RtnMsg:"mem_no == nil"}
	// 	c.ServeJSON()
	// 	return
	// }

	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := mem_no
	//pEntpMemNo := c.GetString("mem_no")
	pRecrutSn := c.GetString("recrut_sn")
	pPpMemNo := c.GetString("pp_mem_no")
	pMsgGbnCd := c.GetString("msg_gbn_cd")
	pMsgCont := c.GetString("msg_cont")
	pLiveItvSdt := c.GetString("live_itv_sdt")
	pArrPpChrgSn := c.GetString("arr_pp_chrg_sn")
	pLiveSn := c.GetString("live_sn")
	//imgServer, _  := beego.AppConfig.String("viewpath")

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
		//emn := pEntpMemNo
		emn := pEntpMemNo.(string)
		rsn := pRecrutSn
		pmn := pPpMemNo
		lsn := pLiveSn

		if pMsgGbnCd == "02" { //대화종료
			gbn := "1004"
			go MessageSendFCM(memNo, gbn, val, emn, rsn, pmn, lsn, pMsgCont)
		} else if pMsgGbnCd == "03" { //다시 대화 시작
			gbn := "1005"
			go MessageSendFCM(memNo, gbn, val, emn, rsn, pmn, lsn, pMsgCont)
		} else if pMsgGbnCd == "06" { //라이브 인터뷰 요청
			gbn := "1006"
			go MessageSendFCM(memNo, gbn, val, emn, rsn, pmn, lsn, pMsgCont)
		} else if pMsgGbnCd == "09" { //라이브 인터뷰 확정 취소
			gbn := "1008"
			go MessageSendFCM(memNo, gbn, val, emn, rsn, pmn, lsn, pMsgCont)
		} else if pMsgGbnCd == "11" { //라이브 인터뷰 요청 취소
			gbn := "1007"
			go MessageSendFCM(memNo, gbn, val, emn, rsn, pmn, lsn, pMsgCont)
		} else if pMsgGbnCd == "99" { //일반대화
			gbn := "1003"
			go MessageSendFCM(memNo, gbn, val, emn, rsn, pmn, lsn, pMsgCont)
		} else if pMsgGbnCd == "04" { //합격
			gbn := "1002"
			go MessageSendFCM(memNo, gbn, val, emn, rsn, pmn, lsn, pMsgCont)
		} else if pMsgGbnCd == "05" { //불합격
			gbn := "1002"
			go MessageSendFCM(memNo, gbn, val, emn, rsn, pmn, lsn, pMsgCont)
		}
	}
	// End : Message Send Process

	c.Data["json"] = &rtnMessageSend
	c.ServeJSON()
}

func MessageSendFCM(memNo string, gbn string, val string, emn string, rsn string, pmn string, lsn string, pMsgCont string) {

	// start : log
	// slog := logs.NewLogger()
	// slog.SetLogger(logs.AdapterConsole)
	// end : log

	//session := c.StartSession()
	pLang, _ := beego.AppConfig.String("lang")
	gbn1 := "01"

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	fmt.Printf(fmt.Sprintf("[Push][Try] memNo:%v, gbn:%v, val:%v, emn:%v, rsn:%v, pmn:%v, lsn:%v, pMsgCont:%v", memNo, gbn, val, emn, rsn, pmn, lsn, pMsgCont))

	// Start : Push Key Info

	fmt.Printf(fmt.Sprintf("CALL ZMSP_PUSH_KEY_INFO_S_R('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, memNo, gbn, val, gbn1, emn, rsn, pmn, lsn))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZMSP_PUSH_KEY_INFO_S_R('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, memNo, gbn, val, gbn1, emn, rsn, pmn, lsn),
		ora.S,   /* PUSH_KEY */
		ora.S,   /* CONT */
		ora.S,   /* BRD_GBN_CD */
		ora.I64, /* SN */
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* LIVE_SN */
		ora.S,   /* push body */
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
		token string
		cont  string
		//brdGbnCd  string
		//sn        int64
		entpMemNo string
		recrutSn  string
		ppMemNo   string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			token = procRset.Row[0].(string)
			cont = procRset.Row[1].(string)
			//brdGbnCd = procRset.Row[2].(string)
			//sn = procRset.Row[3].(int64)
			entpMemNo = procRset.Row[4].(string)
			recrutSn = procRset.Row[5].(string)
			ppMemNo = procRset.Row[6].(string)

			// 합격 불합격은 채용 공고 내용을 바디로 해서 푸쉬 한다.
			if gbn == "1002" {
				pMsgCont = procRset.Row[8].(string)
			}

			if pMsgCont == "" {
				pMsgCont = cont
			}

			fmt.Printf("token : %v", token)
			//slog.Debug("brdGbnCd : %v", brdGbnCd)
			//slog.Debug("sn : %v", sn)

			opt := option.WithCredentialsFile("qrate-2ee14-firebase-adminsdk-64reu-74554f5c44.json")
			app, err := firebase.NewApp(context.Background(), nil, opt)
			if err != nil {
				fmt.Printf("error initializing app: %v\n", err)
			}

			// 비어 있다면 내용은 타이틀로 채워 준다.

			// [START send_to_token_golang]
			// Obtain a messaging.Client from the App.
			ctx := context.Background()
			client, err := app.Messaging(ctx)

			// This registration token comes from the client FCM SDKs.
			registrationToken := token

			// See documentation on defining a message payload.
			message := &messaging.Message{
				Data: map[string]string{
					"type":        gbn,
					"title":       "[직감] " + cont,
					"body":        pMsgCont,
					"entp_mem_no": entpMemNo,
					"recrut_sn":   recrutSn,
					"pp_mem_no":   ppMemNo,
				},
				Android: &messaging.AndroidConfig{
					/*
						Data: map[string]string{
							"type":        gbn,
							"title":       "[직감] " + cont,
							"body":        cont,
							"entp_mem_no": entpMemNo,
							"recrut_sn":   recrutSn,
							"pp_mem_no":   ppMemNo,
						},
							Notification: &messaging.AndroidNotification{
								Title: "[직감] " + cont,
								Body:  cont,
							},
					*/
				},
				APNS: &messaging.APNSConfig{
					/*
						Headers: map[string]string{
							"type":        gbn,
							"title":       "[직감] " + cont,
							"body":        cont,
							"entp_mem_no": entpMemNo,
							"recrut_sn":   recrutSn,
							"pp_mem_no":   ppMemNo,
						},
					*/
					Payload: &messaging.APNSPayload{
						Aps: &messaging.Aps{
							Alert: &messaging.ApsAlert{
								Title: "[직감] " + cont,
								Body:  pMsgCont,
							},
						},
					},
				},
				Token: registrationToken,
			}

			fmt.Printf(fmt.Sprintf("[Push][Complete] memNo:%v, gbn:%v, val:%v, emn:%v, rsn:%v, pmn:%v, lsn:%v, pMsgCont:%v -> Token:%v", memNo, gbn, val, emn, rsn, pmn, lsn, pMsgCont, token))

			// Send a message to the device corresponding to the provided
			// registration token.
			response, err := client.Send(ctx, message)
			if err != nil {
				fmt.Printf("STATUS : ", err)
			}
			// Response is a message ID string.
			fmt.Println("Successfully sent message:", response)
			// [END send_to_token_golang]
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
}
