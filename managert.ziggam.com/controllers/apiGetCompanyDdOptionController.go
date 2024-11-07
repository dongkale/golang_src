package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/astaxie/beego/logs"
	"gopkg.in/rana/ora.v4"
	"managert.ziggam.com/models"
)

type ApiGetCompanyDdOptionController struct {
	BaseController
}

func (c *ApiGetCompanyDdOptionController) Get() {
	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	pLang, _ := beego.AppConfig.String("lang")

	pSearchText := c.GetString("search_text")

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

	log.Debug("CALL MNG_LIST_COMPANY_OPTION('%v', '%v', :1)",
		pLang, pSearchText)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL MNG_LIST_COMPANY_OPTION('%v', '%v', :1)",
		pLang, pSearchText),
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
		entpMemNo string
		entpKoNm  string
	)

	//rtnGroupBannerList := models.GroupBannerDataList{}
	contentClAddItemList := make([]models.ContentClAddItem, 0)

	if procRset.IsOpen() {
		for procRset.Next() {
			entpMemNo = procRset.Row[0].(string)
			entpKoNm = procRset.Row[1].(string)

			contentClAddItemList = append(contentClAddItemList, models.ContentClAddItem{
				EntpMemNo: entpMemNo,
				EntpKoNm: entpKoNm,
			})
		}

		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}

	c.Data["json"] = &contentClAddItemList
	c.ServeJSON()

}
