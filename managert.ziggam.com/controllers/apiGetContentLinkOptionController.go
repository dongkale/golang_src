package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/astaxie/beego/logs"
	"gopkg.in/rana/ora.v4"
	"managert.ziggam.com/models"
)

type ApiGetContentLinkOptionController struct {
	BaseController
}

func (c *ApiGetContentLinkOptionController) Get() {
	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	pLang, _ := beego.AppConfig.String("lang")
	//
	//pEntpMemNo := mem_no

	pOptCode := c.GetString("opt_code")        // 옵션코드
	pOptSubCode := c.GetString("opt_sub_code") // 서브옵션코드
	pBnrSn := c.GetString("bnr_sn")            // 공고번호

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

	log.Debug("CALL MNG_LIST_CNTRL_LNK_OP('%v', '%v', '%v', '%v', :1)",
		pLang, pOptCode, pOptSubCode, pBnrSn)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL MNG_LIST_CNTRL_LNK_OP('%v', '%v', '%v', '%v',:1)",
		pLang, pOptCode, pOptSubCode, pBnrSn),
		ora.S, /* ITEM1 */
		ora.S, /* ITEM2 */
		ora.S, /* ITEM3 */
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
		optionKey   string
		optionValue string
		selectYn    string
	)

	//rtnGroupBannerList := models.GroupBannerDataList{}
	optionList := make([]models.OptionList, 0)

	if procRset.IsOpen() {
		for procRset.Next() {
			optionKey = procRset.Row[0].(string)
			optionValue = procRset.Row[1].(string)
			selectYn = procRset.Row[2].(string)

			optionList = append(optionList, models.OptionList{
				OptionKey:   optionKey,
				OptionValue: optionValue,
				SelectYn:    selectYn,
			})
		}

		if err := procRset.Err(); err != nil {
			panic(err)
		}
		//rtnGroupBannerList = models.GroupBannerDataList{
		//	RtnGroupBannerData: groupBannerData,
		//}
	}

	//log.Debug("rtnGroupBannerList('%v')", rtnGroupBannerList)

	c.Data["json"] = &optionList
	c.ServeJSON()

}
