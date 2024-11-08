package controllers

import (
	"fmt"

	"emst.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type AdminQaHistoryListController struct {
	beego.Controller
}

func (c *AdminQaHistoryListController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	//session := c.StartSession()
	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := c.GetString("entp_mem_no")
	pRecrutSn := c.GetString("recrut_sn")
	pPpMemNo := c.GetString("pp_mem_no")

	cdnPath, _ := beego.AppConfig.String("cdnpath")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Question and Answer History List
	log.Debug("CALL SP_EMS_AM_QA_VD_HIS_LIST_R('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_AM_QA_VD_HIS_LIST_R('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo),
		ora.S, /* QST_SN */
		ora.S, /* VD_TITLE */
		ora.S, /* VD_FILE_PATH */
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

	rtnQaHistoryList := models.RtnQaHistoryList{}
	qaHistoryList := make([]models.QaHistoryList, 0)

	var (
		qstSn      string
		vdTitle    string
		vdFilePath string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			qstSn = procRset.Row[0].(string)
			vdTitle = procRset.Row[1].(string)
			vdFilePath = procRset.Row[2].(string)

			var fullFilePath string

			if vdFilePath == "" {
				fullFilePath = vdFilePath
			} else {
				fullFilePath = cdnPath + vdFilePath
			}

			qaHistoryList = append(qaHistoryList, models.QaHistoryList{
				QstSn:      qstSn,
				VdTitle:    vdTitle,
				VdFilePath: fullFilePath,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnQaHistoryList = models.RtnQaHistoryList{
			RtnQaHistoryListData: qaHistoryList,
		}
	}
	// End : Question and Answer History List

	c.Data["json"] = &rtnQaHistoryList
	c.ServeJSON()
}
