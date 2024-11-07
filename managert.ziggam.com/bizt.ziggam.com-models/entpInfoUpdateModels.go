package models

type EntpInfoUpdate struct {
	SetMemNo string
}

type RtnEntpInfoUpdate struct {
	RtnCd   int64
	RtnMsg  string
	RtnData EntpInfoUpdate
}
