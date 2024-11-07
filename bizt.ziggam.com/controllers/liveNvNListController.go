package controllers

import (
	"fmt"
	"math"
	"strconv"

	"bizt.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

// LiveNvNListController ...
type LiveNvNListController struct {
	BaseController
}

// Get ...
func (c *LiveNvNListController) Get() {

	session := c.StartSession()

	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/login")
		return
	}

	mem_sn := session.Get(c.Ctx.Request.Context(), "mem_sn")
	if mem_sn == nil {
		c.Ctx.Redirect(302, "/login")
		return
	}

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := mem_no.(string)
	pPChrgSn := mem_sn.(string)
	//pStatCd := c.GetString("stat_cd")     //상태코드
	//pKeyword := c.GetString("keyword")    //검색어
	// pViewType := c.GetString("view_type") //보기구분
	// if pViewType == "" {
	// 	pViewType = "G"
	// }
	//pGbnCd1 := c.GetString("gbn_cd1") //구분값(요청)
	//pGbnCd2 := c.GetString("gbn_cd2") //구분값(예정)
	//pGbnCd3 := c.GetString("gbn_cd3") //구분값(종료)
	//pGbnCd4 := c.GetString("gbn_cd4") //구분값(거절.취소)
	//imgServer, _  := beego.AppConfig.String("viewpath")

	pGbnRecrutSn := c.GetString("recrut_sn")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Live Stat
	fmt.Printf(fmt.Sprintf("CALL ZSP_LIVE_NVN_STATS('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pPChrgSn, pGbnRecrutSn))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_NVN_STATS('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pPChrgSn, pGbnRecrutSn),
		ora.I64, /* WAIT_CNT */
		ora.I64, /* ING_CNT */
		ora.I64, /* END_CNT */
		ora.I64, /* CNCL_CNT */
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

	liveStat := make([]models.LiveStat, 0)

	var (
		waitCnt int64
		ingCnt  int64
		endCnt  int64
		cnclCnt int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			waitCnt = procRset.Row[0].(int64)
			ingCnt = procRset.Row[1].(int64)
			endCnt = procRset.Row[2].(int64)
			cnclCnt = procRset.Row[3].(int64)

			liveStat = append(liveStat, models.LiveStat{
				WaitCnt: waitCnt,
				IngCnt:  ingCnt,
				EndCnt:  endCnt,
				CnclCnt: cnclCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Live Stat

	c.Data["WaitCnt"] = waitCnt
	c.Data["IngCnt"] = ingCnt
	c.Data["EndCnt"] = endCnt
	c.Data["CnclCnt"] = cnclCnt

	//	c.Data["LiveList01"] = liveList01
	//	c.Data["LiveList02"] = liveList02
	//	c.Data["LiveList03"] = liveList03

	c.Data["TMenuId"] = "L00"
	c.Data["SMenuId"] = "L01"

	c.TplName = "live/live_nvn_list.html"
}

// Post ...
func (c *LiveNvNListController) Post() {

	session := c.StartSession()

	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Data["json"] = &models.RtnLiveNvNList{}
		c.ServeJSON()
		return
	}

	mem_sn := session.Get(c.Ctx.Request.Context(), "mem_sn")
	if mem_sn == nil {
		c.Data["json"] = &models.RtnLiveNvNList{}
		c.ServeJSON()
		return
	}

	pLang, _ := beego.AppConfig.String("lang")
	imgServer, _  := beego.AppConfig.String("viewpath")

	pEntpMemNo := mem_no.(string)
	pPChrgSn := mem_sn.(string)
	pStatCd := c.GetString("stat_cd")     //상태코드
	pKeyword := c.GetString("keyword")    //검색어
	pViewType := c.GetString("view_type") //보기구분
	pGbnCd1 := c.GetString("gbn_cd1")     //구분값(요청)
	pGbnCd2 := c.GetString("gbn_cd2")     //구분값(예정)
	pGbnCd3 := c.GetString("gbn_cd3")     //구분값(종료)
	pGbnCd4 := c.GetString("gbn_cd4")     //구분값(거절.취소)

	pSortCd := c.GetString("sort_cd") // 01: 시작 시간, 02: 요청일
	if pSortCd == "" {
		pSortCd = "02"
	}

	pSortWay := c.GetString("sort_way")
	if pSortWay == "" {
		pSortWay = "DESC"
	}

	pGbnRecrutSn := c.GetString("recrut_sn")
	pGbnPpChrgSn := c.GetString("pp_chrg_sn")

	pGbnSSdy := c.GetString("sdt_sdy")
	pGbnSEdy := c.GetString("sdt_edy")

	if pGbnSSdy == "" {
		pGbnSSdy = "20200101"
	}

	if pGbnSEdy == "" {
		pGbnSEdy = "20251201"
		// nowTime := time.Now()
		// pGbnEdy = nowTime.Format("20060102")
	}

	// 종료 및 취소
	pGbnESdy := c.GetString("edt_sdy")
	pGbnEEdy := c.GetString("edt_edy")

	if pGbnESdy == "" {
		pGbnESdy = "20200101"
	}

	if pGbnEEdy == "" {
		pGbnEEdy = "20251201"
	}

	// 요청일
	pGbnRSdy := c.GetString("rdt_sdy")
	pGbnREdy := c.GetString("rdt_edy")

	if pGbnRSdy == "" {
		pGbnRSdy = "20200101"
	}

	if pGbnREdy == "" {
		pGbnREdy = "20251201"
	}

	if pGbnPpChrgSn != "" {
		pPChrgSn = pGbnPpChrgSn
	}

	fmt.Printf(fmt.Sprintf("pEntpMemNo:%v, pPChrgSn: %v, pStatCd:%v, pKeyword:%v, pViewType:%v, pGbnCd1:%v, pGbnCd2:%v, pGbnCd3:%v, pGbnCd4:%v, pSortCd:%v, pSortWay:%v, pGbnRecrutSn:%v, pGbnPpChrgSn:%v, pGbnSSdy:%v, pGbnSEdy:%v, pGbnESdy:%v, pGbnEEdy:%v",
		pEntpMemNo, pPChrgSn, pStatCd, pKeyword, pViewType, pGbnCd1, pGbnCd2, pGbnCd3, pGbnCd4, pSortCd, pSortWay, pGbnRecrutSn, pGbnPpChrgSn, pGbnSSdy, pGbnSEdy, pGbnESdy, pGbnEEdy))

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Live Interview List recruit list
	fmt.Printf(fmt.Sprintf("CALL ZSP_LIVE_NVN_LIST_RECRUIT('%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pPChrgSn, pGbnCd1, pGbnCd2, pGbnCd3, pGbnCd4))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_NVN_LIST_RECRUIT('%v', '%v', '%v', '%v', '%v', '%v', '%v',:1)",
		pLang, pEntpMemNo, pPChrgSn, pGbnCd1, pGbnCd2, pGbnCd3, pGbnCd4),
		ora.I64, /* TOT_CNT */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* RECRUT_TITLE */
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

	recruitList := make([]models.LiveNvnListRecruitList, 0)

	var (
		sTotCnt      int64
		sRecrutSn    string
		sRecrutTitle string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			sTotCnt = procRset.Row[0].(int64)
			sRecrutSn = procRset.Row[1].(string)
			sRecrutTitle = procRset.Row[2].(string)

			recruitList = append(recruitList, models.LiveNvnListRecruitList{
				RecrutSn:    sRecrutSn,
				RecrutTitle: sRecrutTitle,
			})
		}

		if err := procRset.Err(); err != nil {
			panic(err)
		}

		fmt.Printf(fmt.Sprintf("sTotCnt: %v", sTotCnt))
	}
	// End : Live Interview List recruit list

	// Start : Live Interview List apply list
	fmt.Printf(fmt.Sprintf("CALL ZSP_LIVE_NVN_LIST_APPLY('%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pPChrgSn, pGbnCd1, pGbnCd2, pGbnCd3, pGbnCd4))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_NVN_LIST_APPLY('%v', '%v', '%v', '%v', '%v', '%v', '%v',:1)",
		pLang, pEntpMemNo, pPChrgSn, pGbnCd1, pGbnCd2, pGbnCd3, pGbnCd4),
		ora.I64, /* TOT_CNT */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* NM */
		ora.S,   /* SEX */
		ora.S,   /* AGE */
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

	applyList := make([]models.LiveNvNSearchApply, 0)

	var (
		aTotCnt  int64
		aPpMemNo string
		aNm      string
		aSex     string
		aAge     string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			aTotCnt = procRset.Row[0].(int64)
			aPpMemNo = procRset.Row[1].(string)
			aNm = procRset.Row[2].(string)
			aSex = procRset.Row[3].(string)
			aAge = procRset.Row[4].(string)

			applyList = append(applyList, models.LiveNvNSearchApply{
				LsaPpMemNo: aPpMemNo,
				LsaNm:      aNm,
				LsaSex:     aSex,
				LsaAge:     aAge,
			})

			//fmt.Printf(fmt.Sprintf("applyList: PpMemNo:%v, Nm:%v, Sex:%v, Age:%v", aPpMemNo, aNm, aSex, aAge))
		}

		if err := procRset.Err(); err != nil {
			panic(err)
		}

		fmt.Printf(fmt.Sprintf("applyList: sTotCnt: %v", aTotCnt))
	}
	// End : Live Interview List apply list

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

	// Start : Live Interview List
	// Get Parameter
	var (
		pageNo    int64
		pageSize  int64
		finalPage int64
		pageList  int64
	)

	pPageNo := c.GetString("pn")
	if pPageNo == "" {
		pPageNo = "1"
	}
	pageNo, err = strconv.ParseInt(pPageNo, 10, 64)
	if err != nil {
		fmt.Printf(fmt.Sprintf("strconv.ParseInt(pPageNo, 10, 64) -> Error"))
	}

	pPageSize := c.GetString("size")
	if pPageSize == "" {
		pPageSize = "100"
	}

	//pPageSize = strconv.FormatInt(imPageNo, 16)

	pageSize, err = strconv.ParseInt(pPageSize, 10, 64)
	if err != nil {
		fmt.Printf(fmt.Sprintf("pageSize, err = strconv.ParseInt(pPageNo, 10, 64) -> Error"))
	}
	pOffSet := (pageNo - 1) * pageSize
	pLimit := pageSize
	pageList = 10

	// fmt.Printf(fmt.Sprintf("CALL ZSP_LIVE_NVN_LIST('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
	// 	pLang, pOffSet, pLimit, pEntpMemNo, pPChrgSn, pStatCd, pKeyword, pViewType, pGbnCd1, pGbnCd2, pGbnCd3, pGbnCd4, pSortCd, pSortWay, pGbnRecrutSn, pGbnSSdy, pGbnSEdy))

	// stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_NVN_LIST('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
	// 	pLang, pOffSet, pLimit, pEntpMemNo, pPChrgSn, pStatCd, pKeyword, pViewType, pGbnCd1, pGbnCd2, pGbnCd3, pGbnCd4, pSortCd, pSortWay, pGbnRecrutSn, pGbnSSdy, pGbnSEdy),
	fmt.Printf(fmt.Sprintf("CALL ZSP_LIVE_NVN_LIST_V2('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pPChrgSn, pStatCd, pKeyword, pViewType, pGbnCd1, pGbnCd2, pGbnCd3, pGbnCd4, pSortCd, pSortWay, pGbnRecrutSn, pGbnSSdy, pGbnSEdy, pGbnESdy, pGbnEEdy, pGbnRSdy, pGbnREdy))

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_NVN_LIST_V2('%v', %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pOffSet, pLimit, pEntpMemNo, pPChrgSn, pStatCd, pKeyword, pViewType, pGbnCd1, pGbnCd2, pGbnCd3, pGbnCd4, pSortCd, pSortWay, pGbnRecrutSn, pGbnSSdy, pGbnSEdy, pGbnESdy, pGbnEEdy, pGbnRSdy, pGbnREdy),
		ora.I64, /* TOT_CNT */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* RECRUT_TITLE */
		ora.S,   /* LIVE_SN */
		ora.S,   /* LIVE_ITV_SDAY */
		ora.S,   /* LIVE_ITV_STIME */
		ora.S,   /* LIVE_ITV_EDAY */
		ora.S,   /* LIVE_ITV_ETIME */
		ora.S,   /* LIVE_REG_DAY */
		ora.S,   /* LIVE_REG_TIME */
		ora.S,   /* LIVE_STAT_CD */
		ora.S,   /* PP_CHRG_NM */
		ora.S,   /* PP_CHRG_BP_NM */
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

	rtnLiveNvnList := models.RtnLiveNvNList{}
	liveNvNList := make([]models.LiveNvNList, 0)

	//searchApplyList := make([]models.LiveNvNApplyList, 0)
	//searchRecruitList := make([]models.LiveNvnListRecruitList, 0)

	var (
		totCnt         int64
		recrutSn       string
		recrutTitle    string
		liveSn         string
		liveItvSday    string
		liveItvStime   string
		liveItvEday    string
		liveItvEtime   string
		liveItvRegDay  string
		liveItvRegTime string
		liveStatCd     string
		ppChrgNm       string
		ppChrgBpNm     string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			totCnt = procRset.Row[0].(int64)
			if totCnt > 0 {
				recrutSn = procRset.Row[1].(string)
				recrutTitle = procRset.Row[2].(string)
				liveSn = procRset.Row[3].(string)
				liveItvSday = procRset.Row[4].(string)
				liveItvStime = procRset.Row[5].(string)
				liveItvEday = procRset.Row[6].(string)
				liveItvEtime = procRset.Row[7].(string)
				liveItvRegDay = procRset.Row[8].(string)
				liveItvRegTime = procRset.Row[9].(string)
				liveStatCd = procRset.Row[10].(string)
				ppChrgNm = procRset.Row[11].(string)
				ppChrgBpNm = procRset.Row[12].(string)

				var (
					prevPageNo  int64 // 이전 페이지 번호
					nextPageNo  int64 // 다음 페이지 번호
					startPageNo int64
					endPageNo   int64
					totalPage   int64
				)

				prevPageNo = 0
				nextPageNo = 0

				finalPage = (totCnt + (pageSize - 1)) / pageSize // 마지막 페이지
				if pageNo > finalPage {                          // 기본값 설정
					pageNo = finalPage
				}

				if pageNo < 0 || pageNo > finalPage { // 현재 페이지 유효성체크
					pageNo = 1
				}

				var (
					isNowFirst bool
					isNowFinal bool
				)

				if pageNo == 1 {
					isNowFirst = true
				} else {
					isNowFirst = false
				}

				if pageNo == finalPage {
					isNowFinal = true
				} else {
					isNowFinal = false
				}

				if isNowFirst { // 이전페이지 계산
					prevPageNo = 1
				} else {
					if (pageNo - 1) < 1 {
						prevPageNo = 1
					} else {
						prevPageNo = pageNo - 1
					}
				}

				if isNowFinal { // 다음페이지 계산
					nextPageNo = finalPage
				} else {
					if (pageNo + 1) > finalPage {
						nextPageNo = finalPage
					} else {
						nextPageNo = pageNo + 1
					}
				}

				d := float64(pageNo) / float64(pageList)
				blockSize := int64(math.Ceil(d))

				startPageNo = ((blockSize - 1) * pageList) + 1
				endPageNo = startPageNo + pageList - 1

				t := float64(totCnt) / float64(pageSize)
				totalPage = int64(math.Ceil(t))

				if endPageNo > totalPage {
					endPageNo = totalPage
				}

				var pagination string

				if totCnt == 0 {
					pagination += "<a href='javascript:void(0);' class='prev disabled'>이전</a>"
					pagination += "<a href='javascript:void(0);' class='disabled'>1</a>"
					pagination += "<a href='javascript:void(0);' class='next disabled'>다음</a>"
				} else {
					if prevPageNo == pageNo {
						pagination += "<a href='javascript:void(0);' class='prev disabled' id='prev' data-page='" + strconv.Itoa(int(prevPageNo)) + "'>이전</a>"
					} else {
						pagination += "<a href='javascript:void(0);' class='prev goPage' id='prev' data-page='" + strconv.Itoa(int(prevPageNo)) + "'>이전</a>"
					}
					for i := startPageNo; i <= endPageNo; i++ {
						if i == pageNo {
							pagination += "<a href='javascript:void(0);' class='active num goPage' id='num' data-page='" + strconv.Itoa(int(i)) + "'>" + strconv.Itoa(int(i)) + "</a>"
						} else {
							pagination += "<a href='javascript:void(0);' class='num goPage' id='num' data-page='" + strconv.Itoa(int(i)) + "'>" + strconv.Itoa(int(i)) + "</a>"
						}
					}
					if nextPageNo == pageNo {
						pagination += "<a href='javascript:void(0);' class='next disabled' id='next' data-page='" + strconv.Itoa(int(nextPageNo)) + "'>다음</a>"
					} else {
						pagination += "<a href='javascript:void(0);' class='next goPage' id='next' data-page='" + strconv.Itoa(int(nextPageNo)) + "'>다음</a>"
					}
				}

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
				)

				fmt.Printf(fmt.Sprintf("CALL ZSP_LIVE_NVN_APPL_MEM_LIST('%v', '%v', '%v', '%v', :1)",
					pLang, pEntpMemNo, recrutSn, liveSn))

				stmtProcCallMem, errMem := ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_NVN_APPL_MEM_LIST('%v', '%v', '%v', '%v', :1)",
					pLang, pEntpMemNo, recrutSn, liveSn),
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

				liveNvNApplyList := make([]models.LiveNvNApplyList, 0)

				if procRsetMem.IsOpen() {
					for procRsetMem.Next() {
						lmPpMemNo = procRsetMem.Row[0].(string)
						lmName = procRsetMem.Row[1].(string)
						lmSex = procRsetMem.Row[2].(string)
						lmAge = procRsetMem.Row[3].(string)
						lmPtoPath = procRsetMem.Row[4].(string)
						lmRecrutSn = procRsetMem.Row[5].(string)
						lmLiveStatCd = procRsetMem.Row[6].(string)
						lmTRC04LiveReqStatCd = procRset.Row[7].(string)
						lmTRC04LiveSn = procRset.Row[8].(string)
						lmMsgYn = procRset.Row[9].(string)
						lmMsgEndYn = procRset.Row[10].(string)
						lmReadEndDay = procRsetMem.Row[11].(string)

						if lmPtoPath == "" {
							lmPtofullPath = lmPtoPath
						} else {
							lmPtofullPath = imgServer + lmPtoPath
						}

						liveNvNApplyList = append(liveNvNApplyList, models.LiveNvNApplyList{
							LmPpMemNo:            lmPpMemNo,
							LmNm:                 lmName,
							LmSex:                lmSex,
							LmAge:                lmAge,
							LmPtoPath:            lmPtofullPath,
							LmRecrutSn:           lmRecrutSn,
							LmRecrutTitle:        recrutTitle,
							LmLiveStatCd:         lmLiveStatCd,
							LmTRC04LiveReqStatCd: lmTRC04LiveReqStatCd,
							LmTRC04LiveSn:        lmTRC04LiveSn,
							LmMsgYn:              lmMsgYn,
							LmMsgEndYn:           lmMsgEndYn,
							LmReadEndDay:         lmReadEndDay,
						})

						// searchApplyList = append(searchApplyList, models.LiveNvNApplyList{
						// 	LmPpMemNo:            lmPpMemNo,
						// 	LmNm:                 lmName,
						// 	LmSex:                lmSex,
						// 	LmAge:                lmAge,
						// 	LmPtoPath:            lmPtofullPath,
						// 	LmRecrutSn:           lmRecrutSn,
						// 	LmRecrutTitle:        recrutTitle,
						// 	LmLiveStatCd:         lmLiveStatCd,
						// 	LmTRC04LiveReqStatCd: lmTRC04LiveReqStatCd,
						// 	LmTRC04LiveSn:        lmTRC04LiveSn,
						// 	LmMsgYn:              lmMsgYn,
						// 	LmMsgEndYn:           lmMsgEndYn,
						// 	LmReadEndDay:         lmReadEndDay,
						// })
					}

					if errMem := procRsetMem.Err(); errMem != nil {
						panic(errMem)
					}
				}
				// End : live apply member list

				// Start : live entp member list
				var (
					lmPpChrgGbnCd  string
					lmPpChrgNm     string
					lmPpChrgBpNm   string
					lmPpChrgSn     string
					lmPpLiveStatCd string
				)

				fmt.Printf(fmt.Sprintf("CALL ZSP_LIVE_NVN_ENTP_MEM_LIST('%v', '%v', '%v', '%v', :1)",
					pLang, pEntpMemNo, recrutSn, liveSn))

				stmtProcCallMem, errMem = ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_NVN_ENTP_MEM_LIST('%v', '%v', '%v', '%v', :1)",
					pLang, pEntpMemNo, recrutSn, liveSn),
					ora.S, /* PP_CHRG_GBN_CD */
					ora.S, /* PP_CHRG_NM */
					ora.S, /* PP_CHRG_BP_NM */
					ora.S, /* PP_CHRG_SN */
					ora.S, /* LIVE_STAT_CD */
				)
				defer stmtProcCallMem.Close()
				if errMem != nil {
					panic(errMem)
				}
				procRsetMem = &ora.Rset{}
				_, errMem = stmtProcCallMem.Exe(procRsetMem)

				if errMem != nil {
					panic(errMem)
				}

				liveNvnMemList := make([]models.LiveNvNMemList, 0)

				if procRsetMem.IsOpen() {
					for procRsetMem.Next() {
						lmPpChrgGbnCd = procRsetMem.Row[0].(string)
						lmPpChrgNm = procRsetMem.Row[1].(string)
						lmPpChrgBpNm = procRsetMem.Row[2].(string)
						lmPpChrgSn = procRsetMem.Row[3].(string)
						lmPpLiveStatCd = procRsetMem.Row[4].(string)

						liveNvnMemList = append(liveNvnMemList, models.LiveNvNMemList{
							LmPpChrgGbnCd:  lmPpChrgGbnCd,
							LmPpChrgNm:     lmPpChrgNm,
							LmPpChrgBpNm:   lmPpChrgBpNm,
							LmPpChrgSn:     lmPpChrgSn,
							LmPpLiveStatCd: lmPpLiveStatCd,
						})
					}
					if errMem := procRsetMem.Err(); errMem != nil {
						panic(errMem)
					}
				}
				// End : live entp member list

				liveNvNList = append(liveNvNList, models.LiveNvNList{
					TotCnt:         totCnt,
					RecrutSn:       recrutSn,
					RecrutTitle:    recrutTitle,
					LiveSn:         liveSn,
					LiveItvSday:    liveItvSday,
					LiveItvStime:   liveItvStime,
					LiveItvEday:    liveItvEday,
					LiveItvEtime:   liveItvEtime,
					LiveItvRegDay:  liveItvRegDay,
					LiveItvRegTime: liveItvRegTime,
					LiveStatCd:     liveStatCd,
					PpChrgNm:       ppChrgNm,
					PpChrgBpNm:     ppChrgBpNm,
					ApplyList:      liveNvNApplyList,
					MemList:        liveNvnMemList,
					Pagination:     pagination,
				})

				// searchRecruitList = append(searchRecruitList, models.LiveNvnListRecruitList{
				// 	RecrutSn:    recrutSn,
				// 	RecrutTitle: recrutTitle,
				// })

				//fmt.Printf(fmt.Sprintf("======= '%v')", liveNvNApplyList.GetApplyInfoStr()))
			}
		}

		if err := procRset.Err(); err != nil {
			panic(err)
		}

		//searchRecruitList = removeDuplicateValues(searchRecruitList)
		//RtnLiveNvnListRecruitListData: searchRecruitList,
		// searchApplyList = removeDuplicateValues(searchApplyList)
		// for _, val := range searchApplyList {
		// 	fmt.Printf(fmt.Sprintf("======= '%v')", val.GetApplyInfoStr()))
		// }

		// var (
		// 	waitCnt int64
		// 	ingCnt  int64
		// 	endCnt  int64
		// 	cnclCnt int64
		// )

		// for _, val := range liveNvNList {

		// 	if val.LiveStatCd == "01" {
		// 		waitCnt++
		// 	} else if val.LiveStatCd == "02" {
		// 		ingCnt++
		// 	} else if val.LiveStatCd == "03" {
		// 		endCnt++
		// 	} else {
		// 		cnclCnt++
		// 	}
		// }

		// fmt.Printf(fmt.Sprintf("waitCnt:%v, ingCnt:%v., endCnt:%v, cnclCnt:%v", waitCnt, ingCnt, endCnt, cnclCnt))

		rtnLiveNvnList = models.RtnLiveNvNList{
			RtnLiveNvNListData: liveNvNList,

			RtnLiveNvnListRecruitData: recruitList,
			RtnGbnRecrutSn:            pGbnRecrutSn,

			RtnGbnPpChrgSn: pGbnPpChrgSn,

			RtnLiveNvnListApplyData: applyList,

			RtnEntpTeamMemberListData: entpTeamMemberList,
		}

		c.Data["json"] = &rtnLiveNvnList
		c.ServeJSON()
	}
	// End : Live Interview List
}

/*
func removeDuplicateValues(checkArray []models.LiveNvNApplyList) []models.LiveNvNApplyList {
	keys := make(map[models.LiveNvNApplyList]bool)
	list := []models.LiveNvNApplyList{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range checkArray {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}
*/

/*
func removeDuplicateValues(checkArray []models.LiveNvnListRecruitList) []models.LiveNvnListRecruitList {
	keys := make(map[models.LiveNvnListRecruitList]bool)
	list := []models.LiveNvnListRecruitList{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range checkArray {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}
*/
