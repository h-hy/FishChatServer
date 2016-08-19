package controllers

import (
	"strconv"
	//	"fmt"
	//	"fmt"
	//	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/monitor/models"
	"github.com/oikomi/FishChatServer/monitor/service"
	"github.com/oikomi/FishChatServer/protocol"
)

type DeviceController struct {
	beego.Controller
	username string
	userId   string
	m        *service.Monitor
}

func (this *DeviceController) Prepare() {
	this.Ctx.ResponseWriter.Header().Add("Access-Control-Allow-Origin", "*")
	ticket := this.GetString("ticket")
	this.username = this.GetString("username")
	log.Info(string(this.Ctx.Input.RequestBody))
	log.Info("ticket=", ticket)
	log.Info("this.username=", this.username)
	if models.UserCheckTicket(this.username, ticket) == false {
		this.Data["json"] = restReturn(44001, "尚未登录或者登陆失效，请重新登陆", map[string]interface{}{})
		this.ServeJSON()
		return
	}
	_, actionName := this.GetControllerAndAction()
	if actionName != "Post" && actionName != "Get" {
		IMEI := this.Ctx.Input.Param(":IMEI")
		if models.CheckBind(this.username, IMEI) == false {
			this.Data["json"] = restReturn(ERROR_DEIVCE_NOT_BIND, ERROR_DEIVCE_NOT_BIND_STRING, map[string]interface{}{})
			this.ServeJSON()
			return
		}
	}
	this.m = service.GetServer()
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
	user, err := models.GetUser(this.username)
	if err != nil {
		this.Data["json"] = restReturn(50000, "获取用户信息失败，请联系管理员", map[string]interface{}{})
		this.ServeJSON()
		return
	}
	o := orm.NewOrm()
	_, err = o.LoadRelated(&user, "Devices")
	if err != nil {
		this.Data["json"] = restReturn(50000, "获取用户绑定设备失败，请联系管理员", map[string]interface{}{})
		this.ServeJSON()
		return
	}

	for _, device := range user.Devices {
		data = append(data, map[string]interface{}{
			"IMEI":           device.IMEI,
			"nick":           device.Nick,
			"status":         device.Alive,
			"work_model":     device.Work_model,
			"volume":         device.Volume,
			"electricity":    device.Energy,
			"emeregncyPhone": device.EmergencyPhone,
		})
	}
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
	o := orm.NewOrm()
	device := models.Device{IMEI: IMEI}

	err := o.Read(&device, "IMEI")
	if err == orm.ErrNoRows {
		this.Data["json"] = restReturn(ERROR_DEIVCE_NOT_EXIST, ERROR_DEIVCE_NOT_EXIST_STRING, map[string]interface{}{})
		this.ServeJSON()
		return
	}

	this.Data["json"] = restReturn(0, "操作成功", map[string]interface{}{
		"IMEI":           device.IMEI,
		"nick":           device.Nick,
		"status":         device.Alive,
		"work_model":     device.Work_model,
		"volume":         device.Volume,
		"electricity":    device.Energy,
		"emeregncyPhone": device.EmergencyPhone,
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

	IMEI := this.GetString("IMEI")
	//检查是否已经绑定
	hadBind := models.CheckBind(this.username, IMEI)
	log.Info("hadBind=", hadBind)
	if hadBind == true {
		this.Data["json"] = restReturn(ERROR_DEIVCE_BINDED, ERROR_DEIVCE_BINDED_STRING, map[string]interface{}{})
		this.ServeJSON()
		return
	}

	//获取设备对象
	o := orm.NewOrm()
	device := models.Device{IMEI: IMEI}

	err := o.Read(&device, "IMEI")
	if err == orm.ErrNoRows {
		this.Data["json"] = restReturn(ERROR_DEIVCE_NOT_EXIST, ERROR_DEIVCE_NOT_EXIST_STRING, map[string]interface{}{})
		this.ServeJSON()
		return
	} else if err != nil {
		this.Data["json"] = restReturn(50000, "获取设备信息失败，请联系管理员", map[string]interface{}{})
		this.ServeJSON()
		return
	}
	//获取用户对象

	user, err := models.GetUser(this.username)
	if err != nil {
		this.Data["json"] = restReturn(50000, "获取用户信息失败，请联系管理员", map[string]interface{}{})
		this.ServeJSON()
		return
	}
	//开始绑定
	m2m := o.QueryM2M(&user, "Devices")

	_, err = m2m.Add(device)
	if err == nil {
		user.Devices = append(user.Devices, &device)
		user.CacheUser()
		this.Data["json"] = restReturn(0, "操作成功", map[string]interface{}{
			"IMEI":           device.IMEI,
			"nick":           device.Nick,
			"status":         device.Alive,
			"work_model":     device.Work_model,
			"volume":         device.Volume,
			"electricity":    device.Energy,
			"emeregncyPhone": device.EmergencyPhone,
		})
		this.ServeJSON()
	}
	this.Data["json"] = restReturn(50000, "绑定失败，请联系管理员", map[string]interface{}{})
	this.ServeJSON()
}

/**
 * @api {delete} /device/:IMEI?username=:user&ticket=:ticket 用户删除绑定设备
 * @apiDescription 特别说明：根据HTTP标准，DELETE方法的身份认证参数务必放在url中而不能放在body中。
 * @apiName deviceDestory
 * @apiGroup Device
 *
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

	user, err := models.GetUser(this.username)
	if err != nil {
		this.Data["json"] = restReturn(50000, "获取用户信息失败，请联系管理员", map[string]interface{}{})
		this.ServeJSON()
		return
	}

	for index, device := range user.Devices {
		if device.IMEI == IMEI {
			newDevices := make([]*models.Device, len(user.Devices)-1)
			copy(newDevices[0:], user.Devices[:index])
			copy(newDevices[index:], user.Devices[index+1:])
			user.Devices = newDevices
			o := orm.NewOrm()
			m2m := o.QueryM2M(&user, "Devices")
			m2m.Remove(device)
			user.CacheUser()
			this.Data["json"] = restReturn(0, "操作成功", map[string]interface{}{})
			this.ServeJSON()
			return
		}
	}
	this.Data["json"] = restReturn(50000, "操作失败，请联系管理员", map[string]interface{}{})
	this.ServeJSON()
}
func (this *DeviceController) sendToDevice(IMEI, cmdName string, Arg1 ...string) {

	sessionCacheData, _ := getCache(this.m, IMEI)
	log.Info(sessionCacheData)
	if sessionCacheData != nil {
		if sessionCacheData.MsgServerAddr != "" &&
			this.m.MsgServerClientMap[sessionCacheData.MsgServerAddr] != nil {
			log.Info("ok")
			cmd := protocol.NewCmdSimple(protocol.ACTION_TRANSFER_TO_DEVICE)
			cmd.Infos["cmdName"] = cmdName
			if len(Arg1) > 0 {
				cmd.AddArg(Arg1[0])

			}
			cmd.Infos["IMEI"] = IMEI
			cmd.Infos["ConnectServerUUID"] = sessionCacheData.ConnectServerUUID
			err := this.m.MsgServerClientMap[sessionCacheData.MsgServerAddr].Session.Send(libnet.Json(cmd))
			if err != nil {
				log.Error(err.Error())
			}
		}
	}
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
	volume := this.GetString("volume")
	nick := this.GetString("nick")
	work_model := this.GetString("work_model")
	//	emeregncyPhone := this.GetString("emeregncyPhone")
	var updateColumn []string
	o := orm.NewOrm()
	device := models.Device{IMEI: IMEI}

	if nick != "" {
		updateColumn = append(updateColumn, "Nick")
		device.Nick = nick
	}
	if volume != "" {
		this.sendToDevice(IMEI, "D"+protocol.DEIVCE_VOLUME_LEVER, volume)
		updateColumn = append(updateColumn, "Volume")
		device.Volume, _ = strconv.Atoi(volume)
	}
	if work_model != "" {
		this.sendToDevice(IMEI, "D"+protocol.DEIVCE_WORK_MODEL_CMD, work_model)
		updateColumn = append(updateColumn, "Work_model")
		device.Work_model, _ = strconv.Atoi(work_model)
	}
	if len(updateColumn) > 0 {
		o.Update(&device, updateColumn...)
	}
	this.Data["json"] = restReturn(0, "操作成功", map[string]interface{}{})
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
	this.sendToDevice(IMEI, "D"+protocol.DEIVCE_LOCATON_CMD)
	this.Data["json"] = restReturn(0, "操作成功", map[string]interface{}{
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
	this.sendToDevice(IMEI, "D"+protocol.DEIVCE_SHUTDOWN_CMD)
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
