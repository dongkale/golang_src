package models

type AdminMemberDetail struct {
	PpMemNo          string
	MemStatCd        string
	MemStatNm        string
	MemStatDt        string
	MemId            string
	Nm               string
	Sex              string
	Email            string
	MoNo             string
	BrthYmd          string
	Age              string
	EmailRecvYn      string
	SmsRecvYn        string
	PtoPath          string
	RegDt            string
	OsGbn            string
	OsVer            string
	MregPrgsStatCd   string
	MregPrgsStatNm   string
	JoinGbnNm        string
	SnsCd            string
	SnsCustNo        string
	LstEdu           string
	CarrGbn          string
	CarrDesc         string
	FrgnLangAbltDesc string
	AtchDataPath     string
	TechQlftKnd      string
	AtchFilePath     string
	TotApplyCnt      int64
	StnbyCnt         int64
	PassCnt          int64
	FailCnt          int64
	MachingRate      float64
}

type MemberVideoProfileList struct {
	VpVdsn        string
	VpThmKndCd    string
	VpThmNm       string
	VpQstCd       string
	VpQstDesc     string
	VpOpnSetCd    string
	VpOpnSetNm    string
	VpVdFilePath  string
	VpVdThumbPath string
	VpRegDt       string
}

type MemberApplyHistoryList struct {
	AhEntpMemNo     string
	AhRecrutSn      string
	AhPpMemNo       string
	AhEntpKoNm      string
	AhRecrutTitle   string
	AhEvlPrgsStatNm string
	AhTotCnt        int64
}
