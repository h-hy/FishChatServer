package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/oikomi/FishChatServer/log"
	//	"github.com/oikomi/FishChatServer/log"
)

type DeviceMessageCenter struct {
	Id           int    `orm:"pk;column(id)"`
	IMEI         string `orm:"column(IMEI)"`
	Action       string
	commandId    int
	UserId       int
	FlagResponse int
	Direction    int
	Status       int
	Content      string
	VoiceUri     string
	CreatedAt    string
}

func init() {
	orm.RegisterModel(new(DeviceMessageCenter))
}

func NewDeviceMessageCenter(userId int, IMEI, action, content, voiceUri string, direction int) (int, error) {

	o := orm.NewOrm()

	var dmc DeviceMessageCenter
	log.Info(userId)
	dmc.UserId = userId
	dmc.IMEI = IMEI
	dmc.Action = action
	dmc.Content = content
	dmc.VoiceUri = voiceUri
	dmc.Direction = direction
	dmc.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	id, err := o.Insert(&dmc)
	if err == nil {
		return int(id), nil
	}
	log.Info(err)
	return 0, errors.New("Unknow Error")
}
