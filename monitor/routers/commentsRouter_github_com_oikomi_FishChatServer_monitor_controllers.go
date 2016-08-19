package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:DeviceController"] = append(beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:DeviceController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:DeviceController"] = append(beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:DeviceController"],
		beego.ControllerComments{
			"Show",
			`/:IMEI`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:DeviceController"] = append(beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:DeviceController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:DeviceController"] = append(beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:DeviceController"],
		beego.ControllerComments{
			"Delete",
			`/:IMEI`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:DeviceController"] = append(beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:DeviceController"],
		beego.ControllerComments{
			"Put",
			`/:IMEI`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:DeviceController"] = append(beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:DeviceController"],
		beego.ControllerComments{
			"PostActionLocation",
			`/:IMEI/action/location`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:DeviceController"] = append(beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:DeviceController"],
		beego.ControllerComments{
			"PostActionShutdown",
			`/:IMEI/action/shutdown`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:DeviceController"] = append(beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:DeviceController"],
		beego.ControllerComments{
			"PostVoice",
			`/:IMEI/voice`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:DeviceController"] = append(beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:DeviceController"],
		beego.ControllerComments{
			"GetChatRecord",
			`/:IMEI/chatRecord`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:SystemController"] = append(beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:SystemController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:UserController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:UserController"],
		beego.ControllerComments{
			"ResetPassword",
			`/:username/resetPassword`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:UserController"],
		beego.ControllerComments{
			"UpdateRongCloudToken",
			`/:username/updateRongCloudToken`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:UserController"],
		beego.ControllerComments{
			"Get",
			`/:username`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:UserController"],
		beego.ControllerComments{
			"Put",
			`/:username`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:UserController"],
		beego.ControllerComments{
			"Login",
			`/:username/login`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:UserController"],
		beego.ControllerComments{
			"Logout",
			`/:username/logout`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:UserController"],
		beego.ControllerComments{
			"GetSMSCode",
			`/:username/SMSCode`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:WechatController"] = append(beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:WechatController"],
		beego.ControllerComments{
			"Authorize",
			`/authorize`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:WechatController"] = append(beego.GlobalControllerRouter["github.com/oikomi/FishChatServer/monitor/controllers:WechatController"],
		beego.ControllerComments{
			"AuthorizeRedirect",
			`/authorizeRedirect`,
			[]string{"get"},
			nil})

}
