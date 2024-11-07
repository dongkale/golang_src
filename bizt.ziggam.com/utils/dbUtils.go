package utils

import (
	"strings"
	"sync"

	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/errgo.v1"
	"gopkg.in/rana/ora.v4"
)

var (
	oraCxMu sync.Mutex
	oraInit sync.Once
)

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
