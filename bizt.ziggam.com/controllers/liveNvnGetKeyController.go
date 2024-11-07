package controllers

import (
	"fmt"
	"strings"

	"bizt.ziggam.com/models"
	"bizt.ziggam.com/utils"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type LiveNvnGetKeyController struct {
	BaseController
}

func (c *LiveNvnGetKeyController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	// 키값은 오면 그냥 입력해 주자.
	// mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	// if mem_no == nil {
	// 	c.Ctx.Redirect(302, "/login")
	// }

	pLang, _ := beego.AppConfig.String("lang")

	//pEntpMemSn := c.GetString("entp_mem_sn")
	pRoomId := c.GetString("room_id")
	pChannelId := c.GetString("channel_id")

	var (
		// lmPpMemNo    string // PP_MEM_NO
		// lmPpChrgSn   string //  PP_CHRG_SN
		lmNm         string // NM
		lmPpChrgBpNm string // PP_CHRG_BP_NM
		lmEntpAdmin  string // PP_CHRG_GBN_CD
		lmPtoPath    string
	)

	key := "room_id:" + pRoomId
	log.Info("remon key : ", key)

	// 레디스에서 리모트 몬스터 키값으로 회원 정보를 가져 온다.
	res, redis_err := utils.RPool.HGet(key, pChannelId)

	if redis_err != nil {
		fmt.Printf(fmt.Sprintf("error: %v", redis_err))
		c.TplName = "error/404.html"
		return
	}

	res_split := strings.Split(res, ":")

	fmt.Printf(fmt.Sprintf("entpGbn: %v ppMemNo: %v ppChrgSn: %v", res_split[0], res_split[1], res_split[2]))

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	fmt.Printf(fmt.Sprintf("CALL ZSP_LIVE_NVN_REMONKEY_INFO('%v', '%v', '%v', :1)",
		pLang, res_split[1], res_split[2]))
	stmtProcCallMem, errMem := ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_NVN_REMONKEY_INFO('%v', '%v', '%v', :1)",
		pLang, res_split[1], res_split[2]),
		// ora.S, /* ENTP_MEM_NO or PP_MEM_NO */
		// ora.S, /* PP_CHRG_SN */
		ora.S, /* NM */
		ora.S, /* PP_CHRG_BP_NM */
		ora.S, /* PP_CHRG_GBN_CD */
		ora.S, /* PTO_PATH */
	)
	defer stmtProcCallMem.Close()
	if errMem != nil {
		panic(errMem)
	}
	procRsetMem := &ora.Rset{}
	_, errMem = stmtProcCallMem.Exe(procRsetMem)

	if errMem != nil {
		panic(errMem)
	}

	imgServer, _  := beego.AppConfig.String("viewpath")
	if procRsetMem.IsOpen() {
		for procRsetMem.Next() {
			// lmPpMemNo = procRsetMem.Row[0].(string)    // PP_MEM_NO
			// lmPpChrgSn = procRsetMem.Row[1].(string)   //  PP_CHRG_SN
			lmNm = procRsetMem.Row[0].(string)         // NM
			lmPpChrgBpNm = procRsetMem.Row[1].(string) // PP_CHRG_BP_NM
			lmEntpAdmin = procRsetMem.Row[2].(string)  // PP_CHRG_GBN_CD
			lmPtoPath = procRsetMem.Row[3].(string)    // PP_CHRG_GBN_CD
			lmPtoPath = imgServer + lmPtoPath
		}
		if errMem := procRsetMem.Err(); errMem != nil {
			panic(errMem)
		}
	}

	rtnNvnChannelId := models.RtnNvnChannelId{}

	rtnNvnChannelId = models.RtnNvnChannelId{
		EntpGbn:   res_split[0],
		PpMemNo:   res_split[1],
		PpChrgSn:  res_split[2],
		ChnnelId:  pChannelId,
		Nm:        lmNm,
		BpNm:      lmPpChrgBpNm,
		EntpAdmin: lmEntpAdmin,
		PtoPath:   lmPtoPath,
	}

	c.Data["json"] = rtnNvnChannelId
	c.ServeJSON()
}
