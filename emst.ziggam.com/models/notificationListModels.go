package models

type NotificationList struct {
	GrpNo     int64
	MemNo     string
	RegDy     string
	RegHm     string
	PtoPath   string
	NotiKndCd string
	NotiCont  string
	Sex       string
	Nm        string
	Age       string
	EntpMemNo string
	RecrutSn  string
	UpJobGrp  string
	JobGrp    string
	CfrmDt    string
	Sn        int64
	RegDt     string
	NextGrpNo int64
	PrevGrpNo int64
	PpMemNo   string
	NewYn     string

	Pagination string
	MenuId     string
}

type RtnNotificationList struct {
	RtnNotificationListData []NotificationList
}
