package controllers

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"bizt.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/rana/ora.v4"
)

type SettingMemberPwdUpdateController struct {
	beego.Controller
}

func (c *SettingMemberPwdUpdateController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/login")
	}
	mem_sn := session.Get(c.Ctx.Request.Context(), "mem_sn")

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := mem_no //c.GetString("entp_mem_no")
	pPpChrgSn := mem_sn  //c.GetString("entp_mem_no")

	pCurrPwd := c.GetString("curr_pwd") //현재 비밀번호
	pPwd := c.GetString("pwd")          //비밀번호

	currsha := sha512.New()
	currsha.Write([]byte(pCurrPwd))
	currsha_pwd := hex.EncodeToString(currsha.Sum(nil))

	sha := sha512.New()
	sha.Write([]byte(pPwd))
	sha_pwd := hex.EncodeToString(sha.Sum(nil))

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Setting Member Password Update Process

	fmt.Printf(fmt.Sprintf("CALL ZSP_MEM_PWD_UPT_PROC('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pPpChrgSn, currsha_pwd, sha_pwd))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_MEM_PWD_UPT_PROC('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pPpChrgSn, currsha_pwd, sha_pwd),
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

	rtnSettingMemberPwdUpdate := models.RtnSettingMemberPwdUpdate{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)

		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnSettingMemberPwdUpdate = models.RtnSettingMemberPwdUpdate{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}

	// End : Setting Member Password Update Process

	c.Data["json"] = &rtnSettingMemberPwdUpdate
	c.ServeJSON()
}
