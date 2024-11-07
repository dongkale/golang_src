package tables

import (
	"encoding/json"
)

// MapEntpTypeCd :
var MapEntpTypeCd = map[string]string{
	"00": "미설정",
	"01": "대기업",
	"02": "중견기업",
	"03": "중소기업",
	"04": "공사/공기업",
	"05": "외국계기업",
	"06": "기타",
}

// MapBizTpyCd :
var MapBizTpyCd = map[string]string{
	"00": "미설정",
	"01": "서비스업",
	"02": "의료/제약/복지",
	"03": "제조/생산/화학",
	"04": "판매/유통",
	"05": "IT/정보통신",
	"06": "건설",
	"07": "교육",
	"08": "미디어/광고",
	"09": "은행/금융",
	"10": "문화/예술/디자인",
	"11": "기관/협회",
	"99": "기타",
}

// MapLstEduGbnCd1 :
var MapLstEduGbnCd1 = map[string]string{
	"00": "미설정",
	"01": "초등학교",
	"02": "중학교",
	"03": "고등학교",
	"04": "대학교 (2, 3년)",
	"05": "대학교 (4년)",
	"06": "대학교 (석사)",
	"07": "대학교 (박사)",
	"99": "기타",
}

// MapLstEduGbnCd2 :
var MapLstEduGbnCd2 = map[string]string{
	"00": "미설정",
	"01": "재학",
	"02": "수료",
	"03": "중퇴",
	"04": "자퇴",
	"05": "졸업예정",
	"06": "졸업",
}

// MapFrgnLangAbltCd1 : 외국어 능력 코드1
var MapFrgnLangAbltCd1 = map[string]string{
	"00": "미설정",
	"01": "영어",
	"02": "일본어",
	"03": "중국어",
	"04": "독일어",
	"05": "불어",
	"06": "스페인어",
	"07": "러시아어",
	"08": "이탈리아어",
	"09": "한국어",
	"10": "아랍어",
	"11": "베트남어",
	"12": "태국어",
	"13": "마인어",
	"14": "그리스어",
	"15": "포르투갈어",
	"16": "네덜란드어",
	"17": "힌디어",
	"18": "노르웨이어",
	"19": "유고어(세르비아)",
	"20": "히브리어",
	"21": "이란(페르시아어)",
	"22": "터키어",
	"23": "체코어",
	"24": "루마니아어",
	"25": "몽골어",
	"26": "스웨덴어",
	"27": "헝가리어",
	"28": "폴란드어",
	"29": "미얀마어",
	"30": "슬로바키아어",
	"31": "세르비아어",
	"99": "기타",
}

// MapFrgnLangAbltCd2 : 외국어 능력 코드2
var MapFrgnLangAbltCd2 = map[string]string{
	"00": "미설정",
	"01": "상(원어민 수준)",
	"02": "중(일상 회화)",
	"03": "하(비즈니스 회화)",
}

// MapCarrGbnCd : 경력 조건
var MapCarrGbnCd = map[string]string{
	"00": "미설정",
	"01": "무관",
	"02": "신입",
	"03": "경력",
	"99": "기타",
}

// MapEmplTypCd : 근무 형태
var MapEmplTypCd = map[string]string{
	"00": "미설정",
	"01": "정규직",
	"02": "계약직",
	"03": "인턴",
	"04": "파견직",
	"05": "도급",
	"06": "프리랜서",
	"07": "아르바이트",
	"08": "연수생・교육생",
	"09": "병역특례",
	"10": "위촉직・개인사업자",
	"99": "기타",
}

// MapLstEduGbnCd : 최종 학력
var MapLstEduGbnCd = map[string]string{
	"00": "미설정",
	"01": "무관",
	"02": "고졸",
	"03": "대졸(2, 3년)",
	"04": "대졸(4년)",
	"05": "대졸(석사)",
	"06": "대졸(박사)",
	"99": "기타",
}

// InviteSmsEmailSelectChar ... -> 같이 수정 invite_send_write.html : LIST_NONE_CHAR
var InviteSmsEmailSelectChar string = "-"

// LiveNvnApplyMaxCount : == static\js\defines\common_defines.js:LIVE_NVN_APPLY_MAX_CNT
var (
	LiveNvnApplyMaxCount int = 4
)

// TableConfig ...(== common_defines.js)
type TableConfig struct {
	LiveNvNApplyMaxCount     int // 인터뷰 가능한 지원자 수
	LiveNvNMemMaxCount       int // 인터뷰 가능한 기업 맴버 수
	LiveNvNViewApplyMaxCount int // 상세 페이지/리스트 페이지 보여지는 지원자 수
	LiveNvNViewMemMaxCount   int // 상세 페이지/리스트 페이지 보여지는 기업 맴버 수
	LiveNvNControlBeforeMin  int // 지원자 컨트롤(추가,삭제,취소) 시작 전 분
	LiveNvNListPagePerCount  int // 리스트에 보여지는 최대 갯수

	OneWayQAMaxLength int    // 원웨이 질문 길이 설정(원웨이 질문 길이, VARCHAR2(2000))
	OneWayQAMaxCount  int    // 원웨이 질문 갯수
	OneWayQASplitChar string // 원웨이 질문 split 문자

	OneWayQACntEntp    string // 기업별 원웨이 질문 갯수
	OneWayQACntRecruit string // 채용 공고별 원웨이 질문 갯수

	EvalItemGradeTbl__ string

	TestInt    int
	TestString string
}

// IsCheckValue ..
func (resp *TableConfig) IsCheckValue() string {
	if resp.LiveNvNApplyMaxCount <= 0 || resp.LiveNvNApplyMaxCount >= 50 {
		return "LiveNvNApplyMaxCount"
	}

	return ""
}

// ToString ...
func (resp *TableConfig) ToString() string {
	e, err := json.Marshal(resp)
	if err != nil {
		return ""
	}
	return string(e)
}

// TableConf ...
var TableConf TableConfig
