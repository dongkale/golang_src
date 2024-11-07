package models

type AdminSimplePassLogin struct {
	MemNo  string
	MemId  string
	MemSn  string
	AuthCd string
}

type RtnAdminSimplePassLogin struct {
	RtnCd   int64
	RtnMsg  string
	RtnData AdminSimplePassLogin
}
