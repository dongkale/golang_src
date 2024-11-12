package controllers

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/rana/ora.v4"
	"managert.ziggam.com/models"
)

type ApiUpdateContentClListController struct {
	BaseController
}

func (c *ApiUpdateContentClListController) Post() {
	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	pLang, _ := beego.AppConfig.String("lang")

	pBnrGrpSubSn := c.GetString("bnr_grp_sn")
	pEntpMemNo := c.GetString("entp_mem_no")
	pFlag := c.GetString("flag")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	log.Debug("CALL MNG_UPDATE_CL_GROUP('%v', '%v', '%v', '%v', :1)",
		pLang,
		pBnrGrpSubSn,
		pEntpMemNo,
		pFlag,
	)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL MNG_UPDATE_CL_GROUP('%v', '%v', '%v', '%v', :1)",
		pLang,
		pBnrGrpSubSn,
		pEntpMemNo,
		pFlag),
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

	rtnResult := models.RtnResult{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnResult = models.RtnResult{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}

	c.Data["json"] = rtnResult
	c.ServeJSON()
}
