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
}

type RtnRecruitPostList struct {
	RtnRecruitPostListData []RecruitPostList
}

type RecruitStatInfo struct {
	RecrutTotCnt  int64
	RecrutIngCnt  int64
	RecrutWaitCnt int64
	RecrutEndCnt  int64
}

type RecruitSubList struct {
	STotCnt      int64
	SEntpMemNo   string
	SRecrutSn    string
	SPrgsStat    string
	SRecrutTitle string
	SEmplTyp     string
	SUpJobGrp    string
	SJobGrp      string
	SRecrutDy    string
	SRecrutEdt   string
	SApplyCnt    int64
	SIngCnt      int64
	SPassCnt     int64
	SFailCnt     int64
	SRegDt       string
	SRegId       string
	SPpChrgBpNm  string
	SPpChrgNm    string
	Pagination   string
}

type RtnRecruitSubList struct {
	RtnRecruitSubListData []RecruitSubList
}

type RecruitMainJobGrpList struct {
	REmplTypCd string
	RJobGrpCd  string
	RJobGrpNm  string
}

/*



type RtnRecruitJobGrpList struct {
	RtnRecruitJobGrpListData []RecruitJobGrpList
}
*/
