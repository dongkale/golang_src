package controllers

import (
	"fmt"
	"strings"

	"emst.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type AdminMemberDetailController struct {
	BaseController
}

func (c *AdminMemberDetailController) Get() {

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
	pPpMemNo := c.GetString("pp_mem_no")

	/* Parameter */
	pmMemStat := c.GetString("p_mem_stat")
	pmKeyword := c.GetString("p_keyword")
	pmSex := c.GetString("p_sex")
	pmAge := c.GetString("p_age")
	pmYpYn := c.GetString("p_vp_yn")
	pmOsGbn := c.GetString("p_os_gbn")
	pmPageNo := c.GetString("p_page_no")
	pmSize := c.GetString("p_size")
	pmSdy := c.GetString("p_sdy")
	pmEdy := c.GetString("p_edy")

	pmLoginSdy := c.GetString("p_login_sdy")
	pmLoginEdy := c.GetString("p_login_edy")

	var fullPtoPath string
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

	// Start : Member Info Detail
	log.Debug("CALL SP_EMS_MEMBER_DTL_INFO_R('%v', '%v', :1)",
		pLang, pPpMemNo)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_MEMBER_DTL_INFO_R('%v', '%v', :1)",
		pLang, pPpMemNo),
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* MEM_STAT_CD */
		ora.S,   /* MEM_STAT_NM */
		ora.S,   /* MEM_STAT_DT */
		ora.S,   /* MEM_ID */
		ora.S,   /* NM */
		ora.S,   /* SEX */
		ora.S,   /* EMAIL */
		ora.S,   /* MO_NO */
		ora.S,   /* BRTH_YMD */
		ora.S,   /* AGE */
		ora.S,   /* EMAIL_RECV_YN */
		ora.S,   /* SMS_RECV_YN */
		ora.S,   /* PTO_PATH */
		ora.S,   /* REG_DT */
		ora.S,   /* OS_GBN */
		ora.S,   /* OS_VER */
		ora.S,   /* M_REG_PRGS_STAT_CD */
		ora.S,   /* M_REG_PRGS_STAT_NM */
		ora.S,   /* JOIN_GBN_NM */
		ora.S,   /* SNS_CD */
		ora.S,   /* SNS_CUST_NO */
		ora.S,   /* LST_EDU */
		ora.S,   /* CARR_GBN */
		ora.S,   /* CARR_DESC */
		ora.S,   /* FRGN_LANG_ABLT_DESC */
		ora.S,   /* ATCH_DATA_PATH */
		ora.S,   /* TECH_QLFT_KND */
		ora.S,   /* ATCH_FILE_PATH */
		ora.I64, /* TOT_APPLY_CNT */
		ora.I64, /* STNBY_CNT */
		ora.I64, /* PASS_CNT */
		ora.I64, /* FAIL_CNT */
		ora.F64, /* MATCHING_RATE */
		ora.S,   /* JOBFAIR_MNG_CDS */
		ora.S,   /* JOBFAIR_LIST */
		ora.S,   /* LOGIN_DT */
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

	adminMemberDetail := make([]models.AdminMemberDetail, 0)

	var (
		ppMemNo          string
		memStatCd        string
		memStatNm        string
		memStatDt        string
		memId            string
		nm               string
		sex              string
		email            string
		moNo             string
		brthYmd          string
		age              string
		emailRecvYn      string
		smsRecvYn        string
		ptoPath          string
		regDt            string
		osGbn            string
		osVer            string
		mregPrgsStatCd   string
		mregPrgsStatNm   string
		joinGbnNm        string
		snsCd            string
		snsCustNo        string
		lstEdu           string
		carrGbn          string
		carrDesc         string
		frgnLangAbltDesc string
		atchDataPath     string
		techQlftKnd      string
		atchFilePath     string
		totApplyCnt      int64
		stnbyCnt         int64
		passCnt          int64
		failCnt          int64
		machingRate      float64
		jobfairCds       string
		jobfairList      string
		loginDt          string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			ppMemNo = procRset.Row[0].(string)
			memStatCd = procRset.Row[1].(string)
			memStatNm = procRset.Row[2].(string)
			memStatDt = procRset.Row[3].(string)
			memId = procRset.Row[4].(string)
			nm = procRset.Row[5].(string)
			sex = procRset.Row[6].(string)
			email = procRset.Row[7].(string)
			moNo = procRset.Row[8].(string)
			brthYmd = procRset.Row[9].(string)
			age = procRset.Row[10].(string)
			emailRecvYn = procRset.Row[11].(string)
			smsRecvYn = procRset.Row[12].(string)
			ptoPath = procRset.Row[13].(string)
			regDt = procRset.Row[14].(string)
			osGbn = procRset.Row[15].(string)
			osVer = procRset.Row[16].(string)
			mregPrgsStatCd = procRset.Row[17].(string)
			mregPrgsStatNm = procRset.Row[18].(string)
			joinGbnNm = procRset.Row[19].(string)
			snsCd = procRset.Row[20].(string)
			snsCustNo = procRset.Row[21].(string)
			lstEdu = procRset.Row[22].(string)
			carrGbn = procRset.Row[23].(string)
			carrDesc = procRset.Row[24].(string)
			frgnLangAbltDesc = procRset.Row[25].(string)
			atchDataPath = procRset.Row[26].(string)
			techQlftKnd = procRset.Row[27].(string)
			atchFilePath = procRset.Row[28].(string)
			totApplyCnt = procRset.Row[29].(int64)
			stnbyCnt = procRset.Row[30].(int64)
			passCnt = procRset.Row[31].(int64)
			failCnt = procRset.Row[32].(int64)
			machingRate = procRset.Row[33].(float64)
			jobfairCds = procRset.Row[34].(string)
			jobfairList = procRset.Row[35].(string)
			loginDt = procRset.Row[36].(string)

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

			adminMemberDetail = append(adminMemberDetail, models.AdminMemberDetail{
				PpMemNo:          ppMemNo,
				MemStatCd:        memStatCd,
				MemStatNm:        memStatNm,
				MemStatDt:        memStatDt,
				MemId:            memId,
				Nm:               nm,
				Sex:              sex,
				Email:            email,
				MoNo:             moNo,
				BrthYmd:          brthYmd,
				Age:              age,
				EmailRecvYn:      emailRecvYn,
				SmsRecvYn:        smsRecvYn,
				PtoPath:          fullPtoPath,
				RegDt:            regDt,
				OsGbn:            osGbn,
				OsVer:            osVer,
				MregPrgsStatCd:   mregPrgsStatCd,
				MregPrgsStatNm:   mregPrgsStatNm,
				JoinGbnNm:        joinGbnNm,
				SnsCd:            snsCd,
				SnsCustNo:        snsCustNo,
				LstEdu:           lstEdu,
				CarrGbn:          carrGbn,
				CarrDesc:         carrDesc,
				FrgnLangAbltDesc: frgnLangAbltDesc,
				AtchDataPath:     atchDataPath,
				TechQlftKnd:      techQlftKnd,
				AtchFilePath:     fullFilePath,
				TotApplyCnt:      totApplyCnt,
				StnbyCnt:         stnbyCnt,
				PassCnt:          passCnt,
				FailCnt:          failCnt,
				MachingRate:      machingRate,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Member Info Detail

	// Start : Profile Video List

	log.Debug("CALL SP_EMS_VP_LIST_R('%v', '%v',:1)",
		pLang, pPpMemNo)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_VP_LIST_R('%v', '%v',:1)",
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
		vpVdsn        string
		vpThmKndCd    string
		vpThmNm       string
		vpQstCd       string
		vpQstDesc     string
		vpOpnSetCd    string
		vpOpnSetNm    string
		vpVdFilePath  string
		vpVdThumbPath string
		vpRegDt       string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			vpVdsn = procRset.Row[0].(string)
			vpThmKndCd = procRset.Row[1].(string)
			vpThmNm = procRset.Row[2].(string)
			vpQstCd = procRset.Row[3].(string)
			vpQstDesc = procRset.Row[4].(string)
			vpOpnSetCd = procRset.Row[5].(string)
			vpOpnSetNm = procRset.Row[6].(string)
			vpVdFilePath = procRset.Row[7].(string)
			vpVdThumbPath = procRset.Row[8].(string)
			vpRegDt = procRset.Row[9].(string)

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

			memberVideoProfileList = append(memberVideoProfileList, models.MemberVideoProfileList{
				VpVdsn:        vpVdsn,
				VpThmKndCd:    vpThmKndCd,
				VpThmNm:       vpThmNm,
				VpQstCd:       vpQstCd,
				VpQstDesc:     vpQstDesc,
				VpOpnSetCd:    vpOpnSetCd,
				VpOpnSetNm:    vpOpnSetNm,
				VpVdFilePath:  fullFilePath,
				VpVdThumbPath: fullThumbPath,
				VpRegDt:       vpRegDt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Profile Video List

	// Start : Member Apply history List
	log.Debug("CALL SP_EMS_AM_HIS_LIST_R('%v', '%v', :1)",
		pLang, pPpMemNo)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_AM_HIS_LIST_R('%v', '%v', :1)",
		pLang, pPpMemNo),
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* ENTP_KO_NM */
		ora.S,   /* RECRUT_TITLE */
		ora.S,   /* EVL_PRGS_STAT_NM */
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

	memberApplyHistoryList := make([]models.MemberApplyHistoryList, 0)

	var (
		ahEntpMemNo     string
		ahRecrutSn      string
		ahPpMemNo       string
		ahEntpKoNm      string
		ahRecrutTitle   string
		ahEvlPrgsStatNm string
		ahTotCnt        int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			ahEntpMemNo = procRset.Row[0].(string)
			ahRecrutSn = procRset.Row[1].(string)
			ahPpMemNo = procRset.Row[2].(string)
			ahEntpKoNm = procRset.Row[3].(string)
			ahRecrutTitle = procRset.Row[4].(string)
			ahEvlPrgsStatNm = procRset.Row[5].(string)
			ahTotCnt = procRset.Row[6].(int64)

			memberApplyHistoryList = append(memberApplyHistoryList, models.MemberApplyHistoryList{
				AhEntpMemNo:     ahEntpMemNo,
				AhRecrutSn:      ahRecrutSn,
				AhPpMemNo:       ahPpMemNo,
				AhEntpKoNm:      ahEntpKoNm,
				AhRecrutTitle:   ahRecrutTitle,
				AhEvlPrgsStatNm: ahEvlPrgsStatNm,
				AhTotCnt:        ahTotCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Member Apply history List

	c.Data["PpMemNo"] = ppMemNo
	c.Data["MemStatCd"] = memStatCd
	c.Data["MemStatNm"] = memStatNm
	c.Data["MemStatDt"] = memStatDt
	c.Data["MemId"] = memId
	c.Data["Nm"] = nm
	c.Data["Sex"] = sex
	c.Data["Email"] = email
	c.Data["MoNo"] = moNo
	c.Data["BrthYmd"] = brthYmd
	c.Data["Age"] = age
	c.Data["EmailRecvYn"] = emailRecvYn
	c.Data["SmsRecvYn"] = smsRecvYn
	c.Data["PtoPath"] = fullPtoPath
	c.Data["RegDt"] = regDt
	c.Data["OsGbn"] = osGbn
	c.Data["OsVer"] = osVer
	c.Data["MregPrgsStatCd"] = mregPrgsStatCd
	c.Data["MregPrgsStatNm"] = mregPrgsStatNm
	c.Data["JoinGbnNm"] = joinGbnNm
	c.Data["SnsCd"] = snsCd
	c.Data["SnsCustNo"] = snsCustNo
	c.Data["LstEdu"] = lstEdu
	c.Data["CarrGbn"] = carrGbn
	c.Data["CarrDesc"] = carrDesc
	c.Data["FrgnLangAbltDesc"] = frgnLangAbltDesc
	c.Data["AtchDataPath"] = atchDataPath
	c.Data["TechQlftKnd"] = techQlftKnd
	c.Data["AtchFilePath"] = fullFilePath
	c.Data["TotApplyCnt"] = totApplyCnt
	c.Data["StnbyCnt"] = stnbyCnt
	c.Data["PassCnt"] = passCnt
	c.Data["FailCnt"] = failCnt
	c.Data["MachingRate"] = machingRate

	c.Data["JobFairCdsArr"] = strings.Split(jobfairCds, ",")
	c.Data["JobFairList"] = strings.Split(jobfairList, ",")

	c.Data["LoginDt"] = loginDt

	c.Data["MenuId"] = "05"

	c.Data["MemberVideoProfileList"] = memberVideoProfileList
	c.Data["MemberApplyHistoryList"] = memberApplyHistoryList
	c.Data["AhTotCnt"] = ahTotCnt

	/* Parameter Value */

	c.Data["pMemStat"] = pmMemStat
	c.Data["pKeyword"] = pmKeyword
	c.Data["pSex"] = pmSex
	c.Data["pAge"] = pmAge
	c.Data["pVpYn"] = pmYpYn
	c.Data["pOsGbn"] = pmOsGbn
	c.Data["pPageNo"] = pmPageNo
	c.Data["pSize"] = pmSize
	c.Data["pSdy"] = pmSdy
	c.Data["pEdy"] = pmEdy

	c.Data["pLoginSdy"] = pmLoginSdy
	c.Data["pLoginEdy"] = pmLoginEdy

	c.TplName = "admin/member_detail.html"
}
