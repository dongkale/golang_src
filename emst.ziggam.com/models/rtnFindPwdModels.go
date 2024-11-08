package models

type RtnFindPwdStep1 struct {
	RtnCd  int64
	RtnMsg string
	CertNo string
}

type RtnFindPwdStep2 struct {
	RtnCd     int64
	RtnMsg    string
	TempMemId string
}

type RtnResetPwd struct {
	RtnCd  int64
	RtnMsg string
}
