package models

type ApplicantScore struct {
	EntpMemNo string
	RecrutSn  string
	PpMemNo   string

	EvalItemString string
	EvalItem       map[string]string
	ResultComment  string
}

type RslApplicantScore struct {
	Rtn           DefaultResult
	RslDetailInfo []ApplicantDetailPopup
	RslScorreInfo []ApplicantScore
}

type ApplicantScoreEvalItem struct {
	Num      int64
	Title    string
	Category string
}

type ApplicantScoreEvalItemCategory struct {
	Category string
	Title    string
}
