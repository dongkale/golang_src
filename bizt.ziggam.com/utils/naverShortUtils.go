package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"	
)

// Naver Short Url
// {"result":{"url":"http://me2.do/GyvykVAu","hash":"GyvykVAu","orgUrl":"http://d2.naver.com/helloworld/4874130"},
//	"message":"ok",
//	"code":"200"}

//  Result:{"errorMessage":"Not Exist Client ID : Authentication failed. (인증에 실패했습니다.)","errorCode":"024"}

// NaverShortUrl ...
type NaverShortUrl struct {
	Url    string `json:"url"`
	Hash   string `json:"hash"`
	OrgUrl string `json:"orgUrl"`
}

// NaverShortUrlResp ...
type NaverShortUrlResp struct {
	Result  NaverShortUrl `json:"result"`
	Message string        `json:"message"`
	Code    string        `json:"code"`
}

func (resp *NaverShortUrlResp) IsOk() bool {
	return resp.Message == "ok" && resp.Code == "200"
}

// NaverShortUrlError ...
// "errorMessage":"Not Exist Client ID : Authentication failed. (인증에 실패했습니다.)","errorCode":"024"
type NaverShortUrlErrorResp struct {
	ErrorMsg  string `json:"errorMessage"`
	ErrorCode string `json:"errorCode"`
}

// NaverShortUrlApiKey ...
const NaverShortUrlApiKey string = "Y9tRnDMGWC"

// NaverShortUrlApiId Api Account
const NaverShortUrlApiId string = "q4CGRH8KUn4V1se1373V"

// NaverShortUrlReq ...
func NaverShortUrlReq(reqUrl string) NaverShortUrlResp {

	const naverShortURL string = "https://openapi.naver.com/v1/util/shorturl"

	//curl "https://openapi.naver.com/v1/util/shorturl"
	// 	   -d "url=http://d2.naver.com/helloworld/4874130"
	//     -H "Content-Type: application/x-www-form-urlencoded; charset=UTF-8"
	//	   -H "X-Naver-Client-Id: q4CGRH8KUn4V1se1373V"
	//	   -H "X-Naver-Client-Secret: Y9tRnDMGWC" -v
	// {"result":{"url":"http://me2.do/GyvykVAu","hash":"GyvykVAu","orgUrl":"http://d2.naver.com/helloworld/4874130"},"message":"ok","code":"200"}

	payload := strings.NewReader(fmt.Sprintf("url=%v", reqUrl))

	req, _ := http.NewRequest("POST", naverShortURL, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("X-Naver-Client-Id", NaverShortUrlApiId)
	req.Header.Add("X-Naver-Client-Secret", NaverShortUrlApiKey)

	resp, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var respData NaverShortUrlResp
	err := json.Unmarshal(body, &respData)
	if err != nil {
		panic(err)
	}

	var respDataError NaverShortUrlErrorResp
	if len(respData.Code) <= 0 {
		err = json.Unmarshal(body, &respDataError)
		if err != nil {
			panic(err)
		}

		respData.Code = respDataError.ErrorCode
		respData.Message = respDataError.ErrorMsg

		fmt.Printf(fmt.Sprintf("[NaverShortUrl][Result] Result:%s, Error:%v, Body:%v", respDataError.ErrorMsg, respDataError.ErrorCode, string(body)))
	}

	jsonBytes, err := json.Marshal(respData)
	if err != nil {
		panic(err)
	}

	fmt.Printf(fmt.Sprintf("[NaverShortUrl][Result] Result:%s, Result:%s, Error:%v", string(jsonBytes), string(body), err))

	return respData
}
