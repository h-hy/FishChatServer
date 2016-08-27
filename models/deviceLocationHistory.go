package models

import (
	"encoding/json"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/oikomi/FishChatServer/provider"
	//	"github.com/oikomi/FishChatServer/log"
)

type DeviceLocationHistory struct {
	Id          int    `orm:"pk;column(id)"`
	IMEI        string `orm:"column(IMEI)"`
	Location    string
	Electricity int
	CreatedAt   string
}

func init() {
	orm.RegisterModel(new(DeviceLocationHistory))
}

func NewDeviceLocationHistory(IMEI string, locationObj provider.Location) error {

	o := orm.NewOrm()

	var dmc DeviceLocationHistory
	dmc.IMEI = IMEI
	dmc.Electricity = locationObj.Energy
	dmc.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	location, err := json.Marshal(locationObj.LocationData)
	if err != nil {
		return err
	}
	dmc.Location = string(location)
	_, err = o.Insert(&dmc)
	return err
}
