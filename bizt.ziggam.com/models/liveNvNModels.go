package models

import (
	"fmt"
)

// LiveNvNApplyMemList ...
type LiveNvNApplyMemList struct {
	RslPpMemNo string `json:"pp_mem_no"`
	RslNm      string `json:"name"`
	RslSex     string `json:"sex"`
	RslAge     int64  `json:"age"`
}

// LiveNvNSearchApply  ...
type LiveNvNSearchApply struct {
	LsaPpMemNo string
	LsaNm      string
	LsaSex     string
	LsaAge     string
}

// LiveNvNApplyList ...
type LiveNvNApplyList struct {
	LmLiveSchedStatCd string

	LmPpMemNo            string
	LmNm                 string
	LmSex                string
	LmAge                string
	LmPtoPath            string
	LmRecrutSn           string
	LmRecrutTitle        string
	LmLiveStatCd         string
	LmTRC04LiveReqStatCd string
	LmTRC04LiveSn        string
	LmMsgYn              string
	LmMsgEndYn           string
	LmReadEndDay         string
	LmItvLink            string
}

// GetApplyInfo ...
func (resp *LiveNvNApplyList) GetApplyInfo() (string, string, string, string) {
	return resp.LmPpMemNo, resp.LmNm, resp.LmSex, resp.LmAge
}

// GetApplyInfoStr ...
func (resp *LiveNvNApplyList) GetApplyInfoStr() string {
	return fmt.Sprintf("%v (%v・%v)[%v]", resp.LmNm, resp.LmSex, resp.LmAge, resp.LmPpMemNo)
}

// RtnLiveNvNApplyList ...
type RtnLiveNvNApplyList struct {
	RtnLiveNvNApplyListData []LiveNvNApplyList
}

// LiveNvNMemList ...
type LiveNvNMemList struct {
	LmLiveSchedStatCd string

	LmPpChrgGbnCd  string
	LmPpChrgNm     string
	LmPpChrgBpNm   string
	LmPpChrgSn     string
	LmPpLiveStatCd string
}

// LiveNvNList ...
type LiveNvNList struct {
	TotCnt         int64
	RecrutSn       string
	RecrutTitle    string
	LiveSn         string
	PtoPath        string
	LiveItvSday    string
	LiveItvStime   string
	LiveItvEday    string
	LiveItvEtime   string
	LiveItvRegDay  string
	LiveItvRegTime string
	LiveStatCd     string
	PpChrgNm       string
	PpChrgBpNm     string
	ApplyList      []LiveNvNApplyList
	MemList        []LiveNvNMemList
	Pagination     string
}

// RtnLiveNvNList ...
type RtnLiveNvNList struct {
	RtnLiveNvNListData        []LiveNvNList
	RtnLiveNvnListRecruitData []LiveNvnListRecruitList
	RtnGbnRecrutSn            string
	RtnGbnPpChrgSn            string

	RtnLiveNvnListApplyData []LiveNvNSearchApply
	RtnGbnPpMemNo           string

	RtnEntpTeamMemberListData []EntpTeamMemberList
}

// LiveNvNDetail ...
type LiveNvNDetail struct {
	RecrutSn     string
	RecrutTitle  string
	LiveSn       string
	LiveItvSdt   string
	LiveItvEdt   string
	LiveItvRegDt string
	LiveStatCd   string
	PpChrgNm     string
	PpChrgBpNm   string
	MemList      string
	ApplyList    string
	RegDt        string
}

// RtnNvnChannelId ...
type RtnNvnChannelId struct {
	EntpGbn   string
	PpMemNo   string
	PpChrgSn  string
	ChnnelId  string
	Nm        string
	BpNm      string
	EntpAdmin string
	PtoPath   string
}

// RtnLiveNvnJoinCnt ...
type RtnLiveNvnJoinCnt struct {
	Rtn     DefaultResult
	JoinCnt int64
}

type LiveNvnHistoryList struct {
	LhEntpMemNo   string
	LhRecrutSn    string
	LhLiveSn      string
	LhPpMemNo     string
	LhPpMemNm     string
	LhLiveStatCd  string
	LhMsgGbnCd    string
	LhMSgSn       string
	LhMsgYn       string
	LhMsgGbnNm    string
	LhMsgEndYn    string
	LhMemGbn      string
	LhMsgRegDtFmt string
	LhNMsgGbnCd   string
	LhPpChrgGbnCd string
	LhPpChrgSn    string
	LhPpChrgNm    string
	LhPpChrgBpNm  string
	LhRegDt       string
}

// LiveNvnRecruitApplyList ...
type LiveNvnRecruitApplyList struct {
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
	RslApplyRegCnt     int64
}

// RtnLiveNvnRecruitApplyList ...
type RtnLiveNvnRecruitApplyList struct {
	RtnLiveNvnRecruitApplyListData []LiveNvnRecruitApplyList
}

// LiveNvnEntpTeamMemberList ...
type LiveNvnEntpTeamMemberList struct {
	EtTotCnt      int64
	EtPpChrgSn    string
	EtPpChrgGbnCd string
	EtPpChrgNm    string
	EtPpChrgBpNm  string
	EtEmail       string
	EtEntpMemId   string
	EtPpChrgTelNo string
	EtRowNo       int64
	EtMemRegCnt   int64
}

// RtnLiveNvnEntpTeamMemberList ...
type RtnLiveNvnEntpTeamMemberList struct {
	RtnLiveNvnEntpTeamMemberListData []LiveNvnEntpTeamMemberList
}

// LiveNvnListRecruitList ...
type LiveNvnListRecruitList struct {
	RecrutSn    string
	RecrutTitle string
}

// RtnLiveNvnListRecruitList ...
type RtnLiveNvnListRecruitList struct {
	RtnLiveNvnListRecruitData []LiveNvnListRecruitList
}
