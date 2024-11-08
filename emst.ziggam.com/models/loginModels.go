package models

type Login struct {
	MemNo string
	MemId string
}

type RtnLogin struct {
	RtnCd   int64
	RtnMsg  string
	RtnData Login
}
