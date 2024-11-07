package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

// InviteRefuseController ...
type InviteRefuseController struct {
	beego.Controller
}

// Post ...
func (c *InviteRefuseController) Post() {

	//session := c.StartSession()

	// memNO := session.Get(c.Ctx.Request.Context(), "mem_no")
	// if memNO == nil {
	// 	c.Ctx.Redirect(302, "/login")
	// }

	// memSn := session.Get(c.Ctx.Request.Context(), "mem_sn")
	// if memSn == nil {
	// 	c.Ctx.Redirect(302, "/login")
	// }

	// http://localhost:7070/bridge?entpmemno=E2019011900001&recruitsn=2019010001&reqname=%EC%9D%B4%EB%8F%99%EA%B4%80&reqmono=010-5226-2107&reqemail=dongkale@naver.com

	// memSn := session.Get(c.Ctx.Request.Context(), "mem_sn")
	// if memSn == nil {
	// 	//rtnData := models.DefaultResult{ RtnCd: 2, RtnMsg: "rtnMsg" }
	// 	c.Data["json"] = &models.DefaultResult{RtnCd: 99, RtnMsg: "/login"}
	// 	c.ServeJSON()
	// 	return
	// }

	pLang, _ := beego.AppConfig.String("lang")

	//pEntpMemNo := memNO // 기업회원번호
	//pRecrutSn := c.GetString("recrut_sn") // 채용일련번호

	pReqType := c.GetString("req_type")

	pName := c.GetString("req_name")
	pMono := c.GetString("req_mono")
	pEmail := c.GetString("req_email")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	var (
		rtnCd  int64
		rtnMsg string
	)

	// Start : Invite Refuse Insert
	fmt.Printf(fmt.Sprintf("CALL ZSP_INVITE_REFUSE_REG('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pReqType, pName, pEmail, pMono))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_INVITE_REFUSE_REG('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pReqType, pName, pEmail, pMono),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
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

	rtnData := models.DefaultResult{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}

		if err := procRset.Err(); err != nil {
			panic(err)
		}

		fmt.Printf(fmt.Sprintf(" ===> rtnCd:%v, rtnMsg:%v", rtnCd, rtnMsg))

		rtnData = models.DefaultResult{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}
	// End : Invite Refuse Insert

	c.Data["json"] = &rtnData
	c.ServeJSON()

}
