package models

type AdminIntroPopUpInsert struct {
	SetIntroSn string
}

type RtnAdminIntroPopUpInsert struct {
	RtnCd   int64
	RtnMsg  string
	RtnData AdminIntroPopUpInsert
}
