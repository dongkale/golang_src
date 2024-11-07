package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/rana/ora.v4"
)

type TeamMemberDeleteController struct {
	beego.Controller
}

func (c *TeamMemberDeleteController) Post() {

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

	pEntpMemNo := mem_no //c.GetString("entp_mem_no")
	pArrPpChrgSn := c.GetString("arr_pp_chrg_sn")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Team Member Delete Process

	fmt.Printf(fmt.Sprintf("CALL ZSP_TEAM_MEM_DEL_PROC('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pArrPpChrgSn))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_TEAM_MEM_DEL_PROC('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pArrPpChrgSn),
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

	rtnTeamMemberDelete := models.RtnTeamMemberDelete{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)

		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnTeamMemberDelete = models.RtnTeamMemberDelete{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}

	// End : Team Member Delete Process

	c.Data["json"] = &rtnTeamMemberDelete
	c.ServeJSON()
}
