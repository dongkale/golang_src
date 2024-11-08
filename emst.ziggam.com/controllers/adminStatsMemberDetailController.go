package controllers

import (
	"fmt"
	"time"

	"emst.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	ora "gopkg.in/rana/ora.v4"
)

type AdminStatsMemberDetailController struct {
	BaseController
}

func (c *AdminStatsMemberDetailController) Get() {

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
	pSdt := c.GetString("sdt")
	pEdt := c.GetString("edt")
	pTGbn := c.GetString("t_gbn")
	pStatCd := c.GetString("stat_cd")
	pGbnVal := c.GetString("gbn_val")
	pMemGbn := c.GetString("mem_gbn")

	nowDate := time.Now()
	dateFmt := fmt.Sprintf(nowDate.Format("20060102150405"))

	if pSdt == "" {
		pSdt = dateFmt[0:8]
	}

	if pEdt == "" {
		pEdt = dateFmt[0:8]
	}

	if pTGbn == "" {
		pTGbn = "M"
	}

	if pStatCd == "" {
		pStatCd = "01"
	}

	if pGbnVal == "" {
		pGbnVal = "0"
	}

	if pMemGbn == "" {
		pMemGbn = "p"
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

	// Start : Admin Member Stats List
	log.Debug("CALL SP_EMS_ADMIN_STATS_MEMBER_R('%v', '%v', '%v', '%v', '%v',  %v, '%v', :1)",
		pLang, pSdt, pEdt, pTGbn, pStatCd, pGbnVal, pMemGbn)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_STATS_MEMBER_R('%v', '%v', '%v', '%v', '%v',  %v, '%v', :1)",
		pLang, pSdt, pEdt, pTGbn, pStatCd, pGbnVal, pMemGbn),
		ora.S,   /* ANAL_DT */
		ora.I64, /* ANAL_CNT */
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

	adminStatsMemberDetail := make([]models.AdminStatsMemberDetail, 0)

	var (
		analDt  string
		analCnt int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			analDt = procRset.Row[0].(string)
			analCnt = procRset.Row[1].(int64)

			adminStatsMemberDetail = append(adminStatsMemberDetail, models.AdminStatsMemberDetail{
				AnalDt:  analDt,
				AnalCnt: analCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Admin Member Stats List

	// Start : Admin Member Stats Top
	log.Debug("CALL SP_EMS_ADMIN_STATS_MEM_TOP_R('%v', '%v', :1)",
		pLang, pMemGbn)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_STATS_MEM_TOP_R('%v', '%v', :1)",
		pLang, pMemGbn),
		ora.I64, /* COM_CNT */
		ora.I64, /* WTD_CNT */
		ora.I64, /* UVF_CNT */
	)
	defer stmtProcCall.Close()
	if err != nil {
		panic(err)
	}
	procRset = &ora.Rset{}
	_, err = stmtProcCall.Exe(procRset)

	if err != nil {
		panic(err)
	}

	adminStatsMemberTopCnt := make([]models.AdminStatsMemberTopCnt, 0)

	var (
		comCnt int64
		wtdCnt int64
		uvfCnt int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			comCnt = procRset.Row[0].(int64)
			wtdCnt = procRset.Row[1].(int64)
			uvfCnt = procRset.Row[2].(int64)

			adminStatsMemberTopCnt = append(adminStatsMemberTopCnt, models.AdminStatsMemberTopCnt{
				ComCnt: comCnt,
				WtdCnt: wtdCnt,
				UvfCnt: uvfCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Admin Member Stats Top

	c.Data["ComCnt"] = comCnt
	c.Data["WtdCnt"] = wtdCnt
	c.Data["UvfCnt"] = uvfCnt
	c.Data["AdminStatsMemberDetail"] = adminStatsMemberDetail
	c.Data["StatCd"] = pStatCd
	c.Data["MemGbn"] = pMemGbn
	c.Data["MenuId"] = "99"

	c.TplName = "admin/stats_member_detail.html"
}

func (c *AdminStatsMemberDetailController) Post() {

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
	pSdt := c.GetString("sdt")
	pEdt := c.GetString("edt")
	pTGbn := c.GetString("t_gbn")
	pStatCd := c.GetString("stat_cd")
	pGbnVal := c.GetString("gbn_val")
	pMemGbn := c.GetString("mem_gbn")

	nowDate := time.Now()
	dateFmt := fmt.Sprintf(nowDate.Format("20060102150405"))

	if pSdt == "" {
		pSdt = dateFmt[0:8]
	}

	if pEdt == "" {
		pEdt = dateFmt[0:8]
	}

	if pTGbn == "" {
		pTGbn = "M"
	}

	if pStatCd == "" {
		pStatCd = "01"
	}

	if pGbnVal == "" {
		pGbnVal = "0"
	}

	if pMemGbn == "" {
		pMemGbn = "P"
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

	// Start : Admin Member Stats List
	log.Debug("CALL SP_EMS_ADMIN_STATS_MEMBER_R('%v', '%v', '%v', '%v', '%v',  %v, '%v', :1)",
		pLang, pSdt, pEdt, pTGbn, pStatCd, pGbnVal, pMemGbn)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_STATS_MEMBER_R('%v', '%v', '%v', '%v', '%v',  %v, '%v', :1)",
		pLang, pSdt, pEdt, pTGbn, pStatCd, pGbnVal, pMemGbn),
		ora.S,   /* ANAL_DT */
		ora.I64, /* ANAL_CNT */
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

	rtnAdminStatsMemberDetail := models.RtnAdminStatsMemberDetail{}
	adminStatsMemberDetail := make([]models.AdminStatsMemberDetail, 0)

	var (
		analDt  string
		analCnt int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			analDt = procRset.Row[0].(string)
			analCnt = procRset.Row[1].(int64)

			adminStatsMemberDetail = append(adminStatsMemberDetail, models.AdminStatsMemberDetail{
				AnalDt:  analDt,
				AnalCnt: analCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnAdminStatsMemberDetail = models.RtnAdminStatsMemberDetail{
			RtnAdminStatsMemberDetailData: adminStatsMemberDetail,
		}
		// End : Admin Member Stats List

		c.Data["json"] = &rtnAdminStatsMemberDetail
		c.ServeJSON()
	}
}
