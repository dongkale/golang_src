package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/rana/ora.v4"
)

type TeamCommentInsertController struct {
	beego.Controller
}

func (c *TeamCommentInsertController) Post() {

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
	pRecrutSn := c.GetString("recrut_sn")
	pPpMemNo := c.GetString("pp_mem_no")
	pPpChrgSn := c.GetString("pp_chrg_sn")
	pPpChrgCmt := c.GetString("pp_chrg_cmt")
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

	// Start : Comment Insert Process

	fmt.Printf(fmt.Sprintf("CALL ZSP_COMMENT_REG_PROC('%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo, pPpChrgSn, pPpChrgCmt))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_COMMENT_REG_PROC('%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo, pPpChrgSn, pPpChrgCmt),
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

	rtnTeamCommentInsert := models.RtnTeamCommentInsert{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnTeamCommentInsert = models.RtnTeamCommentInsert{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}
	// End : Comment Insert Process

	c.Data["json"] = &rtnTeamCommentInsert
	c.ServeJSON()
}
