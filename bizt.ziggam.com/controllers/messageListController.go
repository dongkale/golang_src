package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	"gopkg.in/rana/ora.v4"
)

type MessageListController struct {
	beego.Controller
}

func (c *MessageListController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	mem_id := session.Get(c.Ctx.Request.Context(), "mem_id")
	if mem_id == nil {
		c.Ctx.Redirect(302, "/login")
	}

	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := mem_no
	pRecrutSn := c.GetString("recrut_sn")
	//pRecrutSn = "2019070481"
	pPpMemNo := c.GetString("pp_mem_no")

	imgServer, _  := beego.AppConfig.String("viewpath")
	//cdnPath := beego.AppConfig.String("cdnpath")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Message Detail List
	fmt.Printf(fmt.Sprintf("CALL ZSP_MSG_LIST_R_V2('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_MSG_LIST_R_V2('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo),
		ora.S, /* ENTP_MEM_NO */
		ora.S, /* RECRUT_SN */
		ora.S, /* PP_MEM_NO */
		ora.S, /* MSG_SN */
		ora.S, /* PTO_PATH */
		ora.S, /* LD_YN */
		ora.S, /* LD_DT */
		ora.S, /* LD_DT2 */
		ora.S, /* MEM_GBN */
		ora.S, /* MSG_GBN_CD */
		ora.S, /* MSG_CONT */
		ora.S, /* MSG_CFRM_YN */
		ora.S, /* ENTP_NM */
		ora.S, /* MEM_NM */
		ora.S, /* LIVE_ITV_SDAY */
		ora.S, /* LIVE_ITV_STIME */
		ora.S, /* LIVE_ITV_STIME2 */
		ora.S, /* LIVE_ITV_EDAY */
		ora.S, /* LIVE_ITV_ETIME */
		ora.S, /* LIVE_ITV_ETIME2 */
		ora.S, /* LIVE_ITV_JT */
		ora.S, /* LIVE_SN */
		ora.S, /* LIVE_NVN_YN */
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

	rtnMessageList := models.RtnMessageList{}
	messageList := make([]models.MessageList, 0)

	var (
		mlEntpMemNo     string
		mlRecrutSn      string
		mlPpMemNo       string
		mlMsgSn         string
		mlPtoPath       string
		mlLdYn          string
		mlLdDt          string
		mlLdDt2         string
		mlMemGbn        string
		mlMsgGbnCd      string
		mlMsgCont       string
		mlMsgCfrmYn     string
		mlEntpNm        string
		mlMemNm         string
		mlLiveItvSday   string
		mlLiveItvStime  string
		mlLiveItvStime2 string
		mlLiveItvEday   string
		mlLiveItvEtime  string
		mlLiveItvEtime2 string
		mlLiveItvJt     string
		mlLiveSn        string
		mlLiveNvnYn     string
		mlFullPtoPath   string
	)

	if procRset.IsOpen() {
		for procRset.Next() {

			mlEntpMemNo = procRset.Row[0].(string)
			mlRecrutSn = procRset.Row[1].(string)
			mlPpMemNo = procRset.Row[2].(string)
			mlMsgSn = procRset.Row[3].(string)
			mlPtoPath = procRset.Row[4].(string)
			if mlPtoPath == "" {
				mlFullPtoPath = mlPtoPath
			} else {
				mlFullPtoPath = imgServer + mlPtoPath
			}

			mlLdYn = procRset.Row[5].(string)
			mlLdDt = procRset.Row[6].(string)
			mlLdDt2 = procRset.Row[7].(string)
			mlMemGbn = procRset.Row[8].(string)
			mlMsgGbnCd = procRset.Row[9].(string)
			mlMsgCont = procRset.Row[10].(string)
			mlMsgCfrmYn = procRset.Row[11].(string)
			mlEntpNm = procRset.Row[12].(string)
			mlMemNm = procRset.Row[13].(string)
			mlLiveItvSday = procRset.Row[14].(string)
			mlLiveItvStime = procRset.Row[15].(string)
			mlLiveItvStime2 = procRset.Row[16].(string)
			mlLiveItvEday = procRset.Row[17].(string)
			mlLiveItvEtime = procRset.Row[18].(string)
			mlLiveItvEtime2 = procRset.Row[19].(string)
			mlLiveItvJt = procRset.Row[20].(string)
			mlLiveSn = procRset.Row[21].(string)
			mlLiveNvnYn = procRset.Row[22].(string)

			messageList = append(messageList, models.MessageList{
				MlEntpMemNo:     mlEntpMemNo,
				MlRecrutSn:      mlRecrutSn,
				MlPpMemNo:       mlPpMemNo,
				MlMsgSn:         mlMsgSn,
				MlPtoPath:       mlFullPtoPath,
				MlLdYn:          mlLdYn,
				MlLdDt:          mlLdDt,
				MlLdDt2:         mlLdDt2,
				MlMemGbn:        mlMemGbn,
				MlMsgGbnCd:      mlMsgGbnCd,
				MlMsgCont:       mlMsgCont,
				MlMsgCfrmYn:     mlMsgCfrmYn,
				MlEntpNm:        mlEntpNm,
				MlMemNm:         mlMemNm,
				MlLiveItvSday:   mlLiveItvSday,
				MlLiveItvStime:  mlLiveItvStime,
				MlLiveItvStime2: mlLiveItvStime2,
				MlLiveItvEday:   mlLiveItvEday,
				MlLiveItvEtime:  mlLiveItvEtime,
				MlLiveItvEtime2: mlLiveItvEtime2,
				MlLiveItvJt:     mlLiveItvJt,
				MlLiveSn:        mlLiveSn,
				MlLiveNvnYn:     mlLiveNvnYn,
			})

			fmt.Printf(fmt.Sprintf("[MessageList][Live] EntpMemN:%v, RecrutSn:%v, PpMemNo:%v, LiveSn:%v, LiveNvN:%v",
				pEntpMemNo, pRecrutSn, pPpMemNo, mlLiveSn, mlLiveNvnYn))
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnMessageList = models.RtnMessageList{
			RtnMessageListData: messageList,
		}
	}
	// End : Message Detail List

	c.Data["json"] = &rtnMessageList
	c.ServeJSON()
}
