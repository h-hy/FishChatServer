package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/oikomi/FishChatServer/monitor/controllers"
)

func init() {
	orm.Debug = true
	orm.RegisterDriver("mysql", orm.DRMySQL)

	maxIdle := 30
	maxConn := 30

	mysqlUsername := beego.AppConfig.String("mysqlUsername")
	mysqlPassword := beego.AppConfig.String("mysqlPassword")
	mysqlDatabase := beego.AppConfig.String("mysqlDatabase")

	orm.RegisterDataBase("default", "mysql", mysqlUsername+":"+mysqlPassword+"@/"+mysqlDatabase+"?charset=utf8", maxIdle, maxConn)

	beego.ErrorController(&controllers.ErrorController{})
	beego.Options("*", func(ctx *context.Context) {
		ctx.ResponseWriter.Header().Add("Access-Control-Allow-Origin", "*")
		ctx.ResponseWriter.Header().Add("Access-Control-Allow-Headers", "X-Requested-With, Content-Type")
		ctx.ResponseWriter.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	})

	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/device",
			beego.NSInclude(
				&controllers.DeviceController{},
			),
		),
		beego.NSNamespace("/wechat",
			beego.NSInclude(
				&controllers.WechatController{},
			),
		),
		beego.NSNamespace("/system",
			beego.NSInclude(
				&controllers.SystemController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
