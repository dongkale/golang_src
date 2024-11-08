package controllers

import (
	"fmt"
	"os"
	"time"

	"emst.ziggam.com/models"
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
		c.Ctx.Redirect(302, "/common/login")
	}

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := mem_no
	//pChkEntpMemNo := c.GetString("entp_mem_no")

	//pEntpKoNm := c.GetString("entp_ko_nm")
	pRepNm := c.GetString("rep_nm")
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

	log.Debug("CALL SP_EMS_ENTP_INFO_UPT_PROC('%v','%v','%v','%v', %v ,'%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v', :1)",
		pLang, pEntpMemNo, pRepNm, pEstDy, pEmpCnt, pBizTpy, pBizCond, pPpChrgNm, pPpChrgTelNo, pZip, pAddr, pDtlAddr, pRefAddr, pEntpHtag1, pEntpHtag2, pEntpHtag3, pEntpIntr, pHomePg, pEmail)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ENTP_INFO_UPT_PROC('%v','%v','%v','%v', %v ,'%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v', :1)",
		pLang, pEntpMemNo, pRepNm, pEstDy, pEmpCnt, pBizTpy, pBizCond, pPpChrgNm, pPpChrgTelNo, pZip, pAddr, pDtlAddr, pRefAddr, pEntpHtag1, pEntpHtag2, pEntpHtag3, pEntpIntr, pHomePg, pEmail),
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

					// 기업로고이미지 업로드
					log.Debug(fmt.Sprintf(imgDir+"/ori_%v.%v", setMemNo, pLogoExt))
					// 원본이미지
					c.SaveToFile("entp_logo", fmt.Sprintf(imgDir+"/ori_%v.%v", setMemNo, pLogoExt))

					oriLogoImgPath := "/logo/" + setMemNo + "/ori_" + setMemNo + "." + pLogoExt
					logoImgPath := "/logo/" + setMemNo + "/" + setMemNo + "_" + dateFmt + "." + pLogoExt

					src, err := imaging.Open(uploadPath + oriLogoImgPath)
					if err != nil {
						log.Debug("Open failed: %v", err)
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
						log.Debug("Resising File Remove failed: %v", err200)
					}

					// 200 리사이징 이미지
					rszImg200 := imaging.Resize(src, 200, 0, imaging.Lanczos)
					err = imaging.Save(rszImg200, imgDir+"/"+setMemNo+"_"+dateFmt+"."+pLogoExt)
					if err != nil {
						log.Debug("Save failed rszImg200: %v", err)
					}

					// 기업로고 이미지 등록
					log.Debug("CALL SP_EMS_ENTP_INFO_UPT_SUB_PROC( '%v', '%v', '%v', :1)",
						pLang, setMemNo, logoImgPath)

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
	}

	// End : Entp Info Update Process

	c.Data["json"] = &rtnEntpInfoUpdate
	c.ServeJSON()
}
