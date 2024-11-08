package controllers

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"time"

	"emst.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/disintegration/imaging"
	"gopkg.in/rana/ora.v4"
)

type EntpInfoUpdateV2Controller struct {
	beego.Controller
}

func (c *EntpInfoUpdateV2Controller) Post() {

	/*
		file, header, err := c.GetFile("entp_logo")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer file.Close()
		fileName := header.Filename
		ext := filepath.Ext(fileName)
	*/

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Data["json"] = &models.DefaultResult{RtnCd: 99, RtnMsg: "/login"}
		c.ServeJSON()
		return
	}

	pLang, _ := beego.AppConfig.String("lang")

	//pEntpMemNo := mem_no
	//pChkEntpMemNo := c.GetString("entp_mem_no")
	pEntpMemNo := c.GetString("entp_mem_no")
	//pEntpKoNm := c.GetString("entp_ko_nm")
	pRepNm := c.GetString("rep_nm")
	pTelNo := c.GetString("tel_no")
	pEstDy := c.GetString("est_dy")
	pEmpCnt := c.GetString("emp_cnt")
	pBizTpy := c.GetString("biz_tpy")
	pBizCond := c.GetString("biz_cond")
	pPpChrgNm := c.GetString("pp_chrg_nm")
	pPpChrgTelNo := c.GetString("pp_chrg_tel_no")
	pEmail := c.GetString("email")
	pZip := c.GetString("zip")
	pAddr := c.GetString("addr")
	pDtlAddr := c.GetString("dtl_addr")
	pRefAddr := c.GetString("ref_addr")
	pEntpHtag1 := c.GetString("entp_htag1")
	pEntpHtag2 := c.GetString("entp_htag2")
	pEntpHtag3 := c.GetString("entp_htag3")
	pEntpIntr := c.GetString("entp_intr")
	pHomePg := c.GetString("home_pg")
	pImgYn := c.GetString("img_yn")
	if pImgYn == "" {
		pImgYn = "N"
	}
	pLogoExt := c.GetString("logo_ext")
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

	logs.Debug("pBizTpyCd: " + pBizTpyCd)
	logs.Debug("pBizTpy: " + pBizTpy)
	logs.Debug("pEntpProfile: " + pEntpProfile)
	logs.Debug("pBizIntro: " + pBizIntro)
	logs.Debug("pEntpCapital: " + pEntpCapital)
	logs.Debug("pEntpTotalSales: " + pEntpTotalSales)
	logs.Debug("pEntpTypeCd: " + pEntpTypeCd)
	logs.Debug("pEntpType: " + pEntpType)
	logs.Debug("pLocation: " + pLocation)

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

	logs.Debug(fmt.Sprintf("CALL SP_EMS_ENTP_INFO_UPT_V2_PROC('%v','%v','%v','%v','%v', %v ,'%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v', '%v','%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRepNm, pTelNo, pEstDy, pEmpCnt, pBizTpy, pBizCond, pPpChrgNm, pPpChrgTelNo, pZip, pAddr, pDtlAddr, pRefAddr, pEntpHtag1, pEntpHtag2, pEntpHtag3, pEntpIntr, pHomePg, pEmail, pBizTpyCd, pEntpProfile, pBizIntro, pEntpCapital, pEntpTotalSales, pEntpTypeCd, pLocation))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ENTP_INFO_UPT_V2_PROC('%v','%v','%v','%v','%v', %v ,'%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v', '%v','%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRepNm, pTelNo, pEstDy, pEmpCnt, pBizTpy, pBizCond, pPpChrgNm, pPpChrgTelNo, pZip, pAddr, pDtlAddr, pRefAddr, pEntpHtag1, pEntpHtag2, pEntpHtag3, pEntpIntr, pHomePg, pEmail, pBizTpyCd, pEntpProfile, pBizIntro, pEntpCapital, pEntpTotalSales, pEntpTypeCd, pLocation),
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
						logs.Debug("Origin File Remove failed: %v", errOrg)
					}

					// 기등록된 리사이징 이미지 파일 삭제
					oriLogoFilePath := uploadPath + oriLogoFile
					var err200 = os.Remove(oriLogoFilePath)
					if err200 != nil {
						logs.Debug("Resising File Remove failed: %v", err200)
					}

					// 기업로고이미지 업로드
					logs.Debug(fmt.Sprintf(imgDir+"/ori_%v.%v", setMemNo, pLogoExt))
					// 원본이미지
					c.SaveToFile("entp_logo", fmt.Sprintf(imgDir+"/ori_%v.%v", setMemNo, pLogoExt))

					oriLogoImgPath := "/logo/" + setMemNo + "/ori_" + setMemNo + "." + pLogoExt
					logoImgPath := "/logo/" + setMemNo + "/" + setMemNo + "_" + dateFmt + "." + pLogoExt

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
						logs.Debug("Open failed: %v", err)
					}

					// 200 리사이징 이미지
					// rszImg200 := imaging.Resize(src, 200, 0, imaging.Lanczos)
					// err = imaging.Save(rszImg200, imgDir+"/"+setMemNo+"_"+dateFmt+"."+pLogoExt)
					// if err != nil {
					// 	logs.Debug("Save failed rszImg200: %v", err)
					// }
					// 200 리사이징 이미지
					rszImg200 := imaging.Resize(src, 200, 0, imaging.Lanczos)

					// sbsson 이미지 리사이징
					dst := imaging.New(300, 300, color.RGBA{255, 255, 255, 255})
					dst = imaging.Paste(dst, rszImg200, image.Pt(50, 150-(height/2)))
					dst = imaging.Resize(dst, 200, 0, imaging.Lanczos)

					err = imaging.Save(dst, imgDir+"/"+setMemNo+"_"+dateFmt+"."+pLogoExt)
					if err != nil {
						logs.Debug("Save failed rszImg200: %v", err)
					}

					// 기업로고 이미지 등록
					logs.Debug(fmt.Sprintf("CALL SP_EMS_ENTP_INFO_UPT_SUB_PROC( '%v', '%v', '%v', :1)",
						pLang, setMemNo, logoImgPath))

					stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ENTP_INFO_UPT_SUB_PROC( '%v', '%v', '%v', :1)",
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

		logs.Debug(fmt.Sprintf("[EntpInfo][Save][%v] ===> RtnCd:%v, RtnMsg:%v, RtnData:%v", pEntpMemNo, rtnCd, rtnMsg, entpInfoUpdate))
	}

	// End : Entp Info Update Process

	c.Data["json"] = &rtnEntpInfoUpdate
	c.ServeJSON()
}

// 이미지 썸네일
func ErrorCon(err error) {
	if err != nil {
		logs.Debug(err)
	}
}

func GetImageSize(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		ErrorCon(err)
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		ErrorCon(err)
	}
	defer file.Close()

	return image.Width, image.Height
}
