package controllers

import (
	"fmt"
	"net/url"

	"bizt.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/rana/ora.v4"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

type RecruitInsertController struct {
	beego.Controller
}

func (c *RecruitInsertController) Post() {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/login")
	}

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := mem_no
	pJobGrpCd := c.GetString("job_grp_cd")
	pRecrutGbnCd := c.GetString("recrut_gbn_cd")
	pRecrutCnt := c.GetString("recrut_cnt")
	pRol := c.GetString("rol")                // 역할
	pAplyQufct := c.GetString("aply_qufct")   // 지원자격
	pPerferTrtm := c.GetString("perfer_trtm") // 우대사항
	pRecrutTitle := c.GetString("recrut_title")
	pSdy := c.GetString("sdy")
	pEdy := c.GetString("edy")
	pArrQstTitle := c.GetString("qst_title_arr")
	pDcmntEvlUseCd := c.GetString("dcmnt_evl_use_cd")
	pOnwyIntrvUseCd := c.GetString("onwy_intrv_use_cd")
	pLiveIntrvUseCd := c.GetString("live_intrv_use_cd")

	// LDK 2020/08/24 채용 정보 코드화 -->
	pUpJobGrpCd := c.GetString("up_job_grp_cd")
	pCarrGbnCd := c.GetString("carr_gbn_cd")
	pEntpAddr := c.GetString("entp_addr")
	pEmplTypCd := c.GetString("empl_typ_cd")
	pLstEduGbnCd := c.GetString("lst_edu_gbn_cd")
	pPrgsStatStep := c.GetString("prgs_stat_step") // 전형 절차
	pAnnualSalary := c.GetString("annual_salary")  // 임금(연봉)
	pWorkDays := c.GetString("work_days")          // 근무 요일
	pWelfare := c.GetString("welfare")             // 복리 후생
	pJobfair := c.GetString("jobfair")
	
	pPpChrgSn := session.Get(c.Ctx.Request.Context(), "mem_sn")

	log.Debug(fmt.Sprintf("pEntpMemNo:%s, pPpChrgSn:%s, pJobGrpCd:%s, pUpJobGrpCd:%s, pRecrutGbnCd:%s, pRecrutCnt:%s, pRol:%s, pAplyQufct:%s, pPerferTrtm:%s, pRecrutTitle:%s, pSdy:%s, pEdy:%s, pArrQstTitle:%s, pDcmntEvlUseCd:%s, pOnwyIntrvUseCd:%s, pLiveIntrvUseCd:%s, pCarrGbnCd:%s, pEntpAddr:%s, pEmplTypCd:%s, pLstEduGbnCd:%s, pPrgsStatStep:%s, pAnnualSalary:%s, pWorkDays:%s, pWelfare:%s, pJobfair:%s",
		pEntpMemNo, pPpChrgSn, pJobGrpCd, pUpJobGrpCd, pRecrutGbnCd, pRecrutCnt, pRol, pAplyQufct, pPerferTrtm, pRecrutTitle, pSdy, pEdy, pArrQstTitle, pDcmntEvlUseCd, pOnwyIntrvUseCd, pLiveIntrvUseCd, pCarrGbnCd, pEntpAddr, pEmplTypCd, pLstEduGbnCd, pPrgsStatStep, pAnnualSalary, pWorkDays, pWelfare, pJobfair))

	// 채용 공고 복사하기 위해
	// urlData := fmt.Sprintf("entp_mem_no=%s&mem_sn=%s&job_grp_cd=%s&recrut_gbn_cd=%s&recrut_cnt=%s&rol=%s&aply_qufct=%s&perfer_trtm=%s&recrut_title=%s&sdy=%s&edy=%s&qst_title_arr=%s&dcmnt_evl_use_cd=%s&onwy_intrv_use_cd=%s&live_intrv_use_cd=%s&carr_gbn_cd=%s&entp_addr=%s&empl_typ_cd=%s&lst_edu_gbn_cd=%s&prgs_stat_step=%s&annual_salary=%s&work_days=%s&welfare=%s&jobfair=%s",
	// 	pEntpMemNo,
	// 	pPpChrgSn,
	// 	pJobGrpCd,
	// 	pRecrutGbnCd,
	// 	pRecrutCnt,
	// 	pRol,        // 역할
	// 	pAplyQufct,  // 지원자격
	// 	pPerferTrtm, // 우대사항,
	// 	pRecrutTitle,
	// 	pSdy, pEdy,
	// 	pArrQstTitle,
	// 	pDcmntEvlUseCd,
	// 	pOnwyIntrvUseCd,
	// 	pLiveIntrvUseCd,
	// 	pCarrGbnCd,
	// 	pEntpAddr,
	// 	pEmplTypCd,
	// 	pLstEduGbnCd,
	// 	pPrgsStatStep, // 전형 절차
	// 	pAnnualSalary, // 임금(연봉),
	// 	pWorkDays,     // 근무 요일,
	// 	pWelfare,      // 복리 후생,
	// 	pJobfair)

	// log.Debug(fmt.Sprintf("curl -X POST \"%s/api/recruit/insert\" --data \"%s\"", beego.AppConfig.String("siteurl"), url.QueryEscape(urlData)))

	siteURL, _ := beego.AppConfig.String("siteurl")		

	//curl -X POST "http://localhost:7070/api/recruit/insert" --data "entp_mem_no=E2018102500001&mem_sn=0001&job_grp_cd=01030&recrut_gbn_cd=01&recrut_cnt=0&rol=%EC%97%AD%ED%95%A01%0A%EC%97%AD%ED%95%A02%0A%0A%EC%97%AD%ED%95%A03%0A&aply_qufct=%EC%A7%80%EC%9B%901%0A%0A%EC%A7%80%EC%9B%902&perfer_trtm=%0A%0A%EC%9A%B0%EB%8C%80%EC%82%AC%ED%95%AD&recrut_title=%EA%B3%B5%EB%B3%B5%EC%82%AC&sdy=20210303&edy=20210331&qst_title_arr=&dcmnt_evl_use_cd=1&onwy_intrv_use_cd=0&live_intrv_use_cd=1&carr_gbn_cd=02&entp_addr=%EC%84%9C%EC%9A%B8+%EC%A0%84%EC%A7%80%EC%97%AD&empl_typ_cd=01&lst_edu_gbn_cd=01&prgs_stat_step=1&annual_salary=1&work_days=1&welfare=1&jobfair="
	log.Debug(fmt.Sprintf("curl -X POST \"%s/api/recruit/insert\" --data \"entp_mem_no=%s&mem_sn=%s&job_grp_cd=%s&recrut_gbn_cd=%s&recrut_cnt=%s&rol=%s&aply_qufct=%s&perfer_trtm=%s&recrut_title=%s&sdy=%s&edy=%s&qst_title_arr=%s&dcmnt_evl_use_cd=%s&onwy_intrv_use_cd=%s&live_intrv_use_cd=%s&carr_gbn_cd=%s&entp_addr=%s&empl_typ_cd=%s&lst_edu_gbn_cd=%s&prgs_stat_step=%s&annual_salary=%s&work_days=%s&welfare=%s&jobfair=%s\"",
		siteURL,
		pEntpMemNo,
		pPpChrgSn,
		pJobGrpCd,
		pRecrutGbnCd,
		pRecrutCnt,
		// strings.Replace(pRol, string(10), "%0A", -1),        // 역할
		// strings.Replace(pAplyQufct, string(10), "%0A", -1),  // 지원자격
		// strings.Replace(pPerferTrtm, string(10), "%0A", -1), // 우대사항,
		url.QueryEscape(pRol),         // 역할
		url.QueryEscape(pAplyQufct),   // 지원자격
		url.QueryEscape(pPerferTrtm),  // 우대사항,
		url.QueryEscape(pRecrutTitle), // 채용공고 제목
		pSdy,
		pEdy,
		url.QueryEscape(pArrQstTitle), // 질문 리스트(배열)
		pDcmntEvlUseCd,
		pOnwyIntrvUseCd,
		pLiveIntrvUseCd,
		pCarrGbnCd,
		url.QueryEscape(pEntpAddr), // 소재지
		pEmplTypCd,
		pLstEduGbnCd,
		url.QueryEscape(pPrgsStatStep), // 전형 절차
		url.QueryEscape(pAnnualSalary), // 임금(연봉)
		url.QueryEscape(pWorkDays),     // 근무 요일
		url.QueryEscape(pWelfare),      // 복리 후생
		pJobfair))

	// strings.Replace(pRol, string(10), "%0A", -1)          // 역할
	// strings.Replace(pAplyQufct, string(10), "%0A", -1)    // 지원자격
	// strings.Replace(pPerferTrtm, string(10), "%0A", -1)   // 우대사항
	// strings.Replace(pPrgsStatStep, string(10), "%0A", -1) // 전형 절차
	// strings.Replace(pAnnualSalary, string(10), "%0A", -1) // 임금(연봉)
	// strings.Replace(pWorkDays, string(10), "%0A", -1)     // 근무 요일
	// strings.Replace(pWelfare, string(10), "%0A", -1)      // 복리 후생

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
	// LDK 2020/08/26 : 채용 정보 코드화, 추가 -->
	log.Debug(fmt.Sprintf("CALL ZSP_RECRUIT_REG_PROC_V2('%v', '%v', '%v', '%v',  %v , '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pJobGrpCd, pRecrutGbnCd, pRecrutCnt, pRol, pAplyQufct, pPerferTrtm, pRecrutTitle, pSdy, pEdy, pArrQstTitle, pDcmntEvlUseCd, pOnwyIntrvUseCd, pLiveIntrvUseCd, pPpChrgSn, pCarrGbnCd, pEntpAddr, pEmplTypCd, pLstEduGbnCd, pPrgsStatStep, pAnnualSalary, pWorkDays, pWelfare, pJobfair))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_RECRUIT_REG_PROC_V2('%v', '%v', '%v', '%v',  %v , '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pJobGrpCd, pRecrutGbnCd, pRecrutCnt, pRol, pAplyQufct, pPerferTrtm, pRecrutTitle, pSdy, pEdy, pArrQstTitle, pDcmntEvlUseCd, pOnwyIntrvUseCd, pLiveIntrvUseCd, pPpChrgSn, pCarrGbnCd, pEntpAddr, pEmplTypCd, pLstEduGbnCd, pPrgsStatStep, pAnnualSalary, pWorkDays, pWelfare, pJobfair),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
	)
	// <--

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

	rtnRecruitInsert := models.RtnRecruitInsert{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnRecruitInsert = models.RtnRecruitInsert{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
		// 채용공고등록 : 1001
		entp_memNo := pEntpMemNo.(string)
		gbn := "1001"
		val := pEntpMemNo.(string)

		go RecruitRegFCM(entp_memNo, gbn, val)
	}

	// End : Recruit Insert Process

	c.Data["json"] = &rtnRecruitInsert
	c.ServeJSON()
}

func RecruitRegFCM(entp_memNo string, gbn string, val string) {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

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

	log.Debug(fmt.Sprintf("CALL ZMSP_PUSH_KEY_INFO_S_R('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, entp_memNo, gbn, val, gbn1, entp_memNo, "", "", ""))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZMSP_PUSH_KEY_INFO_S_R('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, entp_memNo, gbn, val, gbn1, entp_memNo, "", "", ""),
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

			log.Debug("token : %v", token)

			opt := option.WithCredentialsFile("qrate-2ee14-firebase-adminsdk-64reu-74554f5c44.json")
			app, err := firebase.NewApp(context.Background(), nil, opt)
			if err != nil {
				log.Debug("error initializing app: %v\n", err)
			}

			// [START send_to_token_golang]
			// Obtain a messaging.Client from the App.
			ctx := context.Background()
			client, err := app.Messaging(ctx)

			// This registration token comes from the client FCM SDKs.
			registrationToken := token

			// See documentation on defining a message payload.
			message := &messaging.Message{
				Android: &messaging.AndroidConfig{
					Data: map[string]string{
						"type":  gbn,
						"title": "[직감]채용공고 등록 알림",
						"body":  cont,
					},
					/*
						Notification: &messaging.AndroidNotification{
							Title: "[직감] " + cont,
							Body:  cont,
						},
					*/
				},
				APNS: &messaging.APNSConfig{
					Headers: map[string]string{
						"type":  gbn,
						"title": "[직감]채용공고 등록 알림",
						"body":  cont,
					},
					Payload: &messaging.APNSPayload{
						Aps: &messaging.Aps{
							Alert: &messaging.ApsAlert{
								Title: "[직감]채용공고 등록 알림",
								Body:  cont,
							},
						},
					},
				},
				Token: registrationToken,
			}
			// Send a message to the device corresponding to the provided
			// registration token.
			response, err := client.Send(ctx, message)
			if err != nil {
				log.Debug("STATUS : ", err)
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
