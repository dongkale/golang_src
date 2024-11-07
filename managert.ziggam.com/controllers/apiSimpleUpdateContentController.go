package controllers

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/astaxie/beego/logs"
	"gopkg.in/rana/ora.v4"
	"managert.ziggam.com/models"
)

type ApiSimpleUpdateContentController struct {
	BaseController
}

func (c *ApiSimpleUpdateContentController) Post() {
	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	pLang, _ := beego.AppConfig.String("lang")

	pBnrGrpSn := c.GetString("bnr_grp_sn")     // 배너그룹코드
	pUseYn := c.GetString("use_yn")     		// 사용여부
	pDeleteYn := c.GetString("del_yn")     		// 삭제여부

	log.Debug("/api/content/simple/update Param ('%v', '%v', '%v')",
		pBnrGrpSn, pUseYn, pDeleteYn)

	if pUseYn == "" {
		pUseYn = "-1"
	}

	if pDeleteYn == "" {
		pDeleteYn = "-1"
	}

	log.Debug("/api/content/simple/update Param ('%v', '%v', '%v')",
		pBnrGrpSn, pUseYn, pDeleteYn)

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	log.Debug("CALL MNG_SIMPLE_UPDATE_CONTENT('%v', '%v', '%v', '%v', :1)",
		pLang,
		pBnrGrpSn,
		pUseYn,
		pDeleteYn,
		)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL MNG_SIMPLE_UPDATE_CONTENT('%v', '%v', '%v', '%v', :1)",
		pLang,
		pBnrGrpSn,
		pUseYn,
		pDeleteYn,
	),
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
