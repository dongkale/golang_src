package models

type RecruitApplyDetail struct {
	EntpMemNo        string
	RecrutSn         string
	PpMemNo          string
	PtoPath          string
	FavrAplyPpYn     string
	Nm               string
	Sex              string
	Age              string
	Email            string
	ApplyDt          string
	LeftDy           string
	ShootTm          string
	ShootCnt         int64
	VpYn             string
	LstEdu           string
	CarrGbn          string
	CarrDesc         string
	FrgnLangAbltDesc string
	AtchDataPath     string
	TechQlftKnd      string
	AtchFilePath     string
	MoNo             string
	EntpGroupCode	 string
}

type RecruitApplyTopInfo struct {
	PrgsStat    string
	RecrutTitle string
	EmplTyp     string
	UpJobGrp    string
	JobGrp      string
	RecrutDy    string
	RecrutEdt   string
	PrgsMsg     string
}

type RecruitApplyMemberAnswerList struct {
	AnsEntpMemNo  string
	AnsRecrutSn   string
	AnsQstSn      string
	AnsVdTitle    string
	AnsVdFilePath string
	AnsTotCnt     int64
}

type VideoProfileList struct {
	VpVdsn        string
	VpVdSec       int64
	VpVdThumbPath string
	VpVdFilePath  string
	VpThmKndCd    string
	VpThmNm       string
	VpQstCd       string
	VpQstDesc     string
	VpOpnSetCd    string
	VpRegDt       string
	VpTotCnt      int64
	VpSn          int64
}
