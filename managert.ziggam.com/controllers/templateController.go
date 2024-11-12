package controllers

import (
	"github.com/beego/beego/v2/core/logs"
)

type TemplateController struct {
	BaseController
}

func (c *TemplateController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	//session := c.StartSession()
	//mem_no := session.Get("mem_no")
	//if mem_no == nil {
	//	c.Ctx.Redirect(302, "/login")
	//}
	//mem_sn := session.Get("mem_sn")

	//pLang, _ := beego.AppConfig.String("lang")
	//pEntpMemNo := mem_no
	//pPpChrgSn := mem_sn
	//pGbnCd := c.GetString("gbn_cd") //구분코드
	//if pGbnCd == "" {
	//	pGbnCd = "Y"
	//}
	//imgServer, _ := beego.AppConfig.String("viewpath")

	// Start : Oracle DB Connection
	//env, srv, ses, err := GetRawConnection()
	//defer env.Close()
	//defer srv.Close()
	//defer ses.Close()
	//if err != nil {
	//	panic(err)
	//}
	//// End : Oracle DB Connection
	//
	//// Start : Main Info
	//
	//log.Debug("CALL ZSP_MAIN_INFO_R('%v', '%v', '%v', :1)",
	//	pLang, pEntpMemNo, pPpChrgSn)
	//
	//stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_MAIN_INFO_R('%v', '%v', '%v', :1)",
	//	pLang, pEntpMemNo, pPpChrgSn),
	//	ora.S,   /* PP_CHRG_BP_NM */
	//	ora.S,   /* PP_CHRG_NM */
	//	ora.S,   /* ENTP_KO_NM */
	//	ora.I64, /* INFO_CNT */
	//	ora.I64, /* VD_CNT */
	//)
	//defer stmtProcCall.Close()
	//if err != nil {
	//	panic(err)
	//}
	//procRset := &ora.Rset{}
	//_, err = stmtProcCall.Exe(procRset)
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//mainInfo := make([]models.MainInfo, 0)
	//
	//var (
	//	mnPpChrgBpNm string
	//	mnPpChrgNm   string
	//	mnEntpKoNm   string
	//	mnInfoCnt    int64
	//	mnVdCnt      int64
	//)
	//
	//if procRset.IsOpen() {
	//	for procRset.Next() {
	//		mnPpChrgBpNm = procRset.Row[0].(string)
	//		mnPpChrgNm = procRset.Row[1].(string)
	//		mnEntpKoNm = procRset.Row[2].(string)
	//		mnInfoCnt = procRset.Row[3].(int64)
	//		mnVdCnt = procRset.Row[4].(int64)
	//
	//		mainInfo = append(mainInfo, models.MainInfo{
	//			MnPpChrgBpNm: mnPpChrgBpNm,
	//			MnPpChrgNm:   mnPpChrgNm,
	//			MnEntpKoNm:   mnEntpKoNm,
	//			MnInfoCnt:    mnInfoCnt,
	//			MnVdCnt:      mnVdCnt,
	//		})
	//	}
	//	if err := procRset.Err(); err != nil {
	//		panic(err)
	//	}
	//}
	//// End : Main Infoc
	//
	//// Start : Main Stat
	//
	//log.Debug("CALL ZSP_MAIN_STATS_R('%v', '%v', :1)",
	//	pLang, pEntpMemNo)
	//
	//stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_MAIN_STATS_R('%v','%v', :1)",
	//	pLang, pEntpMemNo),
	//	ora.I64, /* VIDEO_TODAY_CNT */
	//	ora.I64, /* VIDEO_TOT_CNT */
	//	ora.I64, /* RECRUT_ING_CNT */
	//	ora.I64, /* RECRUT_TOT_CNT */
	//	ora.I64, /* APPLY_TODAY_CNT */
	//	ora.I64, /* APPLY_TOT_CNT */
	//)
	//defer stmtProcCall.Close()
	//if err != nil {
	//	panic(err)
	//}
	//procRset = &ora.Rset{}
	//_, err = stmtProcCall.Exe(procRset)
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//mainStat := make([]models.MainStat, 0)
	//
	//var (
	//	videoTodayCnt int64
	//	videoTotCnt   int64
	//	recrutIngCnt  int64
	//	recrutTotCnt  int64
	//	applyTodayCnt int64
	//	applyTotCnt   int64
	//)
	//
	//if procRset.IsOpen() {
	//	for procRset.Next() {
	//		videoTodayCnt = procRset.Row[0].(int64)
	//		videoTotCnt = procRset.Row[1].(int64)
	//		recrutIngCnt = procRset.Row[2].(int64)
	//		recrutTotCnt = procRset.Row[3].(int64)
	//		applyTodayCnt = procRset.Row[4].(int64)
	//		applyTotCnt = procRset.Row[5].(int64)
	//
	//		mainStat = append(mainStat, models.MainStat{
	//			VideoTodayCnt: videoTodayCnt,
	//			VideoTotCnt:   videoTotCnt,
	//			RecrutIngCnt:  recrutIngCnt,
	//			RecrutTotCnt:  recrutTotCnt,
	//			ApplyTodayCnt: applyTodayCnt,
	//			ApplyTotCnt:   applyTotCnt,
	//		})
	//	}
	//	if err := procRset.Err(); err != nil {
	//		panic(err)
	//	}
	//}
	//// End : Main Stat
	//
	//// Start : Live Interview List
	//
	//log.Debug("CALL ZSP_MAIN_LIVE_LIST_R('%v', '%v', :1)",
	//	pLang, pEntpMemNo)
	//
	//stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_MAIN_LIVE_LIST_R('%v', '%v', :1)",
	//	pLang, pEntpMemNo),
	//	ora.S,   /* RECRUT_SN */
	//	ora.S,   /* PP_MEM_NO */
	//	ora.S,   /* LIVE_SN */
	//	ora.S,   /* LIVE_ITV_SD */
	//	ora.S,   /* LIVE_ITV_ST */
	//	ora.S,   /* PTO_PATH */
	//	ora.S,   /* NM */
	//	ora.S,   /* SEX */
	//	ora.I64, /* AGE */
	//)
	//defer stmtProcCall.Close()
	//if err != nil {
	//	panic(err)
	//}
	//procRset = &ora.Rset{}
	//_, err = stmtProcCall.Exe(procRset)
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//mainLiveList := make([]models.MainLiveList, 0)
	//
	//var (
	//	mnlRecrutSn    string
	//	mnlPpMemNo     string
	//	mnlLiveSn      string
	//	mnlLiveItvSd   string
	//	mnlLiveItvSt   string
	//	mnlPtoPath     string
	//	mnlNm          string
	//	mnlSex         string
	//	mnlAge         int64
	//	mnlFullPtoPath string
	//)
	//
	//if procRset.IsOpen() {
	//	for procRset.Next() {
	//		mnlRecrutSn = procRset.Row[0].(string)
	//		mnlPpMemNo = procRset.Row[1].(string)
	//		mnlLiveSn = procRset.Row[2].(string)
	//		mnlLiveItvSd = procRset.Row[3].(string)
	//		mnlLiveItvSt = procRset.Row[4].(string)
	//		mnlPtoPath = procRset.Row[5].(string)
	//		if mnlPtoPath == "" {
	//			mnlFullPtoPath = mnlPtoPath
	//		} else {
	//			mnlFullPtoPath = imgServer + mnlPtoPath
	//		}
	//		mnlNm = procRset.Row[6].(string)
	//		mnlSex = procRset.Row[7].(string)
	//		mnlAge = procRset.Row[8].(int64)
	//
	//		mainLiveList = append(mainLiveList, models.MainLiveList{
	//			MnlRecrutSn:  mnlRecrutSn,
	//			MnlPpMemNo:   mnlPpMemNo,
	//			MnlLiveSn:    mnlLiveSn,
	//			MnlLiveItvSd: mnlLiveItvSd,
	//			MnlLiveItvSt: mnlLiveItvSt,
	//			MnlPtoPath:   mnlFullPtoPath,
	//			MnlNm:        mnlNm,
	//			MnlSex:       mnlSex,
	//			MnlAge:       mnlAge,
	//		})
	//	}
	//	if err := procRset.Err(); err != nil {
	//		panic(err)
	//	}
	//}
	//// End : Live Interview List
	//
	//// Start : Recruit List
	//
	//log.Debug("CALL ZSP_MAIN_RECRUIT_LIST_R('%v', '%v', :1)",
	//	pLang, pEntpMemNo)
	//
	//stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_MAIN_RECRUIT_LIST_R('%v', '%v', :1)",
	//	pLang, pEntpMemNo),
	//	ora.S,   /* ENTP_MEM_NO */
	//	ora.S,   /* RECRUT_SN */
	//	ora.S,   /* RECRUT_TITLE */
	//	ora.I64, /* NEW_CNT */
	//	ora.S,   /* DDY */
	//)
	//defer stmtProcCall.Close()
	//if err != nil {
	//	panic(err)
	//}
	//procRset = &ora.Rset{}
	//_, err = stmtProcCall.Exe(procRset)
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//mainRecruitList := make([]models.MainRecruitList, 0)
	//
	//var (
	//	mnrEntpMemNo   string
	//	mnrRecrutSn    string
	//	mnrRecrutTitle string
	//	mnrNewCnt      int64
	//	mnrDdy         string
	//)
	//
	//if procRset.IsOpen() {
	//	for procRset.Next() {
	//		mnrEntpMemNo = procRset.Row[0].(string)
	//		mnrRecrutSn = procRset.Row[1].(string)
	//		mnrRecrutTitle = procRset.Row[2].(string)
	//		mnrNewCnt = procRset.Row[3].(int64)
	//		mnrDdy = procRset.Row[4].(string)
	//
	//		mainRecruitList = append(mainRecruitList, models.MainRecruitList{
	//			MnrEntpMemNo:   mnrEntpMemNo,
	//			MnrRecrutSn:    mnrRecrutSn,
	//			MnrRecrutTitle: mnrRecrutTitle,
	//			MnrNewCnt:      mnrNewCnt,
	//			MnrDdy:         mnrDdy,
	//		})
	//	}
	//	if err := procRset.Err(); err != nil {
	//		panic(err)
	//	}
	//}
	//// End : Recruit List
	//
	//// Start : Applicant List
	//
	//log.Debug("CALL ZSP_MAIN_APLY_LIST_R('%v', '%v', '%v', :1)",
	//	pLang, pEntpMemNo, pGbnCd)
	//
	//stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_MAIN_APLY_LIST_R('%v', '%v', '%v', :1)",
	//	pLang, pEntpMemNo, pGbnCd),
	//	ora.S,   /* RECRUT_SN */
	//	ora.S,   /* PP_MEM_NO */
	//	ora.S,   /* FAVR_APLY_PP_YN */
	//	ora.S,   /* PTO_PATH */
	//	ora.S,   /* NM */
	//	ora.S,   /* SEX */
	//	ora.I64, /* AGE */
	//	ora.S,   /* REG_DT */
	//)
	//defer stmtProcCall.Close()
	//if err != nil {
	//	panic(err)
	//}
	//procRset = &ora.Rset{}
	//_, err = stmtProcCall.Exe(procRset)
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//applicantList := make([]models.ApplicantList, 0)
	//
	//var (
	//	mnaRecrutSn     string
	//	mnaPpMemNo      string
	//	mnaFavrAplyPpYn string
	//	mnaPtoPath      string
	//	mnaNm           string
	//	mnaSex          string
	//	mnaAge          int64
	//	mnaRegDt        string
	//	mnaFullPtoPath  string
	//)
	//
	//if procRset.IsOpen() {
	//	for procRset.Next() {
	//		mnaRecrutSn = procRset.Row[0].(string)
	//		mnaPpMemNo = procRset.Row[1].(string)
	//		mnaFavrAplyPpYn = procRset.Row[2].(string)
	//		mnaPtoPath = procRset.Row[3].(string)
	//		if mnaPtoPath == "" {
	//			mnaFullPtoPath = mnaPtoPath
	//		} else {
	//			mnaFullPtoPath = imgServer + mnaPtoPath
	//		}
	//		mnaNm = procRset.Row[4].(string)
	//		mnaSex = procRset.Row[5].(string)
	//		mnaAge = procRset.Row[6].(int64)
	//		mnaRegDt = procRset.Row[7].(string)
	//
	//		applicantList = append(applicantList, models.ApplicantList{
	//			MnaRecrutSn:     mnaRecrutSn,
	//			MnaPpMemNo:      mnaPpMemNo,
	//			MnaFavrAplyPpYn: mnaFavrAplyPpYn,
	//			MnaPtoPath:      mnaFullPtoPath,
	//			MnaNm:           mnaNm,
	//			MnaSex:          mnaSex,
	//			MnaAge:          mnaAge,
	//			MnaRegDt:        mnaRegDt,
	//		})
	//	}
	//	if err := procRset.Err(); err != nil {
	//		panic(err)
	//	}
	//}
	//// End : Applicant List
	//
	//// Start : Main Notice
	//
	//log.Debug("CALL ZSP_MAIN_NOTICE_TOP_R('%v', :1)",
	//	pLang)
	//
	//stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_MAIN_NOTICE_TOP_R('%v', :1)",
	//	pLang),
	//	ora.I64, /* SN */
	//	ora.S,   /* TITLE */
	//	ora.S,   /* REG_DT */
	//)
	//defer stmtProcCall.Close()
	//if err != nil {
	//	panic(err)
	//}
	//procRset = &ora.Rset{}
	//_, err = stmtProcCall.Exe(procRset)
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//nainNotice := make([]models.MainNotice, 0)
	//
	//var (
	//	mnnSn    int64
	//	mnnTitle string
	//	mnnRegDt string
	//)
	//
	//if procRset.IsOpen() {
	//	for procRset.Next() {
	//		mnnSn = procRset.Row[0].(int64)
	//		mnnTitle = procRset.Row[1].(string)
	//		mnnRegDt = procRset.Row[2].(string)
	//
	//		nainNotice = append(nainNotice, models.MainNotice{
	//			MnnSn:    mnnSn,
	//			MnnTitle: mnnTitle,
	//			MnnRegDt: mnnRegDt,
	//		})
	//	}
	//	if err := procRset.Err(); err != nil {
	//		panic(err)
	//	}
	//}
	//// End : Main Notice
	//c.Data["MnnSn"] = mnnSn
	//c.Data["MnnTitle"] = mnnTitle
	//c.Data["MnnRegDt"] = mnnRegDt
	//
	//c.Data["MnPpChrgBpNm"] = mnPpChrgBpNm
	//c.Data["MnPpChrgNm"] = mnPpChrgNm
	//c.Data["MnEntpKoNm"] = mnEntpKoNm
	//c.Data["MnInfoCnt"] = mnInfoCnt
	//c.Data["MnVdCnt"] = mnVdCnt
	//
	//c.Data["VideoTodayCnt"] = videoTodayCnt
	//c.Data["VideoTotCnt"] = videoTotCnt
	//c.Data["RecrutIngCnt"] = recrutIngCnt
	//c.Data["RecrutTotCnt"] = recrutTotCnt
	//c.Data["ApplyTodayCnt"] = applyTodayCnt
	//c.Data["ApplyTotCnt"] = applyTotCnt
	//
	//c.Data["MainLiveList"] = mainLiveList
	//c.Data["MainRecruitList"] = mainRecruitList
	//c.Data["ApplicantList"] = applicantList
	//
	//c.Data["TMenuId"] = "M00"
	//c.Data["SMenuId"] = "M00"

	c.TplName = "main/template.html"
}
