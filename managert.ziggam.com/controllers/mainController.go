package controllers

type MainController struct {
	BaseController
}

func (c *MainController) Get() {

	// // start : log
	// log := logs.NewLogger()
	// log.SetLogger(logs.AdapterConsole)
	// // end : log

	pBnrGrpTypCd := c.GetString("bnr_grp_typ_cd")

	// // "" 이면 "-" 기본값으로 설정
	// if pBnrGrpTypCd == "" || pBnrGrpTypCd == "01" {
	// 	pBnrGrpTypCd = "MN"
	// }

	c.Data["BnrGrpTypCd"] = pBnrGrpTypCd;

	c.TplName = "main/main.html"
}
