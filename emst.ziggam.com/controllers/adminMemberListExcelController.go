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

// AdminMemberListExcelController ...
type AdminMemberListExcelController struct {
	beego.Controller
}

// Post ...
func (c *AdminMemberListExcelController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Data["json"] = &models.RtnAdminMemberListExcel{}
		c.ServeJSON()
		return
	}

	pLang, _ := beego.AppConfig.String("lang")
	imgServer, _ := beego.AppConfig.String("viewpath")

	//pEntpMemNo := c.GetString("entp_mem_no") // 기업회원번호(세션)
	pEntpMemNo := mem_no.(string)

	pSex := c.GetString("sex")                      //성별 (무관:A, 남성:M, 여성:F)
	pAge := c.GetString("age")                      //연령 (전체:00, 19세이하:01, 20~29:02, 30~39:03, 40~49:04, 50~59:05, 60세이상:06)
	pVpYn := c.GetString("vp_yn")                   //영상프로필 (전체:9, 있음:1, 없음:0)
	pKeyword := c.GetString("keyword")              // 검색어
	pSdy := c.GetString("sdy")                      // 검색시작일자
	pEdy := c.GetString("edy")                      // 검색종료일자
	pMemStat := c.GetString("mem_stat")             // 회원상태
	pOsGbn := c.GetString("os_gbn")                 //OS구분 (전체:A, 안드로이드:AD, 애플:IS)
	pMemJoinGbnCd := c.GetString("mem_join_gbn_cd") //회원가입유형 (전체:A, 일반:00, 페이스북:01, 카카오:02)
	pJobFair := c.GetString("jf_mng_cd")

	pLoginSdy := c.GetString("login_sdy") // 로그인 검색시작일자
	pLoginEdy := c.GetString("login_edy") // 로그인 검색종료일자

	logs.Debug("login_sdy:" + pLoginSdy)
	logs.Debug("login_edy:" + pLoginEdy)

	if pMemStat == "" {
		pMemStat = "00"
	}

	if pSex == "" {
		pSex = "A"
	}

	if pAge == "" {
		pAge = "00"
	}

	if pVpYn == "" {
		pVpYn = "9"
	}

	if pOsGbn == "" {
		pOsGbn = "A"
	}

	if pMemJoinGbnCd == "" {
		pMemJoinGbnCd = "A"
	}

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

	// 가입일시	유형	이름	영상	아이디	성별	생년월일	나이	핸드폰번호	이메일	지원	OS	상태	참여 박람회	최근 접속 일시
	// SP_EMS_ADMIN_MEMBER_LST_EXCEL
	// pSex, pAge, pVpYn, pSdy, pEdy, pKeyword, pMemStat, pOsGbn, pMemJoinGbnCd, pJobFair, pLoginSdy, pLoginEdy

	// Start : Recruit Apply Member Excel Download
	log.Debug("CALL SP_EMS_ADMIN_MEMBER_LST_EXCEL('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pSex, pAge, pVpYn, pSdy, pEdy, pKeyword, pMemStat, pOsGbn, pMemJoinGbnCd, pJobFair, pLoginSdy, pLoginEdy)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_MEMBER_LST_EXCEL('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pSex, pAge, pVpYn, pSdy, pEdy, pKeyword, pMemStat, pOsGbn, pMemJoinGbnCd, pJobFair, pLoginSdy, pLoginEdy),
		ora.S, /* REG_DT */
		ora.S, /* MEM_JOIN_GBN_NM */
		ora.S, /* NM */
		ora.S, /* VP_YN */
		ora.S, /* MEM_ID */
		ora.S, /* SEX */
		ora.S, /* BRTH_YMD */
		ora.S, /* AGE */
		ora.S, /* MO_NO */
		ora.S, /* EMAIL */
		ora.S, /* AH_TOT_CNT */
		ora.S, /* OS_GBN */
		ora.S, /* MEM_STAT_NM */
		ora.S, /* JOBFAIR_CODE */
		ora.S, /* LOGIN_DT */
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

	rtnAdminMemberListExcel := models.RtnAdminMemberListExcel{}
	adminMemberListExcel := make([]models.AdminMemberListExcel, 0)

	var (
		applyDt      string
		memJoinGbnNm string
		nm           string
		vpYn         string
		memId        string
		sex          string
		brthYmd      string
		age          string
		email        string
		moNo         string
		ahTotCnt     string
		osGbn        string
		memStatNm    string
		jobFairCode  string
		loginDt      string
		downloadPath string
	)

	row = sheet.AddRow()

	cell = row.AddCell()
	cell.Value = "가입일시"
	cell = row.AddCell()
	cell.Value = "유형"
	cell = row.AddCell()
	cell.Value = "이름"
	cell = row.AddCell()
	cell.Value = "영상"
	cell = row.AddCell()
	cell.Value = "아이디"
	cell = row.AddCell()
	cell.Value = "성별"
	cell = row.AddCell()
	cell.Value = "생년월일"
	cell = row.AddCell()
	cell.Value = "나이"
	cell = row.AddCell()
	cell.Value = "핸드폰번호"
	cell = row.AddCell()
	cell.Value = "이메일"
	cell = row.AddCell()
	cell.Value = "지원"
	cell = row.AddCell()
	cell.Value = "OS"
	cell = row.AddCell()
	cell.Value = "상태"
	cell = row.AddCell()
	cell.Value = "참여 박람회"
	cell = row.AddCell()
	cell.Value = "최근 접속 일시"
	cell = row.AddCell()

	if procRset.IsOpen() {
		for procRset.Next() {
			applyDt = procRset.Row[0].(string)
			memJoinGbnNm = procRset.Row[1].(string)
			nm = procRset.Row[2].(string)
			vpYn = procRset.Row[3].(string)
			memId = procRset.Row[4].(string)
			sex = procRset.Row[5].(string)
			brthYmd = procRset.Row[6].(string)
			age = procRset.Row[7].(string)
			email = procRset.Row[8].(string)
			moNo = procRset.Row[9].(string)
			ahTotCnt = procRset.Row[10].(string)
			osGbn = procRset.Row[11].(string)
			memStatNm = procRset.Row[12].(string)
			jobFairCode = procRset.Row[13].(string)
			loginDt = procRset.Row[14].(string)

			row = sheet.AddRow()

			cell = row.AddCell()
			cell.Value = applyDt
			cell = row.AddCell()
			cell.Value = memJoinGbnNm
			cell = row.AddCell()
			cell.Value = nm
			cell = row.AddCell()
			cell.Value = vpYn
			cell = row.AddCell()
			cell.Value = memId
			cell = row.AddCell()
			cell.Value = sex
			cell = row.AddCell()
			cell.Value = brthYmd
			cell = row.AddCell()
			cell.Value = age
			cell = row.AddCell()
			cell.Value = email
			cell = row.AddCell()
			cell.Value = moNo
			cell = row.AddCell()
			cell.Value = ahTotCnt
			cell = row.AddCell()
			cell.Value = osGbn
			cell = row.AddCell()
			cell.Value = memStatNm
			cell = row.AddCell()
			cell.Value = jobFairCode
			cell = row.AddCell()
			cell.Value = loginDt

			// nowDate := time.Now()
			// dateFmt := fmt.Sprintf(nowDate.Format("20060102150405"))

			// uploadPath, _ := beego.AppConfig.String("uploadpath")
			// imgDir := uploadPath + "/excel/member/" + pEntpMemNo + "/"

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
			downloadPath = imgServer + "/excel/member/" + pEntpMemNo + "/" + dateFmt + ".xlsx"

			adminMemberListExcel = append(adminMemberListExcel, models.AdminMemberListExcel{
				ApplyDt:      applyDt,
				MemJoinGbnNm: memJoinGbnNm,
				Nm:           nm,
				VpYn:         vpYn,
				MemId:        memId,
				Sex:          sex,
				BrthYmd:      brthYmd,
				Age:          age,
				Email:        email,
				MoNo:         moNo,
				AhTotCnt:     ahTotCnt,
				OsGbn:        osGbn,
				MemStatNm:    memStatNm,
				JobFairCode:  jobFairCode,
				LoginDt:      loginDt,
				DownloadPath: downloadPath,
			})

		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		uploadPath, _ := beego.AppConfig.String("uploadpath")
		imgDir := uploadPath + "/excel/member/" + pEntpMemNo + "/"

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

		rtnAdminMemberListExcel = models.RtnAdminMemberListExcel{
			RtnAdminMemberListExcelData: adminMemberListExcel,
		}
	}

	c.Data["json"] = &rtnAdminMemberListExcel
	c.ServeJSON()
}
