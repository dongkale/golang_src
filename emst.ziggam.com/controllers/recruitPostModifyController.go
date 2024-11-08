package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"

	"emst.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

type RecruitPostModifyController struct {
	BaseController
}

func (c *RecruitPostModifyController) Get() {

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

	pEntpMemNo := c.GetString("entp_mem_no")
	pRecrutSn := c.GetString("recrut_sn")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Recruit Modify
	log.Debug("CALL SP_EMS_RECRUIT_INFO_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_RECRUIT_INFO_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn),
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* RECRUT_TITLE */
		ora.S,   /* EMPL_TYP */
		ora.S,   /* UP_JOB_GRP */
		ora.S,   /* JOB_GRP */
		ora.S,   /* SEX */
		ora.I64, /* RECRUT_CNT */
		ora.S,   /* ROL */
		ora.S,   /* APLY_QUFCT */
		ora.S,   /* PRGS_MSG */
		ora.S,   /* ANS_LMT_TM_NM */
		ora.S,   /* RECRUT_DY */
		ora.S,   /* ACPT_PRID_NM */
		ora.S,   /* VD_TITLE1 */
		ora.S,   /* VD_TITLE2 */
		ora.S,   /* VD_TITLE3 */
		ora.S,   /* VD_TITLE4 */
		ora.S,   /* VD_TITLE5 */
		ora.S,   /* VD_TITLE_UPT_YN */
		ora.I64, /* REC_MAX_TM */

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

	recruitPostModify := make([]models.RecruitPostModify, 0)

	var (
		entpMemNo    string
		recrutSn     string
		recrutTitle  string
		emplTyp      string
		upJobGrp     string
		jobGrp       string
		sex          string
		recrutCnt    int64
		rol          string
		aplyQufct    string
		prgsMsg      string
		ansLmtTmNm   string
		recrutDy     string
		acptPridNm   string
		vdTitle1     string
		vdTitle2     string
		vdTitle3     string
		vdTitle4     string
		vdTitle5     string
		vdTitleUptYn string
		recMaxTm     int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			entpMemNo = procRset.Row[0].(string)
			recrutSn = procRset.Row[1].(string)
			recrutTitle = procRset.Row[2].(string)
			emplTyp = procRset.Row[3].(string)
			upJobGrp = procRset.Row[4].(string)
			jobGrp = procRset.Row[5].(string)
			sex = procRset.Row[6].(string)
			recrutCnt = procRset.Row[7].(int64)
			rol = procRset.Row[8].(string)
			aplyQufct = procRset.Row[9].(string)
			prgsMsg = procRset.Row[10].(string)
			ansLmtTmNm = procRset.Row[11].(string)
			recrutDy = procRset.Row[12].(string)
			acptPridNm = procRset.Row[13].(string)
			vdTitle1 = procRset.Row[14].(string)
			vdTitle2 = procRset.Row[15].(string)
			vdTitle3 = procRset.Row[16].(string)
			vdTitle4 = procRset.Row[17].(string)
			vdTitle5 = procRset.Row[18].(string)
			vdTitleUptYn = procRset.Row[19].(string)
			recMaxTm = procRset.Row[20].(int64)

			recruitPostModify = append(recruitPostModify, models.RecruitPostModify{
				EntpMemNo:    entpMemNo,
				RecrutSn:     recrutSn,
				RecrutTitle:  recrutTitle,
				EmplTyp:      emplTyp,
				UpJobGrp:     upJobGrp,
				JobGrp:       jobGrp,
				Sex:          sex,
				RecrutCnt:    recrutCnt,
				Rol:          rol,
				AplyQufct:    aplyQufct,
				PrgsMsg:      prgsMsg,
				AnsLmtTmNm:   ansLmtTmNm,
				RecrutDy:     recrutDy,
				AcptPridNm:   acptPridNm,
				VdTitle1:     vdTitle1,
				VdTitle2:     vdTitle2,
				VdTitle3:     vdTitle3,
				VdTitle4:     vdTitle4,
				VdTitle5:     vdTitle5,
				VdTitleUptYn: vdTitleUptYn,
				RecMaxTm:     recMaxTm,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Recruit Modify

	c.Data["EntpMemNo"] = entpMemNo
	c.Data["RecrutSn"] = recrutSn
	c.Data["RecrutTitle"] = recrutTitle
	c.Data["EmplTyp"] = emplTyp
	c.Data["UpJobGrp"] = upJobGrp
	c.Data["JobGrp"] = jobGrp
	c.Data["Sex"] = sex
	c.Data["RecrutCnt"] = recrutCnt
	c.Data["Rol"] = rol
	c.Data["AplyQufct"] = aplyQufct
	c.Data["PrgsMsg"] = prgsMsg
	c.Data["AnsLmtTmNm"] = ansLmtTmNm
	c.Data["RecrutDy"] = recrutDy
	c.Data["AcptPridNm"] = acptPridNm
	c.Data["VdTitle1"] = vdTitle1
	c.Data["VdTitle2"] = vdTitle2
	c.Data["VdTitle3"] = vdTitle3
	c.Data["VdTitle4"] = vdTitle4
	c.Data["VdTitle5"] = vdTitle5
	c.Data["VdTitleUptYn"] = vdTitleUptYn
	c.Data["RecMaxTm"] = recMaxTm
	c.Data["MenuId"] = "02"

	c.TplName = "recruit/recruit_modify.html"
}
