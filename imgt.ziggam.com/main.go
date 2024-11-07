package main

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	_ "imgt.ziggam.com/routers"
)

// var logger *logs.BeeLogger

func main() {	
	// logs.InitLogs()
	// logs.SetLogger(logs.AdapterFile, `{"filename":"logs/test.log"}`)
	// lotate 설정
	// logs.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/test.log","separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`)

	// ----

	// f := &logs.PatternLogFormatter{
	// 	Pattern:    "%F:%n|%w%t>> %m",
	// 	WhenFormat: "2006-01-02",
	// }
	// logs.RegisterFormatter("pattern", f)

	// _ = logs.SetGlobalFormatter("pattern")

	logs.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/test.log"}`)
	
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	
	beego.Run()
}


// func InitLog() error {
// 	logger = logs.NewLogger()
// 	//logger.EnableFuncCallDepth(true)
// 	makeLog(logger, "system.log")
// 	return nil
// }

// func makeLog(l *logs.BeeLogger, fileName string) error {
// 	//err := l.SetLogger(logs.AdapterMultiFile, fmt.Sprintf(
// 	//	`{"filename":"%s/log/%s", "daily":true, "maxdays":7, "rotate":true}`, GetCurrPath(), fileName))
// 	err := l.SetLogger(logs.AdapterMultiFile, fmt.Sprintf(
// 		`{"filename":"%s/log/%s", "hourly":true, "maxhours":7, "rotate":true}`, GetCurrPath(), fileName))
// 	if err != nil {
// 		return errors.New("init log error:" + err.Error())
// 	}
// 	l.Async(10000)
// 	return nil
// }
