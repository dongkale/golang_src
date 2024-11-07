package routers

import (
	"bizt.ziggam.com/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	/* main */
    beego.Router("/", &controllers.MainController{})
	beego.Router("/main/applicant/list", &controllers.MainApplicantListController{})

	/* Error */
	beego.Router("/error", &controllers.ErrorController{}, "get:Error404") //Error
	beego.Router("/error/404", &controllers.Error404Controller{})          //404 Error

	/* join */
	beego.Router("/join", &controllers.JoinController{})
	beego.Router("/join/insert", &controllers.JoinInsertController{})
	beego.Router("/join/complete", &controllers.JoinCompleteController{})

	/* login */
	beego.Router("/login", &controllers.LoginController{}) 
	beego.Router("/logout", &controllers.LogoutController{})

	/* admin login */
	beego.Router("/admin/simple/pass", &controllers.AdminSimplePassController{})
	beego.Router("/admin/simple/pass/login", &controllers.AdminSimplePassLoginController{})

	/* common */
	beego.Router("/common/dup_chk", &controllers.CommonDupChkController{})
	beego.Router("/common/team/dup_chk", &controllers.CommonTeamDupChkController{})
	beego.Router("/common/id/find", &controllers.IdFindController{})
	beego.Router("/common/id/find/result", &controllers.IdFindResultController{})
	beego.Router("/common/pwd/find", &controllers.PwdFindController{})
	beego.Router("/common/pwd/find/cert", &controllers.PwdFindCertController{})
	beego.Router("/common/pwd/reset", &controllers.PwdResetController{})
	beego.Router("/common/jobgrp", &controllers.CommonJobGrpController{})
	beego.Router("/common/jobgrp2", &controllers.CommonJobGrp2Controller{})

	beego.Router("/common/rgngrp", &controllers.CommonRgnGrpController{})

	/* entp */
	beego.Router("/entp/info/modify", &controllers.EntpInfoModifyController{})
	beego.Router("/entp/info/update", &controllers.EntpInfoUpdateController{})
	beego.Router("/entp/video", &controllers.EntpVideoController{})

	/* message */
	beego.Router("/message", &controllers.MessageController{})
	beego.Router("/message/list", &controllers.MessageListController{})
	beego.Router("/message/list/count", &controllers.MessageListCountController{})
	beego.Router("/message/detail", &controllers.MessageDetailController{})
	beego.Router("/message/confirm/update", &controllers.MessageConfirmUpdateController{})
	beego.Router("/message/send", &controllers.MessageSendController{})

	/* notification */
	beego.Router("/noti/temp/update", &controllers.NotiTempUpdateController{})
	beego.Router("/msg/temp/update", &controllers.MsgTempUpdateController{})
	beego.Router("/notification/list", &controllers.NotificationListController{})
	beego.Router("/notification/update", &controllers.NotificationUpdateController{})

	/* member */
	beego.Router("/favor/member/set", &controllers.FavorMemberUpdateController{})

	/* recruit */
	beego.Router("/recruit/write", &controllers.RecruitWriteController{})                       
	beego.Router("/recruit/insert", &controllers.RecruitInsertController{})                     
	beego.Router("/recruit/modify", &controllers.RecruitModifyController{})                     
	beego.Router("/recruit/update", &controllers.RecruitUpdateController{})                     
	beego.Router("/recruit/copy", &controllers.RecruitCopyController{})                         
	beego.Router("/recruit/post/list", &controllers.RecruitPostListController{})                
	beego.Router("/recruit/post/detail", &controllers.RecruitPostDetailController{})            
	beego.Router("/recruit/post/detail/popup", &controllers.RecruitPostDetailPopupController{}) 
	beego.Router("/recruit/post/detail/excel", &controllers.RecruitPostDetailExcelController{}) 
	beego.Router("/recruit/post/detail/start", &controllers.RecruitPostDetailStartController{}) 
	beego.Router("/recruit/post/detail/end", &controllers.RecruitPostDetailEndController{})     
	beego.Router("/recruit/delete", &controllers.RecruitDeleteController{})                     
	beego.Router("/recruit/main/list", &controllers.RecruitMainListController{})                

	/* applicant */
	beego.Router("/applicant/list", &controllers.ApplicantListController{})                
	beego.Router("/applicant/list/excel", &controllers.ApplicantListExcelController{})     
	beego.Router("/applicant/detail/popup", &controllers.ApplicantDetailPopupController{}) 
	beego.Router("/applicant/delete", &controllers.ApplicantDeleteController{})            
	beego.Router("/applicant/each/delete", &controllers.ApplicantEachDeleteController{})   
	beego.Router("/comment/pop/list", &controllers.CommentPopListController{})             

	/* team member */
	beego.Router("/team/member/list", &controllers.TeamMemberListController{})              
	beego.Router("/team/member/write", &controllers.TeamMemberWriteController{})            
	beego.Router("/team/member/insert", &controllers.TeamMemberInsertController{})          
	beego.Router("/team/member/modify", &controllers.TeamMemberModifyController{})          
	beego.Router("/team/member/update", &controllers.TeamMemberUpdateController{})          
	beego.Router("/team/member/delete", &controllers.TeamMemberDeleteController{})          
	beego.Router("/team/member/auth/update", &controllers.TeamMemberAuthUpdateController{}) 
	beego.Router("/team/member/pwd/modify", &controllers.TeamMemberPwdModifyController{})   
	beego.Router("/team/member/pwd/update", &controllers.TeamMemberPwdUpdateController{})   

	/* team comment */
	beego.Router("/team/comment/list", &controllers.TeamCommentListController{})     
	beego.Router("/team/comment/insert", &controllers.TeamCommentInsertController{}) 
	beego.Router("/team/comment/delete", &controllers.TeamCommentDeleteController{}) 

	/* message verify */
	beego.Router("/message/verify", &controllers.MessageVerifyController{}) 

	/* setting */
	beego.Router("/setting/member/modify", &controllers.SettingMemberModifyController{})        
	beego.Router("/setting/member/pwd/modify", &controllers.SettingMemberPwdModifyController{}) 
	beego.Router("/setting/member/pwd/update", &controllers.SettingMemberPwdUpdateController{}) 

	beego.Router("/setting/notice/list", &controllers.NoticeListController{})     
	beego.Router("/setting/notice/detail", &controllers.NoticeDetailController{}) 

	beego.Router("/setting/inquiry/list", &controllers.InquiryListController{})       
	beego.Router("/setting/inquiry/detail", &controllers.InquiryDetailController{})   
	beego.Router("/setting/inquiry/write", &controllers.InquiryWriteController{})     
	beego.Router("/setting/inquiry/insert", &controllers.InquiryInsertController{})   
	beego.Router("/setting/withdraw", &controllers.WithDrawController{})              
	beego.Router("/setting/withdraw/update", &controllers.WithDrawUpdateController{}) 

	/* live interview */
	beego.Router("/live/intro", &controllers.LiveIntroController{})
	beego.Router("/live/list", &controllers.LiveListController{})
	beego.Router("/live/detail", &controllers.LiveDetailController{})

	beego.Router("/live/itv/popup", &controllers.LiveItvPopupController{})
	beego.Router("/live/itv/exit", &controllers.LiveItvExitController{})
	beego.Router("/live/itv/close", &controllers.LiveItvCloseController{})
	//beego.Router("/live/itv/result", &controllers.LiveItvResultController{})

	beego.Router("/live/conference/popup", &controllers.LiveNvnPopupController{})
	beego.Router("/live/conference/register/channelId", &controllers.LiveNvnRegKeyController{})
	beego.Router("/live/conference/get/channelId", &controllers.LiveNvnGetKeyController{})

	beego.Router("/live/nvn/exit", &controllers.LiveNvnExitController{})
	beego.Router("/live/nvn/close", &controllers.LiveNvnCloseController{})

	beego.Router("/api/apply/result", &controllers.ApiSetApplyResultController{})
	beego.Router("/api/message/send", &controllers.ApiMessageSendController{})
	beego.Router("/api/recruit/insert", &controllers.ApiRecruitInsertController{})
	beego.Router("/api/apisvr/cmd", &controllers.ApiApiServerPostReqController{})

	/* bridge page */
	beego.Router("/bridge", &controllers.BridgeController{})

	beego.Router("/invite/refuse", &controllers.InviteRefuseController{})

	beego.Router("/invite/recurit/list", &controllers.InviteRecruitListController{})

	beego.Router("/invite/send", &controllers.InviteSendController{})
	beego.Router("/invite/send/list", &controllers.InviteSendListController{})
	beego.Router("/invite/send/list/detail/popup", &controllers.InviteSendListDetailPopupController{})
	beego.Router("/invite/send/msg/preview", &controllers.InviteSendMsgPreviewController{})

	beego.Router("/live/nvn/req/popup", &controllers.LiveNvNReqPopupController{})
	beego.Router("/live/nvn/req/added/popup", &controllers.LiveNvNReqAddedPopupController{})

	beego.Router("/live/nvn/recruit/apply/list", &controllers.LiveNvNRecuruitApplyListController{})

	beego.Router("/live/nvn/proc", &controllers.LiveNvNProcController{})
	beego.Router("/live/nvn/list", &controllers.LiveNvNListController{})
	beego.Router("/live/nvn/detail", &controllers.LiveNvNDetailController{})

	beego.Router("/live/nvn/stat/verify", &controllers.LiveNvNStatVerifyController{})
	beego.Router("/live/nvn/joincnt", &controllers.LiveNvNJoinCntController{})

	beego.Router("/applicant/score/popup", &controllers.ApplicantScorePopupController{})
	beego.Router("/applicant/score/reg/update", &controllers.ApplicantScoreRegUpdateController{})

	/* applicant all */
	beego.Router("/applicant/list/all", &controllers.ApplicantListAllController{})
	beego.Router("/applicant/detail/popup/p2", &controllers.ApplicantDetailPopupP2Controller{})

	beego.Router("/ads/popup", &controllers.ADSPopupController{})

	beego.Router("/api/jobfair/recruit/apply", &controllers.JobFairRecruitApplyController{})
	beego.Router("/api/jobfair/live/request", &controllers.JobFairLiveRequestController{})

	/* Utils */
	beego.Router("/utils", &controllers.UtilsController{})
	beego.Router("/reload/template", &controllers.ReloadTemplateController{})
	beego.Router("/reload/config", &controllers.ReloadConfigController{})

	/* Setver to file log */
	beego.Router("/server/filelog", &controllers.ServerFileLogController{})

	/* LDK test */
	beego.Router("/a_test", &controllers.A_TestController{})

	/* temp */
	beego.Router("/temp", &controllers.TempController{})

	/* excel upload */
	beego.Router("/excel", &controllers.ExcelUploadController{})
}
