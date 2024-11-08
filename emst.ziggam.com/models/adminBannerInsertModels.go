package models

type AdminBannerInsert struct {
	SetBnrSn string
}

type RtnAdminBannerInsert struct {
	RtnCd   int64
	RtnMsg  string
	RtnData AdminBannerInsert
}
