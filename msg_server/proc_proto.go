//
// Copyright 2014 Hong Miao (miaohong@miaohong.org). All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strconv"
	"time"
	// "bytes"

	"github.com/astaxie/beego/orm"
	"github.com/oikomi/FishChatServer/base"
	"github.com/oikomi/FishChatServer/models"
	// "github.com/oikomi/FishChatServer/common"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/protocol"
	"github.com/oikomi/FishChatServer/provider"
	"github.com/oikomi/FishChatServer/storage/redis_store"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

type ProtoProc struct {
	msgServer *MsgServer
}

func NewProtoProc(msgServer *MsgServer) *ProtoProc {
	return &ProtoProc{
		msgServer: msgServer,
	}
}

func (self *ProtoProc) procSubscribeChannel(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("procSubscribeChannel")
	channelName := cmd.GetArgs()[0]
	cUUID := cmd.GetArgs()[1]
	log.Info(channelName)
	if self.msgServer.channels[channelName] != nil {
		if channelName == protocol.SYSCTRL_CONNECT_SERVER {
			self.msgServer.sessions[cUUID] = session
			self.msgServer.sessions[cUUID].State = base.NewSessionState(true, cUUID, "connect_server")
		} else if channelName == protocol.SYSCTRL_API_SERVER {
			self.msgServer.sessions[cUUID] = session
			self.msgServer.sessions[cUUID].State = base.NewSessionState(true, cUUID, "api_server")

		} else if channelName == protocol.SYSCTRL_PUSH_SERVER {
			self.msgServer.sessions[cUUID] = session
			self.msgServer.sessions[cUUID].State = base.NewSessionState(true, cUUID, "push_server")
		}
		resp := protocol.NewCmdSimple(protocol.SUBSCRIBE_CHANNEL_CMD_ACK)
		session.Send(libnet.Json(resp))
		self.msgServer.channels[channelName].Channel.Join(session, nil)
		self.msgServer.channels[channelName].ClientIDlist = append(self.msgServer.channels[channelName].ClientIDlist, cUUID)
	} else {
		log.Warning(channelName + " is not exist")
	}
	return nil
	// log.Info(self.msgServer.channels)
}

func (self *ProtoProc) procPing(cmd protocol.Cmd, session *libnet.Session) error {
	//log.Info("procPing")
	//cid := session.State.(*base.SessionState).ClientID
	//self.msgServer.sessions[cid].State.(*base.SessionState).Alive = true
	if session.State != nil {
		// fmt.Printf("session.State != nil")
		self.msgServer.scanSessionMutex.Lock()
		defer self.msgServer.scanSessionMutex.Unlock()
		resp := protocol.NewCmdSimple(protocol.PING_CMD_ACK)
		// log.Info(resp)
		// fmt.Printf("resp=",resp)
		err := session.Send(libnet.Json(resp))
		if err != nil {
			log.Error(err.Error())
			return err
		}
		session.State.(*base.SessionState).Alive = true
	}
	return nil
}

func (self *ProtoProc) procGoOffLine(cmd protocol.Cmd, session *libnet.Session) error {
	IMEI := cmd.GetInfos()["IMEI"]

	sessionCacheData, err := self.msgServer.sessionCache.Get(IMEI)
	if err == nil {
		sessionCacheData.MsgServerAddr = ""
		sessionCacheData.ClientAddr = ""
		sessionCacheData.ConnectServerUUID = ""
		if sessionCacheData.Alive != false {
			o := orm.NewOrm()
			device := models.Device{IMEI: IMEI}
			device.Alive = 0
			o.Update(&device, "Alive")
		}
		sessionCacheData.Alive = false
		sessionCacheData.MaxAge = 600 * time.Second
		self.msgServer.sessionCache.Set(sessionCacheData)
	} else {
		o := orm.NewOrm()
		device := models.Device{IMEI: IMEI}
		device.Alive = 0
		o.Update(&device, "Alive")
	}

	return nil
}

func (self *ProtoProc) procSelectMsgServer(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("procSelectMsgServer")

	IMEI := cmd.GetInfos()["IMEI"]

	session.State.(*base.SessionState).Devices[IMEI] = self.msgServer.cfg.LocalIP

	sessionCacheData, _ := self.checkCache(IMEI, session)
	if sessionCacheData != nil {
		sessionCacheData.MsgServerAddr = self.msgServer.cfg.LocalIP
		sessionCacheData.ConnectServerUUID = session.State.(*base.SessionState).ClientID
		self.msgServer.sessionCache.Set(sessionCacheData)
	}
	return nil
}
func (self *ProtoProc) closeSession(IMEI string, session *libnet.Session) error {
	resp := protocol.NewCmdSimple(protocol.ACTION_DO_CLOSE_SESSION_CMD)
	resp.Infos["IMEI"] = IMEI
	log.Info("Resp | ", resp)

	if session != nil {
		err := session.Send(libnet.Json(resp))
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}
func (self *ProtoProc) checkCache(IMEI string, session *libnet.Session) (*redis_store.SessionCacheData, error) {
	sessionCacheData, err := self.msgServer.sessionCache.Get(IMEI)
	if sessionCacheData == nil {
		log.Warningf("no cache IMEI : %s, err: %s", IMEI, err.Error())
		o := orm.NewOrm()
		device := models.Device{IMEI: IMEI}

		err := o.Read(&device, "IMEI")
		if err == orm.ErrNoRows {
			// not registered
			self.closeSession(IMEI, session)
			log.Warningf("no store IMEI : %s, err: %s", IMEI, err.Error())
			return nil, err
		} else if err == nil {
			UUID := session.State.(*base.SessionState).ClientID
			energy := device.Energy
			work_model := device.Work_model
			volume := device.Volume
			sessionCacheData := redis_store.NewSessionCacheData(IMEI, session.Conn().RemoteAddr().String(), self.msgServer.cfg.LocalIP, UUID, energy, work_model, volume)
			return sessionCacheData, nil
		}
	}
	return sessionCacheData, nil
}
func (self *ProtoProc) procHeartbeat(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("prochHeartbeat")
	resp := protocol.NewCmdSimple("C" + cmd.GetCmdName()[1:])
	IMEI := cmd.GetInfos()["IMEI"]
	if len(cmd.GetArgs()) == 2 {
		energy, err := strconv.Atoi(cmd.GetArgs()[1])
		if err != nil {
			resp.AddArg("2")
		} else {
			resp.AddArg("1")
			//从缓存读取数据
			sessionCacheData, err := self.checkCache(IMEI, session)
			if err != nil {
				return err
			}
			//判断是否需要写缓存
			if sessionCacheData.Energy != energy || sessionCacheData.Alive != true {
				sessionCacheData.Energy = energy
				sessionCacheData.Alive = true
				self.msgServer.sessionCache.Set(sessionCacheData)

				o := orm.NewOrm()
				device := models.Device{IMEI: IMEI}
				device.Energy = energy
				device.Alive = 1
				o.Update(&device, "Energy", "Alive")
			}
		}
	} else {
		resp.AddArg("2")
	}

	resp.Infos["IMEI"] = IMEI

	if session != nil {
		err := session.Send(libnet.Json(resp))
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}

func (self *ProtoProc) procTimeSync(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("prochTimeSync")
	resp := protocol.NewCmdSimple("C" + cmd.GetCmdName()[1:])
	IMEI := cmd.GetInfos()["IMEI"]
	resp.AddArg("1")
	resp.AddArg(time.Now().Format("2006-01-02 15:04:05"))

	resp.Infos["IMEI"] = IMEI

	if session != nil {
		err := session.Send(libnet.Json(resp))
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}

func (self *ProtoProc) procLocation(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("prochLocation")
	resp := protocol.NewCmdSimple("C" + cmd.GetCmdName()[1:])
	IMEI := cmd.GetInfos()["IMEI"]
	log.Info("len(cmd.GetArgs())=", len(cmd.GetArgs()))
	if len(cmd.GetArgs()) == 1 {
		var locationInfo provider.Location
		locationInfo.Parse(cmd.GetArgs()[0])
		log.Info("locationInfo=", locationInfo)
		log.Info("locationInfo.LocationData=", locationInfo.LocationData)

		sessionCacheData, err := self.checkCache(IMEI, session)
		if err != nil {
			return err
		}
		//判断是否需要更新缓存
		if sessionCacheData.Location != locationInfo.LocationData {
			//上报数据与缓存不一致
			sessionCacheData.Location = locationInfo.LocationData
			sessionCacheData.Alive = true
			self.msgServer.sessionCache.Set(sessionCacheData)
			//更新缓存完毕，开始增加历史位置记录
			models.NewDeviceLocationHistory(IMEI, locationInfo)
			//开始更新设备数据表
			location, err := json.Marshal(locationInfo.LocationData)
			log.Info(string(location))
			if err == nil {
				o := orm.NewOrm()
				device := models.Device{IMEI: IMEI}
				device.Location = string(location)
				o.Update(&device, "location")
			}
		}
		resp.AddArg("1")
	} else {
		resp.AddArg("2")
	}

	resp.Infos["IMEI"] = IMEI

	if session != nil {
		err := session.Send(libnet.Json(resp))
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}
func (self *ProtoProc) procLinkDesc(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("prochLinkDesc")
	resp := protocol.NewCmdSimple("C" + cmd.GetCmdName()[1:])
	IMEI := cmd.GetInfos()["IMEI"]
	if len(cmd.GetArgs()) == 2 {
		commId := cmd.GetArgs()[0]
		resp.AddArg(commId)
		resp.AddArg("1")
	} else if len(cmd.GetArgs()) == 1 {
		commId := cmd.GetArgs()[0]
		resp.AddArg(commId)
		resp.AddArg("2")
	} else {
		resp.AddArg("2")
	}

	resp.Infos["IMEI"] = IMEI

	if session != nil {
		err := session.Send(libnet.Json(resp))
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}

func (self *ProtoProc) procVoiceReaded(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("prochVoiceReaded")
	resp := protocol.NewCmdSimple("C" + cmd.GetCmdName()[1:])
	IMEI := cmd.GetInfos()["IMEI"]
	if len(cmd.GetArgs()) == 1 {
		id := cmd.GetArgs()[0]
		resp.AddArg(id)
		resp.AddArg("1")
	} else {
		resp.AddArg("2")
	}

	resp.Infos["IMEI"] = IMEI

	if session != nil {
		err := session.Send(libnet.Json(resp))
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}

func (self *ProtoProc) procLowPower(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("prochLowPower")
	resp := protocol.NewCmdSimple("C" + cmd.GetCmdName()[1:])
	IMEI := cmd.GetInfos()["IMEI"]
	resp.AddArg("1")
	resp.Infos["IMEI"] = IMEI

	if session != nil {
		err := session.Send(libnet.Json(resp))
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}

func (self *ProtoProc) procSOS(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("prochSOS")
	resp := protocol.NewCmdSimple("C" + cmd.GetCmdName()[1:])
	IMEI := cmd.GetInfos()["IMEI"]
	resp.AddArg("1")
	resp.Infos["IMEI"] = IMEI

	if session != nil {
		err := session.Send(libnet.Json(resp))
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}

func (self *ProtoProc) procTransferToDevice(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("prochTransferToDevice")
	cmdName := cmd.GetInfos()["cmdName"]
	cmd.ChangeCmdName(cmdName)
	ConnectServerUUID := cmd.GetInfos()["ConnectServerUUID"]
	if self.msgServer.sessions[ConnectServerUUID] != nil {
		err := self.msgServer.sessions[ConnectServerUUID].Send(libnet.Json(cmd))
		if err != nil {
			log.Error(err.Error())
		}
		log.Info("send")
	}
	return nil
}

func (self *ProtoProc) procupdateSetting(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("prochupdateSetting")
	resp := protocol.NewCmdSimple("C" + cmd.GetCmdName()[1:])
	IMEI := cmd.GetInfos()["IMEI"]
	resp.AddArg("1")
	resp.Infos["IMEI"] = IMEI

	sessionCacheData, err := self.checkCache(IMEI, session)
	if err == nil {
		fmt.Print("begin")
		{
			send := protocol.NewCmdSimple("D2")
			send.AddArg(strconv.Itoa(sessionCacheData.WorkModel))
			send.Infos["IMEI"] = IMEI
			if session != nil {
				session.Send(libnet.Json(send))
			}
		}
		{
			send := protocol.NewCmdSimple("D29")
			send.AddArg(strconv.Itoa(sessionCacheData.Volume))
			send.Infos["IMEI"] = IMEI
			if session != nil {
				session.Send(libnet.Json(send))
			}
		}
	}

	if session != nil {
		err := session.Send(libnet.Json(resp))
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}
