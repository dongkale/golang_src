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

type AdminBannerInsertController struct {
	beego.Controller
}

func (c *AdminBannerInsertController) Post() {

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
	pBnrSn := c.GetString("bnr_sn")
	pBnrTitle := c.GetString("bnr_title")
	pListTitle := c.GetString("list_title")
	pListLinkUrl := c.GetString("list_link_url")
	pSdy := c.GetString("sdy")
	pEdy := c.GetString("edy")
	pLnkGbnCd := c.GetString("lnk_gbn_cd")
	pBnrKndCd := c.GetString("bnr_knd_cd")
	pLnkGbnVal := c.GetString("lnk_gbn_val")
	pLnkGbnListVal := c.GetString("lnk_gbn_list_val")
	pEntpMemNo := c.GetString("entp_mem_no")
	pRecruitSn := c.GetString("recruit_sn")
	pSn := c.GetString("sn")
	pUseYn := c.GetString("use_yn")

	// 배너 영역
	pImgYn := c.GetString("img_yn")
	if pImgYn == "" {
		pImgYn = "N"
	}
	pImgExt := c.GetString("img_ext")
	oriImgFile := c.GetString("ori_img_file")
	oriThumbImgFile := c.GetString("ori_thumb_img_file")
	//oriImgFileExt := c.GetString("ori_img_file_ext")
	// 배너 영역

	// 그룹 이미지 영역
	pListImgYn := c.GetString("list_img_yn")
	if pListImgYn == "" {
		pListImgYn = "N"
	}
	pListImgExt := c.GetString("list_img_ext")
	listOriImgFile := c.GetString("list_ori_img_file")
	listOriThumbImgFile := c.GetString("list_ori_thumb_img_file")

	log.Debug("Image Path Log ( '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pImgExt, oriImgFile, oriThumbImgFile, pListImgExt, listOriImgFile, listOriThumbImgFile)

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Banner Process

	log.Debug("CALL SP_EMS_ADMIN_BANNER_PROC('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', %v, '%v', '%v', '%v', :1)",
		pLang, pCuCd, pBnrTitle, pSdy, pEdy, pLnkGbnCd, pBnrKndCd, pLnkGbnVal, pLnkGbnListVal, pEntpMemNo, pRecruitSn, pSn, pUseYn, pBnrSn, pListTitle, pListLinkUrl)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_BANNER_PROC('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', %v, '%v', '%v', '%v', '%v', :1)",
		pLang, pCuCd, pBnrTitle, pSdy, pEdy, pLnkGbnCd, pBnrKndCd, pLnkGbnVal, pLnkGbnListVal, pEntpMemNo, pRecruitSn, pSn, pUseYn, pBnrSn, pListTitle, pListLinkUrl),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
		ora.S,   /* SET_BNR_SN */
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
		setBnrSn string
	)

	adminBannerInsert := models.AdminBannerInsert{}
	rtnAdminBannerInsert := models.RtnAdminBannerInsert{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)

			if rtnCd == 1 {
				setBnrSn = procRset.Row[2].(string)

				adminBannerInsert = models.AdminBannerInsert{
					SetBnrSn: setBnrSn,
				}

				if pCuCd == "C" {
					// 로고 업로드
					nowDate := time.Now()
					dateFmt := fmt.Sprintf(nowDate.Format("20060102150405"))

					uploadPath, _ := beego.AppConfig.String("uploadpath")
					imgDir := uploadPath + "/banner/" + setBnrSn

					// 폴더가 없을 경우 해당 폴더를 만들어준다.
					if _, err := os.Stat(imgDir); os.IsNotExist(err) {
						err = os.MkdirAll(imgDir, 0755)
						if err != nil {
							panic(err)
						}
					}

					// 이미지 업로드
					log.Debug(fmt.Sprintf("Image Dir Log : " + imgDir+"/%v_%v.%v", setBnrSn, dateFmt, pImgExt))
					// 원본이미지
					c.SaveToFile("bnr_img", fmt.Sprintf(imgDir+"/%v_%v.%v", setBnrSn, dateFmt, pImgExt))

					oriImgPath := "/banner/" + setBnrSn + "/" + setBnrSn + "_" + dateFmt + "." + pImgExt
					imgPath := "/banner/" + setBnrSn + "/" + setBnrSn + "_" + dateFmt + "." + pImgExt
					thumbImgPath := "/banner/" + setBnrSn + "/100_" + setBnrSn + "_" + dateFmt + "." + pImgExt

					src, err := imaging.Open(uploadPath + oriImgPath)
					if err != nil {
						log.Debug("Open failed Banner Image :  %v", err)
					}

					// 100 리사이징 이미지
					rszImg100 := imaging.Resize(src, 100, 0, imaging.Lanczos)
					err = imaging.Save(rszImg100, imgDir+"/100_"+setBnrSn+"_"+dateFmt+"."+pImgExt)
					if err != nil {
						log.Debug("Save failed rszImg100: %v", err)
					}

					var listImgPath = ""
					var listThumbImgPath = ""

					if (pListImgYn == "Y") {
						// 리스트 이미지 업로드

						listImgDir := uploadPath + "/list/" + setBnrSn
						if _, err := os.Stat(listImgDir); os.IsNotExist(err) {
							err = os.MkdirAll(listImgDir, 0755)
							if err != nil {
								panic(err)
							}
						}

						log.Debug(fmt.Sprintf("List Image Dir Log : " + listImgDir + "/list_%v_%v.%v", setBnrSn, dateFmt, pImgExt))
						// 원본이미지
						c.SaveToFile("list_img", fmt.Sprintf(listImgDir+"/list_%v_%v.%v", setBnrSn, dateFmt, pListImgExt))

						listOriImgPath := "/list/" + setBnrSn + "/list_" + setBnrSn + "_" + dateFmt + "." + pListImgExt
						listImgPath = "/list/" + setBnrSn + "/list_" + setBnrSn + "_" + dateFmt + "." + pListImgExt
						listThumbImgPath = "/list/" + setBnrSn + "/list_100_" + setBnrSn + "_" + dateFmt + "." + pListImgExt

						list_src, list_err := imaging.Open(uploadPath + listOriImgPath)
						if err != nil {
							log.Debug("Open failed List Image : %v", err)
						}

						// 100 리사이징 이미지
						listRszImg100 := imaging.Resize(list_src, 100, 0, imaging.Lanczos)
						list_err = imaging.Save(listRszImg100, listImgDir+"/list_100_"+setBnrSn+"_"+dateFmt+"."+pListImgExt)
						if list_err != nil {
							log.Debug("Save failed List rszImg100: %v", list_err)
						}
					}

					// 이미지 등록
					log.Debug("CALL SP_EMS_ADMIN_BANNER_SUB_PROC( '%v', '%v', '%v', '%v', '%v', '%v', :1)",
						pLang, setBnrSn, imgPath, thumbImgPath, listImgPath, listThumbImgPath)

					stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_BANNER_SUB_PROC( '%v', '%v', '%v', '%v',  '%v', '%v', :1)",
						pLang, setBnrSn, imgPath, thumbImgPath, listImgPath, listThumbImgPath),
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
					if pImgYn == "Y" || pListImgYn == "Y" {

						// 배너 이미지 경로
						imgPath := ""
						thumbImgPath := ""
						// 리스트 이미지 경로
						listImgPath := ""
						listThumbImgPath := ""

						if pImgYn == "Y" {
							// 배너 이미지 == "Y"일때 배너 이미지 등록 데이터 생성
							// 로고 업로드
							nowDate := time.Now()
							dateFmt := fmt.Sprintf(nowDate.Format("20060102150405"))

							uploadPath, _ := beego.AppConfig.String("uploadpath")
							imgDir := uploadPath + "/banner/" + setBnrSn

							// 폴더가 없을 경우 해당 폴더를 만들어준다.
							if _, err := os.Stat(imgDir); os.IsNotExist(err) {
								err = os.MkdirAll(imgDir, 0755)
								if err != nil {
									panic(err)
								}
							}

							// 이미지 업로드
							log.Debug(fmt.Sprintf(imgDir+"/%v_%v.%v", setBnrSn, dateFmt, pImgExt))
							// 원본이미지
							c.SaveToFile("bnr_img", fmt.Sprintf(imgDir+"/%v_%v.%v", setBnrSn, dateFmt, pImgExt))

							oriImgPath := "/banner/" + setBnrSn + "/" + setBnrSn + "_" + dateFmt + "." + pImgExt
							imgPath = "/banner/" + setBnrSn + "/" + setBnrSn + "_" + dateFmt + "." + pImgExt
							thumbImgPath = "/banner/" + setBnrSn + "/100_" + setBnrSn + "_" + dateFmt + "." + pImgExt

							src, err := imaging.Open(uploadPath + oriImgPath)
							if err != nil {
								log.Debug("Open failed: %v", err)
							}

							// 기등록된 원본로고 파일 삭제
							//orgFile := imgDir + "/" + setBnrSn + "." + oriImgFileExt
							orgFile := uploadPath + oriImgFile
							var errOrg = os.Remove(orgFile)
							if errOrg != nil {
								log.Debug("Origin File Remove failed: %v", errOrg)
							}

							// 기등록된 리사이징 이미지 파일 삭제
							oriImgFilePath := uploadPath + oriThumbImgFile
							var err100 = os.Remove(oriImgFilePath)
							if err100 != nil {
								log.Debug("Resising File Remove failed: %v", err100)
							}

							// 100 리사이징 이미지
							rszImg100 := imaging.Resize(src, 100, 0, imaging.Lanczos)
							err = imaging.Save(rszImg100, imgDir+"/100_"+setBnrSn+"_"+dateFmt+"."+pImgExt)
							if err != nil {
								log.Debug("Save failed rszImg100: %v", err)
							}
						}
						// 배너 이미지 == "Y"일때 배너 이미지 등록 데이터 생성

						// 그룹 이미지 == "Y"일때 그룹 이미지 등록 데이터 생성

						if pListImgYn == "Y" {
							// 배너 이미지 == "Y"일때 배너 이미지 등록 데이터 생성
							// 로고 업로드
							nowDate := time.Now()
							dateFmt := fmt.Sprintf(nowDate.Format("20060102150405"))

							uploadPath, _ := beego.AppConfig.String("uploadpath")
							imgDir := uploadPath + "/list/" + setBnrSn

							// 폴더가 없을 경우 해당 폴더를 만들어준다.
							if _, err := os.Stat(imgDir); os.IsNotExist(err) {
								err = os.MkdirAll(imgDir, 0755)
								if err != nil {
									panic(err)
								}
							}

							// 이미지 업로드
							log.Debug(fmt.Sprintf(imgDir+"/list_%v_%v.%v", setBnrSn, dateFmt, pListImgExt))
							// 원본이미지
							c.SaveToFile("list_img", fmt.Sprintf(imgDir+"/list_%v_%v.%v", setBnrSn, dateFmt, pListImgExt))

							oriImgPath := "/list/" + setBnrSn + "/list_" + setBnrSn + "_" + dateFmt + "." + pListImgExt
							listImgPath = "/list/" + setBnrSn + "/list_" + setBnrSn + "_" + dateFmt + "." + pListImgExt
							listThumbImgPath = "/list/" + setBnrSn + "/list_100_" + setBnrSn + "_" + dateFmt + "." + pListImgExt

							src, err := imaging.Open(uploadPath + oriImgPath)
							if err != nil {
								log.Debug("Open failed List Image : %v", err)
							}

							// 기등록된 원본로고 파일 삭제
							//orgFile := imgDir + "/" + setBnrSn + "." + oriImgFileExt
							orgFile := uploadPath + listOriImgFile

							var errOrg = os.Remove(orgFile)
							if errOrg != nil {
								log.Debug("Origin List File Remove failed: %v", errOrg)
							}

							// 기등록된 리사이징 이미지 파일 삭제
							oriImgFilePath := uploadPath + listOriThumbImgFile
							var err100 = os.Remove(oriImgFilePath)
							if err100 != nil {
								log.Debug("Resising List File Remove failed: %v", err100)
							}

							// 100 리사이징 이미지
							rszImg100 := imaging.Resize(src, 100, 0, imaging.Lanczos)
							err = imaging.Save(rszImg100, imgDir+"/list_100_"+setBnrSn+"_"+dateFmt+"."+pListImgExt)
							if err != nil {
								log.Debug("Save failed List rszImg100: %v", err)
							}
						}

						// 그룹 이미지 == "Y"일때 그룹 이미지 등록 데이터 생성

						// 이미지 등록
						log.Debug("CALL SP_EMS_ADMIN_BANNER_SUB_PROC( '%v', '%v', '%v', '%v', '%v', '%v', :1)",
							pLang, setBnrSn, imgPath, thumbImgPath, listImgPath, listThumbImgPath)

						stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_BANNER_SUB_PROC( '%v', '%v', '%v', '%v', '%v', '%v', :1)",
							pLang, setBnrSn, imgPath, thumbImgPath, listImgPath, listThumbImgPath),
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

		rtnAdminBannerInsert = models.RtnAdminBannerInsert{
			RtnCd:   rtnCd,
			RtnMsg:  rtnMsg,
			RtnData: adminBannerInsert,
		}
	}

	// End : Banner Process

	c.Data["json"] = &rtnAdminBannerInsert
	c.ServeJSON()
}
