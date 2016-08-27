package controllers

import (
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/chanxuehong/rand"
	"github.com/chanxuehong/wechat.v2/mp/jssdk"
)

var (
	jssdkClient *jssdk.DefaultTicketServer = jssdk.NewDefaultTicketServer(wechatClient)
)

type SystemController struct {
	beego.Controller
}

func (this *SystemController) Prepare() {
	this.Ctx.ResponseWriter.Header().Add("Access-Control-Allow-Origin", "*")
}

/**
 * @api {get} /system 获取系统信息
 * @apiName systemInfo
 * @apiGroup System
 *
 * @apiParam {String} pageUrl 网页地址，用于JSSDK认证
 *
 * @apiParamExample {String} Request-Example:
 * pageUrl=http://www.xxx.com/index.html
 *
 * @apiSuccess {String} appId 公众号的唯一标识
 * @apiSuccess {Number} timestamp 生成签名的时间戳
 * @apiSuccess {String} nonceStr 生成签名的随机串
 * @apiSuccess {String} signature 签名
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 0,
 *         "errmsg": "操作成功",
 *         "data": {
 *             "appId": "1234567",
 *             "timestamp": 1234567
 *             "nonceStr": "nonceStr",
 *             "signature": "signature"
 *         }
 *     }
 *
 * @apiErrorExample Error-Response:
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 403,
 *         "errmsg": "Refer认证失败，请求被拒绝",
 *         "data": {
 *         }
 *     }
 */

// @router / [get]
func (this *SystemController) Get() {

	pageUrl := this.GetString("pageUrl")
	timestamp := time.Now().Unix()
	nonceStr := string(rand.NewHex())
	jssdkClientTicket, err := jssdkClient.Ticket()
	if err != nil {
		this.Abort("500")
	}
	signature := jssdk.WXConfigSign(jssdkClientTicket, nonceStr, string(timestamp), pageUrl)
	fmt.Println("signature =", signature)
	this.Data["json"] = restReturn(0, "操作成功", map[string]interface{}{
		"appId":     wxAppId,
		"timestamp": timestamp,
		"nonceStr":  nonceStr,
		"signature": signature,
	})
	this.ServeJSON()
}
