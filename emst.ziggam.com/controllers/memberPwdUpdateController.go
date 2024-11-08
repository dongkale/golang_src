package controllers

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"
	
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	"gopkg.in/rana/ora.v4"
	"emst.ziggam.com/models"
)

type MemberPwdUpdateController struct {
	beego.Controller
}

func (c *MemberPwdUpdateController) Post() {

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Data["json"] = &models.DefaultResult{RtnCd: 99, RtnMsg: "/login"}
		c.ServeJSON()
		return
	}

	pLang, _ := beego.AppConfig.String("lang")

	pPpMemNo := c.GetString("pp_mem_no")
	pPwd := c.GetString("pwd")

	if len(strings.TrimSpace(pPpMemNo)) == 0 {
		logs.Debug(fmt.Sprintf("[MemberPwdUpdate][Error] pPpMemNo:%v, Pwd:%v --> Invalid PpMemNo", pPpMemNo, pPwd))

		c.Data["json"] = &models.DefaultResult{RtnCd: 2, RtnMsg: "Invalid PpMemNo"}
		c.ServeJSON()
		return
	}

	if len(strings.TrimSpace(pPwd)) == 0 {
		logs.Debug(fmt.Sprintf("[MemberPwdUpdate][Error] pPpMemNo:%v, Pwd:%v --> Invalid Password", pPpMemNo, pPwd))

		c.Data["json"] = &models.DefaultResult{RtnCd: 3, RtnMsg: "Invalid Password"}
		c.ServeJSON()
		return
	}

	sha := sha512.New()
	sha.Write([]byte(pPwd))
	sha_pwd := hex.EncodeToString(sha.Sum(nil))

	logs.Debug(fmt.Sprintf("[MemberPwdUpdate] pPpMemNo:%v, Pwd:%v", pPpMemNo, pPwd))

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Team Member Password Update Process
	logs.Debug(fmt.Sprintf("CALL SP2_MEM_RESET_PWD_PROC('%v', '%v', '%v', :1)",
		pLang, pPpMemNo, sha_pwd))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP2_MEM_RESET_PWD_PROC('%v', '%v', '%v', :1)",
		pLang, pPpMemNo, sha_pwd),
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

	rtnMemberPwdUpdate := models.DefaultResult{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)

		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnMemberPwdUpdate = models.DefaultResult{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}

		logs.Debug(fmt.Sprintf(" ===> RtnCd: %v, RtnMsg: %v", rtnCd, rtnMsg))
	}
	// End : Team Member Password Update Process

	c.Data["json"] = &rtnMemberPwdUpdate
	c.ServeJSON()
}
