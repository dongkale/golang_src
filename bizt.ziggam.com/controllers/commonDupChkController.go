package controllers

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	ora "gopkg.in/rana/ora.v4"
)

type CommonDupChkController struct {
	beego.Controller
}

func (c *CommonDupChkController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	// start : session
	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	mem_sn := session.Get(c.Ctx.Request.Context(), "mem_sn")

	pLang, _ := beego.AppConfig.String("lang")
	// end : session

	pGbnCD := c.GetString("gbn_cd")     //중복구분코드
	pItemVal := c.GetString("item_val") //중복구분값
	pEntpMemNo := mem_no                //기업회원번호
	pPpChrgSn := mem_sn                 //담당자순번

	tmpPpChrgSn := c.GetString("pp_chrg_sn")
	if tmpPpChrgSn != "" {
		pPpChrgSn = tmpPpChrgSn
	}

	if pEntpMemNo == nil {
		pEntpMemNo = ""
	}

	if pPpChrgSn == nil {
		pPpChrgSn = ""
	}

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	fmt.Printf(fmt.Sprintf("CALL ZSP_DUP_CHK_R( '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pGbnCD, pItemVal, pEntpMemNo, pPpChrgSn))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_DUP_CHK_R( '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pGbnCD, pItemVal, pEntpMemNo, pPpChrgSn),
		ora.S, /* DUP_YN */
	)
	defer stmtProcCall.Close()
	if err != nil {
		panic(err)
	}
	procRset := &ora.Rset{}
	rowsAffected, err := stmtProcCall.Exe(procRset)
	fmt.Println(rowsAffected)

	if err != nil {
		panic(err)
	}

	var chkFlagYn string

	if procRset.IsOpen() {
		for procRset.Next() {
			chkFlagYn = procRset.Row[0].(string)
		}
	}
	if err := procRset.Err(); err != nil {
		panic(err)
	}

	var status bool

	if chkFlagYn == "N" {
		status = true
	} else if chkFlagYn == "T" {
		c.Data["json"] = "탈퇴한 회원입니다. 다른 id로 만들어 주세요"
		c.ServeJSON()
	} else {
		status = false
	}

	//c.Data["json"] = &chkFlagYn
	c.Data["json"] = &status
	c.ServeJSON()

}
