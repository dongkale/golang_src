package controllers

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/astaxie/beego/logs"
	"gopkg.in/rana/ora.v4"
	"managert.ziggam.com/models"
)

type ApiUpdateContentClUseYnController struct {
	BaseController
}

func (c *ApiUpdateContentClUseYnController) Post() {
	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	pLang, _ := beego.AppConfig.String("lang")

	pBnrGrpSubSn := c.GetString("bnr_grp_sub_sn")        // 배너SN
	pUseYn := c.GetString("use_yn") // 타이틀명

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	log.Debug("CALL MNG_UPDATE_CONTENT_CL_USE_YN('%v', '%v', %v, :1)",
		pLang,
		pBnrGrpSubSn,
		pUseYn,
	)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL MNG_UPDATE_CONTENT_CL_USE_YN('%v', '%v', %v, :1)",
		pLang,
		pBnrGrpSubSn,
		pUseYn),
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
