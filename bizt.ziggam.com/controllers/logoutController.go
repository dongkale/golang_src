package controllers

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"bizt.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

type LogoutController struct {
	beego.Controller
}

func (c *LogoutController) Get() {

	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	// start : session
	// 1. dvc_id checking
	//pLang, _ := beego.AppConfig.String("lang")
	session := c.StartSession()

	pLang, _ := beego.AppConfig.String("lang")

	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	mem_id := session.Get(c.Ctx.Request.Context(), "mem_id")

	if mem_no == nil {
		c.Ctx.Redirect(302, "/")
	}
	// end : session

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : 로그아웃
	fmt.Printf(fmt.Sprintf("CALL ZSP_LOGOUT_PROC( '%v', '%v', '%v', :1)",
		pLang, mem_no, mem_id))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_LOGOUT_PROC( '%v', '%v', '%v', :1)",
		pLang, mem_no, mem_id),
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

	rtnLogout := make([]models.RtnLogout, 0)
	var (
		rtnCd  string
		rtnMsg string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(string)
			rtnMsg = procRset.Row[1].(string)

			rtnLogout = append(rtnLogout, models.RtnLogout{
				RtnCd:  rtnCd,
				RtnMsg: rtnMsg,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : 로그아웃

	// Delete the session
	session.Delete(c.Ctx.Request.Context(), "mem_no")
	session.Delete(c.Ctx.Request.Context(), "mem_id")
	session.Delete(c.Ctx.Request.Context(), "mem_sn")
	session.Delete(c.Ctx.Request.Context(), "auth_cd")

	// TOKEN KEY Cookie 삭제
	c.Ctx.SetCookie("token", "", 1<<31-1, "/")

	// 메인으로 이동
	c.Ctx.Redirect(302, "/login")
}
