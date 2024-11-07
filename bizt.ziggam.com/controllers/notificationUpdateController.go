package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/rana/ora.v4"
)

type NotificationUpdateController struct {
	beego.Controller
}

func (c *NotificationUpdateController) Post() {

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

	pRecrutSn := c.GetString("recrut_sn")
	pMemNo := c.GetString("mem_no")
	pNotiKndCd := c.GetString("knd_cd")
	pRegDt := c.GetString("reg_dt")
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

	// Start : Notification Confirm Process

	fmt.Printf(fmt.Sprintf("CALL ZSP_NOTIFICATION_CFRM_PROC('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pRecrutSn, pMemNo, pNotiKndCd, pRegDt))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_NOTIFICATION_CFRM_PROC('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pRecrutSn, pMemNo, pNotiKndCd, pRegDt),
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

	rtnNotificationUpdate := models.RtnNotificationUpdate{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnNotificationUpdate = models.RtnNotificationUpdate{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}
	// End : Notification Confirm Process

	c.Data["json"] = &rtnNotificationUpdate
	c.ServeJSON()

}
