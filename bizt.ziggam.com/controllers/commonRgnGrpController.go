package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	"gopkg.in/rana/ora.v4"
)

type CommonRgnGrpController struct {
	beego.Controller
}

func (c *CommonRgnGrpController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	//session := c.StartSession()
	pLang, _ := beego.AppConfig.String("lang")

	pUpRgnGrpCd := c.GetString("up_rgn_grp_cd")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Job Group List
	fmt.Printf(fmt.Sprintf("CALL ZSP_RGN_LIST_R_V2('%v', '%v', :1)",
		pLang, pUpRgnGrpCd))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_RGN_LIST_R_V2('%v', '%v', :1)",
		pLang, pUpRgnGrpCd),
		ora.S, /* RGN_CD */
		ora.S, /* RGN_NM */
		ora.S, /* RGN_F_NM */

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

	rgnGrpList := make([]models.RgnGrp, 0)

	var (
		rgnGrpCd  string
		rgnGrpNm  string
		rgnGrpFNm string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			rgnGrpCd = procRset.Row[0].(string)
			rgnGrpNm = procRset.Row[1].(string)
			rgnGrpFNm = procRset.Row[2].(string)

			rgnGrpList = append(rgnGrpList, models.RgnGrp{
				Code:     rgnGrpCd,
				Name:     rgnGrpNm,
				FullName: rgnGrpFNm,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Job Group List

	c.Data["json"] = &rgnGrpList
	c.ServeJSON()
}
