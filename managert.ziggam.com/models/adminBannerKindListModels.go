package models

type BannerKindList struct {
	EntpMemNo   string
	RecrutSn    string
	EntpKoNm    string
	RecrutTitle string
	SelectYn    string
}

type RtnBannerKindList struct {
	RtnBannerKindListData []BannerKindList
}
