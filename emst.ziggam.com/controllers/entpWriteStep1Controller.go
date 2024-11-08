package controllers

import (
	"fmt"

	"emst.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	ora "gopkg.in/rana/ora.v4"
)

type EntpWriteStep1Controller struct {
	beego.Controller
}

func (c *EntpWriteStep1Controller) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log
	pLang, _ := beego.AppConfig.String("lang")

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
		pLang, "05")

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_SERVICE_DTL_R('%v', '%v', :1)",
		pLang, "05"),
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
		cont string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			cont = procRset.Row[2].(string)

			adminServiceDetail = append(adminServiceDetail, models.AdminServiceDetail{
				Cont: cont,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}

	log.Debug("CALL SP_EMS_SERVICE_DTL_R('%v', '%v', :1)",
		pLang, "04")

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_SERVICE_DTL_R('%v', '%v', :1)",
		pLang, "04"),
		ora.S, /* BRD_GBN_CD */
		ora.S, /* TITLE */
		ora.S, /* CONT */
	)
	defer stmtProcCall.Close()
	if err != nil {
		panic(err)
	}
	procRset = &ora.Rset{}
	_, err = stmtProcCall.Exe(procRset)

	if err != nil {
		panic(err)
	}

	adminServiceDetail = make([]models.AdminServiceDetail, 0)

	var (
		cont1 string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			cont1 = procRset.Row[2].(string)

			adminServiceDetail = append(adminServiceDetail, models.AdminServiceDetail{
				Cont1: cont1,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}

	// End : Service List

	c.Data["Cont"] = cont
	c.Data["Cont1"] = cont1

	c.TplName = "entp/entp_write1.html"
}
