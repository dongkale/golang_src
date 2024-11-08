package controllers

import (
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	"golang.org/x/net/context"
	"google.golang.org/api/option"

	"emst.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

type AdminEventProcessController struct {
	BaseController
}

func (c *AdminEventProcessController) Post() {

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
	pCuCd := c.GetString("cu_cd")   // 처리구분코드
	pSn := c.GetString("sn")        // 일련번호
	pGbnCd := c.GetString("gbn_cd") // 게시구분코드
	pTitle := c.GetString("title")  //제목
	notiDoc1 := c.GetString("notiDoc1")
	notiDoc2 := c.GetString("notiDoc2")
	notiDoc3 := c.GetString("notiDoc3")
	notiDoc4 := c.GetString("notiDoc4")
	notiDoc5 := c.GetString("notiDoc5")
	notiDoc6 := c.GetString("notiDoc6")
	notiDoc7 := c.GetString("notiDoc7")
	notiDoc8 := c.GetString("notiDoc8")
	notiDoc9 := c.GetString("notiDoc9")
	notiDoc10 := c.GetString("notiDoc10")
	notiDoc11 := c.GetString("notiDoc11")
	notiDoc12 := c.GetString("notiDoc12")
	notiDoc13 := c.GetString("notiDoc13")
	notiDoc14 := c.GetString("notiDoc14")
	notiDoc15 := c.GetString("notiDoc15")
	notiDoc16 := c.GetString("notiDoc16")
	notiDoc17 := c.GetString("notiDoc17")
	notiDoc18 := c.GetString("notiDoc18")
	notiDoc19 := c.GetString("notiDoc19")
	notiDoc20 := c.GetString("notiDoc20")
	notiDoc21 := c.GetString("notiDoc21")
	notiDoc22 := c.GetString("notiDoc22")
	notiDoc23 := c.GetString("notiDoc23")
	notiDoc24 := c.GetString("notiDoc24")
	notiDoc25 := c.GetString("notiDoc25")
	notiDoc26 := c.GetString("notiDoc26")
	notiDoc27 := c.GetString("notiDoc27")
	notiDoc28 := c.GetString("notiDoc28")
	notiDoc29 := c.GetString("notiDoc29")
	notiDoc30 := c.GetString("notiDoc30")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Admin Event Process
	log.Debug("CALL SP_EMS_EVENT_PROC( "+
		"'%v', '%v', %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pCuCd, pSn, pGbnCd, pTitle, notiDoc1, notiDoc2, notiDoc3, notiDoc4, notiDoc5, notiDoc6, notiDoc7, notiDoc8, notiDoc9, notiDoc10,
		notiDoc11, notiDoc12, notiDoc13, notiDoc14, notiDoc15, notiDoc16, notiDoc17, notiDoc18, notiDoc19, notiDoc20,
		notiDoc21, notiDoc22, notiDoc23, notiDoc24, notiDoc25, notiDoc26, notiDoc27, notiDoc28, notiDoc29, notiDoc30)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_EVENT_PROC( "+
		"'%v', '%v', %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pCuCd, pSn, pGbnCd, pTitle, notiDoc1, notiDoc2, notiDoc3, notiDoc4, notiDoc5, notiDoc6, notiDoc7, notiDoc8, notiDoc9, notiDoc10,
		notiDoc11, notiDoc12, notiDoc13, notiDoc14, notiDoc15, notiDoc16, notiDoc17, notiDoc18, notiDoc19, notiDoc20,
		notiDoc21, notiDoc22, notiDoc23, notiDoc24, notiDoc25, notiDoc26, notiDoc27, notiDoc28, notiDoc29, notiDoc30),
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

	rtnAdminEventProcess := models.RtnAdminEventProcess{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnAdminEventProcess = models.RtnAdminEventProcess{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}

		if rtnCd == 1 {
			if pCuCd == "C" {
				// 이벤트
				memNo := mem_no
				gbn := "10"
				val := mem_no
				gbn1 := pGbnCd

				go EventFCM(memNo, gbn, val, gbn1)
			}
		}
	}

	// End : Admin Event Process

	c.Data["json"] = &rtnAdminEventProcess
	c.ServeJSON()
}

func EventFCM(memNo interface{}, gbn string, val interface{}, gbn1 interface{}) {

	// start : log
	slog := logs.NewLogger()
	slog.SetLogger(logs.AdapterConsole)
	// end : log

	//session := c.StartSession()
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

	// Start : Certification Key Info

	slog.Debug("CALL SP_EMS_PUSH_KEY_INFO_S_R('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, memNo, gbn, val, gbn1)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_PUSH_KEY_INFO_S_R('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, memNo, gbn, val, gbn1),
		ora.S, /* PUSH_KEY */
		ora.S, /* CONT */
		ora.S, /* BRD_GBN_CD */
		ora.S, /* SN */
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
		token    string
		cont     string
		brdgbncd string
		sn       string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			token = procRset.Row[0].(string)
			cont = procRset.Row[1].(string)
			brdgbncd = procRset.Row[2].(string)
			sn = procRset.Row[3].(string)

			slog.Debug("token : %v", token)

			opt := option.WithCredentialsFile("qrate-2ee14-firebase-adminsdk-64reu-74554f5c44.json")
			app, err := firebase.NewApp(context.Background(), nil, opt)
			if err != nil {
				slog.Debug("error initializing app: %v\n", err)
			}

			// [START send_to_token_golang]
			// Obtain a messaging.Client from the App.
			ctx := context.Background()
			client, err := app.Messaging(ctx)

			// This registration token comes from the client FCM SDKs.
			registrationToken := token

			// See documentation on defining a message payload.
			message := &messaging.Message{
				Data: map[string]string{
					"type":     gbn,
					"title":    "[이벤트]직감에서 알려드립니다.",
					"body":     cont,
					"brdgbncd": brdgbncd,
					"sn":       sn,
				},
				Notification: &messaging.Notification{
					Title: "[이벤트]직감에서 알려드립니다.",
					Body:  cont,
				},
				Token: registrationToken,
			}

			// Send a message to the device corresponding to the provided
			// registration token.
			response, err := client.Send(ctx, message)
			if err != nil {
				slog.Debug("STATUS : ", err)
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
