package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/rana/ora.v4"
)

type WithDrawUpdateController struct {
	beego.Controller
}

func (c *WithDrawUpdateController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/login")
	}
	auth_cd := session.Get(c.Ctx.Request.Context(), "auth_cd")

	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := mem_no
	pPpChrgSn := c.GetString("pp_chrg_sn")
	pAuthCd := auth_cd
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

	// Start : Entp Withdraw Process

	fmt.Printf(fmt.Sprintf("CALL ZSP_WITHDRAW_PROC('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pPpChrgSn, pAuthCd))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_WITHDRAW_PROC('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pPpChrgSn, pAuthCd),
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

	rtnWithDrawUpdate := models.RtnWithDrawUpdate{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		// Delete the session
		session.Delete(c.Ctx.Request.Context(), "mem_no")
		session.Delete(c.Ctx.Request.Context(), "mem_id")
		session.Delete(c.Ctx.Request.Context(), "mem_sn")
		session.Delete(c.Ctx.Request.Context(), "auth_cd")

		// TOKEN KEY Cookie 삭제
		c.Ctx.SetCookie("token", "", 1<<31-1, "/")

		rtnWithDrawUpdate = models.RtnWithDrawUpdate{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}
	// End : Entp Withdraw Process

	c.Data["json"] = &rtnWithDrawUpdate
	c.ServeJSON()

}
