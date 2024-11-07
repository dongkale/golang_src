package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type InviteRecruitListController struct {
	BaseController
}

func (c *InviteRecruitListController) Post() {

	session := c.StartSession()

	// mem_id := session.Get(c.Ctx.Request.Context(), "mem_id")
	// if mem_id == nil {
	// 	c.Ctx.Redirect(302, "/login")
	// }

	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		//c.Ctx.Redirect(302, "/login")

		c.Data["json"] = &models.RtnRecruitSubList{}
		c.ServeJSON()
		return
	}

	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := mem_no
	pGbnCd := c.GetString("recruit_cd") // 구분코드(A:전체, I:채용중, E:종료)

	if pGbnCd == "" {
		pGbnCd = "A"
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

	pOffSet := "0"
	pLimit := "100"
	pSortGbn := "03"

	fmt.Printf(fmt.Sprintf("CALL ZSP_INVITE_RECRUIT_LIST_TYPE('%v', %v, %v, '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pGbnCd, pSortGbn))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_INVITE_RECRUIT_LIST_TYPE('%v', %v, %v, '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pGbnCd, pSortGbn),
		ora.I64, /* TOT_CNT */
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* PRGS_STAT */
		ora.S,   /* RECRUT_TITLE */
		ora.S,   /* UP_JOB_GRP */
		ora.S,   /* JOB_GRP */
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

	rtnRecruitSubList := models.RtnRecruitSubList{}
	recruitSubList := make([]models.RecruitSubList, 0)

	var (
		sTotCnt      int64
		sEntpMemNo   string
		sRecrutSn    string
		sPrgsStat    string
		sRecrutTitle string
		sUpJobGrp    string
		sJobGrp      string
		//sTrimRecrutTitle string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			sTotCnt = procRset.Row[0].(int64)
			sEntpMemNo = procRset.Row[1].(string)
			sRecrutSn = procRset.Row[2].(string)
			sPrgsStat = procRset.Row[3].(string)
			sRecrutTitle = procRset.Row[4].(string)
			sUpJobGrp = procRset.Row[5].(string)
			sJobGrp = procRset.Row[6].(string)

			recruitSubList = append(recruitSubList, models.RecruitSubList{
				STotCnt:      sTotCnt,
				SEntpMemNo:   sEntpMemNo,
				SRecrutSn:    sRecrutSn,
				SPrgsStat:    sPrgsStat,
				SRecrutTitle: sRecrutTitle,
				SUpJobGrp:    sUpJobGrp,
				SJobGrp:      sJobGrp,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnRecruitSubList = models.RtnRecruitSubList{
			RtnRecruitSubListData: recruitSubList,
		}
	}

	c.Data["json"] = &rtnRecruitSubList
	c.ServeJSON()
}
