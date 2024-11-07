package models


// Key - Value Model
type OptionList struct {
	OptionKey   string
	OptionValue string
	SelectYn	string
}

// 결과값 리턴
type RtnResult struct {
	RtnCd   int64
	RtnMsg  string
}


//type GroupBannerDataList struct {
//	RtnGroupBannerData []GroupBannerData
//}
