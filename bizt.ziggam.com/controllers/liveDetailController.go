package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	ora "gopkg.in/rana/ora.v4"
)

type LiveDetailController struct {
	BaseController
}

func (c *LiveDetailController) Get() {

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

	pEntpMemNo := mem_no //c.GetString("entp_mem_no")
	pRecrutSn := c.GetString("recrut_sn")
	pPpMemNo := c.GetString("pp_mem_no")
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

	// Start : Live Interview Detail info
	fmt.Printf(fmt.Sprintf("CALL ZSP_LIVE_DTL_INFO_R('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_DTL_INFO_R('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo),
		ora.S,   /* EVL_PRGS_STAT_CD */
		ora.S,   /* LIVE_REQ_STAT_CD */
		ora.S,   /* REG_DT */
		ora.S,   /* PTO_PATH */
		ora.S,   /* NM */
		ora.S,   /* SEX */
		ora.I64, /* AGE */
		ora.S,   /* UP_JOB_GRP */
		ora.S,   /* JOB_GRP */
		ora.S,   /* RECRUT_TITLE */
		ora.S,   /* LIVE_SN */
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

	liveDetail := make([]models.LiveDetail, 0)

	var (
		evlPrgsStatCd string
		liveReqStatCd string
		regDt         string
		ptoPath       string
		nm            string
		sex           string
		age           int64
		upJobGrp      string
		jobGrp        string
		recrutTitle   string
		liveSn        string
		fullPtoPath   string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			evlPrgsStatCd = procRset.Row[0].(string)
			liveReqStatCd = procRset.Row[1].(string)
			regDt = procRset.Row[2].(string)
			ptoPath = procRset.Row[3].(string)

			if ptoPath == "" {
				fullPtoPath = ptoPath
			} else {
				fullPtoPath = imgServer + ptoPath
			}
			nm = procRset.Row[4].(string)
			sex = procRset.Row[5].(string)
			age = procRset.Row[6].(int64)
			upJobGrp = procRset.Row[7].(string)
			jobGrp = procRset.Row[8].(string)
			recrutTitle = procRset.Row[9].(string)
			liveSn = procRset.Row[10].(string)

			liveDetail = append(liveDetail, models.LiveDetail{
				EvlPrgsStatCd: evlPrgsStatCd,
				LiveReqStatCd: liveReqStatCd,
				RegDt:         regDt,
				PtoPath:       fullPtoPath,
				Nm:            nm,
				Sex:           sex,
				Age:           age,
				UpJobGrp:      upJobGrp,
				JobGrp:        jobGrp,
				RecrutTitle:   recrutTitle,
				LiveSn:        liveSn,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Live Interview Detail info

	// Start : live member list
	var (
		lmPpChrgGbnCd string
		lmPpChrgNm    string
		lmPpChrgBpNm  string
		lmChrgSn      string
	)

	// Start : live member list
	fmt.Printf(fmt.Sprintf("CALL ZSP_LIVE_MEM_LIST_R('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo, liveSn))

	stmtProcCallMem, errMem := ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_MEM_LIST_R('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo, liveSn),
		ora.S, /* PP_CHRG_GBN_CD */
		ora.S, /* PP_CHRG_NM */
		ora.S, /* PP_CHRG_BP_NM */
		ora.S, /* PP_CHRG_SN */
	)
	defer stmtProcCallMem.Close()
	if errMem != nil {
		panic(errMem)
	}
	procRsetMem := &ora.Rset{}
	_, errMem = stmtProcCallMem.Exe(procRsetMem)

	if errMem != nil {
		panic(errMem)
	}

	liveMemList := make([]models.LiveMemList, 0)

	if procRsetMem.IsOpen() {
		for procRsetMem.Next() {
			lmPpChrgGbnCd = procRsetMem.Row[0].(string)
			lmPpChrgNm = procRsetMem.Row[1].(string)
			lmPpChrgBpNm = procRsetMem.Row[2].(string)
			lmChrgSn = procRsetMem.Row[3].(string)

			liveMemList = append(liveMemList, models.LiveMemList{
				LmPpChrgGbnCd: lmPpChrgGbnCd,
				LmPpChrgNm:    lmPpChrgNm,
				LmPpChrgBpNm:  lmPpChrgBpNm,
				LmChrgSn:      lmChrgSn,
			})
		}
		if errMem := procRsetMem.Err(); errMem != nil {
			panic(errMem)
		}
	}
	// End : live member list

	pGbnCd := "L"
	// Start : live history list
	var (
		lhEntpMemNo       string
		lhRecrutSn        string
		lhPpMemNo         string
		lhLiveStatCd      string
		lhMsgGbnCd        string
		lhMSgSn           string
		lhMsgYn           string
		lhLiveSn          string
		lhLiveItvSday     string
		lhLiveItvSTime    string
		lhLiveItvEday     string
		lhLiveItvETime    string
		lhMsgGbnNm        string
		lhMsgEndYn        string
		lhMemGbn          string
		lhLiveItvCnclDay  string
		lhLiveItvCnclTime string
		lhNMsgGbnCd       string
	)

	fmt.Printf(fmt.Sprintf("CALL ZSP_LIVE_HISTORY_LIST_R('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo, pGbnCd))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_HISTORY_LIST_R('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo, pGbnCd),
		ora.S, /* ENTP_MEM_NO */
		ora.S, /* RECRUT_SN */
		ora.S, /* PP_MEM_NO */
		ora.S, /* LIVE_STAT_CD */
		ora.S, /* MSG_GBN_CD */
		ora.S, /* MSG_SN */
		ora.S, /* MSG_YN */
		ora.S, /* LIVE_SN */
		ora.S, /* LIVE_ITV_SDAY */
		ora.S, /* LIVE_ITV_STIME */
		ora.S, /* LIVE_ITV_EDAY */
		ora.S, /* LIVE_ITV_ETIME */
		ora.S, /* MSG_GBN_NM */
		ora.S, /* MSG_END_YN */
		ora.S, /* MEM_GBN */
		ora.S, /* LIVE_ITV_CNCL_DAY */
		ora.S, /* LIVE_ITV_CNCL_TIME */
		ora.S, /* N_MSG_GBN_CD */
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

	liveHistoryList := make([]models.LiveHistoryList, 0)

	if procRset.IsOpen() {
		for procRset.Next() {
			lhEntpMemNo = procRset.Row[0].(string)
			lhRecrutSn = procRset.Row[1].(string)
			lhPpMemNo = procRset.Row[2].(string)
			lhLiveStatCd = procRset.Row[3].(string)
			lhMsgGbnCd = procRset.Row[4].(string)
			lhMSgSn = procRset.Row[5].(string)
			lhMsgYn = procRset.Row[6].(string)
			lhLiveSn = procRset.Row[7].(string)
			lhLiveItvSday = procRset.Row[8].(string)
			lhLiveItvSTime = procRset.Row[9].(string)
			lhLiveItvEday = procRset.Row[10].(string)
			lhLiveItvETime = procRset.Row[11].(string)
			lhMsgGbnNm = procRset.Row[12].(string)
			lhMsgEndYn = procRset.Row[13].(string)
			lhMemGbn = procRset.Row[14].(string)
			lhLiveItvCnclDay = procRset.Row[15].(string)
			lhLiveItvCnclTime = procRset.Row[16].(string)
			lhNMsgGbnCd = procRset.Row[17].(string)

			liveHistoryList = append(liveHistoryList, models.LiveHistoryList{
				LhEntpMemNo:       lhEntpMemNo,
				LhRecrutSn:        lhRecrutSn,
				LhPpMemNo:         lhPpMemNo,
				LhLiveStatCd:      lhLiveStatCd,
				LhMsgGbnCd:        lhMsgGbnCd,
				LhMSgSn:           lhMSgSn,
				LhMsgYn:           lhMsgYn,
				LhLiveSn:          lhLiveSn,
				LhLiveItvSday:     lhLiveItvSday,
				LhLiveItvSTime:    lhLiveItvSTime,
				LhLiveItvEday:     lhLiveItvEday,
				LhLiveItvETime:    lhLiveItvETime,
				LhMsgGbnNm:        lhMsgGbnNm,
				LhMsgEndYn:        lhMsgEndYn,
				LhMemGbn:          lhMemGbn,
				LhLiveItvCnclDay:  lhLiveItvCnclDay,
				LhLiveItvCnclTime: lhLiveItvCnclTime,
				LhNMsgGbnCd:       lhNMsgGbnCd,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : live history list

	pGbnCd = "I"
	// Start : live history info
	var (
		lhiEntpMemNo    string
		lhiRecrutSn     string
		lhiPpMemNo      string
		lhiLiveStatCd   string
		lhiMsgGbnCd     string
		lhiMSgSn        string
		lhiMsgYn        string
		lhiLiveSn       string
		lhiLiveItvSday  string
		lhiLiveItvSTime string
		lhiLiveItvEday  string
		lhiLiveItvETime string
		lhiMsgGbnNm     string
		lhiMsgEndYn     string
	)

	fmt.Printf(fmt.Sprintf("CALL ZSP_LIVE_HISTORY_LIST_R('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo, pGbnCd))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_HISTORY_LIST_R('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo, pGbnCd),
		ora.S, /* ENTP_MEM_NO */
		ora.S, /* RECRUT_SN */
		ora.S, /* PP_MEM_NO */
		ora.S, /* LIVE_STAT_CD */
		ora.S, /* MSG_GBN_CD */
		ora.S, /* MSG_SN */
		ora.S, /* MSG_YN */
		ora.S, /* LIVE_SN */
		ora.S, /* LIVE_ITV_SDAY */
		ora.S, /* LIVE_ITV_STIME */
		ora.S, /* LIVE_ITV_EDAY */
		ora.S, /* LIVE_ITV_ETIME */
		ora.S, /* MSG_GBN_NM */
		ora.S, /* MSG_END_YN */
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

	liveHistoryInfo := make([]models.LiveHistoryInfo, 0)

	if procRset.IsOpen() {
		for procRset.Next() {
			lhiEntpMemNo = procRset.Row[0].(string)
			lhiRecrutSn = procRset.Row[1].(string)
			lhiPpMemNo = procRset.Row[2].(string)
			lhiLiveStatCd = procRset.Row[3].(string)
			lhiMsgGbnCd = procRset.Row[4].(string)
			lhiMSgSn = procRset.Row[5].(string)
			lhiMsgYn = procRset.Row[6].(string)
			lhiLiveSn = procRset.Row[7].(string)
			lhiLiveItvSday = procRset.Row[8].(string)
			lhiLiveItvSTime = procRset.Row[9].(string)
			lhiLiveItvEday = procRset.Row[10].(string)
			lhiLiveItvETime = procRset.Row[11].(string)
			lhiMsgGbnNm = procRset.Row[12].(string)
			lhiMsgEndYn = procRset.Row[13].(string)

			liveHistoryInfo = append(liveHistoryInfo, models.LiveHistoryInfo{
				LhiEntpMemNo:    lhiEntpMemNo,
				LhiRecrutSn:     lhiRecrutSn,
				LhiPpMemNo:      lhiPpMemNo,
				LhiLiveStatCd:   lhiLiveStatCd,
				LhiMsgGbnCd:     lhiMsgGbnCd,
				LhiMSgSn:        lhiMSgSn,
				LhiMsgYn:        lhiMsgYn,
				LhiLiveSn:       lhiLiveSn,
				LhiLiveItvSday:  lhiLiveItvSday,
				LhiLiveItvSTime: lhiLiveItvSTime,
				LhiLiveItvEday:  lhiLiveItvEday,
				LhiLiveItvETime: lhiLiveItvETime,
				LhiMsgGbnNm:     lhiMsgGbnNm,
				LhiMsgEndYn:     lhiMsgEndYn,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : live history info

	c.Data["LhiEntpMemNo"] = lhiEntpMemNo
	c.Data["LhiRecrutSn"] = lhiRecrutSn
	c.Data["LhiPpMemNo"] = lhiPpMemNo
	c.Data["LhiLiveStatCd"] = lhiLiveStatCd
	c.Data["LhiMsgGbnCd"] = lhiMsgGbnCd
	c.Data["LhiMSgSn"] = lhiMSgSn
	c.Data["LhiMsgYn"] = lhiMsgYn
	c.Data["LhiLiveSn"] = lhiLiveSn
	c.Data["LhiLiveItvSday"] = lhiLiveItvSday
	c.Data["LhiLiveItvSTime"] = lhiLiveItvSTime
	c.Data["LhiLiveItvEday"] = lhiLiveItvEday
	c.Data["LhiLiveItvETime"] = lhiLiveItvETime
	c.Data["LhiMsgGbnNm"] = lhiMsgGbnNm
	c.Data["LhiMsgEndYn"] = lhiMsgEndYn

	c.Data["EvlPrgsStatCd"] = evlPrgsStatCd
	c.Data["LiveReqStatCd"] = liveReqStatCd
	c.Data["RegDt"] = regDt
	c.Data["PtoPath"] = fullPtoPath
	c.Data["Nm"] = nm
	c.Data["Sex"] = sex
	c.Data["Age"] = age
	c.Data["UpJobGrp"] = upJobGrp
	c.Data["JobGrp"] = jobGrp
	c.Data["RecrutTitle"] = recrutTitle
	c.Data["RecrutSn"] = pRecrutSn
	c.Data["PpMemNo"] = pPpMemNo

	c.Data["LiveMemList"] = liveMemList
	c.Data["LiveHistoryList"] = liveHistoryList

	c.Data["TMenuId"] = "L00"
	c.Data["SMenuId"] = "L02"

	c.TplName = "live/live_detail.html"
}
