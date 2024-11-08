package models

type AdminIntroPopUpWrite struct {
	IntroTitle   string
	Sdy          string
	Edy          string
	LnkGbnCd     string
	LnkGbnVal    string
	BrdGbnCd     string
	Sn           int64
	EntpMemNo    string
	RecrutSn     string
	DelYn        string
	UseYn        string
	PtoPath      string
	ThumbPtoPath string
	OriImgFile   string
}

type CommonLnkGbnCd struct {
	LgCdId string
	LgCdNm string
}

type CommonLnkMenuCd struct {
	LvCdId string
	LvCdNm string
}
