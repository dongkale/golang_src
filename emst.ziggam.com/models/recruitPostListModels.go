package models

type RecruitPostList struct {
	TotCnt      int64
	EntpMemNo   string
	RecrutSn    string
	PrgsStat    string
	RecrutTitle string
	EmplTyp     string
	UpJobGrp    string
	JobGrp      string
	RecrutDy    string
	RecrutEdt   string
	ApplyCnt    int64
	IngCnt      int64
	PassCnt     int64
	FailCnt     int64

	Pagination string
	MenuId     string
}

type RtnRecruitPostList struct {
	RtnRecruitPostListData []RecruitPostList
}

type RecruitStatInfo struct {
	RecrutTotCnt int64
	RecrutIngCnt int64
	RecrutEndCnt int64
}

type RecruitJobGrpList struct {
	REmplTypCd string
	RJobGrpCd  string
	RJobGrpNm  string
}

type RtnRecruitJobGrpList struct {
	RtnRecruitJobGrpListData []RecruitJobGrpList
}
