package models

type MessageListCount struct {
	MsgCnt int64
}

type RtnMessageListCount struct {
	RtnMessageListCountData []MessageListCount
}
