package controllers

import (
	"fmt"
	"strings"

	"emst.ziggam.com/models"
	"emst.ziggam.com/tables"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	ora "gopkg.in/rana/ora.v4"
)

type AdminEntpInfoController struct {
	BaseController
}

func (c *AdminEntpInfoController) Get() {

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
	pEntpMemNo := c.GetString("entp_mem_no")

	// Parameter
	pPageNo := c.GetString("pn")
	pSize := c.GetString("size")
	pGbnCd := c.GetString("gbn_cd")
	pVdYn := c.GetString("vd_yn")
	pOsGbn := c.GetString("os_gbn")
	pUseYn := c.GetString("use_yn")
	pSdy := c.GetString("sdy")
	pEdy := c.GetString("edy")
	pKeyword := c.GetString("keyword")

	pLoginSdy := c.GetString("login_sdy")
	pLoginEdy := c.GetString("login_edy")

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

	// Start : Entp Info
	logs.Debug(fmt.Sprintf("CALL SP_EMS_ENTP_INFO_R_V2('%v', '%v', :1)",
		pLang, pEntpMemNo))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ENTP_INFO_R_V2('%v', '%v', :1)",
		pLang, pEntpMemNo),
		ora.S,   /* ENTP_KO_NM */
		ora.S,   /* BIZ_REG_NO */
		ora.S,   /* LOGO_PTO_PATH */
		ora.S,   /* REP_NM */
		ora.S,   /* TEL_NO */
		ora.S,   /* EST_DY */
		ora.I64, /* EMP_CNT */
		ora.S,   /* BIZ_TPY */
		ora.S,   /* BIZ_COND */
		ora.S,   /* PP_CHRG_NM */
		ora.S,   /* PP_CHRG_TEL_NO */
		ora.S,   /* ZIP */
		ora.S,   /* ADDR */
		ora.S,   /* DTL_ADDR */
		ora.S,   /* REF_ADDR */
		ora.S,   /* ENTP_HTAG */
		ora.S,   /* ENTP_INTR */
		ora.S,   /* HOME_PG */
		ora.S,   /* VD_TITLE1 */
		ora.S,   /* VD_TITLE2 */
		ora.S,   /* VD_TITLE3 */
		ora.S,   /* VD_TITLE4 */
		ora.S,   /* VD_FILE_PATH1 */
		ora.S,   /* VD_FILE_PATH2 */
		ora.S,   /* VD_FILE_PATH3 */
		ora.S,   /* VD_FILE_PATH4 */
		ora.S,   /* VIDEO_YN */
		ora.I64, /* VIDEO_CNT */
		ora.S,   /* EMAIL */
		ora.S,   /* MEM_STAT_CD */
		ora.S,   /* USE_YN */
		ora.S,   /* BIZ_REG_FILE_PATH */
		ora.S,   /* LOGIN_DT */
		ora.S,   /* BIZ_TPY_CD */
		ora.S,   /* ENTP_PROFILE */
		ora.S,   /* BIZ_INTRO */
		ora.S,   /* ENTP_CAPITAL */
		ora.S,   /* ENTP_TOTAL_SALES */
		ora.S,   /* ENTP_TYPE_CD */
		ora.S,   /* LOCATION */
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

	entpInfo := make([]models.EntpInfo, 0)

	var (
		entpKoNm    string
		bizRegNo    string
		logoPtoPath string
		repNm       string
		telNo          string
		estDy       string
		empCnt      int64
		bizTpy      string
		bizCond        string
		ppChrgNm    string
		ppChrgTelNo string
		zip         string
		addr        string
		dtlAddr     string
		refAddr     string
		entpHtag    string
		entpIntr    string
		homePg      string
		vdTitle1    string
		vdTitle2    string
		vdTitle3    string
		vdTitle4    string
		vdFilePath1 string
		vdFilePath2 string
		vdFilePath3 string
		vdFilePath4 string
		videoYn     string
		videoCnt    int64
		memStatCd   string
		email       string
		useYn       string
		bizRegFilePath string
		jobfairCds     string
		loginDt        string

		fullPtoPath     string
		fullVdFilePath1 string
		fullVdFilePath2 string
		fullVdFilePath3 string
		fullVdFilePath4 string

		entpHtag1 string
		entpHtag2 string
		entpHtag3 string

		bizTpyCd       string
		entpProfile    string
		bizIntro       string
		entpCapital    string
		entpTotalSales string
		entpTypeCd     string
		location       string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			entpKoNm = procRset.Row[0].(string)
			bizRegNo = procRset.Row[1].(string)
			logoPtoPath = procRset.Row[2].(string)

			if logoPtoPath == "" {
				fullPtoPath = logoPtoPath
			} else {
				fullPtoPath = imgServer + logoPtoPath
			}

			repNm = procRset.Row[3].(string)
			telNo = procRset.Row[4].(string)
			estDy = procRset.Row[5].(string)
			empCnt = procRset.Row[6].(int64)
			bizTpy = procRset.Row[7].(string)
			bizCond = procRset.Row[8].(string)
			ppChrgNm = procRset.Row[9].(string)
			ppChrgTelNo = procRset.Row[10].(string)
			zip = procRset.Row[11].(string)
			addr = procRset.Row[12].(string)
			dtlAddr = procRset.Row[13].(string)
			refAddr = procRset.Row[14].(string)
			entpHtag = procRset.Row[15].(string)
			entpIntr = procRset.Row[16].(string)
			homePg = procRset.Row[17].(string)
			vdTitle1 = procRset.Row[18].(string)
			vdTitle2 = procRset.Row[19].(string)
			vdTitle3 = procRset.Row[20].(string)
			vdTitle4 = procRset.Row[21].(string)
			vdFilePath1 = procRset.Row[22].(string)
			vdFilePath2 = procRset.Row[23].(string)
			vdFilePath3 = procRset.Row[24].(string)
			vdFilePath4 = procRset.Row[25].(string)

			entpHtagArr := strings.Split(entpHtag, ",")

			for i := range entpHtagArr {
				fmt.Println(entpHtagArr[i])
			}
			fmt.Println(len(entpHtagArr))

			if len(entpHtagArr) == 1 {
				entpHtag1 = entpHtagArr[0]
				entpHtag2 = ""
				entpHtag3 = ""
			}
			if len(entpHtagArr) == 2 {
				entpHtag1 = entpHtagArr[0]
				entpHtag2 = entpHtagArr[1]
				entpHtag3 = ""
			}
			if len(entpHtagArr) == 3 {
				entpHtag1 = entpHtagArr[0]
				entpHtag2 = entpHtagArr[1]
				entpHtag3 = entpHtagArr[2]
			}

			if vdFilePath1 == "" {
				fullVdFilePath1 = vdFilePath1
			} else {
				fullVdFilePath1 = cdnPath + vdFilePath1
			}

			if vdFilePath2 == "" {
				fullVdFilePath2 = vdFilePath2
			} else {
				fullVdFilePath2 = cdnPath + vdFilePath2
			}

			if vdFilePath3 == "" {
				fullVdFilePath3 = vdFilePath3
			} else {
				fullVdFilePath3 = cdnPath + vdFilePath3
			}

			if vdFilePath4 == "" {
				fullVdFilePath4 = vdFilePath4
			} else {
				fullVdFilePath4 = cdnPath + vdFilePath4
			}

			videoYn = procRset.Row[26].(string)
			videoCnt = procRset.Row[27].(int64)
			email = procRset.Row[28].(string)
			memStatCd = procRset.Row[29].(string)
			useYn = procRset.Row[30].(string)
			bizRegFilePath = procRset.Row[31].(string)
			jobfairCds = procRset.Row[32].(string)
			loginDt = procRset.Row[33].(string)

			bizTpyCd = procRset.Row[34].(string)
			entpProfile = procRset.Row[35].(string)
			bizIntro = procRset.Row[36].(string)
			entpCapital = procRset.Row[37].(string)
			entpTotalSales = procRset.Row[38].(string)
			entpTypeCd = procRset.Row[39].(string)
			location = procRset.Row[40].(string)

			entpInfo = append(entpInfo, models.EntpInfo{
				EntpKoNm:    entpKoNm,
				BizRegNo:    bizRegNo,
				LogoPtoPath: fullPtoPath,
				RepNm:       repNm,
				TelNo:          telNo,
				EstDy:       estDy,
				EmpCnt:      empCnt,
				BizTpy:      bizTpy,
				BizCond:     bizCond,
				PpChrgNm:    ppChrgNm,
				PpChrgTelNo: ppChrgTelNo,
				Zip:         zip,
				Addr:        addr,
				DtlAddr:     dtlAddr,
				RefAddr:     refAddr,
				EntpHtag:    entpHtag,
				EntpHtag1:   entpHtag1,
				EntpHtag2:   entpHtag2,
				EntpHtag3:   entpHtag3,
				EntpIntr:    entpIntr,
				HomePg:      homePg,
				VdTitle1:    vdTitle1,
				VdTitle2:    vdTitle2,
				VdTitle3:    vdTitle3,
				VdTitle4:    vdTitle4,
				VdFilePath1: fullVdFilePath1,
				VdFilePath2: fullVdFilePath2,
				VdFilePath3: fullVdFilePath3,
				VdFilePath4: fullVdFilePath4,
				VideoYn:     videoYn,
				VideoCnt:    videoCnt,
				Email:       email,
				MemStatCd:   memStatCd,
				UseYn:       useYn,
				OriLogoFile: bizRegFilePath,
				LoginDt:     loginDt,
				BizTpyCd:       bizTpyCd,
				EntpProfile:    entpProfile,
				BizIntro:       bizIntro,
				EntpCapital:    entpCapital,
				EntpTotalSales: entpTotalSales,
				EntpTypeCd:     entpTypeCd,
				Location:       location,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Entp Info

	// Start : Entp Team Member List
	pMemberGbnCd := "A"

	logs.Debug(fmt.Sprintf("CALL ZSP_TEAM_MEM_LIST_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pMemberGbnCd))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_TEAM_MEM_LIST_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pMemberGbnCd),
		ora.I64, /* TOT_CNT */
		ora.S,   /* PP_CHRG_SN */
		ora.S,   /* PP_CHRG_GBN_CD */
		ora.S,   /* PP_CHRG_NM */
		ora.S,   /* PP_CHRG_BP_NM */
		ora.S,   /* EMAIL */
		ora.S,   /* ENTP_MEM_ID */
		ora.S,   /* PP_CHRG_TEL_NO */
		ora.I64, /* ROWNO */
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
		etRowNo       int64
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
			etRowNo = procRset.Row[8].(int64)

			entpTeamMemberList = append(entpTeamMemberList, models.EntpTeamMemberList{
				EtTotCnt:      etTotCnt,
				EtPpChrgSn:    etPpChrgSn,
				EtPpChrgGbnCd: etPpChrgGbnCd,
				EtPpChrgNm:    etPpChrgNm,
				EtPpChrgBpNm:  etPpChrgBpNm,
				EtEmail:       etEmail,
				EtEntpMemId:   etEntpMemId,
				EtPpChrgTelNo: etPpChrgTelNo,
				EtRowNo:       etRowNo,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Entp Team Member List

	// Start : Jobfair List
	logs.Debug(fmt.Sprintf("CALL ZSP_ENTP_JOBFAIR_LIST('%v', '%v', :1)", pLang, pEntpMemNo))
	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_ENTP_JOBFAIR_LIST('%v', '%v', :1)",
		pLang, pEntpMemNo),
		ora.S, /* MNG_CD */
		ora.S, /* TITLE */
		ora.S, /* SDY */
		ora.S, /* EDY */
		ora.S, /* HOST_INSTITUTION */
		ora.S, /* MANAGE_AGENCY */
		ora.S, /* AGREEMENT_URL */
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

	jobFairInfoList := make([]models.JobfairInfo, 0)

	var (
		jfMngCd           string
		jfTitle           string
		jfSdy             string
		jfEdy             string
		jfHostInstitution string
		jfManageAgency    string
	)

	if procRset.IsOpen() {
		for procRset.Next() {

			jfMngCd = procRset.Row[0].(string)
			jfTitle = procRset.Row[1].(string)
			jfSdy = procRset.Row[2].(string)
			jfEdy = procRset.Row[3].(string)
			jfHostInstitution = procRset.Row[4].(string)
			jfManageAgency = procRset.Row[5].(string)

			jobFairInfoList = append(jobFairInfoList, models.JobfairInfo{
				MngCd:           jfMngCd,
				Title:           jfTitle,
				Sdy:             jfSdy,
				Edy:             jfEdy,
				HostInstitution: jfHostInstitution,
				ManageAgency:    jfManageAgency,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Jobfair List

	// Start : Jobfair List
	logs.Debug(fmt.Sprintf("CALL ZSP_ENTP_JOBFAIR_LIST('%v', '%v', :1)", pLang, ""))
	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_ENTP_JOBFAIR_LIST('%v', '%v', :1)",
		pLang, ""),
		ora.S, /* MNG_CD */
		ora.S, /* TITLE */
		ora.S, /* SDY */
		ora.S, /* EDY */
		ora.S, /* HOST_INSTITUTION */
		ora.S, /* MANAGE_AGENCY */
		ora.S, /* AGREEMENT_URL */
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

	jobFairNotRegList := make([]models.JobfairInfo, 0)

	if procRset.IsOpen() {
		for procRset.Next() {

			jfMngCd = procRset.Row[0].(string)
			jfTitle = procRset.Row[1].(string)
			jfSdy = procRset.Row[2].(string)
			jfEdy = procRset.Row[3].(string)
			jfHostInstitution = procRset.Row[4].(string)
			jfManageAgency = procRset.Row[5].(string)

			jobFairNotRegList = append(jobFairNotRegList, models.JobfairInfo{
				MngCd:           jfMngCd,
				Title:           jfTitle,
				Sdy:             jfSdy,
				Edy:             jfEdy,
				HostInstitution: jfHostInstitution,
				ManageAgency:    jfManageAgency,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Jobfair not reg List

	c.Data["EntpKoNm"] = entpKoNm
	c.Data["EntpMemNo"] = pEntpMemNo
	c.Data["BizRegNo"] = bizRegNo
	c.Data["LogoPtoPath"] = fullPtoPath
	c.Data["RepNm"] = repNm
	c.Data["TelNo"] = telNo
	c.Data["EstDy"] = estDy
	c.Data["EmpCnt"] = empCnt
	c.Data["BizTpy"] = bizTpy
	c.Data["BizCond"] = bizCond
	c.Data["PpChrgNm"] = ppChrgNm
	c.Data["PpChrgTelNo"] = ppChrgTelNo
	c.Data["Zip"] = zip
	c.Data["Addr"] = addr
	c.Data["DtlAddr"] = dtlAddr
	c.Data["RefAddr"] = refAddr
	c.Data["EntpHtag"] = entpHtag
	c.Data["EntpHtag1"] = entpHtag1
	c.Data["EntpHtag2"] = entpHtag2
	c.Data["EntpHtag3"] = entpHtag3
	c.Data["EntpIntr"] = entpIntr
	c.Data["HomePg"] = homePg
	c.Data["VdTitle1"] = vdTitle1
	c.Data["VdTitle2"] = vdTitle2
	c.Data["VdTitle3"] = vdTitle3
	c.Data["VdTitle4"] = vdTitle4
	c.Data["VdFilePath1"] = fullVdFilePath1
	c.Data["VdFilePath2"] = fullVdFilePath2
	c.Data["VdFilePath3"] = fullVdFilePath3
	c.Data["VdFilePath4"] = fullVdFilePath4
	c.Data["VideoYn"] = videoYn
	c.Data["VideoCnt"] = videoCnt
	c.Data["Email"] = email
	c.Data["MemStatCd"] = memStatCd
	c.Data["UseYn"] = useYn
	c.Data["BizRegFilePath"] = bizRegFilePath
	c.Data["JobFairCdsArr"] = strings.Split(jobfairCds, ",")
	c.Data["LoginDt"] = loginDt

	c.Data["OriLogoFile"] = logoPtoPath

	c.Data["BizTpyCd"] = bizTpyCd
	c.Data["EntpProfile"] = entpProfile
	c.Data["BizIntro"] = bizIntro
	c.Data["EntpCapital"] = entpCapital
	c.Data["EntpTotalSales"] = entpTotalSales
	c.Data["EntpTypeCd"] = entpTypeCd
	c.Data["Location"] = location

	c.Data["MapEntpTypeCd"] = tables.MapEntpTypeCd
	c.Data["MapBizTpyCd"] = tables.MapBizTpyCd

	c.Data["TeamMemberList"] = entpTeamMemberList

	c.Data["JobFairList"] = jobFairInfoList
	c.Data["JobFairNotRegList"] = jobFairNotRegList

	c.Data["MenuId"] = "06"

	/* Parameter Value */
	c.Data["pPageNo"] = pPageNo
	c.Data["pSize"] = pSize
	c.Data["pGbnCd"] = pGbnCd
	c.Data["pVdYn"] = pVdYn
	c.Data["pOsGbn"] = pOsGbn
	c.Data["pUseYn"] = pUseYn
	c.Data["pSdy"] = pSdy
	c.Data["pEdy"] = pEdy
	c.Data["pKeyword"] = pKeyword

	c.Data["pLoginSdy"] = pLoginSdy
	c.Data["pLoginEdy"] = pLoginEdy

	logs.Debug(fmt.Sprintf("[EntpInfo][Last][%v] EntpKoNm:%v, EntpMemNo:%v, BizRegNo:%v, LogoPtoPath:%v, RepNm:%v, TelNo:%v, EstDy:%v, EmpCnt:%v, BizTpy:%v, BizCond:%v, PpChrgNm:%v, PpChrgTelNo:%v, Zip:%v, Addr:%v, DtlAddr:%v, RefAddr:%v, EntpHtag:%v, EntpHtag1:%v, EntpHtag2:%v, EntpHtag3:%v, EntpIntr:%v, HomePg:%v, VdTitle1:%v, VdTitle2:%v, VdTitle3:%v, VdTitle4:%v, VdFilePath1:%v, VdFilePath2:%v, VdFilePath3:%v, VdFilePath4:%v, VideoYn:%v, VideoCnt:%v, Email:%v, MemStatCd:%v, UseYn:%v, BizRegFilePath:%v, JobFairCdsArr:%v, LoginDt:%v, OriLogoFile:%v, BizTpyCd:%v, EntpProfile:%v, BizIntro:%v, EntpCapital:%v, EntpTotalSales:%v, EntpTypeCd:%v, Location:%v",
		pEntpMemNo, entpKoNm, pEntpMemNo, bizRegNo, fullPtoPath, repNm, telNo, estDy, empCnt, bizTpy, bizCond, ppChrgNm, ppChrgTelNo, zip, addr, dtlAddr, refAddr, entpHtag, entpHtag1, entpHtag2, entpHtag3, entpIntr, homePg, vdTitle1, vdTitle2, vdTitle3, vdTitle4, fullVdFilePath1, fullVdFilePath2, fullVdFilePath3, fullVdFilePath4, videoYn, videoCnt, email, memStatCd, useYn, bizRegFilePath, strings.Split(jobfairCds, ","), loginDt, logoPtoPath, bizTpyCd, entpProfile, bizIntro, entpCapital, entpTotalSales, entpTypeCd, location))

	c.TplName = "admin/entp_info.html"
}
