package controllers

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"time"

	"bizt.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/disintegration/imaging"
	"gopkg.in/rana/ora.v4"
)

type EntpInfoUpdateController struct {
	beego.Controller
}

func (c *EntpInfoUpdateController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	mem_id := session.Get(c.Ctx.Request.Context(), "mem_id")
	if mem_id == nil {
		c.Ctx.Redirect(302, "/login")
	}

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := c.GetString("entp_mem_no")
	pRepNm := c.GetString("rep_nm")
	pTelNo := c.GetString("tel_no")
	pSmsRecvYn := c.GetString("sms_recv_yn")
	pEmail := c.GetString("email")
	pEmailRecvYn := c.GetString("email_recv_yn")
	pZip := c.GetString("zip")
	pAddr := c.GetString("addr")
	pDtlAddr := c.GetString("dtl_addr")
	pRefAddr := c.GetString("ref_addr")
	pBizTpy := c.GetString("biz_tpy")
	pbizCond := c.GetString("biz_cond")
	pEntpHtag1 := c.GetString("entp_htag1")
	pEntpHtag2 := c.GetString("entp_htag2")
	pEntpHtag3 := c.GetString("entp_htag3")
	pEntpIntr := c.GetString("entp_intr")
	pHomePg := c.GetString("home_pg")
	pEmpCnt := c.GetString("emp_cnt")
	pEstDy := c.GetString("est_dy")
	pInfoEqYn := c.GetString("info_eq_yn")

	pBizRegYn := c.GetString("biz_reg_yn")
	pImgYn := c.GetString("img_yn")

	pEntpRegNoExt := c.GetString("entp_regno_ext")
	pEntpLogoExt := c.GetString("entp_logo_ext")
	oriBizFile := c.GetString("ori_biz_file")
	//oriBizFileExt := c.GetString("ori_biz_file_ext")
	oriLogoFile := c.GetString("ori_logo_file")
	oriLogoFileExt := c.GetString("ori_logo_file_ext")

	pBizTpyCd := c.GetString("biz_tpy_cd")             // 업종 코드
	pEntpProfile := c.GetString("entp_profile")        // 기업 소개
	pBizIntro := c.GetString("biz_intro")              // 사업 소개
	pEntpCapital := c.GetString("entp_capital")        // 자본금
	pEntpTotalSales := c.GetString("entp_total_sales") // 매출액
	pEntpTypeCd := c.GetString("entp_type_cd")         // 기업 형태 코드
	pEntpType := c.GetString("entp_type")              // 기업 형태(대기업,중견기업,중소기업,공사/공기업,외국계기업,기타)
	pLocation := c.GetString("location")               // 소재지

	log.Debug("pBizTpyCd: " + pBizTpyCd)
	log.Debug("pBizTpy: " + pBizTpy)
	log.Debug("pEntpProfile: " + pEntpProfile)
	log.Debug("pBizIntro: " + pBizIntro)
	log.Debug("pEntpCapital: " + pEntpCapital)
	log.Debug("pEntpTotalSales: " + pEntpTotalSales)
	log.Debug("pEntpTypeCd: " + pEntpTypeCd)
	log.Debug("pEntpType: " + pEntpType)
	log.Debug("pLocation: " + pLocation)

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Entp Info Update Process

	// fmt.Printf(fmt.Sprintf("CALL ZSP_ENTP_UPT_PROC('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', %v, '%v', :1)",
	// 	pLang, pEntpMemNo, pRepNm, pTelNo, pSmsRecvYn, pEmail, pEmailRecvYn, pInfoEqYn, pZip, pAddr, pDtlAddr, pRefAddr, pBizTpy, pbizCond, pEntpHtag1, pEntpHtag2, pEntpHtag3, pEntpIntr, pHomePg, pEmpCnt, pEstDy))
	// stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_ENTP_UPT_PROC('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', %v, '%v', :1)",
	// 	pLang, pEntpMemNo, pRepNm, pTelNo, pSmsRecvYn, pEmail, pEmailRecvYn, pInfoEqYn, pZip, pAddr, pDtlAddr, pRefAddr, pBizTpy, pbizCond, pEntpHtag1, pEntpHtag2, pEntpHtag3, pEntpIntr, pHomePg, pEmpCnt, pEstDy),
	// 	ora.I64, /* RTN_CD */
	// 	ora.S,   /* RTN_MSG */
	// 	ora.S,   /* SET_MEM_NO */
	// )

	log.Debug(fmt.Sprintf("CALL ZSP_ENTP_UPT_PROC_V2('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', %v, '%v','%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRepNm, pTelNo, pSmsRecvYn, pEmail, pEmailRecvYn, pInfoEqYn, pZip, pAddr, pDtlAddr, pRefAddr, pBizTpy, pbizCond, pEntpHtag1, pEntpHtag2, pEntpHtag3, pEntpIntr, pHomePg, pEmpCnt, pEstDy, pBizTpyCd, pEntpProfile, pBizIntro, pEntpCapital, pEntpTotalSales, pEntpTypeCd, pLocation))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_ENTP_UPT_PROC_V2('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', %v, '%v','%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRepNm, pTelNo, pSmsRecvYn, pEmail, pEmailRecvYn, pInfoEqYn, pZip, pAddr, pDtlAddr, pRefAddr, pBizTpy, pbizCond, pEntpHtag1, pEntpHtag2, pEntpHtag3, pEntpIntr, pHomePg, pEmpCnt, pEstDy, pBizTpyCd, pEntpProfile, pBizIntro, pEntpCapital, pEntpTotalSales, pEntpTypeCd, pLocation),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
		ora.S,   /* SET_MEM_NO */
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
		rtnCd    int64
		rtnMsg   string
		setMemNo string
	)

	entpInfoUpdate := models.EntpInfoUpdate{}
	rtnEntpInfoUpdate := models.RtnEntpInfoUpdate{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)

			if rtnCd == 1 {
				setMemNo = procRset.Row[2].(string)

				entpInfoUpdate = models.EntpInfoUpdate{
					SetMemNo: setMemNo,
				}

				// 사업자등록증이 있을 경우
				if pBizRegYn == "Y" {
					// 사업자등록증 업로드
					nowDate := time.Now()
					dateFmt := fmt.Sprintf(nowDate.Format("20060102150405"))

					uploadPath, _ := beego.AppConfig.String("uploadpath")
					imgDir := uploadPath + "/biz/" + setMemNo

					// 폴더가 없을 경우 해당 폴더를 만들어준다.
					if _, err := os.Stat(imgDir); os.IsNotExist(err) {
						err = os.MkdirAll(imgDir, 0755)
						if err != nil {
							panic(err)
						}
					}

					// 기등록된 원본로고 파일 삭제
					orgFile := uploadPath + "/" + oriBizFile
					var errOrg = os.Remove(orgFile)
					if errOrg != nil {
						log.Debug("Origin File Remove failed: %v", errOrg)
					}

					// 사업자등록증 업로드
					log.Debug(fmt.Sprintf(imgDir+"/biz_%v_%v.%v", setMemNo, dateFmt, pEntpRegNoExt))
					// 원본이미지
					c.SaveToFile("entp_regno", fmt.Sprintf(imgDir+"/biz_%v_%v.%v", setMemNo, dateFmt, pEntpRegNoExt))

					bizFilePath := "/biz/" + setMemNo + "/biz_" + setMemNo + "_" + dateFmt + "." + pEntpRegNoExt

					// 사업자등록증 등록
					log.Debug(fmt.Sprintf("CALL ZSP_ENTP_INFO_UPT_SUB_PROC( '%v', '%v', '%v', :1)",
						pLang, setMemNo, bizFilePath))

					stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_ENTP_INFO_UPT_SUB_PROC( '%v', '%v', '%v', :1)",
						pLang, setMemNo, bizFilePath),
						ora.I64, ora.S)

					defer stmtProcCall.Close()
					if err != nil {
						panic(err)
					}
					procRset := &ora.Rset{}
					_, err = stmtProcCall.Exe(procRset)

					if err != nil {
						panic(err)
					}
				}

				// 기업로고가 있을 경우
				if pImgYn == "Y" {
					// 로고 업로드
					nowDate := time.Now()
					dateFmt := fmt.Sprintf(nowDate.Format("20060102150405"))

					uploadPath, _ := beego.AppConfig.String("uploadpath")
					imgDir := uploadPath + "/logo/" + setMemNo

					// 폴더가 없을 경우 해당 폴더를 만들어준다.
					if _, err := os.Stat(imgDir); os.IsNotExist(err) {
						err = os.MkdirAll(imgDir, 0755)
						if err != nil {
							panic(err)
						}
					}

					// 기등록된 원본로고 파일 삭제
					orgFile := imgDir + "/ori_" + setMemNo + "." + oriLogoFileExt
					var errOrg = os.Remove(orgFile)
					if errOrg != nil {
						log.Debug("Origin File Remove failed: %v", errOrg)
					}

					// 기등록된 리사이징 이미지 파일 삭제
					oriLogoFilePath := uploadPath + oriLogoFile
					var err200 = os.Remove(oriLogoFilePath)
					if err200 != nil {
						log.Debug("Resizing File Remove failed: %v", err200)
					}

					// 기업로고이미지 업로드
					log.Debug(fmt.Sprintf(imgDir+"/ori_%v.%v", setMemNo, pEntpLogoExt))
					// 원본이미지
					c.SaveToFile("entp_logo", fmt.Sprintf(imgDir+"/ori_%v.%v", setMemNo, pEntpLogoExt))

					oriLogoImgPath := "/logo/" + setMemNo + "/ori_" + setMemNo + "." + pEntpLogoExt
					logoImgPath := "/logo/" + setMemNo + "/" + setMemNo + "_" + dateFmt + "." + pEntpLogoExt

					// sbsson 이미지 리사이징
					height := 200
					n_w, n_h := GetImageSize(uploadPath + oriLogoImgPath)
					if n_h < 200 {
						height = n_h
					}
					ga := (n_h * 200) / n_w
					height = round(float64(ga))

					src, err := imaging.Open(uploadPath + oriLogoImgPath)
					if err != nil {
						log.Debug("Open failed: %v", err)
					}

					// 200 리사이징 이미지
					rszImg200 := imaging.Resize(src, 200, 0, imaging.Lanczos)

					// sbsson 이미지 리사이징
					dst := imaging.New(300, 300, color.RGBA{255, 255, 255, 255})
					dst = imaging.Paste(dst, rszImg200, image.Pt(50, 150-(height/2)))
					dst = imaging.Resize(dst, 200, 0, imaging.Lanczos)

					err = imaging.Save(dst, imgDir+"/"+setMemNo+"_"+dateFmt+"."+pEntpLogoExt)
					if err != nil {
						log.Debug("Save failed rszImg200: %v", err)
					}

					// 기업로고 이미지 등록
					log.Debug(fmt.Sprintf("CALL ZSP_ENTP_INFO_UPT_SUB2_PROC( '%v', '%v', '%v', :1)",
						pLang, setMemNo, logoImgPath))

					stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_ENTP_INFO_UPT_SUB2_PROC( '%v', '%v', '%v', :1)",
						pLang, setMemNo, logoImgPath),
						ora.I64, ora.S)

					defer stmtProcCall.Close()
					if err != nil {
						panic(err)
					}
					procRset := &ora.Rset{}
					_, err = stmtProcCall.Exe(procRset)

					if err != nil {
						panic(err)
					}
				} else if pImgYn == "D" { // 로고를 삭제했을 경우

					uploadPath, _ := beego.AppConfig.String("uploadpath")
					imgDir := uploadPath + "/logo/" + setMemNo

					// 기등록된 원본로고 파일 삭제
					orgFile := imgDir + "/ori_" + setMemNo + "." + oriLogoFileExt
					var errOrg = os.Remove(orgFile)
					if errOrg != nil {
						log.Debug("Origin File Remove failed: %v", errOrg)
					}

					// 기등록된 리사이징 이미지 파일 삭제
					oriLogoFilePath := uploadPath + oriLogoFile
					var err200 = os.Remove(oriLogoFilePath)
					if err200 != nil {
						log.Debug("Resizing File Remove failed: %v", err200)
					}

					logoImgPath := ""

					// 기업로고 이미지 삭제
					log.Debug(fmt.Sprintf("CALL ZSP_ENTP_INFO_UPT_SUB2_PROC( '%v', '%v', '%v', :1)",
						pLang, setMemNo, logoImgPath))

					stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_ENTP_INFO_UPT_SUB2_PROC( '%v', '%v', '%v', :1)",
						pLang, setMemNo, logoImgPath),
						ora.I64, ora.S)

					defer stmtProcCall.Close()
					if err != nil {
						panic(err)
					}
					procRset := &ora.Rset{}
					_, err = stmtProcCall.Exe(procRset)

					if err != nil {
						panic(err)
					}
				}

			}
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnEntpInfoUpdate = models.RtnEntpInfoUpdate{
			RtnCd:   rtnCd,
			RtnMsg:  rtnMsg,
			RtnData: entpInfoUpdate,
		}
	}

	// End : Entp Info Update Process

	c.Data["json"] = &rtnEntpInfoUpdate
	c.ServeJSON()
}
