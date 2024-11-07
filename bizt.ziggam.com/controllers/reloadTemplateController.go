package controllers

import (
	"fmt"

	"bizt.ziggam.com/models"
	beego "github.com/beego/beego/v2/server/web"
)

// ReloadTemplateController ...
type ReloadTemplateController struct {
	BaseController
}

// Post ...
func (c *ReloadTemplateController) Post() {

	// start : log
	//log := logs.NewLogger()
	//log.SetLogger(logs.AdapterConsole)
	// end : log

	session := c.StartSession()

	memNo := session.Get(c.Ctx.Request.Context(), "mem_no")
	if memNo == nil {
		//c.Ctx.Redirect(302, "/login")
		c.Data["json"] = &models.DefaultResult{
			RtnCd:  99,
			RtnMsg: "FAIL",
		}

		c.ServeJSON()
		return
	}

	fmt.Printf(fmt.Sprintf("[ReloadTemplateController] Start"))

	c.Rerender(beego.BConfig.WebConfig.ViewsPath)
	//c.Rerender("views")

	fmt.Printf(fmt.Sprintf("[ReloadTemplateController] reload Forder: %v", beego.BConfig.WebConfig.ViewsPath))

	rtnData := models.DefaultResult{
		RtnCd:  1,
		RtnMsg: "SUCCESS",
	}

	c.Data["json"] = &rtnData

	c.ServeJSON()

	fmt.Printf(fmt.Sprintf("[ReloadTemplateController] End"))
}

// Rerender :
func (c *ReloadTemplateController) Rerender(viewpath string) {

	beego.BuildTemplate(viewpath) //beego.BConfig.WebConfig.ViewsPath)
	//c.Redirect("/", 302)
}
