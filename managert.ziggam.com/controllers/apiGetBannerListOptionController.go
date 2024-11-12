package controllers

import (
	"fmt"

	"managert.ziggam.com/models"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type ApiGetBannerListOptionController struct {
	BaseController
}

func (c *ApiGetBannerListOptionController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	//session := c.StartSession()
	pLang, _ := beego.AppConfig.String("lang")

	pLnkGbnVal := c.GetString("lnk_gbn_val")
	pBnrSn := c.GetString("bnr_sn")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Banner Kind List
	log.Debug("CALL SP_EMS_ADMIN_BNR_KIND_LIST_R('%v', '%v', '%v', :1)",
		pLang, pLnkGbnVal, pBnrSn)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_BNR_KIND_LIST_R('%v', '%v', '%v', :1)",
		pLang, pLnkGbnVal, pBnrSn),
		ora.S, /* ENTP_MEM_NO */
		ora.S, /* RECRUT_SN */
		ora.S, /* ENTP_KO_NM */
		ora.S, /* RECRUT_TITLE */
		ora.S, /* SELECT_YN */
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

	rtnBannerKindList := models.RtnBannerKindList{}
	bannerKindList := make([]models.BannerKindList, 0)

	var (
		entpMemNo   string
		recrutSn    string
		entpKoNm    string
		recrutTitle string
		selectYn    string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			entpMemNo = procRset.Row[0].(string)
			recrutSn = procRset.Row[1].(string)
			entpKoNm = procRset.Row[2].(string)
			recrutTitle = procRset.Row[3].(string)
			selectYn = procRset.Row[4].(string)

			bannerKindList = append(bannerKindList, models.BannerKindList{
				EntpMemNo:   entpMemNo,
				RecrutSn:    recrutSn,
				EntpKoNm:    entpKoNm,
				RecrutTitle: recrutTitle,
				SelectYn:    selectYn,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnBannerKindList = models.RtnBannerKindList{
			RtnBannerKindListData: bannerKindList,
		}
	}
	// End : Banner Kind List

	c.Data["json"] = &rtnBannerKindList
	c.ServeJSON()
}
