package controllers

import (
	"fmt"
	"strings"

	"emst.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type EntpInfoWriteController struct {
	BaseController
}

func (c *EntpInfoWriteController) Get() {

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
	log.Debug("CALL SP_EMS_ENTP_INFO_R('%v', '%v', :1)",
		pLang, pEntpMemNo)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ENTP_INFO_R('%v', '%v', :1)",
		pLang, pEntpMemNo),
		ora.S,   /* ENTP_KO_NM */
		ora.S,   /* BIZ_REG_NO */
		ora.S,   /* LOGO_PTO_PATH */
		ora.S,   /* REP_NM */
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
		ora.S,   /* MEM_STAT_CD */
		ora.S,   /* EMAIL */
		ora.S,   /* USE_YN */
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
		estDy       string
		empCnt      int64
		bizTpy      string
		bizCond     string
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
		email       string
		useYn       string

		fullPtoPath     string
		fullVdFilePath1 string
		fullVdFilePath2 string
		fullVdFilePath3 string
		fullVdFilePath4 string

		entpHtag1 string
		entpHtag2 string
		entpHtag3 string
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
			estDy = procRset.Row[4].(string)
			empCnt = procRset.Row[5].(int64)
			bizTpy = procRset.Row[6].(string)
			bizCond = procRset.Row[7].(string)
			ppChrgNm = procRset.Row[8].(string)
			ppChrgTelNo = procRset.Row[9].(string)
			zip = procRset.Row[10].(string)
			addr = procRset.Row[11].(string)
			dtlAddr = procRset.Row[12].(string)
			refAddr = procRset.Row[13].(string)
			entpHtag = procRset.Row[14].(string)
			entpIntr = procRset.Row[15].(string)
			homePg = procRset.Row[16].(string)
			vdTitle1 = procRset.Row[17].(string)
			vdTitle2 = procRset.Row[18].(string)
			vdTitle3 = procRset.Row[19].(string)
			vdTitle4 = procRset.Row[20].(string)
			vdFilePath1 = procRset.Row[21].(string)
			vdFilePath2 = procRset.Row[22].(string)
			vdFilePath3 = procRset.Row[23].(string)
			vdFilePath4 = procRset.Row[24].(string)

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

			videoYn = procRset.Row[25].(string)
			videoCnt = procRset.Row[26].(int64)
			email = procRset.Row[27].(string)
			useYn = procRset.Row[29].(string)

			entpInfo = append(entpInfo, models.EntpInfo{
				EntpKoNm:    entpKoNm,
				BizRegNo:    bizRegNo,
				LogoPtoPath: fullPtoPath,
				RepNm:       repNm,
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
				OriLogoFile: logoPtoPath,
				Email:       email,
				UseYn:       useYn,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Entp Info

	c.Data["EntpKoNm"] = entpKoNm
	c.Data["EntpMemNo"] = pEntpMemNo
	c.Data["BizRegNo"] = bizRegNo
	c.Data["LogoPtoPath"] = fullPtoPath
	c.Data["RepNm"] = repNm
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
	c.Data["OriLogoFile"] = logoPtoPath
	c.Data["Email"] = email
	c.Data["UseYn"] = useYn
	c.Data["MenuId"] = "01"

	c.TplName = "entp/entp_info_write.html"
}
