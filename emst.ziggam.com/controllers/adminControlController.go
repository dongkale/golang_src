package controllers

import (
	"crypto/tls"
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"

	"emst.ziggam.com/models"
	gomail "gopkg.in/gomail.v2"
	ora "gopkg.in/rana/ora.v4"
)

type AdminControlController struct {
	beego.Controller
}

func (c *AdminControlController) Post() {

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
	pEntpMemNo := c.GetString("entp_mem_no") // 기업회원번호
	pGbnCd := c.GetString("gbn_cd")          // 제재구분코드
	pChk := c.GetString("chk")               // 제재유형

	// sbsson
	Email := c.GetString("Email")       // 대표자 이메일
	EntpKoNm := c.GetString("EntpKoNm") // 기업명
	BizRegNo := c.GetString("BizRegNo") // 부서/직책
	RepNm := c.GetString("RepNm")       // 대표자명

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Admin Control Process
	log.Debug("CALL SP_EMS_ADMIN_CONTROL_PROC('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pGbnCd, pChk)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_CONTROL_PROC('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pGbnCd, pChk),
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

	rtnAdminControl := models.RtnAdminControl{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnAdminControl = models.RtnAdminControl{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}

	// sbsson
	// Start : 계정관리자 정보
	var pp_chrg_sn = "0001"
	log.Debug("CALL ZSP_TEAM_MEM_INFO_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pp_chrg_sn)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_TEAM_MEM_INFO_R('%v', '%v', '%v', :1)", pLang, pEntpMemNo, pp_chrg_sn),
		ora.S, /* PP_CHRG_NM */
		ora.S, /* PP_CHRG_BP_NM */
		ora.S, /* PP_CHRG_TEL_NO */
		ora.S, /* SMS_RECV_YN */
		ora.S, /* EMAIL */
		ora.S, /* EMAIL_RECV_YN */
		ora.S, /* ENTP_MEM_ID */
	)
	defer stmtProcCall.Close()
	if err != nil {
		panic(err)
	}
	procRset = &ora.Rset{}
	_, err = stmtProcCall.Exe(procRset)

	if err != nil {
		panic(err)
	}

	var (
		PP_CHRG_NM     string
		PP_CHRG_BP_NM  string
		PP_CHRG_TEL_NO string
		EMAIL          string
		ENTP_MEM_ID    string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			PP_CHRG_NM = procRset.Row[0].(string)
			PP_CHRG_BP_NM = procRset.Row[1].(string)
			PP_CHRG_TEL_NO = procRset.Row[2].(string)
			//SMS_RECV_YN = procRset.Row[3].(string)
			EMAIL = Email //procRset.Row[4].(string)
			//EMAIL_RECV_YN = procRset.Row[5].(string)
			ENTP_MEM_ID = procRset.Row[6].(string)
		}

		if err := procRset.Err(); err != nil {
			panic(err)
		}

	}

	// 기업 제재 및 해제 메일 발송
	googleStore, _ := beego.AppConfig.String("googlestore")
	appleStore, _ := beego.AppConfig.String("applestore")
	imgServer, _ := beego.AppConfig.String("viewpath")

	//smtpEmail := beego.AppConfig.String("smtpEmail")
	//emailPwd := beego.AppConfig.String("emailPwd")
	mailTo, _ := beego.AppConfig.String("mailto")
	siteUrl, _ := beego.AppConfig.String("siteurl")

	pEmail := EMAIL

	imgServer = "https://biz.ziggam.com/static/upload"
	if pChk == "0" { //해제
		go SendOk(pEmail, mailTo, "wlrrka0223^^!", mailTo, pEntpMemNo, imgServer, siteUrl, googleStore, appleStore,
			EntpKoNm, BizRegNo, RepNm, PP_CHRG_NM, PP_CHRG_BP_NM, PP_CHRG_TEL_NO, EMAIL, ENTP_MEM_ID)
	} /*else if pChk == "1" { //제재
		go SendBlock(pEmail, "help@qrate.co.kr", "yaoncyfxoroukvzb", mailTo, pEntpMemNo, imgServer, siteUrl, googleStore, appleStore)
	} else {
	}*/
	// End : Admin Control Process

	c.Data["json"] = &rtnAdminControl
	c.ServeJSON()
}

// 기업 제재 및 해제 메일 발송
func SendOk(pEmail string, smtpEmail string, emailPwd string, mailTo string, pEntpKoNm string, imgServer string, siteUrl string, googleStore string, appleStore string,
	EntpKoNm string, BizRegNo string, RepNm string, PP_CHRG_NM string, PP_CHRG_BP_NM string, PP_CHRG_TEL_NO string, EMAIL string, ENTP_MEM_ID string) {

	m := gomail.NewMessage()
	m.SetHeader("From", "support@ziggam.com")
	m.SetHeader("To", pEmail)
	m.SetHeader("Subject", "[큐레잇(직감)] "+EntpKoNm+"님의 기업정보 승인 완료 안내 메일입니다.")
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
									<td style="width:100%;font-size:27px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:43px;color:#171717;letter-spacing:-1px">`+EntpKoNm+`의 <b>직감 회원 가입승인</b>이 완료되었습니다.</td>
								</tr>
								<!-- 내용 텍스트 -->
								<tr>
									<td style="padding:70px 0 100px 0;width:100%;font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:25px;color:#171717;"><b>`+EntpKoNm+`</b>의 직감 가입을 환영합니다.<br>
									직감은 국내 최초 영상 기반의 채용 서비스를 제공하는 플랫폼입니다.<br>
									이제 직감을 통해 채용을 편하게, 면접을 영상으로 진행해 보세요.<br>
									직감에 오는 지원자를 만나실 수 있을거예요!
									</td>
								</tr>
								<!-- 테이블 제목 -->
								<tr>
									<td style="padding:0px 0 20px 0;font-weight:bold;width:100%;font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:25px;color:#171717;">기업 정보</td>
								</tr>
								<!-- 테이블 내용 -->
								<tr>
									<td>
										<table border="0" cellspacing="0" cellpadding="0" style="width:100%;background-color:#fff;border-top:1px solid #eaeeef;border-bottom:1px solid #eaeeef;font-size:14px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;">
											<tr>
												<th style="width:170px;padding:13px 0 13px 24px;background-color:#f8fafb;font-weight:normal;color:#171717;text-align:left;word-break:break-all;white-space:normal;vertical-align:top">기업명</th>
												<td style="padding:13px 24px;color:#171717;word-break:break-all;white-space:normal;line-height:22px">`+EntpKoNm+`</td>
											</tr>
											<tr>
												<th style="width:170px;padding:13px 0 13px 24px;background-color:#f8fafb;font-weight:normal;color:#171717;text-align:left;word-break:break-all;white-space:normal;vertical-align:top;border-top:1px solid #eaeeef">사업자등록번호</th>
												<td style="padding:13px 24px;color:#171717;word-break:break-all;white-space:normal;line-height:22px;border-top:1px solid #eaeeef">`+BizRegNo+`</td>
											</tr>
											<tr>
												<th style="width:170px;padding:13px 0 13px 24px;background-color:#f8fafb;font-weight:normal;color:#171717;text-align:left;word-break:break-all;white-space:normal;vertical-align:top;border-top:1px solid #eaeeef">대표자명</th>
												<td style="padding:13px 24px;color:#171717;word-break:break-all;white-space:normal;line-height:22px;border-top:1px solid #eaeeef">`+RepNm+`</td>
											</tr>
										</table>
									</td>
								</tr>
								<!-- 참고 내용 -->
								<tr>
									<td style="padding:15px 0 15px 0;width:100%;font-size:13px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:25px;color:#4d5256;">
									・ 기업명과 사업자등록번호는 가입 후에는 변경이 불가능합니다.<br>
									</td>
								</tr>

								<!-- 테이블 제목 -->
								<tr>
									<td style="padding:70px 0 20px 0;font-weight:bold;width:100%;font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:25px;color:#171717;">계정 정보</td>
								</tr>
								<!-- 테이블 내용 -->
								<tr>
									<td>
										<table border="0" cellspacing="0" cellpadding="0" style="width:100%;background-color:#fff;border-top:1px solid #eaeeef;border-bottom:1px solid #eaeeef;font-size:14px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;">
											<tr>
												<th style="width:170px;padding:13px 0 13px 24px;background-color:#f8fafb;font-weight:normal;color:#171717;text-align:left;word-break:break-all;white-space:normal;vertical-align:top">이름</th>
												<td style="padding:13px 24px;color:#171717;word-break:break-all;white-space:normal;line-height:22px">`+PP_CHRG_NM+`</td>
											</tr>
											<tr>
												<th style="width:170px;padding:13px 0 13px 24px;background-color:#f8fafb;font-weight:normal;color:#171717;text-align:left;word-break:break-all;white-space:normal;vertical-align:top;border-top:1px solid #eaeeef">부서・직책</th>
												<td style="padding:13px 24px;color:#171717;word-break:break-all;white-space:normal;line-height:22px;border-top:1px solid #eaeeef">`+PP_CHRG_BP_NM+`</td>
											</tr>
											<tr>
												<th style="width:170px;padding:13px 0 13px 24px;background-color:#f8fafb;font-weight:normal;color:#171717;text-align:left;word-break:break-all;white-space:normal;vertical-align:top;border-top:1px solid #eaeeef">전화번호</th>
												<td style="padding:13px 24px;color:#171717;word-break:break-all;white-space:normal;line-height:22px;border-top:1px solid #eaeeef">`+PP_CHRG_TEL_NO+`</td>
											</tr>
											<tr>
												<th style="width:170px;padding:13px 0 13px 24px;background-color:#f8fafb;font-weight:normal;color:#171717;text-align:left;word-break:break-all;white-space:normal;vertical-align:top;border-top:1px solid #eaeeef">이메일</th>
												<td style="padding:13px 24px;color:#171717;word-break:break-all;white-space:normal;line-height:22px;border-top:1px solid #eaeeef">`+EMAIL+`</td>
											</tr>
											<tr>
												<th style="width:170px;padding:13px 0 13px 24px;background-color:#f8fafb;font-weight:normal;color:#171717;text-align:left;word-break:break-all;white-space:normal;vertical-align:top;border-top:1px solid #eaeeef">아이디</th>
												<td style="padding:13px 24px;color:#171717;word-break:break-all;white-space:normal;line-height:22px;border-top:1px solid #eaeeef">`+ENTP_MEM_ID+`</td>
											</tr>																						
										</table>
									</td>
								</tr>
								<!-- 참고 내용 -->
								<tr>
									<td style="padding:15px 0 15px 0;width:100%;font-size:13px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:25px;color:#4d5256;">
									・ 위 계정은 기업의 대표 관리자(어드민)계정으로 사용됩니다.<br>
									・ 팀 멤버 추가 및 어드민 권한 양도는 회원가입 완료 후 기업 관리자 페이지를 이용해 주세요.<br>
									<b>・ 자세한 내용은 서비스 소개서 및 기업관리자 매뉴얼을 확인해 주세요.</b><br>
									</td>
								</tr>
								<!-- 버튼 -->
								<tr>
									<td style="text-align:center;padding:45px 0 10px"><a href="https://biz.ziggam.com" style="border:0"><img src="`+imgServer+`/mail/btn-ems-go-default.png" alt="기업관리자 바로가기"></a></td>
								</tr>

								<tr>
                                <td style="padding:20px 0 0 0;width:100%;font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:25px;color:#171717;"><a href="https://drive.google.com/file/d/1dIqwygK7-47rDxnV-E0E_K7C_wKU5FTz/view?usp=gmail" style="border:0">기업관리자 매뉴얼 다운로드</a></td>
								</tr>
								<tr>
                                <td style="padding:20px 0 0 0;width:100%;font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:25px;color:#171717;"><a href="https://drive.google.com/file/d/1X5EHL6oNkxJ-G-NA01Rdk5ev0-fd2Au_/view?usp=gmail" style="border:0">서비스 제안서 다운로드</a></td>
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
									<td><a href="https://www.ziggam.com/" style="border:0"><img src="`+imgServer+`/mail/btn-zinggam-home.png" alt="HOME"></a></td>
									<td style="width:16px">&nbsp;</td>
									<td><a href="https://biz.ziggam.com" style="border:0"><img src="`+imgServer+`/mail/btn-zinggam-ems.png" alt="EMS"></a></td>
									<td style="width:16px">&nbsp;</td>
									<td><a href="`+googleStore+`" style="border:0"><img src="`+imgServer+`/mail/btn-zinggam-andapp.png" alt="Android"></a></td>
									<td style="width:16px">&nbsp;</td>
									<td><a href="`+appleStore+`" style="border:0"><img src="`+imgServer+`/mail/btn-zinggam-iosapp.png" alt="IOS"></a></td>
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

	//m.Attach(`static\upload\test\[직감]biz 기업 관리자 매뉴얼.pdf`)
	//m.Attach(`static\upload\test\[직감]서비스 제안서.pdf`)
	//m.Attach(`https://www.youtube.com/watch?v=MofE_8ecdkY&authuser=0`)
	//m.Attach(`https://www.youtube.com/watch?v=MxEz00rINRI&authuser=0`)

	d := gomail.NewDialer("smtp.gmail.com", 587, smtpEmail, emailPwd)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func SendBlock(pEmail string, smtpEmail string, emailPwd string, mailTo string, pEntpKoNm string, imgServer string, siteUrl string, googleStore string, appleStore string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "support@ziggam.com")
	m.SetHeader("To", pEmail)
	m.SetHeader("Subject", "[큐레잇(직감)] 기업 정보 승인 취소 안내")
	m.SetBody("text/html", `
	<html>
	<head>
	<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
	<title>직감 채용을 편하게, 면접을 영상으로</title>
	</head>
	<body style="margin:0;padding:0;background-color:#f5f6f9;">

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
