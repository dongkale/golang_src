package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	"gopkg.in/rana/ora.v4"
)

type InquiryInsertController struct {
	beego.Controller
}

func (c *InquiryInsertController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	mem_sn := session.Get(c.Ctx.Request.Context(), "mem_sn")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/common/login")
	}

	pLang, _ := beego.AppConfig.String("lang")

	pEntpMemNo := mem_no                   // 기업회원번호
	pPpChrgSn := mem_sn                    // 담당자순번
	pInqGbnCd := c.GetString("inq_gbn_cd") // 문의종류코드
	pInqTitle := c.GetString("inq_title")  //문의제목
	pInqCont := c.GetString("inq_cont")    // 문의내용
	pEmail := c.GetString("email")         // 이메일

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Inquiry Insert
	fmt.Printf(fmt.Sprintf("CALL ZSP_INQUIRY_REG_PROC('%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pPpChrgSn, pInqGbnCd, pInqTitle, pInqCont, pEmail))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_INQUIRY_REG_PROC('%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pPpChrgSn, pInqGbnCd, pInqTitle, pInqCont, pEmail),
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

	var (
		rtnCd  int64
		rtnMsg string
	)

	rtnInquiry := models.RtnInquiry{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnInquiry = models.RtnInquiry{
			RtnCd:  rtnCd,
			RtnMsg: rtnMsg,
		}
	}
	// End : Inquiry Insert

	c.Data["json"] = &rtnInquiry
	c.ServeJSON()

}
