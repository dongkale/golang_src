package models

type AdminMemberListExcel struct {
	ApplyDt      string
	MemJoinGbnNm string
	Nm           string
	VpYn         string
	MemId        string
	Sex          string
	BrthYmd      string
	Age          string
	Email        string
	MoNo         string
	AhTotCnt     string
	OsGbn        string
	MemStatNm    string
	JobFairCode  string
	LoginDt      string
	DownloadPath string
}

type RtnAdminMemberListExcel struct {
	RtnAdminMemberListExcelData []AdminMemberListExcel
}
