package routers

import (
	"emst.ziggam.com/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	/* main */
	beego.Router("/", &controllers.MainController{})

	/* UEditor*/
	beego.Router("/controller", &controllers.UeditorController{}, "*:ControllerUE")

	/* notice */
	beego.Router("/notice/list", &controllers.NoticeListController{})     //공지사항 리스트
	beego.Router("/notice/detail", &controllers.NoticeDetailController{}) //공지사항 상세

	/* inquiry */
	beego.Router("/inquiry/list", &controllers.InquiryListController{})     //문의사항 리스트
	beego.Router("/inquiry/write", &controllers.InquiryWriteController{})   //문의사항 작성하기
	beego.Router("/inquiry/insert", &controllers.InquiryInsertController{}) //문의사항 등록하기

	/* common */
	beego.Router("/common/login", &controllers.CommonLoginController{})                 //로그인폼
	beego.Router("/common/logout", &controllers.CommonLogoutController{})               //로그아웃
	beego.Router("/common/find/id", &controllers.CommonFindIdController{})              //아이디찾기
	beego.Router("/common/find/pwd", &controllers.CommonFindPwdController{})            //비밀번호찾기
	beego.Router("/common/find/pwd/cert", &controllers.CommonFindPwdCertController{})   //비밀번호찾기 인증번호(이메일 인증)
	beego.Router("/common/pwd/reset", &controllers.CommonPwdResetController{})          //비밀번호 재설정(변경)
	beego.Router("/common/standby/email", &controllers.CommonStandbyEmailController{})  //이메일 승인대기
	beego.Router("/common/cert/resend", &controllers.CommonCertEmailResendController{}) //이메일 재인증
	beego.Router("/common/standby/auth", &controllers.CommonStandbyAuthController{})    //인증대기
	beego.Router("/common/jobgrp", &controllers.CommonJobGrpController{})               //직군
	beego.Router("/common/jobgrp2", &controllers.CommonJobGrp2Controller{})             //직무
	beego.Router("/common/change/pwd", &controllers.CommonChangePwdController{})        //비밀번호변경
	beego.Router("/common/dup_chk", &controllers.CommonDupChkController{})              //항목 중복 체크
	beego.Router("/common/code/list", &controllers.CommonCodeListController{})          //공통코드리스트
	beego.Router("/common/item/list", &controllers.CommonItemListController{})          //공통항목리스트(내부링크용)
	beego.Router("/common/recruit/list", &controllers.CommonRecruitListController{})    //공통채용공고리스트(내부링크용)

	/* Error */
	beego.Router("/error", &controllers.ErrorController{}, "get:Error404") //Error
	beego.Router("/error/404", &controllers.Error404Controller{})          //404 Error

	/* Web */
	beego.Router("/common/email/cert", &controllers.CommonEmailCertController{}) //이메일인증

	/* iframe */
	beego.Router("/entp/service/use", &controllers.EntpServiceUseController{}) //이용약관 iframe

	/* recruit */
	beego.Router("/recruit/post/list", &controllers.RecruitPostListController{})     //채용공고 리스트
	beego.Router("/recruit/post/detail", &controllers.RecruitPostDetailController{}) //채용공고 상세
	beego.Router("/recruit/post/write", &controllers.RecruitPostWriteController{})   //채용공고 작성하기
	beego.Router("/recruit/post/modify", &controllers.RecruitPostModifyController{}) //채용공고 수정하기
	// 안쓴다.. beego.Router("/recruit/post/insert", &controllers.RecruitPostInsertController{})                //채용공고 등록처리
	beego.Router("/recruit/post/update", &controllers.RecruitPostUpdateController{})                //채용공고 수정처리
	beego.Router("/recruit/post/end", &controllers.RecruitPostEndController{})                      //채용공고 종료처리
	beego.Router("/recruit/stat/list", &controllers.RecruitStatListController{})                    //지원현황 리스트
	beego.Router("/recruit/apply/member/delete", &controllers.RecruitApplyMemberDeleteController{}) //지원자 삭제
	beego.Router("/recruit/apply/detail", &controllers.RecruitApplyDetailController{})              //지원자 상세
	// 안쓴다. beego.Router("/recruit/eval/update", &controllers.EvaluateUpdateController{})                   //채용결정/포기
	beego.Router("/recruit/eval/check", &controllers.EvaluateCheckController{})                   //채용결정확인
	beego.Router("/recruit/favor/member/update", &controllers.FavorMemberUpdateController{})      //관심설정/해제
	beego.Router("/recruit/apply/member/excel", &controllers.RecruitApplyMemberExcelController{}) //지원현황 엑셀다운로드
	beego.Router("/recruit/apply/delete", &controllers.RecruitApplyDeleteController{})            //지원현황 엑셀다운로드

	/* applicant */
	beego.Router("/applicant/list", &controllers.ApplicantListController{}) //지원자관리 리스트

	/* management */
	beego.Router("/admin/notice/list", &controllers.AdminNoticeListController{})              //관리자 공지사항 리스트
	beego.Router("/admin/notice/write", &controllers.AdminNoticeWriteController{})            //관리자 공지사항 작성하기
	beego.Router("/admin/notice/process", &controllers.AdminNoticeProcessController{})        //관리자 공지사항 등록/수정처리
	beego.Router("/admin/notice/push", &controllers.NoticePushController{})                   //관리자 공지/이벤트 사항 푸쉬 처리.
	beego.Router("/admin/notice/delete", &controllers.AdminNoticeDeleteController{})          //관리자 공지사항 삭제처리
	beego.Router("/admin/event/content/list", &controllers.AdminEventContentListController{}) //관리자 이벤트(컨텐츠) 리스트
	beego.Router("/admin/event/write", &controllers.AdminEventWriteController{})              //관리자 이벤트 작성하기
	// 안쓴다. adminNoticePro~로 통합됨. beego.Router("/admin/event/process", &controllers.AdminEventProcessController{})                   //관리자 이벤트 등록/수정처리
	beego.Router("/admin/event/delete", &controllers.AdminEventDeleteController{})                     //관리자 이벤트 삭제처리
	beego.Router("/admin/inquiry/list", &controllers.AdminInquiryListController{})                     //관리자 문의사항 리스트
	beego.Router("/admin/inquiry/insert", &controllers.AdminInquiryInsertController{})                 //관리자 문의답변 처리
	beego.Router("/admin/service/list", &controllers.AdminServiceListController{})                     //관리자 서비스정책 리스트
	beego.Router("/admin/service/process", &controllers.AdminServiceProcessController{})               //관리자 서비스정책 수정처리
	beego.Router("/admin/entp/list", &controllers.AdminEntpListController{})                           //관리자 기업 리스트
	beego.Router("/admin/entp/info", &controllers.AdminEntpInfoController{})                           //관리자 기업 상세정보
	beego.Router("/admin/entp/video/insert", &controllers.AdminEntpVideoInsertController{})            //관리자 기업 영상 등록처리
	beego.Router("/admin/entp/video/delete", &controllers.AdminEntpVideoDeleteController{})            //관리자 기업 영상 삭제처리
	beego.Router("/admin/entp/email/cert", &controllers.AdminEntpEmailCertController{})                //관리자 기업 이메일인증
	beego.Router("/admin/login", &controllers.AdminLoginController{})                                  //관리자 간편 로그인처리
	beego.Router("/admin/member/list", &controllers.AdminMemberListController{})                       //관리자 개인 회원가입 리스트
	beego.Router("/admin/version/info", &controllers.AdminVersionInfoController{})                     //관리자 앱 버전관리
	beego.Router("/admin/version/insert", &controllers.AdminVersionInsertController{})                 //관리자 앱 버전 등록처리
	beego.Router("/admin/entp/video/insert/test", &controllers.AdminEntpVideoInsertTestController{})   //관리자 기업 영상 등록처리(테스트)
	beego.Router("/admin/entp/info/test", &controllers.AdminEntpInfoTestController{})                  //관리자 기업 상세정보
	beego.Router("/admin/control", &controllers.AdminControlController{})                              //관리자 제재처리
	beego.Router("/admin/event/list", &controllers.AdminEventListController{})                         //관리자 이벤트 리스트
	beego.Router("/admin/event/video", &controllers.AdminEventVideoController{})                       //관리자 이벤트 영상뷰
	beego.Router("/admin/intro/popup/list", &controllers.AdminIntroPopUpListController{})              //관리자 INTRO팝업 리스트
	beego.Router("/admin/intro/write", &controllers.AdminIntroPopUpWriteController{})                  //관리자 INTRO팝업 등록폼
	beego.Router("/admin/intro/popup/insert", &controllers.AdminIntroPopUpInsertController{})          //관리자 INTRO팝업 등록
	beego.Router("/admin/intro/popup/stat/update", &controllers.AdminIntroPopUpStatUpdateController{}) //관리자 INTRO팝업 활성상태 처리
	beego.Router("/admin/video/use/update", &controllers.AdminEntpVideoUseUpdateController{})          //관리자 기업영상 검증처리
	beego.Router("/admin/member/detail", &controllers.AdminMemberDetailController{})                   //관리자 회원 상세
	beego.Router("/admin/qa/history/list", &controllers.AdminQaHistoryListController{})                //관리자 채용 질문답변 리스트

	beego.Router("/admin/member/list/excel", &controllers.AdminMemberListExcelController{}) //관리자 회원 리스트 엑셀다운로드 LDK 2020/12/08
	beego.Router("/admin/entp/list/excel", &controllers.AdminEntpListExcelController{})     //관리자 기업 리스트 엑셀다운로드 LDK 2020/12/11

	beego.Router("/admin/banner/list", &controllers.AdminBannerListController{})              //관리자 배너 리스트
	beego.Router("/admin/banner/write", &controllers.AdminBannerWriteController{})            //관리자 배너 등록폼
	beego.Router("/admin/banner/insert", &controllers.AdminBannerInsertController{})          //관리자 배너 등록
	beego.Router("/admin/banner/stat/update", &controllers.AdminBannerStatUpdateController{}) //관리자 배너 활성상태 처리
	beego.Router("/admin/banner/kind/list", &controllers.AdminBannerKindListController{})     //관리자 배너 종류 리스트(기업,채용공고)
	beego.Router("/admin/banner/time/update", &controllers.AdminBannerTimeUpdateController{}) //관리자 배너 롤링타임 처리

	/* management > statstics*/
	beego.Router("/admin/stats/main", &controllers.AdminStatsMainController{})                  //관리자 통계 회원 메인
	beego.Router("/admin/stats/member/detail", &controllers.AdminStatsMemberDetailController{}) //관리자 통계 > 회원상세 통계
	beego.Router("/admin/stats/recruit/main", &controllers.AdminStatsRecruitMainController{})   //관리자 통계 채용 메인
	beego.Router("/admin/stats/period/api", &controllers.AdminStatsPeriodApiController{})       // 기간별 지원 API
	beego.Router("/admin/stats/period/main", &controllers.AdminStatsPeriodMainController{})     // 기간별 지원 통계

	/* enterprise */
	beego.Router("/entp/write/step1", &controllers.EntpWriteStep1Controller{})   //기업 회원가입 작성하기 Step1
	beego.Router("/entp/write/step2", &controllers.EntpWriteStep2Controller{})   //기업 회원가입 작성하기 Step2
	beego.Router("/entp/write/step3", &controllers.EntpWriteStep3Controller{})   //기업 회원가입 작성하기 Step3
	beego.Router("/entp/insert", &controllers.EntpInsertController{})            //기업 회원가입 등록처리
	beego.Router("/entp/insert/login", &controllers.EntpInsertLoginController{}) //기업 회원가입 등록 후 로그인처리
	beego.Router("/entp/info/write", &controllers.EntpInfoWriteController{})     //기업정보 작성하기
	beego.Router("/entp/info/update", &controllers.EntpInfoUpdateController{})   //기업정보 수정하기

	beego.Router("/entp/info/update/v2", &controllers.EntpInfoUpdateV2Controller{}) //기업정보 수정하기(일반 기업)

	/* notification */
	beego.Router("/notification/list", &controllers.NotificationListController{}) //알림 리스트
	
	/* job fair */
	beego.Router("/jobfair/entp/reg", &controllers.JobfairEntpRegisterController{}) // 박람회 기업 등록
	beego.Router("/jobfair/entp/del", &controllers.JobfairEntpDeleteController{})   // 박람회 기업 삭제
    /* team member */
	beego.Router("/team/member/pwd/update", &controllers.TeamMemberPwdUpdateController{}) // 팀멤버 비밀번호 변경처리
	/* member */
	beego.Router("/member/pwd/update", &controllers.MemberPwdUpdateController{}) // 개인 회원 비밀번호 변경처리

}
