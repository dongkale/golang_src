package controllers

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"bizt.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

type AdminSimplePassLoginController struct {
	beego.Controller
}

func (c *AdminSimplePassLoginController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()
	pLang, err := beego.AppConfig.String("lang")
	if err != nil {
		panic(err)
	}
	pMemId := c.GetString("mem_id")

	//pSuperAdmin := c.GetString("super_admin")
	//if pSuperAdmin != "" {
	c.SetSession("super_admin", pMemId)
	log.Debug(fmt.Sprintf("[Super Admin] %v", pMemId))
	//}

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	log.Debug(fmt.Sprintf("CALL SP_ZSP_ADMIN_LOGIN_PROC( '%v', '%v', :1)",
		pLang, pMemId))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_ZSP_ADMIN_LOGIN_PROC( '%v', '%v', :1)",
		pLang, pMemId),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
		ora.S,   /* MEM_NO */
		ora.S,   /* MEM_ID */
		ora.S,   /* MEM_SN */
		ora.S,   /* AUTH_CD */
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
		memNo  string
		memId  string
		memSn  string
		authCd string
	)

	adminSimplePassLogin := models.AdminSimplePassLogin{}
	rtnAdminSimplePassLogin := models.RtnAdminSimplePassLogin{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)

			if rtnCd == 1 {
				memNo = procRset.Row[2].(string)
				memId = procRset.Row[3].(string)
				memSn = procRset.Row[4].(string)
				authCd = procRset.Row[5].(string)

				adminSimplePassLogin = models.AdminSimplePassLogin{
					MemNo:  memNo,
					MemId:  memId,
					MemSn:  memSn,
					AuthCd: authCd,
				}

				// Set the session
				session.Set(c.Ctx.Request.Context(), "mem_no", memNo)
				session.Set(c.Ctx.Request.Context(), "mem_id", memId)
				session.Set(c.Ctx.Request.Context(), "mem_sn", memSn)
				session.Set(c.Ctx.Request.Context(), "auth_cd", authCd)

				fmt.Printf(fmt.Sprintf("SessionSet: mem_no:%v, mem_id:%v, mem_sn:%v, auth_cd:%v",
					memNo, memId, memSn, authCd))
			}
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnAdminSimplePassLogin = models.RtnAdminSimplePassLogin{
			RtnCd:   rtnCd,
			RtnMsg:  rtnMsg,
			RtnData: adminSimplePassLogin,
		}
	}

	c.Data["json"] = &rtnAdminSimplePassLogin
	c.ServeJSON()
}
