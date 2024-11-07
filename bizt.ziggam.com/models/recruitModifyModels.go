package models

type RecruitModify struct {
	EntpMemNo      string
	RecrutSn       string
	RecrutTitle    string
	UpJobGrp       string
	JobGrp         string
	RecrutGbnCd    string
	RecrutCnt      int64
	Rol            string
	AplyQufct      string
	PerferTrtm     string
	Sdy            string
	Edy            string
	VdTitleUptYn   string
	DcmntEvlUseCd  string
	OnwyIntrvUseCd string
	LiveIntrvUseCd string
	PrgsStatCd     string
	RegDt          string
	PpChrgBpNm     string
	PpChrgNm       string
	SelUpJobGrpCd  string
	SelJobGrpCd    string
	RecrutProdCd   string
	RecrutEdt      string
}

type RecruitQuestionList struct {
	QstTotCnt int64
	QstSn     string
	VdTitle   string
}
