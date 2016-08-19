package controllers

import (
	"encoding/json"
	"fmt"

	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/monitor/models"
	"github.com/rongcloud/server-sdk-go/RCServerSDK"
)

type UserController struct {
	beego.Controller
	username string
	userId   string
}

func (this *UserController) Prepare() {
	this.Ctx.ResponseWriter.Header().Add("Access-Control-Allow-Origin", "*")
	_, actionName := this.GetControllerAndAction()
	this.username = this.Ctx.Input.Param(":username")
	if actionName != "Login" && actionName != "Post" && actionName != "ResetPassword" && actionName != "GetSMSCode" {
		ticket := this.GetString("ticket")
		if models.UserCheckTicket(this.username, ticket) == false {
			this.Data["json"] = restReturn(ERROR_USER_TICKET, ERROR_USER_TICKET_STRING, map[string]interface{}{})
			this.ServeJSON()
			return
		}
	}
}
func (this *UserController) initRongCloud() (*RCServerSDK.RCServer, error) {
	rongCloudAppKey := beego.AppConfig.String("rongCloudAppKey")
	rongCloudSecret := beego.AppConfig.String("rongCloudSecret")
	rcServer, rcError := RCServerSDK.NewRCServer(rongCloudAppKey, rongCloudSecret, "json")

	return rcServer, rcError

}

func (this *UserController) getRongCloudToken(userId, name, portraitUri string, getNew bool) (string, error) {
	if getNew == false {
		rongCloudToken := redisCache.Get("rongCloudToken_" + userId)
		if rongCloudToken != nil {
			rongCloudTokenString := GetString(rongCloudToken)
			return rongCloudTokenString, nil
		}
	}

	rcServer, rcError := this.initRongCloud()
	if rcError != nil {
		return "", rcError
	}
	byteData, rcError := rcServer.UserGetToken(userId, name, portraitUri)
	if rcError != nil {
		return "", rcError
	}
	var UserGetToken struct {
		Code  int    `json:"code"`
		Token string `json:"token"`
	}
	json.Unmarshal(byteData, &UserGetToken)
	redisCache.Put("rongCloudToken_"+userId, UserGetToken.Token, 30*time.Minute)
	return UserGetToken.Token, nil
}

/**
 * @api {post} /user/ 用户注册
 * @apiName userStore
 * @apiGroup User
 *
 * @apiParam {String} username 用户名
 * @apiParam {String} password 用户密码
 * @apiParam {String} [code] 微信用户授权码
 *
 * @apiParamExample {String} Request-Example:
 * username=13590210000&password=123456
 *
 * @apiSuccess {String} username 用户名
 * @apiSuccess {String} ticket 用户接口调用凭据
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 0,
 *         "errmsg": "注册成功",
 *         "data": {
 *             "username": "13590210000",
 *             "ticket": "abcdefg"
 *         }
 *     }
 *
 * @apiErrorExample 用户名已经存在回复
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 20002,
 *         "errmsg": "用户名已经存在",
 *         "data": {
 *         }
 *     }
 *
 * @apiErrorExample 微信用户授权码已失效回复
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 40001,
 *         "errmsg": "微信用户授权码已失效",
 *         "data": {
 *         }
 *     }
 *     errcode=40001的情况请重新发起微信页面授权
 *
 */

// @router / [post]
func (this *UserController) Post() {

	username := this.GetString("username")
	password := this.GetString("password")
	if username == "" {
		this.Data["json"] = restReturn(ERROR_USER_NAME_NEED, ERROR_USER_NAME_NEED_STRING, map[string]interface{}{})
		this.ServeJSON()
		return
	}
	if password == "" {
		this.Data["json"] = restReturn(ERROR_USER_NAME_PASSWORD_NEED, ERROR_USER_NAME_PASSWORD_NEED_STRING, map[string]interface{}{})
		this.ServeJSON()
		return
	}
	code := this.GetString("code")
	var openid string
	if code != "" {
		openidCache := redisCache.Get("WechatAuthCode_" + code)
		log.Info("WechatAuthCode_" + code)
		if openidCache == nil {
			this.Data["json"] = restReturn(ERROR_WECHAT_CODE_ERROR, ERROR_WECHAT_CODE_ERROR_STRING, map[string]interface{}{})
			this.ServeJSON()
			return
		} else {
			openid = GetString(openidCache)
			models.UserCleanOpenid(openid)
		}
	}
	errcode, errmsg, data := models.UserRegist(username, username, password, openid)
	if errcode == 0 {
		redisCache.Delete("WechatAuthCode_" + code)
	}
	this.Data["json"] = restReturn(errcode, errmsg, data)
	this.ServeJSON()
}

/**
 * @api {post} /user/:username/resetPassword 用户找回密码
 * @apiName userResetPassword
 * @apiGroup User
 *
 * @apiParam {String} password 用户新密码
 * @apiParam {String} SMScode 获取到的短信验证码
 *
 * @apiParamExample {String} Request-Example:
 * username=13590210000&password=123456&SMScode=123456
 *
 * @apiSuccess {String} username 用户名
 * @apiSuccess {String} ticket 用户接口调用凭据
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 0,
 *         "errmsg": "密码重置成功",
 *         "data": {
 *             "username": "13590210000",
 *             "ticket": "abcdefg"
 *         }
 *     }
 *
 * @apiErrorExample 用户名不存在回复
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 20002,
 *         "errmsg": "用户名不存在",
 *         "data": {
 *         }
 *     }
 *
 * @apiErrorExample 短信验证码错误回复
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 20013,
 *         "errmsg": "短信验证码错误",
 *         "data": {
 *         }
 *     }
 *
 */

// @router /:username/resetPassword [post]
func (this *UserController) ResetPassword() {
	username := this.username
	password := this.GetString("password")
	SMScode := this.GetString("SMScode")
	if username == "" || password == "" {
		this.Data["json"] = restReturn(ERROR_USER_NAME_PASSWORD_NEED, ERROR_USER_NAME_PASSWORD_NEED_STRING, map[string]interface{}{})
		this.ServeJSON()
		return
	}
	if SMScode == "" {
		this.Data["json"] = restReturn(ERROR_USER_SMS_CODE_NEED, ERROR_USER_SMS_CODE_NEED_STRING, map[string]interface{}{})
		this.ServeJSON()
		return
	}
	//开始验证码校验
	SMSCodeCache := redisCache.Get("SMSCode_" + username)
	log.Info("SMSCode_" + username)
	if SMSCodeCache == nil {
		this.Data["json"] = restReturn(ERROR_USER_SMS_CODE_NOT_FOUND, ERROR_USER_SMS_CODE_NOT_FOUND_STRING, map[string]interface{}{})
		this.ServeJSON()
		return
	}
	SMSCodeCacheString := GetString(SMSCodeCache)
	if SMSCodeCacheString != SMScode || SMSCodeCacheString == "" {

		this.Data["json"] = restReturn(ERROR_USER_SMS_CODE_ERROR, ERROR_USER_SMS_CODE_ERROR_STRING, map[string]interface{}{})
		this.ServeJSON()
		return
	}
	//开始用户校验

	o := orm.NewOrm()
	user, err := models.GetUser(username)

	log.Info(err)
	if err == orm.ErrNoRows {
		this.Data["json"] = restReturn(ERROR_USER_NOT_FOUND, ERROR_USER_NOT_FOUND_STRING, map[string]interface{}{})
		this.ServeJSON()
		return
	} else if err == nil {

		//开始写入密码
		user.Ticket = GetNewTicket()
		user.CacheUser()
		//		models.UserCacheTicket(username, user.Ticket)
		user.Password = password
		o.Update(&user, "Ticket", "Password")
		redisCache.Delete("SMSCode_" + username)
		this.Data["json"] = restReturn(0, "密码重置成功", map[string]interface{}{
			"username": username,
			"ticket":   user.Ticket,
		})
		this.ServeJSON()
		return
	}
	//	this.ServeJSON()
}

/**
 * @api {get} /user/:username/updateRongCloudToken 刷新融云密钥
 * @apiName userUpdateRongCloudToken
 * @apiGroup User
 *
 * @apiParam {String} ticket 用户接口调用凭据
 *
 * @apiSuccess {String} rongCloudAppKey 融云AppKey
 * @apiSuccess {String} rongCloudToken 融云token
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 0,
 *         "errmsg": "操作成功",
 *         "data": {
 *             "rongCloudAppKey": "rongCloudAppKey"
 *             "rongCloudToken": "rongCloudToken"
 *         }
 *     }
 */

// @router /:username/updateRongCloudToken [get]
func (this *UserController) UpdateRongCloudToken() {
	rongCloudAppKey := beego.AppConfig.String("rongCloudAppKey")
	rongCloudToken, _ := this.getRongCloudToken(this.username, "", "", true)
	this.Data["json"] = restReturn(0, "操作成功", map[string]interface{}{
		"rongCloudAppKey": rongCloudAppKey,
		"rongCloudToken":  rongCloudToken,
	})
	this.ServeJSON()
}

/**
 * @api {get} /user/:username 查看用户信息
 * @apiDescription 接口中的融云token如果失效，需要调用“刷新融云密钥”接口来刷新
 * @apiName userDetail
 * @apiGroup User
 *
 * @apiParam {String} ticket 用户接口调用凭据
 *
 * @apiSuccess {String} username 用户名
 * @apiSuccess {String} rongCloudAppKey 融云AppKey
 * @apiSuccess {String} rongCloudToken 融云token
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 0,
 *         "errmsg": "操作成功",
 *         "data": {
 *             "username": "13590210000"
 *             "rongCloudAppKey": "rongCloudAppKey"
 *             "rongCloudToken": "rongCloudToken"
 *         }
 *     }
 */

// @router /:username [get]
func (this *UserController) Get() {
	rcServer, rcError := this.initRongCloud()
	if rcError == nil {
		data := map[string]interface{}{
			"content": "MessageReceived",
			"extra": map[string]interface{}{
				"messageId":    1234,
				"fromUsername": "123",
				"toUsername":   "123",
				"type":         "voice",
				"mp3Url":       "http://baidu.com/a.mp3",
				"created_at":   "2016-08-08 11:11:11",
			},
		}
		dataString, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
		} else {
			if returnData, returnError := rcServer.MessagePrivatePublish("system", []string{this.username}, "RC:TxtMsg", string(dataString), "", ""); returnError != nil || len(returnData) == 0 {
				fmt.Print("发送单聊消息：测试失败！！！")
			} else {
				fmt.Print("发送单聊消息：测试通过。returnData:", string(returnData))
			}
		}

		dataLocation := map[string]interface{}{
			"content": "LocationUpdated",
			"extra": map[string]interface{}{
				"messageId":    123,
				"IMEI":         "123",
				"nick":         "456",
				"toUsername":   "voice",
				"locationType": "GPS",
				"mapType":      "amap",
				"lat":          "22",
				"lng":          "11",
				"radius":       "11",
				"created_at":   "2016-08-08 11:11:11",
			},
		}
		dataLocationString, err := json.Marshal(dataLocation)
		if err != nil {
			fmt.Println(err)
		} else {
			if returnData, returnError := rcServer.MessagePrivatePublish("system", []string{this.username}, "RC:TxtMsg", string(dataLocationString), "", ""); returnError != nil || len(returnData) == 0 {
				fmt.Print("发送单聊消息：测试失败！！！")
			} else {
				fmt.Print("发送单聊消息：测试通过。returnData:", string(returnData))
			}
		}
	}

	rongCloudAppKey := beego.AppConfig.String("rongCloudAppKey")
	rongCloudToken, _ := this.getRongCloudToken(this.username, "", "", false)
	this.Data["json"] = restReturn(0, "操作成功", map[string]interface{}{
		"username":        this.username,
		"rongCloudAppKey": rongCloudAppKey,
		"rongCloudToken":  rongCloudToken,
	})
	this.ServeJSON()
}

/**
 * @api {put} /user/:username 更新用户信息
 * @apiDescription 本接口只需要传入需要更新的参数即可
 * @apiName userUpdate
 * @apiGroup User
 *
 * @apiParam {String} ticket 用户接口调用凭据
 * @apiParam {String} [oldPassword] 用户旧密码
 * @apiParam {String} [newPassword] 用户新密码
 *
 * @apiParamExample {String} Request-Example:
 * oldPassword=123456&newPassword=111111
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 0,
 *         "errmsg": "操作成功",
 *         "data": {
 *         }
 *     }
 *
 * @apiErrorExample Error-Response:
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 20003,
 *         "errmsg": "原密码正确",
 *         "data": {
 *         }
 *     }
 */

// @router /:username [put]
func (this *UserController) Put() {

	user, err := models.GetUser(this.username)
	log.Info(err)
	if err == orm.ErrNoRows {
		this.Data["json"] = restReturn(ERROR_USER_NOT_FOUND, ERROR_USER_NOT_FOUND_STRING, map[string]interface{}{})
		this.ServeJSON()
		return
	}

	newPassword := this.GetString("newPassword")
	oldPassword := this.GetString("oldPassword")
	if newPassword != "" {

		if user.Password != oldPassword && false {
			this.Data["json"] = restReturn(ERROR_USER_PASSWORD_ERROR, ERROR_USER_PASSWORD_ERROR_STRING, map[string]interface{}{})
			this.ServeJSON()
			return
		} else {
			o := orm.NewOrm()
			user.Password = newPassword
			o.Update(&user, "Password")
			user.CacheUser()
		}
	}

	this.Data["json"] = restReturn(0, "操作成功", map[string]interface{}{})
	this.ServeJSON()
}

/**
 * @api {post} /user/:username/login 用户登录
 * @apiName userLogin
 * @apiGroup User
 *
 * @apiParam {String} password 用户密码
 * @apiParam {String} [code] 微信用户授权码
 *
 * @apiParamExample {String} Request-Example:
 * password=111111&code=123
 *
 * @apiSuccess {String} username 用户名
 * @apiSuccess {String} ticket 用户接口调用凭据
 *
 * @apiSuccessExample 正常回复
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 0,
 *         "errmsg": "操作成功",
 *         "data": {
 *             "username": "13590210000",
 *             "ticket": "abcdefg"
 *         }
 *     }
 *
 * @apiErrorExample 用户名不存在回复
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 20003,
 *         "errmsg": "用户名不存在",
 *         "data": {
 *         }
 *     }
 *
 * @apiErrorExample 用户密码错误回复
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 20004,
 *         "errmsg": "用户密码错误",
 *         "data": {
 *         }
 *     }
 *
 * @apiErrorExample 微信用户授权码已失效回复
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 40001,
 *         "errmsg": "微信用户授权码已失效",
 *         "data": {
 *         }
 *     }
 *     errcode=40001的情况请重新发起微信页面授权
 */
// @router /:username/login [post]
func (this *UserController) Login() {
	username := this.username
	password := this.GetString("password")
	if username == "" || password == "" {
		this.Data["json"] = restReturn(ERROR_USER_NAME_PASSWORD_NEED, ERROR_USER_NAME_PASSWORD_NEED_STRING, map[string]interface{}{})
		this.ServeJSON()
		return
	}
	user, err := models.GetUser(username)

	log.Info(err)
	if err == orm.ErrNoRows {
		this.Data["json"] = restReturn(ERROR_USER_NOT_FOUND, ERROR_USER_NOT_FOUND_STRING, map[string]interface{}{})
		this.ServeJSON()
		return
	} else if err == nil {
		if user.Password != password {
			this.Data["json"] = restReturn(ERROR_USER_PASSWORD_ERROR, ERROR_USER_PASSWORD_ERROR_STRING, map[string]interface{}{})
			this.ServeJSON()
			return
		}
		code := this.GetString("code")
		var openid string = ""
		if code != "" {
			openidCache := redisCache.Get("WechatAuthCode_" + code)
			log.Info("WechatAuthCode_" + code)
			if openidCache == nil {
				this.Data["json"] = restReturn(ERROR_WECHAT_CODE_ERROR, ERROR_WECHAT_CODE_ERROR_STRING, map[string]interface{}{})
				this.ServeJSON()
				return
			} else {
				openid = GetString(openidCache)
				models.UserCleanOpenid(openid)
				user.Openid = openid
			}
			redisCache.Delete("WechatAuthCode_" + code)
		}
		user.Ticket = GetNewTicket()

		o := orm.NewOrm()
		if openid == "" {
			o.Update(&user, "Ticket")
		} else {
			o.Update(&user, "Ticket", "Openid")
		}
		user.CacheUser()
		//		models.UserCacheTicket(username, user.Ticket)
		this.Data["json"] = restReturn(0, "登陆成功", map[string]interface{}{
			"username": username,
			"ticket":   user.Ticket,
		})
		this.ServeJSON()
		return
	}
	log.Info(err)
	this.Data["json"] = restReturn(20009, "登陆失败，请与管理员联系1", map[string]interface{}{})
	this.ServeJSON()
}

/**
 * @api {post} /user/:username/logout 用户退出登录
 * @apiName userLogout
 * @apiGroup User
 *
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
// @router /:username/logout [post]
func (this *UserController) Logout() {

	o := orm.NewOrm()
	user := models.User{Username: this.username}
	err := o.Read(&user, "Username")
	log.Info(err)
	if err == orm.ErrNoRows {
		this.Data["json"] = restReturn(ERROR_USER_NOT_FOUND, ERROR_USER_NOT_FOUND_STRING, map[string]interface{}{})
		this.ServeJSON()
		return
	}

	user.Ticket = ""
	user.Openid = ""
	o.Update(&user, "Ticket", "Openid")
	redisCache.Delete("user_" + this.username)
	//	user.CacheUser()
	//	models.UserCacheTicket(this.username, "")

	this.Data["json"] = restReturn(0, "操作成功", map[string]interface{}{})
	this.ServeJSON()
}

/**
 * @api {get} /user/:username/SMSCode 获取短信验证码
 * @apiName userSMSCode
 * @apiGroup User
 *
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 0,
 *         "errmsg": "获取成功",
 *         "data": {
 *         }
 *     }
 *
 * @apiErrorExample 秒级频率限制提示
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 20010,
 *         "errmsg": "短信验证码已发送，请60秒后再试",
 *         "data": {
 *         }
 *     }
 *
 * @apiErrorExample 天级频率限制提示
 *     HTTP/1.1 200 OK
 *     {
 *         "errcode": 20011,
 *         "errmsg": "每天最多发送10条验证码短信，请明天再试",
 *         "data": {
 *         }
 *     }
 *
 */

// @router /:username/SMSCode [get]
func (this *UserController) GetSMSCode() {
	username := this.username
	if username == "" {
		this.Data["json"] = restReturn(ERROR_USER_NAME_NEED, ERROR_USER_NAME_NEED_STRING, map[string]interface{}{})
		this.ServeJSON()
		return
	}

	o := orm.NewOrm()
	user := models.User{Username: username}

	err := o.Read(&user, "Username")
	log.Info(err)
	if err == orm.ErrNoRows {
		this.Data["json"] = restReturn(ERROR_USER_NOT_FOUND, ERROR_USER_NOT_FOUND_STRING, map[string]interface{}{})
		this.ServeJSON()
		return
	}

	code := string(Krand(6, KC_RAND_KIND_NUM))

	redisCache.Put("SMSCode_"+username, code, 30*time.Minute)

	this.Data["json"] = restReturn(0, "获取成功，验证码是"+code, map[string]interface{}{})
	this.ServeJSON()
}
