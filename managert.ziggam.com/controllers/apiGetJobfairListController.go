package controllers

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/rana/ora.v4"
	"managert.ziggam.com/models"
)

type ApiGetJobfairListController struct {
	BaseController
}

func (c *ApiGetJobfairListController) Get() {
	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

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

	// Start : Applicant Delete Process

	log.Debug("CALL MNG_LIST_JOBFAIR_INFO('%v', :1)",
		pLang)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL MNG_LIST_JOBFAIR_INFO('%v', :1)",
		pLang),
		ora.S, /* MNG_CD */
		ora.S, /* TITLE */
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
		mngCd string
		title string
	)

	//rtnGroupBannerList := models.GroupBannerDataList{}
	jobfairInfoList := make([]models.JobfairInfo, 0)

	if procRset.IsOpen() {
		for procRset.Next() {
			mngCd = procRset.Row[0].(string)
			title = procRset.Row[1].(string)

			jobfairInfoList = append(jobfairInfoList, models.JobfairInfo{
				MngCd: mngCd,
				Title: title,
			})
		}

		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}

	c.Data["json"] = &jobfairInfoList
	c.ServeJSON()

}
