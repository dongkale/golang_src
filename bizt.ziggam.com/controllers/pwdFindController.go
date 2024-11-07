package controllers

import (
	"crypto/tls"
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"bizt.ziggam.com/models"
	gomail "gopkg.in/gomail.v2"
	ora "gopkg.in/rana/ora.v4"
)

type PwdFindController struct {
	BaseController
}

func (c *PwdFindController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	c.TplName = "common/pwd_find.html"
}

func (c *PwdFindController) Post() {

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
	imgServer, _  := beego.AppConfig.String("viewpath")

	googleStore, _ := beego.AppConfig.String("googlestore")
	appleStore, _ := beego.AppConfig.String("applestore")
	mailTo, _ := beego.AppConfig.String("mailto")

	// sbsson
	smtpEmail := "support@ziggam.com" //beego.AppConfig.String("smtpEmail")
	emailPwd := "wlrrka0223^^!"       //beego.AppConfig.String("emailPwd")
	siteUrl, _ := beego.AppConfig.String("siteurl")

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

	fmt.Printf(fmt.Sprintf("CALL ZSP_FIND_PWD_STEP1_R('%v', '%v', '%v', '%v', :1)",
		pLang, pMemId, pPpChrgNm, pEmail))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_FIND_PWD_STEP1_R('%v', '%v', '%v', '%v', :1)",
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
		rtnCd    int64
		rtnMsg   string
		certNo   string
		entpKoNm string
	)

	rtnFindPwdStep1 := models.RtnFindPwdStep1{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
			certNo = procRset.Row[2].(string)
			entpKoNm = procRset.Row[3].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnFindPwdStep1 = models.RtnFindPwdStep1{
			RtnCd:    rtnCd,
			RtnMsg:   rtnMsg,
			CertNo:   certNo,
			EntpKoNm: entpKoNm,
		}

		//pEmail = "mrchang@qrate.co.kr"

		if rtnFindPwdStep1.RtnCd == 1 {
			var certNo = rtnFindPwdStep1.CertNo
			var entpKoNm = rtnFindPwdStep1.EntpKoNm
			go SendMailer(pEmail, imgServer, certNo, mailTo, googleStore, appleStore, smtpEmail, emailPwd, siteUrl, entpKoNm)
		}
	}
	// End : Certification Key Info

	c.Data["json"] = &rtnFindPwdStep1
	c.ServeJSON()
}

// SendMailer ...
func SendMailer(pEmail string, imgServer string, certNo string, mailTo string, googleStore string, appleStore string, smtpEmail string, emailPwd string, siteUrl string, entpKoNm string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "no-reply@ziggam.com")
	m.SetHeader("To", pEmail)
	m.SetHeader("Subject", "[직감:ZIGGAM] 비밀번호 재설정 인증번호입니다.")
	m.SetBody("text/html", `
	<html>
	<head>
	<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
	<title>직감 채용을 편하게, 면접을 영상으로</title>
	</head>
	<body style="margin:0;padding:0;background-color:#f5f6f9;">
	<table border="0" cellspacing="0" cellpadding="0" align="center" style="width:760px;margin:0 auto;font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif">
		<!-- header -->
		<tr>
			<td style="padding:76px 0px 20px;border-bottom:4px solid #2ad0c7;">
				<table border="0" cellspacing="0" cellpadding="0" style="width:100%;">
					<tr>
						<td>
							<a href="https://www.ziggam.com" target="_blank"><img style="display:block;border:0;" src="`+imgServer+`/mail/ic-main-logo.png" alt="직감 채용을 편하게,면접을 영상으로"></a>
						</td>
					</tr>
				</table>
			</td>
		</tr>
		<!-- //header -->
		<!-- contents -->
		<tr>
			<td style="padding:90px 70px;background-color:#ffffff;">
				<table border="0" cellspacing="0" cellpadding="0" style="width:100%;">
					<tr>
						<td style="text-align:left">
							<!-- 내용 -->
							<table border="0" cellspacing="0" cellpadding="0" style="width:100%">
								<!-- 메인 타이틀 -->
								<tr>
									<td style="width:100%;font-size:27px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:43px;color:#171717;letter-spacing:-1px">비밀번호 찾기 인증번호입니다.</td>
								</tr>
								<!-- 내용 텍스트 -->
								<tr>
									<td style="padding:70px 0 60px 0;width:100%;font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:25px;color:#171717;">안녕하세요, `+entpKoNm+`님<br>
									요청하신 비밀번호 찾기 인증번호를 안내 드립니다.<br>
									아래 인증 번호를 입력해주시면 재설정하실 수 있습니다. </td>
								</tr>
								<!-- 인증번호 -->
								<tr>
									<td style="padding:0 0 0 0">
										<table border="0" cellspacing="0" cellpadding="0" style="width:100%;background-color:#f8fafb;border-top:1px solid #eaeeef;border-right:1px solid #eaeeef;border-bottom:1px solid #eaeeef;border-left:1px solid #eaeeef;font-size:14px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center">
											<tr>
												<td style="text-align:center;width:100%;background-color:#f8fafb;font-size:22px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:70px;font-weight:bold;">`+certNo+`</td>
											</tr>
										</table>
									</td>
								</tr>
							</table>
							<!-- //내용 -->
						</td>
					</tr>
				</table>
			</td>
		</tr>
		<!-- //contents -->
		<!-- footer -->
		<tr>
			<td style="padding:20px 0px 20px 0;background-color:#ced3d6;">
				<table border="0" cellspacing="0" cellpadding="0" style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center">
					<tr>
						<td style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center;color:#4d5256">
							본 메일은 발신 전용 메일입니다. 문의 사항은 <a href="mailto:`+mailTo+`" style="color:#4d5256;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;">support@ziggam.com</a>을 이용해주세요.
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
									<td><a href="https://www.ziggam.com" target="_blank" style="border:0"><img src="`+imgServer+`/mail/btn-zinggam-home.png" alt="HOME"></a></td>
									<td style="width:16px">&nbsp;</td>
									<td><a href="`+siteUrl+`" target="_blank" style="border:0"><img src="`+imgServer+`/mail/btn-zinggam-ems.png" alt="EMS"></a></td>
									<td style="width:16px">&nbsp;</td>
									<td><a href="`+googleStore+`" target="_blank" style="border:0"><img src="`+imgServer+`/mail/btn-zinggam-andapp.png" alt="Android"></a></td>
									<td style="width:16px">&nbsp;</td>
									<td><a href="`+appleStore+`" target="_blank" style="border:0"><img src="`+imgServer+`/mail/btn-zinggam-iosapp.png" alt="IOS"></a></td>
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
	</html>
	`)
	d := gomail.NewDialer("smtp.gmail.com", 587, smtpEmail, emailPwd)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
