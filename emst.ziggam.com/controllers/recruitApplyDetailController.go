package controllers

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"emst.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

type RecruitApplyDetailController struct {
	BaseController
}

func (c *RecruitApplyDetailController) Get() {

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
	pEntpMemNo := mem_no
	pChkEntpMemNo := c.GetString("entp_mem_no") // 체크 기업회원번호

	if pChkEntpMemNo != "" {
		if pEntpMemNo != pChkEntpMemNo {
			c.Ctx.Redirect(302, "/error/404")
		}
	}

	pRecrutSn := c.GetString("recrut_sn")
	pPpMemNo := c.GetString("pp_mem_no")

	/* Parameter */
	pmRecrutSn := c.GetString("p_recrut_sn")
	pmPpMemNo := c.GetString("p_pp_mem_no")
	pmEvlPrgsStat := c.GetString("p_evl_prgs_stat")
	pmSex := c.GetString("p_sex")
	pmAge := c.GetString("p_age")
	pmVpYn := c.GetString("p_vp_yn")
	pmFavrAplyPp := c.GetString("p_favr_aply_pp")
	pmSortGbn := c.GetString("p_sort_gbn")
	pmPageNo := c.GetString("p_page_no")
	pmKeyword := c.GetString("p_keyword")
	pmSize := c.GetString("p_size")

	var fullPtoPath string
	var vdFullPtoPath string
	var fullFilePath string
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

	// Start : Recruit Apply Detail
	log.Debug("CALL SP_EMS_AM_DTL_INFO_R('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_AM_DTL_INFO_R('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo),
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* PTO_PATH */
		ora.S,   /* FAVR_APLY_PP_YN */
		ora.S,   /* NM */
		ora.S,   /* SEX */
		ora.S,   /* AGE */
		ora.S,   /* EMAIL */
		ora.S,   /* APPLY_DT */
		ora.S,   /* LEFT_DY */
		ora.S,   /* SHOOT_TM */
		ora.I64, /* SHOOT_CNT */
		ora.S,   /* VP_YN */
		ora.S,   /* LST_EDU */
		ora.S,   /* CARR_GBN */
		ora.S,   /* CARR_DESC */
		ora.S,   /* FRGN_LANG_ABLT_DESC */
		ora.S,   /* ATCH_DATA_PATH */
		ora.S,   /* TECH_QLFT_KND */
		ora.S,   /* ATCH_FILE_PATH */
		ora.S,   /* MO_NO */
		ora.S,   /* GROUP_CODE */

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

	recruitApplyDetail := make([]models.RecruitApplyDetail, 0)

	var (
		entpMemNo        string
		recrutSn         string
		ppMemNo          string
		ptoPath          string
		favrAplyPpYn     string
		nm               string
		sex              string
		age              string
		email            string
		applyDt          string
		leftDy           string
		shootTm          string
		shootCnt         int64
		vpYn             string
		lstEdu           string
		carrGbn          string
		carrDesc         string
		frgnLangAbltDesc string
		atchDataPath     string
		techQlftKnd      string
		atchFilePath     string
		moNo             string
		entpGroupCode    string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			entpMemNo = procRset.Row[0].(string)
			recrutSn = procRset.Row[1].(string)
			ppMemNo = procRset.Row[2].(string)
			ptoPath = procRset.Row[3].(string)
			favrAplyPpYn = procRset.Row[4].(string)
			nm = procRset.Row[5].(string)
			sex = procRset.Row[6].(string)
			age = procRset.Row[7].(string)
			email = procRset.Row[8].(string)
			applyDt = procRset.Row[9].(string)
			leftDy = procRset.Row[10].(string)
			shootTm = procRset.Row[11].(string)
			shootCnt = procRset.Row[12].(int64)
			vpYn = procRset.Row[13].(string)
			lstEdu = procRset.Row[14].(string)
			carrGbn = procRset.Row[15].(string)
			carrDesc = procRset.Row[16].(string)
			frgnLangAbltDesc = procRset.Row[17].(string)
			atchDataPath = procRset.Row[18].(string)
			techQlftKnd = procRset.Row[19].(string)
			atchFilePath = procRset.Row[20].(string)
			moNo = procRset.Row[21].(string)
			entpGroupCode = procRset.Row[22].(string)

			if ptoPath == "" {
				fullPtoPath = ptoPath
			} else {
				fullPtoPath = imgServer + ptoPath
			}

			if atchFilePath == "" {
				fullFilePath = atchFilePath
			} else {
				fullFilePath = imgServer + atchFilePath
			}

			recruitApplyDetail = append(recruitApplyDetail, models.RecruitApplyDetail{
				EntpMemNo:        entpMemNo,
				RecrutSn:         recrutSn,
				PpMemNo:          ppMemNo,
				PtoPath:          fullPtoPath,
				FavrAplyPpYn:     favrAplyPpYn,
				Nm:               nm,
				Sex:              sex,
				Age:              age,
				Email:            email,
				ApplyDt:          applyDt,
				LeftDy:           leftDy,
				ShootTm:          shootTm,
				ShootCnt:         shootCnt,
				VpYn:             vpYn,
				LstEdu:           lstEdu,
				CarrGbn:          carrGbn,
				CarrDesc:         carrDesc,
				FrgnLangAbltDesc: frgnLangAbltDesc,
				AtchDataPath:     atchDataPath,
				TechQlftKnd:      techQlftKnd,
				AtchFilePath:     fullFilePath,
				MoNo:             moNo,
				EntpGroupCode:    entpGroupCode,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Recruit Apply Detail

	// Start : Recruit Apply Top Info
	log.Debug("CALL SP_EMS_AM_RECRUIT_INFO_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_AM_RECRUIT_INFO_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn),
		ora.S, /* PRGS_STAT */
		ora.S, /* RECRUT_TITLE */
		ora.S, /* EMPL_TYP */
		ora.S, /* UP_JOB_GRP */
		ora.S, /* JOB_GRP */
		ora.S, /* RECRUT_DY */
		ora.S, /* RECRUT_EDT */
		ora.S, /* PRGS_MSG */

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

	recruitApplyTopInfo := make([]models.RecruitApplyTopInfo, 0)

	var (
		prgsStat    string
		recrutTitle string
		emplTyp     string
		upJobGrp    string
		jobGrp      string
		recrutDy    string
		recrutEdt   string
		prgsMsg     string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			prgsStat = procRset.Row[0].(string)
			recrutTitle = procRset.Row[1].(string)
			emplTyp = procRset.Row[2].(string)
			upJobGrp = procRset.Row[3].(string)
			jobGrp = procRset.Row[4].(string)
			recrutDy = procRset.Row[5].(string)
			recrutEdt = procRset.Row[6].(string)
			prgsMsg = procRset.Row[7].(string)

			recruitApplyTopInfo = append(recruitApplyTopInfo, models.RecruitApplyTopInfo{
				PrgsStat:    prgsStat,
				RecrutTitle: recrutTitle,
				EmplTyp:     emplTyp,
				UpJobGrp:    upJobGrp,
				JobGrp:      jobGrp,
				RecrutDy:    recrutDy,
				RecrutEdt:   recrutEdt,
				PrgsMsg:     prgsMsg,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Recruit Apply Top Info

	// Start : Recruit Apply Answer List
	log.Debug("CALL SP_EMS_AM_ANSWER_LIST_R('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_AM_ANSWER_LIST_R('%v', '%v', '%v', '%v', :1)",
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

	recruitApplyMemberAnswerList := make([]models.RecruitApplyMemberAnswerList, 0)

	var (
		ansEntpMemNo  string
		ansRecrutSn   string
		ansQstSn      string
		ansVdTitle    string
		ansVdFilePath string
		ansTotCnt     int64
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

	// Start : Profile Video List

	log.Debug("CALL SP_EMS_PROFILE_VIDEO_LIST_R('%v', '%v',:1)",
		pLang, pPpMemNo)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_PROFILE_VIDEO_LIST_R('%v', '%v',:1)",
		pLang, pPpMemNo),
		ora.S,   /* VD_SN */
		ora.I64, /* VD_SEC */
		ora.S,   /* VD_THUMB_PATH */
		ora.S,   /* VD_FILE_PATH */
		ora.S,   /* THM_KND_CD */
		ora.S,   /* THM_NM */
		ora.S,   /* QST_CD */
		ora.S,   /* QST_DESC */
		ora.S,   /* OPN_SET_CD */
		ora.S,   /* REG_DT */
		ora.I64, /* TOT_CNT */
		ora.I64, /* SN */
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

	videoProfileList := make([]models.VideoProfileList, 0)

	var (
		vpVdsn        string
		vpVdSec       int64
		vpVdThumbPath string
		vpVdFilePath  string
		vpThmKndCd    string
		vpThmNm       string
		vpQstCd       string
		vpQstDesc     string
		vpOpnSetCd    string
		vpRegDt       string
		vpTotCnt      int64
		vpSn          int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			vpVdsn = procRset.Row[0].(string)
			vpVdSec = procRset.Row[1].(int64)
			vpVdThumbPath = procRset.Row[2].(string)
			vpVdFilePath = procRset.Row[3].(string)
			vpThmKndCd = procRset.Row[4].(string)
			vpThmNm = procRset.Row[5].(string)
			vpQstCd = procRset.Row[6].(string)
			vpQstDesc = procRset.Row[7].(string)
			vpOpnSetCd = procRset.Row[8].(string)
			vpRegDt = procRset.Row[9].(string)
			vpTotCnt = procRset.Row[10].(int64)
			vpSn = procRset.Row[11].(int64)

			var fullThumbPath string
			var fullFilePath string

			if vpVdThumbPath == "" {
				fullThumbPath = vpVdThumbPath
			} else {
				fullThumbPath = imgServer + vpVdThumbPath
			}

			if vpVdFilePath == "" {
				fullFilePath = vpVdFilePath
			} else {
				fullFilePath = cdnPath + vpVdFilePath
			}

			videoProfileList = append(videoProfileList, models.VideoProfileList{
				VpVdsn:        vpVdsn,
				VpVdSec:       vpVdSec,
				VpVdThumbPath: fullThumbPath,
				VpVdFilePath:  fullFilePath,
				VpThmKndCd:    vpThmKndCd,
				VpThmNm:       vpThmNm,
				VpQstCd:       vpQstCd,
				VpQstDesc:     vpQstDesc,
				VpOpnSetCd:    vpOpnSetCd,
				VpRegDt:       vpRegDt,
				VpTotCnt:      vpTotCnt,
				VpSn:          vpSn,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Profile Video List

	c.Data["PrgsStat"] = prgsStat
	c.Data["RecrutTitle"] = recrutTitle
	c.Data["EmplTyp"] = emplTyp
	c.Data["UpJobGrp"] = upJobGrp
	c.Data["JobGrp"] = jobGrp
	c.Data["RecrutDy"] = recrutDy
	c.Data["RecrutEdt"] = recrutEdt
	c.Data["PrgsMsg"] = prgsMsg

	c.Data["VideoProfileList"] = videoProfileList
	c.Data["ProfileVideoCnt"] = vpTotCnt
	c.Data["RecruitApplyMemberAnswerList"] = recruitApplyMemberAnswerList
	c.Data["VideoCnt"] = ansTotCnt

	c.Data["EntpMemNo"] = pEntpMemNo
	c.Data["RecrutSn"] = recrutSn
	c.Data["PpMemNo"] = ppMemNo
	c.Data["PtoPath"] = fullPtoPath
	c.Data["FavrAplyPpYn"] = favrAplyPpYn
	c.Data["Nm"] = nm
	c.Data["Sex"] = sex
	c.Data["Age"] = age
	c.Data["Email"] = email
	c.Data["ApplyDt"] = applyDt
	c.Data["LeftDy"] = leftDy
	c.Data["ShootTm"] = shootTm
	c.Data["ShootCnt"] = shootCnt
	c.Data["VpYn"] = vpYn
	c.Data["LstEdu"] = lstEdu
	c.Data["CarrGbn"] = carrGbn
	c.Data["CarrDesc"] = carrDesc
	c.Data["FrgnLangAbltDesc"] = frgnLangAbltDesc
	c.Data["AtchDataPath"] = atchDataPath
	c.Data["TechQlftKnd"] = techQlftKnd
	c.Data["AtchFilePath"] = fullFilePath
	c.Data["MoNo"] = moNo
	c.Data["EntpGroupCode"] = entpGroupCode
	c.Data["MenuId"] = "03"

	/* Parameter Value */
	c.Data["pRecrutSn"] = pmRecrutSn
	c.Data["pPpMemNo"] = pmPpMemNo
	c.Data["pEvlPrgsStat"] = pmEvlPrgsStat
	c.Data["pSex"] = pmSex
	c.Data["pAge"] = pmAge
	c.Data["pVpYn"] = pmVpYn
	c.Data["pFavrAplyPp"] = pmFavrAplyPp
	c.Data["pSortGbn"] = pmSortGbn
	c.Data["pPageNo"] = pmPageNo
	c.Data["pKeyword"] = pmKeyword
	c.Data["pSize"] = pmSize

	c.TplName = "recruit/recruit_apply_detail.html"
}
