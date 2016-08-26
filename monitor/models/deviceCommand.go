package models

import (
	"errors"
	"strconv"

	"github.com/astaxie/beego/orm"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/protocol"
	//	"github.com/oikomi/FishChatServer/log"
)

type DeviceCommand struct {
	Id           int    `orm:"pk;column(id)"`
	IMEI         string `orm:"column(IMEI)"`
	Action       string
	UserId       int
	SendTimes    int
	FlagResponse int
	ActionId     int
	Command      string
}

func init() {
	orm.RegisterModel(new(DeviceCommand))
}

func NewDeviceCommand(userId int, IMEI, action, o_actionId string, command string) (int, int, error) {

	actionId, err := strconv.Atoi(o_actionId)
	if err != nil {
		return 0, 0, err
	}
	o := orm.NewOrm()
	var cmd DeviceCommand
	var DMCAction string
	if o_actionId == protocol.DEIVCE_VOICE_DOWN_CMD {

		DMCAction = "voice"
		err = orm.ErrNoRows
	} else {
		cmd = DeviceCommand{IMEI: IMEI, SendTimes: 0, FlagResponse: 0, ActionId: actionId}
		err = o.Read(&cmd, "IMEI", "sendTimes", "flagResponse", "actionId")

		DMCAction = action
		log.Info(err)
	}
	if err == nil {
		cmd.Command = command
		o.Update(&cmd, "command")
		return 0, int(cmd.Id), nil
	} else if err == orm.ErrNoRows {
		var cmd DeviceCommand
		log.Info(userId)
		cmd.UserId = userId
		cmd.IMEI = IMEI
		cmd.Command = command
		cmd.Action = action
		cmd.ActionId = actionId
		cmd.SendTimes = 0
		cmd.FlagResponse = 0

		Id, err := o.Insert(&cmd)
		DMCId, err := NewDeviceMessageCenter(userId, IMEI, DMCAction, "发送中...", command, 2)
		if err == nil {
			return DMCId, int(Id), nil
		}
		log.Info(err)
	}
	return 0, 0, errors.New("Unknow Error")
}
