package models

type AdminStatsMemberDetail struct {
	AnalDt  string
	AnalCnt int64
}

type AdminStatsMemberTopCnt struct {
	ComCnt int64
	WtdCnt int64
	UvfCnt int64
}

type RtnAdminStatsMemberDetail struct {
	RtnAdminStatsMemberDetailData []AdminStatsMemberDetail
}
