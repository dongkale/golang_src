package controllers

import (
	"fmt"

	"emst.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	ora "gopkg.in/rana/ora.v4"
)

type JobfairEntpDeleteController struct {
	beego.Controller
}

func (c *JobfairEntpDeleteController) Post() {

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Data["json"] = &models.DefaultResult{RtnCd: 99, RtnMsg: "/login"}
		c.ServeJSON()
		return
	}

	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := c.GetString("entp_mem_no")   // 기업회원번호
	pJobfairArr := c.GetString("job_fair_arr") // 박람회 번호 배열 리스트

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Jobfair Entp Delete
	logs.Debug(fmt.Sprintf("CALL SP_JOBFAIR_ENTP_DEL('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pJobfairArr))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_JOBFAIR_ENTP_DEL('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pJobfairArr),
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

	defaultResult := models.DefaultResult{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}

		if err := procRset.Err(); err != nil {
			panic(err)
		}

		defaultResult = models.DefaultResult{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}

		logs.Debug(fmt.Sprintf(" ===> rtnCd:%v, rtnMsg:%v", rtnCd, rtnMsg))
	}
	// End : jobfair Entp Delete

	// Start : Jobfair List
	logs.Debug(fmt.Sprintf("CALL ZSP_ENTP_JOBFAIR_LIST('%v', '%v', :1)", pLang, pEntpMemNo))
	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_ENTP_JOBFAIR_LIST('%v', '%v', :1)",
		pLang, pEntpMemNo),
		ora.S, /* MNG_CD */
		ora.S, /* TITLE */
		ora.S, /* SDY */
		ora.S, /* EDY */
		ora.S, /* HOST_INSTITUTION */
		ora.S, /* MANAGE_AGENCY */
		ora.S, /* AGREEMENT_URL */
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

	jobFairInfoList := make([]models.JobfairInfo, 0)

	var (
		jfMngCd           string
		jfTitle           string
		jfSdy             string
		jfEdy             string
		jfHostInstitution string
		jfManageAgency    string
	)

	if procRset.IsOpen() {
		for procRset.Next() {

			jfMngCd = procRset.Row[0].(string)
			jfTitle = procRset.Row[1].(string)
			jfSdy = procRset.Row[2].(string)
			jfEdy = procRset.Row[3].(string)
			jfHostInstitution = procRset.Row[4].(string)
			jfManageAgency = procRset.Row[5].(string)

			jobFairInfoList = append(jobFairInfoList, models.JobfairInfo{
				MngCd:           jfMngCd,
				Title:           jfTitle,
				Sdy:             jfSdy,
				Edy:             jfEdy,
				HostInstitution: jfHostInstitution,
				ManageAgency:    jfManageAgency,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Jobfair List

	c.Data["json"] = &models.JobFairEntpDelete{RslData: defaultResult, RslJobFairList: jobFairInfoList}
	c.ServeJSON()
}
