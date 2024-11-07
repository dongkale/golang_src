package models

type ApplicantListExcel struct {
	AleApplyDt          string
	AleFavrAplyPpYn     string
	AleNm               string
	AleSex              string
	AleBrthYmd          string
	AleAge              string
	AleEmail            string
	AleMoNo             string
	AleEvlPrgsStatNm    string
	AleEvlStatDt        string
	AlePrgsStatCd       string
	AleUpJobGrp         string
	AleJobGrp           string
	AleRecrutTitle      string
	AleRecrutDy         string
	AleLstEdu           string
	AleCarrGbn          string
	AleCarrDesc         string
	AleFrgnLangAbltDesc string
	AleAtchDataPath     string
	AleTechQlftDesc     string
	AleAtchFilePathYn   string
	AleVpYn             string
	AleScoreValue       int64
	DownloadPath        string
}

type RtnApplicantListExcel struct {
	RtnApplicantListExcelData []ApplicantListExcel
}
