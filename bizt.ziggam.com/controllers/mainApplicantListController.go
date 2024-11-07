package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	logs "github.com/beego/beego/v2/core/logs"
	"gopkg.in/rana/ora.v4"
)

type MainApplicantListController struct {
	beego.Controller
}

func (c *MainApplicantListController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/login")
	}

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := mem_no
	pGbnCd := c.GetString("gbn_cd") //구분코드
	if pGbnCd == "" {
		pGbnCd = "Y"
	}

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

	// Start : Applicant List

	fmt.Printf(fmt.Sprintf("CALL ZSP_MAIN_APLY_LIST_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pGbnCd))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_MAIN_APLY_LIST_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pGbnCd),
		ora.S,   /* RECRUT_SN */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* FAVR_APLY_PP_YN */
		ora.S,   /* PTO_PATH */
		ora.S,   /* NM */
		ora.S,   /* SEX */
		ora.I64, /* AGE */
		ora.S,   /* REG_DT */
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

	rtnApplicantList := models.RtnApplicantList{}
	applicantList := make([]models.ApplicantList, 0)

	var (
		mnaRecrutSn     string
		mnaPpMemNo      string
		mnaFavrAplyPpYn string
		mnaPtoPath      string
		mnaNm           string
		mnaSex          string
		mnaAge          int64
		mnaRegDt        string
		mnaFullPtoPath  string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			mnaRecrutSn = procRset.Row[0].(string)
			mnaPpMemNo = procRset.Row[1].(string)
			mnaFavrAplyPpYn = procRset.Row[2].(string)
			mnaPtoPath = procRset.Row[3].(string)
			if mnaPtoPath == "" {
				mnaFullPtoPath = mnaPtoPath
			} else {
				mnaFullPtoPath = imgServer + mnaPtoPath
			}
			mnaNm = procRset.Row[4].(string)
			mnaSex = procRset.Row[5].(string)
			mnaAge = procRset.Row[6].(int64)
			mnaRegDt = procRset.Row[7].(string)

			applicantList = append(applicantList, models.ApplicantList{
				MnaRecrutSn:     mnaRecrutSn,
				MnaPpMemNo:      mnaPpMemNo,
				MnaFavrAplyPpYn: mnaFavrAplyPpYn,
				MnaPtoPath:      mnaFullPtoPath,
				MnaNm:           mnaNm,
				MnaSex:          mnaSex,
				MnaAge:          mnaAge,
				MnaRegDt:        mnaRegDt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnApplicantList = models.RtnApplicantList{
			RtnApplicantListData: applicantList,
		}
	}
	// End : Applicant List

	c.Data["json"] = &rtnApplicantList
	c.ServeJSON()
}
