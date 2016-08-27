package controllers

import "github.com/astaxie/beego"

type ErrorController struct {
	beego.Controller
}

func (self *ErrorController) Prepare() {
	self.Ctx.ResponseWriter.Header().Add("Access-Control-Allow-Origin", "*")
}

func (self *ErrorController) Error404() {
	self.Data["json"] = restReturn(0, "该API不存在或URL参数错误，请检查", map[string]interface{}{})
	self.ServeJSON()
}
