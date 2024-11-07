package main

import (
	"context"
	"runtime"
	"time"

	"fmt"
	"github.com/dustin/go-humanize"
	_ "bizt.ziggam.com/routers"
	"bizt.ziggam.com/utils"
	beego "github.com/beego/beego/v2/server/web"

	//"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/session"
)

func fnIntFmt(in int64) (out string) {
	if in > 0 {
		out = fmt.Sprintf("%s", humanize.Comma(in))
	} else {
		out = "0"
	}
	return
}

func fnTEST(in string) (out string) {
	if in != "" {
		out = "TEST"
	} else {
		out = "TEST0"
	}
	return
}

func fnTEST2(in string) (out string) {
	if in != "" {
		out = "TEST2"
	} else {
		out = "TEST20"
	}
	return
}

func fnRecoverCB(ctx *context.Context) {

	return
}

// func dbError(rw http.ResponseWriter, r *http.Request) {
// 	t, _ := template.ParseFiles(beego.BConfig.WebConfig.ViewsPath + "/error/error.html")
// 	data := make(map[string]interface{})
// 	data["content"] = "database is now down"
// 	t.Execute(rw, data)
// }

func main() {

	//defer beego.re.BConfig.RecoverFunc(fnRecoverCB)

	sessionconf := &session.ManagerConfig{
		CookieName: "begoosessionID",
		Gclifetime: 3600,
	}
	beego.GlobalSessions, _ = session.NewManager("memory", sessionconf)
	go beego.GlobalSessions.GC()

	//beego.ErrorHandler("dbError", dbError)

	runmode, _ := beego.AppConfig.String("runmode")
	appname, _ := beego.AppConfig.String("appname")
	serverStartSlack, _ := beego.AppConfig.String("serverStartSlack")
	serverStartSlackMsg, _ := beego.AppConfig.String("serverStartSlackMsg")
	

	if serverStartSlack == "true" {
		utils.SlackSend(runmode, fmt.Sprintf("[%v] %v", appname, serverStartSlackMsg), "", "")
		//utils.SlackSend("live", fmt.Sprintf("[%v] %v", beego.AppConfig.String("appname"), beego.AppConfig.String("serverStartSlackMsg")), "", "")
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf(fmt.Sprintf("GoMaxProc:%v", runtime.GOMAXPROCS(0)))

	// pool cnt:4, channel cnt: 10, send tyy cnt: 2
	//utils.SendMailPoolMng.Init(4, 10, 2)

	//utils.SendMailDaemonMng.Init()

	// loop time:10, select cnt: 20
	//utils.MailDBPoolMng.Init(utils.MailDBPool_Mail_Pool, 10, 20)

	// beego.ErrorController(&controllers.ErrorController{})
	beego.AddFuncMap("IntFmt", fnIntFmt)
	beego.AddFuncMap("TEST", fnTEST)
	// beego.AddFuncMap("TEST2", fnTEST2)
	beego.Run()
}

func loop() {
	// for {
	// 	select {
	// 	// case <-stop:
	// 	// 	fmt.Println("EXIT: 3 seconds")
	// 	// 	return
	// 	case <-time.After(5 * time.Second):
	// 		fmt.Printf("5 second")
	// 		break
	// 	case <-time.After(1 * time.Second):
	// 		fmt.Printf("1 second")
	// 		break
	// 	}
	// }

	for now := range time.Tick(5 * time.Second) {
		fmt.Printf(now.String())
	}
}