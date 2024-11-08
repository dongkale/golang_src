package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/rana/ora.v4"
	"managert.ziggam.com/models"

	"github.com/beego/beego/v2/core/logs"
)

type ApiGetBannerGroupDetailListController struct {
	BaseController
}

func (c *ApiGetBannerGroupDetailListController) Get() {
	// start : log
	// log := logs.NewLogger()
	// log.SetLogger(logs.AdapterConsole)
	// end : log

	imgServer, _ := beego.AppConfig.String("viewpath")

	pLang, _ := beego.AppConfig.String("lang")
	pBnrGrpSn := c.GetString("bnr_grp_sn") // 사용여부
	pBnrGrpTypCd := c.GetString("bnr_gpr_typ_cd")
	pUseYn := c.GetString("use_yn")
	pPublStat := c.GetString("publ_stat")
	pPublSdy := c.GetString("publ_sdy")
	pPublEdy := c.GetString("publ_edy")
	pRegSdy := c.GetString("reg_sdy")
	pRegEdy := c.GetString("reg_edy")
	pKeyword := c.GetString("keyword")

	// TODO 데이터 로드 로직

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

	// 그룹 배너 단일 내역

	// 그룹 배너 단일 내역

	// 그룹 배너 상세 리스트
	logs.Debug("CALL MNG_LIST_GRP_BAR_DTL_INFO('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pBnrGrpSn, pBnrGrpTypCd, pUseYn, pPublStat, pPublSdy, pPublEdy, pRegSdy, pRegEdy, pKeyword)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL MNG_LIST_GRP_BAR_DTL_INFO('%v', '%v', '%v', '%v','%v', '%v','%v', '%v','%v', '%v',:1)",
		pLang, pBnrGrpSn, pBnrGrpTypCd, pUseYn, pPublStat, pPublSdy, pPublEdy, pRegSdy, pRegEdy, pKeyword),
		ora.I64, /* SW_IDX */
		ora.S,   /* USE_YTN */
		ora.S,   /* BNR_SN */
		ora.S,   /* BNR_TITLE */
		ora.S,   /* PTO_PATH */
		ora.S,   /* THUMB_PTO_PATH */
		ora.S,   /* PUBL_SDY */
		ora.S,   /* PUBL_EDY */
		ora.S,   /* PUBL_STAT */
		ora.S,   /* REG_DT */
		ora.S,   /* BNR_GRP_SUB_SN */
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
		swIdx        int64
		useYn        string
		bnrSn        string
		bnrTitle     string
		ptoPath      string
		thumbPtoPath string
		publSdy      string
		publEdy      string
		publStat     string
		regDt        string
		bnrGrpSubSn  string
	)

	//rtnGroupBannerList := models.GroupBannerDataList{}
	groupBannerDetailDataList := make([]models.GroupBannerDetailData, 0)

	if procRset.IsOpen() {
		for procRset.Next() {
			swIdx = procRset.Row[0].(int64)
			useYn = procRset.Row[1].(string)
			bnrSn = procRset.Row[2].(string)
			bnrTitle = procRset.Row[3].(string)
			ptoPath = procRset.Row[4].(string)
			thumbPtoPath = procRset.Row[5].(string)
			publSdy = procRset.Row[6].(string)
			publEdy = procRset.Row[7].(string)
			publStat = procRset.Row[8].(string)
			regDt = procRset.Row[9].(string)
			bnrGrpSubSn = procRset.Row[10].(string)

			groupBannerDetailDataList = append(groupBannerDetailDataList, models.GroupBannerDetailData{
				SwIdx:        swIdx,
				UseYn:        useYn,
				BnrSn:        bnrSn,
				BnrTitle:     bnrTitle,
				PtoPath:      ptoPath,
				ThumbPtoPath: thumbPtoPath,
				PublSdy:      publSdy,
				PublEdy:      publEdy,
				PublStat:     publStat,
				RegDt:        regDt,
				BnrGrpSubSn: bnrGrpSubSn,
				PtoFullPath: imgServer + ptoPath,
				ThumbPtoFullPath: imgServer + thumbPtoPath,
			})
		}

		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// 그룹 배너 상세 리스트

	logs.Debug("CALL MNG_LIST_GRP_BAR_DTL_INFO RESULT ('%v')", groupBannerDetailDataList)

	c.Data["json"] = &groupBannerDetailDataList
	c.ServeJSON()

}
