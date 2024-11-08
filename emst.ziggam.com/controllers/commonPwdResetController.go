package controllers

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"emst.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/rana/ora.v4"
)

type CommonPwdResetController struct {
	beego.Controller
}

func (c *CommonPwdResetController) Get() {
	c.TplName = "common/reset_pwd.html"
}

func (c *CommonPwdResetController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()
	temp_mem_id := session.Get(c.Ctx.Request.Context(), "temp_mem_id")
	log.Debug("temp_mem_id = %v", temp_mem_id)

	if temp_mem_id == nil {
		c.Ctx.Redirect(302, "/common/find/pwd")
	}

	pLang, _ := beego.AppConfig.String("lang")
	pMemId := temp_mem_id
	pPwd := c.GetString("pwd")
	//imgServer, _ := beego.AppConfig.String("viewpath")

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

	// Start : Password Change Process

	log.Debug("CALL SP_EMS_ENTP_RESET_PWD_PROC('%v', '%v', '%v', :1)",
		pLang, pMemId, sha_pPwd)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ENTP_RESET_PWD_PROC('%v', '%v', '%v', :1)",
		pLang, pMemId, sha_pPwd),
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

	rtnResetPwd := models.RtnResetPwd{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)

			if rtnCd == 1 {
				session.Delete(c.Ctx.Request.Context(), "temp_mem_id")
			}
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnResetPwd = models.RtnResetPwd{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}
	// End : Password Change Process

	c.Data["json"] = &rtnResetPwd
	c.ServeJSON()

}
