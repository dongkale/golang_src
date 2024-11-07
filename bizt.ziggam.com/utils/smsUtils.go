package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"	
)

// {"result_code":"1","message":"success","msg_id":"146645343","success_cnt":1,"error_cnt":0,"msg_type":"SMS"}

// AligoSendSmsResp ...
type AligoSendSmsResp struct {
	ResultCode string `json:"result_code"`
	Message    string `json:"message"`
	MsgID      string `json:"msg_id"`
	SuccessCnt uint32 `json:"success_cnt"`
	ErrorCnt   uint32 `json:"error_cnt"`
	MsgType    string `json:"msg_type"`
}

// IsOk ...
func (resp *AligoSendSmsResp) IsOk() bool {
	return resp.Message == "success" && resp.ResultCode == "1"
}

// AligoSendSmsErrorResp ...
type AligoSendSmsErrorResp struct {
	ResultCode int    `json:"result_code"`
	Message    string `json:"message"`
	MsgID      string `json:"msg_id"`
	SuccessCnt uint32 `json:"success_cnt"`
	ErrorCnt   uint32 `json:"error_cnt"`
	MsgType    string `json:"msg_type"`
}

// {"result_code":1,"message":"success","list":[{"mdid":"3245613052","type":"SMS","sender":"0317391121","receiver":"01052262107","sms_state":"\ubc1c\uc1a1\uc644\ub8cc","reg_date":"2020-10-13 14:29:00","send_date":"2020-10-13 14:29:00","reserve_date":""},{"mdid":"3245613053","type":"SMS","sender":"0317391121","receiver":"01094975702","sms_state":"\ubc1c\uc1a1\uc644\ub8cc","reg_date":"2020-10-13 14:29:00","send_date":"2020-10-13 14:29:00","reserve_date":""}],"next_yn":"N"}

// AligoSendMass ...
type AligoSendMass struct {
	PhoneNum string
	Message  string
}

// AligoSendSmsHisList ...
type AligoSendSmsHisList struct {
	Mdid        string `json:"mdid"`
	Type        string `json:"type"`
	Sender      string `json:"sender"`
	Receiver    string `json:"receiver"`
	SmsState    string `json:"sms_state"`
	RegDate     string `json:"reg_date"`
	SendDate    string `json:"send_date"`
	ReserveDate string `json:"reserve_date"`
}

// AligoSendSmsHisResp ...
type AligoSendSmsHisResp struct {
	ResultCode int                   `json:"result_code"`
	Message    string                `json:"message"`
	List       []AligoSendSmsHisList `json:"list"`
	NextYN     string                `json:"next_yn"`
}

// IsOk ...
func (resp *AligoSendSmsHisResp) IsOk() bool {
	return resp.Message == "success" && resp.ResultCode == 1
}

// {"result_code":1,"message":"success","list":[{"mid":"146719121","type":"SMS","sender":"0317391121","sms_count":"2","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 19:17:11","reserve":""},{"mid":"146718541","type":"SMS","sender":"0317391121","sms_count":"2","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 19:14:50","reserve":""},{"mid":"146717183","type":"SMS","sender":"0317391121","sms_count":"1","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 19:09:40","reserve":""},{"mid":"146708342","type":"SMS","sender":"0317391121","sms_count":"2","reserve_state":"","msg":"API TEST SEND","fail_count":0,"reg_date":"2020-10-13 18:36:33","reserve":""},{"mid":"146706599","type":"SMS","sender":"0317391121","sms_count":"1","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 18:30:19","reserve":""},{"mid":"146706084","type":"SMS","sender":"0317391121","sms_count":"1","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 18:28:41","reserve":""},{"mid":"146706025","type":"SMS","sender":"0317391121","sms_count":"1","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 18:28:23","reserve":""},{"mid":"146706021","type":"SMS","sender":"0317391121","sms_count":"1","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 18:28:23","reserve":""},{"mid":"146703187","type":"SMS","sender":"0317391121","sms_count":"1","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 18:18:27","reserve":""},{"mid":"146697705","type":"SMS","sender":"0317391121","sms_count":"1","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 18:09:40","reserve":""},{"mid":"146696897","type":"SMS","sender":"0317391121","sms_count":"2","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 18:07:05","reserve":""},{"mid":"146696783","type":"SMS","sender":"0317391121","sms_count":"2","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 18:06:40","reserve":""},{"mid":"146695003","type":"SMS","sender":"0317391121","sms_count":"2","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 18:01:41","reserve":""},{"mid":"146694184","type":"SMS","sender":"0317391121","sms_count":"2","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 17:59:52","reserve":""},{"mid":"146694180","type":"SMS","sender":"0317391121","sms_count":"2","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 17:59:50","reserve":""},{"mid":"146693833","type":"SMS","sender":"0317391121","sms_count":"2","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 17:58:32","reserve":""},{"mid":"146693573","type":"SMS","sender":"0317391121","sms_count":"2","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 17:57:43","reserve":""},{"mid":"146692690","type":"SMS","sender":"0317391121","sms_count":"2","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 17:56:08","reserve":""},{"mid":"146692470","type":"SMS","sender":"0317391121","sms_count":"2","reserve_state":"","msg":"API TEST SEND","fail_count":0,"reg_date":"2020-10-13 17:55:30","reserve":""},{"mid":"146692090","type":"SMS","sender":"0317391121","sms_count":"1","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 17:54:27","reserve":""},{"mid":"146689037","type":"SMS","sender":"0317391121","sms_count":"1","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 17:48:07","reserve":""},{"mid":"146688578","type":"SMS","sender":"0317391121","sms_count":"1","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 17:47:08","reserve":""},{"mid":"146687370","type":"SMS","sender":"0317391121","sms_count":"1","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 17:43:30","reserve":""},{"mid":"146686418","type":"SMS","sender":"0317391121","sms_count":"1","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 17:40:22","reserve":""},{"mid":"146685913","type":"SMS","sender":"0317391121","sms_count":"1","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 17:38:46","reserve":""},{"mid":"146654110","type":"SMS","sender":"0317391121","sms_count":"1","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 16:09:25","reserve":""},{"mid":"146651547","type":"SMS","sender":"0317391121","sms_count":"1","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 16:02:42","reserve":""},{"mid":"146645343","type":"SMS","sender":"0317391121","sms_count":"1","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 15:47:21","reserve":""},{"mid":"146643993","type":"SMS","sender":"0317391121","sms_count":"1","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 15:43:24","reserve":""},{"mid":"146643587","type":"SMS","sender":"0317391121","sms_count":"1","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 15:42:21","reserve":""}],"next_yn":"Y"}
// {"mid":"146719121","type":"SMS","sender":"0317391121","sms_count":"2","reserve_state":"","msg":"API TEST","fail_count":0,"reg_date":"2020-10-13 19:17:11","reserve":""}

// AligoSendSmsAllHisList ...
type AligoSendSmsAllHisList struct {
	Mid          string `json:"mid"`
	Type         string `json:"type"`
	Sender       string `json:"sender"`
	SmsCount     string `json:"sms_count"`
	ReserveState string `json:"reserve_state"`
	Msg          string `json:"msg"`
	FailCount    string `json:"fail_count"`
	RegDate      string `json:"reg_date"`
	Reserve      string `json:"reserve"`
}

// AligoSendSmsAllHisResp ...
type AligoSendSmsAllHisResp struct {
	ResultCode int                      `json:"result_code"`
	Message    string                   `json:"message"`
	List       []AligoSendSmsAllHisList `json:"list"`
	NextYN     string                   `json:"next_yn"`
}

// AligoApiKey  Key
const AligoApiKey string = "gd4cj2r68bj7fj0o25cnx4xzpolbdcuo"

// AligoApiId Api Account
const AligoApiId string = "qrateziggam"

// AligoSenderMono Api sender Phone No
const AligoSenderMono string = "0317391121"

// AligoSendSms ...
func AligoSendSms(memList []string, message string, title string) AligoSendSmsResp {

	const aligoURL string = "https://apis.aligo.in/send/"

	var destList string
	destList = strings.Join(memList, ",")

	payload := strings.NewReader(fmt.Sprintf("key=%v&user_id=%v&sender=%v&receiver=%v&msg=%v&title=%v",
		AligoApiKey,                // key
		AligoApiId,                 // user_id
		AligoSenderMono,            // sender
		strings.Join(memList, ","), // receiver
		message,                    // msg
		title))                     // title

	req, _ := http.NewRequest("POST", aligoURL, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Add("cache-control", "no-cache")

	resp, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	//respString := string(body)

	//var respData AligoSendSmsResp
	//err1 := json.Unmarshal([]byte(respString), &respData)

	var respData AligoSendSmsResp
	err := json.Unmarshal(body, &respData)
	// if err != nil {
	// 	panic(err)
	// }

	var respDataError AligoSendSmsErrorResp
	if len(respData.ResultCode) <= 0 {
		err = json.Unmarshal(body, &respDataError)
		// if err != nil {
		// 	panic(err)
		// }

		respData.ResultCode = strconv.Itoa(respDataError.ResultCode)
	}

	// var respData AligoSendSmsResp
	// respData.ResultCode = data["result_code"].(string)
	// respData.Message = data["message"].(string)
	// //if respData.ResultCode == 1.0 {
	// respData.MsgID = data["msg_id"].(string)
	// respData.SuccessCnt = data["success_cnt"].(int)
	// respData.ErrorCnt = data["error_cnt"].(int)
	// respData.MsgType = data["msg_type"].(string)
	// //}

	jsonBytes, err := json.Marshal(respData)
	if err != nil {
		panic(err)
	}

	fmt.Printf(fmt.Sprintf("[AligoSendSms][Result] Result:%s, Result:%s, Error:%v", string(jsonBytes), string(body), err))
	fmt.Printf(fmt.Sprintf("[AligoSendSms][Result] List:%s", destList))

	return respData
}

// AligoSendMassSms ...
func AligoSendMassSms(massList []AligoSendMass, title string, msgType string) AligoSendSmsResp {

	/*
		key				인증용 API Key						O	String
		user_id			사용자id							O	String
		sender			발신자 전화번호 (최대 16bytes)	      O	String

		rec_1			수신자 전화번호1					    O	String
		msg_1			메시지 내용1						   O	String (1~2,000Byte)
		rec_2 ~ rec_500	수신자 전화번호2 ~ 500				    X	String
		msg_2 ~ msg_500	메시지 내용							    X	String (1~2,000Byte)
		cnt				메세지 전송건수(번호,메세지 매칭건수)	 	O	Integer(1~500)
		title			문자제목(LMS,MMS만 허용)				 X	String (1~44Byte)
		msg_type		SMS(단문) , LMS(장문), MMS(그림문자) 구분 O	String
		rdate			예약일 (현재일이상)						 X	YYYYMMDD
		rtime			예약시간 - 현재시간기준 10분이후		   X	HHII
		image			첨부이미지								X	JPEG,PNG,GIF
	*/

	const aligoURL string = "https://apis.aligo.in/send_mass/"

	var destList string

	for index, val := range massList {
		destList += fmt.Sprintf("&rec_%d=%s&msg_%d=%s",
			index+1, val.PhoneNum,
			index+1, val.Message)
	}

	//var msgType string = "SMS"

	var tempString = fmt.Sprintf("key=%s&user_id=%s&sender=%s&title=%s&cnt=%d&msg_type=%s%s",
		AligoApiKey,     // key
		AligoApiId,      // user_id
		AligoSenderMono, // sender
		title,           // Title
		len(massList),
		msgType,
		destList)

	payload := strings.NewReader(tempString)

	req, _ := http.NewRequest("POST", aligoURL, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Add("cache-control", "no-cache")

	resp, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	//respString := string(body)

	//var respData AligoSendSmsResp
	//err1 := json.Unmarshal([]byte(respString), &respData)

	var respData AligoSendSmsResp
	err := json.Unmarshal(body, &respData)
	// if err != nil {
	// 	panic(err)
	// }

	var respDataError AligoSendSmsErrorResp
	if len(respData.ResultCode) <= 0 {
		err = json.Unmarshal(body, &respDataError)
		// if err != nil {
		// 	panic(err)
		// }

		respData.ResultCode = strconv.Itoa(respDataError.ResultCode)
		respData.MsgID = "0"
	}

	// var respData AligoSendSmsResp
	// respData.ResultCode = data["result_code"].(string)
	// respData.Message = data["message"].(string)
	// //if respData.ResultCode == 1.0 {
	// respData.MsgID = data["msg_id"].(string)
	// respData.SuccessCnt = data["success_cnt"].(int)
	// respData.ErrorCnt = data["error_cnt"].(int)
	// respData.MsgType = data["msg_type"].(string)
	// //}

	jsonBytes, err := json.Marshal(respData)
	if err != nil {
		panic(err)
	}

	fmt.Printf(fmt.Sprintf("[AligoSendSms][Result] Result:%s, Result:%s, Error:%v", string(jsonBytes), string(body), err))
	fmt.Printf(fmt.Sprintf("[AligoSendSms][Result] List:%s", destList))

	return respData
}

// AligoSendSmsHis ...
// 146717183
// 146718541
// 146719121
func AligoSendSmsHis(mid string, page int, pageSize int) AligoSendSmsHisResp {

	const aligoURL string = "https://apis.aligo.in/sms_list/"

	payload := strings.NewReader(fmt.Sprintf("key=%v&user_id=%v&mid=%v&page=%v&page_size=%v",
		AligoApiKey, // key
		AligoApiId,  // user_id
		mid,         // 메세지 고유ID
		page,        // 페이지번호
		pageSize))   // 페이지당 출력갯수

	req, _ := http.NewRequest("POST", aligoURL, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Add("cache-control", "no-cache")

	resp, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var respData AligoSendSmsHisResp
	err := json.Unmarshal(body, &respData)
	//if err != nil {
	//	panic(err)
	//}

	// respString := string(body)

	// var data map[string]interface{}
	// err := json.Unmarshal(body, &data)

	// fmt.Printf(data["result_code"])
	// fmt.Printf(data["message"])

	// 146719121
	//[AligoSendSms][Result] Result:{"result_code":1,"message":"success","list":[{"mdid":"3246554255","type":"SMS","sender":"0317391121","receiver":"01052262107","sms_state":"발송완료","reg_date":"2020-10-13 19:17:11","send_date":"2020-10-13 19:17:11","reserve_date":""},{"mdid":"3246554256","type":"SMS","sender":"0317391121","receiver":"010904975702","sms_state":"발/착신 번호 에러","reg_date":"2020-10-13 19:17:11","send_date":"2020-10-13 19:17:11","reserve_date":""}],"next_yn":"N"}({"result_code":1,"message":"success","list":[{"mdid":"3246554255","type":"SMS","sender":"0317391121","receiver":"01052262107","sms_state":"\ubc1c\uc1a1\uc644\ub8cc","reg_date":"2020-10-13 19:17:11","send_date":"2020-10-13 19:17:11","reserve_date":""},{"mdid":"3246554256","type":"SMS","sender":"0317391121","receiver":"010904975702","sms_state":"\ubc1c\/\ucc29\uc2e0 \ubc88\ud638 \uc5d0\ub7ec","reg_date":"2020-10-13 19:17:11","send_date":"2020-10-13 19:17:11","reserve_date":""}],"next_yn":"N"}), Error:<nil>

	jsonBytes, err := json.Marshal(respData)
	if err != nil {
		panic(err)
	}

	fmt.Printf(fmt.Sprintf("[AligoSendSmsHis][Result] Result:%s, Result:%s, Error:%v", string(jsonBytes), string(body), err))

	return respData
}

// AligoSendSmsAllHis ...
func AligoSendSmsAllHis(page int, pageSize int, startDate string, limitDay int) AligoSendSmsAllHisResp {

	const aligoURL string = "https://apis.aligo.in/list/"

	payload := strings.NewReader(fmt.Sprintf("key=%v&user_id=%v&page=%v&page_size=%v&start_date=%v&limit_day=%v",
		AligoApiKey, // key
		AligoApiId,  // user_id
		page,        // 페이지번호
		pageSize,    // 페이지당 출력갯수
		startDate,
		limitDay))

	req, _ := http.NewRequest("POST", aligoURL, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Add("cache-control", "no-cache")

	resp, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var respData AligoSendSmsAllHisResp
	err := json.Unmarshal(body, &respData)
	// if err != nil {
	// 	panic(err)
	// }

	jsonBytes, err := json.Marshal(respData)
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Printf(fmt.Sprintf("[AligoSendSmsAllHisResp][Result] Result:%s, Result:%s, Error:%v", string(jsonBytes), string(body), err))
	// fmt.Printf(fmt.Sprintf("[AligoSendSmsAllHisResp][Result] Result1:%s, Error:%v", string(jsonBytes), err))
	// fmt.Printf(fmt.Sprintf("[AligoSendSmsAllHisResp][Result] Result2:%s, Error:%v", string(body), err))

	return respData
}
