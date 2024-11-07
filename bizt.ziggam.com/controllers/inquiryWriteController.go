package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type InquiryWriteController struct {
	BaseController
}

func (c *InquiryWriteController) Get() {

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
	pGbn := c.GetString("gbn")
	pCdGrpId := "G019"

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Common Code List
	fmt.Printf(fmt.Sprintf("CALL ZSP_COMMON_CD_LIST_R('%v', '%v', :1)",
		pLang, pCdGrpId))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_COMMON_CD_LIST_R('%v', '%v', :1)",
		pLang, pCdGrpId),
		ora.S, /* CD_ID */
		ora.S, /* CD_NM */
		ora.S, /* REMRK */

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

	inqCommonCdList := make([]models.InqCommonCdList, 0)

	var (
		iqCdId  string
		iqCdNm  string
		iqRemrk string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			iqCdId = procRset.Row[0].(string)
			iqCdNm = procRset.Row[1].(string)
			iqRemrk = procRset.Row[2].(string)

			inqCommonCdList = append(inqCommonCdList, models.InqCommonCdList{
				IqCdId:  iqCdId,
				IqCdNm:  iqCdNm,
				IqRemrk: iqRemrk,
				RefGbn:  pGbn,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Common Code List

	c.Data["InqCommonCdList"] = inqCommonCdList

	c.Data["EntpMemNo"] = pEntpMemNo
	c.Data["RefGbn"] = pGbn
	c.Data["TMenuId"] = "M00"
	c.Data["SMenuId"] = "M00"

	c.TplName = "setting/inquiry_write.html"
}
