package controllers

import (
	"fmt"

	"emst.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type AdminEventVideoController struct {
	BaseController
}

func (c *AdminEventVideoController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/common/login")
	}

	pLang, _ := beego.AppConfig.String("lang")
	pPpMemNo := c.GetString("pp_mem_no")
	pKndCd := c.GetString("knd_cd")
	if pKndCd == "" {
		pKndCd = "A00"
	}
	pVdSn := c.GetString("vd_sn")
	if pVdSn == "" {
		pVdSn = "0"
	}

	var vdFullPath string
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

	// Start : Admin Event List
	log.Debug("CALL SP_EMS_ADMIN_EVENT_VIDEO_R('%v', '%v','%v', :1)",
		pLang, pPpMemNo, pKndCd)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_EVENT_VIDEO_R('%v', '%v', '%v', :1)",
		pLang, pPpMemNo, pKndCd),
		ora.I64, /* TOT_CNT */
		ora.I64, /* SN */
		ora.S,   /* REG_DT */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* VD_SN */
		ora.S,   /* VD_FILE_PATH */
		ora.S,   /* THM_KND_CD */
		ora.S,   /* QST_CD */
		ora.S,   /* THM_DESC */
		ora.I64, /* VD_SEC */
		ora.S,   /* OPN_SET_CD */
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

	adminEventVideo := make([]models.AdminEventVideo, 0)

	var (
		totCnt     int64
		sn         int64
		regDt      string
		ppMemNo    string
		vdSn       string
		vdFilePath string
		thmKndCd   string
		qstCd      string
		thmDesc    string
		vdSec      int64
		opnSetCd   string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			totCnt = procRset.Row[0].(int64)
			sn = procRset.Row[1].(int64)
			regDt = procRset.Row[2].(string)
			ppMemNo = procRset.Row[3].(string)
			vdSn = procRset.Row[4].(string)
			vdFilePath = procRset.Row[5].(string)
			thmKndCd = procRset.Row[6].(string)
			qstCd = procRset.Row[7].(string)
			thmDesc = procRset.Row[8].(string)
			vdSec = procRset.Row[9].(int64)
			opnSetCd = procRset.Row[10].(string)

			adminEventVideo = append(adminEventVideo, models.AdminEventVideo{
				TotCnt:     totCnt,
				Sn:         sn,
				RegDt:      regDt,
				PpMemNo:    ppMemNo,
				VdSn:       vdSn,
				VdFilePath: vdFilePath,
				ThmKndCd:   thmKndCd,
				QstCd:      qstCd,
				ThmDesc:    thmDesc,
				VdSec:      vdSec,
				OpnSetCd:   opnSetCd,
				CVdSn:      pVdSn,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}

	// Start : Admin Event View
	log.Debug("CALL SP_EMS_ADMIN_EVENT_VIEW_R('%v', '%v','%v', :1)",
		pLang, pPpMemNo, pVdSn)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_EVENT_VIEW_R('%v', '%v', '%v', :1)",
		pLang, pPpMemNo, pVdSn),
		ora.S,   /* REG_DT */
		ora.S,   /* PP_MEM_NO */
		ora.S,   /* VD_SN */
		ora.S,   /* VD_FILE_PATH */
		ora.S,   /* THM_KND_CD */
		ora.S,   /* QST_CD */
		ora.S,   /* THM_DESC */
		ora.I64, /* VD_SEC */
		ora.S,   /* OPN_SET_CD */
	)
	defer stmtProcCall.Close()
	if err != nil {
		panic(err)
	}
	procRset = &ora.Rset{}
	_, err = stmtProcCall.Exe(procRset)

	if err != nil {
		panic(err)
	}

	adminEventView := make([]models.AdminEventView, 0)

	var (
		vRegDt      string
		vPpMemNo    string
		vVdSn       string
		vVdFilePath string
		vThmKndCd   string
		vQstCd      string
		vThmDesc    string
		vVdSec      int64
		vOpnSetCd	string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			vRegDt = procRset.Row[0].(string)
			vPpMemNo = procRset.Row[1].(string)
			vVdSn = procRset.Row[2].(string)
			vVdFilePath = procRset.Row[3].(string)
			vThmKndCd = procRset.Row[4].(string)
			vQstCd = procRset.Row[5].(string)
			thmDesc = procRset.Row[6].(string)
			vdSec = procRset.Row[7].(int64)
			vOpnSetCd = procRset.Row[8].(string)

			if vVdFilePath == "" {
				vdFullPath = vVdFilePath
			} else {
				vdFullPath = cdnPath + vVdFilePath
			}

			adminEventView = append(adminEventView, models.AdminEventView{
				VRegDt:      vRegDt,
				VPpMemNo:    vPpMemNo,
				VVdSn:       vVdSn,
				VVdFilePath: vdFullPath,
				VThmKndCd:   vThmKndCd,
				VQstCd:      vQstCd,
				VThmDesc:    vThmDesc,
				VVdSec:      vVdSec,
				VOpnSetCd:   vOpnSetCd,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}

	c.Data["TotCnt"] = totCnt
	c.Data["PpMemNo"] = ppMemNo
	c.Data["KndCd"] = pKndCd
	c.Data["AdminEventVideo"] = adminEventVideo

	c.Data["VVdSn"] = pVdSn
	c.Data["VRegDt"] = vRegDt
	c.Data["VVdFilePath"] = vdFullPath
	c.Data["VThmKndCd"] = vThmKndCd
	c.Data["VQstCd"] = vQstCd
	c.Data["VThmDesc"] = vThmDesc
	c.Data["VVdSec"] = vVdSec
	c.Data["VOpnSetCd"] = vOpnSetCd

	c.TplName = "admin/event_video.html"
}
