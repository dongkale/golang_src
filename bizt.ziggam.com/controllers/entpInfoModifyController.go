package controllers

import (
	"fmt"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"bizt.ziggam.com/models"
	"bizt.ziggam.com/tables"
	ora "gopkg.in/rana/ora.v4"
)

type EntpInfoModifyController struct {
	BaseController
}

func (c *EntpInfoModifyController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	mem_id := session.Get(c.Ctx.Request.Context(), "mem_id")
	if mem_id == nil {
		c.Ctx.Redirect(302, "/login")
	}

	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := mem_no

	imgServer, _  := beego.AppConfig.String("viewpath")
	//cdnPath := beego.AppConfig.String("cdnpath")

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
	/*
		fmt.Printf("CALL ZSP_ENTP_INFO_R('%v', '%v', :1)",
			pLang, pEntpMemNo)

		stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_ENTP_INFO_R('%v', '%v', :1)",
			pLang, pEntpMemNo),
			ora.S,   //* MEM_STAT_CD
			ora.S,   //* LOGO_PTO_PATH
			ora.S,   //* ENTP_KO_NM
			ora.S,   //* BIZ_REG_NO
			ora.S,   //* BIZ_REG_FILE_PATH
			ora.S,   //* REP_NM
			ora.S,   //* TEL_NO
			ora.S,   //* SMS_RECV_YN
			ora.S,   //* EMAIL
			ora.S,   //* EMAIL_RECV_YN
			ora.S,   //* ZIP
			ora.S,   //* ADDR
			ora.S,   //* DTL_ADDR
			ora.S,   //* REF_ADDR
			ora.S,   //* BIZ_TPY
			ora.S,   //* BIZ_COND
			ora.S,   //* ENTP_INTR
			ora.S,   //* ENTP_HTAG
			ora.S,   //* HOME_PG
			ora.I64, //* EMP_CNT
			ora.S,   //* EST_DY
			ora.S,   //* INFO_EQ_YN
		)
	*/

	fmt.Printf(fmt.Sprintf("CALL ZSP_ENTP_INFO_R_V2('%v', '%v', :1)",
		pLang, pEntpMemNo))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_ENTP_INFO_R_V2('%v', '%v', :1)",
		pLang, pEntpMemNo),
		ora.S,   /* MEM_STAT_CD */
		ora.S,   /* LOGO_PTO_PATH */
		ora.S,   /* ENTP_KO_NM */
		ora.S,   /* BIZ_REG_NO */
		ora.S,   /* BIZ_REG_FILE_PATH */
		ora.S,   /* REP_NM */
		ora.S,   /* TEL_NO */
		ora.S,   /* SMS_RECV_YN */
		ora.S,   /* EMAIL */
		ora.S,   /* EMAIL_RECV_YN */
		ora.S,   /* ZIP */
		ora.S,   /* ADDR */
		ora.S,   /* DTL_ADDR */
		ora.S,   /* REF_ADDR */
		ora.S,   /* BIZ_TPY */
		ora.S,   /* BIZ_COND */
		ora.S,   /* ENTP_INTR */
		ora.S,   /* ENTP_HTAG */
		ora.S,   /* HOME_PG */
		ora.I64, /* EMP_CNT */
		ora.S,   /* EST_DY */
		ora.S,   /* INFO_EQ_YN */
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
		memStatCd      string
		logoPtoPath    string
		entpKoNm       string
		bizRegNo       string
		bizRegFilePath string
		repNm          string
		telNo          string
		smsRecvYn      string
		email          string
		emailRecvYn    string
		zip            string
		addr           string
		dtlAddr        string
		refAddr        string
		bizTpy         string
		bizCond        string
		entpIntr       string
		entpHtag       string
		homePg         string
		empCnt         int64
		estDy          string
		infoEqYn       string

		bizTpyCd       string
		entpProfile    string
		bizIntro       string
		entpCapital    string
		entpTotalSales string
		entpTypeCd     string
		location       string

		fullPtoPath  string
		fullFilePath string

		entpHtag1 string
		entpHtag2 string
		entpHtag3 string
	)

	if procRset.IsOpen() {
		for procRset.Next() {

			memStatCd = procRset.Row[0].(string)
			logoPtoPath = procRset.Row[1].(string)
			if logoPtoPath == "" {
				fullPtoPath = logoPtoPath
			} else {
				fullPtoPath = imgServer + logoPtoPath
			}
			entpKoNm = procRset.Row[2].(string)
			bizRegNo = procRset.Row[3].(string)
			bizRegFilePath = procRset.Row[4].(string)
			if bizRegFilePath == "" {
				fullFilePath = bizRegFilePath
			} else {
				fullFilePath = imgServer + bizRegFilePath
			}
			repNm = procRset.Row[5].(string)
			telNo = procRset.Row[6].(string)
			smsRecvYn = procRset.Row[7].(string)
			email = procRset.Row[8].(string)
			emailRecvYn = procRset.Row[9].(string)
			zip = procRset.Row[10].(string)
			addr = procRset.Row[11].(string)
			dtlAddr = procRset.Row[12].(string)
			refAddr = procRset.Row[13].(string)
			bizTpy = procRset.Row[14].(string)
			bizCond = procRset.Row[15].(string)
			entpIntr = procRset.Row[16].(string)
			entpHtag = procRset.Row[17].(string)
			homePg = procRset.Row[18].(string)
			empCnt = procRset.Row[19].(int64)
			estDy = procRset.Row[20].(string)
			infoEqYn = procRset.Row[21].(string)

			bizTpyCd = procRset.Row[22].(string)
			entpProfile = procRset.Row[23].(string)
			bizIntro = procRset.Row[24].(string)
			entpCapital = procRset.Row[25].(string)
			entpTotalSales = procRset.Row[26].(string)
			entpTypeCd = procRset.Row[27].(string)
			location = procRset.Row[28].(string)

			entpHtagArr := strings.Split(entpHtag, ",")
			/*
				for i := range entpHtagArr {
					fmt.Println(entpHtagArr[i])
				}
				fmt.Println(len(entpHtagArr))
			*/

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

			entpInfo = append(entpInfo, models.EntpInfo{
				MemStatCd:      memStatCd,
				LogoPtoPath:    fullPtoPath,
				EntpKoNm:       entpKoNm,
				BizRegNo:       bizRegNo,
				BizRegFilePath: fullFilePath,
				RepNm:          repNm,
				TelNo:          telNo,
				SmsRecvYn:      smsRecvYn,
				Email:          email,
				EmailRecvYn:    emailRecvYn,
				Zip:            zip,
				Addr:           addr,
				DtlAddr:        dtlAddr,
				RefAddr:        refAddr,
				BizTpy:         bizTpy,
				BizCond:        bizCond,
				EntpIntr:       entpIntr,
				EntpHtag:       entpHtag,
				EntpHtag1:      entpHtag1,
				EntpHtag2:      entpHtag2,
				EntpHtag3:      entpHtag3,
				HomePg:         homePg,
				EmpCnt:         empCnt,
				EstDy:          estDy,
				InfoEqYn:       infoEqYn,
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

	// Start : Jobfair List
	fmt.Printf(fmt.Sprintf("CALL ZSP_ENTP_JOBFAIR_LIST('%v', '%v', :1)", pLang, pEntpMemNo))
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

	jobFailrInfoList := make([]models.JobfairInfo, 0)

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

			jobFailrInfoList = append(jobFailrInfoList, models.JobfairInfo{
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

	c.Data["MemStatCd"] = memStatCd
	c.Data["LogoPtoPath"] = fullPtoPath
	c.Data["EntpKoNm"] = entpKoNm
	c.Data["BizRegNo"] = bizRegNo
	c.Data["BizRegFilePath"] = fullFilePath
	c.Data["RepNm"] = repNm
	c.Data["TelNo"] = telNo
	c.Data["SmsRecvYn"] = smsRecvYn
	c.Data["Email"] = email
	c.Data["EmailRecvYn"] = emailRecvYn
	c.Data["Zip"] = zip
	c.Data["Addr"] = addr
	c.Data["DtlAddr"] = dtlAddr
	c.Data["RefAddr"] = refAddr
	c.Data["BizTpy"] = bizTpy
	c.Data["BizCond"] = bizCond
	c.Data["EntpIntr"] = entpIntr
	c.Data["EntpHtag"] = entpHtag
	c.Data["EntpHtag1"] = entpHtag1
	c.Data["EntpHtag2"] = entpHtag2
	c.Data["EntpHtag3"] = entpHtag3
	c.Data["HomePg"] = homePg
	c.Data["EmpCnt"] = empCnt
	c.Data["EstDy"] = estDy
	c.Data["InfoEqYn"] = infoEqYn
	c.Data["OriLogoFile"] = logoPtoPath
	c.Data["OriBizFile"] = bizRegFilePath

	c.Data["BizTpyCd"] = bizTpyCd
	c.Data["EntpProfile"] = entpProfile
	c.Data["BizIntro"] = bizIntro
	c.Data["EntpCapital"] = entpCapital
	c.Data["EntpTotalSales"] = entpTotalSales
	c.Data["EntpTypeCd"] = entpTypeCd
	c.Data["Location"] = location

	c.Data["EntpMemNo"] = pEntpMemNo

	// LDK 2020/08/26 : 기업 정보 코드화, 추가 -->
	c.Data["MapEntpTypeCd"] = tables.MapEntpTypeCd
	c.Data["MapBizTpyCd"] = tables.MapBizTpyCd
	// <--

	// LDK 2020/10/17 : 박람회 등록 리스트 -->
	c.Data["JobFailrList"] = jobFailrInfoList
	// <--

	c.Data["TMenuId"] = "E00"
	c.Data["SMenuId"] = "E01"
	c.TplName = "entp/entp_info.html"
}
