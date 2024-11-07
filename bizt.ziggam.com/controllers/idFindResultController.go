package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"

	"bizt.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

type IdFindResultController struct {
	beego.Controller
}

func (c *IdFindResultController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

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

	// Start : Id Find Result List
	fmt.Printf(fmt.Sprintf("CALL ZSP_ID_FIND_RESULT_R('%v', '%v', '%v', :1)",
		pLang, pPpChrgNm, pBizRegNo))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_ID_FIND_RESULT_R('%v', '%v', '%v', :1)",
		pLang, pPpChrgNm, pBizRegNo),
		ora.S, /* MEM_ID */
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

	idFindResult := make([]models.IdFindResult, 0)

	var (
		memId string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			memId = procRset.Row[0].(string)

			idFindResult = append(idFindResult, models.IdFindResult{
				MemId: memId,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}

	// End : Id Find Result List

	c.Data["IdFindResult"] = idFindResult
	c.TplName = "common/id_find_result.html"
}
