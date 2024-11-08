package models

type AdminEventList struct {
	TotCnt      int64
	RegDt       string
	PpMemNo     string
	MemId       string
	Nm          string
	Sex         string
	Email       string
	BrthYmd     string
	Age         int64
	OsGbn       string
	VdCnt       int64
	VdCntString string
	Pagination  string
	MenuId      string
}

type RtnAdminEventList struct {
	RtnAdminEventListData []AdminEventList
}
