package controllers

import (
	"crypto/tls"
	"fmt"

	"emst.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	gomail "gopkg.in/gomail.v2"
	"gopkg.in/rana/ora.v4"
)

type CommonCertEmailResendController struct {
	beego.Controller
}

func (c *CommonCertEmailResendController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	//session := c.StartSession()
	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := c.GetString("mem_no")
	pEmail := c.GetString("email")

	imgServer, _ := beego.AppConfig.String("viewpath")
	siteUrl, _ := beego.AppConfig.String("siteurl")

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

	// Start : Email CertKey Resend

	log.Debug("CALL SP_EMS_CERT_EMAIL_RESEND_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pEmail)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_CERT_EMAIL_RESEND_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pEmail),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
		ora.S,   /* CERT_KEY */
		ora.S,   /* ENTP_KO_NM */
		ora.S,   /* ENTP_MEM_ID */
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
		emailCertKey string
		entpKoNm     string
		entpMemId    string
	)

	rtnCertEmailResend := models.RtnCertEmailResend{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
			emailCertKey = procRset.Row[2].(string)
			entpKoNm = procRset.Row[3].(string)
			entpMemId = procRset.Row[4].(string)

			log.Debug("rtnCd=%v", rtnCd)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnCertEmailResend = models.RtnCertEmailResend{
			RtnCd:        rtnCd,
			RtnMsg:       rtnMsg,
			EmailCertKey: emailCertKey,
			EntpKoNm:     entpKoNm,
			EntpMemId:    entpMemId,
		}

		if rtnCertEmailResend.RtnCd == 1 {
			m := gomail.NewMessage()
			m.SetHeader("From", "no-reply@ziggam.com")
			m.SetHeader("To", pEmail)
			m.SetHeader("Subject", "[직감:ZIGGAM] "+rtnCertEmailResend.EntpKoNm+"님의 기업회원 인증 메일입니다.")
			m.SetBody("text/html", `
			<html>
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=no, maximum-scale=1">
				<script type="text/javascript" src="//code.jquery.com/jquery-3.3.1.min.js"></script>
				<title>[직감]기업회원 이메일 인증</title>
				<style>
					body{font-family: 'Noto Sans KR', sans-serif; color:#808285;}
				</style>
			</head>
			<body>
				<div style="width:100%;padding:30px 0;margin:0;color:#333;">
				<table width="80%" cellspacing="0" cellpadding="0" bgcolor="#ffffff" style="margin:auto;"  summary="이메일 인증을 위한 안내 메일입니다.">
					<caption style="display:none;">기업회원 이메일 인증</caption>
					<thead>
						<tr>
							<th style="padding-top:20px; padding-bottom: 10px; border-bottom:4px solid #00CAED">
								<!-- //직감로고 -->
								<img src="`+imgServer+`/mail/ziggam_logo.png" width="15%">
							</th>
						</tr>
					</thead>
					<tbody>
						<tr>
							<td>					
								<div style="width:86%;margin:auto;padding:40px 15px 10px 15px;font-size:20px; text-align:center;">
									<strong><span style="color:#00CAED;">`+rtnCertEmailResend.EntpKoNm+`</span></strong>의 기업회원 가입을 위해<br>
									이메일 주소를 인증해 주세요.
								</div>
							</td>
						</tr>
						<tr>
							<td>
								<div style="width:80%;margin:auto;padding:10px 15px 30px 15px;font-size:14px; text-align: center;line-height: 1.2em;">
									<p>아래 링크를 눌러 이메일 인증을 완료해주세요. 이메일 인증이 완료되면 내부 인증을 거친 후 기업회원 가입이 완료됩니다.</p>
									<p><b>메일 인증을 완료하신 후에는 앱을 다시 시작해 주시기 바랍니다.</b></p>
								</div>					
							</td>
						</tr>
					</tbody>
				</table>
				<table width="200px" cellspacing="0" cellpadding="0" bgcolor="#ffffff" style="margin:auto;padding:0 30px 30px 30px;">
					<tbody>
						<tr>
							<td height="44" bgcolor="#00CAED" valign="middle" align="center" style="border-radius: 30px;">
								<a href="`+siteUrl+`/common/email/cert?mem_id=`+rtnCertEmailResend.EntpMemId+`&cert_key=`+rtnCertEmailResend.EmailCertKey+`" style="color:#fff;font-size:15px;text-decoration: none;">인증하기</a>
							</td>
						</tr>
					</tbody>
				</table>
			
				<!--footer-->
				<table width="80%" cellspacing="0" cellpadding="0" bgcolor="#ffffff" style="margin:auto;">
					<tbody>
						<tr>
							<th colspan="2" style="font-size:10px;font-weight:normal;padding:30px 15px; text-align: center;line-height: 1.2em;">
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
			d := gomail.NewDialer("smtp.gmail.com", 587, smtpEmail, emailPwd)
			d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

			// Send the email to Bob, Cora and Dan.
			if err := d.DialAndSend(m); err != nil {
				panic(err)
			}
		}
	} // End : Email CertKey Resend

	c.Data["json"] = &rtnCertEmailResend
	c.ServeJSON()
}
