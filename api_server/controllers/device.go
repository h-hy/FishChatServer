package controllers

import (
	"errors"
	"strconv"
	"time"
	//	"fmt"
	//	"fmt"
	//	"encoding/json"
	"path/filepath"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/chanxuehong/wechat.v2/mp/media"
	"github.com/oikomi/FishChatServer/api_server/service"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/models"
	"github.com/oikomi/FishChatServer/protocol"
)

type DeviceController struct {
	beego.Controller
	username string
	userId   int
	m        *service.Monitor
}

func (this *DeviceController) Prepare() {
	this.Ctx.ResponseWriter.Header().Add("Access-Control-Allow-Origin", "*")
	ticket := this.GetString("ticket")
	this.username = this.GetString("username")
	user, err := models.GetUser(this.username)
	if err != nil || user.CheckTicket(ticket) == false {
		this.Data["json"] = restReturn(44001, "尚未登录或者登陆失效，请重新登陆", map[string]interface{}{})
		this.ServeJSON()
		return
	}
	this.userId = user.Id
	_, actionName := this.GetControllerAndAction()
	if actionName != "Post" && actionName != "Get" {
		IMEI := this.Ctx.Input.Param(":IMEI")
		if user.CheckBind(IMEI) == false {
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
func (this *DeviceController) sendToDevice(IMEI, cmdName string, Arg1 ...string) error {

	sessionCacheData, _ := getCache(this.m, IMEI)
	log.Info(sessionCacheData)
	if sessionCacheData != nil && sessionCacheData.MsgServerAddr != "" {
		log.Info("ok")
		cmd := protocol.NewCmdSimple(protocol.ACTION_TRANSFER_TO_DEVICE)
		cmd.Infos["cmdName"] = cmdName
		cmd.Infos["ConnectServerUUID"] = sessionCacheData.ConnectServerUUID
		if len(Arg1) > 0 {
			for _, arg := range Arg1 {
				cmd.AddArg(arg)
			}
		}
		cmd.Infos["IMEI"] = IMEI
		cmd.Infos["ConnectServerUUID"] = sessionCacheData.ConnectServerUUID
		for _, pushServerClient := range this.m.PushServerClientMap {
			err := pushServerClient.Session.Send(libnet.Json(cmd))
			if err != nil {
				log.Error(err.Error())
			}
			return nil
		}
	}
	return errors.New("NOT FOUND")
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
	models.NewDeviceCommand(this.userId, IMEI, "LOCATON", protocol.DEIVCE_LOCATON_CMD, "")
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
	models.NewDeviceCommand(this.userId, IMEI, "SHUTDOWN", protocol.DEIVCE_SHUTDOWN_CMD, "")
	this.sendToDevice(IMEI, "D"+protocol.DEIVCE_SHUTDOWN_CMD)
	this.Data["json"] = restReturn(0, "操作成功", map[string]interface{}{
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
 * @apiParam {String} wechatMediaId 微信提供的mediaId
 *
 * @apiSuccess {Number} id 消息id
 * @apiSuccess {Number} direction 语音方向，1为上行（设备->服务器），2为下行（服务器->设备）
 * @apiSuccess {String} type 消息类型，目前为voice
 * @apiSuccess {String} voiceUrl 语音url
 * @apiSuccess {String} created_at 消息产生时间，格式为Y-m-d H:i:s
 * @apiSuccess {String} status 当前消息状态
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 0,
 *         "errmsg": "操作成功",
 *         "data": {
 *             "id": 2,
 *             "direction": 2,
 *             "type": "voice",
 *             "voiceUrl": "http://xxx.xxx.com/",
 *             "created_at": "2016-01-01 00:00:00",
 *             "status": "发送中..."
 *         }
 *     }
 */

// @router /:IMEI/voice [post]
func (this *DeviceController) PostVoice() {
	IMEI := this.Ctx.Input.Param(":IMEI")
	wechatMediaId := this.GetString("wechatMediaId")
	if wechatMediaId == "" {
		this.Data["json"] = restReturn(50000, "wechatMediaId不能为空", map[string]interface{}{})
		this.ServeJSON()
		return
	}
	//开始保存文件
	amrSaveDir, _ := filepath.Abs(beego.AppConfig.String("amrSaveDir"))
	amrURIPrefix := beego.AppConfig.String("amrURIPrefix")
	timestr := time.Now().Format("2006_01_02_15_04_05")
	filename := IMEI + "_" + timestr + "_" + string(Krand(5, KC_RAND_KIND_LOWER)) + ".amr"
	log.Info(amrSaveDir + filename)
	written, err := media.Download(wechatClient, wechatMediaId, amrSaveDir+filename)
	if err != nil {
		log.Info(err)
		this.Data["json"] = restReturn(50000, "wechatMediaId错误", map[string]interface{}{})
		this.ServeJSON()
		return
	}
	//保存文件完毕，开始写入下行指令
	uri := amrURIPrefix + filename

	DMSId, cmdId, err := models.NewDeviceCommand(this.userId, IMEI, "VOICE_DOWN", protocol.DEIVCE_VOICE_DOWN_CMD, uri)
	if err != nil {
		log.Info(err)
		this.Data["json"] = restReturn(50000, "保存失败", map[string]interface{}{})
		this.ServeJSON()
		return
	}
	//写入下行指令完毕，开始检查客户端是否在线
	voice := service.Voice{Id: cmdId, Uri: uri, Filename: filename, PathFilename: amrSaveDir + filename, Size: int(written)}
	voice.Cache()

	this.sendToDevice(IMEI, "D"+protocol.DEIVCE_SHUTDOWN_CMD, strconv.Itoa(cmdId), "0", "amr", strconv.Itoa(int(written)))
	//全部操作完成
	this.Data["json"] = restReturn(0, "操作成功", map[string]interface{}{
		"id":         DMSId,
		"direction":  2,
		"type":       "voice",
		"created_at": time.Now().Format("2006-01-02 15:04:05"),
		"status":     "发送中...",
		"voiceUrl":   uri,
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
 *             "voiceUrl": "http://xxx.xxx.com/",
 *             "created_at": "2016-01-01 00:00:00",
 *             "status": "已发送到手表"
 *         },{
 *             "id": 1,
 *             "direction": 1,
 *             "type": "voice",
 *             "voiceUrl": "http://xxx.xxx.com/",
 *             "created_at": "2016-01-01 00:00:00",
 *             "status": "已读"
 *         }]
 *     }
 */
// @router /:IMEI/chatRecord [get]
func (this *DeviceController) GetChatRecord() {
	IMEI := this.Ctx.Input.Param(":IMEI")
	startId := this.GetString("startId")
	orderType := this.GetString("orderType")
	if orderType != "" && orderType != "timeDesc" {
		this.Data["json"] = restReturn(50000, "orderType不支持", map[string]interface{}{})
		this.ServeJSON()
		return
	}
	i_startId, err := strconv.Atoi(startId)
	if err != nil {
		i_startId = 0
	}
	length := this.GetString("length")
	i_length, err := strconv.Atoi(length)
	if err != nil {
		i_length = 10
	}
	//参数读取完毕，开始加载数据
	var data []map[string]interface{}
	o := orm.NewOrm()
	var charRecord []models.DeviceMessageCenter
	//	charRecord := models.DeviceMessageCenter{UserId: this.userId}

	q := o.QueryTable(&models.DeviceMessageCenter{}).Filter("UserId", this.userId).Filter("IMEI", IMEI).Filter("action", "voice").OrderBy("-id").Limit(i_length)
	if i_startId > 0 {
		q = q.Filter("id__lt", i_startId)
	}
	q.All(&charRecord)
	//数据加载完毕，开始构建回复
	for _, record := range charRecord {
		data = append(data, map[string]interface{}{
			"id":         record.Id,
			"direction":  record.Direction,
			"type":       "voice",
			"voiceUrl":   record.VoiceUri,
			"created_at": record.CreatedAt,
			"status":     record.Content,
		})
	}
	this.Data["json"] = restReturn(0, "操作成功", data)
	this.ServeJSON()
}
