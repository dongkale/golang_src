package models

type AdminEntpList struct {
	TotCnt      int64
	EntpMemNo   string
	//JobFairCds    string
	JobFairCdsArr []string
	MemStat     string
	MemStatDt   string
	EntpMemId   string
	EntpKoNm    string
	RepNm       string
	RegDt       string
	MemStatCd   string
	TotAplyCnt  int64
	NewAplyCnt  int64
	BizRegNo    string
	Email       string
	PpChrgNm    string
	PpChrgTelNo string
	VdYn        string
	UseYn       string
	OsGbn       string
	LastLogin     string
	Pagination  string
	MenuId      string
}

type AdminEntpTopInfo struct {
	ETotCnt int64
	EComCnt int64
	EStbCnt int64
	EWtdCnt int64
}

type RtnAdminEntpList struct {
	RtnAdminEntpListData []AdminEntpList
}
