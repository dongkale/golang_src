package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type LiveListController struct {
	BaseController
}

func (c *LiveListController) Get() {

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
	//pStatCd := c.GetString("stat_cd")     //상태코드
	//pKeyword := c.GetString("keyword")    //검색어
	pViewType := c.GetString("view_type") //보기구분
	if pViewType == "" {
		pViewType = "G"
	}
	//pGbnCd1 := c.GetString("gbn_cd1") //구분값(요청)
	//pGbnCd2 := c.GetString("gbn_cd2") //구분값(예정)
	//pGbnCd3 := c.GetString("gbn_cd3") //구분값(종료)
	//pGbnCd4 := c.GetString("gbn_cd4") //구분값(거절.취소)
	//imgServer, _  := beego.AppConfig.String("viewpath")

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

	fmt.Printf(fmt.Sprintf("CALL ZSP_LIVE_STATS_R('%v', '%v', :1)",
		pLang, pEntpMemNo))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_STATS_R('%v','%v', :1)",
		pLang, pEntpMemNo),
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

	// // Start : Live Interview List 01
	// pStatCd = "01"
	// log.Debug("CALL ZSP_LIVE_LIST_R('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
	// 	pLang, pEntpMemNo, pStatCd, pKeyword, pViewType, pGbnCd1, pGbnCd2, pGbnCd3, pGbnCd4)

	// stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_LIST_R('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
	// 	pLang, pEntpMemNo, pStatCd, pKeyword, pViewType, pGbnCd1, pGbnCd2, pGbnCd3, pGbnCd4),
	// 	ora.S,   /* RECRUT_SN */
	// 	ora.S,   /* PP_MEM_NO */
	// 	ora.S,   /* LIVE_SN */
	// 	ora.S,   /* PTO_PATH */
	// 	ora.S,   /* NM */
	// 	ora.S,   /* SEX */
	// 	ora.I64, /* AGE */
	// 	ora.S,   /* LIVE_ITV_SDAY */
	// 	ora.S,   /* LIVE_ITV_STIME */
	// 	ora.S,   /* LIVE_ITV_EDAY */
	// 	ora.S,   /* LIVE_ITV_ETIME */
	// )
	// defer stmtProcCall.Close()
	// if err != nil {
	// 	panic(err)
	// }
	// procRset = &ora.Rset{}
	// _, err = stmtProcCall.Exe(procRset)

	// if err != nil {
	// 	panic(err)
	// }

	// liveList01 := make([]models.LiveList01, 0)

	// var (
	// 	s01RecrutSn     string
	// 	s01PpMemNo      string
	// 	s01LiveSn       string
	// 	s01PtoPath      string
	// 	s01Nm           string
	// 	s01Sex          string
	// 	s01Age          int64
	// 	s01LiveItvSday  string
	// 	s01LiveItvStime string
	// 	s01LiveItvEday  string
	// 	s01LiveItvEtime string
	// 	s01FullPtoPath  string
	// )

	// if procRset.IsOpen() {
	// 	for procRset.Next() {
	// 		s01RecrutSn = procRset.Row[0].(string)
	// 		s01PpMemNo = procRset.Row[1].(string)
	// 		s01LiveSn = procRset.Row[2].(string)
	// 		s01PtoPath = procRset.Row[3].(string)
	// 		if s01PtoPath == "" {
	// 			s01FullPtoPath = s01PtoPath
	// 		} else {
	// 			s01FullPtoPath = imgServer + s01PtoPath
	// 		}
	// 		s01Nm = procRset.Row[4].(string)
	// 		s01Sex = procRset.Row[5].(string)
	// 		s01Age = procRset.Row[6].(int64)
	// 		s01LiveItvSday = procRset.Row[7].(string)
	// 		s01LiveItvStime = procRset.Row[8].(string)
	// 		s01LiveItvEday = procRset.Row[9].(string)
	// 		s01LiveItvEtime = procRset.Row[10].(string)

	// 		var (
	// 			lmPpChrgGbnCd string
	// 			lmPpChrgNm    string
	// 			lmPpChrgBpNm  string
	// 		)

	// 		// Start : live member list
	// 		log.Debug("CALL ZSP_LIVE_MEM_LIST_R('%v', '%v', '%v', '%v', '%v', :1)",
	// 			pLang, pEntpMemNo, s01RecrutSn, s01PpMemNo, s01LiveSn)

	// 		stmtProcCallMem, errMem := ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_MEM_LIST_R('%v', '%v', '%v', '%v', '%v', :1)",
	// 			pLang, pEntpMemNo, s01RecrutSn, s01PpMemNo, s01LiveSn),
	// 			ora.S, /* PP_CHRG_GBN_CD */
	// 			ora.S, /* PP_CHRG_NM */
	// 			ora.S, /* PP_CHRG_BP_NM */
	// 		)
	// 		defer stmtProcCallMem.Close()
	// 		if errMem != nil {
	// 			panic(errMem)
	// 		}
	// 		procRsetMem := &ora.Rset{}
	// 		_, errMem = stmtProcCallMem.Exe(procRsetMem)

	// 		if errMem != nil {
	// 			panic(errMem)
	// 		}

	// 		liveMemList := make([]models.LiveMemList, 0)

	// 		if procRsetMem.IsOpen() {
	// 			for procRsetMem.Next() {
	// 				lmPpChrgGbnCd = procRsetMem.Row[0].(string)
	// 				lmPpChrgNm = procRsetMem.Row[1].(string)
	// 				lmPpChrgBpNm = procRsetMem.Row[2].(string)

	// 				liveMemList = append(liveMemList, models.LiveMemList{
	// 					LmPpChrgGbnCd: lmPpChrgGbnCd,
	// 					LmPpChrgNm:    lmPpChrgNm,
	// 					LmPpChrgBpNm:  lmPpChrgBpNm,
	// 				})
	// 			}
	// 			if errMem := procRsetMem.Err(); errMem != nil {
	// 				panic(errMem)
	// 			}
	// 		}
	// 		// End : live member list

	// 		liveList01 = append(liveList01, models.LiveList01{
	// 			S01RecrutSn:     s01RecrutSn,
	// 			S01PpMemNo:      s01PpMemNo,
	// 			S01LiveSn:       s01LiveSn,
	// 			S01PtoPath:      s01FullPtoPath,
	// 			S01Nm:           s01Nm,
	// 			S01Sex:          s01Sex,
	// 			S01Age:          s01Age,
	// 			S01LiveItvSday:  s01LiveItvSday,
	// 			S01LiveItvStime: s01LiveItvStime,
	// 			S01LiveItvEday:  s01LiveItvEday,
	// 			S01LiveItvEtime: s01LiveItvEtime,
	// 			S01SubList:      liveMemList,
	// 		})
	// 	}
	// 	if err := procRset.Err(); err != nil {
	// 		panic(err)
	// 	}
	// }
	// // End : Live Interview List 01

	// // Start : Live Interview List 02
	// pStatCd = "02"
	// log.Debug("CALL ZSP_LIVE_LIST_R('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
	// 	pLang, pEntpMemNo, pStatCd, pKeyword, pViewType, pGbnCd1, pGbnCd2, pGbnCd3, pGbnCd4)

	// stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_LIST_R('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
	// 	pLang, pEntpMemNo, pStatCd, pKeyword, pViewType, pGbnCd1, pGbnCd2, pGbnCd3, pGbnCd4),
	// 	ora.S,   /* RECRUT_SN */
	// 	ora.S,   /* PP_MEM_NO */
	// 	ora.S,   /* LIVE_SN */
	// 	ora.S,   /* PTO_PATH */
	// 	ora.S,   /* NM */
	// 	ora.S,   /* SEX */
	// 	ora.I64, /* AGE */
	// 	ora.S,   /* LIVE_ITV_SDAY */
	// 	ora.S,   /* LIVE_ITV_STIME */
	// 	ora.S,   /* LIVE_ITV_EDAY */
	// 	ora.S,   /* LIVE_ITV_ETIME */
	// )
	// defer stmtProcCall.Close()
	// if err != nil {
	// 	panic(err)
	// }
	// procRset = &ora.Rset{}
	// _, err = stmtProcCall.Exe(procRset)

	// if err != nil {
	// 	panic(err)
	// }

	// liveList02 := make([]models.LiveList02, 0)

	// var (
	// 	s02RecrutSn     string
	// 	s02PpMemNo      string
	// 	s02LiveSn       string
	// 	s02PtoPath      string
	// 	s02Nm           string
	// 	s02Sex          string
	// 	s02Age          int64
	// 	s02LiveItvSday  string
	// 	s02LiveItvStime string
	// 	s02LiveItvEday  string
	// 	s02LiveItvEtime string
	// 	s02FullPtoPath  string
	// )

	// if procRset.IsOpen() {
	// 	for procRset.Next() {
	// 		s02RecrutSn = procRset.Row[0].(string)
	// 		s02PpMemNo = procRset.Row[1].(string)
	// 		s02LiveSn = procRset.Row[2].(string)
	// 		s02PtoPath = procRset.Row[3].(string)
	// 		if s02PtoPath == "" {
	// 			s02FullPtoPath = s02PtoPath
	// 		} else {
	// 			s02FullPtoPath = imgServer + s02PtoPath
	// 		}
	// 		s02Nm = procRset.Row[4].(string)
	// 		s02Sex = procRset.Row[5].(string)
	// 		s02Age = procRset.Row[6].(int64)
	// 		s02LiveItvSday = procRset.Row[7].(string)
	// 		s02LiveItvStime = procRset.Row[8].(string)
	// 		s02LiveItvEday = procRset.Row[9].(string)
	// 		s02LiveItvEtime = procRset.Row[10].(string)

	// 		var (
	// 			lmPpChrgGbnCd string
	// 			lmPpChrgNm    string
	// 			lmPpChrgBpNm  string
	// 		)

	// 		// Start : live member list
	// 		log.Debug("CALL ZSP_LIVE_MEM_LIST_R('%v', '%v', '%v', '%v', '%v', :1)",
	// 			pLang, pEntpMemNo, s02RecrutSn, s02PpMemNo, s02LiveSn)

	// 		stmtProcCallMem, errMem := ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_MEM_LIST_R('%v', '%v', '%v', '%v', '%v', :1)",
	// 			pLang, pEntpMemNo, s02RecrutSn, s02PpMemNo, s02LiveSn),
	// 			ora.S, /* PP_CHRG_GBN_CD */
	// 			ora.S, /* PP_CHRG_NM */
	// 			ora.S, /* PP_CHRG_BP_NM */
	// 		)
	// 		defer stmtProcCallMem.Close()
	// 		if errMem != nil {
	// 			panic(errMem)
	// 		}
	// 		procRsetMem := &ora.Rset{}
	// 		_, errMem = stmtProcCallMem.Exe(procRsetMem)

	// 		if errMem != nil {
	// 			panic(errMem)
	// 		}

	// 		liveMemList := make([]models.LiveMemList, 0)

	// 		if procRsetMem.IsOpen() {
	// 			for procRsetMem.Next() {
	// 				lmPpChrgGbnCd = procRsetMem.Row[0].(string)
	// 				lmPpChrgNm = procRsetMem.Row[1].(string)
	// 				lmPpChrgBpNm = procRsetMem.Row[2].(string)

	// 				liveMemList = append(liveMemList, models.LiveMemList{
	// 					LmPpChrgGbnCd: lmPpChrgGbnCd,
	// 					LmPpChrgNm:    lmPpChrgNm,
	// 					LmPpChrgBpNm:  lmPpChrgBpNm,
	// 				})
	// 			}
	// 			if errMem := procRsetMem.Err(); errMem != nil {
	// 				panic(errMem)
	// 			}
	// 		}
	// 		// End : live member list

	// 		liveList02 = append(liveList02, models.LiveList02{
	// 			S02RecrutSn:     s02RecrutSn,
	// 			S02PpMemNo:      s02PpMemNo,
	// 			S02LiveSn:       s02LiveSn,
	// 			S02PtoPath:      s02FullPtoPath,
	// 			S02Nm:           s02Nm,
	// 			S02Sex:          s02Sex,
	// 			S02Age:          s02Age,
	// 			S02LiveItvSday:  s02LiveItvSday,
	// 			S02LiveItvStime: s02LiveItvStime,
	// 			S02LiveItvEday:  s02LiveItvEday,
	// 			S02LiveItvEtime: s02LiveItvEtime,
	// 			S02SubList:      liveMemList,
	// 		})
	// 	}
	// 	if err := procRset.Err(); err != nil {
	// 		panic(err)
	// 	}
	// }
	// // End : Live Interview List 02

	// // Start : Live Interview List 03
	// pStatCd = "03"
	// log.Debug("CALL ZSP_LIVE_LIST_R('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
	// 	pLang, pEntpMemNo, pStatCd, pKeyword, pViewType, pGbnCd1, pGbnCd2, pGbnCd3, pGbnCd4)

	// stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_LIST_R('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
	// 	pLang, pEntpMemNo, pStatCd, pKeyword, pViewType, pGbnCd1, pGbnCd2, pGbnCd3, pGbnCd4),
	// 	ora.S,   /* RECRUT_SN */
	// 	ora.S,   /* PP_MEM_NO */
	// 	ora.S,   /* LIVE_SN */
	// 	ora.S,   /* PTO_PATH */
	// 	ora.S,   /* NM */
	// 	ora.S,   /* SEX */
	// 	ora.I64, /* AGE */
	// 	ora.S,   /* LIVE_ITV_SDAY */
	// 	ora.S,   /* LIVE_ITV_STIME */
	// 	ora.S,   /* LIVE_ITV_EDAY */
	// 	ora.S,   /* LIVE_ITV_ETIME */
	// )
	// defer stmtProcCall.Close()
	// if err != nil {
	// 	panic(err)
	// }
	// procRset = &ora.Rset{}
	// _, err = stmtProcCall.Exe(procRset)

	// if err != nil {
	// 	panic(err)
	// }

	// liveList03 := make([]models.LiveList03, 0)

	// var (
	// 	s03RecrutSn     string
	// 	s03PpMemNo      string
	// 	s03LiveSn       string
	// 	s03PtoPath      string
	// 	s03Nm           string
	// 	s03Sex          string
	// 	s03Age          int64
	// 	s03LiveItvSday  string
	// 	s03LiveItvStime string
	// 	s03LiveItvEday  string
	// 	s03LiveItvEtime string
	// 	s03FullPtoPath  string
	// )

	// if procRset.IsOpen() {
	// 	for procRset.Next() {
	// 		s03RecrutSn = procRset.Row[0].(string)
	// 		s03PpMemNo = procRset.Row[1].(string)
	// 		s03LiveSn = procRset.Row[2].(string)
	// 		s03PtoPath = procRset.Row[3].(string)
	// 		if s03PtoPath == "" {
	// 			s03FullPtoPath = s03PtoPath
	// 		} else {
	// 			s03FullPtoPath = imgServer + s03PtoPath
	// 		}
	// 		s03Nm = procRset.Row[4].(string)
	// 		s03Sex = procRset.Row[5].(string)
	// 		s03Age = procRset.Row[6].(int64)
	// 		s03LiveItvSday = procRset.Row[7].(string)
	// 		s03LiveItvStime = procRset.Row[8].(string)
	// 		s03LiveItvEday = procRset.Row[9].(string)
	// 		s03LiveItvEtime = procRset.Row[10].(string)

	// 		var (
	// 			lmPpChrgGbnCd string
	// 			lmPpChrgNm    string
	// 			lmPpChrgBpNm  string
	// 		)

	// 		// Start : live member list
	// 		log.Debug("CALL ZSP_LIVE_MEM_LIST_R('%v', '%v', '%v', '%v', '%v', :1)",
	// 			pLang, pEntpMemNo, s03RecrutSn, s03PpMemNo, s03LiveSn)

	// 		stmtProcCallMem, errMem := ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_MEM_LIST_R('%v', '%v', '%v', '%v', '%v', :1)",
	// 			pLang, pEntpMemNo, s03RecrutSn, s03PpMemNo, s03LiveSn),
	// 			ora.S, /* PP_CHRG_GBN_CD */
	// 			ora.S, /* PP_CHRG_NM */
	// 			ora.S, /* PP_CHRG_BP_NM */
	// 		)
	// 		defer stmtProcCallMem.Close()
	// 		if errMem != nil {
	// 			panic(errMem)
	// 		}
	// 		procRsetMem := &ora.Rset{}
	// 		_, errMem = stmtProcCallMem.Exe(procRsetMem)

	// 		if errMem != nil {
	// 			panic(errMem)
	// 		}

	// 		liveMemList := make([]models.LiveMemList, 0)

	// 		if procRsetMem.IsOpen() {
	// 			for procRsetMem.Next() {
	// 				lmPpChrgGbnCd = procRsetMem.Row[0].(string)
	// 				lmPpChrgNm = procRsetMem.Row[1].(string)
	// 				lmPpChrgBpNm = procRsetMem.Row[2].(string)

	// 				liveMemList = append(liveMemList, models.LiveMemList{
	// 					LmPpChrgGbnCd: lmPpChrgGbnCd,
	// 					LmPpChrgNm:    lmPpChrgNm,
	// 					LmPpChrgBpNm:  lmPpChrgBpNm,
	// 				})
	// 			}
	// 			if errMem := procRsetMem.Err(); errMem != nil {
	// 				panic(errMem)
	// 			}
	// 		}
	// 		// End : live member list

	// 		liveList03 = append(liveList03, models.LiveList03{
	// 			S03RecrutSn:     s03RecrutSn,
	// 			S03PpMemNo:      s03PpMemNo,
	// 			S03LiveSn:       s03LiveSn,
	// 			S03PtoPath:      s03FullPtoPath,
	// 			S03Nm:           s03Nm,
	// 			S03Sex:          s03Sex,
	// 			S03Age:          s03Age,
	// 			S03LiveItvSday:  s03LiveItvSday,
	// 			S03LiveItvStime: s03LiveItvStime,
	// 			S03LiveItvEday:  s03LiveItvEday,
	// 			S03LiveItvEtime: s03LiveItvEtime,
	// 			S03SubList:      liveMemList,
	// 		})
	// 	}
	// 	if err := procRset.Err(); err != nil {
	// 		panic(err)
	// 	}
	// }
	// // End : Live Interview List 03

	c.Data["WaitCnt"] = waitCnt
	c.Data["IngCnt"] = ingCnt
	c.Data["EndCnt"] = endCnt
	c.Data["CnclCnt"] = cnclCnt

	//	c.Data["LiveList01"] = liveList01
	//	c.Data["LiveList02"] = liveList02
	//	c.Data["LiveList03"] = liveList03

	c.Data["TMenuId"] = "L00"
	c.Data["SMenuId"] = "L02"

	c.TplName = "live/live_list.html"
}

func (c *LiveListController) Post() {
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
	pStatCd := c.GetString("stat_cd")     //상태코드
	pKeyword := c.GetString("keyword")    //검색어
	pViewType := c.GetString("view_type") //보기구분
	pGbnCd1 := c.GetString("gbn_cd1")     //구분값(요청)
	pGbnCd2 := c.GetString("gbn_cd2")     //구분값(예정)
	pGbnCd3 := c.GetString("gbn_cd3")     //구분값(종료)
	pGbnCd4 := c.GetString("gbn_cd4")     //구분값(거절.취소)
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

	// Start : Live Interview List
	fmt.Printf(fmt.Sprintf("CALL ZSP_LIVE_LIST_R('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pStatCd, pKeyword, pViewType, pGbnCd1, pGbnCd2, pGbnCd3, pGbnCd4))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_LIST_R('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pStatCd, pKeyword, pViewType, pGbnCd1, pGbnCd2, pGbnCd3, pGbnCd4),
		ora.S,   /* RECRUT_SN */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* LIVE_SN */
		ora.S,   /* PTO_PATH */
		ora.S,   /* NM */
		ora.S,   /* SEX */
		ora.I64, /* AGE */
		ora.S,   /* LIVE_ITV_SDAY */
		ora.S,   /* LIVE_ITV_STIME */
		ora.S,   /* LIVE_ITV_EDAY */
		ora.S,   /* LIVE_ITV_ETIME */
		ora.I64, /* TOT_CNT */
		ora.S,   /* LIVE_STAT_CD */
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

	rtnLiveList01 := models.RtnLiveList01{}
	liveList01 := make([]models.LiveList01, 0)

	var (
		s01RecrutSn     string
		s01PpMemNo      string
		s01LiveSn       string
		s01PtoPath      string
		s01Nm           string
		s01Sex          string
		s01Age          int64
		s01LiveItvSday  string
		s01LiveItvStime string
		s01LiveItvEday  string
		s01LiveItvEtime string
		s01TotCnt       int64
		s01LiveStatCd   string
		s01FullPtoPath  string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			s01RecrutSn = procRset.Row[0].(string)
			s01PpMemNo = procRset.Row[1].(string)
			s01LiveSn = procRset.Row[2].(string)
			s01PtoPath = procRset.Row[3].(string)
			if s01PtoPath == "" {
				s01FullPtoPath = s01PtoPath
			} else {
				s01FullPtoPath = imgServer + s01PtoPath
			}
			s01Nm = procRset.Row[4].(string)
			s01Sex = procRset.Row[5].(string)
			s01Age = procRset.Row[6].(int64)
			s01LiveItvSday = procRset.Row[7].(string)
			s01LiveItvStime = procRset.Row[8].(string)
			s01LiveItvEday = procRset.Row[9].(string)
			s01LiveItvEtime = procRset.Row[10].(string)
			s01TotCnt = procRset.Row[11].(int64)
			s01LiveStatCd = procRset.Row[12].(string)

			var (
				lmPpChrgGbnCd string
				lmPpChrgNm    string
				lmPpChrgBpNm  string
			)

			// Start : live member list
			fmt.Printf(fmt.Sprintf("CALL ZSP_LIVE_MEM_LIST_R('%v', '%v', '%v', '%v', '%v', :1)",
				pLang, pEntpMemNo, s01RecrutSn, s01PpMemNo, s01LiveSn))

			stmtProcCallMem, errMem := ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_MEM_LIST_R('%v', '%v', '%v', '%v', '%v', :1)",
				pLang, pEntpMemNo, s01RecrutSn, s01PpMemNo, s01LiveSn),
				ora.S, /* PP_CHRG_GBN_CD */
				ora.S, /* PP_CHRG_NM */
				ora.S, /* PP_CHRG_BP_NM */
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

					liveMemList = append(liveMemList, models.LiveMemList{
						LmPpChrgGbnCd: lmPpChrgGbnCd,
						LmPpChrgNm:    lmPpChrgNm,
						LmPpChrgBpNm:  lmPpChrgBpNm,
					})
				}
				if errMem := procRsetMem.Err(); errMem != nil {
					panic(errMem)
				}
			}
			// End : live member list

			liveList01 = append(liveList01, models.LiveList01{
				S01RecrutSn:     s01RecrutSn,
				S01PpMemNo:      s01PpMemNo,
				S01LiveSn:       s01LiveSn,
				S01PtoPath:      s01FullPtoPath,
				S01Nm:           s01Nm, // 라이브 구직자 이름
				S01Sex:          s01Sex,
				S01Age:          s01Age,
				S01LiveItvSday:  s01LiveItvSday,
				S01LiveItvStime: s01LiveItvStime,
				S01LiveItvEday:  s01LiveItvEday,
				S01LiveItvEtime: s01LiveItvEtime,
				S01TotCnt:       s01TotCnt,
				S01LiveStatCd:   s01LiveStatCd,
				S01SubList:      liveMemList, // 라이브 기업측 맴버
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnLiveList01 = models.RtnLiveList01{
			RtnLiveList01Data: liveList01,
		}
		c.Data["json"] = &rtnLiveList01
		c.ServeJSON()
	}
	// End : Live Interview List

}
