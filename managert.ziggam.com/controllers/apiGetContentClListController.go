package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/astaxie/beego/logs"
	"gopkg.in/rana/ora.v4"
	"managert.ziggam.com/models"
)

type ApiGetContentClListController struct {
	BaseController
}

func (c *ApiGetContentClListController) Get() {
	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	pLang, _ := beego.AppConfig.String("lang")

	pBnrGrpSn := c.GetString("bnr_grp_sn")    // 사용여부
	pUseYn := c.GetString("use_yn")    // 사용여부
	pKeyword := c.GetString("keyword") // 게시상태

	if pUseYn == "" {
		pUseYn = "0"
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

	// Start : Applicant Delete Process

	log.Debug("CALL MNG_LIST_CL_GROUP_INFO('%v', '%v', %v, '%v', :1)",
		pLang, pBnrGrpSn, pUseYn, pKeyword)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL MNG_LIST_CL_GROUP_INFO('%v', '%v', %v, '%v', :1)",
		pLang,  pBnrGrpSn, pUseYn, pKeyword),
		ora.S, /* BNR_GRP_SN */
		ora.S, /* USE_YN */
		ora.S, /* ENTP_MEM_NO */
		ora.S, /* ENTP_KO_NM */
		ora.S, /* UPT_DT */
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
		bnrGrpSn  string
		useYn     string
		entpMemNo string
		entpKoNm  string
		uptDt     string
		bnrGrpSubSn	string
	)

	//rtnGroupBannerList := models.GroupBannerDataList{}
	contentClItemList := make([]models.ContentClItem, 0)

	if procRset.IsOpen() {
		for procRset.Next() {
			bnrGrpSn = procRset.Row[0].(string)
			useYn = procRset.Row[1].(string)
			entpMemNo = procRset.Row[2].(string)
			entpKoNm = procRset.Row[3].(string)
			uptDt = procRset.Row[4].(string)
			bnrGrpSubSn = procRset.Row[5].(string)

			contentClItemList = append(contentClItemList, models.ContentClItem{
				BnrGrpSn: bnrGrpSn,
				UseYn: useYn,
				EntpMemNo: entpMemNo,
				EntpKoNm: entpKoNm,
				UptDt: uptDt,
				BnrGrpSubSn: bnrGrpSubSn,
			})
		}

		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}

	c.Data["json"] = &contentClItemList
	c.ServeJSON()

}
