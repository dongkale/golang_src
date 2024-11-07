package models

type ApplicantDetailPopup struct {
	EntpMemNo              string
	RecrutSn               string
	PpMemNo                string
	PtoPath                string
	Nm                     string
	Sex                    string
	Age                    int64
	Email                  string
	MoNo                   string
	PrgsStatCd             string
	UpJobGrp               string
	JobGrp                 string
	RecrutTitle            string
	Sdy                    string
	Edy                    string
	LstEdu                 string
	CarrGbn                string
	CarrDesc               string
	FrgnLangAbltDesc       string
	AtchDataPath           string
	TechQlftKnd            string
	AtchFilePath           string
	FavrAplyPpYn           string
	EvlPrgsStatCd          string
	EvlStatDt              string
	LiveReqStatCd          string
	MsgYn                  string
	MsgEndYn               string
	ApplyDt                string
	EntpGroupCode          string
	Dcmnt_evl_stat_cd      string
	Onwy_intrv_evl_stat_cd string
	Live_intrv_evl_stat_cd string
	Dcmnt_file_path        string
	Dcmnt_evl_stat_dt      string
	Onwy_intrv_evl_stat_dt string
	Live_intrv_evl_stat_dt string
	Dcmnt_file_name        string
	Recrut_proc_cd         string
	ShootCnt               int64
	CompTm                 string
	CompDT1                string
	CompDT2                string
	LiveSn                 string
	ReadEndDay             string
}

type RecruitApplyMemberAnswerList struct {
	AnsEntpMemNo  string
	AnsRecrutSn   string
	AnsQstSn      string
	AnsVdTitle    string
	AnsVdFilePath string
	AnsTotCnt     int64

	AnsVdFilePathImgSvr string
}

type RecruitApplyCommentList struct {
	CmtTotCnt      int64
	CmtEntpMemNo   string
	CmtRecrutSn    string
	CmtPpMemNo     string
	CmtPpChrgCmtSn string
	CmtPpChrgSn    string
	CmtPpChrgCmt   string
	CmtRegDt       string
	CmtRegId       string
	CmtPpChrgBpNm  string
	CmtPpChrgNm    string
	CmtNewYn       string
	CmtPpChrgGbnCd string
	CmtSMemId      interface{}
	CmtSAuthCd     interface{}
}

type RtnRecruitApplyCommentList struct {
	RtnRecruitApplyCommentListData []RecruitApplyCommentList
}

type MemberVideoProfileList struct {
	MvVdSn        string
	MvThmKndCd    string
	MvThmNm       string
	MvQstSCd      string
	MvQstDesc     string
	MvOpnSetCd    string
	MvOpenSetNm   string
	MvVdFilePath  string
	MvVdThumbPath string
	MvRegDt       string
}
