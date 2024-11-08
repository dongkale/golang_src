package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"

	"emst.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

type CommonFindIdController struct {
	beego.Controller
}

func (c *CommonFindIdController) Get() {
	c.TplName = "common/id_find.html"
}

func (c *CommonFindIdController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	//session := c.StartSession()
	pLang, _ := beego.AppConfig.String("lang")

	pPpChrgNm := c.GetString("pp_chrg_nm")
	pBizRegNo := c.GetString("biz_reg_no")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	log.Debug("CALL SP_EMS_ID_FIND_R( '%v', '%v', '%v', :1)",
		pLang, pPpChrgNm, pBizRegNo)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ID_FIND_R( '%v', '%v', '%v', :1)",
		pLang, pPpChrgNm, pBizRegNo),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
		ora.S,   /* FIND_ID */
		ora.S,   /* REG_DY */
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
		rtnCd    int64
		rtnMsg   string
		resultId string
		resultDy string
	)

	findId := models.FindId{}
	rtnFindId := models.RtnFindId{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)

			if rtnCd == 1 {
				resultId = procRset.Row[2].(string)
				resultDy = procRset.Row[3].(string)

				findId = models.FindId{
					ResultId: resultId,
					ResultDy: resultDy,
				}
			}

		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnFindId = models.RtnFindId{
			RtnCd:   rtnCd,
			RtnMsg:  rtnMsg,
			RtnData: findId,
		}
	}

	c.Data["json"] = &rtnFindId
	c.ServeJSON()

}
