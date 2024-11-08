package models

type MainNotiList struct {
	NotiRegDt string
	NotiTitle string
	NotiSn    int64
	NotiNewYn string
}

type MainStat struct {
	VideoTodaycnt int64
	VideoTotCnt   int64
	RecrutIngCnt  int64
	RecrutTotCnt  int64
	ApplyTodayCnt int64
	ApplyTotCnt   int64
}

type MainRecruitList struct {
	RcEntpMemNo   string
	RcRecrutSn    string
	RcRecrutTitle string
	RcEmplTyp     string
	RcUpJobGrp    string
	RcJobGrp      string
	RcRecrutDy    string
	RcRegDt       string
	RcNewCnt      int64
}

type MainApplyList struct {
	ApEntpMemNo   string
	ApRecrutSn    string
	ApPpMemNo     string
	ApPtoPath     string
	ApNm          string
	ApSex         string
	ApRecrutTitle string
	ApLeftDy      string
	ApRedYn       string
}
