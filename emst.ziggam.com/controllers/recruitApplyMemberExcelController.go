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

type RecruitApplyMemberExcelController struct {
	beego.Controller
}

func (c *RecruitApplyMemberExcelController) Post() {

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
	pEntpMemNo := c.GetString("entp_mem_no")          // 기업회원번호(세션)
	pEvlPrgsStatCd := c.GetString("evl_prgs_stat_cd") // 평가구분코드(00:전체, 02:대기, 03:합격, 04:불합격)
	pRecrutSn := c.GetString("recrut_sn")             // 채용일련번호

	if pEvlPrgsStatCd == "" {
		pEvlPrgsStatCd = "00"
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

	// Start : Recruit Apply Member Excel Download
	log.Debug("CALL SP_EMS_APPLY_MEM_EXCEL_R('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pEvlPrgsStatCd, pRecrutSn)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_APPLY_MEM_EXCEL_R('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pEvlPrgsStatCd, pRecrutSn),
		ora.S, /* RECRUT_TITLE */
		ora.S, /* EMPL_TYP */
		ora.S, /* JOB_GRP_NM */
		ora.S, /* RECRUT_DY */
		ora.S, /* RECRUT_EDT */
		ora.S, /* EVL_PRGS_STAT */
		ora.S, /* FAVR_APLY_PP_YN */
		ora.S, /* NM */
		ora.S, /* SEX */
		ora.S, /* AGE */
		ora.S, /* EMAIL */
		ora.S, /* APPLY_DT */
		ora.S, /* SHOOT_TM */
		ora.S, /* SHOOT_CNT */
		ora.S, /* LEFT_DY */
		ora.S, /* EVL_STAT_DT */
		ora.S, /* VP_YN */
		ora.S, /* LST_EDU */
		ora.S, /* CARR_DESC */
		ora.S, /* FRGN_LANG_ABLT_DESC */
		ora.S, /* ATCH_DATA_PATH */
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

	rtnRecruitApplyMemberExcel := models.RtnRecruitApplyMemberExcel{}
	recruitApplyMemberExcel := make([]models.RecruitApplyMemberExcel, 0)

	var (
		recrutTitle      string
		emplTyp          string
		jobGrpNm         string
		recrutDy         string
		recrutEdt        string
		evlPrgsStat      string
		favrAplyPpYn     string
		nm               string
		sex              string
		age              string
		email            string
		applyDt          string
		shootTm          string
		shootCnt         string
		leftDy           string
		evlStatDt        string
		vpYn             string
		lstEdu           string
		carrDesc         string
		frgnLangAbltDesc string
		atchDataPath     string
		downloadPath     string
	)

	row = sheet.AddRow()

	cell = row.AddCell()
	cell.Value = "채용 공고명"
	cell = row.AddCell()
	cell.Value = "고용 형태"
	cell = row.AddCell()
	cell.Value = "직무"
	cell = row.AddCell()
	cell.Value = "공고 기간"
	cell = row.AddCell()
	cell.Value = "공고 종료일"
	cell = row.AddCell()
	cell.Value = "구분"
	cell = row.AddCell()
	cell.Value = "관심 여부"
	cell = row.AddCell()
	cell.Value = "이름"
	cell = row.AddCell()
	cell.Value = "성별"
	cell = row.AddCell()
	cell.Value = "나이"
	cell = row.AddCell()
	cell.Value = "이메일"
	cell = row.AddCell()
	cell.Value = "최종 지원일자"
	cell = row.AddCell()
	cell.Value = "최종답변 소요시간"
	cell = row.AddCell()
	cell.Value = "지원 시도횟수"
	cell = row.AddCell()
	cell.Value = "결정 마감일"
	cell = row.AddCell()
	cell.Value = "결정일"
	cell = row.AddCell()
	cell.Value = "영상 프로필"
	cell = row.AddCell()
	cell.Value = "최종 학력"
	cell = row.AddCell()
	cell.Value = "경력"
	cell = row.AddCell()
	cell.Value = "보유기술/자격증"
	cell = row.AddCell()
	cell.Value = "외국어 능력"
	cell = row.AddCell()
	cell.Value = "첨부자료 링크"

	if procRset.IsOpen() {
		for procRset.Next() {
			recrutTitle = procRset.Row[0].(string)
			emplTyp = procRset.Row[1].(string)
			jobGrpNm = procRset.Row[2].(string)
			recrutDy = procRset.Row[3].(string)
			recrutEdt = procRset.Row[4].(string)
			evlPrgsStat = procRset.Row[5].(string)
			favrAplyPpYn = procRset.Row[6].(string)
			nm = procRset.Row[7].(string)
			sex = procRset.Row[8].(string)
			age = procRset.Row[9].(string)
			email = procRset.Row[10].(string)
			applyDt = procRset.Row[11].(string)
			shootTm = procRset.Row[12].(string)
			shootCnt = procRset.Row[13].(string)
			leftDy = procRset.Row[14].(string)
			evlStatDt = procRset.Row[15].(string)
			vpYn = procRset.Row[16].(string)
			lstEdu = procRset.Row[17].(string)
			carrDesc = procRset.Row[18].(string)
			frgnLangAbltDesc = procRset.Row[19].(string)
			atchDataPath = procRset.Row[20].(string)

			row = sheet.AddRow()

			cell = row.AddCell()
			cell.Value = recrutTitle
			cell = row.AddCell()
			cell.Value = emplTyp
			cell = row.AddCell()
			cell.Value = jobGrpNm
			cell = row.AddCell()
			cell.Value = recrutDy
			cell = row.AddCell()
			cell.Value = recrutEdt
			cell = row.AddCell()
			cell.Value = evlPrgsStat
			cell = row.AddCell()
			cell.Value = favrAplyPpYn
			cell = row.AddCell()
			cell.Value = nm
			cell = row.AddCell()
			cell.Value = sex
			cell = row.AddCell()
			cell.Value = age
			cell = row.AddCell()
			cell.Value = email
			cell = row.AddCell()
			cell.Value = applyDt
			cell = row.AddCell()
			cell.Value = shootTm
			cell = row.AddCell()
			cell.Value = shootCnt
			cell = row.AddCell()
			cell.Value = leftDy
			cell = row.AddCell()
			cell.Value = evlStatDt
			cell = row.AddCell()
			cell.Value = vpYn
			cell = row.AddCell()
			cell.Value = lstEdu
			cell = row.AddCell()
			cell.Value = carrDesc
			cell = row.AddCell()
			cell.Value = frgnLangAbltDesc
			cell = row.AddCell()
			cell.Value = atchDataPath

			nowDate := time.Now()
			dateFmt := fmt.Sprintf(nowDate.Format("20060102150405"))

			uploadPath, _ := beego.AppConfig.String("uploadpath")
			imgDir := uploadPath + "/excel/apply/" + pEntpMemNo + "/"

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

			/* Image Server Path */
			imgServer, _ := beego.AppConfig.String("viewpath")
			downloadPath = imgServer + "/excel/apply/" + pEntpMemNo + "/" + dateFmt + ".xlsx"

			recruitApplyMemberExcel = append(recruitApplyMemberExcel, models.RecruitApplyMemberExcel{
				RecrutTitle:      recrutTitle,
				EmplTyp:          emplTyp,
				JobGrpNm:         jobGrpNm,
				RecrutDy:         recrutDy,
				RecrutEdt:        recrutEdt,
				EvlPrgsStat:      evlPrgsStat,
				FavrAplyPpYn:     favrAplyPpYn,
				Nm:               nm,
				Sex:              sex,
				Age:              age,
				Email:            email,
				ApplyDt:          applyDt,
				ShootTm:          shootTm,
				ShootCnt:         shootCnt,
				LeftDy:           leftDy,
				EvlStatDt:        evlStatDt,
				VpYn:             vpYn,
				LstEdu:           lstEdu,
				CarrDesc:         carrDesc,
				FrgnLangAbltDesc: frgnLangAbltDesc,
				AtchDataPath:     atchDataPath,
				DownloadPath:     downloadPath,
			})

		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnRecruitApplyMemberExcel = models.RtnRecruitApplyMemberExcel{
			RtnRecruitApplyMemberExcelData: recruitApplyMemberExcel,
		}
	}

	c.Data["json"] = &rtnRecruitApplyMemberExcel
	c.ServeJSON()
}
