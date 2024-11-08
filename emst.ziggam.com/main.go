package main

import (
	"fmt"
	_ "emst.ziggam.com/routers"
	
	beego "github.com/beego/beego/v2/server/web"	
	"github.com/beego/beego/v2/server/web/session"
	"github.com/dustin/go-humanize"
	"emst.ziggam.com/controllers"
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
}
	beego.GlobalSessions, _ = session.NewManager("memory", sessionconf)
	go beego.GlobalSessions.GC()

	beego.ErrorController(&controllers.ErrorController{})
	beego.AddFuncMap("IntFmt", fnIntFmt)
	beego.Run()
}
