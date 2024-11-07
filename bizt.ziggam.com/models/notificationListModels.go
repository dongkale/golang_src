package models

type NotificationList struct {
	NtMemNo     string
	NtDvcId     string
	NtNotiKndCd string
	NtRegDt     string
	NtCfrmDt    string
	NtDtDd      string
	NtDtHh      string
	NtPtoPath   string
	NtNotiCont  string
	NtEntpMemNo string
	NtRecrutSn  string
	NtPpMemNo   string
	NtLiveSn    string
	NtBrdGbnCd  string
	NtSn        int64
	NtLdYn      string
	NtInqGbnCd  string
	NtInqRegDt  string
	NtMsgEndYn  string
	NtLiveNvnYn string
}

type RtnNotificationList struct {
	RtnNotificationListData []NotificationList
}
