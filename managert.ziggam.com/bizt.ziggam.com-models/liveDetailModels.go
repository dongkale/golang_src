package models

type LiveDetail struct {
	EvlPrgsStatCd string
	LiveReqStatCd string
	RegDt         string
	PtoPath       string
	Nm            string
	Sex           string
	Age           int64
	UpJobGrp      string
	JobGrp        string
	RecrutTitle   string
	LiveSn        string
}

type LiveHistoryList struct {
	LhEntpMemNo       string
	LhRecrutSn        string
	LhPpMemNo         string
	LhLiveStatCd      string
	LhMsgGbnCd        string
	LhMSgSn           string
	LhMsgYn           string
	LhLiveSn          string
	LhLiveItvSday     string
	LhLiveItvSTime    string
	LhLiveItvEday     string
	LhLiveItvETime    string
	LhMsgGbnNm        string
	LhMsgEndYn        string
	LhMemGbn          string
	LhLiveItvCnclDay  string
	LhLiveItvCnclTime string
	LhNMsgGbnCd       string
}

type LiveHistoryInfo struct {
	LhiEntpMemNo    string
	LhiRecrutSn     string
	LhiPpMemNo      string
	LhiLiveStatCd   string
	LhiMsgGbnCd     string
	LhiMSgSn        string
	LhiMsgYn        string
	LhiLiveSn       string
	LhiLiveItvSday  string
	LhiLiveItvSTime string
	LhiLiveItvEday  string
	LhiLiveItvETime string
	LhiMsgGbnNm     string
	LhiMsgEndYn     string
}

type LiveInvResult struct {
	LirEntpMemNo    string
	LirRecrutSn     string
	LirPpMemNo      string
	LirPpMemNm      string
	LirLiveSn       string
	LirEntpNm       string
	LirNm           string
	LirLiveItvSday  string
	LirLiveItvSTime string
	LirLiveItvEday  string
	LirLiveItvETime string
	LirLiveItvJt    string
	LirSdtTstmp     string
}

// <--
