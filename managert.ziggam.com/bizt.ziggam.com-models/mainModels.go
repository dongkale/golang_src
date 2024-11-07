package models

type MainInfo struct {
	MnPpChrgBpNm string
	MnPpChrgNm   string
	MnEntpKoNm   string
	MnInfoCnt    int64
	MnVdCnt      int64
}

type MainStat struct {
	VideoTodayCnt int64
	VideoTotCnt   int64
	RecrutIngCnt  int64
	RecrutTotCnt  int64
	ApplyTodayCnt int64
	ApplyTotCnt   int64
}

type MainLiveList struct {
	MnlRecrutSn  string
	MnlPpMemNo   string
	MnlLiveSn    string
	MnlLiveItvSd string
	MnlLiveItvSt string
	MnlPtoPath   string
	MnlNm        string
	MnlSex       string
	MnlAge       int64
}

type MainLiveNvNListApply struct {
	MlnLiveStatCd string
	MlnRecrutSn   string
	MlnPpMemNo    string
	MlnNm         string
	MlnSex        string
	MlnAge        int64
}

type MainLiveNvNList struct {
	MlnRecrutSn          string
	MlnLiveSn            string
	MlnLiveItvSd         string
	MlnLiveItvSt         string
	MainLiveNvNListApply []MainLiveNvNListApply
}

type MainRecruitList struct {
	MnrEntpMemNo   string
	MnrRecrutSn    string
	MnrRecrutTitle string
	MnrNewCnt      int64
	MnrDdy         string
}

type ApplicantList struct {
	MnaRecrutSn     string
	MnaPpMemNo      string
	MnaFavrAplyPpYn string
	MnaPtoPath      string
	MnaNm           string
	MnaSex          string
	MnaAge          int64
	MnaRegDt        string
}

type RtnApplicantList struct {
	RtnApplicantListData []ApplicantList
}

type MainNotice struct {
	MnnSn    int64
	MnnTitle string
	MnnRegDt string
}
