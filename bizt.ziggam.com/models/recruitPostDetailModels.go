package models

type RecruitStatTopInfo struct {
	ApplyCnt     int64
	IngCnt       int64
	PassCnt      int64
	FailCnt      int64
	DcmntPassCnt int64
	DcmntFailCnt int64
}

type RecruitStatList struct {
	RslTotCnt          int64
	RslEntpMemNo       string
	RslRecrutSn        string
	RslFavrAplyPpYn    string
	RslNm              string
	RslSex             string
	RslAge             string
	RslRegDt           string
	RslApplyDt         string
	RslEvlStatDt       string
	RslEvlPrgsStatCd   string
	RslRcrtAplyStatCd  string
	RslEntpCfrmYn      string
	RslVpYn            string
	RslPpMemNo         string
	RslLiveReqStatCd   string
	RslRowNo           int64
	RslPtoPath         string
	DcmntEvlStatCd     string
	OnwyIntrvEvlStatCd string
	LiveIntrvEvlStatCd string
	ReadEndDay         string
	RslZOrder          string
	RslScoreValue      int64

	RslRecrutTitle string

	RslApplyCnt int64

	Pagination string
}

type RtnRecruitStatList struct {
	RtnRecruitStatListData []RecruitStatList

	ApplySortCd  string
	ApplySortWay string
}

type RecruitStatListAll struct {
	RslTotCnt          int64
	RslEntpMemNo       string
	RslRecrutSn        string
	RslFavrAplyPpYn    string
	RslNm              string
	RslSex             string
	RslAge             string
	RslRegDt           string
	RslApplyDt         string
	RslEvlStatDt       string
	RslEvlPrgsStatCd   string
	RslRcrtAplyStatCd  string
	RslEntpCfrmYn      string
	RslVpYn            string
	RslPpMemNo         string
	RslLiveReqStatCd   string
	RslRowNo           int64
	RslPtoPath         string
	DcmntEvlStatCd     string
	OnwyIntrvEvlStatCd string
	LiveIntrvEvlStatCd string
	ReadEndDay         string
	RslZOrder          string
	RslScoreValue      int64

	RslRecrutTitle string

	RslApplyCnt int64

	RslRecruitApplyMemberAnswerList []RecruitApplyMemberAnswerList

	Pagination string
}

type RtnRecruitStatListAll struct {
	Rtn DefaultResult

	RtnRecruitStatListData []RecruitStatListAll

	ApplySortCd  string
	ApplySortWay string
}
