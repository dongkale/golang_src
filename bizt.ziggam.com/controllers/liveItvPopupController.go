package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	ora "gopkg.in/rana/ora.v4"
)

type LiveItvPopupController struct {
	BaseController
}

func (c *LiveItvPopupController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/login")
	}

	mem_sn := session.Get(c.Ctx.Request.Context(), "mem_sn")

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := mem_no
	pRecrutSn := c.GetString("recrut_sn")
	pPpMemNo := c.GetString("pp_mem_no")
	pPpMemNm := c.GetString("pp_mem_nm")
	pPtoPath := c.GetString("pto_path")
	pLiveSn := c.GetString("live_sn")

	// imgServer, _  := beego.AppConfig.String("viewpath")
	// cdnPath := beego.AppConfig.String("cdnpath")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : live member list
	var (
		lmPpChrgGbnCd string
		lmPpChrgNm    string
		lmPpChrgBpNm  string
		lmPpChrgSn    string
	)

	fmt.Printf(fmt.Sprintf("CALL ZSP_LIVE_MEM_LIST_V2_R('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo, pLiveSn))

	stmtProcCallMem, errMem := ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_MEM_LIST_V2_R('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pPpMemNo, pLiveSn),
		ora.S, /* PP_CHRG_GBN_CD */
		ora.S, /* PP_CHRG_NM */
		ora.S, /* PP_CHRG_BP_NM */
		ora.S, /* PP_CHRG_SN */
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

	liveMemList := make([]models.LiveMemList_v2, 0)

	if procRsetMem.IsOpen() {
		for procRsetMem.Next() {
			lmPpChrgGbnCd = procRsetMem.Row[0].(string)
			lmPpChrgNm = procRsetMem.Row[1].(string)
			lmPpChrgBpNm = procRsetMem.Row[2].(string)
			lmPpChrgSn = procRsetMem.Row[3].(string)

			liveMemList = append(liveMemList, models.LiveMemList_v2{
				LmPpChrgGbnCd: lmPpChrgGbnCd,
				LmPpChrgNm:    lmPpChrgNm,
				LmPpChrgBpNm:  lmPpChrgBpNm,
				LmPpChrgSn:    lmPpChrgSn,
			})

			//fmt.Printf("Member = " + lmPpChrgGbnCd + ", " + lmPpChrgNm + ", " + lmPpChrgBpNm)
		}
		if errMem := procRsetMem.Err(); errMem != nil {
			panic(errMem)
		}
	}
	// End : live member list

	//LiveChNo”:“E2018102500001_2019110644_P2019082300298_201912091308002
	//=> 지원하는 기업회원번호_공고번호_개인회원번호_라이브채팅방 번호

	//var isValid bool = true

	// liveInv := c.GetSession("LiveInv")
	// if liveInv == nil {

	// 	c.SetSession("LiveInv", pPpMemNo)
	// 	c.SetSession("LiveInvTime", time.Now())

	// } else {

	// 	if liveInv.(string) == pPpMemNo {
	// 		isValid = false
	// 	} else {
	// 		isValid = true
	// 	}
	// }

	// c.SetSession("XYZ", "WWW")

	//if isValid == true {
	//fmt.Printf("Session[LiveInv]: " + c.GetSession("LiveInv").(string))
	// fmt.Printf("pEntpMemNo: " + pEntpMemNo.(string))
	// fmt.Printf("SMemSn: " + mem_sn.(string))
	// fmt.Printf("pRecrutSn: " + pRecrutSn)
	// fmt.Printf("pPpMemNo: " + pPpMemNo)
	// fmt.Printf("pPpMemNm: " + pPpMemNm)
	// fmt.Printf("pPtoPath: " + pPtoPath)
	// fmt.Printf("pLiveSn: " + pLiveSn)

	//fmt.Printf("Session[XYZ]: " + c.GetSession("XYZ").(string))

	for _, value := range liveMemList {
		fmt.Printf("Member: " + value.LmPpChrgGbnCd + ", " + value.LmPpChrgNm + ", " + value.LmPpChrgBpNm + ", " + value.LmPpChrgSn)
	}

	//c.Data["IsValid"] = isValid
	c.Data["EntpMemNo"] = pEntpMemNo
	c.Data["RecrutSn"] = pRecrutSn
	c.Data["PpMemNo"] = pPpMemNo
	c.Data["PpMemNm"] = pPpMemNm
	c.Data["PtoPath"] = pPtoPath
	c.Data["pLiveSn"] = pLiveSn
	c.Data["SMemSn"] = mem_sn // 맴버 PP_CHRG_SN
	c.Data["LiveMemList"] = liveMemList

	c.TplName = "live/live_itv_popup_new.html"
}
