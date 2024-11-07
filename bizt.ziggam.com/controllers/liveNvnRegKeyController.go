package controllers

import (
	"bizt.ziggam.com/utils"
	"github.com/beego/beego/v2/core/logs"
)

type LiveNvnRegKeyController struct {
	BaseController
}

func (c *LiveNvnRegKeyController) Post() {

	// start : log
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	// end : log

	// 키값은 오면 그냥 입력해 주자.
	// mem_no := session.Get(c.Ctx.Request.Context(), "mem_no")
	// if mem_no == nil {
	// 	c.Ctx.Redirect(302, "/login")
	// }

	//pEntpMemSn := c.GetString("entp_mem_sn")
	pRoomId := c.GetString("room_id")
	pChannelId := c.GetString("channel_id")
	pPpMemNm := c.GetString("my_mem_no")
	gbn_cd := c.GetString("gbn_cd")

	key := "room_id:" + pRoomId

	log.Info("reg remon key : ", key)

	entpGbn := "1"
	if gbn_cd == "" {
		entpGbn = "0"
	}

	_ = utils.RPool.HSet(key, pChannelId, entpGbn+":"+pPpMemNm+":"+gbn_cd)
	// ret := utils.RPool.HSet(key, pChannelId, entpGbn+":"+pPpMemNm+":"+gbn_cd)
	// if ret == nil {
	// 	fmt.Printf("redis Error!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	// 	c.TplName = "error/404.html"
	// 	return
	// }

	utils.RPool.Expire(key, 60*60*24*30)

	c.ServeJSON()
}
