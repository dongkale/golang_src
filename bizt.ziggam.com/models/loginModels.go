package models

type Login struct {
	MemNo  string
	MemId  string
	MemSn  string
	AuthCd string
}

type RtnLogin struct {
	RtnCd   int64
	RtnMsg  string
	RtnData Login
}

type LoginKeepInfo struct {
	RtnCd     string
	RtnMsg    string
	RtnMemNo  string
	RtnMemId  string
	RtnMemSn  string
	RtnAuthCd string
}
