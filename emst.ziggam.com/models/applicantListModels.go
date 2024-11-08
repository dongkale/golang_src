package models

type ApplicantList struct {
	TotCnt         int64
	EntpMemNo      string
	RecrutSn       string
	PpMemNo        string
	RecrutTitle    string
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

	Pagination string
	MenuId     string
}

type RtnApplicantList struct {
	RtnApplicantListData []ApplicantList
}

type CmRecrutList struct {
	CmEntpMemNo   string
	CmRecrutSn    string
	CmRecrutTitle string
}

type ApplicantStatInfo struct {
	ApplyCnt int64
	IngCnt   int64
	PassCnt  int64
	FailCnt  int64
}
