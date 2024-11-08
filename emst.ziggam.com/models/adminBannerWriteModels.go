package models

type AdminBannerWrite struct {
	BnrTitle     			string
	Sdy          			string
	Edy          			string
	LnkGbnCd     			string
	LnkGbnVal    			string
	BrdGbnCd     			string
	Sn           			int64
	EntpMemNo    			string
	RecrutSn     			string
	DelYn        			string
	UseYn        			string
	PtoPath      			string
	ThumbPtoPath 			string
	OriImgFile   			string
	BnrKndCd     			string
	ListBnrSn	 			string
	ListTitle	 			string
	ListPhotoPath			string
	ListThumbPhotoPath		string
	ListLinkUrl				string
}

type CommonBnrKndCd struct {
	BkCdId string
	BkCdNm string
}