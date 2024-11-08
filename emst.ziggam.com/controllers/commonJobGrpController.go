package controllers

import (
	"fmt"

	"emst.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	"gopkg.in/rana/ora.v4"
)

type CommonJobGrpController struct {
	beego.Controller
}

func (c *CommonJobGrpController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	//session := c.StartSession()
	pLang, _ := beego.AppConfig.String("lang")

	pEmplTypCd := c.GetString("empl_typ_cd")
	pUpJobGrpCd := c.GetString("up_job_grp_cd")

	if pEmplTypCd == "" {
		pEmplTypCd = "01"
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

	// Start : Job Group List
	log.Debug("CALL SP_EMS_JOB_GRP_LIST_R('%v', '%v', '%v', :1)",
		pLang, pEmplTypCd, pUpJobGrpCd)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_JOB_GRP_LIST_R('%v', '%v', '%v', :1)",
		pLang, pEmplTypCd, pUpJobGrpCd),
		ora.S, /* JOB_GRP_CD */
		ora.S, /* UP_JOB_GRP_CD */
		ora.S, /* JOB_GRP_NM */

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

	rtnJobGrpList := models.RtnJobGrpList{}
	jobGrpList := make([]models.JobGrpList, 0)

	var (
		jobGrpCd   string
		upJobGrpCd string
		jobGrpNm   string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			jobGrpCd = procRset.Row[0].(string)
			upJobGrpCd = procRset.Row[1].(string)
			jobGrpNm = procRset.Row[2].(string)

			jobGrpList = append(jobGrpList, models.JobGrpList{
				JobGrpCd:   jobGrpCd,
				UpJobGrpCd: upJobGrpCd,
				JobGrpNm:   jobGrpNm,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnJobGrpList = models.RtnJobGrpList{
			RtnJobGrpListData: jobGrpList,
		}
	}
	// End : Job Group List

	c.Data["json"] = &rtnJobGrpList
	c.ServeJSON()
}
