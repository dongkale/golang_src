package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	ora "gopkg.in/rana/ora.v4"
)

type EntpVideoController struct {
	BaseController
}

func (c *EntpVideoController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/login")
	}

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := mem_no
	cdnPath, _ := beego.AppConfig.String("cdnpath")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Entp Video Info
	log.Debug(fmt.Sprintf("CALL ZSP_ENTP_VD_INFO_R('%v', '%v', :1)",
		pLang, pEntpMemNo))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_ENTP_VD_INFO_R('%v', '%v', :1)",
		pLang, pEntpMemNo),
		ora.S,   /* VD_FILE_PATH1 */
		ora.S,   /* VD_FILE_PATH1 */
		ora.S,   /* VD_FILE_PATH1 */
		ora.S,   /* VD_FILE_PATH1 */
		ora.I64, /* VD_CNT */
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

	entpVideo := make([]models.EntpVideo, 0)

	var (
		vdFilePath1     string
		vdFilePath2     string
		vdFilePath3     string
		vdFilePath4     string
		vdCnt           int64
		vdFullFilePath1 string
		vdFullFilePath2 string
		vdFullFilePath3 string
		vdFullFilePath4 string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			vdFilePath1 = procRset.Row[0].(string)
			if vdFilePath1 == "" {
				vdFullFilePath1 = vdFilePath1
			} else {
				vdFullFilePath1 = cdnPath + vdFilePath1
			}
			vdFilePath2 = procRset.Row[1].(string)
			if vdFilePath2 == "" {
				vdFullFilePath2 = vdFilePath2
			} else {
				vdFullFilePath2 = cdnPath + vdFilePath2
			}
			vdFilePath3 = procRset.Row[2].(string)
			if vdFilePath3 == "" {
				vdFullFilePath3 = vdFilePath3
			} else {
				vdFullFilePath3 = cdnPath + vdFilePath3
			}
			vdFilePath4 = procRset.Row[3].(string)
			if vdFilePath4 == "" {
				vdFullFilePath4 = vdFilePath4
			} else {
				vdFullFilePath4 = cdnPath + vdFilePath4
			}
			vdCnt = procRset.Row[4].(int64)

			entpVideo = append(entpVideo, models.EntpVideo{
				VdFilePath1: vdFullFilePath1,
				VdFilePath2: vdFullFilePath2,
				VdFilePath3: vdFullFilePath3,
				VdFilePath4: vdFullFilePath4,
				VdCnt:       vdCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Entp Video Info

	c.Data["VdFilePath1"] = vdFullFilePath1
	c.Data["VdFilePath2"] = vdFullFilePath2
	c.Data["VdFilePath3"] = vdFullFilePath3
	c.Data["VdFilePath4"] = vdFullFilePath4
	c.Data["VdCnt"] = vdCnt

	c.Data["TMenuId"] = "E00"
	c.Data["SMenuId"] = "E02"

	c.TplName = "entp/entp_video.html"
}
