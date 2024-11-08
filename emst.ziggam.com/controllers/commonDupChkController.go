package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"

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
	//session := c.StartSession()
	pLang, _ := beego.AppConfig.String("lang")
	// end : session

	pGbnCD := c.GetString("gbn_cd")          //중복구분코드
	pItemVal := c.GetString("item_val")      //중복구분값
	pEntpMemNo := c.GetString("entp_mem_no") //기업회원번호

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	log.Debug("CALL SP_EMS_DUP_CHK_R( '%v', '%v', '%v', '%v', :1)",
		pLang, pGbnCD, pItemVal, pEntpMemNo)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_DUP_CHK_R( '%v', '%v', '%v', '%v', :1)",
		pLang, pGbnCD, pItemVal, pEntpMemNo),
		ora.S, /* JOIN_YN */
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
	} else {
		status = false
	}

	//c.Data["json"] = &chkFlagYn
	c.Data["json"] = &status
	c.ServeJSON()

}
