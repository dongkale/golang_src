package controllers

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/astaxie/beego/logs"
	"gopkg.in/rana/ora.v4"
	"managert.ziggam.com/models"
)

type ContentsController struct {
	BaseController
}

func (c *ContentsController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	pLang, _ := beego.AppConfig.String("lang")
	pBnrGrpSn := c.GetString("bnr_grp_sn") // 사용여부
	pBnrGrpTypCd := c.GetString("bnr_grp_type_cd") // 사용여부

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
	log.Debug("CALL MNG_GRP_BNR_INFO('%v', '%v', :1)",
		pLang, pBnrGrpSn)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL MNG_GRP_BNR_INFO('%v', '%v',:1)",
		pLang, pBnrGrpSn),
		ora.S,   /* BNR_GRP_TITLE */
		ora.S,   /* BNR_GRP_SN */
		ora.S,   /* BNR_GRP_TYP_CD */
		ora.S,   /* REG_DT */
		ora.S,   /* REG_ID */
		ora.I64, /* ROL_TM */
		ora.S,   /* PUBL_SDY */
		ora.S,   /* PUBL_EDY */
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

	groupBannerData := models.GroupBannerData{}

	if procRset.IsOpen() {
		for procRset.Next() {
			groupBannerData.BnrGrpTitle = procRset.Row[0].(string)
			groupBannerData.BnrGrpSn = procRset.Row[1].(string)
			groupBannerData.BnrGrpTypCd = procRset.Row[2].(string)
			groupBannerData.RegDt = procRset.Row[3].(string)
			groupBannerData.RegId = procRset.Row[4].(string)
			groupBannerData.RolTm = procRset.Row[5].(int64)
			groupBannerData.PublSdy = procRset.Row[6].(string)
			groupBannerData.PublEdy = procRset.Row[7].(string)

		}

		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// 그룹 배너 단일 내역

	c.Data["BnrGrpSn"] = pBnrGrpSn
	c.Data["BnrGrpTypCd"] = pBnrGrpTypCd
	c.Data["GrpBnrData"] = &groupBannerData

	viewName := "contents/contents.html"

	if pBnrGrpTypCd == "CL" {
		viewName = "contents/contents_cl.html"
	}
	c.TplName = viewName
}
