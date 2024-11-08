package models

type AdminIntroPopUpList struct {
	TotCnt      int64
	IntroSn     string
	LnkGbnCd    string
	LnkGbnNm    string
	LnkGbnVal   string
	BrdGbnCd    string
	Sn          int64
	PtoPath     string
	RegDt       string
	IntroTitle  string
	UseYn       string
	LnkGbnValNm string
	Sdy         string
	Edy         string
	StbYn       string
	EndYn       string
	Pagination  string
	MenuId      string
}

type RtnAdminIntroPopUpList struct {
	RtnAdminIntroPopUpListData []AdminIntroPopUpList
}

type AdminIntroPopUpStat struct {
	IntroUseYn string
}
