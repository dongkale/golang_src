package controllers

import (
	"fmt"

	"emst.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	ora "gopkg.in/rana/ora.v4"
)

type AdminEventWriteController struct {
	BaseController
}

func (c *AdminEventWriteController) Get() {

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
	if pSn == "" {
		pSn = "0"
	}
	pPageNo := c.GetString("pn")
	pCuCd := c.GetString("cu_cd")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Event Write
	log.Debug("CALL SP_EMS_EVENT_DTL_R('%v', %v, :1)",
		pLang, pSn)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_EVENT_DTL_R('%v', %v, :1)",
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

	eventDetail := make([]models.EventDetail, 0)

	var (
		sn    int64
		gbnNm string
		title string
		regDt string
		cont  string
		newYn string
		gbnCd string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			sn = procRset.Row[0].(int64)
			gbnNm = procRset.Row[1].(string)
			title = procRset.Row[2].(string)
			regDt = procRset.Row[3].(string)
			cont = procRset.Row[4].(string)
			newYn = procRset.Row[5].(string)
			gbnCd = procRset.Row[6].(string)

			eventDetail = append(eventDetail, models.EventDetail{
				Sn:    sn,
				GbnNm: gbnNm,
				Title: title,
				RegDt: regDt,
				Cont:  cont,
				NewYn: newYn,
				GbnCd: gbnCd,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Event Write

	c.Data["CuCd"] = pCuCd
	c.Data["Sn"] = sn
	c.Data["GbnNm"] = gbnNm
	c.Data["Title"] = title
	c.Data["RegDt"] = regDt
	c.Data["Cont"] = cont
	c.Data["NewYn"] = newYn
	c.Data["GbnCd"] = gbnCd
	c.Data["PageNo"] = pPageNo
	c.Data["MenuId"] = "04"
	c.Data["SubMenuId"] = "02"
	c.TplName = "admin/event_write.html"
}
