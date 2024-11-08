package controllers

import (
	"fmt"

	"emst.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	ora "gopkg.in/rana/ora.v4"
)

type MainController struct {
	BaseController
}

func (c *MainController) Get() {

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

	// Start : Notice List

	log.Debug("CALL SP_EMS_MAIN_NOTICE_LIST_R('%v', :1)",
		pLang)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_MAIN_NOTICE_LIST_R('%v', :1)",
		pLang),
		ora.S,   /* REG_DT */
		ora.S,   /* TITLE */
		ora.I64, /* SN */
		ora.S,   /* NEW_YN */
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

	mainNotiList := make([]models.MainNotiList, 0)

	var (
		notiRegDt string
		notiTitle string
		notiSn    int64
		notiNewYn string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			notiRegDt = procRset.Row[0].(string)
			notiTitle = procRset.Row[1].(string)
			notiSn = procRset.Row[2].(int64)
			notiNewYn = procRset.Row[3].(string)

			mainNotiList = append(mainNotiList, models.MainNotiList{
				NotiRegDt: notiRegDt,
				NotiTitle: notiTitle,
				NotiSn:    notiSn,
				NotiNewYn: notiNewYn,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Notice List

	// Start : Main Stat

	log.Debug("CALL SP_EMS_MAIN_STATS_R('%v', '%v', :1)",
		pLang, pEntpMemNo)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_MAIN_STATS_R('%v','%v', :1)",
		pLang, pEntpMemNo),
		ora.I64, /* VIDEO_TODAY_CNT */
		ora.I64, /* VIDEO_TOT_CNT */
		ora.I64, /* RECRUT_ING_CNT */
		ora.I64, /* RECRUT_TOT_CNT */
		ora.I64, /* APPLY_TODAY_CNT */
		ora.I64, /* APPLY_TOT_CNT */
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

	mainStat := make([]models.MainStat, 0)

	var (
		videoTodaycnt int64
		videoTotCnt   int64
		recrutIngCnt  int64
		recrutTotCnt  int64
		applyTodayCnt int64
		applyTotCnt   int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			videoTodaycnt = procRset.Row[0].(int64)
			videoTotCnt = procRset.Row[1].(int64)
			recrutIngCnt = procRset.Row[2].(int64)
			recrutTotCnt = procRset.Row[3].(int64)
			applyTodayCnt = procRset.Row[4].(int64)
			applyTotCnt = procRset.Row[5].(int64)

			mainStat = append(mainStat, models.MainStat{
				VideoTodaycnt: videoTodaycnt,
				VideoTotCnt:   videoTotCnt,
				RecrutIngCnt:  recrutIngCnt,
				RecrutTotCnt:  recrutTotCnt,
				ApplyTodayCnt: applyTodayCnt,
				ApplyTotCnt:   applyTotCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Main Stat

	// Start : Recruit List

	log.Debug("CALL SP_EMS_MAIN_RECRUIT_LIST_R('%v', '%v', :1)",
		pLang, pEntpMemNo)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_MAIN_RECRUIT_LIST_R('%v', '%v', :1)",
		pLang, pEntpMemNo),
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* RECRUT_TITLE */
		ora.S,   /* EMPL_TYP */
		ora.S,   /* UP_JOB_GRP */
		ora.S,   /* JOB_GRP */
		ora.S,   /* RECRUT_DY */
		ora.S,   /* REG_DT */
		ora.I64, /* NEW_CNT */
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

	mainRecruitList := make([]models.MainRecruitList, 0)

	var (
		rcEntpMemNo   string
		rcRecrutSn    string
		rcRecrutTitle string
		rcEmplTyp     string
		rcUpJobGrp    string
		rcJobGrp      string
		rcRecrutDy    string
		rcRegDt       string
		rcNewCnt      int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			rcEntpMemNo = procRset.Row[0].(string)
			rcRecrutSn = procRset.Row[1].(string)
			rcRecrutTitle = procRset.Row[2].(string)
			rcEmplTyp = procRset.Row[3].(string)
			rcUpJobGrp = procRset.Row[4].(string)
			rcJobGrp = procRset.Row[5].(string)
			rcRecrutDy = procRset.Row[6].(string)
			rcRegDt = procRset.Row[7].(string)
			rcNewCnt = procRset.Row[8].(int64)

			mainRecruitList = append(mainRecruitList, models.MainRecruitList{
				RcEntpMemNo:   rcEntpMemNo,
				RcRecrutSn:    rcRecrutSn,
				RcRecrutTitle: rcRecrutTitle,
				RcEmplTyp:     rcEmplTyp,
				RcUpJobGrp:    rcUpJobGrp,
				RcJobGrp:      rcJobGrp,
				RcRecrutDy:    rcRecrutDy,
				RcRegDt:       rcRegDt,
				RcNewCnt:      rcNewCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Recruit List

	// Start : Apply List

	log.Debug("CALL SP_EMS_MAIN_APPLY_LIST_R('%v', '%v', :1)",
		pLang, pEntpMemNo)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_MAIN_APPLY_LIST_R('%v', '%v', :1)",
		pLang, pEntpMemNo),
		ora.S, /* ENTP_MEM_NO */
		ora.S, /* RECRUT_SN */
		ora.S, /* PP_MEM_NO */
		ora.S, /* PTO_PATH */
		ora.S, /* NM */
		ora.S, /* SEX */
		ora.S, /* RECRUT_TITLE */
		ora.S, /* LEFT_DY */
		ora.S, /* RED_YN */
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

	mainApplyList := make([]models.MainApplyList, 0)

	var (
		apEntpMemNo   string
		apRecrutSn    string
		apPpMemNo     string
		apPtoPath     string
		apNm          string
		apSex         string
		apRecrutTitle string
		apLeftDy      string
		apRedYn       string
		fullPtoPath   string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			apEntpMemNo = procRset.Row[0].(string)
			apRecrutSn = procRset.Row[1].(string)
			apPpMemNo = procRset.Row[2].(string)
			apPtoPath = procRset.Row[3].(string)

			if apPtoPath == "" {
				fullPtoPath = apPtoPath
			} else {
				fullPtoPath = imgServer + apPtoPath
			}

			apNm = procRset.Row[4].(string)
			apSex = procRset.Row[5].(string)
			apRecrutTitle = procRset.Row[6].(string)
			apLeftDy = procRset.Row[7].(string)
			apRedYn = procRset.Row[8].(string)

			mainApplyList = append(mainApplyList, models.MainApplyList{
				ApEntpMemNo:   apEntpMemNo,
				ApRecrutSn:    apRecrutSn,
				ApPpMemNo:     apPpMemNo,
				ApPtoPath:     fullPtoPath,
				ApNm:          apNm,
				ApSex:         apSex,
				ApRecrutTitle: apRecrutTitle,
				ApLeftDy:      apLeftDy,
				ApRedYn:       apRedYn,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Apply List

	c.Data["MainNotiList"] = mainNotiList
	c.Data["VideoTodaycnt"] = videoTodaycnt
	c.Data["VideoTotCnt"] = videoTotCnt
	c.Data["RecrutIngCnt"] = recrutIngCnt
	c.Data["RecrutTotCnt"] = recrutTotCnt
	c.Data["ApplyTodayCnt"] = applyTodayCnt
	c.Data["ApplyTotCnt"] = applyTotCnt
	c.Data["MainRecruitList"] = mainRecruitList
	c.Data["MainApplyList"] = mainApplyList

	c.TplName = "main/main.html"
}
