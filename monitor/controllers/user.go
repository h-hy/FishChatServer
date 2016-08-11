package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
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
		if this.username == "" || ticket == "" {
			this.Data["json"] = restReturn(44001, "尚未登录或者登陆失效["+actionName+"]", map[string]interface{}{})
			this.ServeJSON()
		}
	}
}
func (this *UserController) initRongCloud() (*RCServerSDK.RCServer, error) {
	rongCloudAppKey := beego.AppConfig.String("rongCloudAppKey")
	rongCloudSecret := beego.AppConfig.String("rongCloudSecret")
	rcServer, rcError := RCServerSDK.NewRCServer(rongCloudAppKey, rongCloudSecret, "json")

	return rcServer, rcError

}

func (this *UserController) getRongCloudToken(userId, name, portraitUri string) (string, error) {
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
		this.Data["json"] = restReturn(20003, "用户名不能为空", map[string]interface{}{})
		this.ServeJSON()
	}
	if password == "" {
		this.Data["json"] = restReturn(20003, "密码不能为空", map[string]interface{}{})
		this.ServeJSON()
	}
	//	var user models.User
	//	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	//	uid := models.AddUser(user)
	//	this.Data["json"] = map[string]interface{}{"name": "astaxie"}
	this.Data["json"] = restReturn(0, "注册成功", map[string]interface{}{
		"username": "13590210000",
		"ticket":   "abcdefg",
	})
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
		this.Data["json"] = restReturn(20003, "用户名和密码不能为空", map[string]interface{}{})
		this.ServeJSON()
	}
	if SMScode == "" {
		this.Data["json"] = restReturn(20004, "短信验证码不能为空", map[string]interface{}{})
		this.ServeJSON()
	}
	if SMScode != "123456" {
		this.Data["json"] = restReturn(20013, "短信验证码错误，为123456", map[string]interface{}{})
		this.ServeJSON()
	}
	//	this.Data["json"] = map[string]interface{}{"name": "astaxie"}
	this.Data["json"] = restReturn(0, "密码重置成功", map[string]interface{}{
		"username": "13590210000",
		"ticket":   "abcdefg",
	})
	this.ServeJSON()
}

/**
 * @api {get} /user/:username 查看用户信息
 * @apiName userDetail
 * @apiGroup User
 *
 * @apiParam {String} ticket 用户接口调用凭据
 * @apiParam {boool=true,false} ticket 用户接口调用凭据
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
	rongCloudToken, _ := this.getRongCloudToken(this.username, "", "")
	this.Data["json"] = restReturn(0, "操作成功", map[string]interface{}{
		"username":        "13590210000",
		"rongCloudAppKey": rongCloudAppKey,
		"rongCloudToken":  rongCloudToken,
	})
	this.ServeJSON()
	//	uid := u.GetString(":uid")
	//	if uid != "" {
	//		user, err := models.GetUser(uid)
	//		if err != nil {
	//			u.Data["json"] = err.Error()
	//		} else {
	//			u.Data["json"] = user
	//		}
	//	}
	//	u.ServeJSON()
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
	this.Data["json"] = restReturn(0, "操作成功", map[string]interface{}{})
	this.ServeJSON()
	//	uid := u.GetString(":uid")
	//	if uid != "" {
	//		var user models.User
	//		json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	//		uu, err := models.UpdateUser(uid, &user)
	//		if err != nil {
	//			u.Data["json"] = err.Error()
	//		} else {
	//			u.Data["json"] = uu
	//		}
	//	}
	//	u.ServeJSON()
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
		this.Data["json"] = restReturn(20003, "用户名和密码不能为空", map[string]interface{}{})
		this.ServeJSON()
	}
	this.Data["json"] = restReturn(0, "操作成功", map[string]interface{}{
		"username": "13590210000",
		"ticket":   "abcdefg",
	})
	this.ServeJSON()
	//	username := u.GetString("username")
	//	password := u.GetString("password")
	//	if models.Login(username, password) {
	//		u.Data["json"] = "login success"
	//	} else {
	//		u.Data["json"] = map[string]string{
	//			"abc": "user not exist",
	//		}
	//	}
	//	u.ServeJSON()
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
		this.Data["json"] = restReturn(20003, "用户名不能为空", map[string]interface{}{})
		this.ServeJSON()
	}
	this.Data["json"] = restReturn(0, "获取成功，验证码是123456[测试]", map[string]interface{}{})
	this.ServeJSON()
}
