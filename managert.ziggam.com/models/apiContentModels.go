package models

type ContentItem struct {
	BnrGrpSn     string
	BnrGrpTypCd  string
	BnrGrpIdx    int64
	BnrGrpTitle  string
	BnrGrpSubCn  string
	PtoPath      string
	ThumbPtoPath string
	PublSdy      string
	PublEdy      string
	RdCnt        int64
	ShCnt        int64
	UseYn        string
	DelYn        string
	RegDt        string
	RegId        string
	UptDt        string
	UptId        string
	Expln        string
	RolTm        int64
	SubItem      []ContentSubItem
	EdtFlg		 string
}

type ContentSubItem struct {
	BnrGrpSN    string
	BnrGrpTypCd string
	BnrGrpIdx   int64
	ExtKey      string
	SwIdx       int64
	BnrSubCn01  string
	BnrSubCn02  string
}
