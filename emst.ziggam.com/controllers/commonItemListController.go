package controllers

import (
	"fmt"

	"emst.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	"gopkg.in/rana/ora.v4"
)

type CommonItemListController struct {
	beego.Controller
}

func (c *CommonItemListController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	//session := c.StartSession()
	pLang, _ := beego.AppConfig.String("lang")

	pLnkGbnVal := c.GetString("lnk_gbn_val")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Item List
	log.Debug("CALL SP_EMS_ADMIN_ITEM_LIST_R('%v', '%v', :1)",
		pLang, pLnkGbnVal)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_ITEM_LIST_R('%v', '%v', :1)",
		pLang, pLnkGbnVal),
		ora.S, /* ITEM1 */
		ora.S, /* ITEM2 */

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

	rtnCommonItemList := models.RtnCommonItemList{}
	commonItemList := make([]models.CommonItemList, 0)

	var (
		item1 string
		item2 string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			item1 = procRset.Row[0].(string)
			item2 = procRset.Row[1].(string)

			commonItemList = append(commonItemList, models.CommonItemList{
				Item1: item1,
				Item2: item2,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnCommonItemList = models.RtnCommonItemList{
			RtnCommonItemListData: commonItemList,
		}
	}
	// End : Item List

	c.Data["json"] = &rtnCommonItemList
	c.ServeJSON()
}
