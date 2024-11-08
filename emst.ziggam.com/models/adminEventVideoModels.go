package models

type AdminEventVideo struct {
	TotCnt     int64
	Sn         int64
	RegDt      string
	PpMemNo    string
	VdSn       string
	VdFilePath string
	ThmKndCd   string
	QstCd      string
	ThmDesc    string
	VdSec      int64
	OpnSetCd   string
	CVdSn      string
}

type AdminEventView struct {
	VRegDt      string
	VPpMemNo    string
	VVdSn       string
	VVdFilePath string
	VThmKndCd   string
	VQstCd      string
	VThmDesc    string
	VVdSec      int64
	VOpnSetCd   string
}
