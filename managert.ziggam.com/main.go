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
	sessionconf := &session.ManagerConfig{
		CookieName: "begoosessionID",
		Gclifetime: 3600,
		CookieLifeTime: 10,
	}
	beego.GlobalSessions, _ = session.NewManager("memory", sessionconf)
	go beego.GlobalSessions.GC()

	logs.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/managert.ziggam.com.log"}`)
	
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)

	beego.ErrorController(&controllers.ErrorController{})
	beego.AddFuncMap("IntFmt", fnIntFmt)
	beego.Run()
}

