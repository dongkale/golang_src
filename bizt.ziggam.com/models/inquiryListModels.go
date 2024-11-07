package models

type InquiryList struct {
	TotCnt     int64
	BrdNo      int64
	InqGbnNm   string
	RegDy      string
	InqSn      string
	InqTitle   string
	AnsYn      string
	Pagination string
}

type RtnInquiryList struct {
	RtnInquiryListData []InquiryList
}
