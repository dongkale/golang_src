package controllers

import (
	"fmt"
	"strings"

	"bizt.ziggam.com/models"
	"bizt.ziggam.com/tables"
	"bizt.ziggam.com/utils"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type ApplicantScorePopupController struct {
	BaseController
}

func (c *ApplicantScorePopupController) Get() {

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
	// cdnPath := beego.AppConfig.String("cdnpath")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start :ZSP_APPL_SCORE_BASE_INFO
	fmt.Printf(fmt.Sprintf("CALL ZSP_APPL_SCORE_BASE_INFO('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_APPL_SCORE_BASE_INFO('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo),
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* PTO_PATH */
		ora.S,   /* NM */
		ora.S,   /* SEX */
		ora.I64, /* AGE*/
		ora.S,   /* EAMIL */
		ora.S,   /* MO_NO */
		ora.S,   /* PRGS_STAT_CD */
		ora.S,   /* UP_JOB_GRP */
		ora.S,   /* JOB_GRP */
		ora.S,   /* RECRUT_TITLE */
		ora.S,   /* RECRUT_PROC_CD */
		ora.S,   /* APPLY_DT */
		ora.S,   /* READ_END_DAY */
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
		entpMemNo      string
		recrutSn       string
		ppMemNo        string
		ptoPath        string
		nm             string
		sex            string
		age            int64
		email          string
		moNo           string
		prgsStatCd     string
		upJobGrp       string
		jobGrp         string
		recrutTitle    string
		recrut_proc_cd string
		applyDt        string
		fullPtoPath    string
		readEndDay     string // LDK 2020-11-16 : 합격/불합격 처리 오류, 90일 체크 <-->
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
			recrut_proc_cd = procRset.Row[13].(string)
			applyDt = procRset.Row[14].(string)
			readEndDay = procRset.Row[15].(string) // "N" 정상 "Y" 90일이 넘었므로 비정상

			applicantDetailPopup = append(applicantDetailPopup, models.ApplicantDetailPopup{
				EntpMemNo:      entpMemNo,
				RecrutSn:       recrutSn,
				PpMemNo:        ppMemNo,
				PtoPath:        fullPtoPath,
				Nm:             nm,
				Sex:            sex,
				Age:            age,
				Email:          email,
				MoNo:           moNo,
				PrgsStatCd:     prgsStatCd,
				UpJobGrp:       upJobGrp,
				JobGrp:         jobGrp,
				RecrutTitle:    recrutTitle,
				Recrut_proc_cd: recrut_proc_cd,
				ApplyDt:        applyDt,
				ReadEndDay:     readEndDay,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : ZSP_APPL_SCORE_BASE_INFO

	// Start : ZSP_APPL_SCORE_BASE_INFO
	fmt.Printf(fmt.Sprintf("CALL ZSP_APPL_SCORE_INFO('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_APPL_SCORE_INFO('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo),
		ora.S, /* EVAL_ITEM */
		ora.S, /* RESULT_COMMENT */
		ora.S, /* PP_CHRG_NM */
		ora.S, /* PP_CHRG_BP_NM */
		ora.S, /* UPT_DT */
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

	//applicantScore := make([]models.ApplicantScore, 0)

	var (
		evalItemString string
		resultComment  string
		evalItemMap    map[string]string

		ppChrgNm   string
		ppChrgBpNm string

		uptDt    string
		uptDtFmt string
	)

	//evalItem = make(map[int]int)

	if procRset.IsOpen() {
		for procRset.Next() {

			evalItemString = procRset.Row[0].(string)
			resultComment = procRset.Row[1].(string)

			ppChrgNm = procRset.Row[2].(string)
			ppChrgBpNm = procRset.Row[3].(string)

			uptDt = procRset.Row[4].(string)

			uptDtFmt = utils.DtToFmtDt(uptDt, "06/01/02 15:04")
		}

		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : ZSP_APPL_SCORE_INFO

	if len(evalItemString) > 0 {
		evalItemMap = utils.StringDelimSplit(evalItemString, ",", ":")
	} else {
		evalItemMap = map[string]string{}
	}

	fmt.Printf(fmt.Sprintf("evalItemString:%v(%v), resultComment:%v ", evalItemString, evalItemMap, resultComment))

	// Start : ZSP_APPL_SCORE_EVAL_ITEM
	fmt.Printf(fmt.Sprintf("CALL ZSP_APPL_SCORE_EVAL_ITEM('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_APPL_SCORE_EVAL_ITEM('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn),
		ora.I64, /* NUM */
		ora.S,   /* TITLE */
		ora.S,   /* CATEGORY */
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

	var (
		num      int64
		title    string
		category string
	)

	applicantScoreEvalItem := make([]models.ApplicantScoreEvalItem, 0)

	if procRset.IsOpen() {
		for procRset.Next() {

			num = procRset.Row[0].(int64)
			title = procRset.Row[1].(string)
			category = procRset.Row[2].(string)

			applicantScoreEvalItem = append(applicantScoreEvalItem, models.ApplicantScoreEvalItem{
				Num:      num,
				Title:    title,
				Category: category,
			})
		}

		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : ZSP_APPL_SCORE_EVAL_ITEM

	// Start : ZSP_APPL_SCORE_EVAL_ITEM_CAT
	categoryArray := make([]string, 0)
	categoryArrayDup := make([]string, 0)

	for _, val := range applicantScoreEvalItem {
		categoryArray = append(categoryArray, val.Category)
	}

	err = utils.RemoveDuplicateValues(categoryArray, &categoryArrayDup)
	if err != nil {
		fmt.Printf(fmt.Sprintf("error filtering int slice: %v\n", err))
	}

	categoryString := strings.Join(categoryArrayDup, ",")

	fmt.Printf(fmt.Sprintf("CALL ZSP_APPL_SCORE_EVAL_ITEM_CAT('%v', '%v', :1)",
		pLang, categoryString))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_APPL_SCORE_EVAL_ITEM_CAT('%v', '%v', :1)",
		pLang, categoryString),
		ora.S, /* CATEGORY */
		ora.S, /* TITLE */
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

	var (
		itemCategory string
		itemTitle    string
	)

	applicantScoreEvalItemCategory := make([]models.ApplicantScoreEvalItemCategory, 0)

	if procRset.IsOpen() {
		for procRset.Next() {

			itemCategory = procRset.Row[0].(string)
			itemTitle = procRset.Row[1].(string)

			applicantScoreEvalItemCategory = append(applicantScoreEvalItemCategory, models.ApplicantScoreEvalItemCategory{
				Category: itemCategory,
				Title:    itemTitle,
			})
		}

		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : ZSP_APPL_SCORE_EVAL_ITEM_CAT

	c.Data["EntpMemNo"] = entpMemNo
	c.Data["RecrutSn"] = recrutSn
	c.Data["PpMemNo"] = ppMemNo
	c.Data["PtoPath"] = fullPtoPath
	c.Data["Nm"] = nm
	c.Data["Sex"] = sex
	c.Data["Age"] = age
	c.Data["Email"] = email
	c.Data["Mono"] = moNo
	c.Data["RecrutTitle"] = recrutTitle
	c.Data["UpJobGrp"] = upJobGrp
	c.Data["JobGrp"] = jobGrp
	c.Data["ApplyDt"] = applyDt
	c.Data["ReadEndDay"] = readEndDay

	//c.Data["EvalItemListOld"] = evalItemMap

	c.Data["EvalItemList"] = utils.ToJsonString(evalItemMap)
	c.Data["ResultComment"] = resultComment

	c.Data["PpChrgNm"] = ppChrgNm
	c.Data["PpChrgBpNm"] = ppChrgBpNm

	c.Data["UptDtFmt"] = uptDtFmt

	c.Data["EvalItemGradeTbl__"] = tables.TableConf.EvalItemGradeTbl__

	c.Data["EvalItemTable"] = utils.ToJsonString(applicantScoreEvalItem)
	c.Data["EvalItemCategoryTable"] = utils.ToJsonString(applicantScoreEvalItemCategory)

	c.TplName = "applicant/applicant_score_popup.html"
}

// 사용안함 !!!!
func (c *ApplicantScorePopupController) Post() {

	session := c.StartSession()

	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Data["json"] = &models.DefaultResult{RtnCd: 99, RtnMsg: "/login"}
		c.ServeJSON()
		return
	}

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := mem_no
	pRecrutSn := c.GetString("recrut_sn")
	pPpMemNo := c.GetString("pp_mem_no")

	imgServer, _  := beego.AppConfig.String("viewpath")

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
	fmt.Printf(fmt.Sprintf("CALL ZSP_APPL_SCORE_BASE_INFO('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_APPL_SCORE_BASE_INFO('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo),
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* PTO_PATH */
		ora.S,   /* NM */
		ora.S,   /* SEX */
		ora.I64, /* AGE*/
		ora.S,   /* EAMIL */
		ora.S,   /* MO_NO */
		ora.S,   /* PRGS_STAT_CD */
		ora.S,   /* UP_JOB_GRP */
		ora.S,   /* JOB_GRP */
		ora.S,   /* RECRUT_TITLE */
		ora.S,   /* RECRUT_PROC_CD */
		ora.S,   /* APPLY_DT */
		ora.S,   /* READ_END_DAY */
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
		entpMemNo      string
		recrutSn       string
		ppMemNo        string
		ptoPath        string
		nm             string
		sex            string
		age            int64
		email          string
		moNo           string
		prgsStatCd     string
		upJobGrp       string
		jobGrp         string
		recrutTitle    string
		recrut_proc_cd string
		applyDt        string
		fullPtoPath    string
		readEndDay     string // LDK 2020-11-16 : 합격/불합격 처리 오류, 90일 체크 <-->
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
			recrut_proc_cd = procRset.Row[13].(string)
			applyDt = procRset.Row[14].(string)
			readEndDay = procRset.Row[15].(string) // "N" 정상 "Y" 90일이 넘었므로 비정상

			applicantDetailPopup = append(applicantDetailPopup, models.ApplicantDetailPopup{
				EntpMemNo:      entpMemNo,
				RecrutSn:       recrutSn,
				PpMemNo:        ppMemNo,
				PtoPath:        fullPtoPath,
				Nm:             nm,
				Sex:            sex,
				Age:            age,
				Email:          email,
				MoNo:           moNo,
				PrgsStatCd:     prgsStatCd,
				UpJobGrp:       upJobGrp,
				JobGrp:         jobGrp,
				RecrutTitle:    recrutTitle,
				Recrut_proc_cd: recrut_proc_cd,
				ApplyDt:        applyDt,
				ReadEndDay:     readEndDay,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Applicant Popup Info

	fmt.Printf(fmt.Sprintf("CALL ZSP_APPL_SCORE_INFO('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_APPL_SCORE_INFO('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo),
		ora.S, /* EVAL_ITEM */
		ora.S, /* RESULT_COMMENT */
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

	applicantScore := make([]models.ApplicantScore, 0)

	var (
		evalItemString string
		resultComment  string
		evalItemMap    map[string]string
	)

	//evalItem = make(map[int]int)

	if procRset.IsOpen() {
		for procRset.Next() {

			evalItemString = procRset.Row[0].(string)
			resultComment = procRset.Row[1].(string)

			if len(evalItemString) > 0 {
				evalItemMap = utils.StringDelimSplit(evalItemString, ",", ":")
			} else {
				evalItemMap = map[string]string{}
			}

			applicantScore = append(applicantScore, models.ApplicantScore{
				EntpMemNo: entpMemNo,
				RecrutSn:  recrutSn,
				PpMemNo:   ppMemNo,

				EvalItemString: utils.ToJsonString(evalItemMap),
				ResultComment:  resultComment,
			})
		}

		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}

	fmt.Printf(fmt.Sprintf("evalItemString:%v(%v), resultComment:%v ", evalItemString, evalItemMap, resultComment))

	c.Data["json"] = &models.RslApplicantScore{Rtn: models.DefaultResult{RtnCd: 1, RtnMsg: "SUCCESS"}, RslDetailInfo: applicantDetailPopup, RslScorreInfo: applicantScore}
	c.ServeJSON()
}
