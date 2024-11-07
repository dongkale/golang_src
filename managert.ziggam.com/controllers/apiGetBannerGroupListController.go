package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	"gopkg.in/rana/ora.v4"
	"managert.ziggam.com/models"
)

type ApiGetGroupBannerListController struct {
	BaseController
}

func (c *ApiGetGroupBannerListController) Get() {
	// start : log
	// log := logs.NewLogger()
	// log.SetLogger(logs.AdapterConsole)
	// end : log

	pLang, _ := beego.AppConfig.String("lang")
	//
	//pEntpMemNo := mem_no

	pBnrGrpTypCd := c.GetString("bnr_grp_typ_cd") // 사용여부
	pUseYn := c.GetString("use_yn")               // 사용여부
	pPublStat := c.GetString("publ_stat")         // 게시상태
	pPublSdy := c.GetString("publ_sdy")           // 게시시작일시
	pPublEdy := c.GetString("publ_edy")           // 게시종료일시
	pKeyword := c.GetString("keyword")            // 검색어

	if pBnrGrpTypCd == "" {
		pBnrGrpTypCd = "ALL"
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

	logs.Debug("CALL MNG_LIST_GROUP_BANNER_INFO('%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pBnrGrpTypCd, pUseYn, pPublStat, pPublSdy, pPublEdy, pKeyword)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL MNG_LIST_GROUP_BANNER_INFO('%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pBnrGrpTypCd, pUseYn, pPublStat, pPublSdy, pPublEdy, pKeyword),
		ora.I64, /* RTN_CD */
		ora.S,   /* BNR_GRP_TYP_CD */
		ora.S,   /* BNR_GRP_SN */
		ora.S,   /* BNR_GRP_TITLE */
		ora.S,   /* PHO_PATH */
		ora.S,   /* PUBL_SDY */
		ora.S,   /* PUBL_EDY */
		ora.S,   /* PUBL_STAT */
		ora.S,   /* USE_YN */
		ora.S,   /* REG_DT */
		ora.S,   /* REG_ID */
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
		rtnCd       int64
		bnrGrpTypCd string
		bnrGrpSn    string
		bnrGrpTitle string
		ptoPath     string
		publSdy     string
		publEdy     string
		publStat    string
		useYn       string
		regDt       string
		regId       string
	)

	//rtnGroupBannerList := models.GroupBannerDataList{}
	groupBannerDataList := make([]models.GroupBannerData, 0)

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			bnrGrpTypCd = procRset.Row[1].(string)
			bnrGrpSn = procRset.Row[2].(string)
			bnrGrpTitle = procRset.Row[3].(string)
			ptoPath = procRset.Row[4].(string)
			publSdy = procRset.Row[5].(string)
			publEdy = procRset.Row[6].(string)
			publStat = procRset.Row[7].(string)
			useYn = procRset.Row[8].(string)
			regDt = procRset.Row[9].(string)
			regId = procRset.Row[10].(string)

			groupBannerDataList = append(groupBannerDataList, models.GroupBannerData{
				RtnCd:       rtnCd,
				BnrGrpTypCd: bnrGrpTypCd,
				BnrGrpSn:    bnrGrpSn,
				BnrGrpTitle: bnrGrpTitle,
				PtoPath:     ptoPath,
				PublSdy:     publSdy,
				PublEdy:     publEdy,
				PublStat:    publStat,
				UseYn:       useYn,
				RegDt:       regDt,
				RegId:       regId,
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

	c.Data["json"] = &groupBannerDataList
	c.ServeJSON()

}
