package models

type AdminStatsRC01 struct {
	Rc01TotCnt int64
	Rc01NewCnt int64
	Rc01IngCnt int64
	Rc01EndCnt int64
}

type AdminStatsRC02 struct {
	Rc02TotCnt  int64
	Rc02IngCnt  int64
	Rc02PassCnt int64
	Rc02FailCnt int64
}

type AdminStatsRC03 struct {
	Rc03EntpMemNo    string
	Rc03EntpKoNm     string
	Rc03ApplyCnt     int64
	Rc03PassCnt      int64
	Rc03FailCnt      int64
	Rc03IngCnt       int64
	Rc03MatchingRate float64
}

type AdminStatsRCSub struct {
	SubEntpMemNo string
	SubEntpKoNm  string
}

type RtnAdminStatsRC03 struct {
	RtnAdminStatsRC03Data []AdminStatsRC03
}
