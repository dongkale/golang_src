package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// ret := utils.HttpPostJson("http://localhost:8080/v2/member/normal/insert",
// 	[]byte(`{"sms_recv_yn":"1","os_gbn":"IS","os_ver":"14.4","mem_id":"zixzix88","sex":"F","brth_ymd":"19710101","pwd":"abcd1234","mo_no":"01012125698","email":"ldk88@nate.com","email_recv_yn":"1","nm":"리리88"}`))
func HttpPostJson(url string, toJson []byte) (string, error) {

	// http post Json
	//url := "http://localhost:8080/v2/member/normal/insert"
	fmt.Printf(fmt.Sprintf("Url:%s", url))

	req, err1 := http.NewRequest("POST", url, bytes.NewBuffer(toJson))
	if err1 != nil {
		return "Error 1", fmt.Errorf("http.NewRequest() Error: %v", err1)
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err2 := client.Do(req)
	if err2 != nil {
		return "Error 2", fmt.Errorf("client.Do() Error: %v", err2)
	}
	defer resp.Body.Close()

	body, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		return "Error 3", fmt.Errorf("ioutil.ReadAll() Error: %v", err3)
	}

	fmt.Printf(fmt.Sprintf("Url:%s, Response Status:%s, Response Headers:%s, Response Body:%s", url, resp.Status, resp.Header, string(body)))

	return string(body), nil
}

// values := map[string]string{"sms_recv_yn": "1", "os_gbn": "IS", "os_ver": "14.4", "mem_id": "zixzix55", "sex": "F", "brth_ymd": "19710101", "pwd": "abcd1234", "mo_no": "01012125698", "email": "ldk55@nate.com", "email_recv_yn": "1", "nm": "리리55"}
// ret, err := utils.HttpPostJsonEx("http://localhost:8080/v2/member/normal/insert", values)
// beego.Trace(ret)
// beego.Trace(err)
func HttpPostJsonMap(url string, toJson map[string]string) (string, error) {

	fmt.Printf(fmt.Sprintf("Url:%s", url))

	jsonValue, err1 := json.Marshal(toJson)
	if err1 != nil {
		return "Error 1", fmt.Errorf("json.Marshal() Error: %v", err1)
	}

	resp, err2 := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err2 != nil {
		return "Error 2", fmt.Errorf("http.Post() Error: %v", err2)
	}
	defer resp.Body.Close()

	body, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		return "Error 3", fmt.Errorf("ioutil.ReadAll() Error: %v", err3)
	}

	fmt.Printf(fmt.Sprintf("Url:%s, Response Status:%s, Response Headers:%s, Response Body:%s", url, resp.Status, resp.Header, string(body)))

	return string(body), nil
}

func HttpPostJsonString(url string, toJson string) (string, error) {

	fmt.Printf(fmt.Sprintf("Url:%s", url))

	// jsonValue, err1 := json.Marshal(toJson)
	// if err1 != nil {
	// 	return "Error 1", fmt.Errorf("json.Marshal() Error: %v", err1)
	// }

	resp, err2 := http.Post(url, "application/json", bytes.NewBuffer([]byte(toJson)))
	if err2 != nil {
		return "Error 2", fmt.Errorf("http.Post() Error: %v", err2)
	}
	defer resp.Body.Close()

	body, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		return "Error 3", fmt.Errorf("ioutil.ReadAll() Error: %v", err3)
	}

	fmt.Printf(fmt.Sprintf("Url:%s, Response Status:%s, Response Headers:%s, Response Body:%s", url, resp.Status, resp.Header, string(body)))

	return string(body), nil
}

/*
func SendPostRequest(url string, filename string, filetype string) []byte {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(filetype, filepath.Base(file.Name()))

	if err != nil {
		log.Fatal(err)
	}

	io.Copy(part, file)
	writer.Close()
	request, err := http.NewRequest("POST", url, body)

	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	return content
}
*/

func makeAttachFileRequest2(uri string, params map[string]string, paramsFile, filePath string) (*http.Request, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramsFile, fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return http.NewRequest("POST", uri, body)
}

func makeAttachFileRequest(uri string, params map[string]string, paramsFile, filePath string) (*http.Request, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramsFile, filepath.Base(filePath))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, err
}

// https://gist.github.com/mattetti/5914158/f4d1393d83ebedc682a3c8e7bdc6b49670083b84
// https://matt.aimonetti.net/posts/2013-07-golang-multipart-file-upload-example/
func HttpPostJsonStringAttachFile(url string, toJson string, paramName string, paramFileName string) (string, error) {

	var toJsonMap map[string]string
	err1 := json.Unmarshal([]byte(toJson), &toJsonMap)
	if err1 != nil {
		return "Error 1", fmt.Errorf("json.Unmarshal() Error: %v", err1)
	}

	request, err2 := makeAttachFileRequest(url, toJsonMap, paramName, paramFileName)
	if err2 != nil {
		return "Error 2", fmt.Errorf("makeAttachFileRequest() Error: %v", err2)
	}

	client := &http.Client{}
	resp, err3 := client.Do(request)
	if err3 != nil {
		return "Error 3", fmt.Errorf("client.Do() Error: %v", err3)
	}
	defer resp.Body.Close()

	// body := &bytes.Buffer{}
	// _, err3 := body.ReadFrom(resp.Body)
	// if err3 != nil {
	// 	return "Error 3", fmt.Errorf("body.ReadFrom() Error: %v", err3)
	// }

	body, err4 := ioutil.ReadAll(resp.Body)
	if err4 != nil {
		return "Error 4", fmt.Errorf("ioutil.ReadAll() Error: %v", err4)
	}

	// 결과에 "SetFile":"/resume/P2021021800789/P2021021800789_20210219145346.txt" 이전 화일명, 삭제 용도
	fmt.Printf(fmt.Sprintf("Url:%s, Response Status:%s, Response Headers:%s, Response Body:%s, AttachFile:[%s:%s]", url, resp.Status, resp.Header, string(body), paramName, paramFileName))

	return string(body), nil
}
