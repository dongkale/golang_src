package models

type ContentClItem struct {
	BnrGrpSn    string
	UseYn       string
	EntpMemNo   string
	EntpKoNm    string
	UptDt       string
	BnrGrpSubSn string
}

type ContentClAddItem struct {
	BnrGrpSubSn string
	EntpMemNo   string
	EntpKoNm    string
	RegDt       string
}
