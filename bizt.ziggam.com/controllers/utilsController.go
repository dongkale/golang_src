package controllers

import (
	"fmt"	
)

// UtilsController ...
type UtilsController struct {
	BaseController
}

// Get ...
func (c *UtilsController) Get() {
	session := c.StartSession()

	memNo := session.Get(c.Ctx.Request.Context(), "mem_no")
	if memNo == nil {
		c.Ctx.Redirect(302, "/login")
	} else {
		pPw := c.GetString("pw")

		fmt.Printf(fmt.Sprintf("[UtilsController] Start"))

		c.Data["TMenuId"] = "T00"
		c.Data["SMenuId"] = "T00"

		if pPw == "dlehdrhks" {
			fmt.Printf(fmt.Sprintf("[UtilsController] Entp_mem_No: %v", memNo))
			c.TplName = "utils/utils.html"
		} else {
			fmt.Printf(fmt.Sprintf("[UtilsController] Entp_mem_No: %v Invalid User !!", memNo))
			c.TplName = "error/error.html"
		}

		fmt.Printf(fmt.Sprintf("[UtilsController] End"))
	}
}
