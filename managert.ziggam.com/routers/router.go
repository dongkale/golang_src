package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"managert.ziggam.com/controllers"
)

func init()  {
	/* VIEW */
	/** main **/
	beego.Router("/", &controllers.MainController{})
	beego.Router("/main", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{})

	/** main **/

	/** contents **/
	beego.Router("/contents", &controllers.ContentsController{})	// 배너 그룹 단건 조회
	/** contents **/

	/* API */

	/** banner **/
	beego.Router("/api/banner/list", &controllers.ApiGetGroupBannerListController{}) // 배너그룹목록 가져오기
	beego.Router("/api/banner/save", &controllers.ApiSaveBannerController{}) // 배너그룹목록 가져오기
	beego.Router("/api/banner/group/copy", &controllers.ApiCopyBannerGroupController{}) // 배너그룹목록 가져오기
	beego.Router("/api/banner/detail/list", &controllers.ApiGetBannerGroupDetailListController{}) // <-- 배너그룹소속배너리스트 뿌리기
	beego.Router("/api/banner/detail", &controllers.ApiGetBannerDetailController{}) // <-- 배너 상세 가져오기
	beego.Router("/api/banner/list/option", &controllers.ApiGetBannerListOptionController{}) // <-- 배너 리스트타입 옵션 가져오기
	beego.Router("/api/banner/index/update", &controllers.ApiUpdateBannerIndexController{}) // 배너 순서 업데이트

	/** banner **/

	/** content - 롤링, 그리드 타입 **/
	beego.Router("/api/content/link/option/list", &controllers.ApiGetContentLinkOptionController{}) // 컨텐츠링크옵션
	beego.Router("/api/content/save", &controllers.ApiSaveContentController{}) // 컨텐츠추가
	beego.Router("/api/content/delete", &controllers.ApiDeleteContentController{}) // 컨텐츠 삭제
	beego.Router("/api/content/title/update", &controllers.ApiUpdateContentTitleController{}) // 컨텐츠 타이틀 업데이트
	beego.Router("/api/content/rolling/time/update", &controllers.ApiUpdateContentRollingTimeController{}) // 롤링타임 업데이트
	beego.Router("/api/content/banner/use/update", &controllers.ApiUpdateContentUseYnController{}) // 컨텐츠 사용/미사용
	beego.Router("/api/content/simple/update", &controllers.ApiSimpleUpdateContentController{}) // 컨텐츠링크옵션
	beego.Router("/api/content/index/update", &controllers.ApiUpdateContentIndexController{}) // 배너 순서 업데이트
	beego.Router("/api/content/copy", &controllers.ApiCopyContentController{}) // 배너그룹목록 가져오기

	/** content - 기업 리스트 **/
	beego.Router("/api/content/cl/list", &controllers.ApiGetContentClListController{}) // 기업 리스트 목록
	beego.Router("/api/content/cl/list/update", &controllers.ApiUpdateContentClListController{}) // 기업 리스트 목록 업데이트 API
	beego.Router("/api/content/cl/use/update", &controllers.ApiUpdateContentClUseYnController{}) // 기업 리스트 사용여부 업데이트
	beego.Router("/api/content/cl/date/update", &controllers.ApiUpdateContentClPublDateController{}) // 기업 리스트 사용여부 업데이트
	beego.Router("/api/content/cl/add/list", &controllers.ApiGetContentClAddListController{}) // 기업 추가하기 화면 목록
	beego.Router("/api/content/cl/company/search", &controllers.ApiGetCompanyDdOptionController{}) // 기업리스트 드롭다운 옵션 검색

	beego.Router("/api/jobfair/list", &controllers.ApiGetJobfairListController{}) // 기업 추가하기 화면 목록


	/** content **/
	/* API */
}
