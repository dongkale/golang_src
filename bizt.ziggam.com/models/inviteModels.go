package models

// InviteMember ...
type InviteMember struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// InviteMemberList ...
type InviteMemberList struct {
	List []InviteMember
}

// RtnInvite ...
type RtnInvite struct {
	Name  string
	Email string
	Phone string

	IsEmailRefuse bool
	IsPhoneRefuse bool

	IsEmailSend bool
	IsPhoneSend bool
}

// InviteSendList ...
type InviteSendList struct {
	RslTotCnt      int64
	RslSendDt      string
	RslSendDtFmt   string
	RslRecrutSn    string
	RslRecrutTitle string
	RslSenderName  string
	RslCnt         int64
	Pagination     string
}

// RtnInviteSendList ...
type RtnInviteSendList struct {
	RtnInviteSendListData []InviteSendList
}

// InviteSendListDetail ...
type InviteSendListDetail struct {
	RslTotCnt      int64
	RslRowNo       int64
	RslSendDt      string
	RslRecrutSn    string
	RslName        string
	RslEmail       string
	RslPhone       string
	RslEmailMid    string
	RslEmailResult string
	RslEmailDt     string
	RslSmsMid      string
	RslSmsResult   string
	RslSmsDt       string
	RslListYN      string
	Pagination     string
}

// RtnInviteSendListDetail ...
type RtnInviteSendListDetail struct {
	RtnInviteSendListDetailData []InviteSendListDetail
}

// InviteSendType ...
type InviteSendType int

// InviteSendType ...
const (
	InviteSendTypeMailSms InviteSendType = 0
	InviteSendTypeMail    InviteSendType = 1
	InviteSendTypeSms     InviteSendType = 2
)
