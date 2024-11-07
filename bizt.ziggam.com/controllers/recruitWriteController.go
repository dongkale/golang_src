package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	"bizt.ziggam.com/tables"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type RecruitWriteController struct {
	BaseController
}

func (c *RecruitWriteController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	mem_id := session.Get(c.Ctx.Request.Context(), "mem_id")
	if mem_id == nil {
		c.Ctx.Redirect(302, "/login")
	}

	pLang, _ := beego.AppConfig.String("lang")

	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	pEntpMemNo := mem_no

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
	fmt.Printf(fmt.Sprintf("CALL ZSP_JOB_GRP_LIST_R('%v', '%v', '%v', :1)",
		pLang, pEmplTypCd, pUpJobGrpCd))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_JOB_GRP_LIST_R('%v', '%v', '%v', :1)",
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
	}
	// End : Job Group List

	// Start : Region Group List
	fmt.Printf(fmt.Sprintf("CALL ZSP_RGN_LIST_R_V2('%v', NULL, :1)", pLang))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_RGN_LIST_R_V2('%v', NULL, :1)", pLang),
		ora.S, /* RGN_GRP_CD */
		ora.S, /* RGN_GRP_NM */
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

	rgnGrpList := make([]models.RgnGrp, 0)

	var (
		rgnGrpCd  string
		rgnGrpNm  string
		rgnGrpFNm string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			rgnGrpCd = procRset.Row[0].(string)
			rgnGrpNm = procRset.Row[1].(string)
			rgnGrpFNm = procRset.Row[2].(string)

			rgnGrpList = append(rgnGrpList, models.RgnGrp{
				Code:     rgnGrpCd,
				Name:     rgnGrpNm,
				FullName: rgnGrpFNm,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Job Group List

	// Start : Jobfair List
	fmt.Printf(fmt.Sprintf("CALL ZSP_ENTP_JOBFAIR_LIST('%v', '%v', :1)", pLang, pEntpMemNo))
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

	jobFailrInfoList := make([]models.JobfairInfo, 0)

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

			jobFailrInfoList = append(jobFailrInfoList, models.JobfairInfo{
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

	// LDK 2020/08/24 : 채용 정보 코드화
	c.Data["JobGrpList"] = jobGrpList

	c.Data["RgnGrpList"] = rgnGrpList

	c.Data["JobFairList"] = jobFailrInfoList

	c.Data["CarrGbnMap"] = tables.MapCarrGbnCd
	c.Data["EmplTypMap"] = tables.MapEmplTypCd
	c.Data["LstEduGbnMap"] = tables.MapLstEduGbnCd
	// <--

	// LDK 2020/10/07: 채용 질문 갯수 조정
	// recruit_copy.html
	// recruit_modify.html
	// recruit_write.html
	c.Data["EntpMemNo"] = pEntpMemNo
	// <--

	c.Data["TMenuId"] = "R00"
	c.Data["SMenuId"] = "R00"

	c.TplName = "recruit/recruit_write.html"
}
