package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"

	"emst.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

type AdminEntpVideoUseUpdateController struct {
	beego.Controller
}

func (c *AdminEntpVideoUseUpdateController) Post() {

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
	pEntpMemNo := c.GetString("entp_mem_no") // 기업회원번호
	pUseYn := c.GetString("use_yn")          // 사용여부

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Admin Entp Video Use Process
	log.Debug("CALL SP_EMS_ENTP_VIDEO_USE_PROC('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pUseYn)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ENTP_VIDEO_USE_PROC('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pUseYn),
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

	rtnAdminEntpVideoUseUpdate := models.RtnAdminEntpVideoUseUpdate{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnAdminEntpVideoUseUpdate = models.RtnAdminEntpVideoUseUpdate{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}

	// End : Admin Entp Video Use Process

	c.Data["json"] = &rtnAdminEntpVideoUseUpdate
	c.ServeJSON()
}
