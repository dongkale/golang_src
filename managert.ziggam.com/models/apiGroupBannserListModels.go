package models

type GroupBannerData struct {
	RtnCd       int64
	RtnMsg      string
	BnrGrpTypCd string
	BnrGrpSn    string
	BnrGrpTitle string
	PtoPath     string
	PublSdy     string
	PublEdy     string
	PublStat    string
	UseYn       string
	RegDt       string
	RegId       string
	RolTm       int64
}

type GroupBannerDetailData struct {
	SwIdx            int64
	UseYn            string
	BnrSn            string
	BnrTitle         string
	PtoPath          string
	ThumbPtoPath     string
	PublSdy          string
	PublEdy          string
	PublStat         string
	RegDt            string
	BnrGrpSubSn      string
	PtoFullPath      string
	ThumbPtoFullPath string
}

//type GroupBannerDataList struct {
//	RtnGroupBannerData []GroupBannerData
//}
