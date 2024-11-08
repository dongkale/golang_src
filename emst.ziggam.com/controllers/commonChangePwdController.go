package controllers

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"

	"emst.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

type CommonChangePwdController struct {
	BaseController
}

func (c *CommonChangePwdController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "common/change_pwd.html"
}

func (c *CommonChangePwdController) Post() {
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

	pEntpMemNo := mem_no
	pCurrPwd := c.GetString("curr_pwd")
	pPwd := c.GetString("pwd")
	//imgServer, _ := beego.AppConfig.String("viewpath")

	currsha := sha512.New()
	currsha.Write([]byte(pCurrPwd))
	currsha_pPwd := hex.EncodeToString(currsha.Sum(nil))

	sha := sha512.New()
	sha.Write([]byte(pPwd))
	sha_pPwd := hex.EncodeToString(sha.Sum(nil))

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	log.Debug("CALL SP_EMS_CHANGE_PWD_PROC( '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, currsha_pPwd, sha_pPwd)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_CHANGE_PWD_PROC( '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, currsha_pPwd, sha_pPwd),
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

	rtnCommonChangePwd := models.RtnCommonChangePwd{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnCommonChangePwd = models.RtnCommonChangePwd{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}

	c.Data["json"] = &rtnCommonChangePwd
	c.ServeJSON()
}
