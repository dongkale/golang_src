package controllers

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/rana/ora.v4"
	"managert.ziggam.com/models"
)

type ApiSaveContentController struct {
	BaseController
}

func (c *ApiSaveContentController) Post() {
	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	pLang, _ := beego.AppConfig.String("lang")

	content := models.ContentItem{}
	c.Ctx.Input.Bind(&content, "params")

	// 날짜 형식이 맞는지 검사 하자. content.PublSdy
	// content.PublSdy += "00";
	// content.PublEdy += "59";

	log.Debug("/api/content/save Param ('%v')",
		content)

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	log.Debug("CALL MNG_SAVE_CONTENT('%v', '%v', '%v', %v, '%v', '%v', '%v', '%v', '%v', '%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', %v, '%v', :1)",
		pLang,
		content.BnrGrpSn,
		content.BnrGrpTypCd,
		content.BnrGrpIdx,
		content.BnrGrpTitle,
		content.BnrGrpSubCn,
		content.PtoPath,
		content.ThumbPtoPath,
		content.PublSdy,
		content.PublEdy,
		content.RdCnt,
		content.ShCnt,
		content.UseYn,
		content.DelYn,
		content.RegDt,
		content.RegId,
		content.UptDt,
		content.UptId,
		content.Expln,
		content.RolTm,
		content.EdtFlg)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL MNG_SAVE_CONTENT('%v', '%v', '%v', %v, '%v', '%v', '%v', '%v', '%v', '%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', %v, '%v', :1)",
		pLang,
		content.BnrGrpSn,
		content.BnrGrpTypCd,
		content.BnrGrpIdx,
		content.BnrGrpTitle,
		content.BnrGrpSubCn,
		content.PtoPath,
		content.ThumbPtoPath,
		content.PublSdy,
		content.PublEdy,
		content.RdCnt,
		content.ShCnt,
		content.UseYn,
		content.DelYn,
		content.RegDt,
		content.RegId,
		content.UptDt,
		content.UptId,
		content.Expln,
		content.RolTm,
		content.EdtFlg),
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
