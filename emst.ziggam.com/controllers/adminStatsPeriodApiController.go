package controllers

import (
	"fmt"
	"strings"

	"emst.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	ora "gopkg.in/rana/ora.v4"
)

type AdminStatsPeriodApiController struct {
	BaseController
}

// Post ...
func (c *AdminStatsPeriodApiController) Post() {

	// start : log
	// log := logs.NewLogger()
	// log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/common/login")
	}

	pLang, _ := beego.AppConfig.String("lang")
	pType := c.GetString("search_type")
	pSdy := c.GetString("sdy")
	pEdy := c.GetString("edy")

	pEntpKoNm := c.GetString("entp_ko_nm")
	pJobFairMngCd := c.GetString("jf_mng_cd")

	logs.Debug(fmt.Sprintf("Type: %v, Sdy: %v, Edy: %v, EntpKoNm: %v, JobFairMngCd: %v", pType, pSdy, pEdy, pEntpKoNm, pJobFairMngCd))

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// type 분기
	// 01 : 공고 내용 + 기간 (시작일, 종료일)
	// 02 : 지원 내역(전체)
	// 03 : 원웨이 내역
	// 04 : 라이브 내역(전체)
	// 05 : 라이브 내역(예정/종료)

	if pType == "01" {
		logs.Debug(fmt.Sprintf("CALL SP_EMS_API_PERIOD_TYPE_R('%v', '%v', '%v', '%v', '%v', '%v', :1)",
			pLang, pType, pSdy, pEdy, pEntpKoNm, pJobFairMngCd))

		stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_API_PERIOD_TYPE_R('%v', '%v', '%v', '%v', '%v', '%v', :1)",
			pLang, pType, pSdy, pEdy, pEntpKoNm, pJobFairMngCd),
			ora.S, /* ENTP_KO_NM */
			ora.S, /* ENTP_MEM_NO */
			ora.S, /* JOBFAIR_MNG_CDS */
			ora.S, /* RECRUT_TITLE */
			ora.S, /* RECRUT_SN */
			ora.S, /* JF_MNG_CD */
			ora.S, /* UP_JOB_GRP */
			ora.S, /* JOB_GRP */
			ora.S, /* DCMNT_EVL_USE_CD */
			ora.S, /* ONWY_INTRV_USE_CD */
			ora.S, /* LIVE_INTRV_USE_CD */
			ora.S, /* SDY */
			ora.S, /* EDY */
			ora.S, /* RECRUT_EDT */
			ora.S, /* REG_DT */
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

		rtnAdminStatsPeriod01 := models.RtnAdminStatsPeriod01{}
		adminStatsPeriod01 := make([]models.AdminStatsPeriod01, 0)

		var (
			entpKoNm    string
			entpMemNo   string
			jobfairCds  string
			recrutTitle string
			recrutSn    string
			recrutJfMngCd string
			upJobGrp    string
			jobGrp      string
			dcmntUseCd     string
			onwyUseCd      string
			liveIntrvUseCd string
			sdy         string
			edy         string
			recrutEdt   string
			regDt       string
		)

		if procRset.IsOpen() {
			for procRset.Next() {
				entpKoNm = procRset.Row[0].(string)
				entpMemNo = procRset.Row[1].(string)
				jobfairCds = procRset.Row[2].(string)
				recrutTitle = procRset.Row[3].(string)
				recrutSn = procRset.Row[4].(string)
				recrutJfMngCd = procRset.Row[5].(string)
				upJobGrp = procRset.Row[6].(string)
				jobGrp = procRset.Row[7].(string)
				dcmntUseCd = procRset.Row[8].(string)
				onwyUseCd = procRset.Row[9].(string)
				liveIntrvUseCd = procRset.Row[10].(string)
				sdy = procRset.Row[11].(string)
				edy = procRset.Row[12].(string)
				recrutEdt = procRset.Row[13].(string)
				regDt = procRset.Row[14].(string)

				adminStatsPeriod01 = append(adminStatsPeriod01, models.AdminStatsPeriod01{
					EntpKoNm:    entpKoNm,
					EntpMemNo:   entpMemNo,
					JobFairCds:     strings.ReplaceAll(jobfairCds, ",", "</br>"),
					RecrutTitle: recrutTitle,
					RecrutSn:    recrutSn,
					RecrutJfMngCd: recrutJfMngCd,
					UpJobGrp:    upJobGrp,
					JobGrp:      jobGrp,
					DcmntUseCd:     dcmntUseCd,
					OnwyUseCd:      onwyUseCd,
					LiveIntrvUseCd: liveIntrvUseCd,
					Sdy:         sdy,
					Edy:         edy,
					RecrutEdt:   recrutEdt,
					RegDt:       regDt,
				})
			}
			if err := procRset.Err(); err != nil {
				panic(err)
			}
			rtnAdminStatsPeriod01 = models.RtnAdminStatsPeriod01{
				Result: adminStatsPeriod01,
			}
			// End : Admin Stats Sub2

			c.Data["json"] = &rtnAdminStatsPeriod01
			c.ServeJSON()
		}
	} else if pType == "02" {
		logs.Debug(fmt.Sprintf("CALL SP_EMS_API_PERIOD_TYPE_R('%v', '%v', '%v', '%v', '%v', '%v', :1)",
			pLang, pType, pSdy, pEdy, pEntpKoNm, pJobFairMngCd))

		stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_API_PERIOD_TYPE_R('%v', '%v', '%v', '%v', '%v', '%v', :1)",
			pLang, pType, pSdy, pEdy, pEntpKoNm, pJobFairMngCd),
			ora.S, /* ENTP_KO_NM */
			ora.S, /* ENTP_MEM_NO */
			ora.S, /* RECRUT_TITLE */
			ora.S, /* RECRUT_SN */
			ora.S, /* JF_MNG_CD */
			ora.S, /* UP_JOB_GRP */
			ora.S, /* JOB_GRP */
			ora.S, /* REG_DT */
			ora.S, /* NM */
			ora.S, /* PP_MEM_NO */
			ora.S, /* SEX */
			ora.S, /* BRTH_YMD */
			ora.S, /* MO_NO */
			ora.S, /* EMAIL */
			ora.S, /* LST_EDU */
			ora.S, /* LST_EDU_DESC */
			ora.S, /* CARR_GBN */
			ora.S, /* CARR_GBN_DESC */
			ora.S, /* TECH_QLFT_KND */
			ora.S, /* FRGN_LANG_ABLT_DESC */
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

		rtnAdminStatsPeriod02 := models.RtnAdminStatsPeriod02{}
		adminStatsPeriod02 := make([]models.AdminStatsPeriod02, 0)

		var (
			entpKoNm         string
			entpMemNo        string
			jobfairCds       string
			recrutTitle      string
			recrutSn         string
			recrutJfMngCd    string
			upJobGrp         string
			jobGrp           string
			regDt            string
			nm               string
			ppMemNo          string
			sex              string
			birth            string
			moNo             string
			email            string
			lstEdu           string
			lstEduDesc       string
			carrGbn          string
			carrGbnDesc      string
			techQlftKnd      string
			frgnLangAbltDesc string
		)

		if procRset.IsOpen() {
			for procRset.Next() {
				entpKoNm = procRset.Row[0].(string)
				entpMemNo = procRset.Row[1].(string)
				jobfairCds = procRset.Row[2].(string)
				recrutTitle = procRset.Row[3].(string)
				recrutSn = procRset.Row[4].(string)
				recrutJfMngCd = procRset.Row[5].(string)
				upJobGrp = procRset.Row[6].(string)
				jobGrp = procRset.Row[7].(string)
				regDt = procRset.Row[8].(string)
				nm = procRset.Row[9].(string)
				ppMemNo = procRset.Row[10].(string)
				sex = procRset.Row[11].(string)
				birth = procRset.Row[12].(string)
				moNo = procRset.Row[13].(string)
				email = procRset.Row[14].(string)
				lstEdu = procRset.Row[15].(string)
				lstEduDesc = procRset.Row[16].(string)
				carrGbn = procRset.Row[17].(string)
				carrGbnDesc = procRset.Row[18].(string)
				techQlftKnd = procRset.Row[19].(string)
				frgnLangAbltDesc = procRset.Row[20].(string)

				adminStatsPeriod02 = append(adminStatsPeriod02, models.AdminStatsPeriod02{
					EntpKoNm:         entpKoNm,
					EntpMemNo:        entpMemNo,
					JobFairCds:       strings.ReplaceAll(jobfairCds, ",", "</br>"),
					RecrutTitle:      recrutTitle,
					RecrutSn:         recrutSn,
					RecrutJfMngCd:    recrutJfMngCd,
					UpJobGrp:         upJobGrp,
					JobGrp:           jobGrp,
					RegDt:            regDt,
					Nm:               nm,
					PpMemNo:          ppMemNo,
					Sex:              sex,
					Birth:            birth,
					MoNo:             moNo,
					Email:            email,
					LstEdu:           lstEdu,
					LstEduDesc:       lstEduDesc,
					CarrGbn:          carrGbn,
					CarrGbnDesc:      carrGbnDesc,
					TechQlftKnd:      techQlftKnd,
					FrgnLangAbltDesc: frgnLangAbltDesc,
				})
			}
			if err := procRset.Err(); err != nil {
				panic(err)
			}
			rtnAdminStatsPeriod02 = models.RtnAdminStatsPeriod02{
				Result: adminStatsPeriod02,
			}
			// End : Admin Stats Sub2

			c.Data["json"] = &rtnAdminStatsPeriod02
			c.ServeJSON()
		}
	} else if pType == "03" {
		logs.Debug(fmt.Sprintf("CALL SP_EMS_API_PERIOD_TYPE_R('%v', '%v', '%v', '%v', '%v', '%v', :1)",
			pLang, pType, pSdy, pEdy, pEntpKoNm, pJobFairMngCd))

		stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_API_PERIOD_TYPE_R('%v', '%v', '%v', '%v', '%v', '%v', :1)",
			pLang, pType, pSdy, pEdy, pEntpKoNm, pJobFairMngCd),
			ora.S, /* ENTP_KO_NM */
			ora.S, /* ENTP_MEM_NO */
			ora.S, /* RECRUT_TITLE */
			ora.S, /* RECRUT_SN */
			ora.S, /* JF_MNG_CD */
			ora.S, /* UP_JOB_GRP */
			ora.S, /* JOB_GRP */
			ora.S, /* REG_DT */
			ora.S, /* ONEWAY_DT */
			ora.S, /* ONEWAY_CNT */
			ora.S, /* ONEWAY_TYPE */
			ora.S, /* NM */
			ora.S, /* SEX */
			ora.S, /* BRTH_YMD */
			ora.S, /* MO_NO */
			ora.S, /* EMAIL */
			ora.S, /* LST_EDU */
			ora.S, /* LST_EDU_DESC */
			ora.S, /* CARR_GBN */
			ora.S, /* CARR_GBN_DESC */
			ora.S, /* TECH_QLFT_KND */
			ora.S, /* FRGN_LANG_ABLT_DESC */
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

		rtnAdminStatsPeriod03 := models.RtnAdminStatsPeriod03{}
		adminStatsPeriod03 := make([]models.AdminStatsPeriod03, 0)

		var (
			entpKoNm         string
			entpMemNo        string
			jobfairCds       string
			recrutTitle      string
			recrutSn         string
			recrutJfMngCd    string
			upJobGrp         string
			jobGrp           string
			regDt            string
			oneWayDt         string
			oneWayCnt        string
			oneWayType       string
			nm               string
			sex              string
			birth            string
			moNo             string
			email            string
			lstEdu           string
			lstEduDesc       string
			carrGbn          string
			carrGbnDesc      string
			techQlftKnd      string
			frgnLangAbltDesc string
		)

		if procRset.IsOpen() {
			for procRset.Next() {
				entpKoNm = procRset.Row[0].(string)
				entpMemNo = procRset.Row[1].(string)
				jobfairCds = procRset.Row[2].(string)
				recrutTitle = procRset.Row[3].(string)
				recrutSn = procRset.Row[4].(string)
				recrutJfMngCd = procRset.Row[5].(string)
				upJobGrp = procRset.Row[6].(string)
				jobGrp = procRset.Row[7].(string)
				regDt = procRset.Row[8].(string)
				oneWayDt = procRset.Row[9].(string)
				oneWayCnt = procRset.Row[10].(string)
				oneWayType = procRset.Row[11].(string)
				nm = procRset.Row[12].(string)
				sex = procRset.Row[13].(string)
				birth = procRset.Row[14].(string)
				moNo = procRset.Row[15].(string)
				email = procRset.Row[16].(string)
				lstEdu = procRset.Row[17].(string)
				lstEduDesc = procRset.Row[18].(string)
				carrGbn = procRset.Row[19].(string)
				carrGbnDesc = procRset.Row[20].(string)
				techQlftKnd = procRset.Row[21].(string)
				frgnLangAbltDesc = procRset.Row[22].(string)

				adminStatsPeriod03 = append(adminStatsPeriod03, models.AdminStatsPeriod03{
					EntpKoNm:         entpKoNm,
					EntpMemNo:        entpMemNo,
					JobFairCds:       strings.ReplaceAll(jobfairCds, ",", "</br>"),
					RecrutTitle:      recrutTitle,
					RecrutSn:         recrutSn,
					RecrutJfMngCd:    recrutJfMngCd,
					UpJobGrp:         upJobGrp,
					JobGrp:           jobGrp,
					RegDt:            regDt,
					OneWayDt:         oneWayDt,
					OneWayCnt:        oneWayCnt,
					OneWayType:       oneWayType,
					Nm:               nm,
					Sex:              sex,
					Birth:            birth,
					MoNo:             moNo,
					Email:            email,
					LstEdu:           lstEdu,
					LstEduDesc:       lstEduDesc,
					CarrGbn:          carrGbn,
					CarrGbnDesc:      carrGbnDesc,
					TechQlftKnd:      techQlftKnd,
					FrgnLangAbltDesc: frgnLangAbltDesc,
				})
			}
			if err := procRset.Err(); err != nil {
				panic(err)
			}
			rtnAdminStatsPeriod03 = models.RtnAdminStatsPeriod03{
				Result: adminStatsPeriod03,
			}
			// End : Admin Stats Sub2

			c.Data["json"] = &rtnAdminStatsPeriod03
			c.ServeJSON()
		}
	} else if pType == "04" {
		logs.Debug(fmt.Sprintf("CALL SP_EMS_API_PERIOD_TYPE_R('%v', '%v', '%v', '%v', '%v', '%v', :1)",
			pLang, pType, pSdy, pEdy, pEntpKoNm, pJobFairMngCd))

		stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_API_PERIOD_TYPE_R('%v', '%v', '%v', '%v', '%v', '%v', :1)",
			pLang, pType, pSdy, pEdy, pEntpKoNm, pJobFairMngCd),
			ora.S, /* ENTP_KO_NM */
			ora.S, /* ENTP_MEM_NO */
			ora.S, /* RECRUT_TITLE */
			ora.S, /* RECRUT_SN */
			ora.S, /* JF_MNG_CD */
			ora.S, /* UP_JOB_GRP */
			ora.S, /* JOB_GRP */
			ora.S, /* REQUEST_DT */
			ora.S, /* BEGIN_DT */
			ora.S, /* END_DT */
			ora.S, /* LIVE_STATE */
			ora.S, /* NM */
			ora.S, /* SEX */
			ora.S, /* BRTH_YMD */
			ora.S, /* MO_NO */
			ora.S, /* EMAIL */
			ora.S, /* LST_EDU */
			ora.S, /* LST_EDU_DESC */
			ora.S, /* CARR_GBN */
			ora.S, /* CARR_GBN_DESC */
			ora.S, /* TECH_QLFT_KND */
			ora.S, /* FRGN_LANG_ABLT_DESC */
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

		rtnAdminStatsPeriod04 := models.RtnAdminStatsPeriod04{}
		adminStatsPeriod04 := make([]models.AdminStatsPeriod04, 0)

		var (
			entpKoNm         string
			entpMemNo        string
			jobfairCds       string
			recrutTitle      string
			recrutSn         string
			recrutJfMngCd    string
			upJobGrp         string
			jobGrp           string
			requestDt        string
			beginDt          string
			endDt            string
			liveState        string
			nm               string
			sex              string
			birth            string
			moNo             string
			email            string
			lstEdu           string
			lstEduDesc       string
			carrGbn          string
			carrGbnDesc      string
			techQlftKnd      string
			frgnLangAbltDesc string
		)

		if procRset.IsOpen() {
			for procRset.Next() {
				entpKoNm = procRset.Row[0].(string)
				entpMemNo = procRset.Row[1].(string)
				jobfairCds = procRset.Row[2].(string)
				recrutTitle = procRset.Row[3].(string)
				recrutSn = procRset.Row[4].(string)
				recrutJfMngCd = procRset.Row[5].(string)
				upJobGrp = procRset.Row[6].(string)
				jobGrp = procRset.Row[7].(string)
				requestDt = procRset.Row[8].(string)
				beginDt = procRset.Row[9].(string)
				endDt = procRset.Row[10].(string)
				liveState = procRset.Row[11].(string)
				nm = procRset.Row[12].(string)
				sex = procRset.Row[13].(string)
				birth = procRset.Row[14].(string)
				moNo = procRset.Row[15].(string)
				email = procRset.Row[16].(string)
				lstEdu = procRset.Row[17].(string)
				lstEduDesc = procRset.Row[18].(string)
				carrGbn = procRset.Row[19].(string)
				carrGbnDesc = procRset.Row[20].(string)
				techQlftKnd = procRset.Row[21].(string)
				frgnLangAbltDesc = procRset.Row[22].(string)

				adminStatsPeriod04 = append(adminStatsPeriod04, models.AdminStatsPeriod04{
					EntpKoNm:         entpKoNm,
					EntpMemNo:        entpMemNo,
					JobFairCds:       strings.ReplaceAll(jobfairCds, ",", "</br>"),
					RecrutTitle:      recrutTitle,
					RecrutSn:         recrutSn,
					RecrutJfMngCd:    recrutJfMngCd,
					UpJobGrp:         upJobGrp,
					JobGrp:           jobGrp,
					RequestDt:        requestDt,
					BeginDt:          beginDt,
					EndDt:            endDt,
					LiveState:        liveState,
					Nm:               nm,
					Sex:              sex,
					Birth:            birth,
					MoNo:             moNo,
					Email:            email,
					LstEdu:           lstEdu,
					LstEduDesc:       lstEduDesc,
					CarrGbn:          carrGbn,
					CarrGbnDesc:      carrGbnDesc,
					TechQlftKnd:      techQlftKnd,
					FrgnLangAbltDesc: frgnLangAbltDesc,
				})
			}
			if err := procRset.Err(); err != nil {
				panic(err)
			}
			rtnAdminStatsPeriod04 = models.RtnAdminStatsPeriod04{
				Result: adminStatsPeriod04,
			}
			// End : Admin Stats Sub2

			c.Data["json"] = &rtnAdminStatsPeriod04
			c.ServeJSON()
		}
	} else if pType == "05" {
		logs.Debug(fmt.Sprintf("CALL SP_EMS_API_PERIOD_TYPE_R('%v', '%v', '%v', '%v', '%v', '%v', :1)",
			pLang, pType, pSdy, pEdy, pEntpKoNm, pJobFairMngCd))

		stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_API_PERIOD_TYPE_R('%v', '%v', '%v', '%v', '%v', '%v', :1)",
			pLang, pType, pSdy, pEdy, pEntpKoNm, pJobFairMngCd),
			ora.S, /* ENTP_KO_NM */
			ora.S, /* ENTP_MEM_NO */
			ora.S, /* RECRUT_TITLE */
			ora.S, /* RECRUT_SN */
			ora.S, /* JF_MNG_CD */
			ora.S, /* UP_JOB_GRP */
			ora.S, /* JOB_GRP */
			ora.S, /* CONFIRM_DT */
			ora.S, /* BEGIN_DT */
			ora.S, /* END_DT */
			ora.S, /* NM */
			ora.S, /* SEX */
			ora.S, /* BRTH_YMD */
			ora.S, /* MO_NO */
			ora.S, /* EMAIL */
			ora.S, /* LST_EDU */
			ora.S, /* LST_EDU_DESC */
			ora.S, /* CARR_GBN */
			ora.S, /* CARR_GBN_DESC */
			ora.S, /* TECH_QLFT_KND */
			ora.S, /* FRGN_LANG_ABLT_DESC */
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

		rtnAdminStatsPeriod05 := models.RtnAdminStatsPeriod05{}
		adminStatsPeriod05 := make([]models.AdminStatsPeriod05, 0)

		var (
			entpKoNm         string
			entpMemNo        string
			jobfairCds       string
			recrutTitle      string
			recrutSn         string
			recrutJfMngCd    string
			upJobGrp         string
			jobGrp           string
			confirmDt        string
			beginDt          string
			endDt            string
			nm               string
			sex              string
			birth            string
			moNo             string
			email            string
			lstEdu           string
			lstEduDesc       string
			carrGbn          string
			carrGbnDesc      string
			techQlftKnd      string
			frgnLangAbltDesc string
		)

		if procRset.IsOpen() {
			for procRset.Next() {
				entpKoNm = procRset.Row[0].(string)
				entpMemNo = procRset.Row[1].(string)
				jobfairCds = procRset.Row[2].(string)
				recrutTitle = procRset.Row[3].(string)
				recrutSn = procRset.Row[4].(string)
				recrutJfMngCd = procRset.Row[5].(string)
				upJobGrp = procRset.Row[6].(string)
				jobGrp = procRset.Row[7].(string)
				confirmDt = procRset.Row[8].(string)
				beginDt = procRset.Row[9].(string)
				endDt = procRset.Row[10].(string)
				nm = procRset.Row[11].(string)
				sex = procRset.Row[12].(string)
				birth = procRset.Row[13].(string)
				moNo = procRset.Row[14].(string)
				email = procRset.Row[15].(string)
				lstEdu = procRset.Row[16].(string)
				lstEduDesc = procRset.Row[17].(string)
				carrGbn = procRset.Row[18].(string)
				carrGbnDesc = procRset.Row[19].(string)
				techQlftKnd = procRset.Row[20].(string)
				frgnLangAbltDesc = procRset.Row[21].(string)

				adminStatsPeriod05 = append(adminStatsPeriod05, models.AdminStatsPeriod05{
					EntpKoNm:         entpKoNm,
					EntpMemNo:        entpMemNo,
					JobFairCds:       strings.ReplaceAll(jobfairCds, ",", "</br>"),
					RecrutTitle:      recrutTitle,
					RecrutSn:         recrutSn,
					RecrutJfMngCd:    recrutJfMngCd,
					UpJobGrp:         upJobGrp,
					JobGrp:           jobGrp,
					ConfirmDt:        confirmDt,
					BeginDt:          beginDt,
					EndDt:            endDt,
					Nm:               nm,
					Sex:              sex,
					Birth:            birth,
					MoNo:             moNo,
					Email:            email,
					LstEdu:           lstEdu,
					LstEduDesc:       lstEduDesc,
					CarrGbn:          carrGbn,
					CarrGbnDesc:      carrGbnDesc,
					TechQlftKnd:      techQlftKnd,
					FrgnLangAbltDesc: frgnLangAbltDesc,
				})
			}
			if err := procRset.Err(); err != nil {
				panic(err)
			}
			rtnAdminStatsPeriod05 = models.RtnAdminStatsPeriod05{
				Result: adminStatsPeriod05,
			}
			// End : Admin Stats Sub2

			c.Data["json"] = &rtnAdminStatsPeriod05
			c.ServeJSON()
		}
	} else if pType == "06" {
		logs.Debug(fmt.Sprintf("CALL SP_EMS_API_PERIOD_TYPE_R('%v', '%v', '%v', '%v', '%v', '%v', :1)",
			pLang, pType, pSdy, pEdy, pEntpKoNm, pJobFairMngCd))

		stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_API_PERIOD_TYPE_R('%v', '%v', '%v', '%v', '%v', '%v', :1)",
			pLang, pType, pSdy, pEdy, pEntpKoNm, pJobFairMngCd),
			ora.S, /* ENTP_KO_NM */
			ora.S, /* ENTP_MEM_NO */
			ora.S, /* RECRUT_TITLE */
			ora.S, /* RECRUT_SN */
			ora.S, /* JF_MNG_CD */
			ora.S, /* UP_JOB_GRP */
			ora.S, /* JOB_GRP */
			ora.S, /* REQUEST_DT */
			ora.S, /* BEGIN_DT */
			ora.S, /* END_DT */
			ora.S, /* LIVE_STATE */
			ora.S, /* APPLY_LIST */
			ora.S, /* MEM_LIST */
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

		rtnAdminStatsPeriod06 := models.RtnAdminStatsPeriod06{}
		adminStatsPeriod06 := make([]models.AdminStatsPeriod06, 0)

		var (
			entpKoNm      string
			entpMemNo     string
			jobfairCds    string
			recrutTitle   string
			recrutSn      string
			recrutJfMngCd string
			upJobGrp      string
			jobGrp        string
			requestDt     string
			beginDt       string
			endDt         string
			liveState     string
			applyList     string
			memList       string
		)

		if procRset.IsOpen() {
			for procRset.Next() {
				entpKoNm = procRset.Row[0].(string)
				entpMemNo = procRset.Row[1].(string)
				jobfairCds = procRset.Row[2].(string)
				recrutTitle = procRset.Row[3].(string)
				recrutSn = procRset.Row[4].(string)
				recrutJfMngCd = procRset.Row[5].(string)
				upJobGrp = procRset.Row[6].(string)
				jobGrp = procRset.Row[7].(string)
				requestDt = procRset.Row[8].(string)
				beginDt = procRset.Row[9].(string)
				endDt = procRset.Row[10].(string)
				liveState = procRset.Row[11].(string)
				applyList = procRset.Row[12].(string)
				memList = procRset.Row[13].(string)

				adminStatsPeriod06 = append(adminStatsPeriod06, models.AdminStatsPeriod06{
					EntpKoNm:      entpKoNm,
					EntpMemNo:     entpMemNo,
					JobFairCds:    strings.ReplaceAll(jobfairCds, ",", "</br>"),
					RecrutTitle:   recrutTitle,
					RecrutSn:      recrutSn,
					RecrutJfMngCd: recrutJfMngCd,
					UpJobGrp:      upJobGrp,
					JobGrp:        jobGrp,
					RequestDt:     requestDt,
					BeginDt:       beginDt,
					EndDt:         endDt,
					LiveState:     liveState,
					ApplyList:     strings.ReplaceAll(applyList, ",", "</br>"),
					MemList:       strings.ReplaceAll(memList, ",", "</br>"),
				})
			}
			if err := procRset.Err(); err != nil {
				panic(err)
			}
			rtnAdminStatsPeriod06 = models.RtnAdminStatsPeriod06{
				Result: adminStatsPeriod06,
			}
			// End : Admin Stats Sub2

			c.Data["json"] = &rtnAdminStatsPeriod06
			c.ServeJSON()
		}
	}
}
