package main

import (
	"fmt"

	_ "emst.ziggam.com/routers"

	"emst.ziggam.com/controllers"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/session"
	"github.com/dustin/go-humanize"
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
	beego.BConfig.WebConfig.Session.SessionOn = true
	sessionconf := &session.ManagerConfig{
		CookieName: "begoosessionID",
		Gclifetime: 3600,
		CookieLifeTime: 3600,
	}

	beego.GlobalSessions, _ = session.NewManager("memory", sessionconf)
	go beego.GlobalSessions.GC()

	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "begoosessionID"
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime = 3600
	beego.BConfig.WebConfig.Session.SessionProvider = "memory"
	beego.BConfig.WebConfig.Session.SessionProviderConfig = ""
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime = 3600	


	beego.ErrorController(&controllers.ErrorController{})
	beego.AddFuncMap("IntFmt", fnIntFmt)
	beego.Run()
}
