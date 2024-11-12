package controllers

import (
	"fmt"

	"emst.ziggam.com/models"
	"emst.ziggam.com/utils"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type AdminBannerWriteController struct {
	BaseController
}

func (c *AdminBannerWriteController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/common/login")
	}

	imgServer, _ := beego.AppConfig.String("viewpath")

	pLang, _ := beego.AppConfig.String("lang")
	pBnrSn := c.GetString("bnr_sn")
	if pBnrSn == "" {
		pBnrSn = "0"
	}

	pPageNo := c.GetString("pn")
	pCuCd := c.GetString("cu_cd")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Banner Detail
	log.Debug("CALL SP_EMS_ADMIN_BANNER_DTL_R('%v', '%v', :1)",
		pLang, pBnrSn)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_BANNER_DTL_R('%v', '%v', :1)",
		pLang, pBnrSn),
		ora.S,   /* BNR_TITLE */
		ora.S,   /* PUBL_SDY */
		ora.S,   /* PUBL_EDY */
		ora.S,   /* LNK_GBN_CD */
		ora.S,   /* LNK_GBN_VAL */
		ora.S,   /* BRD_GBN_CD */
		ora.I64, /* SN */
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* DEL_YN */
		ora.S,   /* USE_YN */
		ora.S,   /* PTO_PATH */
		ora.S,   /* THUMB_PTO_PATH */
		ora.S,   /* BNR_KND_CD */
		ora.S,   /* LIST_BNR_SN */
		ora.S,   /* LIST_TITLE */
		ora.S,   /* LIST_PTOTO_PATH */
		ora.S,   /* LIST_THUMB_PTOTO_PATH */
		ora.S,   /* LIST_LINK_URL */
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

	adminBannerWrite := make([]models.AdminBannerWrite, 0)

	var (
		bnrTitle     			string
		sdy          			string
		edy          			string
		lnkGbnCd     			string
		lnkGbnVal    			string
		brdGbnCd     			string
		sn           			int64
		entpMemNo    			string
		recrutSn     			string
		delYn        			string
		useYn        			string
		ptoPath      			string
		thumbPtoPath 			string
		fullPtoPath  			string
		bnrKndCd     			string
		listBnrSn    			string
		listTitle    			string
		listPhotoPath			string
		listThumbPtotoPath		string
		listFullPhotoPath		string
		listLinkUrl				string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			bnrTitle = procRset.Row[0].(string)
			sdy = procRset.Row[1].(string)
			edy = procRset.Row[2].(string)
			lnkGbnCd = procRset.Row[3].(string)
			lnkGbnVal = procRset.Row[4].(string)
			brdGbnCd = procRset.Row[5].(string)
			// sn = procRset.Row[6].(int64)
			sn = utils.DbRowToInt64(procRset.Row[6], 0)
			entpMemNo = procRset.Row[7].(string)
			recrutSn = procRset.Row[8].(string)
			delYn = procRset.Row[9].(string)
			useYn = procRset.Row[10].(string)
			ptoPath = procRset.Row[11].(string)
			thumbPtoPath = procRset.Row[12].(string)
			bnrKndCd = procRset.Row[13].(string)
			listBnrSn = procRset.Row[14].(string)
			listTitle = procRset.Row[15].(string)
			listPhotoPath = procRset.Row[16].(string)
			listThumbPtotoPath = procRset.Row[17].(string)
			listLinkUrl = procRset.Row[18].(string)

			log.Debug("CALL LOG('%v', '%v', '%v', '%v', '%v', '%v', :1)",
				bnrTitle, listTitle, listPhotoPath, listThumbPtotoPath, listLinkUrl, listBnrSn)

			if ptoPath == "" {
				fullPtoPath = ptoPath
			} else {
				fullPtoPath = imgServer + ptoPath
			}

			if listPhotoPath == "" {
				listFullPhotoPath = listPhotoPath
			} else {
				listFullPhotoPath = imgServer + listPhotoPath
			}

			adminBannerWrite = append(adminBannerWrite, models.AdminBannerWrite{
				BnrTitle:     		bnrTitle,
				Sdy:          		sdy,
				Edy:          		edy,
				LnkGbnCd:     		lnkGbnCd,
				LnkGbnVal:    		lnkGbnVal,
				BrdGbnCd:     		brdGbnCd,
				Sn:           		sn,
				EntpMemNo:    		entpMemNo,
				RecrutSn:     		recrutSn,
				DelYn:        		delYn,
				UseYn:        		useYn,
				PtoPath:      		fullPtoPath,
				ThumbPtoPath: 		thumbPtoPath,
				OriImgFile:   		ptoPath,
				BnrKndCd:     		bnrKndCd,
				ListBnrSn:			listBnrSn,
				ListTitle:	  		listTitle,
				ListPhotoPath:		listFullPhotoPath,
				ListThumbPhotoPath:	listThumbPtotoPath,
				ListLinkUrl:		listLinkUrl,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Banner Detail

	// Start : Common Code List - 링크구분
	pCdGrpId := "G023"
	log.Debug("CALL SP_EMS_COMMON_CD_LIST_R('%v', '%v', :1)",
		pLang, pCdGrpId)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_COMMON_CD_LIST_R('%v', '%v', :1)",
		pLang, pCdGrpId),
		ora.S, /* CD_ID */
		ora.S, /* CD_NM */
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

	commonLnkGbnCd := make([]models.CommonLnkGbnCd, 0)

	var (
		lgCdId string
		lgCdNm string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			lgCdId = procRset.Row[0].(string)
			lgCdNm = procRset.Row[1].(string)

			commonLnkGbnCd = append(commonLnkGbnCd, models.CommonLnkGbnCd{
				LgCdId: lgCdId,
				LgCdNm: lgCdNm,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End :  Common Code List

	// Start : Common Code List - 링크메뉴값
	pCdGrpId = "G024"
	log.Debug("CALL SP_EMS_COMMON_CD_LIST_R('%v', '%v', :1)",
		pLang, pCdGrpId)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_COMMON_CD_LIST_R('%v', '%v', :1)",
		pLang, pCdGrpId),
		ora.S, /* CD_ID */
		ora.S, /* CD_NM */
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

	commonLnkMenuCd := make([]models.CommonLnkMenuCd, 0)

	var (
		lvCdId string
		lvCdNm string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			lvCdId = procRset.Row[0].(string)
			lvCdNm = procRset.Row[1].(string)

			commonLnkMenuCd = append(commonLnkMenuCd, models.CommonLnkMenuCd{
				LvCdId: lvCdId,
				LvCdNm: lvCdNm,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End :  Common Code List

	// Start : Common Code List - 배너종류값
	pCdGrpId = "G026"
	log.Debug("CALL SP_EMS_COMMON_CD_LIST_R('%v', '%v', :1)",
		pLang, pCdGrpId)

	stmtProcCall, err = ses.Prep(fmt.Sprintf("CALL SP_EMS_COMMON_CD_LIST_R('%v', '%v', :1)",
		pLang, pCdGrpId),
		ora.S, /* CD_ID */
		ora.S, /* CD_NM */
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

	commonBnrKndCd := make([]models.CommonBnrKndCd, 0)

	var (
		bkCdId string
		bkCdNm string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			bkCdId = procRset.Row[0].(string)
			bkCdNm = procRset.Row[1].(string)

			commonBnrKndCd = append(commonBnrKndCd, models.CommonBnrKndCd{
				BkCdId: bkCdId,
				BkCdNm: bkCdNm,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End :  Common Code List

	if listBnrSn == "" {
		listBnrSn = "0"
	}

	c.Data["CuCd"] = pCuCd
	c.Data["CommonLnkGbnCd"] = commonLnkGbnCd
	c.Data["CommonLnkMenuCd"] = commonLnkMenuCd
	c.Data["CommonBnrKndCd"] = commonBnrKndCd
	c.Data["BnrTitle"] = bnrTitle
	c.Data["Sdy"] = sdy
	c.Data["Edy"] = edy
	c.Data["LnkGbnCd"] = lnkGbnCd
	c.Data["LnkGbnVal"] = lnkGbnVal
	c.Data["BrdGbnCd"] = brdGbnCd
	c.Data["Sn"] = sn
	c.Data["EntpMemNo"] = entpMemNo
	c.Data["RecrutSn"] = recrutSn
	c.Data["DelYn"] = delYn
	c.Data["UseYn"] = useYn
	c.Data["PtoPath"] = fullPtoPath
	c.Data["OriImgFile"] = ptoPath
	c.Data["OriThumbImgFile"] = thumbPtoPath
	c.Data["BnrSn"] = pBnrSn
	c.Data["BnrKndCd"] = bnrKndCd
	c.Data["ListBnrSn"] = listBnrSn
	c.Data["ListTitle"] = listTitle
	c.Data["ListPhotoPath"] = listFullPhotoPath
	c.Data["ListOriImgFile"] = listPhotoPath
	c.Data["ListOriThumbImgFile"] = listThumbPtotoPath
	c.Data["ListLinkUrl"] = listLinkUrl

	c.Data["PageNo"] = pPageNo
	c.Data["MenuId"] = "04"
	c.Data["SubMenuId"] = "07"
	c.TplName = "admin/banner_write.html"
}
