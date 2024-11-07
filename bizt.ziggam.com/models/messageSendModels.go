package models

type RtnMessageSend struct {
	RtnCd   int64
	RtnMsg  string
	RtnData string
}

type RtnMessageSend_v2 struct {
	RtnCd        int64
	RtnMsg       string
	RtnData      string
	RtnLiveSn    string
	RtnErrorCode string
}
