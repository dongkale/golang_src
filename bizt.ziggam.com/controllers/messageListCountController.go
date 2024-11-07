package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/rana/ora.v4"
)

// MessageListCountController ...
type MessageListCountController struct {
	beego.Controller
}

func (c *MessageListCountController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")

	if mem_no == nil {
		c.Data["json"] = &models.RtnMessageListCount{}
		c.ServeJSON()
		return
	}

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := mem_no
	pRecrutSn := c.GetString("recrut_sn")
	pPpMemNo := c.GetString("pp_mem_no")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Message List Count
	fmt.Printf(fmt.Sprintf("CALL ZSP_MSG_LIST_CNT_R('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_MSG_LIST_CNT_R('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo),
		ora.I64, /* MSG_CNT */
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

	rtnMessageListCount := models.RtnMessageListCount{}
	messageListCount := make([]models.MessageListCount, 0)

	var (
		msgCnt int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			msgCnt = procRset.Row[0].(int64)

			messageListCount = append(messageListCount, models.MessageListCount{
				MsgCnt: msgCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnMessageListCount = models.RtnMessageListCount{
			RtnMessageListCountData: messageListCount,
		}
	}
	// End : Message List Count

	c.Data["json"] = &rtnMessageListCount
	c.ServeJSON()
}
