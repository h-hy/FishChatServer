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
	"flag"
	// "bytes"

	"github.com/astaxie/beego/orm"
	"github.com/oikomi/FishChatServer/base"
	"github.com/oikomi/FishChatServer/models"
	// "github.com/oikomi/FishChatServer/common"
	"github.com/oikomi/FishChatServer/libnet"
	"github.com/oikomi/FishChatServer/log"
	"github.com/oikomi/FishChatServer/protocol"
	"github.com/oikomi/FishChatServer/storage/redis_store"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

type ProtoProc struct {
	pushServer *PushServer
}

func NewProtoProc(pushServer *PushServer) *ProtoProc {
	return &ProtoProc{
		pushServer: pushServer,
	}
}

func (self *ProtoProc) procSubscribeChannel(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("procSubscribeChannel")
	channelName := cmd.GetArgs()[0]
	cUUID := cmd.GetArgs()[1]
	log.Info(channelName)
	if self.pushServer.channels[channelName] != nil {
		if channelName == protocol.SYSCTRL_CONNECT_SERVER {
			self.pushServer.sessions[cUUID] = session
			self.pushServer.sessions[cUUID].State = base.NewSessionState(true, cUUID, "connect_server")
		} else if channelName == protocol.SYSCTRL_API_SERVER {
			self.pushServer.sessions[cUUID] = session
			self.pushServer.sessions[cUUID].State = base.NewSessionState(true, cUUID, "api_server")

		}
		resp := protocol.NewCmdSimple(protocol.SUBSCRIBE_CHANNEL_CMD_ACK)
		session.Send(libnet.Json(resp))
		self.pushServer.channels[channelName].Channel.Join(session, nil)
		self.pushServer.channels[channelName].ClientIDlist = append(self.pushServer.channels[channelName].ClientIDlist, cUUID)
	} else {
		log.Warning(channelName + " is not exist")
	}
	return nil
	// log.Info(self.pushServer.channels)
}

func (self *ProtoProc) procPing(cmd protocol.Cmd, session *libnet.Session) error {
	//log.Info("procPing")
	//cid := session.State.(*base.SessionState).ClientID
	//self.pushServer.sessions[cid].State.(*base.SessionState).Alive = true
	if session.State != nil {
		// fmt.Printf("session.State != nil")
		self.pushServer.scanSessionMutex.Lock()
		defer self.pushServer.scanSessionMutex.Unlock()
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
	sessionCacheData, err := self.pushServer.sessionCache.Get(IMEI)
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
			sessionCacheData := redis_store.NewSessionCacheData(IMEI, session.Conn().RemoteAddr().String(), self.pushServer.cfg.LocalIP, UUID, energy, work_model, volume)
			return sessionCacheData, nil
		}
	}
	return sessionCacheData, nil
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

func (self *ProtoProc) procTransferToDevice(cmd protocol.Cmd, session *libnet.Session) error {
	log.Info("prochTransferToDevice")
	cmdName := cmd.GetInfos()["cmdName"]
	cmd.ChangeCmdName(cmdName)
	ConnectServerUUID := cmd.GetInfos()["ConnectServerUUID"]
	if self.pushServer.sessions[ConnectServerUUID] != nil {
		err := self.pushServer.sessions[ConnectServerUUID].Send(libnet.Json(cmd))
		if err != nil {
			log.Error(err.Error())
		}
		log.Info("send")
	}
	return nil
}
