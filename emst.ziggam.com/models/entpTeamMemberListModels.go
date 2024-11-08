package models

type EntpTeamMemberList struct {
	EtTotCnt      int64
	EtPpChrgSn    string
	EtPpChrgGbnCd string
	EtPpChrgNm    string
	EtPpChrgBpNm  string
	EtEmail       string
	EtEntpMemId   string
	EtPpChrgTelNo string
	EtRowNo       int64
}

type RtnEntpTeamMemberList struct {
	RtnEntpTeamMemberListData []EntpTeamMemberList
}
