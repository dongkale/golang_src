package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/rana/ora.v4"
)

type MsgTempUpdateController struct {
	beego.Controller
}

func (c *MsgTempUpdateController) Post() {

	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	// start : session
	session := c.StartSession()
	pLang, _ := beego.AppConfig.String("lang")
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	//mem_id := session.Get(c.Ctx.Request.Context(), "mem_id")
	// end : session

	if mem_no == nil {
		c.Ctx.Redirect(302, "/")
	}

	pEntpMemNo := mem_no

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	fmt.Printf(fmt.Sprintf("CALL ZSP_MSG_TMP_UPT_PROC( '%v', '%v', :1)",
		pLang, pEntpMemNo))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_MSG_TMP_UPT_PROC( '%v', '%v', :1)",
		pLang, pEntpMemNo),
		ora.S, /* RTN_CD */
		ora.S, /* RTN_MSG */
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
		rtnCd  string
		rtnMsg string
	)

	rtnMsgTempUpt := models.RtnMsgTempUpt{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(string)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnMsgTempUpt = models.RtnMsgTempUpt{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}

	c.Data["json"] = &rtnMsgTempUpt
	c.ServeJSON()
}
