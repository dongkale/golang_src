package controllers

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/beego/beego/v2/core/logs"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
	"gopkg.in/errgo.v1"

	"bizt.ziggam.com/models"
	"bizt.ziggam.com/tables"
	"bizt.ziggam.com/utils"
	ora "gopkg.in/rana/ora.v4"
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

	//beego.Trace("[logConf] " + beego.AppConfig.String("logConf"))
	logConf, err := beego.AppConfig.String("logConf")
	if err != nil {
		logs.Error("Error getting logConf:", err)
	} else {
		fmt.Printf("[logConf] " + logConf)
	}

	// 1 Case
	//logs.SetLogger(logs.AdapterFile, `{"filename":"./logs/biz.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":180,"color":true}`)
	/* 2 Case
	if beego.AppConfig.String("logFile") != "" {
		logs.SetLogger(logs.AdapterFile, `{"filename": "`+beego.AppConfig.String("logFile")+`","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":180,"color":true}`)
	} else {
		logs.SetLogger(logs.AdapterFile, `{"filename":"./logs/biz.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":180,"color":true}`)
	}
	*/

	if logConf != "" {
		logs.SetLogger(logs.AdapterFile, logConf)
	} else {
		logs.SetLogger(logs.AdapterFile, `{"filename":"./logs/biz.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":180,"color":true}`)
	}

	// Initialize language type list.
	langTypesStr, err := beego.AppConfig.String("lang_types")
	if err != nil {
		logs.Error("Error getting lang_types:", err)
	} else {
		langTypes = strings.Split(langTypesStr, "|")
	}

	// Load locale files according to language types.
	for _, lang := range langTypes {
		logs.Trace("Loading language: " + lang)
		if err := i18n.SetMessage(lang, "conf/"+"locale_"+lang+".ini"); err != nil {
			logs.Error("Fail to set message file:", err)
			return
		}
	}

	// // 레디스 설정 파일 읽어 오기.
	// redis_addr, err := beego.AppConfig.String("redisAddr")
	// if err != nil {
	// 	logs.Error("Error getting redisAddr:", err)
	// 	return
	// }
	// redis_pass, err := beego.AppConfig.String("redisPass")
	// if err != nil {
	// 	logs.Error("Error getting redisPass:", err)
	// 	return
	// }

	// // 레디스 풀 설정.
	// if utils.RPool.NewPool(redis_addr, redis_pass) == nil {
	// 	logs.Error("redis pool connect config err!!!!!!!!!! redis_addr:", redis_addr, redis_pass)
	// 	return
	// }

	// // 한개 접속 해서 접속이 제대로 되는지 확인 하자.
	// err = utils.RPool.Ping()
	// if err != nil {
	// 	logs.Error("redis connect config err!!!!!!!!!! err:", err)
	// 	return
	// }
	
	// utils.RPool.HSet("asdfasd", "sadfsadf", "1213213")

	// utils.RPool.Expire("asdfasd", 1213213)

	// utils.RPool.HSet("asdfasd1", "sadfsadf", "1213214")

	// utils.RPool.HSet("asdfasd1", "sadfsadf2", "1213215")
	// conn := pool.Get()
	// defer conn.Close()

	setConfFile, err := beego.AppConfig.String("setConfFile")
	if err != nil {
		logs.Error("Error getting setConfFile:", err)
		return
	}
	errCnt := utils.LoadConfigFile(setConfFile, &tables.TableConf)
	if errCnt == 0 {
		result := tables.TableConf.IsCheckValue()
		if result != "" {
			logs.Info(fmt.Sprintf("[LoadConfig] Load Fail!! -> %s InValid Value", result))
		} else {
			logs.Info(fmt.Sprintf("[LoadConfig] Load Ok!!"))
		}
	} else {
		logs.Info(fmt.Sprintf("[LoadConfig] Load Fail!!"))
	}
}

// Prepare implemented Prepare() method for baseController.
// It's used for language option check and setting.
func (c *BaseController) Prepare() {
	// start : log
	//log := logs.NewLogger()
	//logs.SetLogger(logs.AdapterFile,`{"filename":"logs\biz.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
	// end : log

	hasCookie := false

	// Reset language option.
	c.Lang = "" // This field is from i18n.Locale.
	logs.Trace("running prepare")

	// 1. Check URL arguments.
	input, err := c.Input()
	if err != nil {
		logs.Error("Error getting input:", err)
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
	pLang, _ := beego.AppConfig.String("lang")

	/* Image Server Path */
	imgServer, _ := beego.AppConfig.String("viewpath")

	// start : session
	session := c.StartSession()
	mem_id := session.Get(c.Ctx.Request.Context(), "mem_id")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// 로그인 유지 체크(세션이 없을 경우 쿠키 및 토큰으로 유지 체크)
	if mem_id == nil {

		cookieToken := c.Ctx.GetCookie("token")

		if cookieToken != "" {
			// Start : 로그인유지정보
			fmt.Printf(fmt.Sprintf("CALL ZSP_LOGIN_KEEP_R('%v', '%v', :1)",
				pLang, cookieToken))

			stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_LOGIN_KEEP_R('%v', '%v', :1)",
				pLang, cookieToken),
				ora.S, /* RTN_CD */
				ora.S, /* RTN_MSG */
				ora.S, /* MEM_ID */
				ora.S, /* MEM_NO */
				ora.S, /* MEM_SN */
				ora.S, /* AUTH_CD */
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

			loginKeepInfo := make([]models.LoginKeepInfo, 0)

			var (
				rtnCd     string
				rtnMsg    string
				rtnMemNo  string
				rtnMemId  string
				rtnMemSn  string
				rtnAuthCd string
			)

			if procRset.IsOpen() {
				for procRset.Next() {
					rtnCd = procRset.Row[0].(string)
					rtnMsg = procRset.Row[1].(string)

					if rtnCd == "1" {
						rtnMemNo = procRset.Row[2].(string)
						rtnMemId = procRset.Row[3].(string)
						rtnMemSn = procRset.Row[4].(string)
						rtnAuthCd = procRset.Row[5].(string)

						session.Set(c.Ctx.Request.Context(), "mem_id", rtnMemNo)
						session.Set(c.Ctx.Request.Context(), "mem_no", rtnMemId)
						session.Set(c.Ctx.Request.Context(), "mem_sn", rtnMemSn)
						session.Set(c.Ctx.Request.Context(), "auth_cd", rtnAuthCd)

					} else if rtnCd == "9" {

						rtnMemNo = ""
						rtnMemId = ""
						rtnMemSn = ""
						rtnAuthCd = ""

						session.Delete(c.Ctx.Request.Context(), "mem_id")
						session.Delete(c.Ctx.Request.Context(), "mem_no")
						session.Delete(c.Ctx.Request.Context(), "mem_sn")
						session.Delete(c.Ctx.Request.Context(), "auth_cd")

						c.Ctx.SetCookie("token", "", 1<<31-1, "/")
						c.Ctx.Redirect(302, "/login")
					}

					loginKeepInfo = append(loginKeepInfo, models.LoginKeepInfo{
						RtnCd:     rtnCd,
						RtnMsg:    rtnMsg,
						RtnMemNo:  rtnMemNo,
						RtnMemId:  rtnMemId,
						RtnMemSn:  rtnMemSn,
						RtnAuthCd: rtnAuthCd,
					})
				}
				if err := procRset.Err(); err != nil {
					panic(err)
				}
			}
		}
		// End : 로그인유지정보
	}

	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	mem_sn := session.Get(c.Ctx.Request.Context(), "mem_sn")
	auth_cd := session.Get(c.Ctx.Request.Context(), "auth_cd")

	if mem_no != nil {
		c.Data["SMemId"] = mem_id
		c.Data["SMemNo"] = mem_no
		c.Data["SMemSn"] = mem_sn
		c.Data["SAuthCd"] = auth_cd
	} else {
		c.Data["SMemId"] = ""
		c.Data["SMemNo"] = ""
		c.Data["SMemSn"] = ""
		c.Data["SAuthCd"] = ""
	}

	// Start : 기업기본정보
	fmt.Printf(fmt.Sprintf("CALL ZSP_ENTP_BASIC_INFO_R('%v', '%v', '%v', :1)",
		pLang, mem_no, mem_id))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_ENTP_BASIC_INFO_R('%v', '%v', '%v', :1)",
		pLang, mem_no, mem_id),
		ora.S,   /* ENTP_KO_NM */
		ora.S,   /* LOGO_PTO_PATH */
		ora.S,   /* PP_CHRG_NM */
		ora.S,   /* PP_CHRG_BP_NM */
		ora.I64, /* MSG_CNT */
		ora.I64, /* NOTI_CNT */
		ora.S,   /* EMAIL */
		ora.S,   /* BIZ_REG_NO */
		ora.S,   /* REP_NM */
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

	entpBasicInfo := make([]models.EntpBasicInfo, 0)

	var (
		bsEntpKoNm    string
		bsLogoPtoPath string
		bsPpChrgNm    string
		bsPpChrgBpNm  string
		bsMsgCnt      int64
		bsNotiCnt     int64
		bsEmail       string
		bsBizRegNo    string
		bsRepNm       string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			bsEntpKoNm = procRset.Row[0].(string)
			bsLogoPtoPath = procRset.Row[1].(string)
			bsPpChrgNm = procRset.Row[2].(string)
			bsPpChrgBpNm = procRset.Row[3].(string)
			bsMsgCnt = procRset.Row[4].(int64)
			bsNotiCnt = procRset.Row[5].(int64)
			bsEmail = procRset.Row[6].(string)
			bsBizRegNo = procRset.Row[7].(string)
			bsRepNm = procRset.Row[8].(string)

			var fullPtoPath string

			if bsLogoPtoPath == "" {
				fullPtoPath = bsLogoPtoPath
			} else {
				fullPtoPath = imgServer + bsLogoPtoPath
			}

			entpBasicInfo = append(entpBasicInfo, models.EntpBasicInfo{
				BsEntpKoNm:    bsEntpKoNm,
				BsLogoPtoPath: fullPtoPath,
				BsPpChrgNm:    bsPpChrgNm,
				BsPpChrgBpNm:  bsPpChrgBpNm,
				BsMsgCnt:      bsMsgCnt,
				BsNotiCnt:     bsNotiCnt,
				BsEmail:       bsEmail,
				BsBizRegNo:    bsBizRegNo,
				BsRepNm:       bsRepNm,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : 기업기본정보

	c.Data["BsEntpKoNm"] = bsEntpKoNm
	if bsLogoPtoPath == "" {
		c.Data["BsLogoPtoPath"] = bsLogoPtoPath
	} else {
		c.Data["BsLogoPtoPath"] = imgServer + bsLogoPtoPath
	}
	c.Data["BsPpChrgNm"] = bsPpChrgNm
	c.Data["BsPpChrgBpNm"] = bsPpChrgBpNm
	c.Data["BsMsgCnt"] = bsMsgCnt
	c.Data["BsNotiCnt"] = bsNotiCnt
	c.Data["BsEmail"] = bsEmail
	c.Data["BsBizRegNo"] = bsBizRegNo
	c.Data["BsRepNm"] = bsRepNm

	//c.Data["TblConfigJson"] = utils.ToJsonString(tables.TableConf)
	c.Data["TblConfig"] = tables.TableConf

	c.Data["SuperAdmin"] = c.GetSession("super_admin")
	fmt.Printf(fmt.Sprintf("SuperAdmin: %v", c.GetSession("super_admin")))

	c.Ctx.ResponseWriter.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
	c.Ctx.ResponseWriter.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
	c.Ctx.ResponseWriter.Header().Set("Expires", "0")                                         // Proxies.
}

var (
	oraCxMu sync.Mutex
	oraInit sync.Once
)

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

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
