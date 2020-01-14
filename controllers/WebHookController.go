package controllers

import (
	"fmt"

	"BeeCustom/enums"
	"BeeCustom/utils"
)

const SECRETTOKEN = "bee_custom_auto_pull"

// WebHookController handles WebSocket requests.
type WebHookController struct {
	BaseController
}

func (c *WebHookController) Get() {
	signature := c.Ctx.Request.Header.Get("X-Coding-Signature") //获取加密签名
	//res, err := ioutil.ReadAll(c.Ctx.Request.Body) // for application/json
	palyload := c.GetString("payload")
	sha1 := enums.Hmac(SECRETTOKEN, []byte(palyload)) // for application/x-www-form-urlencoded
	calculateSignature := "sha1=" + sha1              // 重新加密内容
	utils.LogDebug(fmt.Sprintf("web_hook同步状态: %v - %v ", calculateSignature == signature, palyload))
	if calculateSignature == signature {
		enums.Cmd("cd", "", []string{"/root/go/src/BeeCustom"})
		enums.Cmd("git", "", []string{"pull"})
		enums.Cmd("go build", "", []string{""})
		enums.Cmd("supervisorctl", "", []string{"restart", "beepkg"})
	}

	c.Data["json"] = palyload
	c.ServeJSON()
}
