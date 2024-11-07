package models

type MessageTopInfo struct {
	IngCnt int64
	EndCnt int64
}

type MessageMemberList struct {
	MmEntpMemNo    string
	MmRecrutSn     string
	MmPpMemNo      string
	MmMsgSn        string
	MmMsgCont      string
	MmRegDt        string
	MmMsgCfrmYn    string
	MmMsgGbnCd     string
	MmPtoPath      string
	MmNm           string
	MmTotCnt       int64
	MmLastMsgGbnCd string
	MmMsgCnt	   int64
	Target         interface{}
}

type RtnMessageMemberList struct {
	RtnMessageMemberListData []MessageMemberList
}
