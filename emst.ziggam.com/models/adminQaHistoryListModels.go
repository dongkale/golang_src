package models

type QaHistoryList struct {
	QstSn      string
	VdTitle    string
	VdFilePath string
}

type RtnQaHistoryList struct {
	RtnQaHistoryListData []QaHistoryList
}
