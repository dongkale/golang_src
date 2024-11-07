package models

type JobGrpList struct {
	JobGrpCd      string
	UpJobGrpCd    string
	JobGrpNm      string
	ChkJobGrpCd   string
	ChkUpJobGrpCd string
}

type RtnJobGrpList struct {
	EmplTypCd  string
	UpJobGrpCd string

	RtnJobGrpListData []JobGrpList
}

type RecruitJobGrpList struct {
	REmplTypCd string
	RJobGrpCd  string
	RJobGrpNm  string
}

type RtnRecruitJobGrpList struct {
	RtnRecruitJobGrpListData []RecruitJobGrpList
}
