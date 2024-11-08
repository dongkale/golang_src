package controllers

import (
	"fmt"

	"emst.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	ora "gopkg.in/rana/ora.v4"
)

type AdminStatsMainController struct {
	BaseController
}

func (c *AdminStatsMainController) Get() {

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
	pEntpMemNo := c.GetString("entp_mem_no")
	if pEntpMemNo == "" {
		pEntpMemNo = "T10"
	}

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Admin Stats Main
	log.Debug("CALL SP_EMS_ADMIN_STATS_MAIN_R('%v', :1)",
		pLang)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_STATS_MAIN_R('%v', :1)",
		pLang),
		ora.I64, /* TOT_MEM_CNT */
		ora.I64, /* TD_MEM_CNT */
		ora.I64, /* WD_MEM_CNT */
		ora.I64, /* TOT_ENTP_CNT */
		ora.I64, /* TD_ENTP_CNT */
		ora.I64, /* WD_ENTP_CNT */
		ora.I64, /* VP_MEM_CNT */
		ora.I64, /* TOT_VP_CNT */
		ora.I64, /* F_VP_CNT */
		ora.I64, /* P_VP_CNT */
		ora.I64, /* C_VP_CNT */
		ora.I64, /* AD_MEM_CNT */
		ora.I64, /* IS_MEM_CNT */
		ora.I64, /* M_MEM_CNT */
		ora.I64, /* F_MEM_CNT */
		ora.I64, /* TOT_EVP_CNT */
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

	adminStatsMain := make([]models.AdminStatsMain, 0)

	var (
		totMemCnt  int64
		tdMemCnt   int64
		wdMemCnt   int64
		totEntpCnt int64
		tdEntpCnt  int64
		wdEntpCnt  int64
		vpMemCnt   int64
		totVpCnt   int64
		fVpCnt     int64
		pVpCnt     int64
		cVpCnt     int64
		adMemCnt   int64
		isMemCnt   int64
		mMemCnt    int64
		fMemCnt    int64
		totEvpCnt  int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			totMemCnt = procRset.Row[0].(int64)
			tdMemCnt = procRset.Row[1].(int64)
			wdMemCnt = procRset.Row[2].(int64)
			totEntpCnt = procRset.Row[3].(int64)
			tdEntpCnt = procRset.Row[4].(int64)
			wdEntpCnt = procRset.Row[5].(int64)
			vpMemCnt = procRset.Row[6].(int64)
			totVpCnt = procRset.Row[7].(int64)
			fVpCnt = procRset.Row[8].(int64)
			pVpCnt = procRset.Row[9].(int64)
			cVpCnt = procRset.Row[10].(int64)
			adMemCnt = procRset.Row[11].(int64)
			isMemCnt = procRset.Row[12].(int64)
			mMemCnt = procRset.Row[13].(int64)
			fMemCnt = procRset.Row[14].(int64)
			totEvpCnt = procRset.Row[15].(int64)

			adminStatsMain = append(adminStatsMain, models.AdminStatsMain{
				TotMemCnt:  totMemCnt,
				TdMemCnt:   tdMemCnt,
				WdMemCnt:   wdMemCnt,
				TotEntpCnt: totEntpCnt,
				TdEntpCnt:  tdEntpCnt,
				WdEntpCnt:  wdEntpCnt,
				VpMemCnt:   vpMemCnt,
				TotVpCnt:   totVpCnt,
				FVpCnt:     fVpCnt,
				PVpCnt:     pVpCnt,
				CVpCnt:     cVpCnt,
				AdMemCnt:   adMemCnt,
				IsMemCnt:   isMemCnt,
				MMemCnt:    mMemCnt,
				FMemCnt:    fMemCnt,
				TotEvpCnt:  totEvpCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Admin Stats Main

	// Start : Admin Stats Sub1
	log.Debug("CALL SP_EMS_ADMIN_STATS_SUB1_R('%v', :1)",
		pLang)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_STATS_SUB1_R('%v', :1)",
		pLang),
		ora.I64, /* P10U */
		ora.I64, /* P16_20 */
		ora.I64, /* P21_25 */
		ora.I64, /* P26_30 */
		ora.I64, /* P31_35 */
		ora.I64, /* P36_40 */
		ora.I64, /* P41_45 */
		ora.I64, /* P46_50 */
		ora.I64, /* P51_55 */
		ora.I64, /* P56_60 */
		ora.I64, /* P61_65 */
		ora.I64, /* P66_70 */
		ora.I64, /* P70O */
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

	adminStatsSub1 := make([]models.AdminStatsSub1, 0)

	var (
		p10U   int64
		p16_20 int64
		p21_25 int64
		p26_30 int64
		p31_35 int64
		p36_40 int64
		p41_45 int64
		p46_50 int64
		p51_55 int64
		p56_60 int64
		p61_65 int64
		p66_70 int64
		p70O   int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			p10U = procRset.Row[0].(int64)
			p16_20 = procRset.Row[1].(int64)
			p21_25 = procRset.Row[2].(int64)
			p26_30 = procRset.Row[3].(int64)
			p31_35 = procRset.Row[4].(int64)
			p36_40 = procRset.Row[5].(int64)
			p41_45 = procRset.Row[6].(int64)
			p46_50 = procRset.Row[7].(int64)
			p51_55 = procRset.Row[8].(int64)
			p56_60 = procRset.Row[9].(int64)
			p61_65 = procRset.Row[10].(int64)
			p66_70 = procRset.Row[11].(int64)
			p70O = procRset.Row[12].(int64)

			adminStatsSub1 = append(adminStatsSub1, models.AdminStatsSub1{
				P10U:   p10U,
				P16_20: p16_20,
				P21_25: p21_25,
				P26_30: p26_30,
				P31_35: p31_35,
				P36_40: p36_40,
				P41_45: p41_45,
				P46_50: p46_50,
				P51_55: p51_55,
				P56_60: p56_60,
				P61_65: p61_65,
				P66_70: p66_70,
				P70O:   p70O,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Admin Stats Sub1

	// Start : Admin Stats Sub2
	log.Debug("CALL SP_EMS_ADMIN_STATS_SUB2_R('%v', '%v', :1)",
		pLang, pEntpMemNo)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_STATS_SUB2_R('%v', '%v', :1)",
		pLang, pEntpMemNo),
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* ENTP_KO_NM */
		ora.I64, /* PV_CNT */
		ora.I64, /* UV_CNT */
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

	adminStatsSub2 := make([]models.AdminStatsSub2, 0)

	var (
		entpMemNo string
		entpKoNm  string
		pvCnt     int64
		uvCnt     int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			entpMemNo = procRset.Row[0].(string)
			entpKoNm = procRset.Row[1].(string)
			pvCnt = procRset.Row[2].(int64)
			uvCnt = procRset.Row[3].(int64)

			adminStatsSub2 = append(adminStatsSub2, models.AdminStatsSub2{
				EntpMemNo: entpMemNo,
				EntpKoNm:  entpKoNm,
				PvCnt:     pvCnt,
				UvCnt:     uvCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Admin Stats Sub2

	// Start : Admin Stats Sub3
	log.Debug("CALL SP_EMS_ADMIN_STATS_SUB3_R('%v', :1)",
		pLang)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_STATS_SUB3_R('%v', :1)",
		pLang),
		ora.S, /* ENTP_MEM_NO */
		ora.S, /* ENTP_KO_NM */
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

	adminStatsSub3 := make([]models.AdminStatsSub3, 0)

	var (
		lstEntpMemNo string
		lstEntpKoNm  string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			lstEntpMemNo = procRset.Row[0].(string)
			lstEntpKoNm = procRset.Row[1].(string)

			adminStatsSub3 = append(adminStatsSub3, models.AdminStatsSub3{
				LstEntpMemNo: lstEntpMemNo,
				LstEntpKoNm:  lstEntpKoNm,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Admin Stats Sub2

	c.Data["AdminStatsSub2"] = adminStatsSub2
	c.Data["AdminStatsSub3"] = adminStatsSub3

	c.Data["P10U"] = p10U
	c.Data["P16_20"] = p16_20
	c.Data["P21_25"] = p21_25
	c.Data["P26_30"] = p26_30
	c.Data["P31_35"] = p31_35
	c.Data["P36_40"] = p36_40
	c.Data["P41_45"] = p41_45
	c.Data["P46_50"] = p46_50
	c.Data["P51_55"] = p51_55
	c.Data["P56_60"] = p56_60
	c.Data["P61_65"] = p61_65
	c.Data["P66_70"] = p66_70
	c.Data["P70O"] = p70O

	c.Data["TotMemCnt"] = totMemCnt
	c.Data["TdMemCnt"] = tdMemCnt
	c.Data["WdMemCnt"] = wdMemCnt
	c.Data["TotEntpCnt"] = totEntpCnt
	c.Data["TdEntpCnt"] = tdEntpCnt
	c.Data["WdEntpCnt"] = wdEntpCnt
	c.Data["VpMemCnt"] = vpMemCnt
	c.Data["TotVpCnt"] = totVpCnt
	c.Data["FVpCnt"] = fVpCnt
	c.Data["PVpCnt"] = pVpCnt
	c.Data["CVpCnt"] = cVpCnt
	c.Data["AdMemCnt"] = adMemCnt
	c.Data["IsMemCnt"] = isMemCnt
	c.Data["MMemCnt"] = mMemCnt
	c.Data["FMemCnt"] = fMemCnt
	c.Data["TotEvpCnt"] = totEvpCnt

	c.Data["MenuId"] = "99"

	c.TplName = "admin/stats_main.html"
}

func (c *AdminStatsMainController) Post() {

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
	pEntpMemNo := c.GetString("entp_mem_no")
	if pEntpMemNo == "" {
		pEntpMemNo = "T10"
	}

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Admin Stats Sub2
	log.Debug("CALL SP_EMS_ADMIN_STATS_SUB2_R('%v', '%v', :1)",
		pLang, pEntpMemNo)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_STATS_SUB2_R('%v', '%v', :1)",
		pLang, pEntpMemNo),
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* ENTP_KO_NM */
		ora.I64, /* PV_CNT */
		ora.I64, /* UV_CNT */
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

	rtnAdminStatsSub2 := models.RtnAdminStatsSub2{}
	adminStatsSub2 := make([]models.AdminStatsSub2, 0)

	var (
		entpMemNo string
		entpKoNm  string
		pvCnt     int64
		uvCnt     int64
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			entpMemNo = procRset.Row[0].(string)
			entpKoNm = procRset.Row[1].(string)
			pvCnt = procRset.Row[2].(int64)
			uvCnt = procRset.Row[3].(int64)

			adminStatsSub2 = append(adminStatsSub2, models.AdminStatsSub2{
				EntpMemNo: entpMemNo,
				EntpKoNm:  entpKoNm,
				PvCnt:     pvCnt,
				UvCnt:     uvCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
		rtnAdminStatsSub2 = models.RtnAdminStatsSub2{
			RtnAdminStatsSub2Data: adminStatsSub2,
		}
		// End : Admin Stats Sub2

		c.Data["json"] = &rtnAdminStatsSub2
		c.ServeJSON()
	}
}
