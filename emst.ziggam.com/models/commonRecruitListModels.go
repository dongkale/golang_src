package models

type CommonRecruitList struct {
	RecrutSn    string
	RecrutTitle string
}

type RtnCommonRecruitList struct {
	RtnCommonRecruitListData []CommonRecruitList
}
