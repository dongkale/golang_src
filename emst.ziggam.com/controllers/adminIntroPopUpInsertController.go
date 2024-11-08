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

type AdminIntroPopUpInsertController struct {
	beego.Controller
}

func (c *AdminIntroPopUpInsertController) Post() {

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
	pCuCd := c.GetString("cu_cd")
	pIntroSn := c.GetString("intro_sn")
	pIntroTitle := c.GetString("intro_title")
	pSdy := c.GetString("sdy")
	pEdy := c.GetString("edy")
	pLnkGbnCd := c.GetString("lnk_gbn_cd")
	pLnkGbnVal := c.GetString("lnk_gbn_val")
	pEntpMemNo := c.GetString("entp_mem_no")
	pRecruitSn := c.GetString("recruit_sn")
	pSn := c.GetString("sn")
	pUseYn := c.GetString("use_yn")

	pImgYn := c.GetString("img_yn")
	if pImgYn == "" {
		pImgYn = "N"
	}
	pImgExt := c.GetString("img_ext")
	oriImgFile := c.GetString("ori_img_file")
	oriThumbImgFile := c.GetString("ori_thumb_img_file")
	//oriImgFileExt := c.GetString("ori_img_file_ext")

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

	log.Debug("CALL SP_EMS_ADMIN_INTRO_PROC('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', %v, '%v', '%v', :1)",
		pLang, pCuCd, pIntroTitle, pSdy, pEdy, pLnkGbnCd, pLnkGbnVal, pEntpMemNo, pRecruitSn, pSn, pUseYn, pIntroSn)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_INTRO_PROC('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', %v, '%v', '%v', :1)",
		pLang, pCuCd, pIntroTitle, pSdy, pEdy, pLnkGbnCd, pLnkGbnVal, pEntpMemNo, pRecruitSn, pSn, pUseYn, pIntroSn),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
		ora.S,   /* SET_INTRO_SN */
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
		rtnCd      int64
		rtnMsg     string
		setIntroSn string
	)

	adminIntroPopUpInsert := models.AdminIntroPopUpInsert{}
	rtnAdminIntroPopUpInsert := models.RtnAdminIntroPopUpInsert{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)

			if rtnCd == 1 {
				setIntroSn = procRset.Row[2].(string)

				adminIntroPopUpInsert = models.AdminIntroPopUpInsert{
					SetIntroSn: setIntroSn,
				}

				if pCuCd == "C" {
					// 로고 업로드
					nowDate := time.Now()
					dateFmt := fmt.Sprintf(nowDate.Format("20060102150405"))

					uploadPath, _ := beego.AppConfig.String("uploadpath")
					imgDir := uploadPath + "/intro/" + setIntroSn

					// 폴더가 없을 경우 해당 폴더를 만들어준다.
					if _, err := os.Stat(imgDir); os.IsNotExist(err) {
						err = os.MkdirAll(imgDir, 0755)
						if err != nil {
							panic(err)
						}
					}

					// 이미지 업로드
					log.Debug(fmt.Sprintf(imgDir+"/%v_%v.%v", setIntroSn, dateFmt, pImgExt))
					// 원본이미지
					c.SaveToFile("intro_img", fmt.Sprintf(imgDir+"/%v_%v.%v", setIntroSn, dateFmt, pImgExt))

					oriImgPath := "/intro/" + setIntroSn + "/" + setIntroSn + "_" + dateFmt + "." + pImgExt
					imgPath := "/intro/" + setIntroSn + "/" + setIntroSn + "_" + dateFmt + "." + pImgExt
					thumbImgPath := "/intro/" + setIntroSn + "/50_" + setIntroSn + "_" + dateFmt + "." + pImgExt

					src, err := imaging.Open(uploadPath + oriImgPath)
					if err != nil {
						log.Debug("Open failed: %v", err)
					}

					// 50 리사이징 이미지
					rszImg50 := imaging.Resize(src, 50, 0, imaging.Lanczos)
					err = imaging.Save(rszImg50, imgDir+"/50_"+setIntroSn+"_"+dateFmt+"."+pImgExt)
					if err != nil {
						log.Debug("Save failed rszImg50: %v", err)
					}

					// 이미지 등록
					log.Debug("CALL SP_EMS_ADMIN_INTRO_SUB_PROC( '%v', '%v', '%v', '%v', :1)",
						pLang, setIntroSn, imgPath, thumbImgPath)

					stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_INTRO_SUB_PROC( '%v', '%v', '%v', '%v', :1)",
						pLang, setIntroSn, imgPath, thumbImgPath),
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
				} else {
					if pImgYn == "Y" {
						// 로고 업로드
						nowDate := time.Now()
						dateFmt := fmt.Sprintf(nowDate.Format("20060102150405"))

						uploadPath, _ := beego.AppConfig.String("uploadpath")
						imgDir := uploadPath + "/intro/" + setIntroSn

						// 폴더가 없을 경우 해당 폴더를 만들어준다.
						if _, err := os.Stat(imgDir); os.IsNotExist(err) {
							err = os.MkdirAll(imgDir, 0755)
							if err != nil {
								panic(err)
							}
						}

						// 이미지 업로드
						log.Debug(fmt.Sprintf(imgDir+"/%v_%v.%v", setIntroSn, dateFmt, pImgExt))
						// 원본이미지
						c.SaveToFile("intro_img", fmt.Sprintf(imgDir+"/%v_%v.%v", setIntroSn, dateFmt, pImgExt))

						oriImgPath := "/intro/" + setIntroSn + "/" + setIntroSn + "_" + dateFmt + "." + pImgExt
						imgPath := "/intro/" + setIntroSn + "/" + setIntroSn + "_" + dateFmt + "." + pImgExt
						thumbImgPath := "/intro/" + setIntroSn + "/50_" + setIntroSn + "_" + dateFmt + "." + pImgExt

						src, err := imaging.Open(uploadPath + oriImgPath)
						if err != nil {
							log.Debug("Open failed: %v", err)
						}

						// 기등록된 원본로고 파일 삭제
						//orgFile := imgDir + "/" + setIntroSn + "." + oriImgFileExt
						orgFile := uploadPath + oriImgFile
						var errOrg = os.Remove(orgFile)
						if errOrg != nil {
							log.Debug("Origin File Remove failed: %v", errOrg)
						}

						// 기등록된 리사이징 이미지 파일 삭제
						oriImgFilePath := uploadPath + oriThumbImgFile
						var err50 = os.Remove(oriImgFilePath)
						if err50 != nil {
							log.Debug("Resising File Remove failed: %v", err50)
						}

						// 50 리사이징 이미지
						rszImg50 := imaging.Resize(src, 50, 0, imaging.Lanczos)
						err = imaging.Save(rszImg50, imgDir+"/50_"+setIntroSn+"_"+dateFmt+"."+pImgExt)
						if err != nil {
							log.Debug("Save failed rszImg50: %v", err)
						}

						// 이미지 등록
						log.Debug("CALL SP_EMS_ADMIN_INTRO_SUB_PROC( '%v', '%v', '%v', '%v', :1)",
							pLang, setIntroSn, imgPath, thumbImgPath)

						stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_INTRO_SUB_PROC( '%v', '%v', '%v', '%v', :1)",
							pLang, setIntroSn, imgPath, thumbImgPath),
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
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnAdminIntroPopUpInsert = models.RtnAdminIntroPopUpInsert{
			RtnCd:   rtnCd,
			RtnMsg:  rtnMsg,
			RtnData: adminIntroPopUpInsert,
		}
	}

	// End : Entp Info Update Process

	c.Data["json"] = &rtnAdminIntroPopUpInsert
	c.ServeJSON()
}
