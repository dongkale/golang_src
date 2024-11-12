package controllers

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/rana/ora.v4"
	"managert.ziggam.com/models"
)

type ApiGetContentClAddListController struct {
	BaseController
}

func (c *ApiGetContentClAddListController) Get() {
	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	pLang, _ := beego.AppConfig.String("lang")

	pBnrGrpSn := c.GetString("bnr_gpr_sn")
	pJobfairCode := c.GetString("jobfair_code")

	log.Debug("Jobfair Code('%v')", pJobfairCode)

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

	log.Debug("CALL MNG_LIST_CL_ADD_LIST_INFO('%v', '%v', :1)",
		pLang, pBnrGrpSn)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL MNG_LIST_CL_ADD_LIST_INFO('%v', '%v', :1)",
		pLang, pBnrGrpSn),
		ora.S, /* BNR_GRP_SUB_SN */
		ora.S, /* ENTP_MEM_NO */
		ora.S, /* ENTP_KO_NM */
		ora.S, /* REG_DT */
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
		bnrGrpSubSn  string
		entpMemNo string
		entpKoNm  string
		regDt     string
	)

	//rtnGroupBannerList := models.GroupBannerDataList{}
	contentClAddItemList := make([]models.ContentClAddItem, 0)

	if procRset.IsOpen() {
		for procRset.Next() {
			bnrGrpSubSn = procRset.Row[0].(string)
			entpMemNo = procRset.Row[1].(string)
			entpKoNm = procRset.Row[2].(string)
			regDt = procRset.Row[3].(string)

			contentClAddItemList = append(contentClAddItemList, models.ContentClAddItem{
				BnrGrpSubSn: bnrGrpSubSn,
				EntpMemNo: entpMemNo,
				EntpKoNm: entpKoNm,
				RegDt: regDt,
			})
		}

		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}

	c.Data["json"] = &contentClAddItemList
	c.ServeJSON()

}
