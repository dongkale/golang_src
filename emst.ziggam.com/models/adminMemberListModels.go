package models

type AdminMemberList struct {
	TotCnt         int64
	PpMemNo        string
	MemStatCd      string
	MemStatNm      string
	MemId          string
	Nm             string
	Sex            string
	Email          string
	BrthYmd        string
	Age            int64
	OsGbn          string
	MregPrgsStatCd string
	VpYn           string
	RegDt          string
	RegDy          string
	ArrVdPath      string
	MemJoinGbnCd   string
	MemJoinGbnNm   string
	MoNo		   string
	AhTotCnt	   int64
	JobFairCdsArr  []string
	LoginDt        string

	Pagination string
	MenuId     string
}

type RtnAdminMemberList struct {
	RtnAdminMemberListData []AdminMemberList
}

type CommonYYList struct {
	YYYY string
}

type CommonMMList struct {
	MM string
}

type CommonDDList struct {
	DD string
}

type MemberTopInfo struct {
	TotMemCnt int64
	RunMemCnt int64
	WtdMemCnt int64
}

const (
	DefaultSdy      string = "20180101"
	DefaultLoginSdy string = "20180101"
)
