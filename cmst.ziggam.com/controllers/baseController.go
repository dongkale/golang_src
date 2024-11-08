package controllers

import (
	"strings"
	"sync"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
	"gopkg.in/errgo.v1"
	"gopkg.in/rana/ora.v4"
)

var langTypes []string // Languages that are supported.

// baseController represents base router for all other app routers.
// It implemented some methods for the same implementation;
// thus, it will be embedded into other routers.
type BaseController struct {
	beego.Controller // Embed struct that has stub implementation of the interface.
	i18n.Locale      // For i18n usage when process data and render template.
}

func init() {
	beego.AddFuncMap("i18n", i18n.Tr)

	// Initialize language type list.
	langTypesStr, err := beego.AppConfig.String("lang_types")
	if err != nil {
		logs.Error("Fail to get lang_types:", err)
		return
	}
	langTypes = strings.Split(langTypesStr, "|")

	// Load locale files according to language types.
	for _, lang := range langTypes {
		logs.Trace("Loading language: " + lang)
		if err := i18n.SetMessage(lang, "conf/"+"locale_"+lang+".ini"); err != nil {
			logs.Error("Fail to set message file:", err)
			return
		}
	}
}

// Prepare implemented Prepare() method for baseController.
// It's used for language option check and setting.
func (c *BaseController) Prepare() {
	// // start : log
	// log := logs.NewLogger()
	// log.SetLogger(logs.AdapterConsole)
	// // end : log

	hasCookie := false

	// Reset language option.
	c.Lang = "" // This field is from i18n.Locale.
	logs.Trace("running prepare")

	// 1. Check URL arguments.
	input, err := c.Input()
	if err != nil {
		logs.Error("Failed to get input:", err)
		return
	}
	lang := input.Get("lang")

	// 2. Get language information from cookies.
	if len(lang) == 0 {
		lang = c.Ctx.GetCookie("lang")
		hasCookie = true
	}

	// Check again in case someone modify on purpose.
	if !i18n.IsExist(lang) {
		lang = ""
		hasCookie = false
	}

	// 3. Get language information from 'Accept-Language'.
	/*if len(lang) == 0 {
		al := c.Ctx.Request.Header.Get("Accept-Language")
		if len(al) > 4 {
			al = al[:5] // Only compare first 5 letters.
			if i18n.IsExist(al) {
				lang = al
			}
		}
		beego.Trace("Accept-Language is " + al)
	}*/

	// 4. Default language is English.
	if len(lang) == 0 {
		lang = "en-US"
	}

	// Save language information in cookies.
	if !hasCookie {
		c.Ctx.SetCookie("lang", lang, 1<<31-1, "/")
	}
	// Set language properties.
	c.Lang = lang

	// Set template level language option.
	c.Data["Lang"] = c.Lang

	// start : session
	// 1. dvc_id checking
	session := c.StartSession()
	//dvc_id := session.Get("dvc_id")
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	admin_super_mem_yn := session.Get(c.Ctx.Request.Context(), "admin_super_mem_yn")
	mem_nm := session.Get(c.Ctx.Request.Context(), "mem_nm")
	email := session.Get(c.Ctx.Request.Context(), "email")

	var loginTitle string
	var loginType string
	if mem_no != nil {
		//loginTitle = i18n.Tr(c.Lang, "mypage.uppercase") //"MYPAGE"
		loginTitle = "마이페이지"
		loginType = "mypage"
		c.Data["MemNo"] = mem_no
		c.Data["AdminSuperMemYn"] = admin_super_mem_yn
		c.Data["MemNm"] = mem_nm
		c.Data["Email"] = email
	} else {
		//loginTitle = i18n.Tr(c.Lang, "login.uppercase") //"LOGIN"
		loginTitle = "로그인"
		loginType = "login"
		c.Data["MemNo"] = ""
		c.Data["AdminSuperMemYn"] = ""
		admin_super_mem_yn = ""
		mem_nm = ""
		email = ""
		mem_no = ""
	}

	c.Data["LoginTitle"] = loginTitle
	c.Data["LoginType"] = loginType

	pRecmdMemNo := c.GetString("rmn")
	c.Data["RecmdMemNo"] = pRecmdMemNo

	/* 언어선택 */
	//pLang := beego.AppConfig.String("lang")
	/* Title명 정의 */
	pWebTitle, _ := beego.AppConfig.String("webtitle")
	c.Data["WebTitle"] = pWebTitle

	/* Site URL 정의 */
	siteUrl, _ := beego.AppConfig.String("siteurl")
	c.Data["SiteUrl"] = siteUrl
	// end : session

	// Start : Oracle DB Connection
	/*
		env, srv, ses, err := GetRawConnection()
		defer env.Close()
		defer srv.Close()
		defer ses.Close()
		if err != nil {
			panic(err)
		}
	*/
	// End : Oracle DB Connection

	c.Ctx.ResponseWriter.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
	c.Ctx.ResponseWriter.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
	c.Ctx.ResponseWriter.Header().Set("Expires", "0")                                         // Proxies.
}

var (
	oraCxMu sync.Mutex
	oraInit sync.Once
)

// GetRawConnection returns a raw (*ora.Ses) connection
func GetRawConnection() (*ora.Env, *ora.Srv, *ora.Ses, error) {
	oraCxMu.Lock()
	defer oraCxMu.Unlock()

	env, err := ora.OpenEnv()
	if err != nil {
		return nil, nil, nil, errgo.Notef(err, "OpenEnv")
	}

	dsn, _ := beego.AppConfig.String("oradsn")
	dsn = strings.TrimSpace(dsn)

	srvCfg := ora.SrvCfg{StmtCfg: env.Cfg()}
	sesCfg := ora.SesCfg{Mode: ora.DSNMode(dsn)}
	sesCfg.Username, sesCfg.Password, srvCfg.Dblink = ora.SplitDSN(dsn)

	srv, err := env.OpenSrv(srvCfg)
	if err != nil {
		return nil, nil, nil, errgo.Notef(err, "OpenSrv(%#v)", srvCfg)
	}

	ses, err := srv.OpenSes(sesCfg)
	if err != nil {
		srv.Close()
		return nil, nil, nil, errgo.Notef(err, "OpenSes(%#v)", sesCfg)
	}

	return env, srv, ses, nil
}
