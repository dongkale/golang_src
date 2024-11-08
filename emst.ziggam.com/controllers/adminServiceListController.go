package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"

	"emst.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

type AdminServiceListController struct {
	BaseController
}

func (c *AdminServiceListController) Get() {

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
	pBrdGbnCd := c.GetString("brd_gbn_cd")
	if pBrdGbnCd == "" {
		pBrdGbnCd = "03"
	}

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Service List
	log.Debug("CALL SP_EMS_SERVICE_DTL_R('%v', '%v', :1)",
		pLang, pBrdGbnCd)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_SERVICE_DTL_R('%v', '%v', :1)",
		pLang, pBrdGbnCd),
		ora.S, /* BRD_GBN_CD */
		ora.S, /* TITLE */
		ora.S, /* CONT */
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

	adminServiceDetail := make([]models.AdminServiceDetail, 0)

	var (
		brdGbnCd string
		title    string
		cont     string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			brdGbnCd = procRset.Row[0].(string)
			title = procRset.Row[1].(string)
			cont = procRset.Row[2].(string)

			adminServiceDetail = append(adminServiceDetail, models.AdminServiceDetail{
				BrdGbnCd: brdGbnCd,
				Title:    title,
				Cont:     cont,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Service List

	c.Data["BrdGbnCd"] = brdGbnCd
	c.Data["Title"] = title
	c.Data["Cont"] = cont
	c.Data["MenuId"] = "04"
	c.Data["SubMenuId"] = "04"
	c.TplName = "admin/service_write.html"
}
