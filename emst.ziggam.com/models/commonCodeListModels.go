package models

type CommonCodeList struct {
	CdId string
	CdNm string
}

type RtnCommonCodeList struct {
	RtnCommonCodeListData []CommonCodeList
}
