package controllers

import (
	"encoding/json"
	"fmt"

	"bizt.ziggam.com/models"
	"bizt.ziggam.com/utils"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type ApiApiServerPostReqController struct {
	ApiBaseController
}

func (c *ApiApiServerPostReqController) Prepare() {
	c.ApiBaseController.Prepare()
}

// curl -X POST "http://localhost:7070/api/apisvr/cmd" --data "cmd=v2/member/normal/insert&data={\"sms_recv_yn\": \"1\", \"os_gbn\": \"IS\", \"os_ver\": \"14.4\", \"mem_id\": \"zixzix22\", \"sex\": \"F\", \"brth_ymd\": \"19710101\", \"pwd\": \"abcd1234\", \"mo_no\": \"01012125698\", \"email\": \"ldk22@nate.com\", \"email_recv_yn\": \"1\", \"nm\": \"리리22\"}"

// curl -X POST "https://bizt.ziggam.com/api/apisvr/cmd" --data "cmd=v2/member/normal/insert&data={\"sms_recv_yn\": \"1\", \"os_gbn\": \"IS\", \"os_ver\": \"14.4\", \"mem_id\": \"sungnamtest01\", \"sex\": \"F\", \"brth_ymd\": \"19710101\", \"pwd\": \"sungnamtest01\", \"mo_no\": \"01057905698\", \"email\": \"ldk1010@nate.com\", \"email_recv_yn\": \"1\", \"nm\": \"성남테스트01\"}"

func (c *ApiApiServerPostReqController) Post() {

	pLang, _ := beego.AppConfig.String("lang")

	pCmd := c.GetString("cmd")
	pData := c.GetString("data")

	if pCmd == "" || pData == "" {
		logs.Debug(fmt.Sprintf("Error: pCmd == nil || pData == nil"))
		return
	}

	apiServerUrl, err := beego.AppConfig.String("apiServerUrl")
	if err != nil || apiServerUrl == "" {
		logs.Debug(fmt.Sprintf("Error: %s", "apiServerUrl"))
		return
	}

	//pData := c.GetString("data")
	// 한글 코드 존재시 변환
	//pData = utils.ConvertEucKR(pData)

	logs.Debug(fmt.Sprintf("lang:%s, cmd:%s, data:%s", pLang, pCmd, pData))

	// Start : Oracle DB Connection
	// env, srv, ses, err := GetRawConnection()
	// defer env.Close()
	// defer srv.Close()
	// defer ses.Close()
	// if err != nil {
	// 	panic(err)
	// }
	// End : Oracle DB Connection

	defaultResult := models.DefaultResult{}

	var convData map[string]string

	// 정상 JSON 테스트 용도
	err = json.Unmarshal([]byte(pData), &convData)
	if err != nil {
		fmt.Printf(fmt.Sprintf("Error: %s", err))
		return
	}

	apiServerUrl, err = beego.AppConfig.String("apiServerUrl")
	if err != nil {
		logs.Debug(fmt.Sprintf("Error: %s", err))
		return
	}
	destUrl := apiServerUrl + "/" + pCmd

	// testData := map[string]string{"sms_recv_yn": "1", "os_gbn": "IS", "os_ver": "14.4", "mem_id": "zixzix55", "sex": "F", "brth_ymd": "19710101", "pwd": "abcd1234", "mo_no": "01012125698", "email": "ldk55@nate.com", "email_recv_yn": "1", "nm": "리리55"}
	// beego.Trace(convData)
	// beego.Trace(testData)

	fmt.Printf(fmt.Sprintf("url:%s, data:%s", destUrl, pData))

	ret, err := utils.HttpPostJsonString(destUrl, pData)
	//ret, err := utils.HttpPostJsonMap(destUrl, convData)
	if err != nil {
		fmt.Printf(fmt.Sprintf("Error: %s, Result: %s", err, ret))
		return
	}

	fmt.Printf(fmt.Sprintf("===> rtnCd: %v, rtnMsg: %v", "1", ret))

	defaultResult = models.DefaultResult{
		RtnCd:  1,
		RtnMsg: ret,
	}

	c.Data["json"] = &defaultResult
	c.ServeJSON()
}
