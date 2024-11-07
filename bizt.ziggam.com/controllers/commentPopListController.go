package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type CommentPopListController struct {
	BaseController
}

func (c *CommentPopListController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/login")
	}

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := mem_no
	pRecrutSn := c.GetString("recrut_sn")
	pPpMemNo := c.GetString("pp_mem_no")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	fmt.Printf(fmt.Sprintf("CALL ZSP_COMMENT_LIST_R('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_COMMENT_LIST_R('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo),
		ora.I64, /* TOT_CNT */
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* PP_CHRG_CMT_SN */
		ora.S,   /* PP_CHRG_SN */
		ora.S,   /* PP_CHRG_CMT */
		ora.S,   /* REG_DT */
		ora.S,   /* REG_ID */
		ora.S,   /* PP_CHRG_BP_NM */
		ora.S,   /* PP_CHRG_NM */
		ora.S,   /* NEW_YN */
		ora.S,   /* PP_CHRG_GBN_CD */
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

	cmtSMemId := session.Get(c.Ctx.Request.Context(), "mem_id")
	cmtSAuthCd := session.Get(c.Ctx.Request.Context(), "auth_cd")

	recruitApplyCommentList := make([]models.RecruitApplyCommentList, 0)

	var (
		cmtTotCnt      int64
		cmtEntpMemNo   string
		cmtRecrutSn    string
		cmtPpMemNo     string
		cmtPpChrgCmtSn string
		cmtPpChrgSn    string
		cmtPpChrgCmt   string
		cmtRegDt       string
		cmtRegId       string
		cmtPpChrgBpNm  string
		cmtPpChrgNm    string
		cmtNewYn       string
		cmtPpChrgGbnCd string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			cmtTotCnt = procRset.Row[0].(int64)
			cmtEntpMemNo = procRset.Row[1].(string)
			cmtRecrutSn = procRset.Row[2].(string)
			cmtPpMemNo = procRset.Row[3].(string)
			cmtPpChrgCmtSn = procRset.Row[4].(string)
			cmtPpChrgSn = procRset.Row[5].(string)
			cmtPpChrgCmt = procRset.Row[6].(string)
			cmtRegDt = procRset.Row[7].(string)
			cmtRegId = procRset.Row[8].(string)
			cmtPpChrgBpNm = procRset.Row[9].(string)
			cmtPpChrgNm = procRset.Row[10].(string)
			cmtNewYn = procRset.Row[11].(string)
			cmtPpChrgGbnCd = procRset.Row[12].(string)

			recruitApplyCommentList = append(recruitApplyCommentList, models.RecruitApplyCommentList{
				CmtTotCnt:      cmtTotCnt,
				CmtEntpMemNo:   cmtEntpMemNo,
				CmtRecrutSn:    cmtRecrutSn,
				CmtPpMemNo:     cmtPpMemNo,
				CmtPpChrgCmtSn: cmtPpChrgCmtSn,
				CmtPpChrgSn:    cmtPpChrgSn,
				CmtPpChrgCmt:   cmtPpChrgCmt,
				CmtRegDt:       cmtRegDt,
				CmtRegId:       cmtRegId,
				CmtPpChrgBpNm:  cmtPpChrgBpNm,
				CmtPpChrgNm:    cmtPpChrgNm,
				CmtNewYn:       cmtNewYn,
				CmtPpChrgGbnCd: cmtPpChrgGbnCd,
				CmtSMemId:      cmtSMemId,
				CmtSAuthCd:     cmtSAuthCd,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Recruit Apply Comment List

	c.Data["RecrutSn"] = pRecrutSn
	c.Data["PpMemNo"] = pPpMemNo
	c.Data["CmtTotCnt"] = cmtTotCnt
	c.Data["RecruitApplyCommentList"] = recruitApplyCommentList

	c.Data["TMenuId"] = "M00"
	c.Data["SMenuId"] = "M00"

	c.TplName = "applicant/comment_pop_list.html"
}
