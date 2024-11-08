package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"

	"emst.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

type AdminServiceProcessController struct {
	BaseController
}

func (c *AdminServiceProcessController) Post() {

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
	pBrdGbnCd := c.GetString("brd_gbn_cd") // 게시구분코드
	pTitle := c.GetString("title")         //제목
	notiDoc1 := c.GetString("notiDoc1")
	notiDoc2 := c.GetString("notiDoc2")
	notiDoc3 := c.GetString("notiDoc3")
	notiDoc4 := c.GetString("notiDoc4")
	notiDoc5 := c.GetString("notiDoc5")
	notiDoc6 := c.GetString("notiDoc6")
	notiDoc7 := c.GetString("notiDoc7")
	notiDoc8 := c.GetString("notiDoc8")
	notiDoc9 := c.GetString("notiDoc9")
	notiDoc10 := c.GetString("notiDoc10")
	notiDoc11 := c.GetString("notiDoc11")
	notiDoc12 := c.GetString("notiDoc12")
	notiDoc13 := c.GetString("notiDoc13")
	notiDoc14 := c.GetString("notiDoc14")
	notiDoc15 := c.GetString("notiDoc15")
	notiDoc16 := c.GetString("notiDoc16")
	notiDoc17 := c.GetString("notiDoc17")
	notiDoc18 := c.GetString("notiDoc18")
	notiDoc19 := c.GetString("notiDoc19")
	notiDoc20 := c.GetString("notiDoc20")
	notiDoc21 := c.GetString("notiDoc21")
	notiDoc22 := c.GetString("notiDoc22")
	notiDoc23 := c.GetString("notiDoc23")
	notiDoc24 := c.GetString("notiDoc24")
	notiDoc25 := c.GetString("notiDoc25")
	notiDoc26 := c.GetString("notiDoc26")
	notiDoc27 := c.GetString("notiDoc27")
	notiDoc28 := c.GetString("notiDoc28")
	notiDoc29 := c.GetString("notiDoc29")
	notiDoc30 := c.GetString("notiDoc30")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Admin Service Process
	log.Debug("CALL SP_EMS_SERVICE_PROC( "+
		"'%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pBrdGbnCd, pTitle, notiDoc1, notiDoc2, notiDoc3, notiDoc4, notiDoc5, notiDoc6, notiDoc7, notiDoc8, notiDoc9, notiDoc10,
		notiDoc11, notiDoc12, notiDoc13, notiDoc14, notiDoc15, notiDoc16, notiDoc17, notiDoc18, notiDoc19, notiDoc20,
		notiDoc21, notiDoc22, notiDoc23, notiDoc24, notiDoc25, notiDoc26, notiDoc27, notiDoc28, notiDoc29, notiDoc30)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_SERVICE_PROC( "+
		"'%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pBrdGbnCd, pTitle, notiDoc1, notiDoc2, notiDoc3, notiDoc4, notiDoc5, notiDoc6, notiDoc7, notiDoc8, notiDoc9, notiDoc10,
		notiDoc11, notiDoc12, notiDoc13, notiDoc14, notiDoc15, notiDoc16, notiDoc17, notiDoc18, notiDoc19, notiDoc20,
		notiDoc21, notiDoc22, notiDoc23, notiDoc24, notiDoc25, notiDoc26, notiDoc27, notiDoc28, notiDoc29, notiDoc30),
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

	rtnAdminServiceProcess := models.RtnAdminServiceProcess{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnAdminServiceProcess = models.RtnAdminServiceProcess{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}

	// End : Admin Service Process

	c.Data["json"] = &rtnAdminServiceProcess
	c.ServeJSON()
}
