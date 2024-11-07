package controllers

import (
	"fmt"
	"unicode/utf8"

	"bizt.ziggam.com/utils"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
)

//var langTypes []string // Languages that are supported.

// baseController represents base router for all other app routers.
// It implemented some methods for the same implementation;
// thus, it will be embedded into other routers.
type ApiBaseController struct {
	beego.Controller // Embed struct that has stub implementation of the interface.
	i18n.Locale      // For i18n usage when process data and render template.
}

func init() {

	fmt.Printf("[ApiBase] Init")

	apiServerUrl, err := beego.AppConfig.String("apiServerUrl")
	if err != nil || apiServerUrl == "" {
		logs.Error(fmt.Sprintf("[ApiBase] Error: %s", "apiServerUrl"))
		return
	}
}

// Prepare implemented Prepare() method for baseController.
// It's used for language option check and setting.
func (c *ApiBaseController) Prepare() {

	hasCookie := false

	// Reset language option.
	c.Lang = "" // This field is from i18n.Locale.
	logs.Trace("[ApiBase] running prepare")

	//beego.Trace(c.Ctx.Input.Params())

	//pam := c.Ctx.Input.Params()
	// beego.Trace(pam)

	logs.Trace(c.Ctx.Input.Header)

	// for i, v := range c.Ctx.Input.Params() {
	// 	beego.Trace(i)
	// 	beego.Trace(v)
	// }

	// 1. Check URL arguments.
	input, _ := c.Input()
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
	//pLang, _ := beego.AppConfig.String("lang")

	/* Image Server Path */
	//imgServer, _  := beego.AppConfig.String("viewpath")

	// start : session
	//session := c.StartSession()
}

func (c *ApiBaseController) Finish() {
	logs.Trace("[ApiBase] running finish")
}

// GetString ...
func (c *ApiBaseController) GetString(key string, def ...string) string {

	getStr := c.Controller.GetString(key)

	//fmt.Printf(fmt.Sprintf("[ApiBase] Key:%s, Value:%s, UTF-8:%t", key, getStr, utf8.Valid([]byte(getStr))))

	if utf8.Valid([]byte(getStr)) == true {
		return getStr
	} else {
		return utils.ConvertEucKR(getStr)
	}

	//fmt.Println(key + ":" + utf8.Valid([]byte(str)))
	// beego.Trace("[ApiBase] Key:%v -> %s", key, utf8.Valid([]byte(str)))
	// return utils.ConvertEucKR(c.Controller.GetString(key))
}

// GetStringKR ...
func (c *ApiBaseController) GetStringKR(key string, def ...string) string {
	return utils.ConvertEucKR(c.GetString(key))
}
