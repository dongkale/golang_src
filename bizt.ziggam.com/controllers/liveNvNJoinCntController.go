package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/rana/ora.v4"
)

type LiveNvNJoinCntController struct {
	beego.Controller
}

func (c *LiveNvNJoinCntController) Post() {

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Data["json"] = &models.RtnLiveNvnJoinCnt{
			Rtn:     models.DefaultResult{RtnCd: 99, RtnMsg: "mem_no == nil"},
			JoinCnt: 0,
		}

		c.ServeJSON()
		return
	}

	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := mem_no
	pRecrutSn := c.GetString("recrut_sn")
	pLiveSn := c.GetString("live_sn")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : join count
	fmt.Printf(fmt.Sprintf("CALL ZSP_LIVE_NVN_SEL_JOIN_CNT('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pLiveSn))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_NVN_SEL_JOIN_CNT('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pLiveSn),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
		ora.I64, /* RTN_JOIN_CNT */
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
		rtnCd      int64
		rtnMsg     string
		rtnJoinCnt int64
	)

	rtnJoinCntInfo := models.RtnLiveNvnJoinCnt{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
			rtnJoinCnt = procRset.Row[2].(int64)
		}

		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnJoinCntInfo = models.RtnLiveNvnJoinCnt{
			Rtn:     models.DefaultResult{RtnCd: rtnCd, RtnMsg: rtnMsg},
			JoinCnt: rtnJoinCnt,
		}

		fmt.Printf(fmt.Sprintf("===> rtnCd:%v, rtnMsg:%v, rtnJoinCnt:%v", rtnCd, rtnMsg, rtnJoinCnt))
	}
	// End : join count

	c.Data["json"] = &rtnJoinCntInfo
	c.ServeJSON()
}
