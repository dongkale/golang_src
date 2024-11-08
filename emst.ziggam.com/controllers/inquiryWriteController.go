package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"

	"emst.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

type InquiryWriteController struct {
	BaseController
}

func (c *InquiryWriteController) Get() {

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/common/login")
	}

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := mem_no
	//pEntpMemNo := "E2018102500001"

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Inquiry Write
	log.Debug("CALL SP_EMS_INQUIRY_INFO_R('%v', '%v', :1)",
		pLang, pEntpMemNo)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_INQUIRY_INFO_R('%v', '%v', :1)",
		pLang, pEntpMemNo),
		ora.S, /* EMAIL */
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

	inquiryWrite := make([]models.InquiryWrite, 0)

	var (
		email string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			email = procRset.Row[0].(string)

			inquiryWrite = append(inquiryWrite, models.InquiryWrite{
				Email: email,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Inquiry Write

	c.Data["Email"] = email

	// End : Inquiry List

	c.TplName = "inquiry/inquiry_write.html"
}
