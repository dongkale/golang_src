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

type TeamMemberInsertController struct {
	beego.Controller
}

func (c *TeamMemberInsertController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/login")
	}

	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := mem_no //c.GetString("entp_mem_no")
	pPpChrgNm := c.GetString("pp_chrg_nm")
	pPpChrgBpNm := c.GetString("pp_chrg_bp_nm")
	pPpChrgTelNo := c.GetString("pp_chrg_tel_no")
	pPpChrgSmsRecvYn := c.GetString("pp_chrg_sms_recv_yn")
	pPpChrgEmail := c.GetString("pp_chrg_email")
	pPpChrgEmailRecvYn := c.GetString("pp_chrg_email_recv_yn")
	pEntpId := c.GetString("entp_mem_id")
	pPwd := c.GetString("pwd")

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

	// Start : Team Member Insert Process

	fmt.Printf(fmt.Sprintf("CALL ZSP_TEAM_MEM_REG_PROC('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pPpChrgNm, pPpChrgBpNm, pPpChrgTelNo, pPpChrgSmsRecvYn, pPpChrgEmail, pPpChrgEmailRecvYn, pEntpId, sha_pwd))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_TEAM_MEM_REG_PROC('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pPpChrgNm, pPpChrgBpNm, pPpChrgTelNo, pPpChrgSmsRecvYn, pPpChrgEmail, pPpChrgEmailRecvYn, pEntpId, sha_pwd),
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

	rtnTeamMemberInsert := models.RtnTeamMemberInsert{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)

		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnTeamMemberInsert = models.RtnTeamMemberInsert{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}

	// End : Team Member Insert Process

	c.Data["json"] = &rtnTeamMemberInsert
	c.ServeJSON()
}
