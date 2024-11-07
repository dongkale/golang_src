package controllers

import (
	"crypto/sha512"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"time"

	"bizt.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/disintegration/imaging"
	gomail "gopkg.in/gomail.v2"
	"gopkg.in/rana/ora.v4"
)

type JoinInsertController struct {
	beego.Controller
}

func (c *JoinInsertController) Post() {

	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	googleStore, _ := beego.AppConfig.String("googlestore")
	appleStore, _ := beego.AppConfig.String("applestore")
	imgServer, _  := beego.AppConfig.String("viewpath")

	// smtpEmail := beego.AppConfig.String("smtpEmail")
	// emailPwd := beego.AppConfig.String("emailPwd")

	smtpEmail := "support@ziggam.com"
	emailPwd := "wlrrka0223^^!"

	mailTo, _ := beego.AppConfig.String("mailto")
	siteUrl, _ := beego.AppConfig.String("siteurl")

	pLang, _ := beego.AppConfig.String("lang")
	pEntpKoNm := c.GetString("entp_ko_nm")
	pRepNm := c.GetString("rep_nm")
	pBizRegNo := c.GetString("biz_reg_no")
	pTelNo := c.GetString("tel_no")
	pSmsRecvYn := c.GetString("sms_recv_yn")
	pEmail := c.GetString("email")
	pEmailRecvYn := c.GetString("email_recv_yn")
	pInfoEqYn := c.GetString("info_eq_yn")
	pPpChrgNm := c.GetString("pp_chrg_nm")
	pPpChrgBpNm := c.GetString("pp_chrg_bp_nm")
	pPpChrgTelNo := c.GetString("pp_chrg_tel_no")
	pPpChrgSmsRecvYn := c.GetString("pp_chrg_sms_recv_yn")
	pPpChrgEmail := c.GetString("pp_chrg_email")
	pPpChrgEmailRecvYn := c.GetString("pp_chrg_email_recv_yn")
	pEntpId := c.GetString("entp_mem_id")
	pPwd := c.GetString("pwd")
	pZip := c.GetString("zip")
	pAddr := c.GetString("addr")
	pDtlAddr := c.GetString("dtl_addr")
	pRefAddr := c.GetString("ref_addr")
	pBizTpy := c.GetString("biz_tpy") // 업종 이름
	pbizCond := c.GetString("biz_cond")
	pEntpHtag1 := c.GetString("entp_htag1")
	pEntpHtag2 := c.GetString("entp_htag2")
	pEntpHtag3 := c.GetString("entp_htag3")
	pEntpIntr := c.GetString("entp_intr")
	pHomePg := c.GetString("home_pg")
	pEmpCnt := c.GetString("emp_cnt")
	pEstDy := c.GetString("est_dy")

	pBizRegYn := c.GetString("biz_reg_yn")
	pImgYn := c.GetString("img_yn")

	pEntpRegNoExt := c.GetString("entp_regno_ext")
	pEntpLogoExt := c.GetString("entp_logo_ext")

	// LDK 2020/08/24: 기업 정보 코드화 -->
	pBizTpyCd := c.GetString("biz_tpy_cd")             // 업종 코드
	pEntpProfile := c.GetString("entp_profile")        // 기업 소개
	pBizIntro := c.GetString("biz_intro")              // 사업 소개
	pEntpCapital := c.GetString("entp_capital")        // 자본금
	pEntpTotalSales := c.GetString("entp_total_sales") // 매출액
	pEntpTypeCd := c.GetString("entp_type_cd")         // 기업 형태 코드
	// pEntpType := c.GetString("entp_type")              // 기업 형태(대기업,중견기업,중소기업,공사/공기업,외국계기업,기타)
	pLocation := c.GetString("location")               // 소재지
	// <--

	pJobfair := c.GetString("jobfair") // LDK 2020/09/23: 박람회 등록 <-->

	sha := sha512.New()
	sha.Write([]byte(pPwd))
	sha_pwd := hex.EncodeToString(sha.Sum(nil))

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection
	
	// Start : Entp Join Process
	
	log.Debug(fmt.Sprintf("CALL ZSP_ENTP_REG_PROC_V2('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', %v , '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpKoNm, pRepNm, pBizRegNo, pTelNo, pSmsRecvYn, pEmail, pEmailRecvYn, pInfoEqYn, pPpChrgNm, pPpChrgBpNm, pPpChrgTelNo, pPpChrgSmsRecvYn, pPpChrgEmail, pPpChrgEmailRecvYn, pEntpId, sha_pwd, pZip, pAddr, pDtlAddr, pRefAddr, pBizTpy, pbizCond, pEntpHtag1, pEntpHtag2, pEntpHtag3, pEntpIntr, pHomePg, pEmpCnt, pEstDy, pBizTpyCd, pEntpProfile, pBizIntro, pEntpCapital, pEntpTotalSales, pEntpTypeCd, pLocation))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_ENTP_REG_PROC_V2('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v',  %v , '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpKoNm, pRepNm, pBizRegNo, pTelNo, pSmsRecvYn, pEmail, pEmailRecvYn, pInfoEqYn, pPpChrgNm, pPpChrgBpNm, pPpChrgTelNo, pPpChrgSmsRecvYn, pPpChrgEmail, pPpChrgEmailRecvYn, pEntpId, sha_pwd, pZip, pAddr, pDtlAddr, pRefAddr, pBizTpy, pbizCond, pEntpHtag1, pEntpHtag2, pEntpHtag3, pEntpIntr, pHomePg, pEmpCnt, pEstDy, pBizTpyCd, pEntpProfile, pBizIntro, pEntpCapital, pEntpTotalSales, pEntpTypeCd, pLocation),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
		ora.S,   /* SET_MEM_NO */
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
		setMemNo string
	)

	joinInsert := models.JoinInsert{}
	rtnJoinInsert := models.RtnJoinInsert{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)

			if rtnCd == 1 {
				setMemNo = procRset.Row[2].(string)

				joinInsert = models.JoinInsert{
					SetMemNo: setMemNo,
				}

				// 사업자등록증이 있을 경우
				if pBizRegYn == "Y" {
					// 사업자등록증 업로드
					nowDate := time.Now()
					dateFmt := fmt.Sprintf(nowDate.Format("20060102150405"))

					uploadPath, _ := beego.AppConfig.String("uploadpath")
					imgDir := uploadPath + "/biz/" + setMemNo

					// 폴더가 없을 경우 해당 폴더를 만들어준다.
					if _, err := os.Stat(imgDir); os.IsNotExist(err) {
						err = os.MkdirAll(imgDir, 0755)
						if err != nil {
							panic(err)
						}
					}

					// 사업자등록증 업로드
					log.Debug(fmt.Sprintf(imgDir+"/biz_%v_%v.%v", setMemNo, dateFmt, pEntpRegNoExt))
					// 원본이미지
					c.SaveToFile("entp_regno", fmt.Sprintf(imgDir+"/biz_%v_%v.%v", setMemNo, dateFmt, pEntpRegNoExt))

					bizFilePath := "/biz/" + setMemNo + "/biz_" + setMemNo + "_" + dateFmt + "." + pEntpRegNoExt

					// 사업자등록증 등록
					log.Debug(fmt.Sprintf("CALL ZSP_ENTP_REG_SUB_PROC( '%v', '%v', '%v', :1)",
						pLang, setMemNo, bizFilePath))

					stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_ENTP_REG_SUB_PROC( '%v', '%v', '%v', :1)",
						pLang, setMemNo, bizFilePath),
						ora.I64, ora.S)

					defer stmtProcCall.Close()
					if err != nil {
						panic(err)
					}
					procRset := &ora.Rset{}
					_, err = stmtProcCall.Exe(procRset)

					if err != nil {
						panic(err)
					}
				}

				// 기업로고가 있을 경우
				if pImgYn == "Y" {
					// 로고 업로드
					nowDate := time.Now()
					dateFmt := fmt.Sprintf(nowDate.Format("20060102150405"))

					uploadPath, _ := beego.AppConfig.String("uploadpath")
					imgDir := uploadPath + "/logo/" + setMemNo

					// 폴더가 없을 경우 해당 폴더를 만들어준다.
					if _, err := os.Stat(imgDir); os.IsNotExist(err) {
						err = os.MkdirAll(imgDir, 0755)
						if err != nil {
							panic(err)
						}
					}

					// 기업로고이미지 업로드
					log.Debug(fmt.Sprintf(imgDir+"/ori_%v.%v", setMemNo, pEntpLogoExt))
					// 원본이미지
					c.SaveToFile("entp_logo", fmt.Sprintf(imgDir+"/ori_%v.%v", setMemNo, pEntpLogoExt))

					oriLogoImgPath := "/logo/" + setMemNo + "/ori_" + setMemNo + "." + pEntpLogoExt
					logoImgPath := "/logo/" + setMemNo + "/" + setMemNo + "_" + dateFmt + "." + pEntpLogoExt

					// sbsson 이미지 리사이징
					height := 200
					n_w, n_h := GetImageSize(uploadPath + oriLogoImgPath)
					if n_h < 200 {
						height = n_h
					}
					ga := (n_h * 200) / n_w
					height = round(float64(ga))

					src, err := imaging.Open(uploadPath + oriLogoImgPath)
					if err != nil {
						log.Debug("Open failed: %v", err)
					}

					// 200 리사이징 이미지
					rszImg200 := imaging.Resize(src, 200, 0, imaging.Lanczos)

					// sbsson 이미지 리사이징
					dst := imaging.New(300, 300, color.RGBA{255, 255, 255, 255})
					dst = imaging.Paste(dst, rszImg200, image.Pt(50, 150-(height/2)))
					dst = imaging.Resize(dst, 200, 0, imaging.Lanczos)

					err = imaging.Save(dst, imgDir+"/"+setMemNo+"_"+dateFmt+"."+pEntpLogoExt)
					if err != nil {
						log.Debug("Save failed rszImg200: %v", err)
					}

					// 기업로고 이미지 등록
					log.Debug(fmt.Sprintf("CALL ZSP_ENTP_REG_SUB2_PROC( '%v', '%v', '%v', :1)",
						pLang, setMemNo, logoImgPath))

					stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_ENTP_REG_SUB2_PROC( '%v', '%v', '%v', :1)",
						pLang, setMemNo, logoImgPath),
						ora.I64, ora.S)

					defer stmtProcCall.Close()
					if err != nil {
						panic(err)
					}
					procRset := &ora.Rset{}
					_, err = stmtProcCall.Exe(procRset)

					if err != nil {
						panic(err)
					}
				}

				// LDK 2020/09/23: 박람회 등록 -->
				if pJobfair != "0" {
					log.Debug(fmt.Sprintf("CALL MSP_JOBFAIR_ENTP_REG_R( '%v', '%v', '%v', :1)",
						pLang, setMemNo, pJobfair))

					stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL MSP_JOBFAIR_ENTP_REG_R( '%v', '%v', '%v', :1)",
						pLang, setMemNo, pJobfair),
						ora.I64, ora.S)

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

					if procRset.IsOpen() {
						for procRset.Next() {
							rtnCd = procRset.Row[0].(int64)
							rtnMsg = procRset.Row[1].(string)
						}
					}

					log.Debug(fmt.Sprintf("CALL MSP_JOBFAIR_ENTP_REG_R( '%v', '%v', '%v', :1) --> rtnCd:%d, rtnMsg:'%v'",
						pLang, setMemNo, pJobfair, rtnCd, rtnMsg))
				}
				// <--
			} else {
				log.Debug(fmt.Sprintf("CALL ZSP_ENTP_REG_PROC_V2('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', %v , '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1) --> rtnCd:%d, rtnMsg:'%v'",
					pLang, pEntpKoNm, pRepNm, pBizRegNo, pTelNo, pSmsRecvYn, pEmail, pEmailRecvYn, pInfoEqYn, pPpChrgNm, pPpChrgBpNm, pPpChrgTelNo, pPpChrgSmsRecvYn, pPpChrgEmail, pPpChrgEmailRecvYn, pEntpId, sha_pwd, pZip, pAddr, pDtlAddr, pRefAddr, pBizTpy, pbizCond, pEntpHtag1, pEntpHtag2, pEntpHtag3, pEntpIntr, pHomePg, pEmpCnt, pEstDy, pBizTpyCd, pEntpProfile, pBizIntro, pEntpCapital, pEntpTotalSales, pEntpTypeCd, pLocation, rtnCd, rtnMsg))
			}
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnJoinInsert = models.RtnJoinInsert{
			RtnCd:   rtnCd,
			RtnMsg:  rtnMsg,
			RtnData: joinInsert,
		}

		log.Debug(fmt.Sprintf("[JoinInsertController] rtnJoinInsert --> rtnCd:%d, rtnMsg:'%v'", rtnCd, rtnMsg))

		if rtnJoinInsert.RtnCd == 1 {
			go SendJoin(pEmail, smtpEmail, emailPwd, mailTo, pEntpKoNm, imgServer, siteUrl, googleStore, appleStore)
		}
	}

	// End : Entp Join Process

	c.Data["json"] = &rtnJoinInsert
	c.ServeJSON()
}

func SendJoin(pEmail string, smtpEmail string, emailPwd string, mailTo string, pEntpKoNm string, imgServer string, siteUrl string, googleStore string, appleStore string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "no-reply@ziggam.com")
	m.SetHeader("To", pEmail)
	m.SetHeader("Subject", "[큐레잇(직감)] "+pEntpKoNm+"님의 기업회원 가입 메일입니다.")
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
									<td style="width:100%;font-size:27px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:43px;color:#171717;letter-spacing:-1px">`+pEntpKoNm+`의 <b>직감 회원 가입신청</b>이 완료되었습니다.</td>
								</tr>
								<!-- 내용 텍스트 -->
								<tr>
									<td style="padding:60px 0 10px 0;width:100%;font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:25px;color:#171717;">안녕하세요, `+pEntpKoNm+`님!<br>
									직감 서비스를 이용해주셔서 감사합니다.<br>
									직감은 구인과 구직의 안전을 위해 가입 정보를 검수하고 있습니다.<br>
									최대 3일까지 소용될 수 있으니, 조금만 기다려주세요.<br>
									가입 정보 검수가 완료되며 메일을 통해 알려드릴게요.<br>
									<br>
									감사합니다.<br>
									<br>
									직감 드림</td>
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
									<!-- <td><a href="`+siteUrl+`" target="_blank" style="border:0"><img src="`+imgServer+`/mail/btn-zinggam-ems.png" alt="EMS"></a></td>
									<td style="width:16px">&nbsp;</td> -->
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

// 이미지 썸네일
func ErrorCon(err error) {
	if err != nil {
		log.Print(err)
	}
}

func GetImageSize(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		ErrorCon(err)
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		ErrorCon(err)
	}
	defer file.Close()

	return image.Width, image.Height
}
