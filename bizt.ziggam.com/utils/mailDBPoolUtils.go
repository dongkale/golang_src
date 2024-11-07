package utils

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/rana/ora.v4"
)

type MailDBPool_Type int

const (
	MailDBPool_Mail_Sender MailDBPool_Type = 0
	MailDBPool_Mail_Daemon MailDBPool_Type = 1
	MailDBPool_Mail_Pool   MailDBPool_Type = 2
)

const (
	MailDBPool_Error_ReTry = "Try again later"
)

type MailDBPoolData struct {
	Sn int64

	ToNm    string
	ToEmail string

	FromNm    string
	FromEmail string

	Subject string

	ContentsParam string
	ContentsFunc  string
}

type MailDBPool struct {
	LoopTimeSec int64
	IsQuit      chan bool

	SelectCount int // 메일 리스트 가져올 갯수

	//SelectSkipMin int // 메일 리스트 가져올때 지정한 시간(등록시간) 지난 것은 skip

	WhereSn int64 // select 기준이 되는 Sn 번호

	SelectZeroCount    int64 // DB Table 갯수 0 체크
	SelectZeroMaxCount int64 // DB Table 갯수 0 체크:  WhereSn 초기화 목적으로 DB 테이블 처음 부터 다시 읽는 기능

	ProcType MailDBPool_Type // 0: MailSender, 1: MailDaemon, 2: MailPool

	LoopCount      int64 // Loop Count
	LoopDelayCount int64 // Loop Skip Count

	MailSendCnter MailSendCounter
}

func (resp *MailDBPool) SetSelectCount(selCount int) {
	resp.SelectCount = selCount
}

// func (resp *MailDBPool) SetSelectSkipMin(skipMinute int) {
// 	resp.SelectSkipMin = skipMinute
// }

// ReInit ...
func (resp *MailDBPool) ReInit(procType MailDBPool_Type, loopTimeSec int64, selectCount int) {

	resp.IsQuit <- true
	resp.Init(procType, loopTimeSec, selectCount)
}

// Init ...
func (resp *MailDBPool) Init(procType MailDBPool_Type, loopTimeSec int64, selectCount int) {
	// for {
	// 	select {
	// 	// case <-stop:
	// 	// 	fmt.Println("EXIT: 3 seconds")
	// 	// 	return
	// 	case <-time.After(5 * time.Second):
	// 		fmt.Printf("5 second")
	// 		break
	// 	case <-time.After(1 * time.Second):
	// 		fmt.Printf("1 second")
	// 		break
	// 	}
	// }

	// go func() {
	// 	loopTick := time.Tick(resp.LoopTimeSec * time.Second)

	// 	for now := range loopTick {
	// 		ret, isExit := resp.DoLoop(now)
	// 		if ret != nil {
	// 			fmt.Printf(fmt.Sprintf("[MailDBPool][Error] %v", ret))
	// 		}

	// 		if isExit == true {
	// 			fmt.Printf(fmt.Sprintf("[MailDBPool][Error] Exit !!"))
	// 			return
	// 		}
	// 	}

	// 	// stop := time.After(3 * time.Second)
	// 	// for {
	// 	// 	select {
	// 	// 	case <-stop:
	// 	// 		fmt.Println("EXIT: 3 seconds")
	// 	// 		return
	// 	// 	case <-time.After(1 * time.Second):
	// 	// 		fmt.Println(i, "second")
	// 	// 	}
	// 	// }
	// }()

	resp.LoopTimeSec = loopTimeSec
	resp.SelectCount = selectCount
	resp.ProcType = procType // 0: MailSender, 1: MailDaemon, 2: MailPool
	resp.WhereSn = 0         // select 할 sn 기준값
	resp.LoopCount = 0
	resp.LoopDelayCount = 0
	resp.SelectZeroCount = 0
	resp.SelectZeroMaxCount = 30

	//resp.LoopTimeOut = time.Now()

	ticker := time.NewTicker(time.Duration(resp.LoopTimeSec) * time.Second)
	resp.IsQuit = make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				ret, isExit := resp.DoLoop(time.Now())
				if ret != nil {
					logs.Error(fmt.Sprintf("[MailDBPool][Error] %v", ret))
				}

				if isExit == true {
					fmt.Printf(fmt.Sprintf("[MailDBPool] Exit !!"))
					return
				}
			case <-resp.IsQuit:
				fmt.Printf(fmt.Sprintf("[MailDBPool] Quit !!"))
				ticker.Stop()
				return
			}
		}
	}()

	resp.MailSendCnter.Load()

	MailDBPoolFuncList.Init()

	fmt.Printf(fmt.Sprintf("[MailDBPool][Init] LoopTimeSec:%v sec, SelectCount:%v, ProcType:%v, SelectZeroMaxCount:%v(%v sec)",
		resp.LoopTimeSec, resp.SelectCount, resp.ProcType, resp.SelectZeroMaxCount, resp.LoopTimeSec*resp.SelectZeroMaxCount))
}

// func (resp *MailDBPool) Add(toMail MailDBPoolData) {

// }

func (resp *MailDBPool) AddArray(toMailArr []MailDBPoolData) error {

	var toName []string
	var toEmail []string
	var fromName []string
	var fromEmail []string
	var subject []string
	var contentsParam []string
	var contentsFunc []string

	for _, val := range toMailArr {
		toName = append(toName, val.ToNm)
		toEmail = append(toEmail, val.ToEmail)
		fromName = append(fromName, val.FromNm)
		fromEmail = append(fromEmail, val.FromEmail)
		subject = append(subject, val.Subject)
		contentsParam = append(contentsParam, val.ContentsParam)
		contentsFunc = append(contentsFunc, val.ContentsFunc)
	}

	resultDB := resp.addDB(toName, toEmail, fromName, fromEmail, subject, contentsParam, contentsFunc)
	if resultDB != nil {
		logs.Error(fmt.Sprintf("[MailDBPool][AddArray][Error] addDB() : %v", resultDB))
	}

	return resultDB
}

func (resp *MailDBPool) DoLoop(now time.Time) (error, bool) {
	resp.LoopCount++

	resp.MailSendCnter.CntCheck(GetSmtpServer())

	fmt.Printf(fmt.Sprintf("[MailDBPool][DoLoop][%v] %v", resp.LoopCount, now.Format(("2006-01-02 15:04:05"))))

	if resp.LoopCount < resp.LoopDelayCount {
		fmt.Printf(fmt.Sprintf("[MailDBPool][DoLoop] Skip[%d:%d] [%s] ", resp.LoopCount, resp.LoopDelayCount, now.Format("2006-01-02 15:04:05")))
		return nil, false
	}

	err, isExit := resp.selectProcess(now)

	return err, isExit
}

// error, bool : 종료 여부
func (resp *MailDBPool) selectProcess(now time.Time) (error, bool) {

	fmt.Printf(fmt.Sprintf("[MailDBPool][selectProcess] ..."))

	pLang, _ := beego.AppConfig.String("lang")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		return err, false
	}
	// End : Oracle DB Connection

	offset := 0
	limit := resp.SelectCount
	whereSn := resp.WhereSn

	// Start : ZSP_MAIL_DB_POOL_LIST
	fmt.Printf(fmt.Sprintf("CALL ZSP_MAIL_DB_POOL_LIST('%v', %v, %v, %v, :1)", pLang, offset, limit, whereSn))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_MAIL_DB_POOL_LIST('%v', %v, %v, %v, :1)", pLang, offset, limit, whereSn),
		ora.I64, /* TOT_CNT */
		ora.I64, /* SN */
		ora.S,   /* TO_NM */
		ora.S,   /* TO_EMAIL */
		ora.S,   /* FROM_NM */
		ora.S,   /* FROM_EMAIL */
		ora.S,   /* SUBJECT */
		ora.S,   /* CONTENTS_PARAM */
		ora.S,   /* CONTENTS_FUNC */
	)
	defer stmtProcCall.Close()
	if err != nil {
		return err, false
	}
	procRset := &ora.Rset{}
	_, err = stmtProcCall.Exe(procRset)

	if err != nil {
		return err, false
	}

	var checkSn []string
	mailDbPoolList := make([]MailDBPoolData, 0)

	var (
		totCnt        int64
		sn            int64
		toNm          string
		toEmail       string
		fromNm        string
		fromEmail     string
		subject       string
		contentsParam string
		contentsFunc  string
	)

	if procRset.IsOpen() {
		for procRset.Next() {

			totCnt = procRset.Row[0].(int64)
			sn = procRset.Row[1].(int64)
			toNm = procRset.Row[2].(string)
			toEmail = procRset.Row[3].(string)
			fromNm = procRset.Row[4].(string)
			fromEmail = procRset.Row[5].(string)
			subject = procRset.Row[6].(string)
			contentsParam = procRset.Row[7].(string)
			contentsFunc = procRset.Row[8].(string)

			mailDbPoolList = append(mailDbPoolList, MailDBPoolData{
				Sn:            sn,
				ToNm:          toNm,
				ToEmail:       toEmail,
				FromNm:        fromNm,
				FromEmail:     fromEmail,
				Subject:       subject,
				ContentsParam: contentsParam,
				ContentsFunc:  contentsFunc,
			})

			checkSn = append(checkSn, fmt.Sprintf("%v", sn))

			resp.WhereSn = sn
		}
		if err := procRset.Err(); err != nil {
			return err, false
		}
	}
	// End : ZSP_MAIL_DB_POOL_LIST

	// 유휴 시간으로 체크!!
	if totCnt == 0 {
		resp.SelectZeroCount++

		if resp.SelectZeroCount >= resp.SelectZeroMaxCount {
			resp.WhereSn = 0
			resp.SelectZeroCount = 0
			fmt.Printf(fmt.Sprintf("[MailDBPool][selectProcess] Reset ==> WhereSn:%v, SelectZeroCount:%v", resp.WhereSn, resp.SelectZeroCount))
		}
	}

	fmt.Printf(fmt.Sprintf("[MailDBPool][selectProcess] TotCnt:%v, WhereSn:%v, SelectZeroCount:%v", totCnt, resp.WhereSn, resp.SelectZeroCount))
	//DebugLog("[MailDBPool][selectProcess] TotCnt:%d, WhereSn:%d, SelectZeroCount:%d", totCnt, resp.WhereSn, resp.SelectZeroCount)

	if len(mailDbPoolList) > 0 {

		fmt.Printf(fmt.Sprintf("[MailDBPool][selectProcess] ListSn ==> [%v]", strings.Join(checkSn, ",")))

		for _, val := range mailDbPoolList {
			fmt.Printf(fmt.Sprintf("[MailDBPool][selectProcess] List ==> Sn:%v, ToNm:%v, ToMail:%v, FromNm:%v, FromMail:%v,  Arg:%v, Arg_Func:%v", val.Sn, val.ToNm, val.ToEmail, val.FromNm, val.FromEmail, val.ContentsParam, val.ContentsFunc))
		}

		if resp.ProcType == 0 {

			// err := MailSender.Connect(GetSmtpServer(), GetSmtpServerPort(), GetReturnEmail(), GetReturnEmailPwd())
			// if err != nil {
			// 	logs.Error(fmt.Sprintf("[MailDBPool][selectProcess][Error] MailSender.Connect() : %v", err))
			// 	return nil, false
			// }

			var snList []string
			var toNameList []string
			var toMailList []string

			for _, val := range mailDbPoolList {

				err := MailSender.Connect(GetSmtpServer(), GetSmtpServerPort(), GetReturnEmail(), GetReturnEmailPwd())
				if err != nil {
					logs.Error(fmt.Sprintf("[MailDBPool][selectProcess][Error] MailSender.Connect() : %v", err))
					break //return nil, false
				}

				to := val.ToEmail
				toName := val.ToNm
				from := val.FromEmail
				fromName := val.FromNm
				subject := val.Subject
				htmlContents, err := Invoke(&MailDBPoolFuncList, val.ContentsFunc, MailDBPoolFuncList.ParamParser(val.ContentsParam))
				if err != nil {
					logs.Error(fmt.Sprintf("[MailDBPool][selectProcess][Error] Invoke() : %v", err))
					continue
				}

				err = MailSender.Send(to, toName, from, fromName, subject, htmlContents.Interface().(string))
				if err != nil {
					logs.Error(fmt.Sprintf("[MailDBPool][selectProcess][Error] MailSender.Send() : %v", err))

					if strings.Contains(string(err.Error()), MailDBPool_Error_ReTry) == true {
						resp.LoopDelayCount = resp.LoopCount + 10 //(10 count == LoopTimeSec * 10 --> time second )
						logs.Error(fmt.Sprintf("[MailDBPool][selectProcess][Set] LoopDelayCount: %v", resp.LoopDelayCount))
						//return nil, false
						break
					}

				} else {
					fmt.Printf(fmt.Sprintf("[MailDBPool][selectProcess][Send] MailSender.Send() : %v", err))

					snList = append(snList, fmt.Sprintf("%v", val.Sn))
					toNameList = append(toNameList, toName)
					toMailList = append(toMailList, to)
				}
			}

			resultDB := resp.updateDB(snList, toNameList, toMailList)
			if resultDB != nil {
				logs.Error(fmt.Sprintf("[MailDBPool][selectProcess][Error] updateDB() : %v", resultDB))
			}
		} else if resp.ProcType == 1 {
			var arrayData []SendMailDaemonPushData

			for i, val := range mailDbPoolList {

				to := val.ToEmail
				toName := val.ToNm
				from := val.FromEmail
				fromName := val.FromNm
				subject := val.Subject
				htmlContents, err := Invoke(&MailDBPoolFuncList, val.ContentsFunc, MailDBPoolFuncList.ParamParser(val.ContentsParam))
				if err != nil {
					logs.Error(fmt.Sprintf("[MailDBPool][selectProcess][Error] Invalid ContentsFunc : %v", err))
					continue
				}

				arrayData = append(arrayData, SendMailDaemonPushData{
					ToSend: fmt.Sprintf("%s[%d]:%s", toName, i, to),
					ToSendData: SendMailDaemonData{
						To:           to,
						ToName:       toName,
						From:         from,
						FromName:     fromName,
						Subject:      subject,
						HtmlContents: htmlContents.Interface().(string),
					},
					Cb:     resp.callbackSend_Daemon,
					CbData: fmt.Sprintf("Sn=%v;ToName=%v;ToMail=%v", val.Sn, toName, to),
				})
			}

			err = SendMailDaemonMng.PushArray(arrayData, 0)
			if err != nil {
				logs.Error(fmt.Sprintf("[MailDBPool][selectProcess][Error] Daemon.PushArray() : %v", err))
			}
		} else {

			var arrayData []SendMailPoolPushData

			for i, val := range mailDbPoolList {
				to := val.ToEmail
				toName := val.ToNm
				from := val.FromEmail
				fromName := val.FromNm
				subject := val.Subject
				htmlContents, err := Invoke(&MailDBPoolFuncList, val.ContentsFunc, MailDBPoolFuncList.ParamParser(val.ContentsParam))
				if err != nil {
					logs.Error(fmt.Sprintf("[MailDBPool][selectProcess][Error] Invalid ContentsFunc : %v", err))
					continue
				}

				arrayData = append(arrayData, SendMailPoolPushData{
					ToSend: fmt.Sprintf("%s[%d]:%s", toName, i, to),
					ToSendData: SendMailPoolData{
						To:           []string{SendMailPoolDestFmt(toName, to)},
						From:         SendMailPoolDestFmt(fromName, from),
						Subject:      subject,
						HtmlContents: htmlContents.Interface().(string),
					},
					Cb:     resp.callbackSend,
					CbData: fmt.Sprintf("Sn=%v;ToName=%v;ToMail=%v", val.Sn, toName, to),
				})
			}

			//err = SendMailPoolMng.PushArray(arrayData, 0)
			err = SendMailPoolMng.PushArrayEx(arrayData, 3, resp.CompleteProcess, fmt.Sprintf("[Time:%v, Cnt:%v]", now.Format("2006-01-02 15:04:05"), len(arrayData)))
			if err != nil {
				logs.Error(fmt.Sprintf("[MailDBPool][selectProcess][Error] Pool.PushArray() : %v", err))
			}
		}
	}

	return nil, false
}

func (resp *MailDBPool) callbackSend_Daemon(err error, cnt int, toSend string, cbData string) {
	resp.callbackSend(err, 0, cnt, toSend, cbData)
}

func (resp *MailDBPool) callbackSend(err error, num int, cnt int, toSend string, cbData string) {

	fmt.Printf(fmt.Sprintf("[MailDBPool][callbackSend] [num:%v, cnt:%v] -> ToSend:%v, CbData: %v, Error:%v", num, cnt, toSend, cbData, err))

	if err != nil {
		logs.Error(fmt.Sprintf("[MailDBPool][callbackSend][Error] ToSend:%v, CbData: %v --> %v", toSend, cbData, err))
		return
	}

	errDB := resp.MailSendCnter.SetCnt(GetSmtpServer(), 1)
	if errDB != nil {
		logs.Error(fmt.Sprintf("[MailDBPool][callbackSend][Error] ToSend:%v, CbData: %v --> SetCnt() : %v", toSend, cbData, errDB))
	}

	result := StringDelimSplit(cbData, ";", "=")

	var snList []string
	var toNameList []string
	var toMailList []string

	snList = append(snList, result["Sn"])
	toNameList = append(toNameList, result["ToName"])
	toMailList = append(toMailList, result["ToMail"])

	resultDB := resp.updateDB(snList, toNameList, toMailList)
	if resultDB != nil {
		logs.Error(fmt.Sprintf("[MailDBPool][callbackSend][Error] updateDB() : %v", resultDB))
	}

	return
}

func (resp *MailDBPool) updateDB(snList []string, toNameList []string, toMailList []string) error {

	pLang, _ := beego.AppConfig.String("lang")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		return err
	}
	// End : Oracle DB Connection

	// Start : ZSP_MAIL_DB_POOL_UPT
	fmt.Printf(fmt.Sprintf("CALL ZSP_MAIL_DB_POOL_UPT('%v', %v, '%v', '%v', '%v', :1)",
		pLang,
		len(snList),
		strings.Join(snList, ","),
		strings.Join(toNameList, ","),
		strings.Join(toMailList, ",")))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_MAIL_DB_POOL_UPT('%v', %v, '%v', '%v', '%v', :1)",
		pLang,
		len(snList),
		strings.Join(snList, ","),
		strings.Join(toNameList, ","),
		strings.Join(toMailList, ",")),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
	)
	defer stmtProcCall.Close()
	if err != nil {
		return err
	}
	procRset := &ora.Rset{}
	_, err = stmtProcCall.Exe(procRset)

	if err != nil {
		return err
	}

	var (
		rtnCd  int64
		rtnMsg string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}

		if err := procRset.Err(); err != nil {
			return err
		}

		fmt.Printf(fmt.Sprintf(" ===> rtnCd:%v, rtnMsg:%v", rtnCd, rtnMsg))
	}
	// Start : ZSP_MAIL_DB_POOL_UPT

	return nil
}

func (resp *MailDBPool) addDB(toNameList []string, toEmailList []string, fromNameList []string, fromEmailList []string, subjectList []string, contentsParamList []string, contentsFuncList []string) error {

	pLang, _ := beego.AppConfig.String("lang")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		logs.Error(fmt.Sprintf("[MailDBPool][callbackSend][Error] %v", err))
		return err
	}
	// End : Oracle DB Connection

	// Start : ZSP_MAIL_DB_POOL_REG
	fmt.Printf(fmt.Sprintf("CALL ZSP_MAIL_DB_POOL_REG('%v', %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang,
		len(toNameList),
		strings.Join(toNameList, ","),
		strings.Join(toEmailList, ","),
		strings.Join(fromNameList, ","),
		strings.Join(fromEmailList, ","),
		strings.Join(subjectList, ","),
		strings.Join(contentsParamList, ","),
		strings.Join(contentsFuncList, ",")))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_MAIL_DB_POOL_REG('%v', %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', :1)",
		pLang,
		len(toNameList),
		strings.Join(toNameList, ","),
		strings.Join(toEmailList, ","),
		strings.Join(fromNameList, ","),
		strings.Join(fromEmailList, ","),
		strings.Join(subjectList, ","),
		strings.Join(contentsParamList, ","),
		strings.Join(contentsFuncList, ",")),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
	)
	defer stmtProcCall.Close()
	if err != nil {
		logs.Error(fmt.Sprintf("[MailDBPool][added][Error] %v", err))
		return err
	}
	procRset := &ora.Rset{}
	_, err = stmtProcCall.Exe(procRset)

	if err != nil {
		logs.Error(fmt.Sprintf("[MailDBPool][added][Error] %v", err))
		return err
	}

	var (
		rtnCd  int64
		rtnMsg string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
		}

		if err := procRset.Err(); err != nil {
			logs.Error(fmt.Sprintf("[MailDBPool][added][Error] %v", err))
			return err
		}

		fmt.Printf(fmt.Sprintf(" ===> rtnCd:%v, rtnMsg:%v", rtnCd, rtnMsg))
	}
	// Start : ZSP_MAIL_DB_POOL_REG

	return nil
}

func (resp *MailDBPool) CompleteProcess(index int, cbData string) {
	fmt.Printf(fmt.Sprintf("[MailDBPool][CompleteProcess] Complete %v", cbData))
}

var MailDBPoolMng MailDBPool

// SmtpSendCounter ---------------------------------------------------------->
type MailSendCount struct {
	SmtpServer string

	AssignDt  time.Time
	CheckHour int32

	Cnt    int32
	MaxCnt int32
}

type MailSendCounter struct {
	MapList map[string]MailSendCount
}

func (resp *MailSendCounter) Load() error {

	resp.MapList = make(map[string]MailSendCount)

	err := resp.loadDB()
	if err != nil {
		return err
	}

	return nil
}

func (resp *MailSendCounter) IsValidCount(smtpServer string) (error, bool) {

	findResult, ok := resp.MapList[smtpServer]
	if ok == false {
		return fmt.Errorf("%v", "Invalid smtpServer"), false
	}

	return nil, findResult.Cnt < findResult.MaxCnt
}

func (resp *MailSendCounter) CntCheck(smtpServer string) error {

	curTime := time.Now()

	//fmt.Printf(fmt.Sprintf("[MailDBPool][MailSendCounter] CheckServer: %v", smtpServer))

	findResult, ok := resp.MapList[smtpServer]
	if ok == false {
		return fmt.Errorf("%v", "Invalid smtpServer")
	}

	then := findResult.AssignDt.Add(time.Duration(findResult.CheckHour) * time.Hour)

	fmt.Printf(fmt.Sprintf("[MailDBPool][MailSendCounter] CheckServer: %v, Cnt: %v, MaxCnt: %v, AssignDt: %v, CheckDt(%v hour): %v",
		smtpServer, findResult.Cnt, findResult.MaxCnt, findResult.AssignDt.Format("06/01/02 15:04"), findResult.CheckHour, then.Format("06/01/02 15:04")))

	if curTime.After(then) {
		err := resp.ResetCnt(smtpServer)
		if err != nil {
			return fmt.Errorf("ResetCnt() Error: %v", err)
			//logs.Error(fmt.Sprintf("[MailDBPool][MailSendCounter][Error] ResetCnt() =>  smtpServer: %v, Error: %v", smtpServer, err))
		}
	}

	return nil
}

func (resp *MailSendCounter) SetCnt(smtpServer string, setCnt int32) error {

	findResult, ok := resp.MapList[smtpServer]
	if ok == false {
		return fmt.Errorf("%v", "Invalid smtpServer")
	}

	err, resultCnt := resp.cntUptDB(smtpServer, setCnt)
	if err != nil {
		return fmt.Errorf("cntUptDB() Error:%v", err)
	}

	findResult.Cnt = resultCnt

	resp.MapList[smtpServer] = findResult

	fmt.Printf(fmt.Sprintf("[MailDBPool][MailSendCounter] SetCnt: %v", findResult.Cnt))

	return nil
}

func (resp *MailSendCounter) ResetCnt(smtpServer string) error {

	findResult, ok := resp.MapList[smtpServer]
	if ok == false {
		return fmt.Errorf("%v", "Invalid smtpServer")
	}

	err, resultCnt, resultAssignDt := resp.resetDB(smtpServer)
	if err != nil {
		return fmt.Errorf("resetDB() Error:%v", err)
	}

	findResult.Cnt = resultCnt
	findResult.AssignDt = resultAssignDt

	resp.MapList[smtpServer] = findResult

	fmt.Printf(fmt.Sprintf("[MailDBPool][MailSendCounter] ResetCnt: %v, ResetAssignDt: %v", findResult.Cnt, findResult.AssignDt))

	return nil
}

func (resp *MailSendCounter) loadDB() error {

	pLang, _ := beego.AppConfig.String("lang")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		return err
	}
	// End : Oracle DB Connection

	// Start : ZSP_MAIL_DB_POOL_LIST
	fmt.Printf(fmt.Sprintf("CALL ZSP_MAIL_SEND_CNT_LIST('%v', :1)", pLang))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_MAIL_SEND_CNT_LIST('%v', :1)", pLang),
		ora.I64, /* TOT_CNT */
		ora.S,   /* SMTP_SERVER */
		ora.S,   /* ASSIGN_DT */
		ora.I32, /* CHECK_HOUR */
		ora.I32, /* CNT */
		ora.I32, /* MAX_CNT */
	)
	defer stmtProcCall.Close()
	if err != nil {
		return err
	}
	procRset := &ora.Rset{}
	_, err = stmtProcCall.Exe(procRset)

	if err != nil {
		return err
	}

	mailSendCount := make([]MailSendCount, 0)

	var (
		//totCnt     int64
		smtpServer string
		assignDt   string
		checkHour  int32
		cnt        int32
		maxCnt     int32
	)

	if procRset.IsOpen() {
		for procRset.Next() {

			_ = procRset.Row[0].(int64)
			smtpServer = procRset.Row[1].(string)
			assignDt = procRset.Row[2].(string)
			checkHour = procRset.Row[3].(int32)
			cnt = procRset.Row[4].(int32)
			maxCnt = procRset.Row[5].(int32)

			conv_assignDt, _ := time.ParseInLocation("20060102150405", assignDt, time.Local)

			mailSendCount = append(mailSendCount, MailSendCount{
				SmtpServer: smtpServer,
				AssignDt:   conv_assignDt,
				CheckHour:  checkHour,
				Cnt:        cnt,
				MaxCnt:     maxCnt,
			})
		}
		if err := procRset.Err(); err != nil {
			return err
		}

		for _, val := range mailSendCount {
			resp.MapList[val.SmtpServer] = val
		}
	}
	// End : ZSP_MAIL_DB_POOL_LIST

	return err
}

func (resp *MailSendCounter) cntUptDB(smtpServer string, uptCount int32) (error, int32) {

	pLang, _ := beego.AppConfig.String("lang")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		return err, 0
	}
	// End : Oracle DB Connection

	// Start : ZSP_MAIL_DB_POOL_UPT
	fmt.Printf(fmt.Sprintf("CALL ZSP_MAIL_SEND_CNT_UPT('%v', '%v', %v, :1)", pLang, smtpServer, uptCount))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_MAIL_SEND_CNT_UPT('%v', '%v', %v, :1)", pLang, smtpServer, uptCount),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
		ora.I32, /* RTN_COUNT */
	)
	defer stmtProcCall.Close()
	if err != nil {
		return err, 0
	}
	procRset := &ora.Rset{}
	_, err = stmtProcCall.Exe(procRset)

	if err != nil {
		return err, 0
	}

	var (
		rtnCd  int64
		rtnMsg string
		rtnCnt int32
	)

	rtnCnt = 0
	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
			rtnCnt = procRset.Row[2].(int32)
		}

		if err := procRset.Err(); err != nil {
			return err, 0
		}

		fmt.Printf(fmt.Sprintf(" ===> rtnCd:%v, rtnMsg:%v, rtnCnt:%v", rtnCd, rtnMsg, rtnCnt))
	}
	// Start : ZSP_MAIL_DB_POOL_UPT

	return err, rtnCnt
}

func (resp *MailSendCounter) resetDB(smtpServer string) (error, int32, time.Time) {

	pLang, _ := beego.AppConfig.String("lang")

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		return err, 0, time.Time{}
	}
	// End : Oracle DB Connection

	// Start : ZSP_MAIL_SEND_CNT_RESET
	fmt.Printf(fmt.Sprintf("CALL ZSP_MAIL_SEND_CNT_RESET('%v', '%v', :1)", pLang, smtpServer))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL ZSP_MAIL_SEND_CNT_RESET('%v', '%v', :1)", pLang, smtpServer),
		ora.I64, /* RTN_CD */
		ora.S,   /* RTN_MSG */
		ora.I32, /* RTN_COUNT */
		ora.S,   /* RTN_ASSIGN_DT */
	)
	defer stmtProcCall.Close()
	if err != nil {
		return err, 0, time.Time{}
	}
	procRset := &ora.Rset{}
	_, err = stmtProcCall.Exe(procRset)

	if err != nil {
		return err, 0, time.Time{}
	}

	var (
		rtnCd       int64
		rtnMsg      string
		rtnCnt      int32
		rtnAssignDt string

		convAssignDt time.Time
	)

	convAssignDt = time.Time{}
	if procRset.IsOpen() {
		for procRset.Next() {
			rtnCd = procRset.Row[0].(int64)
			rtnMsg = procRset.Row[1].(string)
			rtnCnt = procRset.Row[2].(int32)
			rtnAssignDt = procRset.Row[3].(string)

			convAssignDt, _ = time.ParseInLocation("20060102150405", rtnAssignDt, time.Local)
		}

		if err := procRset.Err(); err != nil {
			return err, 0, time.Time{}
		}

		fmt.Printf(fmt.Sprintf(" ===> rtnCd:%v, rtnMsg:%v, rtnCnt:%v, rtnAssignDt:%v", rtnCd, rtnMsg, rtnCnt, rtnAssignDt))
	}
	// Start : ZSP_MAIL_SEND_CNT_RESET

	return err, rtnCnt, convAssignDt
}

// SmtpSendCounter <----------------------------------------------------------

// --------------------------------------------------------------------------------
// --------------------------------------------------------------------------------
// --------------------------------------------------------------------------------
//var CheckIndex int64
//var MailDBPoolFuncList MailDBPoolFunc

func AddArrayTest(cnt int) {

	var toMailArr []MailDBPoolData

	t := time.Now()

	for i := 0; i < cnt; i++ {

		//CheckIndex++

		checkIndex := GetIndexTest()

		toNm := fmt.Sprintf("이동관")
		toEmail := fmt.Sprintf("dongkale@naver.com")
		// fromNm := "큐레잇"
		// fromEmail := "no-reply@ziggam.com"
		fromNm := "큐레잇"                   // naver smtp 서버를 이용하기 위해
		fromEmail := "dongkale@naver.com" // naver smtp 서버를 이용하기 위해
		subject := fmt.Sprintf("이동관:%v [%v]", checkIndex, t.Format("2006-01-02 15:04:05"))
		contentsFunc := "HtmlContents"
		contentsParam := "arg1=ONE;arg2=TWO;arg3=THREE"

		toMailArr = append(toMailArr, MailDBPoolData{
			ToNm:          toNm,
			ToEmail:       toEmail,
			FromNm:        fromNm,
			FromEmail:     fromEmail,
			Subject:       subject,
			ContentsFunc:  contentsFunc,
			ContentsParam: contentsParam,
		})

		checkIndex++
		SetIndexTest(checkIndex)
	}

	err := MailDBPoolMng.AddArray(toMailArr)
	if err != nil {
		logs.Error(fmt.Sprintf("[Error] %v", err))
	}
}

func GetIndexTest() int64 {

	dat, _ := ioutil.ReadFile("conf/maildbpool.conf")
	ret, _ := strconv.ParseInt(string(dat), 10, 64)

	return ret
}

func SetIndexTest(dat int64) {

	err := ioutil.WriteFile("conf/maildbpool.conf", []byte(fmt.Sprintf("%d", dat)), 0644)
	if err != nil {
		logs.Error(fmt.Sprintf("[Error] %v", err))
	}
}
