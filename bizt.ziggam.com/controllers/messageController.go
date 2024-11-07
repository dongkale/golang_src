package controllers

import (
	"fmt"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"

	"bizt.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

type MessageController struct {
	BaseController
}

func (c *MessageController) Get() {
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
	pMsgEndYn := c.GetString("msg_end_yn")
	if pMsgEndYn == "" {
		pMsgEndYn = "0"
	}
	pKeyword := c.GetString("keyword")
	pTarget := c.GetString("target")
	pSn := c.GetString("sn")
	pEndYn := c.GetString("end_yn")
	if pEndYn == "" {
		pEndYn = "0"
	}
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

	// Start : Message Top Info
	fmt.Printf(fmt.Sprintf("CALL ZSP_MSG_TOP_INFO_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_MSG_TOP_INFO_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn),
		ora.I64, /* ING_CNT */
		ora.I64, /* END_CNT */
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

	messageTopInfo := make([]models.MessageTopInfo, 0)

	var (
		ingCnt int64
		endCnt int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {

			ingCnt = procRset.Row[0].(int64)
			endCnt = procRset.Row[1].(int64)

			messageTopInfo = append(messageTopInfo, models.MessageTopInfo{
				IngCnt: ingCnt,
				EndCnt: endCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Message Top Info

	// Start : Message Member List

	fmt.Printf(fmt.Sprintf("CALL ZSP_MSG_MEM_LIST_R('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pMsgEndYn, pKeyword))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_MSG_MEM_LIST_R('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pMsgEndYn, pKeyword),
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* MSG_SN */
		ora.S,   /* MSG_CONT */
		ora.S,   /* REG_DT */
		ora.S,   /* MSG_CFRM_YN */
		ora.S,   /* MSG_GBN_CD */
		ora.S,   /* PTO_PATH */
		ora.S,   /* NM */
		ora.I64, /* TOT_CNT */
		ora.S,   /* LAST_MSG_GBN_CD */
		ora.I64, /* MSG_CNT */
	)
	defer stmtProcCall.Close()
	if err != nil {
		panic(err)
	}
	procRset = &ora.Rset{}
	_, err = stmtProcCall.Exe(procRset)

	if err != nil {
		panic(err)
	}

	messageMemberList := make([]models.MessageMemberList, 0)

	var (
		mmEntpMemNo    string
		mmRecrutSn     string
		mmPpMemNo      string
		mmMsgSn        string
		mmMsgCont      string
		mmRegDt        string
		mmMsgCfrmYn    string
		mmMsgGbnCd     string
		mmPtoPath      string
		mmNm           string
		mmTotCnt       int64
		mmLastMsgGbnCd string
		mmMsgCnt       int64
		fullPtoPath    string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			mmEntpMemNo = procRset.Row[0].(string)
			mmRecrutSn = procRset.Row[1].(string)
			mmPpMemNo = procRset.Row[2].(string)
			mmMsgSn = procRset.Row[3].(string)
			mmMsgCont = procRset.Row[4].(string)
			mmRegDt = procRset.Row[5].(string)
			mmMsgCfrmYn = procRset.Row[6].(string)
			mmMsgGbnCd = procRset.Row[7].(string)
			mmPtoPath = procRset.Row[8].(string)

			if mmPtoPath == "" {
				fullPtoPath = mmPtoPath
			} else {
				fullPtoPath = imgServer + mmPtoPath
			}

			mmNm = procRset.Row[9].(string)
			mmTotCnt = procRset.Row[10].(int64)
			mmLastMsgGbnCd = procRset.Row[11].(string)
			mmMsgCnt = procRset.Row[12].(int64)

			messageMemberList = append(messageMemberList, models.MessageMemberList{
				MmEntpMemNo:    mmEntpMemNo,
				MmRecrutSn:     mmRecrutSn,
				MmPpMemNo:      mmPpMemNo,
				MmMsgSn:        mmMsgSn,
				MmMsgCont:      strings.Replace(mmMsgCont, "<br>", " ", -1),
				MmRegDt:        mmRegDt,
				MmMsgCfrmYn:    mmMsgCfrmYn,
				MmMsgGbnCd:     mmMsgGbnCd,
				MmPtoPath:      fullPtoPath,
				MmNm:           mmNm,
				MmTotCnt:       mmTotCnt,
				MmLastMsgGbnCd: mmLastMsgGbnCd,
				MmMsgCnt:       mmMsgCnt,
				Target:         pTarget,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Message Member List

	// Start : Entp Team Member List

	pGbnCd := "A"

	fmt.Printf(fmt.Sprintf("CALL ZSP_TEAM_MEM_LIST_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pGbnCd))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_TEAM_MEM_LIST_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pGbnCd),
		ora.I64, /* TOT_CNT */
		ora.S,   /* PP_CHRG_SN */
		ora.S,   /* PP_CHRG_GBN_CD */
		ora.S,   /* PP_CHRG_NM */
		ora.S,   /* PP_CHRG_BP_NM */
		ora.S,   /* EMAIL */
		ora.S,   /* ENTP_MEM_ID */
		ora.S,   /* PP_CHRG_TEL_NO */
	)
	defer stmtProcCall.Close()
	if err != nil {
		panic(err)
	}
	procRset = &ora.Rset{}
	_, err = stmtProcCall.Exe(procRset)

	if err != nil {
		panic(err)
	}

	entpTeamMemberList := make([]models.EntpTeamMemberList, 0)

	var (
		etTotCnt      int64
		etPpChrgSn    string
		etPpChrgGbnCd string
		etPpChrgNm    string
		etPpChrgBpNm  string
		etEmail       string
		etEntpMemId   string
		etPpChrgTelNo string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			etTotCnt = procRset.Row[0].(int64)
			etPpChrgSn = procRset.Row[1].(string)
			etPpChrgGbnCd = procRset.Row[2].(string)
			etPpChrgNm = procRset.Row[3].(string)
			etPpChrgBpNm = procRset.Row[4].(string)
			etEmail = procRset.Row[5].(string)
			etEntpMemId = procRset.Row[6].(string)
			etPpChrgTelNo = procRset.Row[7].(string)

			entpTeamMemberList = append(entpTeamMemberList, models.EntpTeamMemberList{
				EtTotCnt:      etTotCnt,
				EtPpChrgSn:    etPpChrgSn,
				EtPpChrgGbnCd: etPpChrgGbnCd,
				EtPpChrgNm:    etPpChrgNm,
				EtPpChrgBpNm:  etPpChrgBpNm,
				EtEmail:       etEmail,
				EtEntpMemId:   etEntpMemId,
				EtPpChrgTelNo: etPpChrgTelNo,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Entp Team Member List

	c.Data["EntpTeamMemberList"] = entpTeamMemberList
	c.Data["MessageMemberList"] = messageMemberList

	c.Data["Target"] = pTarget
	c.Data["TargetSn"] = pSn
	c.Data["TargetEndYn"] = pEndYn

	c.Data["IngCnt"] = ingCnt
	c.Data["EndCnt"] = endCnt
	c.Data["MmRecrutSn"] = pRecrutSn
	c.Data["TMenuId"] = "E00"

	c.TplName = "message/message.html"
}

func (c *MessageController) Post() {
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
	pMsgEndYn := c.GetString("msg_end_yn")
	if pMsgEndYn == "" {
		pMsgEndYn = "N"
	}
	pKeyword := c.GetString("keyword")

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

	// Start : Message Member List

	fmt.Printf(fmt.Sprintf("CALL ZSP_MSG_MEM_LIST_R('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pMsgEndYn, pKeyword))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_MSG_MEM_LIST_R('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pMsgEndYn, pKeyword),
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* MSG_SN */
		ora.S,   /* MSG_CONT */
		ora.S,   /* REG_DT */
		ora.S,   /* MSG_CFRM_YN */
		ora.S,   /* MSG_GBN_CD */
		ora.S,   /* PTO_PATH */
		ora.S,   /* NM */
		ora.I64, /* TOT_CNT */
		ora.S,   /* LAST_MSG_GBN_CD */
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

	rtnMessageMemberList := models.RtnMessageMemberList{}
	messageMemberList := make([]models.MessageMemberList, 0)

	var (
		mmEntpMemNo    string
		mmRecrutSn     string
		mmPpMemNo      string
		mmMsgSn        string
		mmMsgCont      string
		mmRegDt        string
		mmMsgCfrmYn    string
		mmMsgGbnCd     string
		mmPtoPath      string
		mmNm           string
		mmTotCnt       int64
		mmLastMsgGbnCd string
		mmMsgCnt       int64
		fullPtoPath    string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			mmEntpMemNo = procRset.Row[0].(string)
			mmRecrutSn = procRset.Row[1].(string)
			mmPpMemNo = procRset.Row[2].(string)
			mmMsgSn = procRset.Row[3].(string)
			mmMsgCont = procRset.Row[4].(string)
			mmRegDt = procRset.Row[5].(string)
			mmMsgCfrmYn = procRset.Row[6].(string)
			mmMsgGbnCd = procRset.Row[7].(string)
			mmPtoPath = procRset.Row[8].(string)

			if mmPtoPath == "" {
				fullPtoPath = mmPtoPath
			} else {
				fullPtoPath = imgServer + mmPtoPath
			}

			mmNm = procRset.Row[9].(string)
			mmTotCnt = procRset.Row[10].(int64)
			mmLastMsgGbnCd = procRset.Row[11].(string)
			mmMsgCnt = procRset.Row[12].(int64)

			messageMemberList = append(messageMemberList, models.MessageMemberList{
				MmEntpMemNo:    mmEntpMemNo,
				MmRecrutSn:     mmRecrutSn,
				MmPpMemNo:      mmPpMemNo,
				MmMsgSn:        mmMsgSn,
				MmMsgCont:      mmMsgCont,
				MmRegDt:        mmRegDt,
				MmMsgCfrmYn:    mmMsgCfrmYn,
				MmMsgGbnCd:     mmMsgGbnCd,
				MmPtoPath:      fullPtoPath,
				MmNm:           mmNm,
				MmTotCnt:       mmTotCnt,
				MmLastMsgGbnCd: mmLastMsgGbnCd,
				MmMsgCnt:       mmMsgCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnMessageMemberList = models.RtnMessageMemberList{
			RtnMessageMemberListData: messageMemberList,
		}

	}
	// End : Message Member List

	c.Data["json"] = &rtnMessageMemberList
	c.ServeJSON()
}
