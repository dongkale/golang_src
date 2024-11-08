package controllers

import (
	"fmt"

	"emst.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/rana/ora.v4"
)

type CommonFindPwdCertController struct {
	beego.Controller
}

func (c *CommonFindPwdCertController) Post() {
	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	pLang, _ := beego.AppConfig.String("lang")
	pMemId := c.GetString("mem_id")
	pPpChrgNm := c.GetString("pp_chrg_nm")
	pEmail := c.GetString("email")
	pCertNo := c.GetString("cert_no")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Certification Key Info

	log.Debug("CALL SP_EMS_FIND_PWD_STEP2_R('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pMemId, pPpChrgNm, pEmail, pCertNo)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_FIND_PWD_STEP2_R('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pMemId, pPpChrgNm, pEmail, pCertNo),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
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
		rtnCd     int64
		rtnMsg    string
		tempMemId string
	)

	rtnFindPwdStep2 := models.RtnFindPwdStep2{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
			tempMemId = procRset.Row[2].(string)

			// Temp Set the session
			session.Set(c.Ctx.Request.Context(), "temp_mem_id", tempMemId)

		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnFindPwdStep2 = models.RtnFindPwdStep2{
			RtnCd:     rtnCd,
			RtnMsg:    rtnMsg,
			TempMemId: tempMemId,
		}
	}
	// End : Certification Key Info

	c.Data["json"] = &rtnFindPwdStep2
	c.ServeJSON()
}
