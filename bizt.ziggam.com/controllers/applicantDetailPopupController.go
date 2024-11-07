package controllers

import (
	"fmt"
	"time"

	"bizt.ziggam.com/models"
	"bizt.ziggam.com/tables"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type ApplicantDetailPopupController struct {
	BaseController
}

func (c *ApplicantDetailPopupController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/login")
		return
	}

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := mem_no
	pRecrutSn := c.GetString("recrut_sn")
	pPpMemNo := c.GetString("pp_mem_no")

	imgServer, _ := beego.AppConfig.String("viewpath")
	cdnPath, _ := beego.AppConfig.String("cdnpath")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Applicant Popup Info
	log.Debug(fmt.Sprintf("CALL ZSP_MEM_POP_INFO_R_V3('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_MEM_POP_INFO_R_V3('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo),
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* PTO_PATH */
		ora.S,   /* NM */
		ora.S,   /* SEX */
		ora.I64, /* AGE */
		ora.S,   /* EMAIL */
		ora.S,   /* MO_NO */
		ora.S,   /* PRGS_STAT_CD */
		ora.S,   /* UP_JOB_GRP */
		ora.S,   /* JOB_GRP */
		ora.S,   /* RECRUT_TITLE */
		ora.S,   /* SDY */
		ora.S,   /* EDY */
		ora.S,   /* LST_EDU */
		ora.S,   /* CARR_GBN */
		ora.S,   /* CARR_DESC */
		ora.S,   /* FRGN_LANG_ABLT_DESC */
		ora.S,   /* ATCH_DATA_PATH */
		ora.S,   /* TECH_QLFT_KND */
		ora.S,   /* ATCH_FILE_PATH */
		ora.S,   /* FAVR_APLY_PP_YN */
		ora.S,   /* EVL_PRGS_STAT_CD */
		ora.S,   /* EVL_PRGS_DT */
		ora.S,   /* LIVE_REQ_STAT_CD */
		ora.S,   /* MSG_YN */
		ora.S,   /* MSG_END_YN */
		ora.S,   /* APPLY_DT */
		ora.S,   /* ENTP_GROUP_CODE */
		ora.S,   /* DCMNT_EVL_STAT_CD */
		ora.S,   /* ONWY_INTRV_EVL_STAT_CD */
		ora.S,   /* LIVE_INTRV_EVL_STAT_CD */
		ora.S,   /* DCMNT_FILE_PATH */
		ora.S,   /* DCMNT_EVL_STAT_DT */
		ora.S,   /* ONWY_INTRV_EVL_STAT_DT */
		ora.S,   /* LIVE_INTRV_EVL_STAT_DT */
		ora.S,   /* DCMNT_FILE_NAME */
		ora.S,   /* RECRUT_PROC_CD */
		ora.S,   /* LST_EDU_GBN_CD1 */    // 최종 학력 1: 신규
		ora.S,   /* LST_EDU_GBN_CD2 */    // 최종 학력 2: 신규
		ora.S,   /* CARR_YEAR */          // 연차 연도 : 신규
		ora.S,   /* FRGN_LANG_ABLT_CD1 */ // 외국어 능력 1 : 신규
		ora.S,   /* FRGN_LANG_ABLT_CD2 */ // 외국어 능력 2 : 신규
		ora.I64, /* SHOOT_CNT */
		ora.S,   /* COMP_TM */
		ora.S,   /* COMP_DT1 */
		ora.S,   /* COMP_DT2 */
		ora.S,   /* READ_90_OVER */
		ora.S,   /* LIVE_SN */
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

	applicantDetailPopup := make([]models.ApplicantDetailPopup, 0)

	var (
		entpMemNo              string
		recrutSn               string
		ppMemNo                string
		ptoPath                string
		nm                     string
		sex                    string
		age                    int64
		email                  string
		moNo                   string
		prgsStatCd             string
		upJobGrp               string
		jobGrp                 string
		recrutTitle            string
		sdy                    string
		edy                    string
		lstEdu                 string
		carrGbn                string
		carrDesc               string
		frgnLangAbltDesc       string
		atchDataPath           string
		techQlftKnd            string
		atchFilePath           string
		favrAplyPpYn           string
		evlPrgsStatCd          string
		evlStatDt              string
		liveReqStatCd          string
		msgYn                  string
		msgEndYn               string
		applyDt                string
		entpGroupCode          string
		fullPtoPath            string
		fullFilePath           string
		dcmnt_evl_stat_cd      string
		onwy_intrv_evl_stat_cd string
		live_intrv_evl_stat_cd string
		dcmnt_file_path        string
		dcmnt_evl_stat_dt      string
		onwy_intrv_evl_stat_dt string
		live_intrv_evl_stat_dt string
		dcmnt_file_name        string
		recrut_proc_cd         string
		lstEduGbnCd1           string // 최종 학력 1: 신규
		lstEduGbnCd2           string // 최종 학력 2: 신규
		carrYear               string // 연차 연도 : 신규
		frgnLangAbltCd1        string // 외국어 능력 1 : 신규
		frgnLangAbltCd2        string // 외국어 능력 2 : 신규
		shootCnt               int64
		compTm                 string
		compDT1                string
		compDT2                string
		edy_90_over            string
		liveSn                 string // 다대다 라이브
	)

	if procRset.IsOpen() {
		for procRset.Next() {

			entpMemNo = procRset.Row[0].(string)
			recrutSn = procRset.Row[1].(string)
			ppMemNo = procRset.Row[2].(string)
			ptoPath = procRset.Row[3].(string)
			if ptoPath == "" {
				fullPtoPath = ptoPath
			} else {
				fullPtoPath = imgServer + ptoPath
			}
			nm = procRset.Row[4].(string)
			sex = procRset.Row[5].(string)
			age = procRset.Row[6].(int64)

			email = procRset.Row[7].(string)
			moNo = procRset.Row[8].(string)
			prgsStatCd = procRset.Row[9].(string)
			upJobGrp = procRset.Row[10].(string)
			jobGrp = procRset.Row[11].(string)
			recrutTitle = procRset.Row[12].(string)
			sdy = procRset.Row[13].(string)
			edy = procRset.Row[14].(string)
			lstEdu = procRset.Row[15].(string)
			carrGbn = procRset.Row[16].(string)
			carrDesc = procRset.Row[17].(string)
			frgnLangAbltDesc = procRset.Row[18].(string)
			atchDataPath = procRset.Row[19].(string)
			techQlftKnd = procRset.Row[20].(string)
			atchFilePath = procRset.Row[21].(string)
			if atchFilePath == "" {
				fullFilePath = atchFilePath
			} else {
				fullFilePath = imgServer + atchFilePath
			}
			favrAplyPpYn = procRset.Row[22].(string)
			evlPrgsStatCd = procRset.Row[23].(string)
			evlStatDt = procRset.Row[24].(string)
			liveReqStatCd = procRset.Row[25].(string)
			msgYn = procRset.Row[26].(string)
			msgEndYn = procRset.Row[27].(string)
			applyDt = procRset.Row[28].(string)
			entpGroupCode = procRset.Row[29].(string)
			dcmnt_evl_stat_cd = procRset.Row[30].(string)
			onwy_intrv_evl_stat_cd = procRset.Row[31].(string)
			live_intrv_evl_stat_cd = procRset.Row[32].(string)
			dcmnt_file_path = procRset.Row[33].(string)
			dcmnt_evl_stat_dt = procRset.Row[34].(string)
			onwy_intrv_evl_stat_dt = procRset.Row[35].(string)
			live_intrv_evl_stat_dt = procRset.Row[36].(string)
			dcmnt_file_name = procRset.Row[37].(string)
			recrut_proc_cd = procRset.Row[38].(string)

			lstEduGbnCd1 = procRset.Row[39].(string)    // 최종 학력 1: 신규
			lstEduGbnCd2 = procRset.Row[40].(string)    // 최종 학력 2: 신규
			carrYear = procRset.Row[41].(string)        // 연차 연도 : 신규
			frgnLangAbltCd1 = procRset.Row[42].(string) // 외국어 능력 1 : 신규
			frgnLangAbltCd2 = procRset.Row[43].(string) // 외국어 능력 2 : 신규

			shootCnt = procRset.Row[44].(int64)
			compTm = procRset.Row[45].(string)
			compDT1 = procRset.Row[46].(string)
			compDT2 = procRset.Row[47].(string)
			edy_90_over = procRset.Row[48].(string)
			liveSn = procRset.Row[49].(string)

			applicantDetailPopup = append(applicantDetailPopup, models.ApplicantDetailPopup{
				EntpMemNo:              entpMemNo,
				RecrutSn:               recrutSn,
				PpMemNo:                ppMemNo,
				PtoPath:                fullPtoPath,
				Nm:                     nm,
				Sex:                    sex,
				Age:                    age,
				Email:                  email,
				MoNo:                   moNo,
				PrgsStatCd:             prgsStatCd,
				UpJobGrp:               upJobGrp,
				JobGrp:                 jobGrp,
				RecrutTitle:            recrutTitle,
				Sdy:                    sdy,
				Edy:                    edy,
				LstEdu:                 lstEdu,
				CarrGbn:                carrGbn,
				CarrDesc:               carrDesc,
				FrgnLangAbltDesc:       frgnLangAbltDesc,
				AtchDataPath:           atchDataPath,
				TechQlftKnd:            techQlftKnd,
				AtchFilePath:           fullFilePath,
				FavrAplyPpYn:           favrAplyPpYn,
				EvlPrgsStatCd:          evlPrgsStatCd,
				EvlStatDt:              evlStatDt,
				LiveReqStatCd:          liveReqStatCd,
				MsgYn:                  msgYn,
				MsgEndYn:               msgEndYn,
				ApplyDt:                applyDt,
				EntpGroupCode:          entpGroupCode,
				Dcmnt_evl_stat_cd:      dcmnt_evl_stat_cd,
				Onwy_intrv_evl_stat_cd: onwy_intrv_evl_stat_cd,
				Live_intrv_evl_stat_cd: live_intrv_evl_stat_cd,
				Dcmnt_file_path:        dcmnt_file_path,
				Dcmnt_evl_stat_dt:      dcmnt_evl_stat_dt,
				Onwy_intrv_evl_stat_dt: onwy_intrv_evl_stat_dt,
				Live_intrv_evl_stat_dt: live_intrv_evl_stat_dt,
				Dcmnt_file_name:        dcmnt_file_name,
				Recrut_proc_cd:         recrut_proc_cd,
				ShootCnt:               shootCnt,
				CompTm:                 compTm,
				CompDT1:                compDT1,
				CompDT2:                compDT2,
				LiveSn:                 liveSn,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}

	log.Debug(tables.MapLstEduGbnCd1[lstEduGbnCd1])
	log.Debug(tables.MapLstEduGbnCd2[lstEduGbnCd2])
	log.Debug(carrYear)
	log.Debug(tables.MapFrgnLangAbltCd1[frgnLangAbltCd1])
	log.Debug(tables.MapFrgnLangAbltCd2[frgnLangAbltCd2])
	// End : Applicant Popup Info

	var ansTotCnt int64
	recruitApplyMemberAnswerList := make([]models.RecruitApplyMemberAnswerList, 0)

	// 채용 공고가 90일이 지난 경우 영상 답변은 더이상 보여주어서는 안된다.
	if edy_90_over != "Y" {
		// Start : Recruit Apply Answer List
		log.Debug(fmt.Sprintf("CALL ZSP_QST_ANS_LIST_R('%v', '%v', '%v', '%v', :1)",
			pLang, pEntpMemNo, pRecrutSn, pPpMemNo))

		stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_QST_ANS_LIST_R('%v', '%v', '%v', '%v', :1)",
			pLang, pEntpMemNo, pRecrutSn, pPpMemNo),
			ora.S,   /* ENTP_MEM_NO */
			ora.S,   /* RECRUT_SN */
			ora.S,   /* QST_SN */
			ora.S,   /* VD_TITLE */
			ora.S,   /* VD_FILE_PATH */
			ora.I64, /* TOT_CNT */

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

		//	recruitApplyMemberAnswerList := make([]models.RecruitApplyMemberAnswerList, 0)

		var (
			ansEntpMemNo  string
			ansRecrutSn   string
			ansQstSn      string
			ansVdTitle    string
			ansVdFilePath string
			//			ansTotCnt     int64
			vdFullPtoPath string
		)

		if procRset.IsOpen() {
			for procRset.Next() {
				ansEntpMemNo = procRset.Row[0].(string)
				ansRecrutSn = procRset.Row[1].(string)
				ansQstSn = procRset.Row[2].(string)
				ansVdTitle = procRset.Row[3].(string)
				ansVdFilePath = procRset.Row[4].(string)
				ansTotCnt = procRset.Row[5].(int64)

				if ansVdFilePath == "" {
					vdFullPtoPath = ansVdFilePath
				} else {
					vdFullPtoPath = cdnPath + ansVdFilePath
				}

				recruitApplyMemberAnswerList = append(recruitApplyMemberAnswerList, models.RecruitApplyMemberAnswerList{
					AnsEntpMemNo:  ansEntpMemNo,
					AnsRecrutSn:   ansRecrutSn,
					AnsQstSn:      ansQstSn,
					AnsVdTitle:    ansVdTitle,
					AnsVdFilePath: vdFullPtoPath,
					AnsTotCnt:     ansTotCnt,
				})
			}
			if err := procRset.Err(); err != nil {
				panic(err)
			}
		}
		// End : Recruit Apply Answer List
	} else {
		c.TplName = "applicant/applicant_detail_popup_abort.html"
		return
	}

	// Start : Recruit Apply Comment List

	cmtSMemId := session.Get(c.Ctx.Request.Context(), "mem_id")
	cmtSAuthCd := session.Get(c.Ctx.Request.Context(), "auth_cd")

	log.Debug(fmt.Sprintf("CALL ZSP_COMMENT_LIST_R('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_COMMENT_LIST_R('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo),
		ora.I64, /* TOT_CNT */
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* PP_CHRG_CMT_SN */
		ora.S,   /* PP_CHRG_SN */
		ora.S,   /* PP_CHRG_CMT */
		ora.S,   /* REG_DT */
		ora.S,   /* REG_ID */
		ora.S,   /* PP_CHRG_BP_NM */
		ora.S,   /* PP_CHRG_NM */
		ora.S,   /* NEW_YN */
		ora.S,   /* PP_CHRG_GBN_CD */
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

	recruitApplyCommentList := make([]models.RecruitApplyCommentList, 0)

	var (
		cmtTotCnt      int64
		cmtEntpMemNo   string
		cmtRecrutSn    string
		cmtPpMemNo     string
		cmtPpChrgCmtSn string
		cmtPpChrgSn    string
		cmtPpChrgCmt   string
		cmtRegDt       string
		cmtRegId       string
		cmtPpChrgBpNm  string
		cmtPpChrgNm    string
		cmtNewYn       string
		cmtPpChrgGbnCd string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			cmtTotCnt = procRset.Row[0].(int64)
			cmtEntpMemNo = procRset.Row[1].(string)
			cmtRecrutSn = procRset.Row[2].(string)
			cmtPpMemNo = procRset.Row[3].(string)
			cmtPpChrgCmtSn = procRset.Row[4].(string)
			cmtPpChrgSn = procRset.Row[5].(string)
			cmtPpChrgCmt = procRset.Row[6].(string)
			cmtRegDt = procRset.Row[7].(string)
			cmtRegId = procRset.Row[8].(string)
			cmtPpChrgBpNm = procRset.Row[9].(string)
			cmtPpChrgNm = procRset.Row[10].(string)
			cmtNewYn = procRset.Row[11].(string)
			cmtPpChrgGbnCd = procRset.Row[12].(string)

			recruitApplyCommentList = append(recruitApplyCommentList, models.RecruitApplyCommentList{
				CmtTotCnt:      cmtTotCnt,
				CmtEntpMemNo:   cmtEntpMemNo,
				CmtRecrutSn:    cmtRecrutSn,
				CmtPpMemNo:     cmtPpMemNo,
				CmtPpChrgCmtSn: cmtPpChrgCmtSn,
				CmtPpChrgSn:    cmtPpChrgSn,
				CmtPpChrgCmt:   cmtPpChrgCmt,
				CmtRegDt:       cmtRegDt,
				CmtRegId:       cmtRegId,
				CmtPpChrgBpNm:  cmtPpChrgBpNm,
				CmtPpChrgNm:    cmtPpChrgNm,
				CmtNewYn:       cmtNewYn,
				CmtPpChrgGbnCd: cmtPpChrgGbnCd,
				CmtSMemId:      cmtSMemId,
				CmtSAuthCd:     cmtSAuthCd,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Recruit Apply Comment List

	// Start : Member Video Profile List

	log.Debug(fmt.Sprintf("CALL ZSP_VP_LIST_R('%v', '%v', :1)", pLang, pPpMemNo))
	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_VP_LIST_R('%v', '%v', :1)",
		pLang, pPpMemNo),
		ora.S, /* VD_SN */
		ora.S, /* THM_KND_CD */
		ora.S, /* THM_NM */
		ora.S, /* QST_CD */
		ora.S, /* QST_DESC */
		ora.S, /* OPN_SET_CD */
		ora.S, /* OPN_SET_NM */
		ora.S, /* VD_FILE_PATH */
		ora.S, /* VD_THUMB_PATH */
		ora.S, /* REG_DT */
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

	memberVideoProfileList := make([]models.MemberVideoProfileList, 0)

	var (
		mvVdSn        string
		mvThmKndCd    string
		mvThmNm       string
		mvQstSCd      string
		mvQstDesc     string
		mvOpnSetCd    string
		mvOpenSetNm   string
		mvVdFilePath  string
		mvVdThumbPath string
		mvRegDt       string
		vpFullPtoPath string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			mvVdSn = procRset.Row[0].(string)
			mvThmKndCd = procRset.Row[1].(string)
			mvThmNm = procRset.Row[2].(string)
			mvQstSCd = procRset.Row[3].(string)
			mvQstDesc = procRset.Row[4].(string)
			mvOpnSetCd = procRset.Row[5].(string)
			mvOpenSetNm = procRset.Row[6].(string)
			mvVdFilePath = procRset.Row[7].(string)
			if mvVdFilePath == "" {
				vpFullPtoPath = mvVdFilePath
			} else {
				vpFullPtoPath = cdnPath + mvVdFilePath
			}
			mvVdThumbPath = procRset.Row[8].(string)
			mvRegDt = procRset.Row[9].(string)

			memberVideoProfileList = append(memberVideoProfileList, models.MemberVideoProfileList{
				MvVdSn:        mvVdSn,
				MvThmKndCd:    mvThmKndCd,
				MvThmNm:       mvThmNm,
				MvQstSCd:      mvQstSCd,
				MvQstDesc:     mvQstDesc,
				MvOpnSetCd:    mvOpnSetCd,
				MvOpenSetNm:   mvOpenSetNm,
				MvVdFilePath:  vpFullPtoPath,
				MvVdThumbPath: mvVdThumbPath,
				MvRegDt:       mvRegDt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Member Video Profile List

	// Start : Entp Team Member List

	pGbnCd := "A"

	log.Debug(fmt.Sprintf("CALL ZSP_TEAM_MEM_LIST_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pGbnCd))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_TEAM_MEM_LIST_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pGbnCd),
		ora.I64, /* TOT_CNT */
		ora.S,   /* PP_CHRG_SN */
		ora.S,   /* PP_CHRG_GBN_CD */
		ora.S,   /* PP_CHRG_NM */
		ora.S,   /* PP_CHRG_BP_NM */
		ora.S,   /* EMAIL */
		ora.S,   /* ENTP_MEM_ID */
		ora.S,   /* PP_CHRG_TEL_NO */
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

	entpTeamMemberList := make([]models.EntpTeamMemberList, 0)

	var (
		etTotCnt      int64
		etPpChrgSn    string
		etPpChrgGbnCd string
		etPpChrgNm    string
		etPpChrgBpNm  string
		etEmail       string
		etEntpMemId   string
		etPpChrgTelNo string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			etTotCnt = procRset.Row[0].(int64)
			etPpChrgSn = procRset.Row[1].(string)
			etPpChrgGbnCd = procRset.Row[2].(string)
			etPpChrgNm = procRset.Row[3].(string)
			etPpChrgBpNm = procRset.Row[4].(string)
			etEmail = procRset.Row[5].(string)
			etEntpMemId = procRset.Row[6].(string)
			etPpChrgTelNo = procRset.Row[7].(string)

			entpTeamMemberList = append(entpTeamMemberList, models.EntpTeamMemberList{
				EtTotCnt:      etTotCnt,
				EtPpChrgSn:    etPpChrgSn,
				EtPpChrgGbnCd: etPpChrgGbnCd,
				EtPpChrgNm:    etPpChrgNm,
				EtPpChrgBpNm:  etPpChrgBpNm,
				EtEmail:       etEmail,
				EtEntpMemId:   etEntpMemId,
				EtPpChrgTelNo: etPpChrgTelNo,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Entp Team Member List

	c.Data["EntpTeamMemberList"] = entpTeamMemberList
	c.Data["RecruitApplyMemberAnswerList"] = recruitApplyMemberAnswerList
	c.Data["RecruitApplyCommentList"] = recruitApplyCommentList
	c.Data["MemberVideoProfileList"] = memberVideoProfileList

	c.Data["AnsTotCnt"] = ansTotCnt

	c.Data["CmtTotCnt"] = cmtTotCnt
	c.Data["EntpMemNo"] = entpMemNo
	c.Data["RecrutSn"] = recrutSn
	c.Data["PpMemNo"] = ppMemNo
	c.Data["PtoPath"] = fullPtoPath
	c.Data["Nm"] = nm
	c.Data["Sex"] = sex
	c.Data["Age"] = age
	c.Data["Email"] = email
	c.Data["MoNo"] = moNo
	c.Data["PrgsStatCd"] = prgsStatCd
	c.Data["UpJobGrp"] = upJobGrp
	c.Data["JobGrp"] = jobGrp
	c.Data["RecrutTitle"] = recrutTitle
	c.Data["Sdy"] = sdy
	c.Data["Edy"] = edy
	c.Data["LstEdu"] = lstEdu
	c.Data["CarrGbn"] = carrGbn
	c.Data["CarrDesc"] = carrDesc
	c.Data["FrgnLangAbltDesc"] = frgnLangAbltDesc
	c.Data["AtchDataPath"] = atchDataPath
	c.Data["TechQlftKnd"] = techQlftKnd
	c.Data["AtchFilePath"] = fullFilePath
	c.Data["FavrAplyPpYn"] = favrAplyPpYn
	c.Data["EvlPrgsStatCd"] = evlPrgsStatCd
	c.Data["EvlStatDt"] = evlStatDt
	c.Data["LiveReqStatCd"] = liveReqStatCd
	c.Data["MsgYn"] = msgYn
	c.Data["MsgEndYn"] = msgEndYn
	c.Data["ApplyDt"] = applyDt
	c.Data["EntpGroupCode"] = entpGroupCode

	c.Data["Dcmnt_evl_stat_cd"] = dcmnt_evl_stat_cd
	c.Data["Onwy_intrv_evl_stat_cd"] = onwy_intrv_evl_stat_cd
	c.Data["Live_intrv_evl_stat_cd"] = live_intrv_evl_stat_cd
	c.Data["Dcmnt_file_path"] = imgServer + dcmnt_file_path
	c.Data["Dcmnt_evl_stat_dt"] = dcmnt_evl_stat_dt
	c.Data["Onwy_intrv_evl_stat_dt"] = onwy_intrv_evl_stat_dt
	c.Data["Live_intrv_evl_stat_dt"] = live_intrv_evl_stat_dt
	c.Data["Dcmnt_file_name"] = dcmnt_file_name
	c.Data["Recrut_proc_cd"] = recrut_proc_cd

	convTime, _ := time.Parse("20060102150405", dcmnt_evl_stat_dt)
	dcmntEvlStatDt := fmt.Sprintf(convTime.Format("06/01/02 15:04"))
	c.Data["DcmntEvlStatDtFmt"] = dcmntEvlStatDt

	// LDK 2020/08/28: 개인 정보 코드화
	c.Data["LstEduGbnCd1"] = tables.MapLstEduGbnCd1[lstEduGbnCd1]
	c.Data["LstEduGbnCd2"] = tables.MapLstEduGbnCd2[lstEduGbnCd2]

	c.Data["CarrYear"] = carrYear

	c.Data["FrgnLangAbltCd1"] = tables.MapFrgnLangAbltCd1[frgnLangAbltCd1]
	c.Data["FrgnLangAbltCd2"] = tables.MapFrgnLangAbltCd2[frgnLangAbltCd2]

	c.Data["FrgnLangAbltDesc"] = frgnLangAbltDesc
	// <--

	c.Data["ShootCnt"] = shootCnt
	c.Data["CompTm"] = compTm
	c.Data["CompDT1"] = compDT1
	c.Data["CompDT2"] = compDT2

	// LDK 2021/01/12: 다대다 라이브
	c.Data["LiveSn"] = liveSn
	// <--

	c.TplName = "applicant/applicant_detail_popup.html"
}
