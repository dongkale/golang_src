package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// AligoSendSms ...
//ALARM("https://hooks.slack.com/services/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", "alarm", "ALARM", ":bot_alarm:"),
//TEST_ALARM("https://hooks.slack.com/services/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", "test_alarm", "[TEST]ALARM", ":wowwow:"),

//ERROR("https://hooks.slack.com/services/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", "live-alram", "ERROR", ":bot_error:"),
//TEST_ERROR("https://hooks.slack.com/services/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", "test_alarm", "[TEST]ERROR", ":zzzz:"),
//Live_DB_Error("https://hooks.slack.com/services/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", "live-db-alram", "alram", ":bot_alarm:");

// SlackMessageContents ...
type SlackMessageContents struct {
	Color     string `json:"color"`
	PreText   string `json:"pretext"`
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Text      string `json:"text"`
}

// SlackMessage ...
type SlackMessage struct {
	Text      string                 `json:"text"`
	Channel   string                 `json:"channel"`
	Username  string                 `json:"username"`
	IconEmoji string                 `json:"icon_emoji"`
	List      []SlackMessageContents `json:"attachments"`
}

// SlackSend ...
func SlackSend(runMode string, preText string, title string, message string) {

	//runMode := beego.AppConfig.String("runmode")

	var slackURL string
	var slackChannel string
	var slackUsername string
	var slackIcon string
	if runMode == "live" {
		// slackURL = "https://hooks.slack.com/services/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"
		// slackChannel = "live-alram"
		// slackUsername = "ERROR"
		// slackIcon = ":bot_error:"
		///////////////////////////////////////////////////////////////////////////////////////////////
		// slackURL = "https://hooks.slack.com/services/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"
		// slackChannel = "alarm_biz_live"
		// slackUsername = "[Live]Alarm"
		// slackIcon = ":bot_error:"
		///////////////////////////////////////////////////////////////////////////////////////////////
		slackURL = "https://hooks.slack.com/services/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"
		slackChannel = "alarm_live_svr"
		slackUsername = "[Live]Alarm"
		slackIcon = ":bot_error:"
	} else if runMode == "dev" {
		// slackURL = "https://hooks.slack.com/services/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"
		// slackChannel = "test_alarm"
		// slackUsername = "[TEST]ERROR"
		// slackIcon = ":zzzz:"
		///////////////////////////////////////////////////////////////////////////////////////////////
		//slackURL = "https://hooks.slack.com/services/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"
		//slackURL = "https://hooks.slack.com/services/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"
		// slackChannel = "alarm_biz_dev"
		// slackUsername = "[Dev]Alarm"
		// slackIcon = ":wowwow:"
		///////////////////////////////////////////////////////////////////////////////////////////////
		slackURL = "https://hooks.slack.com/services/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"
		slackChannel = "alarm_dev_svr"
		slackUsername = "[Dev]Alarm"
		slackIcon = ":bot_error:"
	} else {
		fmt.Printf("[SlackSend] Invalid RunMode: " + runMode)
		return
	}

	jsonData := SlackMessage{
		Text:      "[" + time.Now().Format("2006-01-02 15:04:05") + "]",
		Channel:   slackChannel,
		Username:  slackUsername,
		IconEmoji: slackIcon,
		List:      make([]SlackMessageContents, 0),
	}

	jsonData.List = append(jsonData.List, SlackMessageContents{
		Color:     "",
		PreText:   preText,
		Title:     title,
		TitleLink: "",
		Text:      message,
	})

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		panic(err)
	}

	payload := strings.NewReader(string(jsonBytes))

	req, _ := http.NewRequest("POST", slackURL, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Add("cache-control", "no-cache")

	resp, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Printf(fmt.Sprintf("[SlackSend][Result] SendData:%v, Result:%s", string(jsonBytes), string(body)))
}

func SlackSendDBError(runMode string, preText string, title string, message string) {

	//runMode := beego.AppConfig.String("runmode")
	//https://hooks.slack.com/services/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", "live-db-alram", "alram", ":bot_alarm:");

	var slackURL string
	var slackChannel string
	var slackUsername string
	var slackIcon string
	if runMode == "live" {
		slackURL = "https://hooks.slack.com/services/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"
		slackChannel = "live-db-alram"
		slackUsername = "alram"
		slackIcon = ":bot_alarm:"
	} else if runMode == "dev" {
		slackURL = "https://hooks.slack.com/services/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"
		slackChannel = "live-db-alram"
		slackUsername = "alram"
		slackIcon = ":bot_alarm:"
	} else {
		fmt.Printf("[SlackSend] Invalid RunMode: " + runMode)
		return
	}

	jsonData := SlackMessage{
		Text:      "[" + time.Now().Format("2006-01-02 15:04:05") + "]",
		Channel:   slackChannel,
		Username:  slackUsername,
		IconEmoji: slackIcon,
		List:      make([]SlackMessageContents, 0),
	}

	jsonData.List = append(jsonData.List, SlackMessageContents{
		Color:     "",
		PreText:   preText,
		Title:     title,
		TitleLink: "",
		Text:      message,
	})

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		panic(err)
	}

	payload := strings.NewReader(string(jsonBytes))

	req, _ := http.NewRequest("POST", slackURL, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Add("cache-control", "no-cache")

	resp, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Printf(fmt.Sprintf("[SlackSend(DB)][Result] SendData:%v, Result:%s", string(jsonBytes), string(body)))
}
