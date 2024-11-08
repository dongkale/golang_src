package controllers

import (
	"fmt"
	"log"
	"strconv"

	"emst.ziggam.com/models"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	"golang.org/x/net/context"
	"google.golang.org/api/option"

	ora "gopkg.in/rana/ora.v4"
)

type NoticePushController struct {
	BaseController
}

func (c *NoticePushController) Post() {

	// start : log
	// log := logs.NewLogger()
	// log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/common/login")
	}

	pLang, _ := beego.AppConfig.String("lang")

	pSn := c.GetString("sn") // 일련번호
	//pGbnCd := c.GetString("gbn_cd") // 게시구분코드 (02면 이벤트성 광고)
	//pMemCd := c.GetString("mem_cd") //

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Admin Event Process
	// 게시물이 있는지 검사 하자.
	// 게시물 제목이 일치 하는지도 검사해주자.
	// Start : Notice List
	logs.Debug(fmt.Sprintf("CALL SP_EMS_NOTICE_DTL_R('%v', %v, :1)",
		pLang, pSn))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_NOTICE_DTL_R('%v', %v, :1)",
		pLang, pSn),
		ora.I64, /* SN */
		ora.S,   /* GBN_NM */
		ora.S,   /* TITLE */
		ora.S,   /* CONT */
		ora.S,   /* REG_DT */
		ora.S,   /* NEW_YN */
		ora.S,   /* GBN_CD */
		ora.S,   /* MEM_CD */
		ora.S,   /* EPS_YN */
	)
	defer stmtProcCall.Close()
	if err != nil {
		panic(err)
	}
	procRset := &ora.Rset{}
	_, err = stmtProcCall.Exe(procRset)

	if err != nil {
		panic(err)
	}

	var (
		rtnCd  int64
		rtnMsg string
		sn     int64
		title  string
		//regDt  string
		//cont string
		//newYn  string
		gbnCd string
		memCd string
		//epsYn  string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			sn = procRset.Row[0].(int64)
			//memNm = procRset.Row[1].(string)
			title = procRset.Row[2].(string)
			//regDt = procRset.Row[3].(string)
			//cont = procRset.Row[4].(string)
			//newYn = procRset.Row[5].(string)
			gbnCd = procRset.Row[6].(string)
			memCd = procRset.Row[7].(string)
			//epsYn = procRset.Row[8].(string)
		}

		if err := procRset.Err(); err != nil {
			panic(err)
		}

		rtnCd = 1
	}

	rtnNoticePush := models.RtnNoticePush{
		RtnCd:  rtnCd,
		RtnMsg: rtnMsg,
	}

	if rtnCd == 1 {

		gbn := "10"

		// 푸쉬 발송 시간과 몇번째 보내는지를 적자.
		SetPushSendInfoRegister(ses, pSn, "S", -1, 0)

		go NoticeFCM(gbn, memCd, title, sn, gbnCd)
	}

	// End : Admin Event Process

	c.Data["json"] = &rtnNoticePush
	c.ServeJSON()
	//c.TplName = "admin/notice_write.html"
}

func SetPushSendInfoRegister(ses *ora.Ses, sn string, gbn string, pushSendResCnt int, pushSendResTotalCnt int) {
	// 이미지 등록
	pLang, _ := beego.AppConfig.String("lang")

	logs.Debug(fmt.Sprintf("CALL SP_EMS_PUSH_SEND_INFO_UPDATE( '%v', '%v', %v, %v, %v, :1)",
		pLang, gbn, sn, pushSendResCnt, pushSendResTotalCnt))

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_PUSH_SEND_INFO_UPDATE( '%v', '%v', %v, %v, %v, :1)",
		pLang, gbn, sn, pushSendResCnt, pushSendResTotalCnt),
		ora.I64, ora.S)

	defer stmtProcCall.Close()
	if err != nil {
		panic(err)
	}
	procRset := &ora.Rset{}
	_, err = stmtProcCall.Exe(procRset)

	if err != nil {
		panic(err)
	}
}

func SendPushKeyToUser(registrationTokens []string, contMsg string, gbn string, mem_cd string, cont string, temp_Sn string, brdgbncd string) int {

	opt := option.WithCredentialsFile("qrate-2ee14-firebase-adminsdk-64reu-74554f5c44.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logs.Debug("error initializing app: %v\n", err)
	}

	// [START send_to_token_golang]
	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	client, err := app.Messaging(ctx)

	message := &messaging.MulticastMessage{
		Data: map[string]string{
			"type":     gbn,
			"body":     contMsg,
			"brdgbncd": brdgbncd,
			"sn":       temp_Sn,
		},
		Notification: &messaging.Notification{
			Body: contMsg,
		},
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				Icon:  "stock_ticker_update",
				Color: "#2ad0c7",
			},
		},
		Tokens: registrationTokens,
	}

	br, err := client.SendMulticast(context.Background(), message)
	if err != nil {
		log.Fatalln(err)
	}

	if br.FailureCount > 0 {
		var failedTokens []string
		var errorTokens []string
		for idx, resp := range br.Responses {
			if !resp.Success {
				// The order of responses corresponds to the order of the registration tokens.
				failedTokens = append(failedTokens, registrationTokens[idx]+"\n")
				//fmt.Printf(resp.Error.Error())
				errorTokens = append(errorTokens, resp.Error.Error()+"\n")
			}
		}

		logs.Debug("List of tokens that caused failures: ", failedTokens, "\n")
		logs.Debug("List of ErrorCode : ", errorTokens, "\n")
	}
	// See the BatchResponse reference documentation
	// for the contents of response.
	return br.SuccessCount
}

func NoticeFCM(gbn string, mem_cd string, cont string, sn int64, brdgbncd string) {

	// start : log
	slog := logs.NewLogger()
	slog.SetLogger(logs.AdapterConsole)
	// end : log

	//session := c.StartSession()
	pLang, _ := beego.AppConfig.String("lang")

	var (
		// token *string
		//pushagreeyn string
		pushResCnt int
		contMsg    string
	)

	// brdgbncd: '01' 이면 공지 02면 이벤트
	if brdgbncd == "02" {
		contMsg = "(광고) " + cont + "\n수신거부:설정에서 변경 가능"
	} else {
		contMsg = cont
	}

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	logs.Debug(contMsg)

	// Start : Certification Key Info

	logs.Debug(fmt.Sprintf("CALL SP_EMS_PUSH_KEY_INFO_S_R('%v', '%v', '%v', '%v', :1)",
		pLang, gbn, brdgbncd, mem_cd))
	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_PUSH_KEY_INFO_S_R('%v', '%v', '%v', '%v', :1)",
		pLang, gbn, brdgbncd, mem_cd),
		ora.S, /* PUSH_KEY */
		// ora.S, /* CONT */
	)
	defer stmtProcCall.Close()
	if err != nil {
		panic(err)
	}
	procRset := &ora.Rset{}
	_, err = stmtProcCall.Exe(procRset)

	if err != nil {
		panic(err)
	}

	temp_Sn := strconv.Itoa(int(sn))

	//rset_len = procRset.Len()
	// 그냥 10만건 만들자.
	push_key_array := make([]*string, 0, 100000)

	push_key_map := map[string]int{}

	if procRset.IsOpen() {
		for procRset.Next() {
			token := procRset.Row[0].(string)
			// pushagreeyn = procRset.Row[1].(string)
			//slog.Debug("token : %v ", token)

			_, ok := push_key_map[token] // 중복 체크 한다.
			if !ok {
				push_key_array = append(push_key_array, &token)
				tmp_array_size := len(push_key_array)
				push_key_map[token] = tmp_array_size
			}
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}

	array_size := len(push_key_array)
	map_size := len(push_key_map)

	logs.Debug("array_size: ", array_size, " map_size: ", map_size)

	def_send_size := 100

	registrationTokens := make([]string, 0, def_send_size)
	// 여기 까지 왔다면 100개씩 푸쉬를 보내면 된다.
	for i := range push_key_array {
		fmt.Printf("%v\n", *push_key_array[i])

		registrationTokens = append(registrationTokens, *push_key_array[i])

		if (i+1)%def_send_size == 0 {
			//			slog.Info("registrationTokens: 100 %v ", registrationTokens)
			pushResCnt += SendPushKeyToUser(registrationTokens, contMsg, gbn, mem_cd, cont, temp_Sn, brdgbncd)
			registrationTokens = make([]string, 0, def_send_size)
		}
	}

	if len(registrationTokens) > 0 {
		pushResCnt += SendPushKeyToUser(registrationTokens, contMsg, gbn, mem_cd, cont, temp_Sn, brdgbncd)
	}

	logs.Debug(pushResCnt, "/", array_size, " messages were sent successfully\n")

	SetPushSendInfoRegister(ses, temp_Sn, "E", pushResCnt, array_size)
}
