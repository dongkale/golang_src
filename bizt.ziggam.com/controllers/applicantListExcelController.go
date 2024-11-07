package controllers

import (
	"fmt"
	"os"
	"time"

	"bizt.ziggam.com/models"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/tealeg/xlsx"
	ora "gopkg.in/rana/ora.v4"
)

type ApplicantListExcelController struct {
	beego.Controller
}

func (c *ApplicantListExcelController) Post() {

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
	imgServer, _ := beego.AppConfig.String("viewpath")

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := mem_no                              //c.GetString("entp_mem_no")
	pRecrutSn := c.GetString("recrut_sn")             // 채용일련번호
	pEvlPrgsStat := c.GetString("evl_prgs_stat")      // 평가진행상태 (00:전체, 02:대기, 03:합격, 04:불합격)
	pSex := c.GetString("sex")                        //성별 (무관:A, 남성:M, 여성:F)
	pAge := c.GetString("age")                        //연령 (전체:00, 19세이하:01, 20~29:02, 30~39:03, 40~49:04, 50~59:05, 60세이상:06)
	pVpYn := c.GetString("vp_yn")                     //영상프로필 (전체:9, 있음:1, 없음:0)
	pFavrAplyPp := c.GetString("favr_aply_pp")        //관심지원자 (전체:9, 있음:1, 없음:0)
	pSortGbn := c.GetString("sort_gbn")               // 정렬구분(01:신규순, 02:마감순)
	pKeyword := c.GetString("keyword")                // 검색어
	pLiveReqStatCd := c.GetString("live_req_stat_cd") // 라이브요청상태코드(전체:A, 03:예정)
	pJobGrpCd := c.GetString("job_grp_cd")            // 직군코드

	if pEvlPrgsStat == "" {
		pEvlPrgsStat = "00"
	}

	if pSortGbn == "" {
		pSortGbn = "01"
	}

	if pRecrutSn == "" {
		pRecrutSn = "A"
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

	if pFavrAplyPp == "" {
		pFavrAplyPp = "9"
	}

	if pLiveReqStatCd == "" {
		pLiveReqStatCd = "A"
	}

	if pJobGrpCd == "" {
		pJobGrpCd = "A"
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

	// Start : Applicant List Excel
	fmt.Printf(fmt.Sprintf("CALL ZSP_APPLICANT_LIST_EXCEL_R('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pEvlPrgsStat, pSex, pAge, pVpYn, pFavrAplyPp, pSortGbn, pKeyword, pLiveReqStatCd, pJobGrpCd))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_APPLICANT_LIST_EXCEL_R('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pEvlPrgsStat, pSex, pAge, pVpYn, pFavrAplyPp, pSortGbn, pKeyword, pLiveReqStatCd, pJobGrpCd),
		ora.S, /* APPLY_DT */
		ora.S, /* FAVR_APLY_PP_YN */
		ora.S, /* NM */
		ora.S, /* SEX */
		ora.S, /* BRTH_YMD */
		ora.S, /* AGE */
		ora.S, /* EMAIL */
		ora.S, /* MO_NO */
		ora.S, /* EVL_PRGS_STAT_NM */
		ora.S, /* EVL_STAT_DT */
		ora.S, /* PRGS_STAT_CD */
		ora.S, /* UP_JOB_GRP */
		ora.S, /* JOB_GRP */
		ora.S, /* RECRUT_TITLE */
		ora.S, /* RECRUT_DY */
		ora.S, /* LST_EDU */
		ora.S, /* CARR_GBN */
		ora.S, /* CARR_DESC */
		ora.S, /* FRGN_LANG_ABLT_DESC */
		ora.S, /* ATCH_DATA_PATH */
		ora.S, /* TECH_QLFT_DESC */
		ora.S, /* ATCH_FILE_PATH_YN */
		ora.S, /* VP_YN */
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

	rtnApplicantListExcel := models.RtnApplicantListExcel{}
	applicantListExcel := make([]models.ApplicantListExcel, 0)

	var (
		aleApplyDt          string
		aleFavrAplyPpYn     string
		aleNm               string
		aleSex              string
		aleBrthYmd          string
		aleAge              string
		aleEmail            string
		aleMoNo             string
		aleEvlPrgsStatNm    string
		aleEvlStatDt        string
		alePrgsStatCd       string
		aleUpJobGrp         string
		aleJobGrp           string
		aleRecrutTitle      string
		aleRecrutDy         string
		aleLstEdu           string
		aleCarrGbn          string
		aleCarrDesc         string
		aleFrgnLangAbltDesc string
		aleAtchDataPath     string
		aleTechQlftDesc     string
		aleAtchFilePathYn   string
		aleVpYn             string
		downloadPath        string
	)

	row = sheet.AddRow()

	cell = row.AddCell()
	cell.Value = "지원일시"
	cell = row.AddCell()
	cell.Value = "관심 여부"
	cell = row.AddCell()
	cell.Value = "이름"
	cell = row.AddCell()
	cell.Value = "성별"
	cell = row.AddCell()
	cell.Value = "생년월일"
	cell = row.AddCell()
	cell.Value = "이메일"
	cell = row.AddCell()
	cell.Value = "휴대폰 번호"
	cell = row.AddCell()
	cell.Value = "지원자 상태"
	cell = row.AddCell()
	cell.Value = "처리일시"
	cell = row.AddCell()
	cell.Value = "공고상태"
	cell = row.AddCell()
	cell.Value = "지원직무"
	cell = row.AddCell()
	cell.Value = "지원공고"
	cell = row.AddCell()
	cell.Value = "모집기간"
	cell = row.AddCell()
	cell.Value = "최종학력"
	cell = row.AddCell()
	cell.Value = "경력구분"
	cell = row.AddCell()
	cell.Value = "경력사항"
	cell = row.AddCell()
	cell.Value = "보유기술/자격증정보"
	cell = row.AddCell()
	cell.Value = "외국어 능력"
	cell = row.AddCell()
	cell.Value = "첨부자료링크"
	cell = row.AddCell()
	cell.Value = "첨부파일여부"
	cell = row.AddCell()
	cell.Value = "영상프로필여부"

	if procRset.IsOpen() {
		for procRset.Next() {
			aleApplyDt = procRset.Row[0].(string)           // 지원 일시
			aleFavrAplyPpYn = procRset.Row[1].(string)      // 관심 여부
			aleNm = procRset.Row[2].(string)                // 이름
			aleSex = procRset.Row[3].(string)               // 성별
			aleBrthYmd = procRset.Row[4].(string)           // 생년월일
			aleAge = procRset.Row[5].(string)               // 나이
			aleEmail = procRset.Row[6].(string)             // email
			aleMoNo = procRset.Row[7].(string)              // 폰번호
			aleEvlPrgsStatNm = procRset.Row[8].(string)     // 지원자 상태
			aleEvlStatDt = procRset.Row[9].(string)         // 처리일시
			alePrgsStatCd = procRset.Row[10].(string)       // 공고상태
			aleUpJobGrp = procRset.Row[11].(string)         // 지원직무
			aleJobGrp = procRset.Row[12].(string)           // 지원직군
			aleRecrutTitle = procRset.Row[13].(string)      // 지원공고
			aleRecrutDy = procRset.Row[14].(string)         // 모집기간
			aleLstEdu = procRset.Row[15].(string)           // 최종학력
			aleCarrGbn = procRset.Row[16].(string)          // 경력구분
			aleCarrDesc = procRset.Row[17].(string)         // 경력사항
			aleFrgnLangAbltDesc = procRset.Row[18].(string) // 보유기술/자격증정보
			aleAtchDataPath = procRset.Row[19].(string)     // 외국어 능력
			aleTechQlftDesc = procRset.Row[20].(string)     // 첨부자료링크
			aleAtchFilePathYn = procRset.Row[21].(string)   // 첨부파일여부
			aleVpYn = procRset.Row[22].(string)             // 영상프로필여부

			row = sheet.AddRow()

			cell = row.AddCell()
			cell.Value = aleApplyDt
			cell = row.AddCell()
			cell.Value = aleFavrAplyPpYn
			cell = row.AddCell()
			cell.Value = aleNm
			cell = row.AddCell()
			cell.Value = aleSex
			cell = row.AddCell()
			cell.Value = aleBrthYmd + " (" + aleAge + ")"
			cell = row.AddCell()
			cell.Value = aleEmail
			cell = row.AddCell()
			cell.Value = aleMoNo
			cell = row.AddCell()
			cell.Value = aleEvlPrgsStatNm
			cell = row.AddCell()
			cell.Value = aleEvlStatDt
			cell = row.AddCell()
			cell.Value = alePrgsStatCd
			cell = row.AddCell()
			cell.Value = aleUpJobGrp + " > " + aleJobGrp
			cell = row.AddCell()
			cell.Value = aleRecrutTitle
			cell = row.AddCell()
			cell.Value = aleRecrutDy
			cell = row.AddCell()
			cell.Value = aleLstEdu
			cell = row.AddCell()
			cell.Value = aleCarrGbn
			cell = row.AddCell()
			cell.Value = aleCarrDesc
			cell = row.AddCell()
			cell.Value = aleTechQlftDesc
			cell = row.AddCell()
			cell.Value = aleFrgnLangAbltDesc
			cell = row.AddCell()
			cell.Value = aleAtchDataPath
			cell = row.AddCell()
			cell.Value = aleAtchFilePathYn
			cell = row.AddCell()
			cell.Value = aleVpYn

			downloadPath = imgServer + "/excel/apply/" + pEntpMemNo.(string) + "/" + dateFmt + ".xlsx"

			applicantListExcel = append(applicantListExcel, models.ApplicantListExcel{
				AleApplyDt:          aleApplyDt,
				AleFavrAplyPpYn:     aleFavrAplyPpYn,
				AleNm:               aleNm,
				AleSex:              aleSex,
				AleBrthYmd:          aleBrthYmd,
				AleAge:              aleAge,
				AleEmail:            aleEmail,
				AleMoNo:             aleMoNo,
				AleEvlPrgsStatNm:    aleEvlPrgsStatNm,
				AleEvlStatDt:        aleEvlStatDt,
				AlePrgsStatCd:       alePrgsStatCd,
				AleUpJobGrp:         aleUpJobGrp,
				AleJobGrp:           aleJobGrp,
				AleRecrutTitle:      aleRecrutTitle,
				AleRecrutDy:         aleRecrutDy,
				AleLstEdu:           aleLstEdu,
				AleCarrGbn:          aleCarrGbn,
				AleCarrDesc:         aleCarrDesc,
				AleFrgnLangAbltDesc: aleFrgnLangAbltDesc,
				AleAtchDataPath:     aleAtchDataPath,
				AleTechQlftDesc:     aleTechQlftDesc,
				AleAtchFilePathYn:   aleAtchFilePathYn,
				AleVpYn:             aleVpYn,
				DownloadPath:        downloadPath,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		uploadPath, _ := beego.AppConfig.String("uploadpath")
		imgDir := uploadPath + "/excel/apply/" + pEntpMemNo.(string) + "/"

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

		rtnApplicantListExcel = models.RtnApplicantListExcel{
			RtnApplicantListExcelData: applicantListExcel,
		}
	}

	c.Data["json"] = &rtnApplicantListExcel
	c.ServeJSON()
}
