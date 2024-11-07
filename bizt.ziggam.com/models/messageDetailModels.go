package models

type MessageDetailTopInfo struct {
	EvlPrgsStatCd string
	EvlPrgsStatNm string
	FavrAplyPpYn  string
	PtoPath       string
	Nm            string
	Sex           string
	Age           int64
	LiveReqStatCd string
	UpJobGrp      string
	JobGrp        string
	RecrutTitle   string
	RegDt         string
}

type MessageList struct {
	MlEntpMemNo     string
	MlRecrutSn      string
	MlPpMemNo       string
	MlMsgSn         string
	MlPtoPath       string
	MlLdYn          string
	MlLdDt          string
	MlLdDt2         string
	MlMemGbn        string
	MlMsgGbnCd      string
	MlMsgCont       string
	MlMsgCfrmYn     string
	MlEntpNm        string
	MlMemNm         string
	MlLiveItvSday   string
	MlLiveItvStime  string
	MlLiveItvStime2 string
	MlLiveItvEday   string
	MlLiveItvEtime  string
	MlLiveItvEtime2 string
	MlLiveItvJt     string
	MlLiveSn        string
	MlLiveNvnYn     string
	MlMsgCnt        int64
}

type RtnMessageList struct {
	RtnMessageListData []MessageList
}
