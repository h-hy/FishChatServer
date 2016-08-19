package models

import (
	"github.com/astaxie/beego/orm"
	//	"github.com/oikomi/FishChatServer/log"
)

type UserDevice struct {
	Id     int     `orm:"pk;column(id)"`
	User   *User   `orm:"rel(one);column(user_id)"`
	Device *Device `orm:"rel(one);column(device_IMEI)"`
}

type UserDeviceRel struct {
	Id          int
	UserId      int
	Device_IMEI string
	Device      Device
}

func (u *UserDevice) TableName() string {
	return "user_devices"
}

func init() {
	orm.RegisterModel(new(UserDevice))
}
