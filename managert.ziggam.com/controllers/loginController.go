package controllers

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	// "bizt.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
	models "managert.ziggam.com/bizt.ziggam.com-models"
)

type LoginController struct {
	BaseController
}

func (c *LoginController) Get() {
	// start : log
	//log := logs.NewLogger()
	//log.SetLogger(logs.AdapterConsole)

	//log2 := logs.NewLogger()
	//log2.SetLogger(logs.AdapterFile, `{"filename":"test.log"}`)
	//
	//log2.Alert("Test Log!!!!")
	// end : log

	session := c.StartSession()

	mem_id := session.Get(c.Ctx.Request.Context(), "mem_id")
	if mem_id != nil {
		c.Ctx.Redirect(302, "/")
		return 
	}

	// mem_id__ := c.GetSession("mem_id")
	// fmt.Println(mem_id__)

	// // TEST
	// session.Set(c.Ctx.Request.Context(), "mem_no", "E2020062600783")
	// session.Set(c.Ctx.Request.Context(), "mem_id", "likemandoo")
	// session.Set(c.Ctx.Request.Context(), "mem_sn", "0001")
	// session.Set(c.Ctx.Request.Context(), "auth_cd", "01")
	// // TEST

	rtnUrl := c.GetString("rtnurl")
	rtnUrl = strings.Replace(rtnUrl, "＆", "&", -1)

	c.Data["RtnUrl"] = rtnUrl

	c.TplName = "main/login.html"
}

func (c *LoginController) Post() {
	// // start : log
	// log := logs.NewLogger()
	// log.SetLogger(logs.AdapterConsole)
	// // end : log

	session := c.Ctx.Input.CruSession
	//session := c.StartSession()

	mem_id := session.Get(c.Ctx.Request.Context(), "mem_id")
	if mem_id != nil {
		c.Ctx.Redirect(302, "/")
		return 
	}

	pLang, _ := beego.AppConfig.String("lang")

	pMemTypCd := "E"
	pMemId := c.GetString("mem_id")
	pPwd := c.GetString("pwd")
	pLoginMtnYn := c.GetString("login_mtn_yn")
	pIpAddress := c.GetString("ip_addr")
	pOsGbn := c.GetString("os_gbn")

	var shaTknKey string
	if pLoginMtnYn == "Y" {
		sha_tkn := sha512.New()
		sha_tkn.Write([]byte(pMemId))
		shaTknKey = hex.EncodeToString(sha_tkn.Sum(nil))
		c.Ctx.SetCookie("token", shaTknKey, 1<<31-1, "/")
	} else {
		c.Ctx.SetCookie("token", "", 1<<31-1, "/")
		shaTknKey = ""
	}

	sha := sha512.New()
	sha.Write([]byte(pPwd))
	sha_pPwd := hex.EncodeToString(sha.Sum(nil))

	logs.Debug(sha_pPwd)

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	logs.Debug(fmt.Sprintf("CALL ZSP_LOGIN_PROC( '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pMemTypCd, pMemId, sha_pPwd, pLoginMtnYn, shaTknKey, pOsGbn, pIpAddress))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_LOGIN_PROC( '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pMemTypCd, pMemId, sha_pPwd, pLoginMtnYn, shaTknKey, pOsGbn, pIpAddress),
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

	login := models.Login{}
	rtnLogin := models.RtnLogin{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)

			if rtnCd == 1 {
				memNo = procRset.Row[2].(string)
				memId = procRset.Row[3].(string)
				memSn = procRset.Row[4].(string)
				authCd = procRset.Row[5].(string)

				login = models.Login{
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
			} else if rtnCd == 5 {
				memNo = procRset.Row[2].(string)
				memId = procRset.Row[3].(string)
				memSn = procRset.Row[4].(string)
				authCd = procRset.Row[5].(string)

				login = models.Login{
					MemNo:  memNo,
					MemId:  memId,
					MemSn:  memSn,
					AuthCd: authCd,
				}
			}
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnLogin = models.RtnLogin{
			RtnCd:   rtnCd,
			RtnMsg:  rtnMsg,
			RtnData: login,
		}
	}

	c.Data["json"] = &rtnLogin
	c.TplName = "main/main.html"
//	c.ServeJSON()
}
