package controllers

import beego "github.com/beego/beego/v2/server/web"

type ADSPopupController struct {
	beego.Controller
}

func (c *ADSPopupController) Get() {

	// session := c.StartSession()

	// mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")

	// pLang, _ := beego.AppConfig.String("lang")
	// pEntpMemNo := mem_no

	// imgServer, _  := beego.AppConfig.String("viewpath")

	// fmt.Printf(mem_no)

	// Start : Oracle DB Connection
	env, srv, ses, err := GetRawConnection()
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
	if err != nil {
		panic(err)
	}
	// End : Oracle DB Connection

	c.TplName = "ads/ads_popup.html"
}
