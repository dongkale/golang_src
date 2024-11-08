package controllers

import (
	"fmt"

	"emst.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	ora "gopkg.in/rana/ora.v4"
)

type AdminVersionInfoController struct {
	BaseController
}

func (c *AdminVersionInfoController) Get() {

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

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Version List
	log.Debug("CALL SP_EMS_ADMIN_APP_VER_INFO_R('%v' :1)",
		pLang)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_APP_VER_INFO_R('%v', :1)",
		pLang),
		ora.S,   /* APP_GBN */
		ora.I64, /* APP_VER_CD */
		ora.S,   /* FRC_UPT_YN */
		ora.S,   /* DSTB_SDT */
		ora.S,   /* APP_VER */
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

	versionList := make([]models.VersionList, 0)

	var (
		appGbn   string
		appVerCd int64
		frcUptYn string
		dstbSdt  string
		appVer   string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			appGbn = procRset.Row[0].(string)
			appVerCd = procRset.Row[1].(int64)
			frcUptYn = procRset.Row[2].(string)
			dstbSdt = procRset.Row[3].(string)
			appVer = procRset.Row[4].(string)

			versionList = append(versionList, models.VersionList{
				AppGbn:   appGbn,
				AppVerCd: appVerCd,
				FrcUptYn: frcUptYn,
				DstbSdt:  dstbSdt,
				AppVer:   appVer,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}

	// End : Version List

	// Start : Common Count99 List
	pGbnCd := "99"

	log.Debug("CALL SP_EMS_COUNT_LIST_R('%v', '%v', :1)",
		pLang, pGbnCd)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_COUNT_LIST_R('%v', '%v', :1)",
		pLang, pGbnCd),
		ora.S, /* COUNT_VAL */
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

	commonCount99List := make([]models.CommonCount99List, 0)

	var (
		count99 string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			count99 = procRset.Row[0].(string)

			commonCount99List = append(commonCount99List, models.CommonCount99List{
				Count99: count99,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Common Count99 List

	c.Data["VersionList"] = versionList
	c.Data["CommonCount99List"] = commonCount99List
	c.Data["MenuId"] = "04"
	c.Data["SubMenuId"] = "05"
	c.TplName = "admin/version_info.html"
}
