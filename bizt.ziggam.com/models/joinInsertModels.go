package models

type JoinInsert struct {
	SetMemNo string
}

type RtnJoinInsert struct {
	RtnCd   int64
	RtnMsg  string
	RtnData JoinInsert
}
