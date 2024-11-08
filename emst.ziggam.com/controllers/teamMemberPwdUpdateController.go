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

type TeamMemberPwdUpdateController struct {
	beego.Controller
}

func (c *TeamMemberPwdUpdateController) Post() {

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Data["json"] = &models.DefaultResult{RtnCd: 99, RtnMsg: "/login"}
		c.ServeJSON()
		return
	}

	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := c.GetString("entp_mem_no")
	pPpChrgSn := c.GetString("pp_chrg_sn")
	pPwd := c.GetString("pwd")

	if len(strings.TrimSpace(pEntpMemNo)) == 0 {
		logs.Debug(fmt.Sprintf("[EntpTeamMemberPwdUpdate][Error] EntpMemNo:%v, Pwd:%v --> Invalid pEntpMemNo", pEntpMemNo, pPwd))

		c.Data["json"] = &models.DefaultResult{RtnCd: 2, RtnMsg: "Invalid EntpMemNo"}
		c.ServeJSON()
		return
	}

	if len(strings.TrimSpace(pPwd)) == 0 {
		logs.Debug(fmt.Sprintf("[EntpTeamMemberPwdUpdate][Error] EntpMemNo:%v, Pwd:%v --> Invalid Password", pEntpMemNo, pPwd))

		c.Data["json"] = &models.DefaultResult{RtnCd: 3, RtnMsg: "Invalid Password"}
		c.ServeJSON()
		return
	}

	sha := sha512.New()
	sha.Write([]byte(pPwd))
	sha_pwd := hex.EncodeToString(sha.Sum(nil))

	// sha_pwd2, _ := hex.DecodeString(sha_pwd)
	// logs.Debug(string(sha_pwd2))

	logs.Debug(fmt.Sprintf("[EntpTeamMemberPwdUpdate] EntpMemNo:%v, PpChrgSn:%v, Pwd:%v", pEntpMemNo, pPpChrgSn, pPwd))

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
	logs.Debug(fmt.Sprintf("CALL ZSP_TEAM_MEM_PWD_UPT_PROC('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pPpChrgSn, sha_pwd))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_TEAM_MEM_PWD_UPT_PROC('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pPpChrgSn, sha_pwd),
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

	rtnTeamMemberPwdUpdate := models.DefaultResult{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)

		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnTeamMemberPwdUpdate = models.DefaultResult{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}

		logs.Debug(fmt.Sprintf(" ===> RtnCd:%v, RtnMsg:%v", rtnCd, rtnMsg))
	}
	// End : Team Member Password Update Process

	c.Data["json"] = &rtnTeamMemberPwdUpdate
	c.ServeJSON()
}
