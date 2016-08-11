package controllers

import (
	//	"fmt"
	//	"encoding/json"

	"github.com/astaxie/beego"
	//	"github.com/oikomi/FishChatServer/monitor/models"
)

type DeviceController struct {
	beego.Controller
	username string
	userId   string
}

func (this *DeviceController) Prepare() {
	//	_, actionName := this.GetControllerAndAction()
	this.Ctx.ResponseWriter.Header().Add("Access-Control-Allow-Origin", "*")
	if true {
		ticket := this.GetString("ticket")
		this.username = this.GetString("username")
		if this.username == "" || ticket == "" {
			this.Data["json"] = restReturn(44001, "尚未登录或者登陆失效", map[string]interface{}{})
			this.ServeJSON()
		}
	}
}

/**
 * @api {get} /device 查看用户设备列表
 * @apiName deviceList
 * @apiGroup Device
 *
 *
 * @apiParam {String} username 用户名
 * @apiParam {String} ticket 用户接口调用凭据
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 0,
 *         "errmsg": "操作成功",
 *         "data": [{
 *             "IMEI": "123456789101112",
 *             "nick": "123",
 *             "status": 1,
 *             "work_model": 1,
 *             "volume": 6,
 *             "electricity": 100,
 *             "emeregncyPhone": "13590210000",
 *         }]
 *     }
 */

// @router / [get]
func (this *DeviceController) Get() {
	var data []map[string]interface{}
	data = append(data, map[string]interface{}{
		"IMEI":           "123456789101112",
		"nick":           "123",
		"status":         1,
		"work_model":     1,
		"volume":         6,
		"electricity":    100,
		"emeregncyPhone": "13590210000",
	})
	data = append(data, map[string]interface{}{
		"IMEI":           "12345678910555",
		"nick":           "123",
		"status":         1,
		"work_model":     1,
		"volume":         6,
		"electricity":    100,
		"emeregncyPhone": "13590210000",
	})
	this.Data["json"] = restReturn(0, "操作成功", data)
	this.ServeJSON()
}

/**
 * @api {get} /device/:IMEI 查看用户设备详情
 * @apiName deviceDetail
 * @apiGroup Device
 *
 *
 * @apiParam {String} username 用户名
 * @apiParam {String} ticket 用户接口调用凭据
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 0,
 *         "errmsg": "操作成功",
 *         "data": {
 *             "IMEI": "123456789101112",
 *             "nick": "123",
 *             "status": 1
 *             "work_model": 1,
 *             "volume": 6,
 *             "electricity": 100,
 *             "emeregncyPhone": "13590210000",
 *         }
 *     }
 */

// @router /:IMEI [get]
func (this *DeviceController) Show() {
	IMEI := this.Ctx.Input.Param(":IMEI")
	this.Data["json"] = restReturn(0, "操作成功", map[string]interface{}{
		"IMEI":           IMEI,
		"nick":           "123",
		"status":         1,
		"work_model":     1,
		"volume":         6,
		"electricity":    100,
		"emeregncyPhone": "13590210000",
	})
	this.ServeJSON()
}

/**
* @api {post} /device 用户绑定设备
* @apiName deviceBinding
* @apiGroup Device
*
* @apiParam {String} username 用户名
* @apiParam {String} ticket 用户接口调用凭据
* @apiParam {String} IMEI 设备IMEI
* @apiParam {String} nick 设备昵称
*
* @apiParamExample {String} Request-Example:
* IMEI=1234567891011&nick=abc
*
* @apiSuccessExample Success-Response:
*     HTTP/1.1 200 OK
*     {
*         "errcode": 0,
*         "errmsg": "操作成功",
*         "data": {
*             "IMEI": "123456789101112",
*             "nick": "123",
*             "status": 1
*             "work_model": 1,
*             "volume": 6,
*             "electricity": 100,
*             "emeregncyPhone": "13590210000",
*         }
*     }
 */

// @router / [post]
func (this *DeviceController) Post() {
	this.Data["json"] = restReturn(0, "操作成功", map[string]interface{}{
		"IMEI":           "123456789101112",
		"nick":           "123",
		"status":         1,
		"work_model":     1,
		"volume":         6,
		"electricity":    100,
		"emeregncyPhone": "13590210000",
	})
	this.ServeJSON()
}

/**
 * @api {delete} /device/:IMEI 用户删除绑定设备
 * @apiName deviceDestory
 * @apiGroup Device
 *
 * @apiParam {String} username 用户名
 * @apiParam {String} ticket 用户接口调用凭据
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 0,
 *         "errmsg": "操作成功",
 *         "data": {
 *         }
 *     }
 */

// @router /:IMEI [delete]
func (this *DeviceController) Delete() {
	IMEI := this.Ctx.Input.Param(":IMEI")
	this.Data["json"] = restReturn(0, "Delete操作成功"+IMEI, map[string]interface{}{})
	this.ServeJSON()
}

/**
 * @api {put} /device/:IMEI 更新设备信息
 * @apiDescription 本接口只需要传入需要更新的参数即可，无需更新的无需传入
 * @apiName DeivceUpdate
 * @apiGroup Device
 *
 * @apiParam {String} username 用户名
 * @apiParam {String} ticket 用户接口调用凭据
 * @apiParam {Number} [work_model] 工作模式
 * @apiParam {String} [emeregncyPhone] 设备紧急号码
 * @apiParam {Number} [volume] 设备音量
 * @apiParam {String} [nick] 设备昵称
 *
 * @apiSuccess {Number} messageId 消息ID
 *
 * @apiParamExample {String} Request-Example:
 * 传入需要更新的参数即可，无需更新的无需传入
 * work_model=1234567891011&nick=abc&volume=6&emeregncyPhone=13590210000
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 0,
 *         "errmsg": "操作成功",
 *         "data": {
 *             "messageId": 123,
 *         }
 *     }
 */

// @router /:IMEI [put]
func (this *DeviceController) Put() {
	IMEI := this.Ctx.Input.Param(":IMEI")
	this.Data["json"] = restReturn(0, "Put操作成功"+IMEI, map[string]interface{}{})
	this.ServeJSON()
}

/**
 * @api {post} /device/:IMEI/action/location 设备实时定位
 * @apiName DeivceUpdateLocation
 * @apiGroup Device
 *
 * @apiParam {String} username 用户名
 * @apiParam {String} ticket 用户接口调用凭据
 *
 * @apiSuccess {Number} messageId 消息ID
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 0,
 *         "errmsg": "操作成功",
 *         "data": {
 *             "messageId": 123,
 *         }
 *     }
 */

// @router /:IMEI/action/location [post]
func (this *DeviceController) PostActionLocation() {
	IMEI := this.Ctx.Input.Param(":IMEI")
	this.Data["json"] = restReturn(0, "操作成功设备实时定位"+IMEI, map[string]interface{}{
		"messageId": 123,
	})
	this.ServeJSON()
}

/**
 * @api {post} /device/:IMEI/action/shutdown 设备关机
 * @apiName DeivceShutdown
 * @apiGroup Device
 *
 * @apiParam {String} username 用户名
 * @apiParam {String} ticket 用户接口调用凭据
 *
 * @apiSuccess {Number} messageId 消息ID
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 0,
 *         "errmsg": "操作成功",
 *         "data": {
 *             "messageId": 123,
 *         }
 *     }
 */
// @router /:IMEI/action/shutdown [post]
func (this *DeviceController) PostActionShutdown() {
	IMEI := this.Ctx.Input.Param(":IMEI")
	this.Data["json"] = restReturn(0, "操作成功设备关机"+IMEI, map[string]interface{}{
		"messageId": 123,
	})
	this.ServeJSON()
}

/**
 * @api {post} /device/:IMEI/voice 发送聊天语音
 * @apiName DeivceSendVoice
 * @apiGroup Device
 *
 * @apiParam {String} username 用户名
 * @apiParam {String} ticket 用户接口调用凭据
 * @apiParam {String} [mp3Url] mp3地址（和wechatMediaId二选一）
 * @apiParam {String} [wechatMediaId] 微信提供的mediaId（和mp3Url二选一）
 *
 * @apiSuccess {Number} messageId 消息ID
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 0,
 *         "errmsg": "操作成功",
 *         "data": {
 *             "messageId": 123,
 *         }
 *     }
 */

// @router /:IMEI/voice [post]
func (this *DeviceController) PostVoice() {
	IMEI := this.Ctx.Input.Param(":IMEI")
	this.Data["json"] = restReturn(0, "发送聊天语音"+IMEI, map[string]interface{}{
		"messageId": 123,
	})
	this.ServeJSON()
}

/**
 * @api {get} /device/:IMEI/chatRecord 拉取聊天记录
 * @apiName DeivceChatRecord
 * @apiGroup Device
 *
 * @apiParam {String} username 用户名
 * @apiParam {String} ticket 用户接口调用凭据
 * @apiParam {String="timeDesc"} [orderType=timeDesc] 顺序类型，默认为timeDesc
 * @apiParam {Number} [startId=0] 开始拉取id
 * @apiParam {Number} [length=10] 拉取长度
 *
 * @apiSuccess {Number} id 消息id
 * @apiSuccess {Number} direction 语音方向，1为上行（设备->服务器），2为下行（服务器->设备）
 * @apiSuccess {String} type 消息类型，目前为voice
 * @apiSuccess {String} [voiceUrl] 语音url，消息类型为voice时返回
 * @apiSuccess {String} created_at 消息产生时间，格式为Y-m-d H:i:s
 * @apiSuccess {String} status 当前消息状态
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 0,
 *         "errmsg": "操作成功",
 *         "data": [{
 *             "id": 2,
 *             "direction": 2,
 *             "type": "voice",
 *             "url": "http://xxx.xxx.com/",
 *             "created_at": "2016-01-01 00:00:00",
 *             "status": "已发送到手表"
 *         },{
 *             "id": 1,
 *             "direction": 1,
 *             "type": "voice",
 *             "url": "http://xxx.xxx.com/",
 *             "created_at": "2016-01-01 00:00:00",
 *             "status": "已读"
 *         }]
 *     }
 */
// @router /:IMEI/chatRecord [get]
func (this *DeviceController) GetChatRecord() {
	var data []map[string]interface{}
	data = append(data, map[string]interface{}{
		"id":         3,
		"direction":  1,
		"type":       "voice",
		"url":        "http://xxx.xxx.com/",
		"created_at": "2016-01-01 00:00:02",
		"status":     "已读",
	})
	data = append(data, map[string]interface{}{
		"id":         2,
		"direction":  2,
		"type":       "voice",
		"url":        "http://xxx.xxx.com/",
		"created_at": "2016-01-01 00:00:01",
		"status":     "已发送到手表",
	})
	data = append(data, map[string]interface{}{
		"id":         1,
		"direction":  1,
		"type":       "voice",
		"url":        "http://xxx.xxx.com/",
		"created_at": "2016-01-01 00:00:00",
		"status":     "已读",
	})
	this.Data["json"] = restReturn(0, "操作成功", data)
	this.ServeJSON()
}
