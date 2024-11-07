package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/rana/ora.v4"
)

// LiveNvNStatVerifyController ...
type LiveNvNStatVerifyController struct {
	beego.Controller
}

// Post ...
func (c *LiveNvNStatVerifyController) Post() {

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/common/login")
	}

	pLang, _ := beego.AppConfig.String("lang")

	pGbnCd := c.GetString("gbn_cd")
	pVal := c.GetString("val")
	pEntpMemNo := mem_no
	pRecrutSn := c.GetString("recrut_sn")
	pPpMemNo := c.GetString("pp_mem_no")
	//imgServer, _  := beego.AppConfig.String("viewpath")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Message Verify Process

	fmt.Printf(fmt.Sprintf("CALL ZSP_LIVE_STAT_VERIFY('%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pGbnCd, pVal, pEntpMemNo, pRecrutSn, pPpMemNo))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_STAT_VERIFY('%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pGbnCd, pVal, pEntpMemNo, pRecrutSn, pPpMemNo),
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

	rtnMessageVerify := models.RtnMessageVerify{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnMessageVerify = models.RtnMessageVerify{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}
	// End : Message Verify Process

	c.Data["json"] = &rtnMessageVerify
	c.ServeJSON()

}
