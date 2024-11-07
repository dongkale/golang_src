package controllers

import (
	"encoding/base64"
	"fmt"

	"bizt.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"

	// "github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/logs"
	ora "gopkg.in/rana/ora.v4"
)

type LiveNvnPopupController struct {
	BaseController
}

func (c *LiveNvnPopupController) Get() {
	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	pPpMemGbn := c.GetString("my_mem_no")

	var (
		lmPpMemNo     string
		lmPpChrgGbnCd string
		lmPpChrgNm    string
		lmPpChrgBpNm  string
		lmPpChrgSn    string
		lmEntpGbn     string
		lmLiveStatCd  string
		lmPtoPath     string
		pEntpGbn      string
		memSn         string
		entpAdmin     string
		pEntpMemNo    string
		pPpMemNm      string
		pPpMemNo      string
		pPtoPath      string
		tempValue     string
		pLang         string
		pRecrutSn     string
		pLiveSn       string
		isPPLink      bool
	)

	imgServer, _  := beego.AppConfig.String("viewpath")

	tempValue = c.GetString("pp")
	isPPLink = false
	if tempValue != "" {
		// tempValue = DecodeB64(tempValue)
		fmt.Printf("" + tempValue)

		pEntpMemNo, pRecrutSn, pLiveSn, pPpMemGbn = getPPLinkData(pLang, tempValue)

		if pPpMemGbn == "" {
			fmt.Printf("not found itv link : " + tempValue)
			c.TplName = "error/400.html"
			return
		}

			pEntpGbn = "P"

		isPPLink = true

	} else {
		pEntpGbn = "E"

		if pPpMemGbn == "" {
			fmt.Printf("not found itv link : " + tempValue)
			c.TplName = "error/400.html"
			return
		}
		
		if pPpMemGbn[0] == 'P' { // 개인 회원	<-- 여기로 들어오면 상용 서버에서는 삭제 해야 된다.
			pEntpGbn = "P"
			pEntpMemNo = c.GetString("entp_mem_no")
		} else { // 기업 회원
			mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
			if mem_no == nil {
				c.Ctx.Redirect(302, "/login")
			}

			memSn = session.Get(c.Ctx.Request.Context(), "mem_sn").(string)
			pEntpMemNo = mem_no.(string)
		}

		pLang, _ = beego.AppConfig.String("lang")

		pRecrutSn = c.GetString("recrut_sn")
		// pPpMemNo := c.GetString("pp_mem_no")
		// pPpMemNm := c.GetString("pp_mem_nm")
		// pPtoPath := c.GetString("pto_path")
		pLiveSn = c.GetString("live_sn")
	}

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

	fmt.Printf(fmt.Sprintf("CALL ZSP_LIVE_NVN_MEM_LIST_V2_R('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pLiveSn))

	stmtProcCallMem, errMem := ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_NVN_MEM_LIST_V2_R('%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn, pLiveSn),
		ora.S, /* PP_MEM_NO */
		ora.S, /* PP_CHRG_GBN_CD */
		ora.S, /* PP_CHRG_NM */
		ora.S, /* PP_CHRG_BP_NM */
		ora.S, /* PP_CHRG_SN */
		ora.S, /* ENTP_GBN */
		ora.S, /* LIVE_STAT_CD */
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

	liveMemList := make([]models.LiveMemList_v3, 0)

	if procRsetMem.IsOpen() {
		for procRsetMem.Next() {
			lmPpMemNo = procRsetMem.Row[0].(string)
			lmPpChrgGbnCd = procRsetMem.Row[1].(string)
			lmPpChrgNm = procRsetMem.Row[2].(string)
			lmPpChrgBpNm = procRsetMem.Row[3].(string)
			lmPpChrgSn = procRsetMem.Row[4].(string)
			lmEntpGbn = procRsetMem.Row[5].(string)
			lmLiveStatCd = procRsetMem.Row[6].(string)
			lmPtoPath = procRsetMem.Row[7].(string)
			if lmPtoPath == "" {
				lmPtoPath = ""
			} else {
				lmPtoPath = imgServer + lmPtoPath
			}
			// 들어온 유저 정보 검색후 세팅
			if pEntpGbn == "E" && lmPpChrgSn == pPpMemGbn { // 개인 회원
				pPpMemNm = lmPpChrgNm
				memSn = lmPpChrgSn
				entpAdmin = lmPpChrgGbnCd
				pPpMemNo = lmPpMemNo
			} else if pEntpGbn == "P" && lmPpMemNo == pPpMemGbn {
				pPpMemNm = lmPpChrgNm
				memSn = ""
				entpAdmin = "00"
				pPpMemNo = lmPpMemNo
				pPtoPath = lmPtoPath

				if lmLiveStatCd != "01" && lmLiveStatCd != "04" {
					// 요청중,수락 확정. 만 들어 올수 있다.
					fmt.Printf("LiveStatCd : " + lmLiveStatCd)
					c.TplName = "error/400.html"
					return
				}
			}

			liveMemList = append(liveMemList, models.LiveMemList_v3{
				LmPpMemNo:     lmPpMemNo,
				LmPpChrgGbnCd: lmPpChrgGbnCd,
				LmPpChrgNm:    lmPpChrgNm,
				LmPpChrgBpNm:  lmPpChrgBpNm,
				LmPpChrgSn:    lmPpChrgSn,
				LmEntpGbn:     lmEntpGbn,
				LmLiveStatCd:  lmLiveStatCd,
				LmPtoPath:     lmPtoPath,
			})

			//fmt.Printf("Member = " + lmPpChrgGbnCd + ", " + lmPpChrgNm + ", " + lmPpChrgBpNm)
		}
		if errMem := procRsetMem.Err(); errMem != nil {
			panic(errMem)
		}
	}
	// End : live member list

	if pPpMemNm == "" { // 자기 자신을 못찾은 경우.
		fmt.Printf("자기 자신을 못찾은 경우. my_mem_no: " + pPpMemGbn)
		c.TplName = "error/404.html"
		return
	}

	// 면접자가 링크로 들어 온 경우는 넘어간다.
	if isPPLink == false {
		// Start : live nvn join count
		fmt.Printf(fmt.Sprintf("CALL ZSP_LIVE_NVN_UPT_JOIN_CNT('%v', '%v', '%v', '%v', %v, :1)",
			pLang, pEntpMemNo, pRecrutSn, pLiveSn, 1))
		stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_NVN_UPT_JOIN_CNT('%v', '%v', '%v', '%v', %v, :1)",
			pLang, pEntpMemNo, pRecrutSn, pLiveSn, 1),
			ora.I64, /* RTN_CD */
			ora.S,   /* RTN_MSG */
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
	
		var (
			rtnCd  int64
			rtnMsg string
		)
	
		if procRset.IsOpen() {
			for procRset.Next() {
				rtnCd = procRset.Row[0].(int64)
				rtnMsg = procRset.Row[1].(string)
			}
	
			if err := procRset.Err(); err != nil {
				panic(err)
			}
	
			fmt.Printf(fmt.Sprintf(" ===> rtnCd:%v, rtnMsg:%v", rtnCd, rtnMsg))
		}
		// End : live nvn join count
	}

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
	c.Data["SMemSn"] = memSn        // 맴버 PP_CHRG_SN
	c.Data["EntpAdmin"] = entpAdmin // 맴버 PP_CHRG_SN
	c.Data["EntpGbn"] = pEntpGbn    // 기업회원인지 개인 회원인지 구분
	c.Data["LiveMemList"] = liveMemList

	c.TplName = "live/live_conference_popup.html"
}

func EncodeB64(message string) (retour string) {
	base64Text := make([]byte, base64.StdEncoding.EncodedLen(len(message)))
	base64.StdEncoding.Encode(base64Text, []byte(message))
	return string(base64Text)
}

func DecodeB64(message string) (retour string) {
	base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(message)))
	base64.StdEncoding.Decode(base64Text, []byte(message))
	fmt.Printf("base64: %s\n", base64Text)
	return string(base64Text)
}

func getPPLinkData(pLang, itv_link string) (entpMemNo string, recrutSn string, liveSn string, ppMemNo string) {
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : live member list

	fmt.Printf(fmt.Sprintf("CALL ZSP_LIVE_NVN_ITV_LINK('%v', '%v', :1)",
		pLang, itv_link))

	stmtProcCallMem, errMem := ses.Prep(fmt.Sprintf("CALL ZSP_LIVE_NVN_ITV_LINK('%v', '%v', :1)",
		pLang, itv_link),
		ora.S, /* ENTP_MEM_NO */
		ora.S, /* RECRUT_SN */
		ora.S, /* LIVE_SN */
		ora.S, /* PP_MEM_NO */
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

	if procRsetMem.IsOpen() {
		for procRsetMem.Next() {
			entpMemNo = procRsetMem.Row[0].(string)
			recrutSn = procRsetMem.Row[1].(string)
			liveSn = procRsetMem.Row[2].(string)
			ppMemNo = procRsetMem.Row[3].(string)

			return entpMemNo, recrutSn, liveSn, ppMemNo
		}
	}

	return
}
