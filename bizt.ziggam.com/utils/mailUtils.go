package utils

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"time"
	
	"github.com/jordan-wright/email"
	"github.com/knadh/smtppool"
	"gopkg.in/gomail.v2"
)

type MailStruct struct {
	To     []string
	ToName []string

	From     string
	FromName string

	Subject  string
	Contents string
}

//--------------
//https://kinsta.com/blog/gmail-smtp-server/#gmail-smtp-server-faqs
//https://blog.naver.com/gnltns/222271831827
var smtpEmail = "support@ziggam.com" //beego.AppConfig.String("smtpEmail")
var emailPwd = "wlrrka0223^^!"       //beego.AppConfig.String("emailPwd")

var smtpServer = "smtp.gmail.com"
var smtpServerPort = 587

//--------------

//--------------
// var smtpServer = "smtp.mail.yahoo.com"
// var smtpServerPort = 587

// var smtpEmail = "dongkale@yahoo.com"
// var emailPwd = "crraynkxvocakhmd"

//crraynkxvocakhmd
//--------------

//--------------
// var smtpServer = "smtp.naver.com"
// var smtpServerPort = 465

// var smtpEmail = "dongkale@naver.com"
// var emailPwd = "doublespy71!"

//--------------

func GetReturnEmail() string {
	return smtpEmail
}

func GetReturnEmailPwd() string {
	return emailPwd
}

func GetSmtpServer() string {
	return smtpServer
}

func GetSmtpServerPort() int {
	return smtpServerPort
}

// SendMail ...
func SendMail(to string, toName string, from string, fromName string, subject string, htmlContents string) error {

	// https://stackoverrun.com/ko/q/7480544
	m := gomail.NewMessage()
	// m.SetHeader("From", from)
	// m.SetHeader("To", to)
	// m.SetHeader("Subject", subject)
	m.SetAddressHeader("From", from, fromName)
	m.SetAddressHeader("To", to, toName)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlContents)

	d := gomail.NewDialer(smtpServer, smtpServerPort, smtpEmail, emailPwd)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		//panic(err)
		fmt.Printf(fmt.Sprintf("[SendMail][Error] %s", err))
		return fmt.Errorf("Could not dial email to %q(%v): %v", to, toName, err)
	}

	return nil
}

// SendMailEx ...
func SendMailEx(to string, toName string, from string, fromName string, subject string, htmlContents string) error {

	d := gomail.NewDialer(smtpServer, smtpServerPort, smtpEmail, emailPwd)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	s, err := d.Dial()
	if err != nil {
		//return err
		fmt.Printf(fmt.Sprintf("[SendMailEx][1][Error] %s", err))
		return fmt.Errorf("Could not dial email to %q(%v): %v", to, toName, err)
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", from, fromName)
	m.SetAddressHeader("To", to, toName)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlContents)

	if err := gomail.Send(s, m); err != nil {
		fmt.Printf(fmt.Sprintf("[SendMailEx][2][Error] %s", err))
		return fmt.Errorf("Could not send email to %q(%v): %v", to, toName, err)
	}

	m.Reset()

	return nil
}

// SendMailExEx ...
func SendMailExEx(to string, toName string, from string, fromName string, subject string, htmlContents string) error {

	auth := smtp.PlainAuth("", smtpEmail, emailPwd, smtpServer)

	// toArray := []string{to}
	// headerBlank := "\r\n"
	// msg := []byte(subject + headerBlank + htmlContents)
	// err := smtp.SendMail("smtp.gmail.com:smtpServerPort", auth, from, toArray, msg)
	// if err != nil {
	// 	panic(err)
	// }

	message := fmt.Sprintf("To: '%v' <%v>\r\nFrom: '%v' <%v>\r\nSubject:%v\r\n%v", toName, to, fromName, from, subject, htmlContents)

	err := smtp.SendMail(smtpServer+":587", auth, from, []string{to}, []byte(message))
	if err != nil {
		fmt.Printf(fmt.Sprintf("[SendMail][Error] %s", err))
		return fmt.Errorf("Could not send email to %q(%v): %v", to, toName, err)
		//panic(err)
	}

	return nil
}

// SendMailer  ---------------------------------------------------------->

// SendMailer ...
type StMailSender struct {
	gomailSendCloser gomail.SendCloser
}

// Connect ...
func (resp *StMailSender) Connect(mailServer string, mailServerPort int, returnEmail string, returnEmailPwd string) error {
	var err error

	d := gomail.NewDialer(mailServer, mailServerPort, returnEmail, returnEmailPwd)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	resp.gomailSendCloser, err = d.Dial()
	if err != nil {
		return fmt.Errorf("Could not dial email to %v", err)
	}

	return nil
}

// Send ...
func (resp *StMailSender) Send(to string, toName string, from string, fromName string, subject string, htmlContents string) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", from, fromName)
	m.SetAddressHeader("To", to, toName)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlContents)

	if err := gomail.Send(resp.gomailSendCloser, m); err != nil {
		return fmt.Errorf("Could not send email to %v : %q(%v)", err, to, toName)
	}

	m.Reset()

	fmt.Printf(fmt.Sprintf("[MailSender][Send] %s:%s", toName, to))

	return nil
}

var MailSender StMailSender

// SendMailPool ---------------------------------------------------------->

// SendMailPool ---------------------------------------------------------->
// go get github.com/jordan-wright/email

type SendMailPoolCallback func(error, int, int, string, string)
type SendMailPoolCompleteCallback func(int, string)

// type SendMailPoolDestFmt struct {
// 	Dest     string
// 	DestName string
// }

// func (resp *SendMailPoolDestFmt) Get() string {
// 	return fmt.Sprintf("%v <%v>", resp.DestName, resp.Dest)
// }

func SendMailPoolDestFmt(destName string, dest string) string {
	return fmt.Sprintf("%v <%v>", destName, dest)
}

type SendMailPoolData struct {
	To   []string
	From string

	Subject      string
	HtmlContents string

	Bcc []string
	Cc  []string

	//Result string
}

type SendMailPoolPushData struct {
	ToSend string

	ToSendData SendMailPoolData
	Cb         SendMailPoolCallback
	CbData     string
}

type SendMailPool struct {
	//PushChannel chan email.Email
	//ResultChannel chan SendMailPoolData

	PushChannel chan SendMailPoolPushData
	SendIndex   int

	SendTryCount int

	CompleteIndex  int
	CompleteCb     SendMailPoolCompleteCallback
	CompleteCbData string
}

// Init ...
func (resp *SendMailPool) ClearSendIndex() {
	resp.SendIndex = 0
}

func (resp *SendMailPool) GetSendIndex() int {
	return resp.SendIndex
}

func (resp *SendMailPool) ClearComplete() {
	resp.CompleteIndex = 0
}

func (resp *SendMailPool) IsComplete() bool {
	return resp.CompleteIndex == 0
}

// func (resp *SendMailPool) SetCheckIndex(CheckIndex int, CheckIndexCb SendMailPoolCheckIndexCallback, CheckIndexCbData string) error {
// 	if resp.CheckIndex > 0 {
// 		return fmt.Errorf("Already Set:%d", resp.CheckIndex)
// 	}

// 	resp.CheckIndex = CheckIndex
// 	resp.CheckIndexCb = CheckIndexCb
// 	resp.CheckIndexCbData = CheckIndexCbData

// 	return nil
// }

func (resp *SendMailPool) SetCompleteCb(setCount int, completeCb SendMailPoolCompleteCallback, completeCbData string) error {
	if resp.IsComplete() == false {
		return fmt.Errorf("Not Complte:%d", resp.CompleteIndex)
	}

	resp.CompleteIndex = resp.GetSendIndex() + setCount
	resp.CompleteCb = completeCb
	resp.CompleteCbData = completeCbData

	fmt.Printf(fmt.Sprintf("=========================================== %v ===========================================", resp.CompleteIndex))

	return nil
}

// Init ...
func (resp *SendMailPool) Init(poolCount int, channelCount int, sendTryCount int) error {

	resp.ClearSendIndex()

	resp.SendTryCount = sendTryCount

	//resp.PushChannel = make(chan email.Email, channelCount)
	//resp.ResultChannel = make(chan SendMailPoolData, channelCount)
	resp.PushChannel = make(chan SendMailPoolPushData, channelCount)

	pool, _ := email.NewPool(
		fmt.Sprintf("%v:%d", smtpServer, smtpServerPort),
		poolCount,
		smtp.PlainAuth("", smtpEmail, emailPwd, smtpServer),
	)

	for i := 0; i < poolCount; i++ {
		go func(n int) {
			for m := range resp.PushChannel {

				// var convTo []string
				// for _, val := range m.ToSendData.To {
				// 	convTo = append(convTo, val.Get())
				// }

				d := email.Email{
					// From: fmt.Sprintf("%v <%v>", m.ToSendData.FromName, m.ToSendData.From),
					// To:   []string{fmt.Sprintf("%v <%v>", m.ToSendData.ToName, m.ToSendData.To)},
					//To:   []string{m.ToSendData.To.Get()},

					From: m.ToSendData.From,
					To:   m.ToSendData.To,

					Bcc: m.ToSendData.Bcc,
					Cc:  m.ToSendData.Cc,

					Subject: m.ToSendData.Subject,
					HTML:    []byte(m.ToSendData.HtmlContents),
				}

				var err error
				for i := 0; i < resp.SendTryCount; i++ {
					err = pool.Send(&d, 30*time.Second)
					if err != nil {
						// 접속 종료시: write tcp 192.168.1.207:49810->108.177.97.109:smtpServerPort: wsasend: An established connection was aborted by the software in your host machine
						// --> 이때만 재시도 ?!
						fmt.Printf(fmt.Sprintf("[SendMailPool][%v][%v][Send][Error] %s: %s", n, resp.SendIndex, m.ToSend, err))

						// if m.Cb != nil {
						// 	m.Cb(err, n, resp.SendIndex, m.ToSend, m.CbData)
						// }
					} else {
						fmt.Printf(fmt.Sprintf("[SendMailPool][%v][%v][Send] %s", n, resp.SendIndex, m.ToSend))

						// if m.Cb != nil {
						// 	m.Cb(err, n, resp.SendIndex, m.ToSend, m.CbData)
						// }

						break
					}
				}

				if m.Cb != nil {
					m.Cb(err, n, resp.SendIndex, m.ToSend, m.CbData)
				}

				resp.SendIndex++

				if resp.CompleteIndex > 0 {
					if resp.CompleteIndex == resp.SendIndex {
						if resp.CompleteCb != nil {
							resp.CompleteCb(resp.SendIndex, resp.CompleteCbData)
						}
						resp.CompleteIndex = 0
					}
				}
			}
		}(i)
	}

	return nil
}

// // Push ... 사용금지
// func (resp *SendMailPool) Push(toSend string, to string, toName string, from string, fromName string, subject string, htmlContents string, bcc []string, cc []string, cb SendMailPoolCallback, cbData string) error {

// 	// e := email.Email{
// 	// 	From: fmt.Sprintf("%v <%v>", fromName, from),
// 	// 	To:   []string{fmt.Sprintf("%v <%v>", toName, to)},

// 	// 	Bcc: bcc,
// 	// 	Cc:  cc,

// 	// 	Subject: subject,
// 	// 	//Text:    []byte("This is a test e-mail"),
// 	// 	HTML: []byte(htmlContents),
// 	// }

// 	d := SendMailPoolPushData{
// 		ToSend: toSend,
// 		ToSendData: SendMailPoolData{
// 			To:   []string{SendMailPoolDestFmt(toName, to)},
// 			From: SendMailPoolDestFmt(fromName, from),

// 			//ToName:       toName,
// 			//FromName:     fromName,
// 			Bcc:          bcc,
// 			Cc:           cc,
// 			Subject:      subject,
// 			HtmlContents: htmlContents,
// 		},
// 		Cb:     cb,
// 		CbData: cbData,
// 	}

// 	resp.PushChannel <- d

// 	return nil
// }

// PushArray ...
func (resp *SendMailPool) PushArray(arrayData []SendMailPoolPushData, delaySec time.Duration) error {

	// for _, val := range arrayData {
	// 	e := email.Email{
	// 		From: fmt.Sprintf("%v <%v>", val.FromName, val.From),
	// 		To:   []string{fmt.Sprintf("%v <%v>", val.ToName, val.To)},

	// 		Bcc: val.Bcc,
	// 		Cc:  val.Cc,

	// 		Subject: val.Subject,
	// 		//Text:    []byte("This is a test e-mail"),
	// 		HTML: []byte(val.HtmlContents),
	// 	}

	// 	resp.PushChannel <- e

	// 	time.Sleep(time.Second * delaySec)
	// }

	for _, val := range arrayData {
		resp.PushChannel <- val
		time.Sleep(time.Second * delaySec)
	}

	return nil
}

// PushArrayEx ...
func (resp *SendMailPool) PushArrayEx(arrayData []SendMailPoolPushData, delaySec time.Duration, completeCb SendMailPoolCompleteCallback, completeCbData string) error {
	var count = len(arrayData)

	err := resp.SetCompleteCb(count, completeCb, completeCbData)
	if err != nil {
		return err
	}

	for _, val := range arrayData {
		resp.PushChannel <- val

		if count > 1 {
			time.Sleep(time.Second * delaySec)
		}
	}

	return nil
}

var SendMailPoolMng SendMailPool

// SendMailPool <----------------------------------------------------------

// SendMailDaemon ---------------------------------------------------------->

type SendMailDaemonData struct {
	To           string
	ToName       string
	From         string
	FromName     string
	Subject      string
	HtmlContents string

	// Bcc []string
	// Cc  []string
	// Result string
}

// https://pythonq.com/so/go/1592087
type SendMailDaemonCallback func(error, int, string, string)

type SendMailDaemonPushData struct {
	ToSend string
	//Message gomail.Message
	ToSendData SendMailDaemonData
	Cb         SendMailDaemonCallback
	CbData     string
}

// func (resp *SendMailDaemonPushData) Set(toSend string, cb SendMailDaemonCallback, cbData string, to string, toName string, from string, fromName string, subject string, htmlContents string) {
// 	resp.ToSend = toSend
// 	resp.ToSendData = SendMailDaemonData{
// 		To:           to,
// 		ToName:       toName,
// 		From:         from,
// 		FromName:     fromName,
// 		Subject:      subject,
// 		HtmlContents: htmlContents,
// 	}
// 	resp.Cb = cb
// 	resp.CbData = cbData
// }

type SendMailDaemon struct {
	//PushChannel chan *gomail.Message

	PushChannel chan SendMailDaemonPushData

	SendIndex int
}

// Init ...
func (resp *SendMailDaemon) ClearSendIndex() {
	resp.SendIndex = 0
}

func (resp *SendMailDaemon) GetSendIndex() int {
	return resp.SendIndex
}

// Init ...
func (resp *SendMailDaemon) Init() error {

	resp.ClearSendIndex()

	//resp.PushChannel = make(chan *gomail.Message)
	resp.PushChannel = make(chan SendMailDaemonPushData)

	go func() {
		d := gomail.NewDialer(smtpServer, smtpServerPort, smtpEmail, emailPwd)
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

		var s gomail.SendCloser
		var err error
		open := false
		for {
			select {
			case m, ok := <-resp.PushChannel:
				if !ok {
					return
				}
				if !open {
					if s, err = d.Dial(); err != nil {
						if m.Cb != nil {
							m.Cb(err, resp.SendIndex, m.ToSend, m.CbData)
						}
						fmt.Printf(fmt.Sprintf("[SendMailDaemon][Error:1] %s", err))
						break
					}

					open = true
				}

				d := gomail.NewMessage()
				d.SetAddressHeader("From", m.ToSendData.From, m.ToSendData.FromName)
				d.SetAddressHeader("To", m.ToSendData.To, m.ToSendData.ToName)
				d.SetHeader("Subject", m.ToSendData.Subject)
				d.SetBody("text/html", m.ToSendData.HtmlContents)

				err := gomail.Send(s, d)
				if err != nil {
					fmt.Printf(fmt.Sprintf("[SendMailDaemon][%v][Send][Error:2] %s: %s", resp.SendIndex, m.ToSend, err))
				} else {
					fmt.Printf(fmt.Sprintf("[SendMailDaemon][%v][Send] %s", resp.SendIndex, m.ToSend))
				}

				if m.Cb != nil {
					m.Cb(err, resp.SendIndex, m.ToSend, m.CbData)
				}

				resp.SendIndex++

			// Close the connection to the SMTP server if no email was sent in
			// the last 30 seconds.
			case <-time.After(30 * time.Second):
				if open {
					if err := s.Close(); err != nil {
						fmt.Printf(fmt.Sprintf("[SendMailDaemon][Error:3] %s", err))
						break
					}
					open = false
				}
			}
		}
	}()

	return nil
}

// Push ...
func (resp *SendMailDaemon) Push(toSend string, to string, toName string, from string, fromName string, subject string, htmlContents string, cb SendMailDaemonCallback, cbData string) error {

	// m := gomail.NewMessage()
	// m.SetAddressHeader("From", from, fromName)
	// m.SetAddressHeader("To", to, toName)
	// m.SetHeader("Subject", subject)
	// m.SetBody("text/html", htmlContents)

	//resp.PushChannel <- m

	d := SendMailDaemonPushData{
		ToSend: toSend,
		ToSendData: SendMailDaemonData{
			To:           to,
			ToName:       toName,
			From:         from,
			FromName:     fromName,
			Subject:      subject,
			HtmlContents: htmlContents,
		},
		Cb:     cb,
		CbData: cbData,
	}

	resp.PushChannel <- d

	return nil
}

// PushArray ...
func (resp *SendMailDaemon) PushArray(arrayData []SendMailDaemonPushData, delaySec time.Duration) error {

	// for _, val := range arrayData {
	// 	m := gomail.NewMessage()
	// 	m.SetAddressHeader("From", val.From, val.FromName)
	// 	m.SetAddressHeader("To", val.To, val.ToName)
	// 	m.SetHeader("Subject", val.Subject)
	// 	m.SetBody("text/html", val.HtmlContents)

	// 	resp.PushChannel <- m

	// 	time.Sleep(time.Second * delaySec)
	// }

	// for _, val := range arrayData {
	// 	m := gomail.NewMessage()
	// 	m.SetAddressHeader("From", val.From, val.FromName)
	// 	m.SetAddressHeader("To", val.To, val.ToName)
	// 	m.SetHeader("Subject", val.Subject)
	// 	m.SetBody("text/html", val.HtmlContents)

	// 	d := SendMailDaemonPushData{
	// 		To:      fmt.Sprintf("%s:%s", val.ToName, val.To),
	// 		Message: *m,
	// 	}

	// 	resp.PushChannel <- d

	// 	time.Sleep(time.Second * delaySec)
	// }

	var count = len(arrayData)

	for _, val := range arrayData {
		resp.PushChannel <- val

		if count > 1 {
			time.Sleep(time.Second * delaySec)
		}
	}

	return nil
}

var SendMailDaemonMng SendMailDaemon

// SendMailDaemon <----------------------------------------------------------

// ------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------

func SendMailDaemon_Test(to string, toName string, from string, fromName string, subject string, htmlContents string) error {

	ch := make(chan *gomail.Message)

	go func() {
		d := gomail.NewDialer(smtpServer, smtpServerPort, smtpEmail, emailPwd)
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

		var s gomail.SendCloser
		var err error
		open := false
		for {
			select {
			case m, ok := <-ch:
				if !ok {
					return
				}
				if !open {
					if s, err = d.Dial(); err != nil {
						//panic(err)
						fmt.Printf(fmt.Sprintf("[SendMailDaemon][1][Error] %s", err))
						break
					}
					// s, err := d.Dial()
					// if err != nil {
					// 	//return err
					// 	fmt.Printf(fmt.Sprintf("[SendMailEx][1][Error] %s", err))
					// 	break //return fmt.Errorf("Could not dial email to %q(%v): %v", to, toName, err)
					// }

					open = true
				}
				if err := gomail.Send(s, m); err != nil {
					//log.Print(err)
					fmt.Printf(fmt.Sprintf("[SendMailDaemon][2][Error] %s", err))
					break
				}

				//fmt.Printf(fmt.Sprintf("[SendMailDaemon][Send] %s", m))
			// Close the connection to the SMTP server if no email was sent in
			// the last 30 seconds.
			case <-time.After(30 * time.Second):
				if open {
					if err := s.Close(); err != nil {
						//panic(err)
						fmt.Printf(fmt.Sprintf("[SendMailDaemon][3][Error] %s", err))
						break
					}
					open = false
				}
			}
		}
	}()

	// Use the channel in your program to send emails.
	m := gomail.NewMessage()
	m.SetAddressHeader("From", from, fromName)
	m.SetAddressHeader("To", to, toName)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlContents)

	ch <- m

	// m = gomail.NewMessage()
	// m.SetAddressHeader("From", from, fromName)
	// m.SetAddressHeader("To", to, toName)
	// m.SetHeader("Subject", subject)
	// m.SetBody("text/html", htmlContents)

	// ch <- m

	// m = gomail.NewMessage()
	// m.SetAddressHeader("From", from, fromName)
	// m.SetAddressHeader("To", to, toName)
	// m.SetHeader("Subject", subject)
	// m.SetBody("text/html", htmlContents)

	// ch <- m

	// for i := 0; i < 10; i++ {
	// 	m := gomail.NewMessage()
	// 	m.SetAddressHeader("From", from, fromName)
	// 	m.SetAddressHeader("To", to, toName)
	// 	m.SetHeader("Subject", subject)
	// 	m.SetBody("text/html", htmlContents)

	// 	ch <- m
	// }

	// Close the channel to stop the mail daemon.
	close(ch)

	return nil
}

// go get github.com/knadh/smtppool
func SendMailPool_Test(to string, toName string, from string, fromName string, subject string, htmlContents string) error {

	//d := gomail.NewDialer(smtpServer, smtpServerPort, smtpEmail, emailPwd)

	pool, err := smtppool.New(smtppool.Opt{
		Host:            smtpServer,
		Port:            smtpServerPort,
		MaxConns:        10,
		IdleTimeout:     time.Second * 10,
		PoolWaitTimeout: time.Second * 3,
		//Auth:            smtp.PlainAuth("", smtpEmail, emailPwd, smtpServer),
	})
	if err != nil {
		log.Fatalf("error creating pool: %v", err)
	}

	e := smtppool.Email{
		From: from,
		To:   []string{to},

		// Optional.
		// Bcc: []string{"doebcc@example.com"},
		// Cc:  []string{"doecc@example.com"},

		Subject: subject,
		//Text:    []byte("This is a test e-mail"),
		HTML: []byte(htmlContents),
	}

	// Add attachments.
	// if _, err := e.AttachFile("test.txt"); err != nil {
	// 	log.Fatalf("error attaching file: %v", err)
	// }

	if err := pool.Send(e); err != nil {
		log.Fatalf("error sending e-mail: %v", err)
	}

	return nil
}

func SendMailPoolEx(to string, toName string, from string, fromName string, subject string, htmlContents string) error {
	//var ch <-chan *email.Email
	ch := make(chan *email.Email)

	p, _ := email.NewPool(
		fmt.Sprintf("%v:%d", smtpServer, smtpServerPort),
		4,
		smtp.PlainAuth("", smtpEmail, emailPwd, smtpServer),
	)

	for i := 0; i < 4; i++ {
		go func() {
			for e := range ch {
				err := p.Send(e, 10*time.Second)
				if err != nil {
					fmt.Printf(fmt.Sprintf("[SendMailPoolEx][3][Error] %s", err))
				}
			}
		}()
	}

	// e := &email.Email{
	// 	From: fmt.Sprintf("%v <%v>", fromName, from),
	// 	To:   []string{fmt.Sprintf("%v <%v>", toName, to)},

	// 	// Optional.
	// 	// Bcc: []string{"doebcc@example.com"},
	// 	// Cc:  []string{"doecc@example.com"},

	// 	Subject: subject,
	// 	//Text:    []byte("This is a test e-mail"),
	// 	HTML: []byte(htmlContents),
	// }

	// ch <- e

	for i := 0; i < 100; i++ {
		e := &email.Email{
			From: fmt.Sprintf("%v <%v>", fromName, from),
			To:   []string{fmt.Sprintf("%v <%v>", toName, to)},

			Subject: subject,
			Text:    []byte("This is a test e-mail"),
			HTML:    []byte(htmlContents),
		}

		ch <- e
	}

	return nil
}

var chSendMailPool = make(chan email.Email)

func SendMailPoolStart(count int) error {
	pool, _ := email.NewPool(
		fmt.Sprintf("%v:%d", smtpServer, smtpServerPort),
		count,
		smtp.PlainAuth("", smtpEmail, emailPwd, smtpServer),
	)

	for i := 0; i < count; i++ {
		go func() {
			for e := range chSendMailPool {
				err := pool.Send(&e, 30*time.Second)
				if err != nil {
					fmt.Printf(fmt.Sprintf("[SendMailPool][Send][Error] %s: %s", e.To, err))
				} else {
					fmt.Printf(fmt.Sprintf("[SendMailPool][Send] %s", e.To))
				}
			}

			fmt.Printf(fmt.Sprintf("[SendMailPool]------------"))
		}()
	}

	return nil
}

func SendMailPoolPush(to string, toName string, from string, fromName string, subject string, htmlContents string) error {

	for i := 0; i < 10; i++ {
		e := email.Email{
			From: fmt.Sprintf("%v <%v>", fromName, from),
			To:   []string{fmt.Sprintf("%v <%v>", toName, to)},

			Subject: subject,
			Text:    []byte("This is a test e-mail"),
			HTML:    []byte(htmlContents),
		}

		chSendMailPool <- e
		time.Sleep(time.Second * 2)
	}

	return nil
}

// ------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------
