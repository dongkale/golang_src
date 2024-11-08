package controllers

import (
	"fmt"

	"emst.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	"gopkg.in/rana/ora.v4"
)

type CommonEmailCertController struct {
	beego.Controller
}

func (c *CommonEmailCertController) Get() {
	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	//session := c.StartSession()
	pLang, _ := beego.AppConfig.String("lang")
	pMemId := c.GetString("mem_id")
	pCertKey := c.GetString("cert_key")

	//siteUrl := beego.AppConfig.String("siteurl")
	imgServer, _ := beego.AppConfig.String("viewpath")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Certification Key Info

	log.Debug("CALL SP_EMS_EMAIL_CERT_CHK_R('%v', '%v', '%v', :1)",
		pLang, pMemId, pCertKey)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_EMAIL_CERT_CHK_R('%v', '%v', '%v', :1)",
		pLang, pMemId, pCertKey),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
		ora.S,   /* ENTP_NM */
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

	emailCert := make([]models.EmailCert, 0)

	var (
		rtnCd  int64
		rtnMsg string
		entpNm string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
			entpNm = procRset.Row[2].(string)

			emailCert = append(emailCert, models.EmailCert{
				RtnCd:     rtnCd,
				RtnMsg:    rtnMsg,
				EntpNm:    entpNm,
				ImgServer: imgServer,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}

	c.Data["RtnCd"] = rtnCd
	c.Data["RtnMsg"] = rtnMsg
	c.Data["EntpNm"] = entpNm
	c.Data["ImgServer"] = imgServer
	c.TplName = "cert/cert_form.html"
}
