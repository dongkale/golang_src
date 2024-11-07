package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type LiveNvNRecuruitApplyListController struct {
	BaseController
}

func (c *LiveNvNRecuruitApplyListController) Post() {

	session := c.StartSession()

	// mem_id := session.Get(c.Ctx.Request.Context(), "mem_id")
	// if mem_id == nil {
	// 	c.Ctx.Redirect(302, "/login")
	// }

	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		//c.Ctx.Redirect(302, "/login")

		c.Data["json"] = &models.RtnRecruitSubList{}
		c.ServeJSON()
		return
	}

	pLang, _ := beego.AppConfig.String("lang")
	imgServer, _  := beego.AppConfig.String("viewpath")

	pEntpMemNo := mem_no
	pRecrutSn := c.GetString("recruit_sn") // 구분코드(A:전체, I:채용중, E:종료)

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Recruit Stat List
	pLiveSn := ""
	pOffSet := "0"
	pLimit := "100"
	pEvlPrgsStat := "00"
	pSortGbn := "01"
	pLiveReqStatCd := "01"
	pApplySortCd := "01"
	pApplySortWay := "DESC"

	fmt.Printf(fmt.Sprintf("CALL ZSP_LIVE_NVN_APPL_ADD_LIST('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pRecrutSn, pLiveSn, pEvlPrgsStat, pSortGbn, pLiveReqStatCd, pApplySortCd, pApplySortWay))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_NVN_APPL_ADD_LIST('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
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
	procRset := &ora.Rset{}
	_, err = stmtProcCall.Exe(procRset)

	if err != nil {
		panic(err)
	}

	rtnrecruitApplyList := models.RtnLiveNvnRecruitApplyList{}
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

		rtnrecruitApplyList = models.RtnLiveNvnRecruitApplyList{
			RtnLiveNvnRecruitApplyListData: recruitApplyList,
		}
	}
	//End : Recruit Apply List

	c.Data["json"] = &rtnrecruitApplyList
	c.ServeJSON()
}
