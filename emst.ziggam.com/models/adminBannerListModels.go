package models

type AdminBannerList struct {
	TotCnt      int64
	BnrSn       string
	BnrKndCd    string
	LnkGbnCd    string
	LnkGbnNm    string
	LnkGbnVal   string
	BrdGbnCd    string
	Sn          int64
	PtoPath     string
	RegDt       string
	BnrTitle    string
	UseYn       string
	LnkGbnValNm string
	Sdy         string
	Edy         string
	StbYn       string
	EndYn       string
	BnrKndNm    string
	Pagination  string
	MenuId      string
}

type RtnAdminBannerList struct {
	RtnAdminBannerListData []AdminBannerList
}

type AdminBannerStat struct {
	BnrUseYn string
	RolTime  int64
}
