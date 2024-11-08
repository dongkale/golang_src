package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"crypto/hmac"
	"crypto/sha256"

	"emst.ziggam.com/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AdminEntpVideoInsertController struct {
	beego.Controller
}

type EntpVideoResultJson struct {
	Jobs []struct {
		JobId string
	}
	Error struct {
		ErrorCode int64
		Message   string
	}
}

func (c *AdminEntpVideoInsertController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	tsession := c.StartSession()
	mem_no := tsession.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/common/login")
	}

	file, header, err := c.GetFile("entp_video")
	fileName := header.Filename
	ext := filepath.Ext(fileName)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	pLang, _ := beego.AppConfig.String("lang")

	// 기업회원번호
	pEnptMemNo := c.GetString("entp_mem_no")
	pGbnCd := c.GetString("gbn_cd")
	pVdNo := c.GetString("vd_no")

	if pVdNo == "" {
		pVdNo = "1"
	}

	defer file.Close()

	nowDate := time.Now()
	dateFmt := fmt.Sprintf(nowDate.Format("20060102150405"))

	uploadPath, _ := beego.AppConfig.String("uploadpath")

	// 영상프로필 업로드
	tempDir := uploadPath + "/entp_video/" + pEnptMemNo + "/" + pVdNo
	//filePath := "/" + pEntpMemNo + "/" + pRecrutSn + "/" + pPpMemNo + "/" + pQstSn + "/480p.mp4"
	//viewpath := beego.AppConfig.String("viewpath")

	// 폴더가 없을 경우 해당 폴더를 만들어준다.
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		err = os.MkdirAll(tempDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	log.Debug(fmt.Sprintf(tempDir+"/evp_%v%v", dateFmt, ext))
	// 원본파일 - Rename
	c.SaveToFile("entp_video", fmt.Sprintf(tempDir+"/evp_%v%v", dateFmt, ext))

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection
	objectDBPath := "/entp_video/" + pEnptMemNo + "/" + pVdNo + "/"
	outputFileNm := dateFmt[4:len(dateFmt)]

	rtnEntpVideoUploadFile := models.RtnEntpVideoUploadFile{}

	//filePathNm := objectDBPath + "480p.mp4"
	filePathNm := objectDBPath + outputFileNm + ".mp4"

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

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
			rtnData = procRset.Row[2].(string)

			if rtnCd == 1 {
				// Object Storage Upload Process
				go updateEvp(rtnData)
				go evp_upload(pEnptMemNo, pVdNo, dateFmt, ext, outputFileNm)
			}
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnEntpVideoUploadFile = models.RtnEntpVideoUploadFile{
			RtnCd:   rtnCd,
			RtnMsg:  rtnMsg,
			RtnData: rtnData,
		}
	}

	c.Data["json"] = &rtnEntpVideoUploadFile
	c.ServeJSON()
}

func evp_upload(pEnptMemNo string, pVdNo string, dateFmt string, ext string, outputFileNm string) {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	uploadPath, _:= beego.AppConfig.String("uploadpath")
	evp_presetId, _ := beego.AppConfig.String("evp_presetId")

	// Object Storage 파일 물리적 경로
	objectFilePath := uploadPath + "/entp_video/" + pEnptMemNo + "/" + pVdNo + "/evp_" + dateFmt + ext

	awsAccessKeyID, err := beego.AppConfig.String("accessKey")  //"6JSYw9NY9JAsO7XynmbO"
	if err != nil {
		log.Error("Failed to get accessKey from config: %v", err)
		return
	}
	awsSecretAccessKey, err := beego.AppConfig.String("secretKey")  //"uARPb8klZto1zXTIC4G0TXt0aFrsAM9FkYzlX1vT"
	if err != nil {
		log.Error("Failed to get secretKey from config: %v", err)
		return
	}
	awsBucketName, err := beego.AppConfig.String("bucketName") //"qrate-interview-service-test"
	if err != nil {
		log.Error("Failed to get bucketName from config: %v", err)
		return
	}
	

	f, e := os.Open(objectFilePath)
	if e != nil {
		fmt.Printf("err opening file: %s", e)
		return
	}
	defer f.Close()

	fileInfo, _ := f.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	f.Read(buffer)

	objectPath := "/input/entp_video/" + pEnptMemNo + "/" + pVdNo + "/evp_" + dateFmt + ext
	objectOutPath := "/output/entp_video/" + pEnptMemNo + "/" + pVdNo + "/"

	//outputFileNm := "480p" //dateFmt + ext

	//fmt.Println("objectPath : ", objectPath)
	//fmt.Println("objectOutPath : ", objectOutPath)

	key := aws.String(objectPath)
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

	// Upload a new object "uploadFile" with the string "Ziggam" to our "bucket".
	ret, err := s3Client.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(buffer), //strings.NewReader(objectFilePath),
		Bucket: bucket,
		Key:    key,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		fmt.Printf("Failed to upload data to %s/%s, %s\n", *bucket, *key, err.Error())
		return
	} else {
		fmt.Printf("Successfully created Object %s and uploaded data with key %s\n", *bucket, *key)
	}

	fmt.Println("ret : ", ret)

	// TransCode Job Api Request
	// 1. makeSignature

	now := time.Now()
	// convert time to millisecond(TimeStamp) 구하기
	unixNano := now.UnixNano()
	umillisec := unixNano / 1000000
	//fmt.Println("(correct)Millisecond : ", umillisec)

	var (
		space   = " "
		newLine = "\n"

		timestamp = umillisec
		apiKey, _    = beego.AppConfig.String("apiKey")
		accessKey, _ = beego.AppConfig.String("accessKey")
		secretKey, _ = beego.AppConfig.String("secretKey")
		//presetUrl    = beego.AppConfig.String("presetUrl")
		jobCreateUrl, _ = beego.AppConfig.String("jobCreateUrl")
	)

	/*
		var buf bytes.Buffer
		buf.WriteString("GET")
		buf.WriteString(space)
		buf.WriteString("/api/v1/presets")
		buf.WriteString(newLine)
		buf.WriteString(strconv.Itoa(int(timestamp)))
		buf.WriteString(newLine)
		buf.WriteString(apiKey)
		buf.WriteString(newLine)
		buf.WriteString(accessKey)
		resultPresetString := buf.String()

		presetApiSignature := VdTestComputeHmac256(resultPresetString, secretKey)
		//fmt.Println("presetApiSignature : ", presetApiSignature)

		// 2. Preset Operation
		req, err := http.NewRequest("GET", presetUrl, nil)

		req.Header.Add("x-ncp-apigw-timestamp", strconv.Itoa(int(timestamp)))
		req.Header.Add("x-ncp-apigw-api-key", apiKey)
		req.Header.Add("x-ncp-iam-access-key", accessKey)
		req.Header.Add("x-ncp-apigw-signature-v1", presetApiSignature)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("error : ", err)
		}
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Println("StringBody : ", string(body))
		presetApiSignature = ""

	*/

	//presetId := "9b7e8e2d-04c9-11e8-8379-00505685080f" // "480 SD"
	//presetId := "703d7fc0-961c-11e9-b460-005056855f38" // ziggam:480p
	//presetId := "a2c0aa06-987d-11e9-b460-005056855f38" // ziggam:720p
	//presetId := "dd814b45-987e-11e9-b460-005056855f38" // ziggam:1080p
	//presetId  "a2c0aa06-987d-11e9-b460-005056855f38"
	//presetId := "9b7e8e2d-04c9-11e8-8379-00505685080f"
	//presetId := "ab639551-2b7c-11e9-af94-005056851dca" // 720 HD

	// 3. Job Operation

	var jobbuf bytes.Buffer
	jobbuf.WriteString("POST")
	jobbuf.WriteString(space)
	jobbuf.WriteString("/api/v2/jobs")
	jobbuf.WriteString(newLine)
	jobbuf.WriteString(strconv.Itoa(int(timestamp)))
	jobbuf.WriteString(newLine)
	jobbuf.WriteString(apiKey)
	jobbuf.WriteString(newLine)
	jobbuf.WriteString(accessKey)
	resultJobString := jobbuf.String()

	// Job생성 Signature 만들기
	jobApiSignature := VdTestComputeHmac256(resultJobString, secretKey)
	//fmt.Println("jobApiSignature : ", jobApiSignature)

	// Job생성 API Request
	jobNm := "evp-" + dateFmt
	var jsonStr = []byte(`
	{
		"jobName": "` + jobNm + `",
		"inputs": [{
			"inputBucketName": "` + awsBucketName + `",
			"inputFilePath": "` + objectPath + `"
		}],
		"output": {
			"outputBucketName": "` + awsBucketName + `",
			"outputFilePath": "` + objectOutPath + `",
			"thumbnailOn": "false",
			"thumbnailBucketName": "` + awsBucketName + `",
			"thumbnailFilePath": "` + objectOutPath + `",
			"outputFiles": [{
				"presetId": "` + evp_presetId + `",
				"outputFileName": "` + outputFileNm + `",
				"accessControl": "PUBLIC-READ"
			}]
		}
	}	
	`)
	reqJob, errJob := http.NewRequest("POST", jobCreateUrl, bytes.NewBuffer(jsonStr))
	reqJob.Header.Add("Content-Type", "application/json")
	reqJob.Header.Add("x-ncp-apigw-timestamp", strconv.Itoa(int(timestamp)))
	reqJob.Header.Add("x-ncp-apigw-api-key", apiKey)
	reqJob.Header.Add("x-ncp-iam-access-key", accessKey)
	reqJob.Header.Add("x-ncp-apigw-signature-v1", jobApiSignature)

	resJob, errJob := http.DefaultClient.Do(reqJob)
	if errJob != nil {
		fmt.Println("error : ", errJob)
	}
	defer resJob.Body.Close()
	bodyJob, _ := ioutil.ReadAll(resJob.Body)

	// Job수행 후 결과값 Json Return
	var resultJson EntpVideoResultJson
	err = json.Unmarshal([]byte(bodyJob), &resultJson)
	if err != nil {
		fmt.Println("There was an error:", err)
	}
	fmt.Println("JobId : ", resultJson.Jobs[0].JobId)
	fmt.Println("ErrorCode : ", resultJson.Error.ErrorCode)
	fmt.Println("Message : ", resultJson.Error.Message)

	var resultCD int64

	resultCD = resultJson.Error.ErrorCode
	fmt.Println("resultCD : ", resultCD)

	if resultCD == 0 {
		//go deleteObject(pEnptMemNo, pVdNo, dateFmt, ext)
	}
}

func updateEvp(prevFilePath string) {
	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	//uploadPath, _ := beego.AppConfig.String("uploadpath")

	// Object Storage 파일 물리적 경로
	//objectFilePath := uploadPath + "/entp_video/" + pEnptMemNo + "/" + pVdNo + "/evp_" + dateFmt + ext

	awsAccessKeyID, _ := beego.AppConfig.String("accessKey")  //"6JSYw9NY9JAsO7XynmbO"
	awsSecretAccessKey, _ := beego.AppConfig.String("secretKey")  //"uARPb8klZto1zXTIC4G0TXt0aFrsAM9FkYzlX1vT"
	awsBucketName, _ := beego.AppConfig.String("bucketName") //"qrate-interview-service-test"

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

	_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		fmt.Printf("Failed to delete data to %s/%s, %s\n", *bucket, *key, err.Error())
		return
	} else {
		fmt.Printf("Successfully Delete Object %s and deleted data with key %s\n", *bucket, *key)
	}

	/*
		// 업로드 된 파일 서버에서 삭제처리
		var errDel = os.Remove(objectFilePath)
		if errDel != nil {
			log.Debug("Remove failed: %v", errDel)
		}
	*/
}

func deleteObject(pEnptMemNo string, pVdNo string, dateFmt string, ext string) {
	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	//uploadPath, _ := beego.AppConfig.String("uploadpath")

	// Object Storage 파일 물리적 경로
	//objectFilePath := uploadPath + "/entp_video/" + pEnptMemNo + "/" + pVdNo + "/evp_" + dateFmt + ext

	awsAccessKeyID, _     := beego.AppConfig.String("accessKey")  //"6JSYw9NY9JAsO7XynmbO"
	awsSecretAccessKey, _ := beego.AppConfig.String("secretKey")  //"uARPb8klZto1zXTIC4G0TXt0aFrsAM9FkYzlX1vT"
	awsBucketName, _      := beego.AppConfig.String("bucketName") //"qrate-interview-service-test"	

	objectPath := "/input/entp_video/" + pEnptMemNo + "/" + pVdNo + "/evp_" + dateFmt + ext
	//objectOutPath := "/output/entp_video/" + pEnptMemNo + "/" + pVdNo + "/"

	//outputFileNm := "480p" //dateFmt + ext

	//fmt.Println("objectPath : ", objectPath)
	//fmt.Println("objectOutPath : ", objectOutPath)

	key := aws.String(objectPath)
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

	_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		fmt.Printf("Failed to delete data to %s/%s, %s\n", *bucket, *key, err.Error())
		return
	} else {
		fmt.Printf("Successfully Delete Object %s and deleted data with key %s\n", *bucket, *key)
	}

	/*
		// 업로드 된 파일 서버에서 삭제처리
		var errDel = os.Remove(objectFilePath)
		if errDel != nil {
			log.Debug("Remove failed: %v", errDel)
		}
	*/
}

// Hmac256 Algorythm Encrypt
func VdComputeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
