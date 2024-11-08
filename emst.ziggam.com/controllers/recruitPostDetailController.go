package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"

	"emst.ziggam.com/models"
	ora "gopkg.in/rana/ora.v4"
)

type RecruitPostDetailController struct {
	BaseController
}

func (c *RecruitPostDetailController) Get() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()
	mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	if mem_no == nil {
		c.Ctx.Redirect(302, "/common/login")
	}

	pLang, _ := beego.AppConfig.String("lang")
	pEntpMemNo := mem_no                        //c.GetString("entp_mem_no")
	pChkEntpMemNo := c.GetString("entp_mem_no") // 체크 기업회원번호

	if pChkEntpMemNo != "" {
		if pEntpMemNo != pChkEntpMemNo {
			c.Ctx.Redirect(302, "/error/404")
		}
	}
	pRecrutSn := c.GetString("recrut_sn")

	/* Parameter */
	pmKeyword := c.GetString("p_keyword")     // 검색어
	pmEmplTyp := c.GetString("p_empl_typ")    // 고용형태코드
	pmJobGrpCd := c.GetString("p_job_grp_cd") // 직군코드
	pmSortGbn := c.GetString("p_sort_gbn")    // 정렬구분
	pmGbnCd := c.GetString("p_gbn_cd")        // 구분코드
	pmPageNo := c.GetString("p_page_no")      // 페이지번호

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	// Start : Recruit Info Detail
	log.Debug("CALL SP_EMS_RECRUIT_DTL_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn)

	stmtProcCall, err := ses.Prep(fmt.Sprintf("CALL SP_EMS_RECRUIT_DTL_R('%v', '%v', '%v', :1)",
		pLang, pEntpMemNo, pRecrutSn),
		ora.S,   /* ENTP_MEM_NO */
		ora.S,   /* RECRUT_SN */
		ora.S,   /* PRGS_STAT */
		ora.S,   /* RECRUT_TITLE */
		ora.S,   /* EMPL_TYP */
		ora.S,   /* UP_JOB_GRP */
		ora.S,   /* JOB_GRP */
		ora.S,   /* SEX */
		ora.I64, /* RECRUT_CNT */
		ora.S,   /* ANS_LMT_TM_NM */
		ora.S,   /* ACPT_PRID_NM */
		ora.S,   /* RECRUT_DY */
		ora.S,   /* ROL */
		ora.S,   /* APLY_QUFCT */
		ora.S,   /* PERFER_TRTM */
		ora.S,   /* PRGS_MSG */
		ora.S,   /* VD_TITLE1 */
		ora.S,   /* VD_TITLE2 */
		ora.S,   /* VD_TITLE3 */
		ora.S,   /* VD_TITLE4 */
		ora.S,   /* VD_TITLE5 */
		ora.I64, /* PV */
		ora.I64, /* UV */
		ora.I64, /* APPLY_CNT */
		ora.I64, /* ING_CNT */
		ora.I64, /* PASS_CNT */
		ora.I64, /* FAIL_CNT */
		ora.S,   /* RECRUT_EDT */

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

	recruitPostDetail := make([]models.RecruitPostDetail, 0)

	var (
		entpMemNo   string
		recrutSn    string
		prgsStat    string
		recrutTitle string
		emplTyp     string
		upJobGrp    string
		jobGrp      string
		sex         string
		recrutCnt   int64
		ansLmtTmNm  string
		acptPridNm  string
		recrutDy    string
		rol         string
		aplyQufct   string
		perferTrtm  string
		prgsMsg     string
		vdTitle1    string
		vdTitle2    string
		vdTitle3    string
		vdTitle4    string
		vdTitle5    string
		pv          int64
		uv          int64
		applyCnt    int64
		ingCnt      int64
		passCnt     int64
		failCnt     int64
		recrutEdt   string
	)

	if procRset.IsOpen() {
		for procRset.Next() {
			entpMemNo = procRset.Row[0].(string)
			recrutSn = procRset.Row[1].(string)
			prgsStat = procRset.Row[2].(string)
			recrutTitle = procRset.Row[3].(string)
			emplTyp = procRset.Row[4].(string)
			upJobGrp = procRset.Row[5].(string)
			jobGrp = procRset.Row[6].(string)
			sex = procRset.Row[7].(string)
			recrutCnt = procRset.Row[8].(int64)
			ansLmtTmNm = procRset.Row[9].(string)
			acptPridNm = procRset.Row[10].(string)
			recrutDy = procRset.Row[11].(string)
			rol = procRset.Row[12].(string)
			aplyQufct = procRset.Row[13].(string)
			perferTrtm = procRset.Row[14].(string)
			prgsMsg = procRset.Row[15].(string)
			vdTitle1 = procRset.Row[16].(string)
			vdTitle2 = procRset.Row[17].(string)
			vdTitle3 = procRset.Row[18].(string)
			vdTitle4 = procRset.Row[19].(string)
			vdTitle5 = procRset.Row[20].(string)
			pv = procRset.Row[21].(int64)
			uv = procRset.Row[22].(int64)
			applyCnt = procRset.Row[23].(int64)
			ingCnt = procRset.Row[24].(int64)
			passCnt = procRset.Row[25].(int64)
			failCnt = procRset.Row[26].(int64)
			recrutEdt = procRset.Row[27].(string)

			recruitPostDetail = append(recruitPostDetail, models.RecruitPostDetail{
				EntpMemNo:   entpMemNo,
				RecrutSn:    recrutSn,
				PrgsStat:    prgsStat,
				RecrutTitle: recrutTitle,
				EmplTyp:     emplTyp,
				UpJobGrp:    upJobGrp,
				JobGrp:      jobGrp,
				Sex:         sex,
				RecrutCnt:   recrutCnt,
				AnsLmtTmNm:  ansLmtTmNm,
				AcptPridNm:  acptPridNm,
				RecrutDy:    recrutDy,
				Rol:         rol,
				AplyQufct:   aplyQufct,
				PerferTrtm:  perferTrtm,
				PrgsMsg:     prgsMsg,
				VdTitle1:    vdTitle1,
				VdTitle2:    vdTitle2,
				VdTitle3:    vdTitle3,
				VdTitle4:    vdTitle4,
				VdTitle5:    vdTitle5,
				Pv:          pv,
				Uv:          uv,
				ApplyCnt:    applyCnt,
				IngCnt:      ingCnt,
				PassCnt:     passCnt,
				FailCnt:     failCnt,
				RecrutEdt:   recrutEdt,
			})
		}
		if err := procRset.Err(); err != nil {
			panic(err)
		}
	}
	// End : Recruit Info Detail

	c.Data["EntpMemNo"] = entpMemNo
	c.Data["RecrutSn"] = recrutSn
	c.Data["PrgsStat"] = prgsStat
	c.Data["RecrutTitle"] = recrutTitle
	c.Data["EmplTyp"] = emplTyp
	c.Data["UpJobGrp"] = upJobGrp
	c.Data["JobGrp"] = jobGrp
	c.Data["Sex"] = sex
	c.Data["RecrutCnt"] = recrutCnt
	c.Data["AnsLmtTmNm"] = ansLmtTmNm
	c.Data["AcptPridNm"] = acptPridNm
	c.Data["RecrutDy"] = recrutDy
	c.Data["Rol"] = rol
	c.Data["AplyQufct"] = aplyQufct
	c.Data["PerferTrtm"] = perferTrtm
	c.Data["PrgsMsg"] = prgsMsg
	c.Data["VdTitle1"] = vdTitle1
	c.Data["VdTitle2"] = vdTitle2
	c.Data["VdTitle3"] = vdTitle3
	c.Data["VdTitle4"] = vdTitle4
	c.Data["VdTitle5"] = vdTitle5
	c.Data["Pv"] = pv
	c.Data["Uv"] = uv
	c.Data["ApplyCnt"] = applyCnt
	c.Data["IngCnt"] = ingCnt
	c.Data["PassCnt"] = passCnt
	c.Data["FailCnt"] = failCnt
	c.Data["RecrutEdt"] = recrutEdt
	c.Data["MenuId"] = "02"

	/* Parameter */
	c.Data["pKeyword"] = pmKeyword
	c.Data["pEmplTyp1"] = pmEmplTyp
	c.Data["pJobGrpCd"] = pmJobGrpCd
	c.Data["pSortGbn"] = pmSortGbn
	c.Data["pGbnCd"] = pmGbnCd
	c.Data["pPageNo"] = pmPageNo

	c.TplName = "recruit/recruit_detail.html"
}
