package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
)

// ServerFileLogController ...
type ServerFileLogController struct {
	beego.Controller
}

// Post ...
func (c *ServerFileLogController) Post() {

	pTitle := c.GetString("title")
	pMessage := c.GetString("message")

	fmt.Printf(fmt.Sprintf("[ServerLog][%v] %v", pTitle, pMessage))

	c.Data["json"] = &models.DefaultResult{RtnCd: 1, RtnMsg: pMessage}
	c.ServeJSON()
}
