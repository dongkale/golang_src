package models

type LiveStat struct {
	WaitCnt int64
	IngCnt  int64
	EndCnt  int64
	CnclCnt int64
}

type LiveMemList struct {
	LmPpChrgGbnCd string
	LmPpChrgNm    string
	LmPpChrgBpNm  string
	LmChrgSn      string
}

type LiveMemList_v2 struct {
	LmPpChrgGbnCd string
	LmPpChrgNm    string
	LmPpChrgBpNm  string
	LmPpChrgSn    string
}

type LiveMemList_v3 struct {
	LmPpMemNo     string
	LmPpChrgGbnCd string
	LmPpChrgNm    string
	LmPpChrgBpNm  string
	LmPpChrgSn    string
	LmEntpGbn     string
	LmLiveStatCd  string
	LmPtoPath     string
}

type LiveList01 struct {
	S01RecrutSn     string
	S01PpMemNo      string
	S01LiveSn       string
	S01PtoPath      string
	S01Nm           string
	S01Sex          string
	S01Age          int64
	S01LiveItvSday  string
	S01LiveItvStime string
	S01LiveItvEday  string
	S01LiveItvEtime string
	S01TotCnt       int64
	S01LiveStatCd   string
	S01SubList      []LiveMemList
}

type RtnLiveList01 struct {
	RtnLiveList01Data []LiveList01
}

type LiveList02 struct {
	S02RecrutSn     string
	S02PpMemNo      string
	S02LiveSn       string
	S02PtoPath      string
	S02Nm           string
	S02Sex          string
	S02Age          int64
	S02LiveItvSday  string
	S02LiveItvStime string
	S02LiveItvEday  string
	S02LiveItvEtime string
	S02TotCnt       int64
	S02SubList      []LiveMemList
}

type RtnLiveList02 struct {
	RtnLiveList02Data []LiveList02
}

type LiveList03 struct {
	S03RecrutSn     string
	S03PpMemNo      string
	S03LiveSn       string
	S03PtoPath      string
	S03Nm           string
	S03Sex          string
	S03Age          int64
	S03LiveItvSday  string
	S03LiveItvStime string
	S03LiveItvEday  string
	S03LiveItvEtime string
	S03TotCnt       int64
	S03SubList      []LiveMemList
}

type RtnLiveList03 struct {
	RtnLiveList03Data []LiveList03
}

type LiveList struct {
	LsRecrutSn     string
	LsPpMemNo      string
	LsLiveSn       string
	LsPtoPath      string
	LsNm           string
	LsSex          string
	LsAge          int64
	LsLiveItvSday  string
	LsLiveItvStime string
	LsLiveItvEday  string
	LsLiveItvEtime string
	LsTotCnt       int64
	LsLiveStatCd   string
}

type RtnLiveList struct {
	RtnLiveListData []LiveList
}
