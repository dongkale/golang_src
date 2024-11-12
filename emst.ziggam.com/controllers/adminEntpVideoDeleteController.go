package controllers

import (
	"fmt"

	"emst.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AdminEntpVideoDeleteController struct {
	beego.Controller
}

func (c *AdminEntpVideoDeleteController) Post() {

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
	pEnptMemNo := c.GetString("entp_mem_no")
	pGbnCd := c.GetString("gbn_cd")
	pVdNo := c.GetString("vd_no")
	filePathNm := ""

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Admin Entp Video Delete Process
	log.Debug("CALL SP_EMS_ADMIN_ENTP_VIDEO_PROC('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pGbnCd, pEnptMemNo, pVdNo, filePathNm)
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_ADMIN_ENTP_VIDEO_PROC('%v', '%v', '%v', '%v', '%v', :1)",
		pLang, pGbnCd, pEnptMemNo, pVdNo, filePathNm),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
		ora.S,   /* RTN_DATA */
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
		rtnCd   int64
		rtnMsg  string
		rtnData string
	)

	rtnEntpVideoDelete := models.RtnEntpVideoDelete{}

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
			rtnData = procRset.Row[2].(string)

			if rtnCd == 1 {
				// 기존 파일 삭제처리
				go deleteEvp(rtnData)
			}

		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnEntpVideoDelete = models.RtnEntpVideoDelete{
			RtnCd:   rtnCd,
			RtnMsg:  rtnMsg,
			RtnData: rtnData,
		}
	}

	// End : Admin Entp Video Delete Process

	c.Data["json"] = &rtnEntpVideoDelete
	c.ServeJSON()
}

func deleteEvp(prevFilePath string) {
	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	//uploadPath := beego.AppConfig.String("uploadpath")

	// Object Storage 파일 물리적 경로
	//objectFilePath := uploadPath + "/entp_video/" + pEnptMemNo + "/" + pVdNo + "/evp_" + dateFmt + ext
	
	awsAccessKeyID, err := beego.AppConfig.String("accessKey")  //"6JSYw9NY9JAsO7XynmbO"
	if err != nil {
		log.Error("Failed to get accessKey: %v", err)
		return
	}
	awsSecretAccessKey, err := beego.AppConfig.String("secretKey")  //"uARPb8klZto1zXTIC4G0TXt0aFrsAM9FkYzlX1vT"
	if err != nil {
		log.Error("Failed to get secretKey: %v", err)
		return
	}
	awsBucketName, err := beego.AppConfig.String("bucketName") //"qrate-interview-service-test"
	if err != nil {
		log.Error("Failed to get bucketName: %v", err)
		return
	}

	//objectPath := "/input/entp_video/" + pEnptMemNo + "/" + pVdNo + "/evp_" + dateFmt + ext
	//objectOutPath := "/output/entp_video/" + pEnptMemNo + "/" + pVdNo + "/"
	objectPrevFile := "/output" + prevFilePath

	//outputFileNm := "480p" //dateFmt + ext

	//fmt.Println("objectPath : ", objectPath)
	fmt.Println("objectPrevFile : ", objectPrevFile)

	key := aws.String(objectPrevFile)
	bucket := aws.String(awsBucketName)
	//key := aws.String(objectPath + fileName)

	// Configure to use nCloud Server
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
		Endpoint:         aws.String("https://kr.object.ncloudstorage.com"),
		Region:           aws.String("kr-standard"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession := session.New(s3Config)

	s3Client := s3.New(newSession)

	ret, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		fmt.Printf("Failed to delete data to %s/%s, %s\n", *bucket, *key, err.Error())
		return
	} else {
		fmt.Printf("Successfully Delete Object %s and deleted data with key %s\n", *bucket, *key)
	}

	fmt.Println("ret : ", ret)

	/*
		// 업로드 된 파일 서버에서 삭제처리
		var errDel = os.Remove(objectFilePath)
		if errDel != nil {
			log.Debug("Remove failed: %v", errDel)
		}
	*/
}
