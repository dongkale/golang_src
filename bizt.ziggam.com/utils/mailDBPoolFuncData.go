package utils

import (
	"fmt"
	"strings"
)

type MailDBPoolFunc struct {
	paramChar      string
	paramValueChar string
}

func (resp *MailDBPoolFunc) Init() {
	resp.paramChar = ";"
	resp.paramValueChar = "="
}

func (resp *MailDBPoolFunc) ParamParser(paramStrings string) map[string]string {
	entries := strings.Split(paramStrings, resp.paramChar)

	m := make(map[string]string)
	for _, e := range entries {
		parts := strings.Split(e, resp.paramValueChar)
		if len(parts) > 1 {
			m[parts[0]] = parts[1]
		} else {
			m[parts[0]] = ""
		}
	}

	return m
	//return utils.StringDelimSplit(paramStrings, resp.paramChar, resp.paramValueChar)
}

// 함수명 대문자 !!!!!
func (resp *MailDBPoolFunc) ContentsTemplate(params map[string]string) string {
	arg1 := params["arg1"]
	arg2 := params["arg2"]

	var htmlContents = fmt.Sprintf(" HTML ==> arg1:%v, arg2:%v", arg1, arg2)

	return htmlContents
}

// users function...
func (resp *MailDBPoolFunc) HtmlContents(params map[string]string) string {
	arg1 := params["arg1"]
	arg2 := params["arg2"]
	arg3 := params["arg3"]

	//var htmlContents = fmt.Sprintf(" HTML ==> arg1:%v, arg2:%v, arg3:%v", arg1, arg2, arg3)
	var htmlContents = fmt.Sprintf("<html><body>Hello....<p>%v</p><p>%v</p><p>%v</p></body></html>", arg1, arg2, arg3)

	return htmlContents
}

var MailDBPoolFuncList MailDBPoolFunc
