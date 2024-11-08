package models

type AdminEntpListExcel struct {
	ApplyDt      string
	ApplyStat    string
	EntpMemNo    string
	JobFairCode  string
	EntpKoNm     string
	BizRegNo     string
	RepNm        string
	EntpMemId    string
	PpChrgNm     string
	PpChrgTelNo  string
	TotAplyCnt   int64
	NewAplyCnt   int64
	VpYn         string
	OsGbn        string
	VerifyStat   string
	LastLoginDt  string
	DownloadPath string
}

type RtnAdminEntpListExcel struct {
	RtnAdminEntpListExcelData []AdminEntpListExcel
}
