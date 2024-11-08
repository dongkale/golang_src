package models

type RecruitStatList struct {
	TotCnt         int64
	EntpMemNo      string
	RecrutSn       string
	FavrAplyPpYn   string
	Nm             string
	Sex            string
	Age            string
	RegDt          string
	ApplyDt        string
	EvlStatDt      string
	EvlPrgsStatCd  string
	RcrtAplyStatCd string
	EntpCfrmYn     string
	LeftDy         string
	Tm             string
	VpYn           string
	PpMemNo        string

	Pagination string
	MenuId     string
}

type RtnRecruitStatList struct {
	RtnRecruitStatListData []RecruitStatList
}

type RecruitStatTopInfo struct {
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
