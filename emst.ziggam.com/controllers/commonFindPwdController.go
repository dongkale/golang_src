package controllers

import (
	"crypto/tls"
	"fmt"

	"emst.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	gomail "gopkg.in/gomail.v2"
	"gopkg.in/rana/ora.v4"
)

type CommonFindPwdController struct {
	beego.Controller
}

func (c *CommonFindPwdController) Get() {
	c.TplName = "common/pwd_find.html"
}

func (c *CommonFindPwdController) Post() {
	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	//session := c.StartSession()
	pLang, _ := beego.AppConfig.String("lang")
	pMemId := c.GetString("mem_id")
	pPpChrgNm := c.GetString("pp_chrg_nm")
	pEmail := c.GetString("email")

	//siteUrl := beego.AppConfig.String("siteurl")
	imgServer, _ := beego.AppConfig.String("viewpath")

	googleStore, _ := beego.AppConfig.String("googlestore")
	appleStore, _ := beego.AppConfig.String("applestore")
	mailTo, _ := beego.AppConfig.String("mailto")

	smtpEmail, _ := beego.AppConfig.String("smtpEmail")
	emailPwd, _ := beego.AppConfig.String("emailPwd")

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

	log.Debug("CALL SP_EMS_FIND_PWD_STEP1_R('%v', '%v', '%v', '%v', :1)",
		pLang, pMemId, pPpChrgNm, pEmail)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_FIND_PWD_STEP1_R('%v', '%v', '%v', '%v', :1)",
		pLang, pMemId, pPpChrgNm, pEmail),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
		ora.S,   /* CERT_NO */
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
		certNo string
	)

	rtnFindPwdStep1 := models.RtnFindPwdStep1{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
			certNo = procRset.Row[2].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnFindPwdStep1 = models.RtnFindPwdStep1{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
			CertNo: certNo,
		}

		//pEmail = "mrchang@qrate.co.kr"

		if rtnFindPwdStep1.RtnCd == 1 {
			m := gomail.NewMessage()
			m.SetHeader("From", "no-reply@ziggam.com")
			m.SetHeader("To", pEmail)
			m.SetHeader("Subject", "[직감:ZIGGAM] 비밀번호 재설정 인증번호입니다.")
			m.SetBody("text/html", `
			<html>
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=no, maximum-scale=1">
				<link href="https://fonts.googleapis.com/css?family=Noto+Sans+KR:400,700" rel="stylesheet">
				<title>[직감]비밀번호 재설정 인증번호 발송</title>
				<style>
					body{font-family: 'Noto Sans KR', sans-serif; color:#808285;}
				</style>
			</head>
			<body>
				<div style="width:100%;padding:30px 0;margin:0;color:#333;">
				<table width="80%" cellspacing="0" cellpadding="0" bgcolor="#ffffff" style="margin:auto;"  summary="이메일 인증을 위한 안내 메일입니다.">
					<caption style="display:none;">비밀번호 재설정 인증번호 발송</caption>
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
								<strong>비밀번호 재설정 인증번호</strong>
							</div>
						</td>
					</tr>
					<tr>
						<td>
							<div style="width:80%;margin:auto;padding:10px 15px 30px 15px;font-size:14px; text-align: center;line-height: 1.2em;">
								<p>요청하신 비밀번호 재설정을 위한 인증코드입니다. <br>아래 인증코드를 사용하여 비밀번호를 재설정하실 수 있습니다</p>
							</div>					
						</td>
					</tr>
					</tbody>
				</table>
				<table width="250px" cellspacing="0" cellpadding="0" bgcolor="#ffffff" style="margin:auto;padding:0 30px 30px 30px;">
					<tbody>
						<tr>
							<td height="44" bgcolor="#eeeeee" valign="middle" align="center" style="border-radius: 30px;">
								<a href="javascript:void(0);" style="color:#000;font-size:15px;text-decoration: none;">`+rtnFindPwdStep1.CertNo+`</a>
							</td>
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
			d := gomail.NewDialer("smtp.gmail.com", 587, smtpEmail, emailPwd)
			d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

			// Send the email to Bob, Cora and Dan.
			if err := d.DialAndSend(m); err != nil {
				panic(err)
			}
		}
	}
	// End : Certification Key Info

	c.Data["json"] = &rtnFindPwdStep1
	c.ServeJSON()
}
