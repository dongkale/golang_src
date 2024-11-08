package models

type RecruitApplyMemberExcel struct {
	RecrutTitle      string
	EmplTyp          string
	JobGrpNm         string
	RecrutDy         string
	RecrutEdt        string
	EvlPrgsStat      string
	FavrAplyPpYn     string
	Nm               string
	Sex              string
	Age              string
	Email            string
	ApplyDt          string
	ShootTm          string
	ShootCnt         string
	LeftDy           string
	EvlStatDt        string
	VpYn             string
	LstEdu           string
	CarrDesc         string
	FrgnLangAbltDesc string
	AtchDataPath     string
	DownloadPath     string
}

type RtnRecruitApplyMemberExcel struct {
	RtnRecruitApplyMemberExcelData []RecruitApplyMemberExcel
}
