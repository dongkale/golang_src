package models

// JobfairInfo ...
type JobfairInfo struct {
	MngCd           string
	Title           string
	Sdy             string
	Edy             string
	HostInstitution string
	ManageAgency    string
}

// AdminStatsPeriod01 ...
type AdminStatsPeriod01 struct {
	EntpKoNm    string
	EntpMemNo   string
	JobFairCds  string
	RecrutTitle string
	RecrutSn    string
	RecrutJfMngCd string
	UpJobGrp    string
	JobGrp      string
	DcmntUseCd     string
	OnwyUseCd      string
	LiveIntrvUseCd string
	Sdy         string
	Edy         string
	RecrutEdt   string
	RegDt       string
}

// RtnAdminStatsPeriod01 ...
type RtnAdminStatsPeriod01 struct {
	Result []AdminStatsPeriod01
}

// AdminStatsPeriod02 ...
type AdminStatsPeriod02 struct {
	EntpKoNm         string
	EntpMemNo        string
	JobFairCds       string
	RecrutTitle      string
	RecrutSn         string
	RecrutJfMngCd    string
	UpJobGrp         string
	JobGrp           string
	RegDt            string
	Nm               string
	PpMemNo          string
	Sex              string
	Birth            string
	MoNo             string
	Email            string
	LstEdu           string
	LstEduDesc       string
	CarrGbn          string
	CarrGbnDesc      string
	FrgnLangAbltDesc string
	TechQlftKnd      string
}

// RtnAdminStatsPeriod02 ...
type RtnAdminStatsPeriod02 struct {
	Result []AdminStatsPeriod02
}

// AdminStatsPeriod03 ...
type AdminStatsPeriod03 struct {
	EntpKoNm         string
	EntpMemNo        string
	JobFairCds       string
	RecrutTitle      string
	RecrutSn         string
	RecrutJfMngCd    string
	UpJobGrp         string
	JobGrp           string
	RegDt            string
	OneWayDt         string
	OneWayCnt        string
	OneWayType       string
	Nm               string
	Sex              string
	Birth            string
	MoNo             string
	Email            string
	LstEdu           string
	LstEduDesc       string
	CarrGbn          string
	CarrGbnDesc      string
	FrgnLangAbltDesc string
	TechQlftKnd      string
}

// RtnAdminStatsPeriod03 ...
type RtnAdminStatsPeriod03 struct {
	Result []AdminStatsPeriod03
}

// AdminStatsPeriod04 ...
type AdminStatsPeriod04 struct {
	EntpKoNm         string
	EntpMemNo        string
	JobFairCds       string
	RecrutTitle      string
	RecrutSn         string
	RecrutJfMngCd    string
	UpJobGrp         string
	JobGrp           string
	RequestDt        string
	BeginDt          string
	EndDt            string
	LiveState        string
	Nm               string
	Sex              string
	Birth            string
	MoNo             string
	Email            string
	LstEdu           string
	LstEduDesc       string
	CarrGbn          string
	CarrGbnDesc      string
	FrgnLangAbltDesc string
	TechQlftKnd      string
}

// RtnAdminStatsPeriod04 ...
type RtnAdminStatsPeriod04 struct {
	Result []AdminStatsPeriod04
}

// AdminStatsPeriod05 ...
type AdminStatsPeriod05 struct {
	EntpKoNm         string
	EntpMemNo        string
	JobFairCds       string
	RecrutTitle      string
	RecrutSn         string
	RecrutJfMngCd    string
	UpJobGrp         string
	JobGrp           string
	ConfirmDt        string
	BeginDt          string
	EndDt            string
	Nm               string
	Sex              string
	Birth            string
	MoNo             string
	Email            string
	LstEdu           string
	LstEduDesc       string
	CarrGbn          string
	CarrGbnDesc      string
	FrgnLangAbltDesc string
	TechQlftKnd      string
}

// RtnAdminStatsPeriod05 ...
type RtnAdminStatsPeriod05 struct {
	Result []AdminStatsPeriod05
}

// AdminStatsPeriod06 ...
type AdminStatsPeriod06 struct {
	EntpKoNm      string
	EntpMemNo     string
	JobFairCds    string
	RecrutTitle   string
	RecrutSn      string
	RecrutJfMngCd string
	UpJobGrp      string
	JobGrp        string
	RequestDt     string
	BeginDt       string
	EndDt         string
	LiveState     string
	ApplyList     string
	MemList       string
}

// RtnAdminStatsPeriod06 ...
type RtnAdminStatsPeriod06 struct {
	Result []AdminStatsPeriod06
}
