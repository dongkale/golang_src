package controllers

import (
	"crypto/tls"
	"fmt"

	gomail "gopkg.in/gomail.v2"

	"emst.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/rana/ora.v4"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

type EvaluateUpdateController struct {
	beego.Controller
}

func (c *EvaluateUpdateController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	//session := c.StartSession()
	pLang, _ := beego.AppConfig.String("lang")

	pGbnCd := c.GetString("gbn_cd")
	pEntpMemNo := c.GetString("entp_mem_no") //E2019011900006 다날
	pRecrutSn := c.GetString("recrut_sn")
	pPpMemNo := c.GetString("pp_mem_no")
	imgServer, _ := beego.AppConfig.String("viewpath")

	googleStore, _ := beego.AppConfig.String("googlestore")
	appleStore, _ := beego.AppConfig.String("applestore")

	smtpEmail, _ := beego.AppConfig.String("smtpEmail")
	emailPwd, _ := beego.AppConfig.String("emailPwd")
	mailTo, _ := beego.AppConfig.String("mailto")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Recruit Evaluation Update

	log.Debug("CALL SP_EMS_EVAL_UPT_PROC('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pGbnCd, pEntpMemNo, pRecrutSn, pPpMemNo)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_EVAL_UPT_PROC('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pGbnCd, pEntpMemNo, pRecrutSn, pPpMemNo),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
		ora.S,   /* ENTP_NM */
		ora.S,   /* EMAIL */
		ora.S,   /* BIZ_TPY */
		ora.S,   /* PRGS_MSG */
		ora.S,   /* EVL_STAT_DT */
		ora.S,   /* APPLY_DT */
		ora.S,   /* NM */
		ora.S,   /* M_EMAIL */
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
		entpNm    string
		email     string
		bizTpy    string
		prgsMsg   string
		evlStatDt string
		applyDt   string
		nm        string
		memail    string
	)

	rtnEvaluateUpdate := models.RtnEvaluateUpdate{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
			entpNm = procRset.Row[2].(string)
			email = procRset.Row[3].(string)
			bizTpy = procRset.Row[4].(string)
			prgsMsg = procRset.Row[5].(string)
			evlStatDt = procRset.Row[6].(string)
			applyDt = procRset.Row[7].(string)
			nm = procRset.Row[8].(string)
			memail = procRset.Row[9].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnEvaluateUpdate = models.RtnEvaluateUpdate{
			RtnCd:     rtnCd,
			RtnMsg:    rtnMsg,
			EntpNm:    entpNm,
			Email:     email,
			BizTpy:    bizTpy,
			PrgsMsg:   prgsMsg,
			EvlStatDt: evlStatDt,
			ApplyDt:   applyDt,
			Nm:        nm,
			Memail:    memail,
		}

		/*
			fmt.Println("EntpNm : ", rtnEvaluateUpdate.EntpNm)
			fmt.Println("Email : ", rtnEvaluateUpdate.Email)
			fmt.Println("BizTpy : ", rtnEvaluateUpdate.BizTpy)
			fmt.Println("PrgsMsg : ", rtnEvaluateUpdate.PrgsMsg)
			fmt.Println("EvlStatDt : ", rtnEvaluateUpdate.EvlStatDt)
			fmt.Println("ApplyDt : ", rtnEvaluateUpdate.ApplyDt)
			fmt.Println("Nm : ", rtnEvaluateUpdate.Nm)
			fmt.Println("Memail : ", rtnEvaluateUpdate.Memail)
		*/
		// 채용심사결과 : 02
		memNo := pPpMemNo
		gbn := "02"
		val := rtnEvaluateUpdate.Nm

		if pEntpMemNo != "E2019011900006" { //다날일 경우 푸시 발송 제외
			go EvalFCM(memNo, gbn, val)
		}

		if pEntpMemNo != "E2019011900006" { // 다날일 경우 메일 발송 제외
			if rtnEvaluateUpdate.RtnCd == 1 {
				m := gomail.NewMessage()
				m.SetHeader("From", "no-reply@ziggam.com")
				m.SetHeader("To", rtnEvaluateUpdate.Memail)
				m.SetHeader("Subject", "[직감:ZIGGAM] "+rtnEvaluateUpdate.Nm+" 님의 "+rtnEvaluateUpdate.EntpNm+" 입사지원 결과입니다.")

				if pGbnCd == "03" {
					m.SetBody("text/html", `
			<html>
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=no, maximum-scale=1">
				<link href="https://fonts.googleapis.com/css?family=Noto+Sans+KR:400,700" rel="stylesheet">
				<title>ZIGGAM::직감</title>
				<style>
					body{font-family: 'Noto Sans KR', sans-serif; color:#808285;}
				</style>
			</head>
			<body>
				<div style="width:100%;padding:30px 0;margin:0;color:#333;">
				<table width="80%" cellspacing="0" cellpadding="0" bgcolor="#ffffff" style="margin:auto;">
					<thead>
						<tr>
							<th style="padding-top:20px; padding-bottom: 10px; border-bottom:4px solid #00CAED">
								<!-- //직감로고 -->
								<img src="`+imgServer+`/mail/logo_color_RGB_200.png" width="15%">
							</th>
						</tr>
					</thead>
					<tbody>
					<tr>
						<td>					
							<div style="width:86%;margin:auto;padding:40px 15px 10px 15px;font-size:18px; text-align:center;">직감을 통해 지원하신 <strong><span style="color:#00EFD5;">`+rtnEvaluateUpdate.EntpNm+`</span></strong>에서 메세지를 보내셨습니다.
							</div>
						</td>
					</tr>
					<tr>
						<th style="padding:10px;">
							<table width="90%" border="0" cellspacing="0" cellpadding="0" style="margin:auto;">
								<tr>
									<!-- //축하합격 메세지 -->
									<td align="center">
										<div style="background-image: url(`+imgServer+`/mail/pop_bg@2x.png);background-size:100% 100%;width:100%;max-width:420px;min-width: 300px;">
											<div style="padding:5% 15% 15% 15%;width:70%; min-height:340px;">
												<p><img src="`+imgServer+`/mail/icon_balloons.png" width="25%"></p>
												<p style="font-size:18pt; font-weight: 400;color:#00CAED;">축하합니다.</p>
												<p style="font-size:14px;">`+rtnEvaluateUpdate.PrgsMsg+`</p>
												<p style="color:#ff5400;font-size:12px;">채용 진행 결정 : `+rtnEvaluateUpdate.EvlStatDt+`</p>
											</div>
										</div>
										<div style="width:84%;margin:20px auto;max-width:420px;text-align: center;line-height: 1.2;">
										<p style="font-size:10px;color:#a0a0a0;">채용 진행 결정일로부터 3일이 경과될때까지 회신이 없는 지원에 대해서는 해당 기업의 내규에 따라 결정이 취소될 수도 있습니다. </p>
										</div>
									</td>
								</tr>
							</table>
						</th>
					</tr>
				</tbody>
				</table>
				<table width="70%" cellspacing="0" cellpadding="0" bgcolor="#ffffff" style="margin:auto;padding:5px; font-size:15px;margin-top:20px;">
					<tbody>
						<tr>
							<td style="font-weight: 600;">지원내역</td>
						</tr>
					</tbody>
				</table>
				<table width="70%" cellspacing="0" cellpadding="0" bgcolor="#ffffff" style="margin:auto;border-bottom:1px solid #777777;border-top:1px solid #777777;padding:20px 0; font-size:14px;">
					<tbody>
						<tr>
							<td style="width:35%;padding-left:15px;">기업명</td>
							<td><strong><span style="color:#00CAED;">`+rtnEvaluateUpdate.EntpNm+`</span></strong></td>
						</tr>
						<tr>
							<td style="padding-left:15px;">업종</td>
							<td><strong><span style="color:#00CAED;">`+rtnEvaluateUpdate.BizTpy+`</span></strong></td>
						</tr>
						<tr>
							<td style="padding-left:15px;">지원날짜</td>
							<td>`+rtnEvaluateUpdate.ApplyDt+`</td>
						</tr>
					</tbody>
				</table>
				<!--footer-->
				<table width="80%" cellspacing="0" cellpadding="0" bgcolor="#ffffff" style="margin:auto;">
					<tbody>
						<tr>
							<th colspan="2" style="font-size:10px;font-weight:normal;padding:30px 15px; text-align: center;line-height: 1.2em;">
								<p style="margin-bottom:10px">본 메일은 가입 시 입력하신 메일로 발송되는 발신전용 메일입니다.</p>
								<p>문의사항은 <a href="#">support@ziggam.com</a>을 이용해 주세요.</p>
							</th>
						</tr>
					</tbody>
				</table>
				<table width="80%" cellspacing="0" cellpadding="0" bgcolor="#ffffff" style="margin:auto;">
					<tfoot style="background: #eee;">
						<tr>
							<td style="padding:15px 0 0 15px;">
								<img src="`+imgServer+`/mail/Logo-Qrate_grey_150.png" width="120px">
							</td>
							<td style="text-align:right;padding:15px 10px 0 0;">
								<span><a href="https://www.facebook.com/Guziggamsung" target="_blank"><img src="`+imgServer+`/mail/icon_facebook.png" style="height:30px;"></a></span>
								<span><a href="`+googleStore+`" target="_blank"><img src="`+imgServer+`/mail/icon_googlestore.png" style="height:30px;"></a></span>
								<span><a href="`+appleStore+`" target="_blank"><img src="`+imgServer+`/mail/icon_appstore.png" style="height:30px;"></span>
							</td>
						</tr>
						<tr>
							<td colspan="2" style="font-size:9px;font-weight:normal;padding:5px 15px 20px 20px;line-height: 1.2em; color:#a8a9ad">㈜큐레잇 대표:박혁재<br>
							경기도 성남시 분당구 황새울로 311번 길 36, 2층<br>
							+82 31 739 1130<br>
							<a href="https://www.ziggam.co/" target="_blank">www.ziggam.com</a><br>
							<p>Copyright © 2019. ㈜큐레잇 All rights reserved</p>
							</td>
							<td></td>
						</tr>
					</tfoot>
				</table>
			</body>
			</html>
			`)
				} else {
					m.SetBody("text/html", `
				<html>
				<head>
					<meta charset="UTF-8">
					<meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=no, maximum-scale=1">
					<link href="https://fonts.googleapis.com/css?family=Noto+Sans+KR:400,700" rel="stylesheet">
					<title>ZIGGAM::직감</title>
					<style>
						body{font-family: 'Noto Sans KR', sans-serif; color:#808285;}
					</style>
				</head>
				<body>
					<div style="width:100%;padding:30px 0;margin:0;color:#333;">
					<table width="80%" cellspacing="0" cellpadding="0" bgcolor="#ffffff" style="margin:auto;">
						<thead>
							<tr>
								<th style="padding-top:20px; padding-bottom: 10px; border-bottom:4px solid #00CAED">
									<!-- //직감로고 -->
									<img src="`+imgServer+`/mail/logo_color_RGB_200.png" width="15%">
								</th>
							</tr>
						</thead>
						<tbody>
						<tr>
							<td>					
								<div style="width:86%;margin:auto;padding:40px 15px 10px 15px;font-size:18px; text-align:center;">직감을 통해 지원하신 <strong><span style="color:#00CAED;">`+rtnEvaluateUpdate.EntpNm+`</span></strong>에서 메세지를 보내셨습니다.
								</div>
							</td>
						</tr>
						<tr>
							<th style="padding:10px;">
								<table width="90%" border="0" cellspacing="0" cellpadding="0" style="margin:auto;">
									<tr>
									<!-- //채용 포기 메세지 -->
									<td align="center">
										<div style="background-image: url(`+imgServer+`/mail/pop_bg@2x.png);background-size:100% 100%;width:100%;max-width:420px;min-width: 300px;">
											<div style="padding:5% 15% 15% 15%;width:70%; min-height:340px;">
												<p><img src="`+imgServer+`/mail/icon_down_gr.png" width="25%"></p>
												<p style="font-size:18pt; font-weight: 400;">죄송합니다.</p>
												<p style="font-size:14px;">귀하의 채용 진행을 포기하셨습니다. <br>다른 업체의 채용 결과를 기다려 주세요. <br>좋은 결과가 있으실 거에요.</p>
											</div>
										</div>
									</td>
									</tr>
								</table>
							</th>
						</tr>
					</tbody>
					</table>
					<table width="70%" cellspacing="0" cellpadding="0" bgcolor="#ffffff" style="margin:auto;padding:5px; font-size:15px;margin-top:20px;">
						<tbody>
							<tr>
								<td style="font-weight: 600;">지원내역</td>
							</tr>
						</tbody>
					</table>
					<table width="70%" cellspacing="0" cellpadding="0" bgcolor="#ffffff" style="margin:auto;border-bottom:1px solid #777777;border-top:1px solid #777777;padding:20px 0; font-size:14px;">
						<tbody>
							<tr>
								<td style="width:35%;padding-left:15px;">기업명</td>
								<td><strong><span style="color:#00CAED;">`+rtnEvaluateUpdate.EntpNm+`</span></strong></td>
							</tr>
							<tr>
								<td style="padding-left:15px;">업종</td>
								<td><strong><span style="color:#00CAED;">`+rtnEvaluateUpdate.BizTpy+`</span></strong></td>
							</tr>
							<tr>
								<td style="padding-left:15px;">지원날짜</td>
								<td>`+rtnEvaluateUpdate.ApplyDt+`</td>
							</tr>
						</tbody>
					</table>
					<!--footer-->
					<table width="80%" cellspacing="0" cellpadding="0" bgcolor="#ffffff" style="margin:auto;">
						<tbody>
							<tr>
								<th colspan="2" style="font-size:12px;font-weight:normal;padding:30px 15px; text-align: center;line-height: 1.2em;">
									<p style="margin-bottom:10px">본 메일은 가입 시 입력하신 메일로 발송되는 발신전용 메일입니다.</p>
									<p>문의사항은 <a href="`+mailTo+`">`+mailTo+`</a>을 이용해 주세요.</p>
								</th>
							</tr>
						</tbody>
					</table>
					<table width="80%" cellspacing="0" cellpadding="0" bgcolor="#ffffff" style="margin:auto;">
						<tfoot style="background: #eee;">
							<tr>
								<td style="padding:15px 0 0 15px;">
									<img src="`+imgServer+`/mail/Logo-Qrate_grey_150.png" width="120px">
								</td>
								<td style="text-align:right;padding:15px 10px 0 0;">
									<span><a href="https://www.facebook.com/Guziggamsung" target="_blank"><img src="`+imgServer+`/mail/icon_facebook.png" style="height:30px;"></a></span>
									<span><a href="`+googleStore+`" target="_blank"><img src="`+imgServer+`/mail/icon_googlestore.png" style="height:30px;"></a></span>
									<span><a href="`+appleStore+`" target="_blank"><img src="`+imgServer+`/mail/icon_appstore.png" style="height:30px;"></span>
								</td>
							</tr>
							<tr>
								<td colspan="2" style="font-size:9px;font-weight:normal;padding:5px 15px 20px 20px;line-height: 1.2em; color:#a8a9ad">㈜큐레잇 대표:박혁재<br>
								경기도 성남시 분당구 황새울로 311번 길 36, 2층<br>
								+82 31 739 1130<br>
								<a href="https://www.ziggam.co/" target="_blank">www.ziggam.com</a><br>
								<p>Copyright © 2019. ㈜큐레잇 All rights reserved</p>
								</td>
								<td></td>
							</tr>
						</tfoot>
					</table>
				</body>
				</html>
				`)
				}
				d := gomail.NewDialer("smtp.gmail.com", 587, smtpEmail, emailPwd)
				d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

				// Send the email to Bob, Cora and Dan.
				if err := d.DialAndSend(m); err != nil {
					panic(err)
				}
			}
		}
	}

	// End : Recruit Evaluation Update

	c.Data["json"] = &rtnEvaluateUpdate
	c.ServeJSON()
}

func EvalFCM(memNo string, gbn string, val string) {

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
					"title": "[직감]지원 결과 알림",
					"body":  cont,
				},
				Notification: &messaging.Notification{
					Title: "[직감]지원 결과 알림",
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
