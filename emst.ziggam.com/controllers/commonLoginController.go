package controllers

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"emst.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

type CommonLoginController struct {
	beego.Controller
}

func (c *CommonLoginController) Get() {
	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log
	session := c.StartSession()

	mem_id := session.Get(c.Ctx.Request.Context(), "mem_id")
	if mem_id != nil {
		c.Ctx.Redirect(302, "/")
	}

	var saveId string

	cookieMemId := c.Ctx.GetCookie("ck_mem_id")
	if cookieMemId == "" {
		saveId = "N"
	} else {
		saveId = "Y"
	}
	//saveId := c.GetString("save_id")

	rtnUrl := c.GetString("rtnurl")
	rtnUrl = strings.Replace(rtnUrl, "＆", "&", -1)
	log.Debug("return URL = %v", rtnUrl)

	c.Data["SaveId"] = saveId
	c.Data["RtnUrl"] = rtnUrl

	c.TplName = "common/login.html"
}

func (c *CommonLoginController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()
	pLang, _ := beego.AppConfig.String("lang")

	pMemTypCd := "E"
	pMemId := c.GetString("mem_id")
	pPwd := c.GetString("pwd")
	pSaveId := c.GetString("save_id")

	sha := sha512.New()
	sha.Write([]byte(pPwd))
	sha_pPwd := hex.EncodeToString(sha.Sum(nil))

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	log.Debug("CALL SP_EMS_LOGIN_PROC( '%v', '%v', '%v', '%v', :1)",
		pLang, pMemTypCd, pMemId, sha_pPwd)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_LOGIN_PROC( '%v', '%v', '%v', '%v', :1)",
		pLang, pMemTypCd, pMemId, sha_pPwd),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
		ora.S,   /* MEM_NO */
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

				login = models.Login{
					MemNo: memNo,
					MemId: memId,
				}

				// Set the session
				session.Set(c.Ctx.Request.Context(), "mem_no", memNo)
				session.Set(c.Ctx.Request.Context(), "mem_id", memId)

				if pSaveId == "Y" {
					// Set the Cookie
					c.Ctx.SetCookie("ck_mem_id", memId, 1<<31-1, "/")
				} else {
					c.Ctx.SetCookie("ck_mem_id", "", 1<<31-1, "/")
				}

			} else if rtnCd == 5 {
				memNo = procRset.Row[2].(string)
				memId = procRset.Row[3].(string)

				login = models.Login{
					MemNo: memNo,
					MemId: memId,
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
	c.ServeJSON()
}
