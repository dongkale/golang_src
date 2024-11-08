package controllers

import (
	"fmt"

	"emst.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	ora "gopkg.in/rana/ora.v4"
)

type AdminStatsRecruitMainController struct {
	BaseController
}

func (c *AdminStatsRecruitMainController) Get() {

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
	pEntpMemNo := c.GetString("entp_mem_no")

	if pEntpMemNo == "" {
		pEntpMemNo = "T10"
	}

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Admin Member Stats List
	log.Debug("CALL SP_EMS_ADMIN_STATS_RC01_R('%v', :1)",
		pLang)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_STATS_RC01_R('%v', :1)",
		pLang),
		ora.I64, /* RC01_TOT_CNT */
		ora.I64, /* RC01_NEW_CNT */
		ora.I64, /* RC01_ING_CNT */
		ora.I64, /* RC01_END_CNT */
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

	adminStatsRC01 := make([]models.AdminStatsRC01, 0)

	var (
		rc01TotCnt int64
		rc01NewCnt int64
		rc01IngCnt int64
		rc01EndCnt int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			rc01TotCnt = procRset.Row[0].(int64)
			rc01NewCnt = procRset.Row[1].(int64)
			rc01IngCnt = procRset.Row[2].(int64)
			rc01EndCnt = procRset.Row[3].(int64)

			adminStatsRC01 = append(adminStatsRC01, models.AdminStatsRC01{
				Rc01TotCnt: rc01TotCnt,
				Rc01NewCnt: rc01NewCnt,
				Rc01IngCnt: rc01IngCnt,
				Rc01EndCnt: rc01EndCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Admin Member Stats List

	// Start : Admin Member Stats List
	log.Debug("CALL SP_EMS_ADMIN_STATS_RC02_R('%v', :1)",
		pLang)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_STATS_RC02_R('%v', :1)",
		pLang),
		ora.I64, /* RC02_TOT_CNT */
		ora.I64, /* RC02_ING_CNT */
		ora.I64, /* RC02_PASS_CNT */
		ora.I64, /* RC02_FAIL_CNT */
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

	adminStatsRC02 := make([]models.AdminStatsRC02, 0)

	var (
		rc02TotCnt  int64
		rc02IngCnt  int64
		rc02PassCnt int64
		rc02FailCnt int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			rc02TotCnt = procRset.Row[0].(int64)
			rc02IngCnt = procRset.Row[1].(int64)
			rc02PassCnt = procRset.Row[2].(int64)
			rc02FailCnt = procRset.Row[3].(int64)

			adminStatsRC02 = append(adminStatsRC02, models.AdminStatsRC02{
				Rc02TotCnt:  rc02TotCnt,
				Rc02IngCnt:  rc02IngCnt,
				Rc02PassCnt: rc02PassCnt,
				Rc02FailCnt: rc02FailCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Admin Member Stats List

	// Start : Admin Member Stats List
	log.Debug("CALL SP_EMS_ADMIN_STATS_RC03_R('%v', '%v', :1)",
		pLang, pEntpMemNo)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_STATS_RC03_R('%v', '%v', :1)",
		pLang, pEntpMemNo),
		ora.S,   /* RC03_ENTP_MEM_NO */
		ora.S,   /* RC03_ENTP_KO_NM */
		ora.I64, /* RC03_APPLY_CNT */
		ora.I64, /* RC03_PASS_CNT */
		ora.I64, /* RC03_FAIL_CNT */
		ora.I64, /* RC03_ING_CNT */
		ora.F64, /* RC03_MATCHING_RATE */
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

	adminStatsRC03 := make([]models.AdminStatsRC03, 0)

	var (
		rc03EntpMemNo    string
		rc03EntpKoNm     string
		rc03ApplyCnt     int64
		rc03PassCnt      int64
		rc03FailCnt      int64
		rc03IngCnt       int64
		rc03MatchingRate float64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			rc03EntpMemNo = procRset.Row[0].(string)
			rc03EntpKoNm = procRset.Row[1].(string)
			rc03ApplyCnt = procRset.Row[2].(int64)
			rc03PassCnt = procRset.Row[3].(int64)
			rc03FailCnt = procRset.Row[4].(int64)
			rc03IngCnt = procRset.Row[5].(int64)
			rc03MatchingRate = procRset.Row[6].(float64)

			adminStatsRC03 = append(adminStatsRC03, models.AdminStatsRC03{
				Rc03EntpMemNo:    rc03EntpMemNo,
				Rc03EntpKoNm:     rc03EntpKoNm,
				Rc03ApplyCnt:     rc03ApplyCnt,
				Rc03PassCnt:      rc03PassCnt,
				Rc03FailCnt:      rc03FailCnt,
				Rc03IngCnt:       rc03IngCnt,
				Rc03MatchingRate: rc03MatchingRate,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Admin Member Stats List

	// Start : Admin Member Stats List
	log.Debug("CALL SP_EMS_ADMIN_STATS_RC_SUB_R('%v', :1)",
		pLang)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_STATS_RC_SUB_R('%v', :1)",
		pLang),
		ora.S, /* SUB_ENTP_MEM_NO */
		ora.S, /* SUB_ENTP_KO_NM */
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

	adminStatsRCSub := make([]models.AdminStatsRCSub, 0)

	var (
		subEntpMemNo string
		subEntpKoNm  string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			subEntpMemNo = procRset.Row[0].(string)
			subEntpKoNm = procRset.Row[1].(string)

			adminStatsRCSub = append(adminStatsRCSub, models.AdminStatsRCSub{
				SubEntpMemNo: subEntpMemNo,
				SubEntpKoNm:  subEntpKoNm,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Admin Member Stats List

	c.Data["Rc01TotCnt"] = rc01TotCnt
	c.Data["Rc01NewCnt"] = rc01NewCnt
	c.Data["Rc01IngCnt"] = rc01IngCnt
	c.Data["Rc01EndCnt"] = rc01EndCnt

	c.Data["Rc02TotCnt"] = rc02TotCnt
	c.Data["Rc02IngCnt"] = rc02IngCnt
	c.Data["Rc02PassCnt"] = rc02PassCnt
	c.Data["Rc02FailCnt"] = rc02FailCnt

	c.Data["AdminStatsRC03"] = adminStatsRC03
	c.Data["AdminStatsRCSub"] = adminStatsRCSub

	c.Data["MenuId"] = "99"
	c.TplName = "admin/stats_recruit_main.html"
}

func (c *AdminStatsPeriodMainController) Post() {

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
	pEntpMemNo := c.GetString("entp_mem_no")
	if pEntpMemNo == "" {
		pEntpMemNo = "T10"
	}

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Admin Member Stats List
	log.Debug("CALL SP_EMS_ADMIN_STATS_RC03_R('%v', '%v', :1)",
		pLang, pEntpMemNo)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_STATS_RC03_R('%v', '%v', :1)",
		pLang, pEntpMemNo),
		ora.S,   /* RC03_ENTP_MEM_NO */
		ora.S,   /* RC03_ENTP_KO_NM */
		ora.I64, /* RC03_APPLY_CNT */
		ora.I64, /* RC03_PASS_CNT */
		ora.I64, /* RC03_FAIL_CNT */
		ora.I64, /* RC03_ING_CNT */
		ora.F64, /* RC03_MATCHING_RATE */
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

	rtnAdminStatsRC03 := models.RtnAdminStatsRC03{}
	adminStatsRC03 := make([]models.AdminStatsRC03, 0)

	var (
		rc03EntpMemNo    string
		rc03EntpKoNm     string
		rc03ApplyCnt     int64
		rc03PassCnt      int64
		rc03FailCnt      int64
		rc03IngCnt       int64
		rc03MatchingRate float64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			rc03EntpMemNo = procRset.Row[0].(string)
			rc03EntpKoNm = procRset.Row[1].(string)
			rc03ApplyCnt = procRset.Row[2].(int64)
			rc03PassCnt = procRset.Row[3].(int64)
			rc03FailCnt = procRset.Row[4].(int64)
			rc03IngCnt = procRset.Row[5].(int64)
			rc03MatchingRate = procRset.Row[6].(float64)

			adminStatsRC03 = append(adminStatsRC03, models.AdminStatsRC03{
				Rc03EntpMemNo:    rc03EntpMemNo,
				Rc03EntpKoNm:     rc03EntpKoNm,
				Rc03ApplyCnt:     rc03ApplyCnt,
				Rc03PassCnt:      rc03PassCnt,
				Rc03FailCnt:      rc03FailCnt,
				Rc03IngCnt:       rc03IngCnt,
				Rc03MatchingRate: rc03MatchingRate,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnAdminStatsRC03 = models.RtnAdminStatsRC03{
			RtnAdminStatsRC03Data: adminStatsRC03,
		}
		// End : Admin Stats Sub2

		c.Data["json"] = &rtnAdminStatsRC03
		c.ServeJSON()
	}
}
