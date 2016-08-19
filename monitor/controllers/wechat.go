package controllers

import (
	"fmt"
	"log"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/chanxuehong/rand"
	"github.com/chanxuehong/wechat.v2/mp/core"
	mpoauth2 "github.com/chanxuehong/wechat.v2/mp/oauth2"
	"github.com/chanxuehong/wechat.v2/oauth2"
	"github.com/oikomi/FishChatServer/monitor/models"
)

type WechatController struct {
	beego.Controller
}

var (
	wxAppId           = beego.AppConfig.String("wechatAppId")
	wxAppSecret       = beego.AppConfig.String("wechatAppSecret")
	wxOriId           = beego.AppConfig.String("wechatOriId")
	wxToken           = beego.AppConfig.String("wechatToken")
	wxEncodedAESKey   = beego.AppConfig.String("wechatEncodedAESKey")
	oauth2RedirectURI = beego.AppConfig.String("wechatOauth2RedirectURI")
)

var (
	accessTokenServer core.AccessTokenServer = core.NewDefaultAccessTokenServer(wxAppId, wxAppSecret, nil)
	wechatClient      *core.Client           = core.NewClient(accessTokenServer, nil)
)

/**
 * @api {get} /wechat/authorize 微信请求授权
 * @apiName wechatAuthorize
 * @apiGroup Wechat
 * @apiSampleRequest off
 * @apiParam {String} redirectUri 授权后跳转地址（返回参数格式待商议）
 *
 * @apiSuccess {String} code 用户授权码（注册和登陆时需带上）
 * @apiSuccess {bool} isLogined 用户是否已经登陆（true/false）
 * @apiSuccess {String} [username] 用户名（isLogined==true带上）
 * @apiSuccess {String} [ticket] 用户接口调用凭据（isLogined==true带上）
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 302 Found
 *     Redirect to {redirectUri}&code={code}&isLogined={isLogined}&username={username}&ticket={ticket}
 *
 * @apiErrorExample Error-Response:
 *     HTTP/1.1 403 Forbidden
 */

// @router /authorize [get]
func (this *WechatController) Authorize() {
	redirectUri := this.GetString("redirectUri")
	if redirectUri == "" {
		this.Abort("403")
	}
	oauth2Scope := "snsapi_userinfo"
	state := string(rand.NewHex())
	this.SetSession("state", state)
	this.SetSession("redirectUri", redirectUri)
	AuthCodeURL := mpoauth2.AuthCodeURL(wxAppId, oauth2RedirectURI, oauth2Scope, state)
	fmt.Println(AuthCodeURL)
	this.Ctx.Redirect(302, AuthCodeURL)
}

// AuthorizeRedirect方法为微信回调
// @router /authorizeRedirect [get]
func (this *WechatController) AuthorizeRedirect() {
	//判断请求是否合法（state是否和session中一致）
	sessionState := this.GetSession("state")
	getState := this.GetString("state")
	if sessionState == nil || getState == "" || sessionState != getState {
		this.Abort("403")
	}
	//判断请求redirectUri是否存在
	redirectUri := this.GetSession("redirectUri")
	if redirectUri == nil {
		this.Abort("403")
	}
	log.Printf("redirectUri: %+v\r\n", redirectUri)
	//从微信获取AccessToken
	var oauth2Endpoint oauth2.Endpoint = mpoauth2.NewEndpoint(wxAppId, wxAppSecret)
	oauth2Client := oauth2.Client{
		Endpoint: oauth2Endpoint,
	}
	code := this.GetString("code")
	token, err := oauth2Client.ExchangeToken(code)
	if err != nil {
		this.Abort(err.Error())
	}
	log.Printf("token: %+v\r\n", token)

	//从微信拉去授权信息
	//	userinfo, err := mpoauth2.GetUserInfo(token.AccessToken, token.OpenId, "", nil)

	//	log.Printf("userinfo:", userinfo)

	wechatAuthorizeCode := string(Krand(32, KC_RAND_KIND_LOWER))
	err = redisCache.Put("WechatAuthCode_"+wechatAuthorizeCode, token.OpenId, 30*time.Minute)

	log.Println(err)
	o := orm.NewOrm()

	user := models.User{Openid: token.OpenId}

	err = o.Read(&user, "Openid")
	if err == orm.ErrNoRows {
		appRedirectUri := redirectUri.(string) + "&code=" + wechatAuthorizeCode +
			"&isLogined=false"
		this.Ctx.Redirect(302, appRedirectUri)
	} else if err == nil {
		user.Ticket = GetNewTicket()
		o.Update(&user, "Ticket")
		appRedirectUri := redirectUri.(string) + "&code=" + wechatAuthorizeCode +
			"&username=" + user.Username +
			"&isLogined=true" +
			"&ticket=" + user.Ticket
		this.Ctx.Redirect(302, appRedirectUri)
	} else {
		this.Data["json"] = restReturn(44002, "用户授权失败，请与管理员联系", map[string]interface{}{})
		this.ServeJSON()
	}

}
