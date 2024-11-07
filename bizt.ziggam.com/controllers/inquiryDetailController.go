package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	ora "gopkg.in/rana/ora.v4"
)

type InquiryDetailController struct {
	BaseController
}

func (c *InquiryDetailController) Get() {

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
	pSn := c.GetString("sn")
	pPageNo := c.GetString("pn")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Inquiry Detail
	fmt.Printf(fmt.Sprintf("CALL ZSP_INQUIRY_DTL_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pSn))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_INQUIRY_DTL_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pSn),
		ora.S, /* REG_DT */
		ora.S, /* PP_CHRG_NM */
		ora.S, /* PP_CHRG_BP_NM */
		ora.S, /* INQ_GBN_NM */
		ora.S, /* IQN_TITLE */
		ora.S, /* IQN_CONT */
		ora.S, /* EMAIL */
		ora.S, /* ANS_YN */
		ora.S, /* ANS_DT */
		ora.S, /* ANS_CONT */
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

	inquiryDetail := make([]models.InquiryDetail, 0)

	var (
		regDt      string
		ppChrgNm   string
		ppChrgBpNm string
		inqGbnNm   string
		inqTitle   string
		inqCont    string
		email      string
		ansYn      string
		ansDt      string
		ansCont    string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			regDt = procRset.Row[0].(string)
			ppChrgNm = procRset.Row[1].(string)
			ppChrgBpNm = procRset.Row[2].(string)
			inqGbnNm = procRset.Row[3].(string)
			inqTitle = procRset.Row[4].(string)
			inqCont = procRset.Row[5].(string)
			email = procRset.Row[6].(string)
			ansYn = procRset.Row[7].(string)
			ansDt = procRset.Row[8].(string)
			ansCont = procRset.Row[9].(string)

			inquiryDetail = append(inquiryDetail, models.InquiryDetail{
				RegDt:      regDt,
				PpChrgNm:   ppChrgNm,
				PpChrgBpNm: ppChrgBpNm,
				InqGbnNm:   inqGbnNm,
				InqTitle:   inqTitle,
				InqCont:    inqCont,
				Email:      email,
				AnsYn:      ansYn,
				AnsDt:      ansDt,
				AnsCont:    ansCont,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Inquiry Detail

	c.Data["RegDt"] = regDt
	c.Data["PpChrgNm"] = ppChrgNm
	c.Data["PpChrgBpNm"] = ppChrgBpNm
	c.Data["InqGbnNm"] = inqGbnNm
	c.Data["InqTitle"] = inqTitle
	c.Data["InqCont"] = inqCont
	c.Data["Email"] = email
	c.Data["AnsYn"] = ansYn
	c.Data["AnsDt"] = ansDt
	c.Data["AnsCont"] = ansCont
	c.Data["PageNo"] = pPageNo

	c.Data["TMenuId"] = "M00"
	c.Data["SMenuId"] = "M00"

	c.TplName = "setting/inquiry_detail.html"
}
