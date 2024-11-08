package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"

	"emst.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

type AdminVersionInsertController struct {
	beego.Controller
}

func (c *AdminVersionInsertController) Post() {

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
	pAppVer := c.GetString("app_ver")      // 앱버전
	pOsGbn := c.GetString("os_gbn")        //os구분
	pAppVerCd := c.GetString("app_ver_cd") // 앱버전코드
	pFrcUptYn := c.GetString("frc_upt_yn") // 강제업데이트여부

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Version Insert
	log.Debug("CALL SP_EMS_ADMIN_VER_REG_PROC('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pAppVer, pOsGbn, pAppVerCd, pFrcUptYn)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_VER_REG_PROC('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pAppVer, pOsGbn, pAppVerCd, pFrcUptYn),
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

	rtnAdminVersionInsert := models.RtnAdminVersionInsert{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnAdminVersionInsert = models.RtnAdminVersionInsert{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}

	c.Data["json"] = &rtnAdminVersionInsert
	c.ServeJSON()
}
