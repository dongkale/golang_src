package models

type FindId struct {
	ResultId string
	ResultDy string
}

type RtnFindId struct {
	RtnCd   int64
	RtnMsg  string
	RtnData FindId
}
