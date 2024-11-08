package controllers

import (
	"fmt"
	"os"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/tealeg/xlsx"

	"emst.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

// AdminEntpListExcelController ...
type AdminEntpListExcelController struct {
	beego.Controller
}

// Post ...
func (c *AdminEntpListExcelController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Data["json"] = &models.RtnAdminEntpListExcel{}
		c.ServeJSON()
		return
	}

	pLang, _ := beego.AppConfig.String("lang")
	imgServer, _ := beego.AppConfig.String("viewpath")

	//pEntpMemNo := c.GetString("entp_mem_no") // 기업회원번호(세션)
	pEntpMemNo := mem_no.(string)

	pGbnCd := c.GetString("gbn_cd")
	if pGbnCd == "" {
		pGbnCd = "A"
	}
	pKeyword := c.GetString("keyword") // 검색어
	pSdy := c.GetString("sdy")         // 검색시작일자
	pEdy := c.GetString("edy")         // 검색종료일자
	pVdYn := c.GetString("vd_yn")      //영상프로필 (전체:A, 있음:Y, 없음:N)
	pUseYn := c.GetString("use_yn")    //영상검증여부 (전체:A, 완료:1, 대기:0)
	pOsGbn := c.GetString("os_gbn")    //유입경로 (전체:A, 웹:WB, 안드로이드:AD, 아이폰:IS)

	pJobFair := c.GetString("jf_mng_cd")

	pLoginSdy := c.GetString("login_sdy") // 로그인 검색시작일자
	pLoginEdy := c.GetString("login_edy") // 로그인 검색종료일자

	if pSdy == "" {
		pSdy = models.DefaultSdy
	}

	if pEdy == "" {
		nowTime := time.Now()
		pEdy = nowTime.Format("20060102")
	}

	if pLoginSdy == "" {
		pLoginSdy = models.DefaultLoginSdy
	}

	if pLoginEdy == "" {
		nowTime := time.Now()
		pLoginEdy = nowTime.Format("20060102")
	}

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// -- 가입일시(미검증기업/탈퇴기업/), 기업관리코드, 참여 박람회, 	회사명, 사업자번호,	대표자명, 아이디, 담당자이름, 담당자연락처, 총지원수(신규), 영상, 유입, 검증(대기/완료/-), 최근 접속 일시

	// Start : Recruit Apply Entp Excel Download
	logs.Debug(fmt.Sprintf("CALL SP_EMS_ADMIN_ENTP_LST_EXCEL('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pGbnCd, pSdy, pEdy, pVdYn, pUseYn, pKeyword, pOsGbn, pJobFair, pLoginSdy, pLoginEdy))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_ENTP_LST_EXCEL('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pGbnCd, pSdy, pEdy, pVdYn, pUseYn, pKeyword, pOsGbn, pJobFair, pLoginSdy, pLoginEdy),
		ora.S,   /* APPLY_DT */
		ora.S,   /* APPLY_STAT */
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* JOBFAIR_CODE */
		ora.S,   /* ENTP_KO_NM */
		ora.S,   /* BIZ_REG_NO */
		ora.S,   /* REP_NM */
		ora.S,   /* ENTP_MEM_ID */
		ora.S,   /* PP_CHRG_NM */
		ora.S,   /* PP_CHRG_TEL_NO */
		ora.I64, /* TOT_APLY_CNT */
		ora.I64, /* NEW_APLYCNT */
		ora.S,   /* VD_YN */
		ora.S,   /* OS_GBN */
		ora.S,   /* VERIFY_STAT */
		ora.S,   /* LAST_LOGIN */
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

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}

	nowDate := time.Now()
	dateFmt := fmt.Sprintf(nowDate.Format("20060102150405"))

	rtnAdminEntpListExcel := models.RtnAdminEntpListExcel{}
	adminEntpListExcel := make([]models.AdminEntpListExcel, 0)

	var (
		applyDt      string
		applyStat    string
		entpMemNo    string
		jobFairCode  string
		entpKoNm     string
		bizRegNo     string
		repNm        string
		entpMemId    string
		ppChrgNm     string
		ppChrgTelNo  string
		totAplyCnt   int64
		newAplyCnt   int64
		vpYn         string
		osGbn        string
		verifyStat   string
		lastLoginDt  string
		downloadPath string
	)

	row = sheet.AddRow()

	cell = row.AddCell()
	cell.Value = "가입일시"
	cell = row.AddCell()
	cell.Value = "가입상태"
	cell = row.AddCell()
	cell.Value = "기업관리코드"
	cell = row.AddCell()
	cell.Value = "참여 박람회"
	cell = row.AddCell()
	cell.Value = "회사명"
	cell = row.AddCell()
	cell.Value = "사업자번호"
	cell = row.AddCell()
	cell.Value = "대표자명"
	cell = row.AddCell()
	cell.Value = "아이디"
	cell = row.AddCell()
	cell.Value = "담당자이름"
	cell = row.AddCell()
	cell.Value = "담당자연락처"
	cell = row.AddCell()
	cell.Value = "총지원수"
	cell = row.AddCell()
	cell.Value = "신규지원수"
	cell = row.AddCell()
	cell.Value = "영상"
	cell = row.AddCell()
	cell.Value = "유입"
	cell = row.AddCell()
	cell.Value = "검증"
	cell = row.AddCell()
	cell.Value = "최근 접속 일시"
	cell = row.AddCell()

	if procRset.IsOpen() {
		for procRset.Next() {
			applyDt = procRset.Row[0].(string)
			applyStat = procRset.Row[1].(string)
			entpMemNo = procRset.Row[2].(string)
			jobFairCode = procRset.Row[3].(string)
			entpKoNm = procRset.Row[4].(string)
			bizRegNo = procRset.Row[5].(string)
			repNm = procRset.Row[6].(string)
			entpMemId = procRset.Row[7].(string)
			ppChrgNm = procRset.Row[8].(string)
			ppChrgTelNo = procRset.Row[9].(string)
			totAplyCnt = procRset.Row[10].(int64)
			newAplyCnt = procRset.Row[11].(int64)
			vpYn = procRset.Row[12].(string)
			osGbn = procRset.Row[13].(string)
			verifyStat = procRset.Row[14].(string)
			lastLoginDt = procRset.Row[15].(string)

			row = sheet.AddRow()

			cell = row.AddCell()
			cell.Value = applyDt
			cell = row.AddCell()
			cell.Value = applyStat
			cell = row.AddCell()
			cell.Value = entpMemNo
			cell = row.AddCell()
			cell.Value = jobFairCode
			cell = row.AddCell()
			cell.Value = entpKoNm
			cell = row.AddCell()
			cell.Value = bizRegNo
			cell = row.AddCell()
			cell.Value = repNm
			cell = row.AddCell()
			cell.Value = entpMemId
			cell = row.AddCell()
			cell.Value = ppChrgNm
			cell = row.AddCell()
			cell.Value = ppChrgTelNo
			cell = row.AddCell()
			cell.Value = fmt.Sprintf("%v", totAplyCnt)
			cell = row.AddCell()
			cell.Value = fmt.Sprintf("%v", newAplyCnt)
			cell = row.AddCell()
			cell.Value = vpYn
			cell = row.AddCell()
			cell.Value = osGbn
			cell = row.AddCell()
			cell.Value = verifyStat
			cell = row.AddCell()
			cell.Value = lastLoginDt
			cell = row.AddCell()

			// nowDate := time.Now()
			// dateFmt := fmt.Sprintf(nowDate.Format("20060102150405"))

			// uploadPath, _ := beego.AppConfig.String("uploadpath")
			// imgDir := uploadPath + "/excel/entp/" + pEntpMemNo + "/"

			// // 폴더가 없을 경우 해당 폴더를 만들어준다.
			// if _, err := os.Stat(imgDir); os.IsNotExist(err) {
			// 	err = os.MkdirAll(imgDir, 0755)
			// 	if err != nil {
			// 		panic(err)
			// 	}
			// }

			// err = file.Save(imgDir + dateFmt + ".xlsx")
			// if err != nil {
			// 	fmt.Printf(err.Error())
			// }

			// /* Image Server Path */
			// imgServer, _ := beego.AppConfig.String("viewpath")
			downloadPath = imgServer + "/excel/entp/" + pEntpMemNo + "/" + dateFmt + ".xlsx"

			adminEntpListExcel = append(adminEntpListExcel, models.AdminEntpListExcel{
				ApplyDt:      applyDt,
				ApplyStat:    applyStat,
				EntpMemNo:    entpMemNo,
				JobFairCode:  jobFairCode,
				EntpKoNm:     entpKoNm,
				BizRegNo:     bizRegNo,
				RepNm:        repNm,
				EntpMemId:    entpMemId,
				PpChrgNm:     ppChrgNm,
				PpChrgTelNo:  ppChrgTelNo,
				TotAplyCnt:   totAplyCnt,
				NewAplyCnt:   newAplyCnt,
				VpYn:         vpYn,
				OsGbn:        osGbn,
				VerifyStat:   verifyStat,
				LastLoginDt:  lastLoginDt,
				DownloadPath: downloadPath,
			})

		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		uploadPath, _ := beego.AppConfig.String("uploadpath")
		imgDir := uploadPath + "/excel/entp/" + pEntpMemNo + "/"

		// 폴더가 없을 경우 해당 폴더를 만들어준다.
		if _, err := os.Stat(imgDir); os.IsNotExist(err) {
			err = os.MkdirAll(imgDir, 0755)
			if err != nil {
				panic(err)
			}
		}

		err = file.Save(imgDir + dateFmt + ".xlsx")
		if err != nil {
			fmt.Printf(err.Error())
		}

		rtnAdminEntpListExcel = models.RtnAdminEntpListExcel{
			RtnAdminEntpListExcelData: adminEntpListExcel,
		}
	}

	c.Data["json"] = &rtnAdminEntpListExcel
	c.ServeJSON()
}
