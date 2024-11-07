package controllers

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"bizt.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

type TeamMemberListController struct {
	BaseController
}

func (c *TeamMemberListController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	mem_id := session.Get(c.Ctx.Request.Context(), "mem_id")
	if mem_id == nil {
		c.Ctx.Redirect(302, "/login")
	}
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	pEntpMemNo := mem_no
	//auth_cd := session.Get(c.Ctx.Request.Context(), "auth_cd")
	//sAuthCd := auth_cd
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

	// Start : Team Member Admin Info
	fmt.Printf(fmt.Sprintf("CALL ZSP_TEAM_ADMIN_R('%v', '%v', :1)",
		pLang, pEntpMemNo))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_TEAM_ADMIN_R('%v', '%v', :1)",
		pLang, pEntpMemNo),
		ora.S, /* PP_CHRG_SN */
		ora.S, /* PP_CHRG_NM */
		ora.S, /* PP_CHRG_BP_NM */
		ora.S, /* EMAIL */
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

	teamMemberAdmin := make([]models.TeamMemberAdmin, 0)

	var (
		taPpChrgSn   string
		taPpChrgNm   string
		taPpChrgBpNm string
		taEmail      string
		taEntpMemId  string
		taRegDt      string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			taPpChrgSn = procRset.Row[0].(string)
			taPpChrgNm = procRset.Row[1].(string)
			taPpChrgBpNm = procRset.Row[2].(string)
			taEmail = procRset.Row[3].(string)
			taEntpMemId = procRset.Row[4].(string)
			taRegDt = procRset.Row[5].(string)

			teamMemberAdmin = append(teamMemberAdmin, models.TeamMemberAdmin{
				TaPpChrgSn:   taPpChrgSn,
				TaPpChrgNm:   taPpChrgNm,
				TaPpChrgBpNm: taPpChrgBpNm,
				TaEmail:      taEmail,
				TaEntpMemId:  taEntpMemId,
				TaRegDt:      taRegDt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Team Member Admin Info

	// Start : Entp Team Member List

	pGbnCd := "M"

	fmt.Printf(fmt.Sprintf("CALL ZSP_TEAM_MEM_LIST_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pGbnCd))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_TEAM_MEM_LIST_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pGbnCd),
		ora.I64, /* TOT_CNT */
		ora.S,   /* PP_CHRG_SN */
		ora.S,   /* PP_CHRG_GBN_CD */
		ora.S,   /* PP_CHRG_NM */
		ora.S,   /* PP_CHRG_BP_NM */
		ora.S,   /* EMAIL */
		ora.S,   /* ENTP_MEM_ID */
		ora.S,   /* PP_CHRG_TEL_NO */
		ora.I64, /* ROWNO */
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

	entpTeamMemberList := make([]models.EntpTeamMemberList, 0)

	var (
		etTotCnt      int64
		etPpChrgSn    string
		etPpChrgGbnCd string
		etPpChrgNm    string
		etPpChrgBpNm  string
		etEmail       string
		etEntpMemId   string
		etPpChrgTelNo string
		etRowNo       int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			etTotCnt = procRset.Row[0].(int64)
			etPpChrgSn = procRset.Row[1].(string)
			etPpChrgGbnCd = procRset.Row[2].(string)
			etPpChrgNm = procRset.Row[3].(string)
			etPpChrgBpNm = procRset.Row[4].(string)
			etEmail = procRset.Row[5].(string)
			etEntpMemId = procRset.Row[6].(string)
			etPpChrgTelNo = procRset.Row[7].(string)
			etRowNo = procRset.Row[8].(int64)

			entpTeamMemberList = append(entpTeamMemberList, models.EntpTeamMemberList{
				EtTotCnt:      etTotCnt,
				EtPpChrgSn:    etPpChrgSn,
				EtPpChrgGbnCd: etPpChrgGbnCd,
				EtPpChrgNm:    etPpChrgNm,
				EtPpChrgBpNm:  etPpChrgBpNm,
				EtEmail:       etEmail,
				EtEntpMemId:   etEntpMemId,
				EtPpChrgTelNo: etPpChrgTelNo,
				EtRowNo:       etRowNo,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Entp Team Member List

	c.Data["EntpTeamMemberList"] = entpTeamMemberList

	c.Data["TaPpChrgSn"] = taPpChrgSn
	c.Data["TaPpChrgNm"] = taPpChrgNm
	c.Data["TaPpChrgBpNm"] = taPpChrgBpNm
	c.Data["TaEmail"] = taEmail
	c.Data["TaEntpMemId"] = taEntpMemId
	c.Data["TaRegDt"] = taRegDt

	c.Data["TMenuId"] = "T00"
	c.Data["SMenuId"] = "T00"

	c.TplName = "team/team_member_list.html"
}
