package models

type AdminStatsMain struct {
	TotMemCnt  int64
	TdMemCnt   int64
	WdMemCnt   int64
	TotEntpCnt int64
	TdEntpCnt  int64
	WdEntpCnt  int64
	VpMemCnt   int64
	TotVpCnt   int64
	FVpCnt     int64
	PVpCnt     int64
	CVpCnt     int64
	AdMemCnt   int64
	IsMemCnt   int64
	MMemCnt    int64
	FMemCnt    int64
	TotEvpCnt  int64
}

type AdminStatsSub1 struct {
	P10U   int64
	P16_20 int64
	P21_25 int64
	P26_30 int64
	P31_35 int64
	P36_40 int64
	P41_45 int64
	P46_50 int64
	P51_55 int64
	P56_60 int64
	P61_65 int64
	P66_70 int64
	P70O   int64
}

type AdminStatsSub2 struct {
	EntpMemNo string
	EntpKoNm  string
	PvCnt     int64
	UvCnt     int64
}

type AdminStatsSub3 struct {
	LstEntpMemNo string
	LstEntpKoNm  string
}

type RtnAdminStatsSub2 struct {
	RtnAdminStatsSub2Data []AdminStatsSub2
}
