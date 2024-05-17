package controllers

import (
	"go-beego-exapmple-02/models"

	beego "github.com/beego/beego/v2/server/web"
)

type TestController struct {
	beego.Controller
}

// URLMapping ...
// func (c *TestController) URLMapping() {	
// 	c.Mapping("Get", c.Get)	
// }

// @Title hello
// @Description Tests the API
// @Success 200 {object} models.Test
// @Failure 403 body is empty
// @router /hello [get]
func (t *TestController) Hello() {
	Response := models.TestFunction()
	t.Data["json"] = Response
	
	t.ServeJSON()  
}