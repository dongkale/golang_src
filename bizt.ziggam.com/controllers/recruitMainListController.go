package controllers

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"bizt.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

type RecruitMainListController struct {
	beego.Controller
}

func (c *RecruitMainListController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	mem_id := session.Get(c.Ctx.Request.Context(), "mem_id")
	if mem_id == nil {
		c.Ctx.Redirect(302, "/login")
	}
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := mem_no // 기업회원번호(세션)

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Recruit Main List
	fmt.Printf(fmt.Sprintf("CALL ZSP_RECRUIT_MAIN_LIST_R('%v', '%v', :1)",
		pLang, pEntpMemNo))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_RECRUIT_MAIN_LIST_R('%v', '%v', :1)",
		pLang, pEntpMemNo),
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* PRGS_STAT */
		ora.S,   /* RECRUT_TITLE */
		ora.S,   /* EMPL_TYP */
		ora.S,   /* UP_JOB_GRP */
		ora.S,   /* JOB_GRP */
		ora.S,   /* RECRUT_DY */
		ora.S,   /* RECRUT_EDT */
		ora.I64, /* APPLY_CNT */
		ora.I64, /* ING_CNT */
		ora.I64, /* PASS_CNT */
		ora.I64, /* FAIL_CNT */
		ora.I64, /* TOT_CNT */
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

	rtnRecruitPostList := models.RtnRecruitPostList{}
	recruitPostList := make([]models.RecruitPostList, 0)

	var (
		entpMemNo   string
		recrutSn    string
		prgsStat    string
		recrutTitle string
		emplTyp     string
		upJobGrp    string
		jobGrp      string
		recrutDy    string
		recrutEdt   string
		applyCnt    int64
		ingCnt      int64
		passCnt     int64
		failCnt     int64
		totCnt      int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			entpMemNo = procRset.Row[0].(string)
			recrutSn = procRset.Row[1].(string)
			prgsStat = procRset.Row[2].(string)
			recrutTitle = procRset.Row[3].(string)
			emplTyp = procRset.Row[4].(string)
			upJobGrp = procRset.Row[5].(string)
			jobGrp = procRset.Row[6].(string)
			recrutDy = procRset.Row[7].(string)
			recrutEdt = procRset.Row[8].(string)
			applyCnt = procRset.Row[9].(int64)
			ingCnt = procRset.Row[10].(int64)
			passCnt = procRset.Row[11].(int64)
			failCnt = procRset.Row[12].(int64)
			totCnt = procRset.Row[13].(int64)

			recruitPostList = append(recruitPostList, models.RecruitPostList{
				EntpMemNo:   entpMemNo,
				RecrutSn:    recrutSn,
				PrgsStat:    prgsStat,
				RecrutTitle: recrutTitle,
				EmplTyp:     emplTyp,
				UpJobGrp:    upJobGrp,
				JobGrp:      jobGrp,
				RecrutDy:    recrutDy,
				RecrutEdt:   recrutEdt,
				ApplyCnt:    applyCnt,
				IngCnt:      ingCnt,
				PassCnt:     passCnt,
				FailCnt:     failCnt,
				TotCnt:      totCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnRecruitPostList = models.RtnRecruitPostList{
			RtnRecruitPostListData: recruitPostList,
		}

		c.Data["json"] = &rtnRecruitPostList
		c.ServeJSON()
	}

}
