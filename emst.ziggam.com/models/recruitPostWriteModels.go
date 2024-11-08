package models

type JobGrpList struct {
	JobGrpCd   string
	UpJobGrpCd string
	JobGrpNm   string
}

type RtnJobGrpList struct {
	RtnJobGrpListData []JobGrpList
}
