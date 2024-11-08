package controllers

import (
	"fmt"

	"emst.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	"gopkg.in/rana/ora.v4"
)

type CommonCodeListController struct {
	beego.Controller
}

func (c *CommonCodeListController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	//session := c.StartSession()
	pLang, _ := beego.AppConfig.String("lang")

	pCdGrpId := c.GetString("cd_grp_id")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Group List
	log.Debug("CALL SP_EMS_COMMON_CD_LIST_R('%v', '%v', :1)",
		pLang, pCdGrpId)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_COMMON_CD_LIST_R('%v', '%v', :1)",
		pLang, pCdGrpId),
		ora.S, /* CD_ID */
		ora.S, /* CD_NM */

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

	rtnCommonCodeList := models.RtnCommonCodeList{}
	commonCodeList := make([]models.CommonCodeList, 0)

	var (
		cdId string
		cdNm string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			cdId = procRset.Row[0].(string)
			cdNm = procRset.Row[1].(string)

			commonCodeList = append(commonCodeList, models.CommonCodeList{
				CdId: cdId,
				CdNm: cdNm,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnCommonCodeList = models.RtnCommonCodeList{
			RtnCommonCodeListData: commonCodeList,
		}
	}
	// End : Group List

	c.Data["json"] = &rtnCommonCodeList
	c.ServeJSON()
}
