package controllers

import (
	"fmt"

	"emst.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	"gopkg.in/rana/ora.v4"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

type RecruitPostInsertController struct {
	beego.Controller
}

func (c *RecruitPostInsertController) Post() {

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
	pEntpMemNo := mem_no
	pJobGrpCd := c.GetString("job_grp_cd")
	pEmplTypCd := c.GetString("empl_typ_cd")
	pRecrutCnt := c.GetString("recrut_cnt")
	pSex := c.GetString("sex")
	pRol := c.GetString("rol")
	pAplyQufct := c.GetString("aply_qufct")
	pRecrutTitle := c.GetString("recrut_title")
	pAcptPridCd := c.GetString("acpt_prid_cd")
	pPrgsMsg := c.GetString("prgs_msg")
	pAnsLmtTmCd := c.GetString("ans_lmt_tm_cd")
	pRecMaxTm := c.GetString("rec_max_tm")
	pArrQstTitle := c.GetString("qst_title_arr")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Recruit Insert Process

	log.Debug("CALL SP_EMS_RECRUIT_REG_PROC('%v', '%v', '%v', '%v',  %v , '%v', '%v', '%v', '%v', '%v', '%v', '%v', %v, '%v', :1)",
		pLang, pEntpMemNo, pJobGrpCd, pEmplTypCd, pRecrutCnt, pSex, pRol, pAplyQufct, pRecrutTitle, pAcptPridCd, pPrgsMsg, pAnsLmtTmCd, pRecMaxTm, pArrQstTitle)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_RECRUIT_REG_PROC('%v', '%v', '%v', '%v',  %v , '%v', '%v', '%v', '%v', '%v', '%v', '%v', %v, '%v', :1)",
		pLang, pEntpMemNo, pJobGrpCd, pEmplTypCd, pRecrutCnt, pSex, pRol, pAplyQufct, pRecrutTitle, pAcptPridCd, pPrgsMsg, pAnsLmtTmCd, pRecMaxTm, pArrQstTitle),
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

	rtnRecruitPostInsert := models.RtnRecruitPostInsert{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnRecruitPostInsert = models.RtnRecruitPostInsert{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
		// 채용공고등록 : 01
		memNo := pEntpMemNo
		gbn := "01"
		val := pEntpMemNo

		go RecruitRegFCM(memNo, gbn, val)
	}

	// End : Recruit Insert Process

	c.Data["json"] = &rtnRecruitPostInsert
	c.ServeJSON()
}

func RecruitRegFCM(memNo interface{}, gbn string, val interface{}) {

	// start : log
	slog := logs.NewLogger()
	slog.SetLogger(logs.AdapterConsole)
	// end : log

	//session := c.StartSession()
	pLang, _ := beego.AppConfig.String("lang")
	gbn1 := "00"

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
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			token = procRset.Row[0].(string)
			cont = procRset.Row[1].(string)

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
					"type":  gbn,
					"title": "[직감]채용공고 등록 알림",
					"body":  cont,
				},
				Notification: &messaging.Notification{
					Title: "[직감]채용공고 등록 알림",
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
