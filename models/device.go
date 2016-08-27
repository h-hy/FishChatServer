package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/oikomi/FishChatServer/log"
)

type Device struct {
	IMEI           string `orm:"pk;column(IMEI)"`
	Work_model     int
	Volume         int
	Energy         int
	Alive          int
	EmergencyPhone string
	Nick           string
	Location       string
	User           []*User `orm:"reverse(many)"`
}

//func (u *Device) TableName() string {
//	return "devices"
//}
func init() {
	orm.RegisterModel(new(Device))
}
func CheckBind(username, IMEI string) bool {
	cacheExist := redisCache.IsExist("user_" + username)
	user, err := GetUser(username)
	if err != nil {
		return false
	}

	for _, device := range user.Devices {
		if device.IMEI == IMEI {
			return true
		}
	}
	if cacheExist == true {
		//重新获取一次试试
		//		redisCache.Delete("user_" + username)
		user.UpdateDevice()
		log.Info(IMEI)
		for _, device := range user.Devices {
			log.Info(device.IMEI)
			if device.IMEI == IMEI {
				return true
			}
		}
	}
	return false
}
