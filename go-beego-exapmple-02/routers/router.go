// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"go-beego-exapmple-02/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",
		// beego.NSNamespace("/object",
		// 	beego.NSInclude(
		// 		&controllers.ObjectController{},
		// 	),
		// ),
		// beego.NSNamespace("/user",
		// 	beego.NSInclude(
		// 		&controllers.UserController{},
		// 	),
		// ),
		beego.NSNamespace("/test",
			beego.NSInclude(
				&controllers.TestController{},				
			),
		),
	)
	beego.AddNamespace(ns)

	// beego.Router("/v1/test", &controllers.TestController{}, "get:Get")

	// ns := beego.NewNamespace("/v1",
	// 	beego.NSNamespace("/test",
	// 		beego.NSInclude(
	// 			&controllers.TestController{},
	// 		),
	// 	),
	// )

	// beego.AddNamespace(ns);

	beego.Router("/v1/test/hello", &controllers.TestController{}, "get:Hello")
}
