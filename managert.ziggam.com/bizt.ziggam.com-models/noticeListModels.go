package models

type NoticeList struct {
	TotCnt     int64
	Sn         int64
	GbnNm      string
	Title      string
	RegDt      string
	NewYn      string
	Pagination string
}

type RtnNoticeList struct {
	RtnNoticeListData []NoticeList
}
