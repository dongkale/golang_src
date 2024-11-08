package routers

import (
	"cmst.ziggam.com/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	/* Login */
	beego.Router("/login", &controllers.LoginController{})

	/* main */
    beego.Router("/", &controllers.MainController{})
}
