package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type NotificationListController struct {
	BaseController
}

func (c *NotificationListController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/login")
	}

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := mem_no
	pGbnCd := c.GetString("gbn_cd") //구분코드
	if pGbnCd == "" {
		pGbnCd = "02" //페이지
	}
	imgServer, _  := beego.AppConfig.String("viewpath")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Notification List

	fmt.Printf(fmt.Sprintf("CALL ZSP_NOTIFICATION_LIST_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pGbnCd))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_NOTIFICATION_LIST_R('%v','%v', '%v', :1)",
		pLang, pEntpMemNo, pGbnCd),
		ora.S,   /* MEM_NO */
		ora.S,   /* DVC_ID */
		ora.S,   /* NOTI_KND_CD */
		ora.S,   /* REG_DT */
		ora.S,   /* CFRM_DT */
		ora.S,   /* DT_DD */
		ora.S,   /* DT_HH */
		ora.S,   /* PTO_PATH */
		ora.S,   /* NOTI_CONT */
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* LIVE_SN */
		ora.S,   /* BRD_GBN_CD */
		ora.I64, /* SN */
		ora.S,   /* LD_YN */
		ora.S,   /* INQ_GBN_CD */
		ora.S,   /* INQ_REG_DT */
		ora.S,   /* MSG_END_YN */
		ora.S,   /* LIVE_NVN_YN */
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

	notificationList := make([]models.NotificationList, 0)

	var (
		ntMemNo       string
		ntDvcId       string
		ntNotiKndCd   string
		ntRegDt       string
		ntCfrmDt      string
		ntDtDd        string
		ntDtHh        string
		ntPtoPath     string
		ntNotiCont    string
		ntEntpMemNo   string
		ntRecrutSn    string
		ntPpMemNo     string
		ntLiveSn      string
		ntBrdGbnCd    string
		ntSn          int64
		ntLdYn        string
		ntInqGbnCd    string
		ntInqRegDt    string
		ntMsgEndYn    string
		ntLiveNvnYn   string
		ntFullPtoPath string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			ntMemNo = procRset.Row[0].(string)
			ntDvcId = procRset.Row[1].(string)
			ntNotiKndCd = procRset.Row[2].(string)
			ntRegDt = procRset.Row[3].(string)
			ntCfrmDt = procRset.Row[4].(string)
			ntDtDd = procRset.Row[5].(string)
			ntDtHh = procRset.Row[6].(string)
			ntPtoPath = procRset.Row[7].(string)
			if ntPtoPath == "" {
				ntFullPtoPath = ntPtoPath
			} else {
				ntFullPtoPath = imgServer + ntPtoPath
			}
			ntNotiCont = procRset.Row[8].(string)
			ntEntpMemNo = procRset.Row[9].(string)
			ntRecrutSn = procRset.Row[10].(string)
			ntPpMemNo = procRset.Row[11].(string)
			ntLiveSn = procRset.Row[12].(string)
			ntBrdGbnCd = procRset.Row[13].(string)
			ntSn = procRset.Row[14].(int64)
			ntLdYn = procRset.Row[15].(string)
			ntInqGbnCd = procRset.Row[16].(string)
			ntInqRegDt = procRset.Row[17].(string)
			ntMsgEndYn = procRset.Row[18].(string)
			ntLiveNvnYn = procRset.Row[19].(string)

			notificationList = append(notificationList, models.NotificationList{
				NtMemNo:     ntMemNo,
				NtDvcId:     ntDvcId,
				NtNotiKndCd: ntNotiKndCd,
				NtRegDt:     ntRegDt,
				NtCfrmDt:    ntCfrmDt,
				NtDtDd:      ntDtDd,
				NtDtHh:      ntDtHh,
				NtPtoPath:   ntFullPtoPath,
				NtNotiCont:  ntNotiCont,
				NtEntpMemNo: ntEntpMemNo,
				NtRecrutSn:  ntRecrutSn,
				NtPpMemNo:   ntPpMemNo,
				NtLiveSn:    ntLiveSn,
				NtBrdGbnCd:  ntBrdGbnCd,
				NtSn:        ntSn,
				NtLdYn:      ntLdYn,
				NtInqGbnCd:  ntInqGbnCd,
				NtInqRegDt:  ntInqRegDt,
				NtMsgEndYn:  ntMsgEndYn,
				NtLiveNvnYn: ntLiveNvnYn,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Notification List

	c.Data["NotificationList"] = notificationList
	c.Data["TMenuId"] = "M00"
	c.Data["SMenuId"] = "M00"

	c.TplName = "notification/notification_list.html"
}

func (c *NotificationListController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/login")
	}

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := mem_no
	pGbnCd := c.GetString("gbn_cd") //구분코드
	if pGbnCd == "" {
		pGbnCd = "01" //팝업
	}
	imgServer, _  := beego.AppConfig.String("viewpath")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Notification PopUp List

	fmt.Printf(fmt.Sprintf("CALL ZSP_NOTIFICATION_LIST_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pGbnCd))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_NOTIFICATION_LIST_R('%v','%v', '%v', :1)",
		pLang, pEntpMemNo, pGbnCd),
		ora.S,   /* MEM_NO */
		ora.S,   /* DVC_ID */
		ora.S,   /* NOTI_KND_CD */
		ora.S,   /* REG_DT */
		ora.S,   /* CFRM_DT */
		ora.S,   /* DT_DD */
		ora.S,   /* DT_HH */
		ora.S,   /* PTO_PATH */
		ora.S,   /* NOTI_CONT */
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* LIVE_SN */
		ora.S,   /* BRD_GBN_CD */
		ora.I64, /* SN */
		ora.S,   /* LD_YN */
		ora.S,   /* INQ_GBN_CD */
		ora.S,   /* INQ_REG_DT */
		ora.S,   /* MSG_END_YN */
		ora.S,   /* LIVE_NVN_YN */
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

	rtnNotificationList := models.RtnNotificationList{}
	notificationList := make([]models.NotificationList, 0)

	var (
		ntMemNo       string
		ntDvcId       string
		ntNotiKndCd   string
		ntRegDt       string
		ntCfrmDt      string
		ntDtDd        string
		ntDtHh        string
		ntPtoPath     string
		ntNotiCont    string
		ntEntpMemNo   string
		ntRecrutSn    string
		ntPpMemNo     string
		ntLiveSn      string
		ntBrdGbnCd    string
		ntSn          int64
		ntLdYn        string
		ntInqGbnCd    string
		ntInqRegDt    string
		ntMsgEndYn    string
		ntLiveNvnYn   string
		ntFullPtoPath string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			ntMemNo = procRset.Row[0].(string)
			ntDvcId = procRset.Row[1].(string)
			ntNotiKndCd = procRset.Row[2].(string)
			ntRegDt = procRset.Row[3].(string)
			ntCfrmDt = procRset.Row[4].(string)
			ntDtDd = procRset.Row[5].(string)
			ntDtHh = procRset.Row[6].(string)
			ntPtoPath = procRset.Row[7].(string)
			if ntPtoPath == "" {
				ntFullPtoPath = ntPtoPath
			} else {
				ntFullPtoPath = imgServer + ntPtoPath
			}
			ntNotiCont = procRset.Row[8].(string)
			ntEntpMemNo = procRset.Row[9].(string)
			ntRecrutSn = procRset.Row[10].(string)
			ntPpMemNo = procRset.Row[11].(string)
			ntLiveSn = procRset.Row[12].(string)
			ntBrdGbnCd = procRset.Row[13].(string)
			ntSn = procRset.Row[14].(int64)
			ntLdYn = procRset.Row[15].(string)
			ntInqGbnCd = procRset.Row[16].(string)
			ntInqRegDt = procRset.Row[17].(string)
			ntMsgEndYn = procRset.Row[18].(string)
			ntLiveNvnYn = procRset.Row[19].(string)

			notificationList = append(notificationList, models.NotificationList{
				NtMemNo:     ntMemNo,
				NtDvcId:     ntDvcId,
				NtNotiKndCd: ntNotiKndCd,
				NtRegDt:     ntRegDt,
				NtCfrmDt:    ntCfrmDt,
				NtDtDd:      ntDtDd,
				NtDtHh:      ntDtHh,
				NtPtoPath:   ntFullPtoPath,
				NtNotiCont:  ntNotiCont,
				NtEntpMemNo: ntEntpMemNo,
				NtRecrutSn:  ntRecrutSn,
				NtPpMemNo:   ntPpMemNo,
				NtLiveSn:    ntLiveSn,
				NtBrdGbnCd:  ntBrdGbnCd,
				NtSn:        ntSn,
				NtLdYn:      ntLdYn,
				NtInqGbnCd:  ntInqGbnCd,
				NtInqRegDt:  ntInqRegDt,
				NtMsgEndYn:  ntMsgEndYn,
				NtLiveNvnYn: ntLiveNvnYn,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnNotificationList = models.RtnNotificationList{
			RtnNotificationListData: notificationList,
		}
	}
	// End : Notification PopUp List

	c.Data["json"] = &rtnNotificationList
	c.ServeJSON()
}
