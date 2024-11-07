package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type ApplicantScoreRegUpdateController struct {
	beego.Controller
}

func (c *ApplicantScoreRegUpdateController) Post() {

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Data["json"] = &models.DefaultResult{RtnCd: 99, RtnMsg: "/login"}
		c.ServeJSON()
		return
	}

	memSn := session.Get(c.Ctx.Request.Context(), "mem_sn")
	if memSn == nil {
		c.Data["json"] = &models.DefaultResult{RtnCd: 99, RtnMsg: "/login"}
		c.ServeJSON()
		return
	}

	memId := session.Get(c.Ctx.Request.Context(), "mem_id")
	if memId == nil {
		c.Data["json"] = &models.DefaultResult{RtnCd: 99, RtnMsg: "/login"}
		c.ServeJSON()
		return
	}

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := mem_no
	pRecrutSn := c.GetString("recrut_sn")
	pPpMemNo := c.GetString("pp_mem_no")

	pEvalItem := c.GetString("eval_item")
	pResultValue := c.GetString("result_value")
	pResultComment := c.GetString("result_comment")

	pPpChrgSn := memSn
	//pMemId := memId
	//fmt.Printf(fmt.Sprintf("pPpChrgSn: %v, pMemId:%v", pPpChrgSn, pMemId))

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : ZSP_APPL_SCORE_REG_UPT
	fmt.Printf(fmt.Sprintf("CALL ZSP_APPL_SCORE_REG_UPT('%v', '%v', '%v', '%v', '%v', %v, '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo, pEvalItem, pResultValue, pResultComment, pPpChrgSn))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_APPL_SCORE_REG_UPT('%v', '%v', '%v', '%v', '%v', %v, '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo, pEvalItem, pResultValue, pResultComment, pPpChrgSn),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
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

	var (
		rtnCd  int64
		rtnMsg string
	)

	rtnRegUpdate := models.DefaultResult{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		fmt.Printf(fmt.Sprintf(" ===> rtnCd:%v, rtnMsg:%v", rtnCd, rtnMsg))

		fmt.Printf(fmt.Sprintf("[ApplicantScore][Save] EntpMemNo: %v, RecrutSn: %v, PpMemNo: %v, EvalItem: %v, ResultValue: %v, ResultComment: %v, PpChrgSn: %v",
			pEntpMemNo, pRecrutSn, pPpMemNo, pEvalItem, pResultValue, pResultComment, pPpChrgSn))

		rtnRegUpdate = models.DefaultResult{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}
	// End : ZSP_APPL_SCORE_REG_UPT

	c.Data["json"] = &rtnRegUpdate
	c.ServeJSON()
}
