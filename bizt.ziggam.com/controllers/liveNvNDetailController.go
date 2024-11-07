package controllers

import (
	"fmt"
	"time"

	"bizt.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

// LiveNvNDetailController ...
type LiveNvNDetailController struct {
	BaseController
}

// Get ...
func (c *LiveNvNDetailController) Get() {	
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	session := c.StartSession()

	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/login")
	}

	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := mem_no
	pRecrutSn := c.GetString("recrut_sn")
	//pPpMemNo := c.GetString("pp_mem_no")
	pLiveSn := c.GetString("live_sn")

	pConfirmPopup := c.GetString("confirm_popup")

	imgServer, _  := beego.AppConfig.String("viewpath")
	nvnLinkPath,_ := beego.AppConfig.String("nvn_link_path")
	
	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : live detail info
	log.Debug(fmt.Sprintf("CALL ZSP_LIVE_NVN_DETAIL_INFO('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pLiveSn))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_NVN_DETAIL_INFO('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pLiveSn),
		ora.S,   /* RECRUT_SN */
		ora.S,   /* RECRUT_TITLE */
		ora.S,   /* LIVE_SN */
		ora.S,   /* LIVE_ITV_SDT */
		ora.S,   /* LIVE_ITV_EDT */
		ora.S,   /* LIVE_ITV_SDAY */
		ora.S,   /* LIVE_ITV_STIME */
		ora.S,   /* LIVE_ITV_EDAY */
		ora.S,   /* LIVE_ITV_ETIME */
		ora.S,   /* LIVE_REG_DAY */
		ora.S,   /* LIVE_REG_TIME */
		ora.S,   /* LIVE_STAT_CD */
		ora.S,   /* PP_CHRG_NM */
		ora.S,   /* PP_CHRG_BP_NM */
		ora.S,   /* MEM_LIST */
		ora.S,   /* APPLY_LIST */
		ora.I64, /* LIVE_ITV_JOIN_CNT */
		ora.S,   /* REG_DT */
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

	var (
		recrutSn       string
		recrutTitle    string
		liveSn         string
		liveItvSdt     string
		liveItvEdt     string
		liveItvSday    string
		liveItvStime   string
		liveItvEday    string
		liveItvEtime   string
		liveItvRday    string
		liveItvRtime   string
		liveStatCd     string
		ppChrgNm       string
		ppChrgBpNm     string
		memList        string
		applyList      string
		liveItvJoinCnt int64
		regDt          string
		liveItvSdtFmt  string
		liveItvEdtFmt  string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			recrutSn = procRset.Row[0].(string)
			recrutTitle = procRset.Row[1].(string)
			liveSn = procRset.Row[2].(string)
			liveItvSdt = procRset.Row[3].(string)
			liveItvEdt = procRset.Row[4].(string)
			liveItvSday = procRset.Row[5].(string)
			liveItvStime = procRset.Row[6].(string)
			liveItvEday = procRset.Row[7].(string)
			liveItvEtime = procRset.Row[8].(string)
			liveItvRday = procRset.Row[9].(string)
			liveItvRtime = procRset.Row[10].(string)
			liveStatCd = procRset.Row[11].(string)
			ppChrgNm = procRset.Row[12].(string)
			ppChrgBpNm = procRset.Row[13].(string)
			memList = procRset.Row[14].(string)
			applyList = procRset.Row[15].(string)
			liveItvJoinCnt = procRset.Row[16].(int64)
			regDt = procRset.Row[17].(string)

			log.Debug(fmt.Sprintf("recrutSn:%v, recrutTitle:%v, liveSn:%v, liveItvSdt:%v, liveItvEdt:%v, liveItvSday:%v, liveItvStime:%v, liveItvEday:%v, liveItvEtime:%v, liveItvRday:%v, liveItvRtime:%v, liveStatCd:%v, ppChrgNm:%v, ppChrgBpNm:%v, memList:%v, applyList:%v, liveItvJoinCnt:%v, regDt:%v",
				recrutSn, recrutTitle, liveSn, liveItvSdt, liveItvEdt, liveItvSday, liveItvStime, liveItvEday, liveItvEtime, liveItvRday, liveItvRtime, liveStatCd, ppChrgNm, ppChrgBpNm, memList, applyList, liveItvJoinCnt, regDt))
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}

	convTime, _ := time.Parse("20060102150405", liveItvSdt)
	liveItvSdtFmt = fmt.Sprintf(convTime.Format("2006.01.02 15:04"))

	if liveItvEdt != "" {
		convTime, _ = time.Parse("20060102150405", liveItvEdt)
		liveItvEdtFmt = fmt.Sprintf(convTime.Format("2006.01.02 15:04"))
	}

	// End : live history info

	// Start : live member list
	var (
		lmPpChrgGbnCd  string
		lmPpChrgNm     string
		lmPpChrgBpNm   string
		lmPpChrgSn     string
		lmPpLiveStatCd string
	)

	// Start : live member list
	log.Debug(fmt.Sprintf("CALL ZSP_LIVE_NVN_ENTP_MEM_LIST('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pLiveSn))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_NVN_ENTP_MEM_LIST('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pLiveSn),
		ora.S, /* PP_CHRG_GBN_CD */
		ora.S, /* PP_CHRG_NM */
		ora.S, /* PP_CHRG_BP_NM */
		ora.S, /* PP_CHRG_SN */
		ora.S, /* LIVE_STAT_CD */
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

	liveNvnMemList := make([]models.LiveNvNMemList, 0)

	if procRset.IsOpen() {
		for procRset.Next() {
			lmPpChrgGbnCd = procRset.Row[0].(string)
			lmPpChrgNm = procRset.Row[1].(string)
			lmPpChrgBpNm = procRset.Row[2].(string)
			lmPpChrgSn = procRset.Row[3].(string)
			lmPpLiveStatCd = procRset.Row[4].(string)

			liveNvnMemList = append(liveNvnMemList, models.LiveNvNMemList{
				LmLiveSchedStatCd: liveStatCd,

				LmPpChrgGbnCd:  lmPpChrgGbnCd,
				LmPpChrgNm:     lmPpChrgNm,
				LmPpChrgBpNm:   lmPpChrgBpNm,
				LmPpChrgSn:     lmPpChrgSn,
				LmPpLiveStatCd: lmPpLiveStatCd,
			})
		}
		if errMem := procRset.Err(); errMem != nil {
			panic(errMem)
		}
	}
	// End : live member list

	// Start : live apply member list
	var (
		lmPpMemNo            string
		lmName               string
		lmSex                string
		lmAge                string
		lmPtoPath            string
		lmPtofullPath        string
		lmRecrutSn           string
		lmLiveStatCd         string
		lmTRC04LiveReqStatCd string
		lmTRC04LiveSn        string
		lmMsgYn              string
		lmMsgEndYn           string
		lmReadEndDay         string
		lmItvLink            string
		lmItvLinkfullPath    string
	)

	log.Debug(fmt.Sprintf("CALL ZSP_LIVE_NVN_APPL_MEM_LIST('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pLiveSn))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_NVN_APPL_MEM_LIST('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pLiveSn),
		ora.S, /* PP_MEM_NO */
		ora.S, /* NAME */
		ora.S, /* SEX */
		ora.S, /* AGE */
		ora.S, /* PTO_PATH */
		ora.S, /* RECRUT_SN */
		ora.S, /* LIVE_STAT_CD */
		ora.S, /* TRC04.LIVE_REQ_STAT_CD */
		ora.S, /* TRC04.LIVE_SN */
		ora.S, /* MSG_YN */
		ora.S, /* MSG_END_YN */
		ora.S, /* READ_END_DAY */ // 90 체크
		ora.S, /* ITV_LINK */     // PP_LINK

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

	liveNvNApplyList := make([]models.LiveNvNApplyList, 0)

	if procRset.IsOpen() {
		for procRset.Next() {
			lmPpMemNo = procRset.Row[0].(string)
			lmName = procRset.Row[1].(string)
			lmSex = procRset.Row[2].(string)
			lmAge = procRset.Row[3].(string)
			lmPtoPath = procRset.Row[4].(string)
			lmRecrutSn = procRset.Row[5].(string)
			lmLiveStatCd = procRset.Row[6].(string)
			lmTRC04LiveReqStatCd = procRset.Row[7].(string)
			lmTRC04LiveSn = procRset.Row[8].(string)
			lmMsgYn = procRset.Row[9].(string)
			lmMsgEndYn = procRset.Row[10].(string)
			lmReadEndDay = procRset.Row[11].(string)
			lmItvLink = procRset.Row[12].(string)

			if lmPtoPath == "" {
				lmPtofullPath = lmPtoPath
			} else {
				lmPtofullPath = imgServer + lmPtoPath
			}

			if nvnLinkPath == "" {
				lmItvLinkfullPath = ""
			} else {
				lmItvLinkfullPath = nvnLinkPath + lmItvLink
			}

			liveNvNApplyList = append(liveNvNApplyList, models.LiveNvNApplyList{
				LmLiveSchedStatCd:    liveStatCd,
				LmPpMemNo:            lmPpMemNo,
				LmNm:                 lmName,
				LmSex:                lmSex,
				LmAge:                lmAge,
				LmPtoPath:            lmPtofullPath,
				LmRecrutSn:           lmRecrutSn,
				LmLiveStatCd:         lmLiveStatCd,
				LmTRC04LiveReqStatCd: lmTRC04LiveReqStatCd,
				LmTRC04LiveSn:        lmTRC04LiveSn,
				LmMsgYn:              lmMsgYn,
				LmMsgEndYn:           lmMsgEndYn,
				LmReadEndDay:         lmReadEndDay,
				LmItvLink:            lmItvLinkfullPath,
			})
		}
		if errMem := procRset.Err(); errMem != nil {
			panic(errMem)
		}
	}
	// End : live apply member list

	// Start : live history list
	var (
		lhEntpMemNo   string
		lhRecrutSn    string
		lhLiveSn      string
		lhPpMemNo     string
		lhPpMemNm     string
		lhLiveStatCd  string
		lhMsgGbnCd    string
		lhMSgSn       string
		lhMsgYn       string
		lhMsgGbnNm    string
		lhMsgEndYn    string
		lhMemGbn      string
		lhMsgRegDtFmt string
		lhNMsgGbnCd   string
		lhPpChrgGbnCd string
		lhPpChrgSn    string
		lhPpChrgNm    string
		lhPpChrgBpNm  string
		lhRegDt       string
	)

	log.Debug(fmt.Sprintf("CALL ZSP_LIVE_NVN_HISTORY_LIST('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pLiveSn))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_NVN_HISTORY_LIST('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pLiveSn),
		ora.S, /* ENTP_MEM_NO */
		ora.S, /* RECRUT_SN */
		ora.S, /* LIVE_SN */
		ora.S, /* PP_MEM_NO */
		ora.S, /* PP_MEM_NM */
		ora.S, /* LIVE_STAT_CD */
		ora.S, /* MSG_GBN_CD */
		ora.S, /* MSG_SN */
		ora.S, /* MSG_YN */
		ora.S, /* MSG_GBN_NM */
		ora.S, /* MSG_END_YN */
		ora.S, /* MEM_GBN */
		ora.S, /* MSG_REG_DT_FMT */
		ora.S, /* N_MSG_GBN_CD */
		ora.S, /* PP_CHRG_GBN_CD */
		ora.S, /* PP_CHRG_SN */
		ora.S, /* PP_CHRG_NM */
		ora.S, /* PP_CHRG_BP_NM */
		ora.S, /* REG_DT */
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

	liveHistoryList := make([]models.LiveNvnHistoryList, 0)

	if procRset.IsOpen() {
		for procRset.Next() {
			lhEntpMemNo = procRset.Row[0].(string)
			lhRecrutSn = procRset.Row[1].(string)
			lhLiveSn = procRset.Row[2].(string)
			lhPpMemNo = procRset.Row[3].(string)
			lhPpMemNm = procRset.Row[4].(string)
			lhLiveStatCd = procRset.Row[5].(string)
			lhMsgGbnCd = procRset.Row[6].(string)
			lhMSgSn = procRset.Row[7].(string)
			lhMsgYn = procRset.Row[8].(string)
			lhMsgGbnNm = procRset.Row[9].(string)
			lhMsgEndYn = procRset.Row[10].(string)
			lhMemGbn = procRset.Row[11].(string)
			lhMsgRegDtFmt = procRset.Row[12].(string)
			lhNMsgGbnCd = procRset.Row[13].(string)
			lhPpChrgGbnCd = procRset.Row[14].(string)
			lhPpChrgSn = procRset.Row[15].(string)
			lhPpChrgNm = procRset.Row[16].(string)
			lhPpChrgBpNm = procRset.Row[17].(string)
			lhRegDt = procRset.Row[18].(string)

			liveHistoryList = append(liveHistoryList, models.LiveNvnHistoryList{
				LhEntpMemNo:   lhEntpMemNo,
				LhRecrutSn:    lhRecrutSn,
				LhLiveSn:      lhLiveSn,
				LhPpMemNo:     lhPpMemNo,
				LhPpMemNm:     lhPpMemNm,
				LhLiveStatCd:  lhLiveStatCd,
				LhMsgGbnCd:    lhMsgGbnCd,
				LhMSgSn:       lhMSgSn,
				LhMsgYn:       lhMsgYn,
				LhMsgGbnNm:    lhMsgGbnNm,
				LhMsgEndYn:    lhMsgEndYn,
				LhMemGbn:      lhMemGbn,
				LhMsgRegDtFmt: lhMsgRegDtFmt,
				LhNMsgGbnCd:   lhNMsgGbnCd,
				LhPpChrgGbnCd: lhPpChrgGbnCd,
				LhPpChrgSn:    lhPpChrgSn,
				LhPpChrgNm:    lhPpChrgNm,
				LhPpChrgBpNm:  lhPpChrgBpNm,
				LhRegDt:       lhRegDt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : live history list

	c.Data["EntpMemNo"] = pEntpMemNo
	c.Data["RecrutSn"] = pRecrutSn
	c.Data["LiveSn"] = pLiveSn

	c.Data["RecrutTitle"] = recrutTitle

	c.Data["LiveStatCd"] = liveStatCd

	c.Data["LiveItvSdt"] = liveItvSdt
	c.Data["LiveItvEdt"] = liveItvEdt

	c.Data["LiveItvSday"] = liveItvSday
	c.Data["LiveItvStime"] = liveItvStime

	c.Data["LiveItvEday"] = liveItvEday
	c.Data["LiveItvEtime"] = liveItvEtime

	c.Data["LiveItvSdtFmt"] = liveItvSdtFmt
	c.Data["LiveItvEdtFmt"] = liveItvEdtFmt

	c.Data["LiveItvRday"] = liveItvRday
	c.Data["LiveItvRtime"] = liveItvRtime

	c.Data["LiveItvJoinCnt"] = liveItvJoinCnt

	c.Data["LiveHistoryList"] = liveHistoryList

	c.Data["PpChrgNm"] = ppChrgNm
	c.Data["PpChrgBpNm"] = ppChrgBpNm

	c.Data["LhiMsgGbnCd"] = "99"

	c.Data["LiveApplyList"] = liveNvNApplyList
	c.Data["LiveMemList"] = liveNvnMemList

	c.Data["ConfirmPopup"] = pConfirmPopup

	//c.Data["LiveApplyListJS"] = &liveNvNApplyList
	// doc, _ := json.Marshal(liveNvNApplyList)
	// c.Data["LiveApplyListJS"] = string(doc)

	c.Data["TMenuId"] = "L00"
	c.Data["SMenuId"] = "L01"

	c.TplName = "live/live_nvn_detail.html"
}
