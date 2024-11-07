package controllers

import (
	"encoding/json"
	"fmt"

	"bizt.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

// LiveItvNvNReqPopupController ...
type LiveNvNReqPopupController struct {
	BaseController
}

// Get ...
func (c *LiveNvNReqPopupController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/login")
	}

	mem_sn := session.Get(c.Ctx.Request.Context(), "mem_sn")

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := mem_no
	pRecrutSn := c.GetString("recrut_sn")
	pApplyMemList := c.GetString("apply_mem_list")

	pRecrutChange := c.GetString("recrut_change")
	pReqToDetail := c.GetString("req_to_detail")

	// pPpMemNo := c.GetString("pp_mem_no")
	// pPpMemNm := c.GetString("pp_mem_nm")
	// pPtoPath := c.GetString("pto_path")
	// pLiveSn := c.GetString("live_sn")

	fmt.Printf(pApplyMemList)
	fmt.Printf(pRecrutChange)
	fmt.Printf(pReqToDetail)

	// 리스트 Json 데이터 -> array 로
	var retApplyMemList []models.LiveNvNApplyMemList

	err := json.Unmarshal([]byte(pApplyMemList), &retApplyMemList)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", retApplyMemList)

	imgServer, _  := beego.AppConfig.String("viewpath")
	// cdnPath := beego.AppConfig.String("cdnpath")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Entp Team Member List

	pGbnCd := "A"
	pLiveSn := ""

	fmt.Printf(fmt.Sprintf("CALL ZSP_LIVE_NVN_ENTP_ADD_LIST('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pLiveSn, pGbnCd))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_NVN_ENTP_ADD_LIST('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pLiveSn, pGbnCd),
		ora.I64, /* TOT_CNT */
		ora.S,   /* PP_CHRG_SN */
		ora.S,   /* PP_CHRG_GBN_CD */
		ora.S,   /* PP_CHRG_NM */
		ora.S,   /* PP_CHRG_BP_NM */
		ora.S,   /* EMAIL */
		ora.S,   /* ENTP_MEM_ID */
		ora.S,   /* PP_CHRG_TEL_NO */
		ora.I64, /* ROWNO */
		ora.I64, /* MEM_REG_CNT */
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

	entpTeamMemberList := make([]models.LiveNvnEntpTeamMemberList, 0)

	var (
		etTotCnt      int64
		etPpChrgSn    string
		etPpChrgGbnCd string
		etPpChrgNm    string
		etPpChrgBpNm  string
		etEmail       string
		etEntpMemId   string
		etPpChrgTelNo string
		etRowNo       int64
		etMemRegCnt   int64
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

			etRowNo = procRset.Row[8].(int64)
			etMemRegCnt = procRset.Row[9].(int64)

			entpTeamMemberList = append(entpTeamMemberList, models.LiveNvnEntpTeamMemberList{
				EtTotCnt:      etTotCnt,
				EtPpChrgSn:    etPpChrgSn,
				EtPpChrgGbnCd: etPpChrgGbnCd,
				EtPpChrgNm:    etPpChrgNm,
				EtPpChrgBpNm:  etPpChrgBpNm,
				EtEmail:       etEmail,
				EtEntpMemId:   etEntpMemId,
				EtPpChrgTelNo: etPpChrgTelNo,
				EtRowNo:       etRowNo,
				EtMemRegCnt:   etMemRegCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Entp Team Member List

	// Start : invite recruit List
	// pGbnCd := c.GetString("gbn_cd") // 구분코드(A:전체, I:채용중, E:종료)

	// if pGbnCd == "" {
	// 	pGbnCd = "02"
	// }

	pGbnCd = "02"
	pOffSet := "0"
	pLimit := "100"
	pSortGbnR := "03"

	fmt.Printf(fmt.Sprintf("CALL ZSP_RECRUIT_LIST_TYPE('%v', %v, %v, '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, "", pGbnCd, pSortGbnR))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_RECRUIT_LIST_TYPE('%v', %v, %v, '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, "", pGbnCd, pSortGbnR),
		ora.I64, /* TOT_CNT */
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* PRGS_STAT */
		ora.S,   /* RECRUT_TITLE */
		ora.S,   /* UP_JOB_GRP */
		ora.S,   /* JOB_GRP */
		ora.S,   /* RECRUT_EDT */
		ora.S,   /* SDY */
		ora.S,   /* EDY */
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

	recruitSubList := make([]models.RecruitSubList, 0)

	var (
		sTotCnt      int64
		sEntpMemNo   string
		sRecrutSn    string
		sPrgsStat    string
		sRecrutTitle string
		sUpJobGrp    string
		sJobGrp      string
		//sTrimRecrutTitle string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			sTotCnt = procRset.Row[0].(int64)
			sEntpMemNo = procRset.Row[1].(string)
			sRecrutSn = procRset.Row[2].(string)
			sPrgsStat = procRset.Row[3].(string)
			sRecrutTitle = procRset.Row[4].(string)
			sUpJobGrp = procRset.Row[5].(string)
			sJobGrp = procRset.Row[6].(string)

			recruitSubList = append(recruitSubList, models.RecruitSubList{
				STotCnt:      sTotCnt,
				SEntpMemNo:   sEntpMemNo,
				SRecrutSn:    sRecrutSn,
				SPrgsStat:    sPrgsStat,
				SRecrutTitle: sRecrutTitle,
				SUpJobGrp:    sUpJobGrp,
				SJobGrp:      sJobGrp,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : invite recruit List

	// Start : Recruit Stat List
	pLiveSn = ""
	pOffSet = "0"
	pLimit = "100"
	pEvlPrgsStat := "00"
	pSortGbn := "01"
	pLiveReqStatCd := "01"
	pApplySortCd := "01"
	pApplySortWay := "DESC"

	fmt.Printf(fmt.Sprintf("CALL ZSP_LIVE_NVN_APPL_ADD_LIST('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pRecrutSn, pLiveSn, pEvlPrgsStat, pSortGbn, pLiveReqStatCd, pApplySortCd, pApplySortWay))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_NVN_APPL_ADD_LIST('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pRecrutSn, pLiveSn, pEvlPrgsStat, pSortGbn, pLiveReqStatCd, pApplySortCd, pApplySortWay),
		ora.I64, /* TOT_CNT */
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* NM */
		ora.S,   /* SEX */
		ora.S,   /* AGE */
		ora.S,   /* REG_DT*/
		ora.S,   /* APPLY_DT */
		ora.S,   /* EVL_STAT_DT*/
		ora.S,   /* EVL_PRGS_STAT_CD */
		ora.S,   /* RCRT_APLY_STAT_CD */
		ora.S,   /* ENTP_CFRM_YN */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* LIVE_REQ_STAT_CD */
		ora.I64, /* ROWNO */
		ora.S,   /* PTO_PATH */
		ora.S,   /* DCMNT_EVL_STAT_CD */
		ora.S,   /* ONWY_INTRV_EVL_STAT_CD */
		ora.S,   /* LIVE_INTRV_EVL_STAT_CD */
		ora.S,   /* READ_END_DAY */
		ora.I64, /* APPLY_REG_CNT */ // pLiveSn = "" 이므로 무조건 0
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

	recruitApplyList := make([]models.LiveNvnRecruitApplyList, 0)

	var (
		rslTotCnt          int64
		rslEntpMemNo       string
		rslRecrutSn        string
		rslNm              string
		rslSex             string
		rslAge             string
		rslRegDt           string
		rslApplyDt         string
		rslEvlStatDt       string
		rslEvlPrgsStatCd   string
		rslRcrtAplyStatCd  string
		rslEntpCfrmYn      string
		rslPpMemNo         string
		rslLiveReqStatCd   string
		rslRowNo           int64
		rslPtoPath         string
		fullPtoPath        string
		dcmntEvlStatCd     string
		onwyIntrvEvlStatCd string
		liveIntrvEvlStatCd string
		readEndDay         string
		rslApplyRegCnt     int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			rslTotCnt = procRset.Row[0].(int64)
			rslEntpMemNo = procRset.Row[1].(string)
			rslRecrutSn = procRset.Row[2].(string)
			rslNm = procRset.Row[3].(string)
			rslSex = procRset.Row[4].(string)
			rslAge = procRset.Row[5].(string)
			rslRegDt = procRset.Row[6].(string)
			rslApplyDt = procRset.Row[7].(string)
			rslEvlStatDt = procRset.Row[8].(string)
			rslEvlPrgsStatCd = procRset.Row[9].(string)
			rslRcrtAplyStatCd = procRset.Row[10].(string)
			rslEntpCfrmYn = procRset.Row[11].(string)
			//rslVpYn = procRset.Row[13].(string)
			rslPpMemNo = procRset.Row[12].(string)
			rslLiveReqStatCd = procRset.Row[13].(string)
			rslRowNo = procRset.Row[14].(int64)
			rslPtoPath = procRset.Row[15].(string)

			if rslPtoPath == "" {
				fullPtoPath = rslPtoPath
			} else {
				fullPtoPath = imgServer + rslPtoPath
			}

			dcmntEvlStatCd = procRset.Row[16].(string)
			onwyIntrvEvlStatCd = procRset.Row[17].(string)
			liveIntrvEvlStatCd = procRset.Row[18].(string)

			readEndDay = procRset.Row[19].(string)    // "N" 상 "Y" 90일이 넘었므로 비정상
			rslApplyRegCnt = procRset.Row[20].(int64) // pLiveSn 있으면 해당 라이브에 포함된 지원자수

			recruitApplyList = append(recruitApplyList, models.LiveNvnRecruitApplyList{
				RslTotCnt:          rslTotCnt,
				RslEntpMemNo:       rslEntpMemNo,
				RslRecrutSn:        rslRecrutSn,
				RslNm:              rslNm,
				RslSex:             rslSex,
				RslAge:             rslAge,
				RslRegDt:           rslRegDt,
				RslApplyDt:         rslApplyDt,
				RslEvlStatDt:       rslEvlStatDt,
				RslEvlPrgsStatCd:   rslEvlPrgsStatCd,
				RslRcrtAplyStatCd:  rslRcrtAplyStatCd,
				RslEntpCfrmYn:      rslEntpCfrmYn,
				RslPpMemNo:         rslPpMemNo,
				RslLiveReqStatCd:   rslLiveReqStatCd,
				RslRowNo:           rslRowNo,
				RslPtoPath:         fullPtoPath,
				DcmntEvlStatCd:     dcmntEvlStatCd,
				OnwyIntrvEvlStatCd: onwyIntrvEvlStatCd,
				LiveIntrvEvlStatCd: liveIntrvEvlStatCd,
				ReadEndDay:         readEndDay,
				RslApplyRegCnt:     rslApplyRegCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	//End : Recruit Apply List

	c.Data["EntpTeamMemberList"] = entpTeamMemberList
	c.Data["EntpMemNo"] = pEntpMemNo
	c.Data["RecrutSn"] = pRecrutSn
	c.Data["RecruitApplyList"] = recruitApplyList

	c.Data["RecruitList"] = recruitSubList

	c.Data["ApplyMemList"] = retApplyMemList

	c.Data["RecruitChange"] = pRecrutChange
	c.Data["ReqToDetail"] = pReqToDetail

	// c.Data["PpMemNo"] = pPpMemNo
	// c.Data["PpMemNm"] = pPpMemN
	// c.Data["PtoPath"] = pPtoPat
	// c.Data["pLiveSn"] = pLiveSn

	c.Data["SMemSn"] = mem_sn // 맴버PP_CHGS
	//c.Data["LiveMemList"] = liveMemList

	c.TplName = "live/live_nvn_req_popup.html"
}
