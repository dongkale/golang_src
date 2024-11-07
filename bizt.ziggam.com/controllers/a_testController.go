package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"bizt.ziggam.com/models"
	"bizt.ziggam.com/utils"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	ora "gopkg.in/rana/ora.v4"
)

type AligoSendSmsResp struct {
	result_code string
	message     string
	msg_id      string
	success_cnt int
	error_cnt   int
	msg_type    string
}

// A_TestController
type A_TestController struct {
	BaseController
}

// Get()
func (c *A_TestController) Get() {

	// start : log
	//log := logs.NewLogger()
	//log.SetLogger(logs.AdapterConsole)
	// end : log

	// var (
	// 	mnnSn    int64
	// 	mnnTitle string
	// 	mnnRegDt string
	// )

	//session := c.StartSession()

	fmt.Println("Start...")

	c.Rerender()

	//ret := mem_sn_sn.(string)

	//fmt.Printf(ret)

	//c.Abort("DataBase")

	// HTTP Request...
	/*
		// https://jeonghwan-kim.github.io/dev/2019/02/07/go-net-http.html
		url := "https://google.com/robots.txt"

		resp, _ := http.Get(url)
		robots, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		fmt.Printf("%s\n", robots)
	*/
	// HTTP Request...

	// https://www.codershood.info/2017/06/25/http-curl-request-golang/
	/*
		url := "https://google.com/robots.txt"

		req, _ := http.NewRequest("GET", url, nil)

		//req.Header.Add("cache-control", "no-cache")
		res, _ := http.DefaultClient.Do(req)

		res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		fmt.Printf("%s\n", body)
	*/

	// GET with header
	/*
		url := "https://google.com/robots.txt"

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("set-cookie", "foo=bar") // 헤더값 설정

		resp, _ := http.DefaultClient.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		fmt.Printf("%s\n", body)
	*/

	// POST with header
	/*
		url := "https://reqres.in/api/users"

		payload := strings.NewReader("name=test&jab=teacher")

		req, _ := http.NewRequest("POST", url, payload)

		req.Header.Add("content-type", "application/x-www-form-urlencoded")
		req.Header.Add("cache-control", "no-cache")

		resp, _ := http.DefaultClient.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		fmt.Printf("%s\n", body)
	*/

	//url
	/*
		easy := curl.EasyInit()
		defer easy.Cleanup()

		easy.Setopt(curl.OPT_URL, "http://www.baidu.com/")

		// make a callback function
		fooTest := func(buf []byte, userdata interface{}) bool {
			println("DEBUG: size=>", len(buf))
			println("DEBUG: content=>", string(buf))
			return true
		}

		easy.Setopt(curl.OPT_WRITEFUNCTION, fooTest)

		if err := easy.Perform(); err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
	*/
	//url

	// aligo
	/*
		url := "https://apis.aligo.in/send/"

		//payload := strings.NewReader("key=gd4cj2r68bj7fj0o25cnx4xzpolbdcuo&user_id=qrateziggam&sender=01052262107&receiver=01052262107,01052262107&destination=01052262107|Lee,01052262107|Lee&msg=API TEST&title=API TEST&rdate=20201013&rtime=1210&testmode_yn=Y")

		//payload := strings.NewReader(fmt.Sprintf("key=%v&user_id=%v&sender=%v&receiver=%v&destination=%v&msg=%v&title=%v&rdate=%v&rtime=%v&testmode_yn=%v",
		payload := strings.NewReader(fmt.Sprintf("key=%v&user_id=%v&sender=%v&receiver=%v&msg=%v&title=%v&testmode_yn=%v",
			"gd4cj2r68bj7fj0o25cnx4xzpolbdcuo", // key
			"qrateziggam",                      // user_id
			"0317391121",                       // sender
			"01052262107,01052262107",          // receiver
			"API TEST",                         // msg
			"API TEST(title)",                  // title
			"Y"))                               // testmode_yn

		req, _ := http.NewRequest("POST", url, payload)

		req.Header.Add("content-type", "application/x-www-form-urlencoded;charset=UTF-8")
		req.Header.Add("cache-control", "no-cache")

		// curl -X POST "https://apis.aligo.in/send/" --data-urlencode "key=gd4cj2r68bj7fj0o25cnx4xzpolbdcuo" --data-urlencode "user_id=qrateziggam" --data-urlencode "sender=0000000000" --data-urlencode "receiver=01052262107,01052262107" --data-urlencode "destination=01052262107|Lee,01052262107|Lee" --data-urlencode "msg=API TEST SEND" --data-urlencode "title=API TEST 입니다" --data-urlencode "rdate=20201013" --data-urlencode "rtime=1210" --data-urlencode "testmode_yn=Y" --data-urlencode "charset=utf-8"
		// req.Header.Add("key", "gd4cj2r68bj7fj0o25cnx4xzpolbdcuo")
		// req.Header.Add("user_id", "qrateziggam")
		// req.Header.Add("sender", "01052262107")
		// req.Header.Add("receiver", "01052262107,01052262107")
		// req.Header.Add("destination", "01052262107|Lee,01052262107|Lee")
		// req.Header.Add("msg", "API TEST SEND")
		// req.Header.Add("title", "API TEST")
		// req.Header.Add("rdate", "20201013")
		// req.Header.Add("rtime", "1210")
		// req.Header.Add("testmode_yn", "Y")
		//req.Header.Add("charset", "UTF-8")

		resp, _ := http.DefaultClient.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		fmt.Printf("%s\n", body)
	*/

	/*
		var respString = "{\"result_code\":\"1\",\"message\":\"success\",\"msg_id\":\"146654110\",\"success_cnt\":1,\"error_cnt\":0,\"msg_type\":\"SMS\"}"

		var data map[string]interface{}
		err := json.Unmarshal([]byte(respString), &data)
		if err != nil {
			panic(err)
		}

		var respData utils.AligoSendSmsResp
		err2 := json.Unmarshal([]byte(respString), &respData)
		if err2 != nil {
			panic(err2)
		}

		fmt.Printf("[AligoSendSms][Result] %v", respData)
		fmt.Printf(fmt.Sprintf("[AligoSendSms][Result] %v", data["result_code"]))
		fmt.Printf(fmt.Sprintf("[AligoSendSms][Result] %v", data["message"]))

		beego.Trace(err)
		beego.Trace(err2)
	*/

	//memList2 := []string{"01052262107,010904975702"}

	/*
		var memList []string
		//memList := make([]string)

		memList = append(memList, "01052262107")
		memList = append(memList, "010904975702")
	*/

	/*
		/// https://ko.wikipedia.org/wiki/KS_X_1001%EC%9D%98_%ED%8A%B9%EC%88%98_%EB%AC%B8%EC%9E%90
		var MsgFmt = fmt.Sprintf("%v 님, 안녕하세요.\\n"+
			"%v 에서 직감을 통한 영상 지원을 요청하셨습니다.\\n"+
			"\\n"+
			"√ ¶ 하단 링크를 통해 앱 설치 후 자세한 채용공고 내용을 확인하신 후 지원하실 수 있습니다.\\n"+
			"√ ¶ 채용공고 마감일 이후에는 지원이 불가하니 유의해주세요.\\n"+
			"\\n"+
			"▶ 기업명: %v\\n"+
			"▶ 채용공고: %v\\n"+
			"▶ 바로가기: %v\\n"+
			"\\n"+
			"* 직감이란? 직감은 영상 기반의 채용 서비스를 제공하는 플랫폼입니다. 모바일 어플리케이션을 다운로드하여 진행하실 수 있습니다.\\n"+
			"직감 영상 인터뷰 가이드 바로가기: %v\\n"+
			"* 영상 인터뷰 관련 문의 : support@ziggam.com\\n"+
			"* 수신자 정보는 30일 후 자동 삭제됩니다. 즉시 삭제를 원하실 경우 support@ziggam.com을 통해 문의해주세요.",
			"이동&관",
			"큐레잇",
			"큐레잇",
			"채용공고",
			"https://localhost:7070?ddd=10&&dddd=10",
			"https://www.notion.so/qrate/534a4dc9166d4baebcf701fc6a8c392c?v=35c06f627823428a93e4df954cdeac98")

		//massList := make([]utils.AligoSendMass, 0)
		massList := []utils.AligoSendMass{}

		massList = append(massList, utils.AligoSendMass{
			PhoneNum: "01052262107",
			Message:  MsgFmt, //"jfTitle1\\njfTitle2\\njfTitle3",
		})

		massList = append(massList, utils.AligoSendMass{
			PhoneNum: "010904975702",
			Message:  "jfTitle1\\njfTitle2\\njfTitle3",
		})

		var respData0 utils.AligoSendSmsResp

		respData0 = utils.AligoSendMassSms(massList, "[기업명] 직감을 통해 채용공고 지원을 요청드립니다.", "LMS")

		beego.Trace(respData0)
	*/

	// const testURL string = "localhost:8080"

	// payload := strings.NewReader(fmt.Sprintf("nm=%s", "한글")) // title

	// req, _ := http.NewRequest("POST", testURL, payload)

	// req.Header.Add("content-type", "application/x-www-form-urlencoded;charset=UTF-8")
	// req.Header.Add("cache-control", "no-cache")

	// resp, _ := http.DefaultClient.Do(req)
	// body, _ := ioutil.ReadAll(resp.Body)
	// resp.Body.Close()

	// beego.Trace(body)

	// resp, err := http.PostForm("http://localhost:8080/v2/member/normal/insert", url.Values{"nm": {"한글"}})
	// if err != nil {
	// 	panic(err)
	// }

	// defer resp.Body.Close()

	// // Response 체크.
	// respBody, err := ioutil.ReadAll(resp.Body)
	// if err == nil {
	// 	str := string(respBody)
	// 	println(str)
	// }

	// http post Json
	// url := "http://localhost:8080/v2/member/normal/insert"
	// fmt.Println("URL:>", url)

	// var jsonStr = []byte(`{"sms_recv_yn":"1","os_gbn":"IS","os_ver":"14.4","mem_id":"zixzix77","sex":"F","brth_ymd":"19710101","pwd":"abcd1234","mo_no":"01012125698","email":"ldk77@nate.com","email_recv_yn":"1","nm":"리리77"}`)

	// req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	// req.Header.Set("X-Custom-Header", "myvalue")
	// req.Header.Set("Content-Type", "application/json")

	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))

	// ret := utils.HttpPostJson("http://localhost:8080/v2/member/normal/insert",
	// 	[]byte(`{"sms_recv_yn":"1","os_gbn":"IS","os_ver":"14.4","mem_id":"zixzix88","sex":"F","brth_ymd":"19710101","pwd":"abcd1234","mo_no":"01012125698","email":"ldk88@nate.com","email_recv_yn":"1","nm":"리리88"}`))

	// fmt.Println(ret)

	//person := []byte(`{"sms_recv_yn":"1","os_gbn":"IS","os_ver":"14.4","mem_id":"zixzix88","sex":"F","brth_ymd":"19710101","pwd":"abcd1234","mo_no":"01012125698","email":"ldk88@nate.com","email_recv_yn":"1","nm":"리리88"}`))
	//strp := `{"sms_recv_yn":"1","os_gbn":"IS","os_ver":"14.4","mem_id":"zixzix99","sex":"F","brth_ymd":"19710101","pwd":"abcd1234","mo_no":"01012125698","email":"ldk99@nate.com","email_recv_yn":"1","nm":"리리99"}`

	//values := map[string]string{"sms_recv_yn": "1", "os_gbn": "IS", "os_ver": "14.4", "mem_id": "zixzix55", "sex": "F", "brth_ymd": "19710101", "pwd": "abcd1234", "mo_no": "01012125698", "email": "ldk55@nate.com", "email_recv_yn": "1", "nm": "리리55"}

	// jsonValue, _ := json.Marshal(values)
	// //buff := bytes.NewBuffer(pbytes)
	// resp, err := http.Post("http://localhost:8080/v2/member/normal/insert", "application/json", bytes.NewBuffer(jsonValue))
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))

	// beego.Trace(resp)
	// beego.Trace(err)

	// ret, err := utils.HttpPostJsonMap("http://localhost:8080/v2/member/normal/insert", values)
	// beego.Trace(ret)
	// beego.Trace(err)

	//values := map[string]string{"my_mem_no": "P2021021800789"}
	//ret, err := utils.HttpPostJsonAttachFile("http://localhost:8080/v2/join/resume/insert", values, "attachFile", "aaa.txt")

	// values := `{"my_mem_no": "P2021021800789"}`

	// ret, err := utils.HttpPostJsonStringAttachFile("http://localhost:8080/v2/join/resume/insert", values, "attachFile", "aaa.txt")
	// beego.Trace(ret)
	// beego.Trace(err)

	/*
		memList := make([]string, 2)

		memList[0] = "01052262107"
		memList[1] = "010904975702"

		var respData1 utils.AligoSendSmsResp

		respData1 = utils.AligoSendSms(memList, "API TEST", "API TEST(Title)")

		beego.Trace(memList)
		//beego.Trace(memList2)
		beego.Trace(respData1)
	*/

	// var respData2 utils.AligoSendSmsHisResp

	// respData2 = utils.AligoSendSmsHis("151566693", 1, 30)

	// beego.Trace(respData2)

	// for _, val := range respData2.List {
	// 	fmt.Println(val)
	// 	fmt.Println(val.SmsState)
	// }

	/*
		var respData3 utils.AligoSendSmsAllHisResp

		respData3 = utils.AligoSendSmsAllHis(1, 30, "20201013", 1)

		beego.Trace(respData3)
	*/

	// for _, val := range respData3.List {
	// 	fmt.Println(val)
	// 	fmt.Println(val.Mid)
	// }

	//utils.SendMail("dongkale@naver.com", "no-reply@ziggam.com", "subject", "htmlContents")

	// err := utils.SendMailEx("dongkale@jk.com", "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")
	// err := utils.SendMail("dongkale@jk.com", "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")

	//beego.Trace(err)

	// err0 := utils.SendMail("dongkale@naver.com", "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")
	// err1 := utils.SendMailEx("dongkale@naver.com", "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")
	// err2 := utils.SendMailExEx("dongkale@naver.com", "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")

	// for i := 0; i < 100; i++ {
	// 	// err1 := utils.SendMailEx(fmt.Sprintf("fghjklasd%d@naver.com", i), "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")
	// 	// beego.Trace(fmt.Sprintf("%v -> %v", fmt.Sprintf("fghjklasd%d@naver.com", i), err1))

	// }

	go TestSendEmail()

	//utils.MailDBPoolMng.ReInit(30, 20)

	// utils.MailSender.Connect(utils.GetSmtpServer(), utils.GetSmtpServerPort(), utils.GetReturnEmail(), utils.GetReturnEmailPwd())

	// err := utils.MailSender.Send("dongkale@naver.com", "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")
	// beego.Trace(err)

	// err = utils.MailSender.Send("dongkale@qrate.co.kr", "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")
	// beego.Trace(err)

	// for i := 0; i < 100; i++ {
	// 	err1 := utils.SendMailDaemon("dongkale@naver.com", "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")
	// 	beego.Trace(fmt.Sprintf("%v -> %v", fmt.Sprintf("fghjklasd%d@naver.com", i), err1))
	// }

	//err3 := utils.SendMailDaemon("dongkale@naver.com", "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")

	// err3 = utils.SendMailExExEx("dongkale@naver.com", "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")
	// err3 = utils.SendMailExExEx("dongkale@naver.com", "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")
	// err3 = utils.SendMailExExEx("dongkale@naver.com", "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")
	// err3 = utils.SendMailExExEx("dongkale@naver.com", "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")

	// beego.Trace(err0)
	// beego.Trace(err1)
	// beego.Trace(err2)
	//beego.Trace(err3)

	/*
		-- 네이버 단축 URL
		curl "https://openapi.naver.com/v1/util/shorturl" -d "url=http://d2.naver.com/helloworld/4874130" -H "Content-Type: application/x-www-form-urlencoded; charset=UTF-8" -H "X-Naver-Client-Id: q4CGRH8KUn4V1se1373V" -H "X-Naver-Client-Secret: Y9tRnDMGWC" -v
	*/

	/*
		var respData3 utils.NaverShortUrlResp

		respData3 = utils.NaverShortUrlReq("https://bizt.ziggam.com")

		beego.Trace(respData3.IsOk())

		beego.Trace(respData3)
	*/

	// imgServer, _  := beego.AppConfig.String("viewpath")

	// googleStore := beego.AppConfig.String("googlestore")
	// appleStore := beego.AppConfig.String("applestore")
	// mailTo := beego.AppConfig.String("mailto")

	// siteUrl := beego.AppConfig.String("siteurl")

	// bridgeUrl := beego.AppConfig.String("bridgeUrl")

	// var retSendData []models.InviteMember

	// retSendData = append(retSendData, models.InviteMember{
	// 	Name:  "이동관",
	// 	Email: "dongkale@naver.com",
	// 	Phone: "010-5226-2107",
	// })

	// retSendData = append(retSendData, models.InviteMember{
	// 	Name:  "이동관2",
	// 	Email: "dongkale@qrate.co.kr",
	// 	Phone: "010-5226-2107",
	// })

	// retSendData = append(retSendData, models.InviteMember{
	// 	Name:  "이수진",
	// 	Email: "ojustone@naver.com",
	// 	Phone: "010-5021-7520",
	// })

	// TestInviteSendEmailAll(bridgeUrl, imgServer, mailTo, "support@ziggam.com", googleStore, appleStore, siteUrl, "entpKoNm", "title", "msg1\\nmsg2\\nmsg3\\n", "recruitTitle", "upJobGrp", "jobGrp", "resultEntpVdUrl", retSendData)

	//InviteSendEmailAll(bridgeUrl string, imgServer string, mailTo string, supportMail string, googleStore string, appleStore string, siteUrl string, sendTime string, entpMemNo string, entpKoNm string, title string, msg string, recruitSn string, recruitTitle string, upJobGrp string, jobGrp string, entpVdUrl string, sendList []models.InviteMember) {

	// utils.SlackSend(beego.AppConfig.String("runmode"), "PreText1", "[직감]오류 발생1", "Message1\nMessage2\nMessage3\n")
	// utils.SlackSend(beego.AppConfig.String("runmode"), "PreText2", "[직감]오류 발생2", "Message1\nMessage2\nMessage3\n")

	//utils.SendMail("dongkale@naver.com", "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")

	c.Data["TMenuId"] = "M00"
	c.Data["SMenuId"] = "M00"

	models.Room1.AuthCd++

	rrr1 := models.RoomTest{}

	fmt.Printf(models.TestID["Created"])
	//rrr2 := models.Room1
	//rrr2.AuthCd++

	//beego.RenderForm()

	//beego.BuildTemplate(beego.ViewsPath, "index.tpl")

	beego.BuildTemplate("")

	logs.Info(rrr1)
	//beego.Info(rrr2)
	logs.Info(strconv.Itoa(models.Room1.AuthCd))

	c.Data["MnnTitle"] = strconv.Itoa(models.Room1.AuthCd) //"..........................."

	// 제대로 작동이 안됨
	v1 := c.GetSession("SAVE1")
	if v1 == nil {
		c.Data["SAVE1"] = "AAA"
	} else {
		c.Data["SAVE1"] = "FFF"
	}

	v2 := c.GetSession("SAVE2")
	if v2 == nil {
		c.SetSession("SAVE2", "ABC")
	} else {
		c.SetSession("SAVE2", "CBA")
	}

	v3 := c.GetSession("SAVE3")
	if v3 == nil {
		c.SetSession("SAVE3", int(0))
	} else {
		c.SetSession("SAVE3", v3.(int)+1)
	}

	fmt.Printf("%v", c.GetSession("SAVE1"))
	fmt.Printf("%v", c.GetSession("SAVE2"))
	fmt.Printf("%v", c.GetSession("SAVE3"))

	c.TplName = "ldk_test/a_test.html"

	//c.TplName = "ldk_test/remote_monster.html"

	//panic("a problem")
}

// Post()
func (c *A_TestController) Post() {

	session := c.StartSession()

	mem_id := session.Get(c.Ctx.Request.Context(), "mem_id")
	if mem_id != nil {
		c.Ctx.Redirect(302, "/")
	}

	pMemId := c.GetString("mem_id")
	pPwd := c.GetString("pwd")

	fmt.Printf(pMemId)
	fmt.Printf(pPwd)
}

func (c *A_TestController) Rerender() {

	beego.BuildTemplate(beego.BConfig.WebConfig.ViewsPath)
	//c.Redirect("/", 302)
}

func TestInviteSendEmailAll(bridgeUrl string, imgServer string, mailTo string, supportMail string, googleStore string, appleStore string, siteUrl string, entpKoNm string, title string, msg string, recruitTitle string, upJobGrp string, jobGrp string, entpVdUrl string, sendList []models.InviteMember) {

	fmt.Printf(bridgeUrl)
	fmt.Printf(imgServer)
	fmt.Printf(mailTo)
	fmt.Printf(supportMail)
	fmt.Printf(googleStore)
	fmt.Printf(appleStore)
	fmt.Printf(entpKoNm)
	fmt.Printf(title)
	fmt.Printf(msg)
	fmt.Printf(recruitTitle)
	fmt.Printf(upJobGrp)
	fmt.Printf(jobGrp)
	fmt.Printf(entpVdUrl)

	for _, val := range sendList {

		var resultRecruitUrl string

		var recruitUrl = bridgeUrl + "?" //+ url.QueryEscape(fmt.Sprintf("entpmemno=%v&recruitsn=%v&reqname=%v&reqemail=%v&reqmono=%v", entpMemNo, recruitSn, val.Name, val.Email, val.Phone))

		var respData utils.NaverShortUrlResp
		respData = utils.NaverShortUrlReq(recruitUrl)
		if respData.IsOk() == false {

			resultRecruitUrl = recruitUrl
		} else {
			resultRecruitUrl = respData.Result.Url
		}

		var convMsg = strings.Replace(msg, "{지원자명}", val.Name, 1)
		convMsg = strings.Replace(convMsg, "\\n", "<br>", 10)

		fmt.Printf(resultRecruitUrl)

		/*
			var htmlString = `
				<html>
				<head>
				<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
				<meta http-equiv="X-UA-Compatible" content="ie=edge">
				<title>직감 채용을 편하게, 면접을 영상으로</title>
				</head>
				<body style="margin:0;padding:0;background-color:#f5f6f9;">
				<table border="0" cellspacing="0" cellpadding="0" align="center" style="font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif">
					<!-- header -->
					<tr>
						<td style="padding:26px 0px 20px;border-bottom:4px solid #2ad0c7;">
							<table border="0" cellspacing="0" cellpadding="0" style="width:100%;">
								<tr>
									<td style="float:left">
										<a href="https://www.ziggam.com" target="_blank"><img style="display:block;border:0;" src="` + imgServer + `/mail/ic-main-logo.png" alt="직감 채용을 편하게,면접을 영상으로"></a>
									</td>
								</tr>
							</table>
						</td>
					</tr>
					<!-- //header -->
					<!-- contents -->
					<tr>
						<td style="padding:40px 70px;background-color:#ffffff;">
							<table border="0" cellspacing="0" cellpadding="0" style="width:100%;">
								<tr>
									<td style="text-align:left">
										<!-- 내용 -->
										<table border="0" cellspacing="0" cellpadding="0" style="width:100%">
											<!-- 메인 타이틀 -->
											<tr>
												<td style="width:100%;font-size:27px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:43px;color:#171717;letter-spacing:-1px">` + title + `</td>
											</tr>
											<!-- 내용 텍스트 -->
											<tr>
												<td style="padding:35px 0 1px 0;width:100%;font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:25px;color:#171717;">` + convMsg + `</td>
											</tr>
										</table>
										<!-- //내용 -->
									</td>
								</tr>
							</table>
							<br><br>
							<table border="1" style="border:1px solid darkgray;width:100%;border-collapse:collapse;">
								<p style="font-size:25px;float:left">채용정보</p>
								<tr>
									<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>기업명</b></td>
									<td style="padding:5px;font-size:16px;letter-spacing:0px">` + entpKoNm + `</td>
								</tr>
								<tr>
									<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>채용공고</b></td>
									<td style="padding:5px;font-size:16px;letter-spacing:0px">` + recruitTitle + `</td>
								</tr>
								<tr>
									<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>지원방법</b></td>
									<td style="padding:5px;font-size:16px;letter-spacing:0.3px"><p>1. 아래 ‘바로 지원하기’ 버튼을 클릭하시거나 구글플레이/앱스토어/원스토어를 통해 ‘직감’어플리케이션을 설치합니다.</p>
										<p>2. 앱을 실행 후 회원가입을 진행합니다.</p>
										<p>3. 메일로 돌아오셔서 아래 ‘직감 바로가기’ 버튼을 다시 클릭하시거나 채용공고 제목을 검색합니다.</p>
										<p>4. 채용공고 내용을 확인 후 하단 [지원하기]를 선택해주세요.</p>
									</td>
								</tr>
								<tr>
									<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>유의사항</b></td>
									<td style="padding:5px;font-size:16px;letter-spacing:0.3px"><p>• 직감의 영상 지원은 모바일 어플리케이션을 통해서만 진행하실 수 있습니다.</p>
										<p>• 영상 촬영은 세로 화면으로만 진행됩니다. 모바일을 가로로 돌리지 않도록 유의해주세요.</p>
										<p>• 채용공고 마감일 이후에는 지원이 불가능하니 유의해주세요.</p>
										<p>* 수신자 정보는 30일 후 자동 삭제됩니다. 즉시 삭제를 원하실 경우 [바로가기] 링크를 통해 삭제 요청을 하실 수 있습니다.</p>
										<p>* 영상 인터뷰 관련 문의 : <a href="mailto:` + mailTo + `">` + supportMail + `"</a></p>
									</td>
								</tr>
							</table>
							<br>
							<div style="display: flex;align-items: center;justify-content: center;">
								<a href="` + entpVdUrl + `"
								style="height:11px;border: 1px solid darkgray;width:180px;border-radius:3px;background-color: lightgray;color: black; padding: 20px 34px; text-align: center; text-decoration: none;display: inline-block;font-size: 16px; margin: 4px 2px; cursor: pointer;">
									직감 영상 인터뷰 가이드
								</a>
								<a href=` + resultRecruitUrl + `
								style="height:11px;border: 1px solid darkgray;width:180px;border-radius:3px;background-color: darkgray; color: black; padding: 20px 34px; text-align: center; text-decoration: none;display: inline-block;font-size: 16px; margin: 4px 2px; cursor: pointer;">
									바로 지원하기
								</a>
							</div>
						</td>
					</tr>
					<!-- //contents -->
					<!-- footer -->
					<tr>
						<td style="padding:20px 0px 20px 0;background-color:#ced3d6;">
							<table border="0" cellspacing="0" cellpadding="0" style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center">
								<tr>
									<td style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center;color:#4d5256">
										본 메일은 발신 전용 메일입니다. 문의 사항은 <a href="mailto:` + mailTo + `" style="color:#4d5256;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;">support@ziggam.com</a>을 이용해주세요.
									</td>
								</tr>
								<tr>
									<td style="padding:10px 0 20px 0;width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:20px;text-align:center;color:#4d5256">
										©<b>Qrate</b> Corp. All rights reserved&nbsp;&nbsp;&nbsp;I&nbsp;&nbsp;&nbsp;경기도 성남시 분당구 황새울로 311번길 36, 2F<br>
										사업자등록번호 : 850-88-00704&nbsp;&nbsp;&nbsp;I&nbsp;&nbsp;&nbsp;직업정보제공사업신고 : J151602019001
									</td>
								</tr>
								<tr>
									<td style="text-align:center">
										<table border="0" cellspacing="0" cellpadding="0" style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center">
											<tr>
												<td style="width:266px">&nbsp;</td>
												<td><a href="https://www.ziggam.com" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-home.png" alt="HOME"></a></td>
												<td style="width:16px">&nbsp;</td>
												<td><a href="` + siteUrl + `" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-ems.png" alt="EMS"></a></td>
												<td style="width:16px">&nbsp;</td>
												<td><a href="` + googleStore + `" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-andapp.png" alt="Android"></a></td>
												<td style="width:16px">&nbsp;</td>
												<td><a href="` + appleStore + `" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-iosapp.png" alt="IOS"></a></td>
												<td style="width:266px">&nbsp;</td>
											</tr>
										</table>
									</td>
								</tr>
							</table>
						</td>
					</tr>
					<!-- //footer -->
				</table>
				</body>
				</html>`
		*/

		/*
			var htmlString = `
						<html>
						<head>
						<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
						<meta http-equiv="X-UA-Compatible" content="ie=edge">
						<title>직감 채용을 편하게, 면접을 영상으로</title>
						</head>
						<body style="margin:0;padding:0;background-color:#f5f6f9;">
						<table border="0" cellspacing="0" cellpadding="0" align="center" style="font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif">
							<!-- header -->
							<tr>
								<td style="padding:26px 0px 20px;border-bottom:4px solid #2ad0c7;">
									<table border="0" cellspacing="0" cellpadding="0" style="width:100%;">
										<tr>
											<td style="float:left">
												<a href="https://www.ziggam.com" target="_blank"><img style="display:block;border:0;" src="` + imgServer + `/mail/ic-main-logo.png" alt="직감 채용을 편하게,면접을 영상으로"></a>
											</td>
										</tr>
									</table>
								</td>
							</tr>
							<!-- //header -->
							<!-- contents -->
							<tr>
								<td style="padding:40px 70px;background-color:#ffffff;">
									<table border="0" cellspacing="0" cellpadding="0" style="width:100%;">
										<tr>
											<td style="text-align:left">
												<!-- 내용 -->
												<table border="0" cellspacing="0" cellpadding="0" style="width:100%">
													<!-- 메인 타이틀 -->
													<tr>
														<td style="width:100%;font-size:27px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:43px;color:#171717;letter-spacing:-1px">` + title + `</td>
													</tr>
													<!-- 내용 텍스트 -->
													<tr>
														<td style="padding:35px 0 1px 0;width:100%;font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:25px;color:#171717;">` + convMsg + `</td>
													</tr>
												</table>
												<!-- //내용 -->
											</td>
										</tr>
									</table>
									<br><br>
									<p style="font-size:25px;float:left">채용정보</p>
									<table border="1" style="border:1px solid darkgray;width:100%;border-collapse:collapse;">
										<tr>
											<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>기업명</b></td>
											<td style="padding:5px;font-size:16px;letter-spacing:0px">` + entpKoNm + `</td>
										</tr>
										<tr>
											<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>채용공고</b></td>
											<td style="padding:5px;font-size:16px;letter-spacing:0px">` + recruitTitle + `</td>
										</tr>
										<tr>
											<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>직무</b></td>
											<td style="padding:5px;font-size:16px;letter-spacing:0px">` + upJobGrp + `>` + jobGrp + `</td>
										</tr>
										<tr>
											<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>지원방법</b></td>
											<td style="padding:5px;font-size:16px;letter-spacing:0.3px"><p>1. 아래 ‘바로 지원하기’ 버튼을 클릭하시거나 구글플레이/앱스토어/원스토어를 통해 ‘직감’어플리케이션을 설치합니다.</p>
												<p>2. 앱을 실행 후 회원가입을 진행합니다.</p>
												<p>3. 메일로 돌아오셔서 아래 ‘직감 바로가기’ 버튼을 다시 클릭하시거나 채용공고 제목을 검색합니다.</p>
												<p>4. 채용공고 내용을 확인 후 하단 [지원하기]를 선택해주세요.</p>
											</td>
										</tr>
										<tr>
											<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>유의사항</b></td>
											<td style="padding:5px;font-size:16px;letter-spacing:0.3px"><p>• 직감의 영상 지원은 모바일 어플리케이션을 통해서만 진행하실 수 있습니다.</p>
												<p>• 영상 촬영은 세로 화면으로만 진행됩니다. 모바일을 가로로 돌리지 않도록 유의해주세요.</p>
												<p>• 채용공고 마감일 이후에는 지원이 불가능하니 유의해주세요.</p>
												<p>* 수신자 정보는 30일 후 자동 삭제됩니다. 즉시 삭제를 원하실 경우 [바로가기] 링크를 통해 삭제 요청을 하실 수 있습니다.</p>
												<p>* 영상 인터뷰 관련 문의 : <a href="mailto:` + mailTo + `">` + supportMail + `</a></p>
											</td>
										</tr>
									</table>
									<br>
									<div style="display: flex;text-align:center;justify-content: center;">
										<div></div>
										<a href="` + entpVdUrl + `"
										style="height:20px;border: 1px solid darkgray;width:180px;border-radius:3px;background-color: lightgray;color: black; padding: 10px 10px; text-align: center; text-decoration: none;font-size: 16px; margin: 4px 2px; cursor: pointer;">
											직감 영상 인터뷰 가이드
										</a>
										<a href=` + resultRecruitUrl + `
										style="height:20px;border: 1px solid darkgray;width:180px;border-radius:3px;background-color: darkgray; color: black; padding: 10px 10px; text-align: center; text-decoration: none;font-size: 16px; margin: 4px 2px; cursor: pointer;">
											바로 지원하기
										</a>
										<div></div>
									</div>
								</td>
							</tr>
							<!-- //contents -->
							<!-- footer -->
							<tr>
								<td style="padding:20px 0px 20px 0;background-color:#ced3d6;">
									<table border="0" cellspacing="0" cellpadding="0" style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center">
										<tr>
											<td style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center;color:#4d5256">
												본 메일은 발신 전용 메일입니다. 문의 사항은 <a href="mailto:` + mailTo + `" style="color:#4d5256;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;">support@ziggam.com</a>을 이용해주세요.
											</td>
										</tr>
										<tr>
											<td style="padding:10px 0 20px 0;width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:20px;text-align:center;color:#4d5256">
												©<b>Qrate</b> Corp. All rights reserved&nbsp;&nbsp;&nbsp;I&nbsp;&nbsp;&nbsp;경기도 성남시 분당구 황새울로 311번길 36, 2F<br>
												사업자등록번호 : 850-88-00704&nbsp;&nbsp;&nbsp;I&nbsp;&nbsp;&nbsp;직업정보제공사업신고 : J151602019001
											</td>
										</tr>
										<tr>
											<td style="text-align:center">
												<table border="0" cellspacing="0" cellpadding="0" style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center">
													<tr>
														<td style="width:266px">&nbsp;</td>
														<td><a href="https://www.ziggam.com" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-home.png" alt="HOME"></a></td>
														<td style="width:16px">&nbsp;</td>
														<td><a href="` + siteUrl + `" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-ems.png" alt="EMS"></a></td>
														<td style="width:16px">&nbsp;</td>
														<td><a href="` + googleStore + `" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-andapp.png" alt="Android"></a></td>
														<td style="width:16px">&nbsp;</td>
														<td><a href="` + appleStore + `" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-iosapp.png" alt="IOS"></a></td>
														<td style="width:266px">&nbsp;</td>
													</tr>
												</table>
											</td>
										</tr>
									</table>
								</td>
							</tr>
							<!-- //footer -->
						</table>
						</body>
						</html>`
		*/

		/*
			var htmlString = `
			<html>
			<head>
			<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
			<meta http-equiv="X-UA-Compatible" content="ie=edge">
			<title>직감 채용을 편하게, 면접을 영상으로</title>
			</head>
			<body style="margin:0;padding:0;background-color:#f5f6f9;">
			<table border="0" cellspacing="0" cellpadding="0" align="center" style="font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif">
				<!-- header -->
				<tr>
					<td style="padding:26px 0px 20px;border-bottom:4px solid #2ad0c7;">
						<table border="0" cellspacing="0" cellpadding="0" style="width:100%;">
							<tr>
								<td style="float:left">
									<a href="https://www.ziggam.com" target="_blank"><img style="display:block;border:0;" src="` + imgServer + `/mail/ic-main-logo.png" alt="직감 채용을 편하게,면접을 영상으로"></a>
								</td>
							</tr>
						</table>
					</td>
				</tr>
				<!-- //header -->
				<!-- contents -->
				<tr>
					<td style="padding:40px 70px;background-color:#ffffff;">
						<table border="0" cellspacing="0" cellpadding="0" style="width:100%;">
							<tr>
								<td style="text-align:left">
									<!-- 내용 -->
									<table border="0" cellspacing="0" cellpadding="0" style="width:100%">
										<!-- 메인 타이틀 -->
										<tr>
											<td style="width:100%;font-size:27px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:43px;color:#171717;letter-spacing:-1px">` + entpKoNm + `에서 영상 지원 요청이 도착했습니다.</td>
										</tr>
										<!-- 내용 텍스트 -->
										<tr>
											<td style="padding:35px 0 1px 0;width:100%;font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:25px;color:#171717;">` + val.Name + `님, 안녕하세요.<br>
											` + entpKoNm + `에서 직감을 통한 영상 지원 요청이 도착했습니다..<br>
											아래 상세 내용을 확인해주세요.</td>
										</tr>
									</table>
									<!-- //내용 -->
								</td>
							</tr>
						</table>
						<br><br>
						<table border="1" style="border:1px solid darkgray;width:100%;border-collapse:collapse;">
							<p style="font-size:25px;float:left">채용정보</p>
							<tr>
								<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>기업명</b></td>
								<td style="padding:5px;font-size:16px;letter-spacing:0px">` + entpKoNm + `</td>
							</tr>
							<tr>
								<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>채용공고</b></td>
								<td style="padding:5px;font-size:16px;letter-spacing:0px">` + recruitTitle + `</td>
							</tr>
							<tr>
								<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>지원방법</b></td>
								<td style="padding:5px;font-size:16px;letter-spacing:0.3px"><p>1. 아래 ‘바로 지원하기’ 버튼을 클릭하시거나 구글플레이/앱스토어/원스토어를 통해 ‘직감’어플리케이션을 설치합니다.</p>
									<p>2. 앱을 실행 후 회원가입을 진행합니다.</p>
									<p>3. 메일로 돌아오셔서 아래 ‘직감 바로가기’ 버튼을 다시 클릭하시거나 채용공고 제목을 검색합니다.</p>
									<p>4. 채용공고 내용을 확인 후 하단 [지원하기]를 선택해주세요.</p>
								</td>
							</tr>
							<tr>
								<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>유의사항</b></td>
								<td style="padding:5px;font-size:16px;letter-spacing:0.3px"><p>• 직감의 영상 지원은 모바일 어플리케이션을 통해서만 진행하실 수 있습니다.</p>
									<p>• 영상 촬영은 세로 화면으로만 진행됩니다. 모바일을 가로로 돌리지 않도록 유의해주세요.</p>
									<p>• 채용공고 마감일 이후에는 지원이 불가능하니 유의해주세요.</p>
									<p>* 영상 인터뷰 관련 문의 : <a href="mailto:` + mailTo + `">support@ziggam.com</a></p>
									<p>* 수신자 정보는 30일 후 자동 삭제됩니다. 즉시 삭제를 원하실 경우 <a href="mailto:` + mailTo + `">support@ziggam.com</a>을 통해 문의해주세요.</p>
								</td>
							</tr>
						</table>
						<br>
						<div style="display: flex;align-items: center;justify-content: center;">
							<a href="` + entpVdUrl + `"
							style="height:11px;border: 1px solid darkgray;width:180px;border-radius:3px;background-color: lightgray;color: black; padding: 20px 34px; text-align: center; text-decoration: none;display: inline-block;font-size: 16px; margin: 4px 2px; cursor: pointer;">
								직감 영상 인터뷰 가이드
							</a>
							<a href=` + resultRecruitUrl + `
							style="height:11px;border: 1px solid darkgray;width:180px;border-radius:3px;background-color: darkgray; color: black; padding: 20px 34px; text-align: center; text-decoration: none;display: inline-block;font-size: 16px; margin: 4px 2px; cursor: pointer;">
								바로 지원하기
							</a>
						</div>
					</td>
				</tr>
				<!-- //contents -->
				<!-- footer -->
				<tr>
					<td style="padding:20px 0px 20px 0;background-color:#ced3d6;">
						<table border="0" cellspacing="0" cellpadding="0" style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center">
							<tr>
								<td style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center;color:#4d5256">
									본 메일은 발신 전용 메일입니다. 문의 사항은 <a href="mailto:` + mailTo + `" style="color:#4d5256;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;">support@ziggam.com</a>을 이용해주세요.
								</td>
							</tr>
							<tr>
								<td style="padding:10px 0 20px 0;width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:20px;text-align:center;color:#4d5256">
									©<b>Qrate</b> Corp. All rights reserved&nbsp;&nbsp;&nbsp;I&nbsp;&nbsp;&nbsp;경기도 성남시 분당구 황새울로 311번길 36, 2F<br>
									사업자등록번호 : 850-88-00704&nbsp;&nbsp;&nbsp;I&nbsp;&nbsp;&nbsp;직업정보제공사업신고 : J151602019001
								</td>
							</tr>
							<tr>
								<td style="text-align:center">
									<table border="0" cellspacing="0" cellpadding="0" style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center">
										<tr>
											<td style="width:266px">&nbsp;</td>
											<td><a href="https://www.ziggam.com" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-home.png" alt="HOME"></a></td>
											<td style="width:16px">&nbsp;</td>
											<td><a href="` + siteUrl + `" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-ems.png" alt="EMS"></a></td>
											<td style="width:16px">&nbsp;</td>
											<td><a href="` + googleStore + `" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-andapp.png" alt="Android"></a></td>
											<td style="width:16px">&nbsp;</td>
											<td><a href="` + appleStore + `" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-iosapp.png" alt="IOS"></a></td>
											<td style="width:266px">&nbsp;</td>
										</tr>
									</table>
								</td>
							</tr>
						</table>
					</td>
				</tr>
				<!-- //footer -->
			</table>
			</body>
			</html>`
		*/

		var htmlString = `
					<html>
					<head>
					<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
					<meta name="viewport" content="width=device-width, initial-scale=1.0" />
					<title>직감 채용을 편하게, 면접을 영상으로</title>
					</head>
					<style type="text/css">
						/* CLIENT-SPECIFIC RESETS */
						/* Outlook.com(Hotmail)의 전체 너비 및 적절한 줄 높이를 허용 */
						.ReadMsgBody{ width: 100%; }
						.ExternalClass{ width: 100%; }
						.ExternalClass, .ExternalClass p, .ExternalClass span, .ExternalClass font, .ExternalClass td, .ExternalClass div { line-height: 100%; }
						/* Outlook 2007 이상에서 Outlook이 추가하는 테이블 주위의 간격을 제거 */
						table, td { mso-table-lspace: 0pt; mso-table-rspace: 0pt; }
						/* Internet Explorer에서 크기가 조정된 이미지를 렌더링하는 방식을 수정 */
						img { -ms-interpolation-mode: bicubic; }
						/* Webkit 및 Windows 기반 클라이언트가 텍스트 크기를 자동으로 조정하지 않도록 수정 */
						body, table, td, p, a, li, blockquote { -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%;}
					</style>
					<body>
					<!-- OUTERMOST CONTAINER TABLE -->
					<table border="0" cellpadding="0" cellspacing="0" width="100%" id="bodyTable" style="font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif">
						<!-- header -->
						<tr>
							<td style="padding:26px 0px 20px;border-bottom:4px solid #2ad0c7;">
								<table border="0" cellspacing="0" cellpadding="0" style="width:100%;">
									<tr>
										<td style="float:left">
											<a href="https://www.ziggam.com" target="_blank"><img style="display:block;border:0;" src="` + imgServer + `/mail/ic-main-logo.png" alt="직감 채용을 편하게,면접을 영상으로"></a>
										</td>
									</tr>
								</table>
							</td>
						</tr>
						<!-- //header -->
						<tr>
						<td>
							<!-- 600px - 800px CONTENTS CONTAINER TABLE -->
							<table border="0" cellpadding="0" cellspacing="0" width="600">
							<tr>
								<td style="text-align:left">
									<!-- 내용 -->
									<table border="0" cellspacing="0" cellpadding="0" style="width:100%">
										<!-- 메인 타이틀 -->											
										<!-- 내용 텍스트 -->										
										<tr>
											<td style="padding-top:20px;width:100%;font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:25px;color:#171717;">` + convMsg + `</td>
										</tr>
										<tr>
											<td style="padding-top:20px;padding-bottom:10px;width:100%;font-size:27px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:43px;color:#171717;letter-spacing:-1px">채용 정보</td>
										</tr>
										<table border="1" style="border:1px solid darkgray;width:100%;border-collapse:collapse;">
											<tr>
												<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>기업명</b></td>
												<td style="padding:5px;font-size:16px;letter-spacing:0px">` + entpKoNm + `</td>
											</tr>
											<tr>
												<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>채용공고</b></td>
												<td style="padding:5px;font-size:16px;letter-spacing:0px">` + recruitTitle + `</td>
											</tr>
											<tr>
												<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>직무</b></td>
												<td style="padding:5px;font-size:16px;letter-spacing:0px">` + upJobGrp + `>` + jobGrp + `</td>
											</tr>
											<tr>
												<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>지원방법</b></td>
												<td style="padding:5px;font-size:16px;letter-spacing:0.3px"><p>1. 아래 ‘바로 지원하기’ 버튼을 클릭하시거나 구글플레이/앱스토어/원스토어를 통해 ‘직감’어플리케이션을 설치합니다.</p>
													<p>2. 앱을 실행 후 회원가입을 진행합니다.</p>
													<p>3. 메일로 돌아오셔서 아래 ‘직감 바로가기’ 버튼을 다시 클릭하시거나 채용공고 제목을 검색합니다.</p>
													<p>4. 채용공고 내용을 확인 후 하단 [지원하기]를 선택해주세요.</p>
												</td>
											</tr>
											<tr>
												<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>유의사항</b></td>
												<td style="padding:5px;font-size:16px;letter-spacing:0.3px"><p>• 직감의 영상 지원은 모바일 어플리케이션을 통해서만 진행하실 수 있습니다.</p>
													<p>• 영상 촬영은 세로 화면으로만 진행됩니다. 모바일을 가로로 돌리지 않도록 유의해주세요.</p>
													<p>• 채용공고 마감일 이후에는 지원이 불가능하니 유의해주세요.</p>
													<p>* 수신자 정보는 30일 후 자동 삭제됩니다. 즉시 삭제를 원하실 경우 [바로가기] 링크를 통해 삭제 요청을 하실 수 있습니다.</p>
													<p>* 영상 인터뷰 관련 문의 : <a href="mailto:` + mailTo + `">` + supportMail + `</a></p>
												</td>
											</tr>
										</table>
										<table border="0" style="border:0px solid darkgray;width:100%;border-collapse:collapse;text-align:center;height:100px">
										<tr>
											<td> 
											<a href="` + entpVdUrl + `"
												style="height:20px;border: 1px solid darkgray;width:180px;border-radius:3px;background-color: lightgray;color: black; padding: 10px 10px; text-align: center; text-decoration: none;font-size: 16px; margin: 4px 2px; cursor: pointer;">
												직감 영상 인터뷰 가이드
											</a>
											<a href=` + resultRecruitUrl + `
												style="height:20px;border: 1px solid darkgray;width:180px;border-radius:3px;background-color: darkgray; color: black; padding: 10px 10px; text-align: center; text-decoration: none;font-size: 16px; margin: 4px 2px; cursor: pointer;">
												바로 지원하기
											</a>
											</td> 
										<tr>
										</table>
									</table>
									<!-- //내용 -->
								</td>
							</tr>
							</table>
						</td>
						</tr>						
					</table>
					<table border="0" cellpadding="0" cellspacing="0" width="600">
					<!-- footer -->
						<tr>
							<td style="padding:20px 0px 20px 0;background-color:#ced3d6;">
								<table border="0" cellspacing="0" cellpadding="0" style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center">
									<tr>
										<td style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center;color:#4d5256">
											본 메일은 발신 전용 메일입니다. 문의 사항은 <a href="mailto:` + mailTo + `" style="color:#4d5256;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;">support@ziggam.com</a>을 이용해주세요.
										</td>
									</tr>
									<tr>
										<td style="padding:10px 0 20px 0;width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:20px;text-align:center;color:#4d5256">
											©<b>Qrate</b> Corp. All rights reserved&nbsp;&nbsp;&nbsp;I&nbsp;&nbsp;&nbsp;경기도 성남시 분당구 황새울로 311번길 36, 2F<br>
											사업자등록번호 : 850-88-00704&nbsp;&nbsp;&nbsp;I&nbsp;&nbsp;&nbsp;직업정보제공사업신고 : J151602019001
										</td>
									</tr>
									<tr>
										<td style="text-align:center">
											<table border="0" cellspacing="0" cellpadding="0" style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center">
												<tr>
													<td style="width:266px">&nbsp;</td>
													<td><a href="https://www.ziggam.com" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-home.png" alt="HOME"></a></td>
													<td style="width:16px">&nbsp;</td>
													<td><a href="` + siteUrl + `" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-ems.png" alt="EMS"></a></td>
													<td style="width:16px">&nbsp;</td>
													<td><a href="` + googleStore + `" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-andapp.png" alt="Android"></a></td>
													<td style="width:16px">&nbsp;</td>
													<td><a href="` + appleStore + `" target="_blank" style="border:0"><img src="` + imgServer + `/mail/btn-zinggam-iosapp.png" alt="IOS"></a></td>
													<td style="width:266px">&nbsp;</td>
												</tr>
											</table>
										</td>
									</tr>
								</table>
							</td>
						</tr>
						<!-- //footer -->
						</table>
					</body>
					</html>`

		isEmailSend := utils.SendMail(val.Email, val.Name, "no-reply@ziggam.com", fmt.Sprintf("%v by 직감", entpKoNm), title, htmlString)

		fmt.Printf("%v", isEmailSend)
	}

}

func TestSendEmail() {

	t := time.Now()

	fmt.Printf(t.String())

	//utils.SendMailPoolEx_Start()

	//utils.SendMailPoolPush(fmt.Sprintf("dongkale@naver.com"), "이동관", "no-reply@ziggam.com", "큐레잇", fmt.Sprintf("subject:%v", t.Format("2006-01-02 15:04:05")), "htmlContents")
	//utils.SendMailPoolEx_Push(fmt.Sprintf("dongkale@naver.com"), "이동관", "no-reply@ziggam.com", "큐레잇", fmt.Sprintf("subject:%v", t.Format("2006-01-02 15:04:05")), "htmlContents")

	// SendMailPool -->
	// var arrayData1 []utils.SendMailPoolPushData

	// for i := 0; i < 50; i++ {

	// 	to := fmt.Sprintf("dongkale@naver.com")
	// 	toName := fmt.Sprintf("이동관")
	// 	from := "no-reply@ziggam.com"
	// 	fromName := "큐레잇"
	// 	subject := fmt.Sprintf("SendMailPool %v:%v", t.Format("2006-01-02 15:04:05"), i)
	// 	htmlContents := fmt.Sprintf("<html><body>Hello....<p>%v</p></body></html>", i)

	// 	// to2 := fmt.Sprintf("dongkale@qrate.co.kr")
	// 	// toName2 := fmt.Sprintf("이동관")

	// 	arrayData1 = append(arrayData1, utils.SendMailPoolPushData{
	// 		ToSend: fmt.Sprintf("%s=>%s[%d]:%s", "SendMailPool", toName, i, to),
	// 		ToSendData: utils.SendMailPoolData{
	// 			//To: []string{utils.SendMailPoolDestFmt(toName, to), utils.SendMailPoolDestFmt(toName2, to2)},
	// 			To:           []string{utils.SendMailPoolDestFmt(toName, to)},
	// 			From:         utils.SendMailPoolDestFmt(fromName, from),
	// 			Subject:      subject,
	// 			HtmlContents: htmlContents,
	// 		},
	// 		Cb:     JOB_FAIR,         //TestPool_Callback,
	// 		CbData: "E2018102500001", //fmt.Sprintf("CallbackData%v", i),
	// 	})
	// }

	// //utils.SendMailPoolMng.PushArray(arrayData1, 3)
	// err := utils.SendMailPoolMng.PushArrayEx(arrayData1, 3, CallbackComplete, "=========================================================")
	// fmt.Printf(err)

	// to := fmt.Sprintf("dongkale@naver.com")
	// toName := fmt.Sprintf("이동관")
	// from := "no-reply@ziggam.com"
	// fromName := "큐레잇"
	// subject := fmt.Sprintf("subject:%v:%v", t.Format("2006-01-02 15:04:05"))
	// htmlContents := fmt.Sprintf("<html><body>Hello....<p>%v</p><body></html>", 1)

	// utils.SendMailPoolMng.Push(fmt.Sprintf("%s:%s", toName, to),
	// 	to,
	// 	toName,
	// 	from,
	// 	fromName,
	// 	subject,
	// 	htmlContents,
	// 	nil,
	// 	nil,
	// 	TestPool_Callback,
	// 	fmt.Sprintf("CallbackData"))

	// SendMailPool <--

	// SendMailDaemom --->

	// utils.SendMailDaemonMng.Push(fmt.Sprintf("%s:%s", "dongkale@naver.com", "이동관1"), fmt.Sprintf("dongkale@naver.com"), "이동관1", "no-reply@ziggam.com", "큐레잇", fmt.Sprintf("subject:%v", t.Format("2006-01-02 15:04:05")), "<html><body>Hello....<p>World</p><body></htm>", TestCallback, "cbData1")
	// utils.SendMailDaemonMng.Push(fmt.Sprintf("%s:%s", "dongkale@naver.com", "이동관2"), fmt.Sprintf("dongkale@naver.com"), "이동관2", "no-reply@ziggam.com", "큐레잇", fmt.Sprintf("subject:%v", t.Format("2006-01-02 15:04:05")), "<html><body>Hello....<p>World</p><body></htm>", TestCallback1, "cbData2")
	// utils.SendMailDaemonMng.Push(fmt.Sprintf("%s:%s", "dongkale@naver.com", "이동관3"), fmt.Sprintf("dongkale@naver.com"), "이동관3", "no-reply@ziggam.com", "큐레잇", fmt.Sprintf("subject:%v", t.Format("2006-01-02 15:04:05")), "<html><body>Hello....<p>World</p><body></htm>", TestCallback2, "cbData3")

	// var arrayData2 []utils.SendMailDaemonPushData

	// for i := 0; i < 1; i++ {

	// 	to := fmt.Sprintf("dongkale@naver.com")
	// 	toName := fmt.Sprintf("이동관")
	// 	from := "no-reply@ziggam.com"
	// 	fromName := "큐레잇"
	// 	subject := fmt.Sprintf("SendMailDaemon %v:%v", t.Format("2006-01-02 15:04:05"), i)
	// 	htmlContents := fmt.Sprintf("<html><body>Hello....<p>%v</p></body></html>", i)

	// 	arrayData2 = append(arrayData2, utils.SendMailDaemonPushData{
	// 		ToSend: fmt.Sprintf("%s=>%s[%d]:%s", "SendMailDaemon", toName, i, to),
	// 		ToSendData: utils.SendMailDaemonData{
	// 			To:           to,
	// 			ToName:       toName,
	// 			From:         from,
	// 			FromName:     fromName,
	// 			Subject:      subject,
	// 			HtmlContents: htmlContents,
	// 		},
	// 		Cb:     JOB_FAIR2,
	// 		CbData: "E2018102500001",
	// 	})
	// }

	// utils.SendMailDaemonMng.PushArray(arrayData2, 3)

	// SendMailDaemom --->

	// mailDbPool -->
	utils.AddArrayTest(1)
	// mailDbPool -->

	// err2 := utils.SendMailEx("dongkale@naver.com", "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")
	// fmt.Printf(err2)

	// to2 := fmt.Sprintf("dongkale@naver.com")
	// toName2 := fmt.Sprintf("이동관")
	// from2 := "no-reply@ziggam.com"
	// fromName2 := "큐레잇"
	// subject2 := fmt.Sprintf("No SMTP %v:%v", t.Format("2006-01-02 15:04:05"), 0)
	// htmlContents2 := fmt.Sprintf("<html><body>Hello....<p>%v</p></body></html>", 0)

	// ret := utils.MailSend_Test(to2, toName2, from2, fromName2, subject2, htmlContents2)

	// fmt.Printf(ret)

	// ----------------------
	// ----------------------
	// ----------------------

	// err1 := utils.SendMailPoolEx(fmt.Sprintf("dongkale@naver.com"), "이동관", "no-reply@ziggam.com", "큐레잇", fmt.Sprintf("subject:%v", t.Format("2006-01-02 15:04:05")), "htmlContents")
	// beego.Trace(err1)

	// err2 := utils.SendMailDaemon(fmt.Sprintf("dongkale@naver.com"), "이동관", "no-reply@ziggam.com", "큐레잇", fmt.Sprintf("subject:%v", t.Format("2006-01-02 15:04:05")), "htmlContents")
	// beego.Trace(err2)

	//for i := 0; i < 100; i++ {
	//err1 := utils.SendMailEx(fmt.Sprintf("fghjklasd%d@naver.com", i), "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")
	//err1 := utils.SendMailDaemon(fmt.Sprintf("fghjklasd%d@naver.com", i), "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")
	//err1 := utils.SendMailDaemon(fmt.Sprintf("dongkale@naver.com"), "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")
	//beego.Trace(fmt.Sprintf("%v -> %v", fmt.Sprintf("fghjklasd%d@naver.com", i), err1))
	//beego.Trace(err1)

	// err2 := utils.SendMailDaemon(fmt.Sprintf("dongkale@naver.com"), "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")
	// beego.Trace(err2)

	// err3 := utils.SendMailDaemon(fmt.Sprintf("dongkale@naver.com"), "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")
	// beego.Trace(err3)

	// err4 := utils.SendMailDaemon(fmt.Sprintf("dongkale@naver.com"), "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")
	// beego.Trace(err4)

	// err5 := utils.SendMailDaemon(fmt.Sprintf("dongkale@naver.com"), "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")
	// beego.Trace(err5)

	// err6 := utils.SendMailDaemon(fmt.Sprintf("dongkale@naver.com"), "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")
	// beego.Trace(err6)

	// err7 := utils.SendMailDaemon(fmt.Sprintf("dongkale@naver.com"), "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")
	// beego.Trace(err7)

	// err8 := utils.SendMailDaemon(fmt.Sprintf("dongkale@naver.com"), "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")
	// beego.Trace(err8)

	// err9 := utils.SendMailDaemon(fmt.Sprintf("dongkale@naver.com"), "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")
	// beego.Trace(err9)

	// err10 := utils.SendMailDaemon(fmt.Sprintf("dongkale@naver.com"), "이동관", "no-reply@ziggam.com", "큐레잇", "subject", "htmlContents")
	// beego.Trace(err10)

	//time.Sleep(1 * time.Second)
	//}
}

func TestCallback(err error, cnt int, to string, cbData string) {

	fmt.Printf(fmt.Sprintf("Callback -> To:%v, CbData: %v, Error:%v", to, cbData, err))
}

func TestCallback1(err error, cnt int, to string, cbData string) {

	fmt.Printf(fmt.Sprintf("Callback1 -> To:%v, CbData: %v, Error:%v", to, cbData, err))
}

func TestCallback2(err error, cnt int, to string, cbData string) {

	fmt.Printf(fmt.Sprintf("Callback2 -> To:%v, CbData: %v, Error:%v", to, cbData, err))
}

func TestPool_Callback(err error, num int, cnt int, to string, cbData string) {

	fmt.Printf(fmt.Sprintf("Pool Callback[num:%v][Cnt:%v] -> To:%v, CbData: %v, Error:%v", num, cnt, to, cbData, err))
}

func JOB_FAIR(err error, num int, cnt int, toSend string, cbData string) {

	fmt.Printf(fmt.Sprintf("JOB_FAIR() [num:%v][cnt:%v] -> ToSend:%v, CbData: %v, Error:%v", num, cnt, toSend, cbData, err))

	if err != nil {
		return
	}

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := cbData

	// Start : Oracle DB Connection
	env, srv, ses, err := utils.GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Jobfair List
	fmt.Printf(fmt.Sprintf("CALL ZSP_ENTP_JOBFAIR_LIST('%v', '%v', :1)", pLang, pEntpMemNo))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_ENTP_JOBFAIR_LIST('%v', '%v', :1)",
		pLang, pEntpMemNo),
		ora.S, /* MNG_CD */
		ora.S, /* TITLE */
		ora.S, /* SDY */
		ora.S, /* EDY */
		ora.S, /* HOST_INSTITUTION */
		ora.S, /* MANAGE_AGENCY */
		ora.S, /* AGREEMENT_URL */
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

	jobFailrInfoList := make([]models.JobfairInfo, 0)

	var (
		jfMngCd           string
		jfTitle           string
		jfSdy             string
		jfEdy             string
		jfHostInstitution string
		jfManageAgency    string
	)

	if procRset.IsOpen() {
		for procRset.Next() {

			jfMngCd = procRset.Row[0].(string)
			jfTitle = procRset.Row[1].(string)
			jfSdy = procRset.Row[2].(string)
			jfEdy = procRset.Row[3].(string)
			jfHostInstitution = procRset.Row[4].(string)
			jfManageAgency = procRset.Row[5].(string)

			jobFailrInfoList = append(jobFailrInfoList, models.JobfairInfo{
				MngCd:           jfMngCd,
				Title:           jfTitle,
				Sdy:             jfSdy,
				Edy:             jfEdy,
				HostInstitution: jfHostInstitution,
				ManageAgency:    jfManageAgency,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}

	// End : Jobfair List

	fmt.Printf("%v", jobFailrInfoList)
}

func JOB_FAIR2(err error, cnt int, toSend string, cbData string) {

	fmt.Printf(fmt.Sprintf("JOB_FAIR2() [num:0][cnt:%v] -> ToSend:%v, CbData: %v, Error:%v", cnt, toSend, cbData, err))

	if err != nil {
		return
	}

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := cbData

	// Start : Oracle DB Connection
	env, srv, ses, err := utils.GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Jobfair List
	fmt.Printf(fmt.Sprintf("CALL ZSP_ENTP_JOBFAIR_LIST('%v', '%v', :1)", pLang, pEntpMemNo))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_ENTP_JOBFAIR_LIST('%v', '%v', :1)",
		pLang, pEntpMemNo),
		ora.S, /* MNG_CD */
		ora.S, /* TITLE */
		ora.S, /* SDY */
		ora.S, /* EDY */
		ora.S, /* HOST_INSTITUTION */
		ora.S, /* MANAGE_AGENCY */
		ora.S, /* AGREEMENT_URL */
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

	jobFailrInfoList := make([]models.JobfairInfo, 0)

	var (
		jfMngCd           string
		jfTitle           string
		jfSdy             string
		jfEdy             string
		jfHostInstitution string
		jfManageAgency    string
	)

	if procRset.IsOpen() {
		for procRset.Next() {

			jfMngCd = procRset.Row[0].(string)
			jfTitle = procRset.Row[1].(string)
			jfSdy = procRset.Row[2].(string)
			jfEdy = procRset.Row[3].(string)
			jfHostInstitution = procRset.Row[4].(string)
			jfManageAgency = procRset.Row[5].(string)

			jobFailrInfoList = append(jobFailrInfoList, models.JobfairInfo{
				MngCd:           jfMngCd,
				Title:           jfTitle,
				Sdy:             jfSdy,
				Edy:             jfEdy,
				HostInstitution: jfHostInstitution,
				ManageAgency:    jfManageAgency,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}

	// End : Jobfair List

	fmt.Printf("%v", jobFailrInfoList)

}

func CallbackComplete(index int, cbData string) {
	fmt.Printf(fmt.Sprintf("Complte:%v !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!", cbData))
}

// type HtmlContents struct {
// 	utils.MailDBPoolFunc
// }

// func (resp *utils.MailDBPoolFunc) HtmlContents(values int) string {

// 	return "HTML"
// }

/*
	htmlContents := `
					<html>
					<head>
					<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
					<meta http-equiv="X-UA-Compatible" content="ie=edge">
					<title>직감 채용을 편하게, 면접을 영상으로</title>
					</head>
					<body style="margin:0;padding:0;background-color:#f5f6f9;">
					<table border="0" cellspacing="0" cellpadding="0" align="center" style="font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif">
						<!-- header -->
						<tr>
							<td style="padding:26px 0px 20px;border-bottom:4px solid #2ad0c7;">
								<table border="0" cellspacing="0" cellpadding="0" style="width:100%;">
									<tr>
										<td style="float:left">
											<a href="https://www.ziggam.com" target="_blank"><img style="display:block;border:0;" src="` + `imgServer` + `/mail/ic-main-logo.png" alt="직감 채용을 편하게,면접을 영상으로"></a>
										</td>
									</tr>
								</table>
							</td>
						</tr>
						<!-- //header -->
						<!-- contents -->
						<tr>
							<td style="padding:40px 70px;background-color:#ffffff;">
								<table border="0" cellspacing="0" cellpadding="0" style="width:100%;">
									<tr>
										<td style="text-align:left">
											<!-- 내용 -->
											<table border="0" cellspacing="0" cellpadding="0" style="width:100%">
												<!-- 메인 타이틀 -->
												<tr>
													<td style="width:100%;font-size:27px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:43px;color:#171717;letter-spacing:-1px">` + `title` + `</td>
												</tr>
												<!-- 내용 텍스트 -->
												<tr>
													<td style="padding:35px 0 1px 0;width:100%;font-size:16px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:25px;color:#171717;">` + `convMsg` + `</td>
												</tr>
											</table>
											<!-- //내용 -->
										</td>
									</tr>
								</table>
								<br><br>
								<p style="font-size:25px;float:left">채용정보</p>
								<br>
								<table border="1" style="border:1px solid darkgray;width:100%;border-collapse:collapse;">
									<tr>
										<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>기업명</b></td>
										<td style="padding:5px;font-size:16px;letter-spacing:0px">` + `entpKoNm` + `</td>
									</tr>
									<tr>
										<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>채용공고</b></td>
										<td style="padding:5px;font-size:16px;letter-spacing:0px">` + `recruitTitle` + `</td>
									</tr>
									<tr>
										<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>직무</b></td>
										<td style="padding:5px;font-size:16px;letter-spacing:0px">` + `upJobGrp` + `>` + `jobGrp` + `</td>
									</tr>
									<tr>
										<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>지원방법</b></td>
										<td style="padding:5px;font-size:16px;letter-spacing:0.3px"><p>1. 아래 ‘바로 지원하기’ 버튼을 클릭하시거나 구글플레이/앱스토어/원스토어를 통해 ‘직감’어플리케이션을 설치합니다.</p>
											<p>2. 앱을 실행 후 회원가입을 진행합니다.</p>
											<p>3. 메일로 돌아오셔서 아래 ‘직감 바로가기’ 버튼을 다시 클릭하시거나 채용공고 제목을 검색합니다.</p>
											<p>4. 채용공고 내용을 확인 후 하단 [지원하기]를 선택해주세요.</p>
										</td>
									</tr>
									<tr>
										<td bgcolor="lightgray" style="padding:5px;width:100px;font-size:16px;color:black;text-align: center;"><b>유의사항</b></td>
										<td style="padding:5px;font-size:16px;letter-spacing:0.3px"><p>• 직감의 영상 지원은 모바일 어플리케이션을 통해서만 진행하실 수 있습니다.</p>
											<p>• 영상 촬영은 세로 화면으로만 진행됩니다. 모바일을 가로로 돌리지 않도록 유의해주세요.</p>
											<p>• 채용공고 마감일 이후에는 지원이 불가능하니 유의해주세요.</p>
											<p>* 수신자 정보는 30일 후 자동 삭제됩니다. 즉시 삭제를 원하실 경우 [바로가기] 링크를 통해 삭제 요청을 하실 수 있습니다.</p>
											<p>* 영상 인터뷰 관련 문의 : <a href="mailto:` + `mailTo` + `">` + `supportMail` + `</a></p>
										</td>
									</tr>
								</table>
								<br>
								<div style="display: flex;text-align:center;justify-content: center;">
									<div></div>
									<a href="` + `entpVdUrl` + `"
									style="height:20px;border: 1px solid darkgray;width:180px;border-radius:3px;background-color: lightgray;color: black; padding: 10px 10px; text-align: center; text-decoration: none;font-size: 16px; margin: 4px 2px; cursor: pointer">
										직감 영상 인터뷰 가이드
									</a>
									<a href=` + `resultRecruitUrl` + `
									style="height:20px;border: 1px solid darkgray;width:180px;border-radius:3px;background-color: darkgray; color: black; padding: 10px 10px; text-align: center; text-decoration: none;font-size: 16px; margin: 4px 2px; cursor: pointer">
										바로 지원하기
									</a>
									<div></div>
								</div>
							</td>
						</tr>
						<!-- //contents -->
						<!-- footer -->
						<tr>
							<td style="padding:20px 0px 20px 0;background-color:#ced3d6;">
								<table border="0" cellspacing="0" cellpadding="0" style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center">
									<tr>
										<td style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center;color:#4d5256">
											본 메일은 발신 전용 메일입니다. 문의 사항은 <a href="mailto:` + `mailTo` + `" style="color:#4d5256;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;">support@ziggam.com</a>을 이용해주세요.
										</td>
									</tr>
									<tr>
										<td style="padding:10px 0 20px 0;width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:20px;text-align:center;color:#4d5256">
											©<b>Qrate</b> Corp. All rights reserved&nbsp;&nbsp;&nbsp;I&nbsp;&nbsp;&nbsp;경기도 성남시 분당구 황새울로 311번길 36, 2F<br>
											사업자등록번호 : 850-88-00704&nbsp;&nbsp;&nbsp;I&nbsp;&nbsp;&nbsp;직업정보제공사업신고 : J151602019001
										</td>
									</tr>
									<tr>
										<td style="text-align:center">
											<table border="0" cellspacing="0" cellpadding="0" style="width:100%;font-size:12px;font-family:'맑은 고딕','Malgun Gothic','돋움',dotum,sans-serif;line-height:1.5;text-align:center">
												<tr>
													<td style="width:266px">&nbsp;</td>
													<td><a href="https://www.ziggam.com" target="_blank" style="border:0"><img src="` + `imgServer` + `/mail/btn-zinggam-home.png" alt="HOME"></a></td>
													<td style="width:16px">&nbsp;</td>
													<td><a href="` + `siteUrl` + `" target="_blank" style="border:0"><img src="` + `imgServer` + `/mail/btn-zinggam-ems.png" alt="EMS"></a></td>
													<td style="width:16px">&nbsp;</td>
													<td><a href="` + `googleStore` + `" target="_blank" style="border:0"><img src="` + `imgServer` + `/mail/btn-zinggam-andapp.png" alt="Android"></a></td>
													<td style="width:16px">&nbsp;</td>
													<td><a href="` + `appleStore` + `" target="_blank" style="border:0"><img src="` + `imgServer` + `/mail/btn-zinggam-iosapp.png" alt="IOS"></a></td>
													<td style="width:266px">&nbsp;</td>
												</tr>
											</table>
										</td>
									</tr>
								</table>
							</td>
						</tr>
						<!-- //footer -->
					</table>
					</body>
					</html>`
*/
