package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"

	"emst.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

type NoticeDetailController struct {
	BaseController
}

func (c *NoticeDetailController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/common/login")
	}

	pLang, _ := beego.AppConfig.String("lang")
	pSn := c.GetString("sn")
	pPageNo := c.GetString("pn")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Notice List
	log.Debug("CALL SP_EMS_NOTICE_DTL_R('%v', %v, :1)",
		pLang, pSn)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_NOTICE_DTL_R('%v', %v, :1)",
		pLang, pSn),
		ora.I64, /* SN */
		ora.S,   /* GBN_NM */
		ora.S,   /* TITLE */
		ora.S,   /* CONT */
		ora.S,   /* REG_DT */
		ora.S,   /* NEW_YN */
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

	noticeDetail := make([]models.NoticeDetail, 0)

	var (
		sn    int64
		memNm string
		title string
		regDt string
		cont  string
		newYn string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			sn = procRset.Row[0].(int64)
			memNm = procRset.Row[1].(string)
			title = procRset.Row[2].(string)
			regDt = procRset.Row[3].(string)
			cont = procRset.Row[4].(string)
			newYn = procRset.Row[5].(string)

			noticeDetail = append(noticeDetail, models.NoticeDetail{
				Sn:    sn,
				MemNm: memNm,
				Title: title,
				RegDt: regDt,
				Cont:  cont,
				NewYn: newYn,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Notice List

	c.Data["Sn"] = sn
	c.Data["MemNm"] = memNm
	c.Data["Title"] = title
	c.Data["RegDt"] = regDt
	c.Data["Cont"] = cont
	c.Data["NewYn"] = newYn
	c.Data["PageNo"] = pPageNo

	c.TplName = "notice/notice_detail.html"
}
