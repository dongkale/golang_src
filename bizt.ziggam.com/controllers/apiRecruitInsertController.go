package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	"bizt.ziggam.com/utils"
	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/rana/ora.v4"
)

// ApiRecruitInsertController ...
type ApiRecruitInsertController struct {
	//beego.Controller
	ApiBaseController
}

// Prepare ...
func (c *ApiRecruitInsertController) Prepare() {
	c.ApiBaseController.Prepare()

	// Prepare...
}

// Post ...
// curl -X POST "http://localhost:7070/api/recruit/insert" --data "entp_mem_no=E2018102500001&mem_sn=0001&job_grp_cd=01017&recrut_gbn_cd=01&recrut_cnt=0&rol=역할&aply_qufct=지원자격&perfer_trtm=우대사항&recrut_title=채용 제목&sdy=20210218&edy=20210228&qst_title_arr=질문11,질문22&dcmnt_evl_use_cd=1&onwy_intrv_use_cd=1&live_intrv_use_cd=1&carr_gbn_cd=01&entp_addr=서울 전지역&empl_typ_cd=02&lst_edu_gbn_cd=03&prgs_stat_step=전형절차&annual_salary=임금&work_days=근무요일&welfare=복리후생&jobfair="
// curl -X POST "localhost:7070/api/recruit/insert" --data "entp_mem_no=E2018102500001&mem_sn=0001&recrut_title=채용제목&job_grp_cd=01017&recrut_gbn_cd=01&recrut_cnt=0&rol=역할&aply_qufct=지원자격&perfer_trtm=우대사항&sdy=20210216&edy=20210220&qst_title_arr=질문1,질문2&dcmnt_evl_use_cd=1&onwy_intrv_use_cd=1&live_intrv_use_cd=1&carr_gbn_cd=02&entp_addr=서울 전지역&empl_typ_cd=03&lst_edu_gbn_cd=04&prgs_stat_step=전형절차&annual_salary=임금&work_days=근무요일&welfare=복리후생&jobfair="
// ucurl -X POST "localhost:7070/api/recruit/insert" --data "entp_mem_no=E2018102500001&mem_sn=0001&recrut_title=채용제목&job_grp_cd=01017&recrut_gbn_cd=01&recrut_cnt=0&rol=역할&aply_qufct=지원자격&perfer_trtm=우대사항&sdy=20210216&edy=20210220&qst_title_arr=질문1,질문2&dcmnt_evl_use_cd=1&onwy_intrv_use_cd=1&live_intrv_use_cd=1&carr_gbn_cd=02&entp_addr=서울 전지역&empl_typ_cd=03&lst_edu_gbn_cd=04&prgs_stat_step=전형절차&annual_salary=임금&work_days=근무요일&welfare=복리후생&jobfair="
//// curl --data-urlencode "entp_mem_no=E2018102500001&mem_sn=0001&recrut_title=채용제목&job_grp_cd=01017&recrut_gbn_cd=01&recrut_cnt=0&rol=역할&aply_qufct=지원자격&perfer_trtm=우대사항&sdy=20210216&edy=20210220&qst_title_arr=질문1,질문2&dcmnt_evl_use_cd=1&onwy_intrv_use_cd=1&live_intrv_use_cd=1&carr_gbn_cd=02&entp_addr=서울 전지역&empl_typ_cd=03&lst_edu_gbn_cd=04&prgs_stat_step=전형절차&annual_salary=임금&work_days=근무요일&welfare=복리후생&jobfair=&remote=TRUE" localhost:7070/api/recruit/insert

// --> 최종(개행 처리:https://stackoverflow.com/questions/3872427/how-to-send-line-break-with-curl)
// curl -X POST "http://localhost:7070/api/recruit/insert" --data "entp_mem_no=E2018102500001&mem_sn=0001&job_grp_cd=01017&recrut_gbn_cd=01&recrut_cnt=0&rol=역할1%0A역할2%0A역할3%0A역할4%0A역할5%0A%0A역할6%0A%0A 역할7&aply_qufct=지원자격1%0A%0A지원자격2%0A%0A지원자격3&perfer_trtm=우대사항1%0A%0A우대사항2%0A%0A우대사항3%0A%0A우대사항4%0A%0A우대사항5%0A&recrut_title=공고복사&sdy=20210303&edy=20210331&qst_title_arr=질문1,질문2&dcmnt_evl_use_cd=1&onwy_intrv_use_cd=1&live_intrv_use_cd=1&carr_gbn_cd=01&entp_addr=서울 전지역&empl_typ_cd=01&lst_edu_gbn_cd=01&prgs_stat_step=전형절차1%0A전형절차2%0A&annual_salary=임금1%0A임금2&work_days=근무요일1%0A근무요일2&welfare=복리후생1%0A%0A복리후생2%0A%0A복리후생3&jobfair="
// url.QueryEscape() 사용 :  curl -X POST "http://localhost:7070/api/recruit/insert" --data "entp_mem_no=E2018102500001&mem_sn=0001&job_grp_cd=01017&recrut_gbn_cd=01&recrut_cnt=0&rol=%EC%97%AD%ED%95%A01%0A%EC%97%AD%ED%95%A02%0A%0A%EC%97%AD%ED%95%A03%0A%0A%EC%97%AD%ED%95%A04%0A%0A%EC%97%AD%ED%95%A05%0A%0A&aply_qufct=%0A%EC%A7%80%EC%9B%90%EC%9E%90%EA%B2%A91%0A%EC%A7%80%EC%9B%90%EC%9E%90%EA%B2%A92%0A%EC%A7%80%EC%9B%90%EC%9E%90%EA%B2%A93%0A%0A%EC%A7%80%EC%9B%90%EC%9E%90%EA%B2%A94%0A%EC%A7%80%EC%9B%90%EC%9E%90%EA%B2%A95%0A&perfer_trtm=-+%EC%9A%B0%EB%8C%80%EC%82%AC%ED%95%AD1%0A-+%EC%9A%B0%EB%8C%80%EC%82%AC%ED%95%AD2%0A%2A%EC%9A%B0%EB%8C%80%EC%82%AC%ED%95%AD%2A%0A%21%EC%9A%B0%EB%8C%80%EC%82%AC%ED%95%AD%21%0A%40%EC%9A%B0%EB%8C%80%EC%82%AC%ED%95%AD%40%0A~%21%40%23%24%25%5E%26%2A%28%29&recrut_title=%EA%B3%B5%EA%B3%A0%EB%B3%B5%EC%82%AC&sdy=20210304&edy=20210331&qst_title_arr=%EC%A7%88%EB%AC%B81%3F%2C%EC%A7%88%EB%AC%B82%3F&dcmnt_evl_use_cd=1&onwy_intrv_use_cd=1&live_intrv_use_cd=1&carr_gbn_cd=01&entp_addr=%EC%84%9C%EC%9A%B8+%EC%A0%84%EC%A7%80%EC%97%AD&empl_typ_cd=01&lst_edu_gbn_cd=01&prgs_stat_step=%EC%A0%84%ED%98%95%EC%A0%88%EC%B0%A8&annual_salary=%EC%9E%84%EA%B8%88&work_days=%EA%B7%BC%EB%AC%B4%EC%9A%94%EC%9D%BC&welfare=%EB%B3%B5%EB%A6%AC%ED%9B%84%EC%83%9D%0A%EB%B3%B5%EB%A6%AC%ED%9B%84%EC%83%9D&jobfair="
func (c *ApiRecruitInsertController) Post() {

	//session := c.StartSession()
	pLang, _ := beego.AppConfig.String("lang")

	pRecrutTitle2 := c.Controller.GetString("recrut_title")
	fmt.Printf(pRecrutTitle2)

	pEntpMemNo := c.GetString("entp_mem_no")
	pPpChrgSn := c.GetString("mem_sn")
	pJobGrpCd := c.GetString("job_grp_cd")
	pRecrutGbnCd := c.GetString("recrut_gbn_cd")
	pRecrutCnt := c.GetString("recrut_cnt")
	pRol := c.GetString("rol")
	pAplyQufct := c.GetString("aply_qufct")
	pPerferTrtm := c.GetString("perfer_trtm")
	pRecrutTitle := c.GetString("recrut_title")
	pSdy := c.GetString("sdy")
	pEdy := c.GetString("edy")
	pArrQstTitle := c.GetString("qst_title_arr")
	pDcmntEvlUseCd := c.GetString("dcmnt_evl_use_cd")
	pOnwyIntrvUseCd := c.GetString("onwy_intrv_use_cd")
	pLiveIntrvUseCd := c.GetString("live_intrv_use_cd")

	//pUpJobGrpCd := c.GetString("up_job_grp_cd")
	pCarrGbnCd := c.GetString("carr_gbn_cd")
	pEntpAddr := c.GetString("entp_addr")
	pEmplTypCd := c.GetString("empl_typ_cd")
	pLstEduGbnCd := c.GetString("lst_edu_gbn_cd")
	pPrgsStatStep := c.GetString("prgs_stat_step")
	pAnnualSalary := c.GetString("annual_salary")
	pWorkDays := c.GetString("work_days")
	pWelfare := c.GetString("welfare")
	pJobfair := c.GetString("jobfair")
	// pRemote := c.GetString("remote")

	// if pRemote != "" {
	// 	// 한글 변환
	// 	pRol = utils.ConvertEucKR(pRol)
	// 	pAplyQufct = utils.ConvertEucKR(pAplyQufct)
	// 	pPerferTrtm = utils.ConvertEucKR(pPerferTrtm)
	// 	pRecrutTitle = utils.ConvertEucKR(pRecrutTitle)
	// 	pArrQstTitle = utils.ConvertEucKR(pArrQstTitle)
	// 	pEntpAddr = utils.ConvertEucKR(pEntpAddr)
	// 	pPrgsStatStep = utils.ConvertEucKR(pPrgsStatStep)
	// 	pAnnualSalary = utils.ConvertEucKR(pAnnualSalary)
	// 	pWorkDays = utils.ConvertEucKR(pWorkDays)
	// 	pWelfare = utils.ConvertEucKR(pWelfare)
	// }

	if utils.IsDateDtFmt(pSdy) == false || utils.IsDateDtFmt(pEdy) == false {
		fmt.Printf(fmt.Sprintf("[Error] pEntpMemNo:%s, pPpChrgSn:%s, pJobGrpCd:%s, pRecrutGbnCd:%s, pRecrutCnt:%s, pRol:%s, pAplyQufct:%s, pPerferTrtm:%s, pRecrutTitle:%s, pSdy:%s, pEdy:%s, pArrQstTitle:%s, pDcmntEvlUseCd:%s, pOnwyIntrvUseCd:%s, pLiveIntrvUseCd:%s, pCarrGbnCd:%s, pEntpAddr:%s, pEmplTypCd:%s, pLstEduGbnCd:%s, pPrgsStatStep:%s, pAnnualSalary:%s, pWorkDays:%s, pWelfare:%s, pJobfair:%s --> Invalid Sdy, Edy",
			pEntpMemNo, pPpChrgSn, pJobGrpCd, pRecrutGbnCd, pRecrutCnt, pRol, pAplyQufct, pPerferTrtm, pRecrutTitle, pSdy, pEdy, pArrQstTitle, pDcmntEvlUseCd, pOnwyIntrvUseCd, pLiveIntrvUseCd, pCarrGbnCd, pEntpAddr, pEmplTypCd, pLstEduGbnCd, pPrgsStatStep, pAnnualSalary, pWorkDays, pWelfare, pJobfair))

		return
	}

	fmt.Printf(fmt.Sprintf("pEntpMemNo:%s, pPpChrgSn:%s, pJobGrpCd:%s, pRecrutGbnCd:%s, pRecrutCnt:%s, pRol:%s, pAplyQufct:%s, pPerferTrtm:%s, pRecrutTitle:%s, pSdy:%s, pEdy:%s, pArrQstTitle:%s, pDcmntEvlUseCd:%s, pOnwyIntrvUseCd:%s, pLiveIntrvUseCd:%s, pCarrGbnCd:%s, pEntpAddr:%s, pEmplTypCd:%s, pLstEduGbnCd:%s, pPrgsStatStep:%s, pAnnualSalary:%s, pWorkDays:%s, pWelfare:%s, pJobfair:%s",
		pEntpMemNo, pPpChrgSn, pJobGrpCd, pRecrutGbnCd, pRecrutCnt, pRol, pAplyQufct, pPerferTrtm, pRecrutTitle, pSdy, pEdy, pArrQstTitle, pDcmntEvlUseCd, pOnwyIntrvUseCd, pLiveIntrvUseCd, pCarrGbnCd, pEntpAddr, pEmplTypCd, pLstEduGbnCd, pPrgsStatStep, pAnnualSalary, pWorkDays, pWelfare, pJobfair))

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
	fmt.Printf(fmt.Sprintf("CALL ZSP_RECRUIT_REG_PROC_V2('%v', '%v', '%v', '%v',  %v , '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pJobGrpCd, pRecrutGbnCd, pRecrutCnt, pRol, pAplyQufct, pPerferTrtm, pRecrutTitle, pSdy, pEdy, pArrQstTitle, pDcmntEvlUseCd, pOnwyIntrvUseCd, pLiveIntrvUseCd, pPpChrgSn, pCarrGbnCd, pEntpAddr, pEmplTypCd, pLstEduGbnCd, pPrgsStatStep, pAnnualSalary, pWorkDays, pWelfare, pJobfair))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_RECRUIT_REG_PROC_V2('%v', '%v', '%v', '%v',  %v , '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pJobGrpCd, pRecrutGbnCd, pRecrutCnt, pRol, pAplyQufct, pPerferTrtm, pRecrutTitle, pSdy, pEdy, pArrQstTitle, pDcmntEvlUseCd, pOnwyIntrvUseCd, pLiveIntrvUseCd, pPpChrgSn, pCarrGbnCd, pEntpAddr, pEmplTypCd, pLstEduGbnCd, pPrgsStatStep, pAnnualSalary, pWorkDays, pWelfare, pJobfair),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
		ora.S,   /* RTN_RECRUT_SN */
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
		rtnCd       int64
		rtnMsg      string
		rtnRecrutSn string
	)

	rtnRecruitInsert := models.RtnRecruitInsert{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
			rtnRecrutSn = procRset.Row[2].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnRecruitInsert = models.RtnRecruitInsert{
			RtnCd:       rtnCd,
			RtnMsg:      rtnMsg,
			RtnRecrutSn: rtnRecrutSn,
		}

		fmt.Printf(fmt.Sprintf("===> rtnCd:%v, rtnMsg:%v, rtnRecrutSn:%v", rtnCd, rtnMsg, rtnRecrutSn))

		// 채용공고등록 : 1001
		entp_memNo := pEntpMemNo
		gbn := "1001"
		val := pEntpMemNo

		go RecruitRegFCM(entp_memNo, gbn, val)
	}
	// End : Recruit Insert Process

	c.Data["json"] = &rtnRecruitInsert
	c.ServeJSON()
}

/*
func RecruitRegFCM(entp_memNo string, gbn string, val string) {

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

	fmt.Printf(fmt.Sprintf("CALL ZMSP_PUSH_KEY_INFO_S_R('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, entp_memNo, gbn, val, gbn1, entp_memNo, "", "", ""))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZMSP_PUSH_KEY_INFO_S_R('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, entp_memNo, gbn, val, gbn1, entp_memNo, "", "", ""),
		ora.S, // PUSH_KEY
		ora.S, // CONT
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

			fmt.Printf("token : %v", token)

			opt := option.WithCredentialsFile("qrate-2ee14-firebase-adminsdk-64reu-74554f5c44.json")
			app, err := firebase.NewApp(context.Background(), nil, opt)
			if err != nil {
				fmt.Printf("error initializing app: %v\n", err)
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
					//
					// 	Notification: &messaging.AndroidNotification{
					// 		Title: "[직감] " + cont,
					// 		Body:  cont,
					// 	},
					//
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
*/
