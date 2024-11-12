package main

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/session"
	"github.com/dustin/go-humanize"
	"managert.ziggam.com/controllers"
	_ "managert.ziggam.com/routers"
)

func fnIntFmt(in int64) (out string) {
	if in > 0 {
		out = fmt.Sprintf("%s", humanize.Comma(in))
	} else {
		out = "0"
	}
	return
}

func main() {
	beegoInitialize()

	beego.Run()
}

// func beegoInitialize__() {
// 	// sessionconf := &session.ManagerConfig{		
// 	// 	CookieName: "begoosessionID",
// 	// 	Gclifetime: 3600,
// 	// 	CookieLifeTime: 10,		
// 	// }
// 	beego.BConfig.WebConfig.Session.SessionOn = true

// 	sessionconf := &session.ManagerConfig{				
// 		CookieName: "begoosessionID",
// 		Gclifetime: 3600,
// 		CookieLifeTime: 3600,
// 	}
	
// 	beego.GlobalSessions, _ = session.NewManager("memory", sessionconf)	
// 	go beego.GlobalSessions.GC()

// 	logs.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/managert.ziggam.com.log", "maxlines":0, "maxsize":0, "daily":true, "maxdays":10, "rotate":true, "level":7}`)
	
// 	logs.EnableFuncCallDepth(true)
// 	logs.SetLogFuncCallDepth(3)

// 	beego.BConfig.WebConfig.Session.SessionOn = true
// 	beego.BConfig.WebConfig.Session.SessionName = "begoosessionID"
// 	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600
// 	beego.BConfig.WebConfig.Session.SessionCookieLifeTime = 3600
// 	beego.BConfig.WebConfig.Session.SessionProvider = "memory"
// 	beego.BConfig.WebConfig.Session.SessionProviderConfig = ""
// 	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600
// 	beego.BConfig.WebConfig.Session.SessionCookieLifeTime = 3600	

// 	beego.ErrorController(&controllers.ErrorController{})
// 	beego.AddFuncMap("IntFmt", fnIntFmt)	
// }

func beegoInitialize() {
	beego.BConfig.WebConfig.Session.SessionOn = true

    var structdata session.ManagerConfig
    structdata.CookieName = "gosessionid"
    structdata.EnableSetCookie = true
    structdata.Gclifetime = 3600
    structdata.Maxlifetime = 3600
    structdata.Secure = false
    structdata.CookieLifeTime = 3600
    structdata.ProviderConfig = ""
    beego.GlobalSessions, _ = session.NewManager("memory", &structdata)
    go beego.GlobalSessions.GC()
	
	logs.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/managert.ziggam.com.log", "maxlines":0, "maxsize":0, "daily":true, "maxdays":10, "rotate":true, "level":7}`)
	
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)

	beego.ErrorController(&controllers.ErrorController{})
	beego.AddFuncMap("IntFmt", fnIntFmt)	
}