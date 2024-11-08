package controllers

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"emst.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

type AdminLoginController struct {
	beego.Controller
}

func (c *AdminLoginController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()
	pLang, _ := beego.AppConfig.String("lang")

	pMemId := c.GetString("mem_id")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	log.Debug("CALL SP_EMS_ADMIN_LOGIN_PROC( '%v', '%v', :1)",
		pLang, pMemId)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_LOGIN_PROC( '%v', '%v', :1)",
		pLang, pMemId),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
		ora.S,   /* MEM_NO */
		ora.S,   /* MEM_ID */
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
		memNo  string
		memId  string
	)

	login := models.Login{}
	rtnLogin := models.RtnLogin{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)

			if rtnCd == 1 {
				memNo = procRset.Row[2].(string)
				memId = procRset.Row[3].(string)

				login = models.Login{
					MemNo: memNo,
					MemId: memId,
				}

				// Set the session
				session.Set(c.Ctx.Request.Context(), "mem_no", memNo)
				session.Set(c.Ctx.Request.Context(), "mem_id", memId)

			}
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnLogin = models.RtnLogin{
			RtnCd:   rtnCd,
			RtnMsg:  rtnMsg,
			RtnData: login,
		}
	}

	c.Data["json"] = &rtnLogin
	c.ServeJSON()
}
