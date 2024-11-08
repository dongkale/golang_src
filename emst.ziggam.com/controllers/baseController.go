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

	"emst.ziggam.com/models"
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
		logs.Error("Failed to get logConf:", err)
	} else {
		logs.Debug("[logConf] " + logConf)
	}

	if logConf != "" {
		logs.SetLogger(logs.AdapterFile, logConf)
	} else {
		logs.SetLogger(logs.AdapterFile, `{"filename":"./logs/ems.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":180,"color":true}`)
	}

	// Initialize language type list.
	langTypesStr, err := beego.AppConfig.String("lang_types")
	if err != nil {
		logs.Error("Failed to get lang_types:", err)
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
}

// Prepare implemented Prepare() method for baseController.
// It's used for language option check and setting.
func (c *BaseController) Prepare() {
	// start : log
	// log := logs.NewLogger()
	// log.SetLogger(logs.AdapterConsole)
	// end : log

	hasCookie := false

	// Reset language option.
	c.Lang = "" // This field is from i18n.Locale.
	logs.Trace("running prepare")

	var lang string

	// 1. Check URL arguments.
	input, err := c.Input()
	if err != nil {
		logs.Error("Failed to get input:", err)
		lang = ""
	} else {
		lang = input.Get("lang")
	}

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

	//start:session
	session := c.StartSession()

	// 세션 기업회원번호 가져오기
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	mem_id := session.Get(c.Ctx.Request.Context(), "mem_id")

	logs.Debug(fmt.Sprintf("mem_no:%v, mem_id:%v", mem_no, mem_id))

	if mem_no != nil {
		c.Data["SMemNo"] = mem_no
		c.Data["SMemId"] = mem_id
	} else {
		c.Data["SMemNo"] = ""
		c.Data["SMemId"] = ""
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

	// Start : 기업기본정보
	logs.Debug("CALL SP_EMS_ENTP_BASIC_INFO_R('%v', '%v', :1)",
		pLang, mem_no)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ENTP_BASIC_INFO_R('%v', '%v', :1)",
		pLang, mem_no),
		ora.S,   /* ENTP_KO_NM */
		ora.S,   /* LOGO_PTO_PATH */
		ora.S,   /* NOTI_YN */
		ora.S,   /* ZIP */
		ora.S,   /* EST_DY */
		ora.I64, /* EMP_CNT */
		ora.S,   /* BIZ_TPY */
		ora.S,   /* BIZ_COND */
		ora.S,   /* PP_CHRG_NM */
		ora.S,   /* PP_CHRG_TEL_NO */
		ora.S,   /* ENTP_HTAG */
		ora.S,   /* ENTP_INTR */
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
		entpKoNm         string
		logoPtoPath      string
		notificationYn   string
		basicZip         string
		basicEstDy       string
		basicEmpCnt      int64
		basicBizTpy      string
		basicBizCond     string
		basicPpChrgNm    string
		basicPpChrgTelNo string
		basicEntpTag     string
		basicEntpIntr    string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			entpKoNm = procRset.Row[0].(string)
			logoPtoPath = procRset.Row[1].(string)
			notificationYn = procRset.Row[2].(string)
			basicZip = procRset.Row[3].(string)
			basicEstDy = procRset.Row[4].(string)
			basicEmpCnt = procRset.Row[5].(int64)
			basicBizTpy = procRset.Row[6].(string)
			basicBizCond = procRset.Row[7].(string)
			basicPpChrgNm = procRset.Row[8].(string)
			basicPpChrgTelNo = procRset.Row[9].(string)
			basicEntpTag = procRset.Row[10].(string)
			basicEntpIntr = procRset.Row[11].(string)

			var fullPtoPath string

			if logoPtoPath == "" {
				fullPtoPath = logoPtoPath
			} else {
				fullPtoPath = imgServer + logoPtoPath
			}

			entpBasicInfo = append(entpBasicInfo, models.EntpBasicInfo{
				EntpKoNm:         entpKoNm,
				LogoPtoPath:      fullPtoPath,
				NotificationYn:   notificationYn,
				BasicZip:         basicZip,
				BasicEstDy:       basicEstDy,
				BasicEmpCnt:      basicEmpCnt,
				BasicBizTpy:      basicBizTpy,
				BasicBizCond:     basicBizCond,
				BasicPpChrgNm:    basicPpChrgNm,
				BasicPpChrgTelNo: basicPpChrgTelNo,
				BasicEntpTag:     basicEntpTag,
				BasicEntpIntr:    basicEntpIntr,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : 기업기본정보

	c.Data["BasicEntpKoNm"] = entpKoNm
	c.Data["BasicNotificationYn"] = notificationYn
	c.Data["BasicLogoPtoPath"] = imgServer + logoPtoPath

	c.Data["BasicZip"] = basicZip
	c.Data["BasicEstDy"] = basicEstDy
	c.Data["BasicEmpCnt"] = basicEmpCnt
	c.Data["BasicBizTpy"] = basicBizTpy
	c.Data["BasicBizCond"] = basicBizCond
	c.Data["BasicPpChrgNm"] = basicPpChrgNm
	c.Data["BasicPpChrgTelNo"] = basicPpChrgTelNo
	c.Data["BasicEntpTag"] = basicEntpTag
	c.Data["BasicEntpIntr"] = basicEntpIntr

	c.Data["MenuId"] = "00"

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
