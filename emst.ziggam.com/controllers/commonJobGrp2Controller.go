package controllers

import (
	"fmt"

	"emst.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	"gopkg.in/rana/ora.v4"
)

type CommonJobGrp2Controller struct {
	beego.Controller
}

func (c *CommonJobGrp2Controller) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	//session := c.StartSession()
	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := c.GetString("entp_mem_no")
	pEmplTypCd := c.GetString("empl_typ_cd")

	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Job Group List
	log.Debug("CALL SP_EMS_RECRUIT_JOB_GRP_LIST_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pEmplTypCd)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_RECRUIT_JOB_GRP_LIST_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pEmplTypCd),
		ora.S, /* EMPL_TYP_CD */
		ora.S, /* JOB_GRP_CD */
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

	rtnRecruitJobGrpList := models.RtnRecruitJobGrpList{}
	recruitJobGrpList := make([]models.RecruitJobGrpList, 0)

	var (
		rEmplTypCd string
		rJobGrpCd  string
		rJobGrpNm  string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			rEmplTypCd = procRset.Row[0].(string)
			rJobGrpCd = procRset.Row[1].(string)
			rJobGrpNm = procRset.Row[2].(string)

			recruitJobGrpList = append(recruitJobGrpList, models.RecruitJobGrpList{
				REmplTypCd: rEmplTypCd,
				RJobGrpCd:  rJobGrpCd,
				RJobGrpNm:  rJobGrpNm,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnRecruitJobGrpList = models.RtnRecruitJobGrpList{
			RtnRecruitJobGrpListData: recruitJobGrpList,
		}

		c.Data["json"] = &rtnRecruitJobGrpList
		c.ServeJSON()
	}
	// End : Job Group List
}
