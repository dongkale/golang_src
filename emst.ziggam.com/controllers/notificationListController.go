package controllers

import (
	"fmt"

	"emst.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
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
		c.Ctx.Redirect(302, "/common/login")
	}

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := mem_no //"E2018102500001"
	imgServer, _ := beego.AppConfig.String("viewpath")

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

	log.Debug("CALL SP_EMS_NOTIFICATION_LIST_R('%v','%v', :1)",
		pLang, pEntpMemNo)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_NOTIFICATION_LIST_R('%v','%v', :1)",
		pLang, pEntpMemNo),
		ora.I64, /* GRP_NO */
		ora.S,   /* MEM_NO */
		ora.S,   /* REG_DY */
		ora.S,   /* REG_HM */
		ora.S,   /* PTO_PATH */
		ora.S,   /* NOTI_KND_CD */
		ora.S,   /* NOTI_CONT */
		ora.S,   /* SEX */
		ora.S,   /* NM */
		ora.S,   /* AGE */
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* UP_JOB_GRP */
		ora.S,   /* JOB_GRP */
		ora.S,   /* CFRM_DT */
		ora.I64, /* SN */
		ora.S,   /* REG_DT */
		ora.I64, /* NEXT_GRP_NO */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* NEW_YN */
		ora.I64, /* PREV_GRP_NO */
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
		grpNo     int64
		memNo     string
		regDy     string
		regHm     string
		ptoPath   string
		notiKndCd string
		notiCont  string
		sex       string
		nm        string
		age       string
		entpMemNo string
		recrutSn  string
		upJobGrp  string
		jobGrp    string
		cfrmDt    string
		sn        int64
		regDt     string
		nextGrpNo int64
		ppMemNo   string
		newYn     string
		prevGrpNo int64

		fullPtoPath string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			grpNo = procRset.Row[0].(int64)
			memNo = procRset.Row[1].(string)
			regDy = procRset.Row[2].(string)
			regHm = procRset.Row[3].(string)
			ptoPath = procRset.Row[4].(string)

			if ptoPath == "" {
				fullPtoPath = ptoPath
			} else {
				fullPtoPath = imgServer + ptoPath
			}

			notiKndCd = procRset.Row[5].(string)
			notiCont = procRset.Row[6].(string)
			sex = procRset.Row[7].(string)
			nm = procRset.Row[8].(string)
			age = procRset.Row[9].(string)
			entpMemNo = procRset.Row[10].(string)
			recrutSn = procRset.Row[11].(string)
			upJobGrp = procRset.Row[12].(string)
			jobGrp = procRset.Row[13].(string)
			cfrmDt = procRset.Row[14].(string)
			sn = procRset.Row[15].(int64)
			regDt = procRset.Row[16].(string)
			nextGrpNo = procRset.Row[17].(int64)
			ppMemNo = procRset.Row[18].(string)
			newYn = procRset.Row[19].(string)
			prevGrpNo = procRset.Row[20].(int64)

			notificationList = append(notificationList, models.NotificationList{
				GrpNo:     grpNo,
				MemNo:     memNo,
				RegDy:     regDy,
				RegHm:     regHm,
				PtoPath:   fullPtoPath,
				NotiKndCd: notiKndCd,
				NotiCont:  notiCont,
				Sex:       sex,
				Nm:        nm,
				Age:       age,
				EntpMemNo: entpMemNo,
				RecrutSn:  recrutSn,
				UpJobGrp:  upJobGrp,
				JobGrp:    jobGrp,
				CfrmDt:    cfrmDt,
				Sn:        sn,
				RegDt:     regDt,
				NextGrpNo: nextGrpNo,
				PpMemNo:   ppMemNo,
				NewYn:     newYn,
				PrevGrpNo: prevGrpNo,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Notification List

	c.Data["NotificationList"] = notificationList

	c.TplName = "notification/notification_list.html"
}
