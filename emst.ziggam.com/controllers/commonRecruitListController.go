package controllers

import (
	"fmt"

	"emst.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	"gopkg.in/rana/ora.v4"
)

type CommonRecruitListController struct {
	beego.Controller
}

func (c *CommonRecruitListController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	//session := c.StartSession()
	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := c.GetString("entp_mem_no")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Recruit List
	log.Debug("CALL SP_EMS_ADMIN_ITEM2_LIST_R('%v', '%v', :1)",
		pLang, pEntpMemNo)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_ITEM2_LIST_R('%v', '%v', :1)",
		pLang, pEntpMemNo),
		ora.S, /* RECRUT_SN */
		ora.S, /* RECRUT_TITLE */

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

	rtnCommonRecruitList := models.RtnCommonRecruitList{}
	commonRecruitList := make([]models.CommonRecruitList, 0)

	var (
		recrutSn    string
		recrutTitle string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			recrutSn = procRset.Row[0].(string)
			recrutTitle = procRset.Row[1].(string)

			commonRecruitList = append(commonRecruitList, models.CommonRecruitList{
				RecrutSn:    recrutSn,
				RecrutTitle: recrutTitle,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnCommonRecruitList = models.RtnCommonRecruitList{
			RtnCommonRecruitListData: commonRecruitList,
		}
	}
	// End : Recruit List

	c.Data["json"] = &rtnCommonRecruitList
	c.ServeJSON()
}
